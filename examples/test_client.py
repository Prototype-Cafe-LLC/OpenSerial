#!/usr/bin/env python3
"""
Simple test client for OpenSerial bridge.
This script connects to the OpenSerial TCP server and demonstrates
bidirectional communication.
"""

import socket
import time
import threading
import sys

def receive_data(sock):
    """Receive data from the server and print it."""
    while True:
        try:
            data = sock.recv(1024)
            if not data:
                break
            print(f"Received: {data.decode('utf-8', errors='ignore')}")
        except Exception as e:
            print(f"Error receiving data: {e}")
            break

def send_data(sock):
    """Send data to the server."""
    while True:
        try:
            message = input("Enter message to send (or 'quit' to exit): ")
            if message.lower() == 'quit':
                break
            sock.send(message.encode('utf-8'))
        except Exception as e:
            print(f"Error sending data: {e}")
            break

def main():
    if len(sys.argv) != 3:
        print("Usage: python test_client.py <host> <port>")
        print("Example: python test_client.py localhost 8080")
        sys.exit(1)

    host = sys.argv[1]
    port = int(sys.argv[2])

    try:
        # Connect to the OpenSerial server
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.connect((host, port))
        print(f"Connected to {host}:{port}")

        # Start receive thread
        receive_thread = threading.Thread(target=receive_data, args=(sock,))
        receive_thread.daemon = True
        receive_thread.start()

        # Start send thread
        send_thread = threading.Thread(target=send_data, args=(sock,))
        send_thread.daemon = True
        send_thread.start()

        # Wait for threads to complete
        send_thread.join()

    except KeyboardInterrupt:
        print("\nShutting down...")
    except Exception as e:
        print(f"Error: {e}")
    finally:
        sock.close()
        print("Disconnected")

if __name__ == "__main__":
    main()
