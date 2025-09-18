#!/bin/bash

# OpenSerial TCP Bridge Setup Script
# This script helps set up the TCP bridge solution for Mac and Windows

set -e

echo "OpenSerial TCP Bridge Setup"
echo "=========================="
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
if [ "$PLATFORM" == "mac" ]; then
    echo "Building for macOS..."
    go build -o openserial-mac ./cmd/openserial
    go build -o tcpbridge-mac ./cmd/tcpbridge
    echo "Built: openserial-mac, tcpbridge-mac"
else
    echo "Building for Windows..."
    go build -o openserial-windows.exe ./cmd/openserial
    go build -o tcpbridge-windows.exe ./cmd/tcpbridge
    echo "Built: openserial-windows.exe, tcpbridge-windows.exe"
fi

echo ""

# Create configuration files
echo "Creating configuration files..."

if [ "$PLATFORM" == "mac" ]; then
    echo "Setting up Mac configuration..."
    cp configs/server-mac.yaml config.yaml
    cp configs/tcpbridge.yaml tcpbridge.yaml
    echo "Created: config.yaml (server mode), tcpbridge.yaml"
    echo ""
    echo "To run on Mac:"
    echo "1. Start OpenSerial server: ./openserial-mac -config config.yaml"
    echo "2. Start TCP bridge: ./tcpbridge-mac -config tcpbridge.yaml"
    echo "3. Connect Cursor Serial Monitor to localhost:8081"
else
    echo "Setting up Windows configuration..."
    cp configs/client-windows.yaml config.yaml
    echo "Created: config.yaml (client mode)"
    echo ""
    echo "To run on Windows:"
    echo "1. Update config.yaml with Mac's IP address"
    echo "2. Start OpenSerial client: ./openserial-windows.exe -client -config config.yaml"
fi

echo ""
echo "Setup complete!"
echo ""
echo "Next steps:"
echo "1. Ensure both machines are connected to the same network"
echo "2. Find Mac's IP address (ifconfig on Mac, ipconfig on Windows)"
echo "3. Update Windows config.yaml with Mac's IP address"
echo "4. Start the applications in the correct order"
echo ""
echo "For more information, see README-TCP-BRIDGE.md"
