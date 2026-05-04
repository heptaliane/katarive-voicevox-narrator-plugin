# katarive-voicevox-narrator-plugin

A plugin for [katarive](https://github.com/heptaliane/katarive-server) that enables narration using the [VOICEVOX](https://voicevox.hiroshiba.jp/) engine.

## Prerequisites

- **Go**: 1.26 or later
- **Docker**: For running the VOICEVOX engine
- **ffmpeg**: For audio processing and conversion

## Setup

### 1. Start VOICEVOX Engine

This plugin requires a running VOICEVOX engine. You can start it using Docker Compose:

```bash
make run-voicevox
```

By default, this uses the CPU version of the VOICEVOX engine.

### 2. Build the Plugin

Build the plugin binary:

```bash
make build
```

This will generate a binary named `katarive-voicevox-narrator-plugin`.

## Configuration

You can configure the plugin using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `KATARIVE_VOICEVOX_SERVER` | URL of the VOICEVOX engine server | `http://localhost:50021` |

## Usage

This is a HashiCorp plugin. It is designed to be used as a backend for `katarive`.

```bash
# Example usage (assuming katarive is installed and configured)
./katarive-voicevox-narrator-plugin
```

## Development

### Code Generation

To regenerate the VOICEVOX client from the OpenAPI specification:

```bash
make generate
```

### Testing

Run the tests (requires VOICEVOX engine to be running):

```bash
make test
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
