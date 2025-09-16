# OpenSerial

OpenSerial is a UART-to-TCP bridge agent that enables network access to serial devices.

<!-- Workflow testing branch -->

## Overview

OpenSerial creates a bidirectional bridge between serial ports and TCP network
connections, allowing remote access to serial devices over the network. It's
designed as a lightweight, cross-platform solution for IoT, embedded systems,
and industrial applications.

## Features

- **Bidirectional forwarding**: Serial port ↔ TCP port (TX/RX)
- **Full serial configuration**: Baud rate, data bits, stop bits, parity, flow control
- **TCP server mode**: Listen for incoming connections on specified ports
- **Network interface binding**: Bind to all network interfaces (0.0.0.0)
- **Configuration files**: Store settings in YAML format
- **Cross-platform**: Windows, macOS, and Linux support
- **Automatic reconnection**: Handles serial port and network disconnections

## Installation

### Binary Releases

Download the latest release for your platform from the [releases page](https://github.com/Prototype-Cafe-LLC/OpenSerial/releases).

### Build from Source

1. Install Go 1.25+
2. Clone the repository:

   ```bash
   git clone https://github.com/Prototype-Cafe-LLC/OpenSerial
   cd openserial
   ```

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Build:

   ```bash
   make build
   ```

## Usage

### Basic Usage

```bash
# Run with default configuration
./openserial

# Run with custom configuration file
./openserial -config /path/to/config.yaml

# Show help
./openserial -help

# Show version
./openserial -version
```

### Configuration

Create a configuration file (YAML format):

```yaml
serial:
  port: "/dev/ttyUSB0"
  baud_rate: 115200
  data_bits: 8
  stop_bits: 1
  parity: "none"
  flow_control: "none"

network:
  listen_port: 8080
  bind_address: "0.0.0.0"
```

#### Configuration Options

**Serial Configuration:**

- `port`: Serial port device path (e.g., "/dev/ttyUSB0", "COM3")
- `baud_rate`: Baud rate (300, 600, 1200, 2400, 4800, 9600, 19200, 38400,
  57600, 115200, 230400, 460800, 921600)
- `data_bits`: Data bits (5-8)
- `stop_bits`: Stop bits (1-2)
- `parity`: Parity (none, odd, even, mark, space)
- `flow_control`: Flow control (none, rts_cts, xon_xoff)

**Network Configuration:**

- `listen_port`: TCP port to listen on (1-65535)
- `bind_address`: IP address to bind to (e.g., "0.0.0.0" for all interfaces)

### Configuration File Locations

The application looks for configuration files in the following order:

1. File specified with `-config` flag
2. `./config.yaml`
3. `./configs/config.yaml`
4. `/etc/openserial/config.yaml`
5. `$HOME/.openserial/config.yaml`

## Examples

### Arduino/ESP32 Remote Access

```yaml
serial:
  port: "/dev/ttyUSB0"
  baud_rate: 115200
  data_bits: 8
  stop_bits: 1
  parity: "none"
  flow_control: "none"

network:
  listen_port: 8080
  bind_address: "0.0.0.0"
```

Connect to your Arduino via TCP:

```bash
telnet localhost 8080
```

### Industrial Equipment Monitoring

```yaml
serial:
  port: "/dev/ttyS0"
  baud_rate: 9600
  data_bits: 8
  stop_bits: 1
  parity: "even"
  flow_control: "rts_cts"

network:
  listen_port: 502
  bind_address: "192.168.1.100"
```

## Building

### Prerequisites

- Go 1.25 or later
- Make (optional, for using Makefile)

### macOS Security Warning

If you get a security warning when opening the macOS binary, this is normal for
unsigned binaries. To resolve:

```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine openserial-darwin-amd64
# or
xattr -d com.apple.quarantine openserial-darwin-arm64
```

Alternatively, right-click the binary and select "Open" from the context menu.

### Build Commands

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Lint code
make lint

# Clean build artifacts
make clean
```

### Cross-Platform Building

The Makefile supports building for multiple platforms:

- **Linux**: amd64, arm64
- **Windows**: amd64, arm64
- **macOS**: amd64, arm64

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/config
```

## Architecture

```text
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Serial Port   │◄──►│   OpenSerial     │◄──►│   TCP Client    │
│   (Device)      │    │   Bridge         │    │   (Network)     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

The bridge consists of:

- **Serial Handler**: Manages serial port communication
- **Network Handler**: Manages TCP server and client connections
- **Bridge Core**: Coordinates bidirectional data forwarding
- **Configuration**: Manages settings and validation

## Error Handling

OpenSerial includes comprehensive error handling:

- **Serial Port Errors**: Automatic reconnection on device disconnection
- **Network Errors**: Client reconnection handling
- **Buffer Management**: Overflow protection and data integrity
- **Configuration Validation**: Validates all configuration parameters

## Performance

- **Latency**: < 10ms for local network connections
- **Throughput**: Support up to 1Mbps serial data rates
- **Memory Usage**: < 50MB per active bridge

## Use Cases

### Industrial Applications

- Remote monitoring of industrial equipment
- SCADA system integration
- PLC communication bridges

### IoT Development

- Remote debugging of embedded devices
- IoT device management
- Sensor data collection

### Educational & Hobby

- Arduino/Raspberry Pi remote access
- Electronics project development
- Learning serial communication protocols

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

### Development Setup

```bash
git clone https://github.com/Prototype-Cafe-LLC/OpenSerial
cd openserial
go mod download
make build
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

For issues and questions:

- Create an issue on GitHub
- Check the documentation
- Review the test cases for usage examples
