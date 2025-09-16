# OpenSerial

[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)](https://github.com/Prototype-Cafe-LLC/OpenSerial)

OpenSerial is a lightweight, cross-platform UART-to-TCP bridge that enables network access to serial devices. Perfect for IoT development, industrial applications, and remote debugging of embedded systems.

## 🚀 Quick Start

```bash
# Download the latest release for your platform
# Or build from source:
git clone https://github.com/Prototype-Cafe-LLC/OpenSerial.git
cd OpenSerial
go build -o openserial

# Run with configuration file
./openserial --config config.yaml
```

## ✨ Features

### Core Functionality

- **Bidirectional Serial Bridge**: Real-time serial port ↔ TCP port forwarding
- **Cross-Platform**: Windows, macOS (Intel/Apple Silicon), Linux (x64/ARM64)
- **Single Binary**: No external dependencies, easy deployment
- **Configuration File**: YAML/JSON configuration support

### Serial Communication

- Full serial parameter configuration (baud rate, data bits, stop bits, parity, flow control)
- Automatic serial port reconnection on device disconnection
- Buffer management with overflow protection

### Network Features

- TCP server mode with configurable listening port
- Bind to all network interfaces (0.0.0.0)
- Automatic client connection handling and cleanup
- Real-time connection status monitoring

## 📋 System Requirements

| Resource | Minimum | Recommended |
|----------|---------|-------------|
| **RAM** | 512MB | 1GB |
| **Storage** | 50MB | 100MB |
| **OS** | Windows 10+, macOS 10.15+, Linux (kernel 3.10+) | |

## ⚙️ Configuration

Create a `config.yaml` file:

```yaml
serial:
  port: "/dev/ttyUSB0"        # Serial port path
  baud_rate: 115200           # Baud rate
  data_bits: 8                # Data bits (5-8)
  stop_bits: 1                # Stop bits (1-2)
  parity: "none"              # Parity: none, odd, even
  flow_control: "none"        # Flow control: none, rts/cts, xon/xoff

network:
  listen_port: 8080           # TCP listening port
  bind_address: "0.0.0.0"     # Network interface binding
```

## 🛠️ Installation

### Binary Releases

Download pre-compiled executables from the [Releases](https://github.com/Prototype-Cafe-LLC/OpenSerial/releases) page.

### Build from Source

```bash
# Prerequisites: Go 1.25+
git clone https://github.com/Prototype-Cafe-LLC/OpenSerial.git
cd OpenSerial
go mod download
go build -o openserial
```

## 📖 Usage

### Basic Usage

```bash
# Run with configuration file
./openserial --config config.yaml

# On Windows
openserial.exe --config config.yaml
```

### Connecting to the Bridge

Once running, connect to the serial device via TCP:

```bash
# Using telnet
telnet localhost 8080

# Using netcat
nc localhost 8080

# Using Python
import socket
s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
s.connect(('localhost', 8080))
```

## 🎯 Use Cases

### Industrial Applications

- Remote monitoring of PLCs and industrial equipment
- SCADA system integration
- Equipment diagnostics and maintenance

### IoT Development

- Remote debugging of Arduino/Raspberry Pi projects
- IoT device management and configuration
- Sensor data collection and monitoring

### Educational & Hobby

- Learning serial communication protocols
- Electronics project development
- Remote access to embedded systems

## 🔧 Development

### Project Structure

```text
OpenSerial/
├── cmd/           # Main application entry point
├── internal/      # Internal packages
│   ├── bridge/    # Serial-TCP bridge logic
│   ├── config/    # Configuration management
│   └── serial/    # Serial port handling
├── pkg/           # Public packages
├── configs/       # Example configurations
└── docs/          # Documentation
```

### Building and Testing

```bash
# Install dependencies
go mod download

# Build
go build -o openserial ./cmd/openserial

# Run tests
go test ./...

# Format code
gofmt -w .

# Lint code
golangci-lint run
```

## 📊 Performance

- **Latency**: < 10ms for local network connections
- **Throughput**: Up to 1Mbps serial data rates
- **Memory**: < 50MB per active bridge
- **CPU**: Minimal overhead, optimized for efficiency

## 🔒 Security Considerations

⚠️ **Important**: OpenSerial currently runs without authentication or encryption. Consider your network security when deploying:

- Use in trusted networks only
- Consider VPN or SSH tunneling for remote access
- Monitor network traffic and connections
- Future versions will include TLS encryption and authentication

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
# Clone your fork
git clone git@github.com:yourusername/OpenSerial.git
cd OpenSerial

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o openserial ./cmd/openserial
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/Prototype-Cafe-LLC/OpenSerial/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Prototype-Cafe-LLC/OpenSerial/discussions)
- **Documentation**: [Wiki](https://github.com/Prototype-Cafe-LLC/OpenSerial/wiki)

## 🗺️ Roadmap

### Current Scope (v1.0)

- ✅ Basic serial-to-TCP forwarding
- ✅ Cross-platform support
- ✅ Configuration file support
- ✅ Single client connection

### Future Enhancements

- 🔄 Multiple serial port support
- 🔄 Web-based management interface
- 🔄 TLS encryption support
- 🔄 Authentication and authorization
- 🔄 Advanced logging and monitoring

---

**Made with ❤️ by [Prototype Cafe LLC](https://prototypecafe.com)**
