import React, { ChangeEvent } from 'react';
import { InlineField, Input } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions, MySecureJsonData } from '../types';

interface Props extends DataSourcePluginOptionsEditorProps<MyDataSourceOptions, MySecureJsonData> {}

export function ConfigEditor(props: Props) {
  const { onOptionsChange, options } = props;
  const { jsonData } = options;

  const onUrlChange = (event: ChangeEvent<HTMLInputElement>) => {
    onOptionsChange({
      ...options,
      jsonData: {
        ...jsonData,
        url: event.target.value,
      },
    });
  };

  return (
    <>
      <InlineField label="URL" labelWidth={14} interactive tooltip={'Full path of Centrifuge Websocket url'}>
        <Input
          id="config-editor-url"
          onChange={onUrlChange}
          value={jsonData.url}
          placeholder="Enter the url, e.g. 'ws://localhost:8000/api-ws"
          width={50}
        />
      </InlineField>
    </>
  );
}
