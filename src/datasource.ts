import {
  CoreApp,
  DataQueryRequest,
  DataQueryResponse,
  DataSourceApi,
  DataSourceInstanceSettings,
  CircularDataFrame,
  FieldType,
  LoadingState,
} from '@grafana/data';

import {MyQuery, MyDataSourceOptions} from './types';
import {Observable, merge, Subscriber} from 'rxjs';

import {Centrifuge, Subscription} from 'centrifuge';

type SubscriberWithFrame = {
  subscriber: Subscriber<DataQueryResponse>;
  frame: CircularDataFrame;
};

type ChannelSubscriber = {
  centrifugeSubscription: Subscription;
  subscribers: SubscriberWithFrame[];
}

export class DataSource extends DataSourceApi<MyQuery, MyDataSourceOptions> {
  centrifuge: Centrifuge;
  private channelSubscribers: Map<string, ChannelSubscriber> = new Map();
  errorMessage: string | undefined;

  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
    const settings = instanceSettings.jsonData

    this.centrifuge = new Centrifuge(settings.url);
    this.centrifuge.on("connected", (ctx) => {
      console.log("Centrifuge Connected", ctx)
    })

    this.centrifuge.on("disconnected", (ctx) => {
      console.warn("Centrifuge Disconnected", ctx)
    })

    this.centrifuge.on("error", (ctx) => {
      console.error("Centrifuge Error", ctx.error.message)
      this.errorMessage = ctx.error.message
    })

    this.centrifuge.connect();
  }

  getDefaultQuery(_: CoreApp): Partial<MyQuery> {
    return {channelName: ""};
  }

  filterQuery(query: MyQuery): boolean {
    return !!query.channelName;
  }

  query(options: DataQueryRequest<MyQuery>): Observable<DataQueryResponse> {
    const observables = options.targets.map((query) => {
      return new Observable<DataQueryResponse>((subscriber) => {
        const channelName = query.channelName

        if (!this.channelSubscribers.has(channelName)) {
          this.channelSubscribers.set(channelName, {
            centrifugeSubscription: this.centrifuge.newSubscription(channelName),
            subscribers: [],
          });

          const channel = this.channelSubscribers.get(channelName)!;
          channel.centrifugeSubscription.on('publication', (ctx) => {
            channel.subscribers.forEach(function (sub) {
              sub.frame.add({time: Date.now(), value: ctx.data});
              sub.subscriber.next({
                data: [sub.frame],
                key: sub.frame.refId,
                state: LoadingState.Streaming,
              });
            })
          });

          channel.centrifugeSubscription.subscribe();
        }

        const channel = this.channelSubscribers.get(channelName)!;

        let frame = new CircularDataFrame({
          append: 'tail',
          capacity: 1000,
        });

        frame.refId = query.refId;
        frame.addField({name: 'time', type: FieldType.time});
        frame.addField({name: 'value', type: FieldType.number});

        channel.subscribers.push({subscriber, frame});

        return () => {
          this.cleanupChannel(channelName, subscriber)
        }
      });
    });

    return merge(...observables);
  }

  private cleanupChannel(channelName: string, subscriber: Subscriber<DataQueryResponse>): void {
    const channel = this.channelSubscribers.get(channelName);
    if (!channel) {
      return
    }

    channel.subscribers = channel.subscribers.filter((s) => s.subscriber !== subscriber);

    if (channel.subscribers.length === 0) {
      channel.centrifugeSubscription.unsubscribe();
      this.centrifuge.removeSubscription(channel.centrifugeSubscription)
      this.channelSubscribers.delete(channelName);
    }
  }

  /**
   * Checks whether we can connect to the API.
   */
  async testDatasource() {
    let errorMessage = 'Cannot connect to API';

    try {
      await this.centrifuge.ready(1000);
      return {
        status: 'success',
        message: 'Success',
      };
    } catch (err) {
      if (this.errorMessage !== "") {
        errorMessage += ": " + this.errorMessage
      }
      return {
        status: 'error',
        message: errorMessage,
      };
    }
  }
}
