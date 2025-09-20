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
- **TCP client mode**: Connect to remote servers (new in v0.2.0)
- **TCP Bridge server**: Relay data between external clients and local services
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

## Quick Start Guide

**Choose your scenario:**

| Scenario | Tool | Configuration | Description |
|----------|------|---------------|-------------|
| **macOS with serial device** | `openserial` | `configs/server-mac.yaml` | UART-to-TCP bridge |
| **macOS without serial device** | `tcpbridge` | *(default config)* | Pure TCP relay server |
| **Windows client** | `openserial -client` | `configs/client-windows.yaml` | Connect to Mac server |
| **Development/testing** | `tcpbridge` | `configs/tcpbridge.yaml` | Forward TCP connections |

## Usage

### Basic Usage

```bash
# Run with default configuration (server mode)
./openserial

# Run with custom configuration file
./openserial -config /path/to/config.yaml

# Run in client mode (connect to remote server)
./openserial -client -config /path/to/client-config.yaml

# Show help
./openserial -help

# Show version
./openserial -version
```

### TCP Bridge Mode Usage

The TCP Bridge mode provides a pure TCP-to-TCP relay server that doesn't require serial ports. This is ideal for:

- **macOS without serial hardware**: Run a TCP server without needing physical serial devices
- **Cross-platform access**: When Windows machines don't have admin privileges
- **Development and testing**: Forward connections between different TCP services
- **Load balancing**: Distribute connections to multiple backend servers

#### macOS TCP Server (No Serial Port Required)

#### Option 1: Standalone TCP Bridge Server (Default Configuration)

```bash
# Start TCP bridge server with default configuration (no config file needed)
./tcpbridge

# This will:
# - Listen on port 8080 for incoming connections
# - Forward all data to localhost:8081
# - Handle up to 10 concurrent connections
# - Use default timeout of 5 minutes per connection
```

#### Option 1b: Custom Configuration

```bash
# Start TCP bridge server with custom configuration
./tcpbridge -config configs/tcpbridge.yaml
```

#### Option 2: Combined Setup with Serial Bridge

```bash
# Terminal 1: Start TCP bridge server
./tcpbridge -config configs/tcpbridge.yaml

# Terminal 2: Start serial bridge (if you have serial hardware)
./openserial -config configs/server-mac.yaml
```

#### Windows Client (Connects to Mac Server)

```bash
# Update config with Mac's IP address, then:
./openserial -client -config configs/client-windows.yaml
```

#### Pure TCP-to-TCP Forwarding

For scenarios where you need to forward TCP connections without any serial involvement:

```bash
# Start TCP bridge with default configuration (forwards port 8080 → 8081)
./tcpbridge

# Or with custom configuration
./tcpbridge -config configs/tcpbridge.yaml

# Clients connect to localhost:8080
# Data gets forwarded to localhost:8081
# Perfect for development, testing, or service proxying
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

**TCP Bridge Configuration:**

- `server.port`: Port for external client connections (default: 8080)
- `server.host`: Host to bind to (default: "0.0.0.0")
- `target.host`: Target server hostname (default: "localhost")
- `target.port`: Target server port (default: 8081)
- `clients.max_connections`: Maximum concurrent connections (default: 10)
- `clients.connection_timeout`: Connection timeout (default: "5m")

**Default Configuration:**

Both `openserial` and `tcpbridge` work without configuration files using sensible defaults:

- **openserial**: Requires serial port configuration, will prompt for missing values
- **tcpbridge**: Works immediately with default values (port 8080 → localhost:8081)

### Configuration File Locations

The application looks for configuration files in the following order:

1. File specified with `-config` flag
2. `./config.yaml`
3. `./configs/config.yaml`
4. `/etc/openserial/config.yaml`
5. `$HOME/.openserial/config.yaml`

## Examples

### Serial Device Remote Access

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

Connect to your serial device via TCP:

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

### TCP Bridge Server Configuration

```yaml
server:
  port: 8080
  host: "0.0.0.0"

target:
  host: "localhost"
  port: 8081

clients:
  max_connections: 10
  connection_timeout: "5m"
```

### macOS Without Serial Hardware

For macOS users who don't have serial devices but need to run a TCP server:

**Usage (Default Configuration):**

```bash
# Build the TCP bridge
go build -o build/tcpbridge ./cmd/tcpbridge

# Start the server with default configuration
./build/tcpbridge
```

**Usage (Custom Configuration):**

```bash
# Start with custom configuration
./build/tcpbridge -config configs/tcpbridge.yaml
```

**What this does:**

- Listens on `localhost:8080` for incoming connections
- Forwards all data to `localhost:8081` (your target service)
- Handles up to 10 concurrent connections
- No serial port required - pure TCP relay
- Perfect for development, testing, or service proxying

**Use Cases:**

- Forward test clients to your development server
- Bridge between different TCP services
- Load balance connections to multiple backends
- Network debugging and monitoring

## Building

### Prerequisites

- Go 1.25 or later
- Make (optional, for using Makefile)

### macOS Security Warning

If you get a security warning when opening the macOS binary, this is normal for
unsigned binaries. To resolve:

```bash
# Make the binary executable
chmod +x openserial-darwin-amd64
# or
chmod +x openserial-darwin-arm64

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

### Standard Mode (Direct TCP Bridge)

```text
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Serial Port   │◄──►│   OpenSerial     │◄──►│   TCP Client    │
│   (Device)      │    │   Bridge         │    │   (Network)     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### TCP Bridge Mode (Cross-Platform Access)

```text
Windows Machine (No Admin)          Mac Machine (Admin)
┌─────────────────────────┐        ┌─────────────────────────┐
│  UART Device            │        │  Serial Monitor         │
│  (Serial)               │        │  (Port 8081 TX/RX)      │
└─────────┬───────────────┘        └─────────┬───────────────┘
          │                                  │
          │                                  │
┌─────────▼───────────────┐        ┌─────────▼───────────────┐
│  OpenSerial Client      │◄──────►│  TCP Bridge Server      │
│  (Connects to Mac:8080) │        │  (Listens: 8080)        │
└─────────────────────────┘        └─────────────────────────┘
```

### TCP-Only Mode (No Serial Hardware Required)

```text
Client Application          Mac Machine (No Serial Hardware)
┌─────────────────────────┐        ┌─────────────────────────┐
│  Any TCP Client         │        │  Target Service         │
│  (Port 8080)            │        │  (Port 8081)            │
└─────────┬───────────────┘        └─────────┬───────────────┘
          │                                  │
          │                                  │
┌─────────▼───────────────┐        ┌─────────▼───────────────┐
│  TCP Bridge Server      │◄──────►│  Your Application       │
│  (Pure TCP Relay)       │        │  (Any TCP Service)      │
└─────────────────────────┘        └─────────────────────────┘
```

**Perfect for:**

- macOS without serial devices
- Development and testing
- Service proxying and load balancing
- Network debugging

The bridge consists of:

- **Serial Handler**: Manages serial port communication
- **Network Handler**: Manages TCP server and client connections
- **Bridge Core**: Coordinates bidirectional data forwarding
- **Configuration**: Manages settings and validation
- **TCP Bridge Server**: Relays data between external clients and local services

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

- Serial device remote access
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
