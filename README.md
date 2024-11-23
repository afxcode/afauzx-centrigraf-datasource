# Centrigraf - Grafana DataSource Plugin for Centrifugo

This Grafana DataSource plugin allows you to connect Grafana to a Centrifugo backend, enabling real-time data streaming from Centrifugo channels into Grafana dashboards. With this plugin, you can display live data streams from your Centrifugo instance in Grafana without needing a backend.

## Key Features:
- **Real-time Data**: Subscribe to Centrifugo channels and stream live data into Grafana.
- **Frontend-Only**: The plugin operates entirely on the frontend, leveraging WebSocket connections for real-time updates.
- **Easy Configuration**: Simple setup via Grafana's native data source configuration UI.

## Installation Instructions:

### 1. Download the Plugin:
- Download the latest `.zip` release from the [GitHub releases page](https://github.com/afxcode/afauzx-centrigraf-datasource/releases).

### 2. Install the Plugin:
- Extract the `.zip` file into the Grafana plugin directory:
   - **Linux**: `/var/lib/grafana/plugins/`
   - **Windows**: `C:\Program Files\GrafanaLabs\grafana\data\plugins\`
   - **Docker**: `/var/lib/grafana/plugins/`

### 3. Restart Grafana:
- After placing the plugin in the correct directory, restart Grafana:
   - **Linux**: `sudo systemctl restart grafana-server`
   - **Windows**: Restart the Grafana service.

### 4. Configure the DataSource:
- Navigate to **Configuration > Data Sources** in Grafana.
- Add a new data source and select **Centrifugo** from the list.
- Provide the URL of your Centrifugo server.

### 5. Test the Connection:
- Use the **Save & Test** button in the data source configuration to verify the connection. The plugin will attempt to connect to the Centrifugo backend and ensure that real-time data streaming is working.

## How It Works:
- The plugin connects to your Centrifugo server over WebSockets to receive live updates from channels.
- Once the connection is established, it subscribes to specified channels and sends the incoming data to Grafana.
- The plugin requires no backend server and runs entirely on the frontend, making it simple to deploy and configure.

## Notes:
- **Centrifugo Instance**: Ensure that your Centrifugo server is running and accessible from your Grafana instance. You’ll need the URL of your Centrifugo instance to configure the plugin.
- **Plugin Status**: This plugin is currently a work in progress and may not be fully tested. Please report any issues or feature requests on the [GitHub issues page](https://github.com/yourusername/your-repo/issues).

## Future Improvements:
- Enhanced error handling and automatic reconnection logic for more robust operation.
- Support for private channels and authentication mechanisms to secure your real-time data streams.
- Expanded configuration options for advanced use cases.

## License:
This plugin is open-source and licensed under the [Apache License](LICENSE).

---

### Contributing:
Feel free to open an issue or submit a pull request if you would like to contribute to the development of this plugin. Your contributions are welcome!

