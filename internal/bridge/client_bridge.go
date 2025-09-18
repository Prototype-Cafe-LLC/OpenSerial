package bridge

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/config"
	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/serial"
	"github.com/Prototype-Cafe-LLC/OpenSerial/pkg/logger"
)

// ClientBridge manages the bidirectional data forwarding between serial and TCP client connection
type ClientBridge struct {
	config            *config.Config
	serialHandler     *serial.Handler
	conn              net.Conn
	mu                sync.RWMutex
	isRunning         bool
	stopChan          chan struct{}
	wg                sync.WaitGroup
	reconnectAttempts int
}

// NewClientBridge creates a new client bridge instance
func NewClientBridge(cfg *config.Config) *ClientBridge {
	return &ClientBridge{
		config:        cfg,
		serialHandler: serial.NewHandler(&cfg.Serial),
		stopChan:      make(chan struct{}),
	}
}

// Start starts the client bridge
func (b *ClientBridge) Start() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.isRunning {
		return fmt.Errorf("client bridge is already running")
	}

	// Open serial port
	if err := b.serialHandler.Open(); err != nil {
		return fmt.Errorf("failed to open serial port: %w", err)
	}

	// Start connection loop
	b.wg.Add(1)
	go b.connectionLoop()

	b.isRunning = true
	return nil
}

// Stop stops the client bridge
func (b *ClientBridge) Stop() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.isRunning {
		return nil
	}

	// Signal stop
	close(b.stopChan)

	// Wait for goroutines to finish
	b.wg.Wait()

	// Close connections
	b.serialHandler.Close()
	b.closeConnection()

	b.isRunning = false
	return nil
}

// connectionLoop handles connection and reconnection
func (b *ClientBridge) connectionLoop() {
	defer b.wg.Done()

	for {
		select {
		case <-b.stopChan:
			return
		default:
			if err := b.connect(); err != nil {
				logger.Default.Error("Failed to connect: %v", err)

				// Check if we should retry
				maxAttempts := 10 // Default max attempts
				if b.reconnectAttempts >= maxAttempts {
					logger.Default.Fatal("Max reconnection attempts reached")
					return
				}

				// Wait before retry
				select {
				case <-b.stopChan:
					return
				case <-time.After(5 * time.Second):
					b.reconnectAttempts++
					continue
				}
			}

			// Reset reconnect attempts on successful connection
			b.reconnectAttempts = 0

			// Start data forwarding
			b.wg.Add(2)
			go b.forwardSerialToNetwork()
			go b.forwardNetworkToSerial()

			// Wait for connection to close
			select {
			case <-b.stopChan:
				return
			case <-time.After(time.Second):
				// Check if connection is still alive
				if !b.isConnected() {
					logger.Default.Info("Connection lost, attempting to reconnect...")
					b.closeConnection()
					break
				}
			}
		}
	}
}

// connect establishes connection to the server
func (b *ClientBridge) connect() error {
	// Connect to server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", b.config.Network.BindAddress, b.config.Network.ListenPort))
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	b.mu.Lock()
	b.conn = conn
	b.mu.Unlock()

	logger.Default.Info("Connected to server: %s:%d", b.config.Network.BindAddress, b.config.Network.ListenPort)
	return nil
}

// closeConnection closes the TCP connection
func (b *ClientBridge) closeConnection() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.conn != nil {
		b.conn.Close()
		b.conn = nil
	}
}

// isConnected returns whether the client is connected
func (b *ClientBridge) isConnected() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.conn != nil
}

// forwardSerialToNetwork forwards data from serial port to network
func (b *ClientBridge) forwardSerialToNetwork() {
	defer b.wg.Done()

	for {
		select {
		case <-b.stopChan:
			return
		default:
			if !b.isConnected() {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// Read from serial port
			data, err := b.serialHandler.Read()
			if err != nil {
				// Try to reconnect serial port
				if reconnectErr := b.serialHandler.Reconnect(); reconnectErr != nil {
					time.Sleep(1 * time.Second)
					continue
				}
				continue
			}

			if len(data) == 0 {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			// Write to network
			b.mu.RLock()
			conn := b.conn
			b.mu.RUnlock()

			if conn != nil {
				if _, err := conn.Write(data); err != nil {
					logger.Default.Error("Failed to write to network: %v", err)
					b.closeConnection()
					return
				}
			}
		}
	}
}

// forwardNetworkToSerial forwards data from network to serial port
func (b *ClientBridge) forwardNetworkToSerial() {
	defer b.wg.Done()

	buffer := make([]byte, 4096)
	for {
		select {
		case <-b.stopChan:
			return
		default:
			if !b.isConnected() {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// Read from network
			b.mu.RLock()
			conn := b.conn
			b.mu.RUnlock()

			if conn == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// Set read timeout
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))

			n, err := conn.Read(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				logger.Default.Error("Failed to read from network: %v", err)
				b.closeConnection()
				return
			}

			if n == 0 {
				continue
			}

			// Write to serial port
			if err := b.serialHandler.Write(buffer[:n]); err != nil {
				logger.Default.Error("Failed to write to serial port: %v", err)
				// Try to reconnect serial port
				if reconnectErr := b.serialHandler.Reconnect(); reconnectErr != nil {
					time.Sleep(1 * time.Second)
				}
				continue
			}
		}
	}
}

// IsRunning returns whether the bridge is running
func (b *ClientBridge) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.isRunning
}

// GetStatus returns the current status of the bridge
func (b *ClientBridge) GetStatus() ClientStatus {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return ClientStatus{
		IsRunning:         b.isRunning,
		SerialOpen:        b.serialHandler.IsOpen(),
		NetworkConnected:  b.conn != nil,
		ReconnectAttempts: b.reconnectAttempts,
	}
}

// ClientStatus represents the current status of the client bridge
type ClientStatus struct {
	IsRunning         bool `json:"is_running"`
	SerialOpen        bool `json:"serial_open"`
	NetworkConnected  bool `json:"network_connected"`
	ReconnectAttempts int  `json:"reconnect_attempts"`
}
