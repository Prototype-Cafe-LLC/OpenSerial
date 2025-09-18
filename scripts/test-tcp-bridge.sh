#!/bin/bash

# OpenSerial TCP Bridge Test Script
# This script demonstrates the TCP bridge solution

set -e

echo "OpenSerial TCP Bridge Test"
echo "========================="
echo ""

# Check if we're on Mac or Windows
if [[ "$OSTYPE" == "darwin"* ]]; then
    PLATFORM="mac"
    echo "Detected platform: macOS"
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
    PLATFORM="windows"
    echo "Detected platform: Windows"
else
    echo "Unsupported platform: $OSTYPE"
    exit 1
fi

echo ""

# Build the applications
echo "Building applications..."
go build -o build/openserial ./cmd/openserial
go build -o build/tcpbridge ./cmd/tcpbridge
echo "Build complete!"
echo ""

if [ "$PLATFORM" == "mac" ]; then
    echo "Mac Setup Instructions:"
    echo "======================"
    echo ""
    echo "1. Start OpenSerial server (Terminal 1):"
    echo "   ./build/openserial -config configs/server-mac.yaml"
    echo ""
    echo "2. Start TCP bridge (Terminal 2):"
    echo "   ./build/tcpbridge -config configs/tcpbridge.yaml"
    echo ""
    echo "3. Connect Cursor Serial Monitor to localhost:8081"
    echo ""
    echo "4. Find your VPN IP address:"
    echo "   ifconfig | grep 'inet ' | grep -v 127.0.0.1"
    echo ""
    echo "5. Share the VPN IP with Windows user"
    echo ""
else
    echo "Windows Setup Instructions:"
    echo "=========================="
    echo ""
    echo "1. Update configs/client-windows.yaml with Mac's VPN IP address"
    echo ""
    echo "2. Start OpenSerial client:"
    echo "   ./build/openserial.exe -client -config configs/client-windows.yaml"
    echo ""
    echo "3. Find your VPN IP address:"
    echo "   ipconfig | findstr 'IPv4'"
    echo ""
fi

echo ""
echo "Test the connection:"
echo "==================="
echo ""
echo "From Windows, test connectivity to Mac:"
echo "telnet <mac-vpn-ip> 8080"
echo ""
echo "If connection is successful, you should see:"
echo "- OpenSerial client connects to bridge"
echo "- Bridge connects to OpenSerial server"
echo "- Data flows between serial ports"
echo ""
echo "Troubleshooting:"
echo "==============="
echo "- Ensure both machines are on the same VPN"
echo "- Check firewall settings on Mac (allow port 8080)"
echo "- Verify serial ports are correct and devices connected"
echo "- Check that OpenSerial server is running before bridge"
echo ""
echo "For more information, see README-TCP-BRIDGE.md"
