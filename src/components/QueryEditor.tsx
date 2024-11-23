import React, { ChangeEvent } from 'react';
import { InlineField, Input, Stack } from '@grafana/ui';
import { QueryEditorProps } from '@grafana/data';
import { DataSource } from '../datasource';
import { MyDataSourceOptions, MyQuery } from '../types';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

export function QueryEditor({ query, onChange, onRunQuery }: Props) {
  const onChannelNameChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({ ...query, channelName: event.target.value });
    onRunQuery();
  };

  const { channelName } = query;

  return (
    <Stack gap={0}>
      <InlineField label="Channel Name">
        <Input
          id="query-editor-channel-name"
          onChange={onChannelNameChange}
          value={channelName}
          placeholder="Enter a channel name"
        />
      </InlineField>
    </Stack>
  );
}
