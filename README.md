# Xray Exporter

Xray Exporter is a tool for exporting metrics from Xray to Prometheus. It collects user statistics (downlinks and uplinks) from the Xray API and provides them in a Prometheus-compatible format.

## Description

The exporter periodically queries user statistics from the Xray API endpoint `/debug/vars` and exports metrics to Prometheus. Supported metrics:

- `xray_user_downlinks`: The number of active user downlinks (by user_id)
- `xray_user_uplinks`: The number of active user uplinks (by user_id)

## Installation

### Building from source

1. Ensure you have Go version 1.25 or higher installed.
2. Clone the repository or download the source code.
3. Navigate to the project directory:
   ```
   cd /path/to/xray-exporter
   ```
4. Build the binary:
   ```
   go build -o xray-exporter .
   ```

### Using Docker

1. Build the Docker image:
   ```
   docker build -t xray-exporter .
   ```
2. Run the container:
   ```
   docker run -p 9595:9595 xray-exporter
   ```

## Running

Run the exporter with parameters:

```
./xray-exporter -xray-endpoint=localhost:11111 -port=9595
```

### Command-line parameters

- `-xray-endpoint`: The address of the Xray API endpoint (default: `localhost:11111`)
- `-port`: The port for the Prometheus metrics HTTP server (default: `9595`)

## Prometheus Configuration

Add to your Prometheus configuration:

```yaml
scrape_configs:
  - job_name: 'xray'
    static_configs:
      - targets: ['localhost:9595']
```

## Metrics

### xray_user_downlinks

- **Type**: Gauge
- **Description**: The number of active downlinks for the user
- **Labels**: `user_id`

### xray_user_uplinks

- **Type**: Gauge
- **Description**: The number of active uplinks for the user
- **Labels**: `user_id`

## Dependencies

- Go 1.25+
- Prometheus client library for Go

## License

MIT License