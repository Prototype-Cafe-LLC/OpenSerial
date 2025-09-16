# OpenSerial

OpenSerial is a UART-to-TCP bridge agent that enables network access to serial devices.

## Overview

OpenSerial creates a bidirectional bridge between serial ports and TCP network connections, allowing remote access to serial devices over the network. It's designed as a lightweight, cross-platform solution for IoT, embedded systems, and industrial applications.

## Core Features

### Serial Communication

- **Bidirectional forwarding**: Serial port ↔ TCP port (TX/RX)
- **Full serial configuration**: Baud rate, data bits, stop bits, parity, flow control
- NOT in the scope: **Multiple serial ports**: Support for multiple concurrent serial port bridges
- NOT in the scope: **Hot-plug support**: Dynamic serial port detection and connection

### Network Features

- **TCP server mode**: Listen for incoming connections on specified ports
- NOT in the scope: **Multiple clients**: Support multiple concurrent TCP connections per serial port
- **Network interface binding**: Bind to all network interfaces (0.0.0.0)
- NOT in the scope: **Network interface binding**: Bind to specific interfaces
- **Port specification**: Configurable local port/listening port
- **Connection management**: Automatic client connection handling and cleanup

### Configuration Management

- **Configuration files**: Store settings in JSON/YAML format
- NOT in the scope: **Command-line options**: Override configuration via CLI arguments
- NOT in the scope: **Environment variables**: Support for environment-based configuration
- NOT in the scope: **Runtime reconfiguration**: Reload configuration without restart

## Technical Specifications

### Platform Support

- **Windows**: x64, ARM64
- **macOS**: Intel (x64), Apple Silicon (ARM64)
- **Linux**: Intel/AMD64, ARM64

### Architecture

- **Language**: Go (Golang)
- **Dependencies**:
  - `go-serial` for cross-platform serial communication
  - `net` package for TCP networking
  - `viper` for configuration management
- **Build**: Single binary, no external dependencies
- **Distribution**: Standalone executable

### Performance Requirements

- **Latency**: < 10ms for local network connections
- **Throughput**: Support up to 1Mbps serial data rates
- NOT in the scope: **Concurrent connections**: Up to 100 TCP clients per serial port
- **Memory usage**: < 50MB per active bridge

## Security Features

### Access Control

- NOT in the scope: **IP whitelist**: Restrict connections to specific IP addresses
- NOT in the scope: **Port authentication**: Optional password-based authentication
- NOT in the scope: **Connection limits**: Maximum concurrent connections per serial port

### Network Security

- NOT in the scope: **TLS support**: Optional encrypted connections (TLS 1.2+)
- NOT in the scope: **Firewall friendly**: Configurable port ranges and protocols

## Error Handling & Monitoring

### Error Recovery

- **Serial port reconnection**: Automatic reconnection on device disconnection
- **Network reconnection**: Client reconnection handling
- **Buffer management**: Overflow protection and data integrity

### Logging & Monitoring

- NOT in the scope: **Structured logging**: JSON-formatted logs with configurable levels
- **Connection status**: Real-time connection monitoring
- NOT in the scope: **Performance metrics**: Throughput, latency, and error statistics
- NOT in the scope: **Health checks**: Built-in health monitoring endpoints

## Configuration Examples

### Basic Configuration

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
  #max_connections: 10 # NOT in the scope

# NOT in the scope
# logging:
#   level: "info"
#   file: "/var/log/openserial.log"
```

### Command Line Usage

```bash
# Basic usage
# NOT in the scope
# openserial --port /dev/ttyUSB0 --baud 115200 --listen 8080

# With configuration file
openserial --config /etc/openserial/config.yaml

# Multiple serial ports
# NOT in the scope
# openserial --config /etc/openserial/multi-port.yaml
```

## Installation & Distribution

### Installation Methods

- **Binary releases**: Pre-compiled executables for all platforms
- NOT in the scope: **Package managers**: Homebrew (macOS), apt/yum (Linux), Chocolatey (Windows)
- NOT in the scope: **Docker**: Containerized deployment option
- **Source build**: Go module-based build system

### System Requirements

- **Minimum**: 512MB RAM, 50MB disk space
- **Recommended**: 1GB RAM, 100MB disk space
- **OS**: Windows 10+, macOS 10.15+, Linux (kernel 3.10+)

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

### Development Setup

1. Install Go 1.25+
2. Clone repository: `git clone https://github.com/yourorg/openserial`
3. Install dependencies: `go mod download`
4. Build: `go build -o openserial`
5. Test: `go test ./...`

### Code Standards

- Follow Go standard formatting (`gofmt`)
- Include comprehensive tests
- Document all public APIs
- Follow semantic versioning

## Testing

### Test Plan

A comprehensive test plan has been developed to ensure OpenSerial meets all functional, performance, and quality requirements. The test plan covers:

- **Unit Tests**: Configuration management, serial port handling, network management, and data forwarding
- **Integration Tests**: End-to-end bridge functionality and configuration integration
- **Performance Tests**: Latency (< 10ms), throughput (up to 1Mbps), and memory usage (< 50MB)
- **Cross-Platform Tests**: Windows, macOS, and Linux compatibility validation
- **Error Handling Tests**: Serial port errors, network errors, and application error recovery
- **Security Tests**: Network security and input validation
- **Stress Tests**: Long-running and high-load scenarios

The complete test plan is documented in [TEST_PLAN.md](TEST_PLAN.md) and includes 113 specific test cases with unique identifiers, a 7-week phased execution strategy, and comprehensive success criteria.

### Test Execution

```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test package
go test ./internal/bridge
```

## License

MIT License - see LICENSE file for details.
