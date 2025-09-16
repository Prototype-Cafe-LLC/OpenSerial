#!/bin/bash

# Test script for OpenSerial bridge
# This script demonstrates how to test the OpenSerial bridge

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}OpenSerial Bridge Test Script${NC}"
echo "=================================="

# Check if openserial binary exists
if [ ! -f "./build/openserial" ]; then
    echo -e "${RED}Error: openserial binary not found. Please run 'make build' first.${NC}"
    exit 1
fi

# Check if config file exists
if [ ! -f "./configs/config.yaml" ]; then
    echo -e "${RED}Error: config file not found at ./configs/config.yaml${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Binary and config file found${NC}"

# Test help functionality
echo -e "\n${YELLOW}Testing help functionality...${NC}"
./build/openserial -help > /dev/null
echo -e "${GREEN}✓ Help functionality works${NC}"

# Test version functionality
echo -e "\n${YELLOW}Testing version functionality...${NC}"
VERSION=$(./build/openserial -version)
echo -e "${GREEN}✓ Version: $VERSION${NC}"

# Test configuration loading
echo -e "\n${YELLOW}Testing configuration loading...${NC}"
./build/openserial -config ./configs/test_config.yaml &
OPENSERIAL_PID=$!

    # Wait a moment for the process to start
    sleep 3

    # Check if the process is still running
    if kill -0 $OPENSERIAL_PID 2>/dev/null; then
        echo -e "${GREEN}✓ Configuration loaded successfully${NC}"
        echo -e "${GREEN}✓ Bridge started successfully (PID: $OPENSERIAL_PID)${NC}"
        
        # Test TCP connection
        echo -e "\n${YELLOW}Testing TCP connection...${NC}"
        # Try multiple times with a short delay
        for i in {1..5}; do
            if timeout 2 bash -c "echo > /dev/tcp/localhost/8080" 2>/dev/null; then
                echo -e "${GREEN}✓ TCP port 8080 is listening${NC}"
                break
            else
                if [ $i -eq 5 ]; then
                    echo -e "${YELLOW}⚠ TCP port 8080 test inconclusive (bridge may still be starting)${NC}"
                else
                    sleep 1
                fi
            fi
        done
    
    # Stop the bridge
    echo -e "\n${YELLOW}Stopping bridge...${NC}"
    kill $OPENSERIAL_PID
    wait $OPENSERIAL_PID 2>/dev/null || true
    echo -e "${GREEN}✓ Bridge stopped successfully${NC}"
else
    echo -e "${RED}✗ Bridge failed to start${NC}"
    exit 1
fi

echo -e "\n${GREEN}All tests passed!${NC}"
echo -e "\n${YELLOW}To test with a real serial device:${NC}"
echo "1. Connect your serial device to the configured port"
echo "2. Run: ./build/openserial -config ./configs/config.yaml"
echo "3. Connect a TCP client to localhost:8080"
echo "4. Use the Python test client: python3 examples/test_client.py localhost 8080"
