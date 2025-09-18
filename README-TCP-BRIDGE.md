# OpenSerial TCP Bridge Solution

This solution enables OpenSerial to work across network connections when Windows machines don't have admin privileges for firewall configuration.

## Architecture

```text
Windows Machine (No Admin)          Mac Machine (Admin)
┌─────────────────────────┐        ┌─────────────────────────┐
│  UART Device            │        │  UART Device            │
│  (Arduino/Serial)       │        │  (Arduino/Serial)       │
└─────────┬───────────────┘        └─────────┬───────────────┘
          │                                  │
          │                                  │
┌─────────▼───────────────┐        ┌─────────▼───────────────┐
│  OpenSerial Client      │◄──────►│  TCP Bridge Server      │
│  (Connects to Mac)      │        │  (Runs on Mac)          │
└─────────────────────────┘        └─────────┬───────────────┘
                                             │
                                             │
                                    ┌─────────▼───────────────┐
                                    │  OpenSerial Server      │
                                    │  (Listens on Mac)       │
                                    └─────────┬───────────────┘
                                             │
                                    ┌─────────▼───────────────┐
                                    │  Cursor IDE             │
                                    │  Serial Monitor         │
                                    └─────────────────────────┘
```

## How It Works

1. **Mac Side**:
   - OpenSerial runs in **server mode** (listens on port 8081)
   - TCP Bridge Server runs and connects to OpenSerial's port 8081
   - Bridge Server accepts connections from Windows on port 8080
   - Cursor Serial Monitor connects to OpenSerial's port 8081

2. **Windows Side**:
   - OpenSerial runs in **client mode** (connects to Mac's bridge server)
   - No need to listen on any port (outbound connection only)
   - Forwards UART data to/from the TCP connection

## Setup Instructions

### Prerequisites

- Both machines connected to the same network
- Go 1.23+ installed on both machines
- Admin privileges on Mac machine
- No admin privileges required on Windows

### Mac Setup (Server Side)

1. **Build the applications**:

   ```bash
   go build -o openserial-mac ./cmd/openserial
   go build -o tcpbridge-mac ./cmd/tcpbridge
   ```

2. **Configure OpenSerial server**:

   ```bash
   cp configs/server-mac.yaml config.yaml
   # Edit config.yaml if needed (serial port, baud rate, etc.)
   ```

3. **Configure TCP bridge**:

   ```bash
   cp configs/tcpbridge.yaml tcpbridge.yaml
   # Edit tcpbridge.yaml if needed
   ```

4. **Start the services**:

   ```bash
   # Terminal 1: Start OpenSerial server
   ./openserial-mac -config config.yaml
   
   # Terminal 2: Start TCP bridge
   ./tcpbridge-mac -config tcpbridge.yaml
   ```

5. **Connect Cursor Serial Monitor**:
   - Open Cursor IDE
   - Open Serial Monitor
   - Connect to `localhost:8081`

### Windows Setup (Client Side)

1. **Build the application**:

   ```bash
   go build -o openserial-windows.exe ./cmd/openserial
   ```

2. **Configure OpenSerial client**:

   ```bash
   cp configs/client-windows.yaml config.yaml
   ```

3. **Update configuration**:
   - Edit `config.yaml`
   - Change `bind_address` to Mac's IP address
   - Update `listen_port` to 8080 (bridge server port)
   - Update serial port (e.g., "COM3")

4. **Start OpenSerial client**:

   ```bash
   ./openserial-windows.exe -client -config config.yaml
   ```

## Configuration Files

### Mac Server Configuration (`config.yaml`)

```yaml
serial:
  port: "/dev/ttyUSB0"
  baud_rate: 115200
  data_bits: 8
  stop_bits: 1
  parity: "none"
  flow_control: "none"

network:
  listen_port: 8081
  bind_address: "0.0.0.0"
```

### TCP Bridge Configuration (`tcpbridge.yaml`)

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

### Windows Client Configuration (`config.yaml`)

```yaml
serial:
  port: "COM3"
  baud_rate: 115200
  data_bits: 8
  stop_bits: 1
  parity: "none"
  flow_control: "none"

network:
  listen_port: 8080
  bind_address: "192.168.1.100"  # Mac's IP address
```

## Finding IP Addresses

### On Mac

```bash
ifconfig | grep "inet " | grep -v 127.0.0.1
```

### On Windows

```cmd
ipconfig | findstr "IPv4"
```

## Troubleshooting

### Connection Issues

- Verify both machines are on the same network
- Check firewall settings on Mac (allow port 8080)
- Ensure OpenSerial server is running before starting bridge
- Check that serial ports are correct and devices are connected

### Serial Port Issues

- On Mac: Check `/dev/tty*` devices
- On Windows: Check Device Manager for COM ports
- Ensure devices are not in use by other applications

### Network Issues

- Test connectivity: `telnet <mac-ip> 8080` from Windows
- Check network connection status
- Verify IP addresses are correct

## Benefits

- **No firewall issues**: Windows only makes outbound connections
- **No admin required on Windows**: Just run the client
- **Simple implementation**: Minimal changes to existing code
- **Direct connection**: No cloud costs or complexity
- **Works over network**: Both machines on same network
- **Low latency**: Direct TCP connection

## Alternative Solutions

If this approach doesn't work, consider:

1. **Cloud Bridge**: More complex but works over internet
2. **SSH Tunneling**: If SSH access is available
3. **Router Port Forwarding**: Configure router for port forwarding
4. **USB over Network**: Use USB-over-IP solutions

## Files Created

- `cmd/tcpbridge/` - TCP bridge server application
- `internal/tcpbridge/` - TCP bridge server implementation
- `internal/bridge/client_bridge.go` - OpenSerial client mode
- `configs/tcpbridge.yaml` - Bridge server configuration
- `configs/server-mac.yaml` - Mac server configuration
- `configs/client-windows.yaml` - Windows client configuration
- `scripts/setup-tcp-bridge.sh` - Setup script
