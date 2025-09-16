package bridge

import (
	"fmt"
	"sync"
	"time"

	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/config"
	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/network"
	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/serial"
)

// Bridge manages the bidirectional data forwarding between serial and network
type Bridge struct {
	config         *config.Config
	serialHandler  *serial.Handler
	networkHandler *network.Handler
	mu             sync.RWMutex
	isRunning      bool
	stopChan       chan struct{}
	wg             sync.WaitGroup
}

// NewBridge creates a new bridge instance
func NewBridge(cfg *config.Config) *Bridge {
	return &Bridge{
		config:         cfg,
		serialHandler:  serial.NewHandler(&cfg.Serial),
		networkHandler: network.NewHandler(&cfg.Network),
		stopChan:       make(chan struct{}),
	}
}

// Start starts the bridge
func (b *Bridge) Start() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.isRunning {
		return fmt.Errorf("bridge is already running")
	}

	// Open serial port
	if err := b.serialHandler.Open(); err != nil {
		return fmt.Errorf("failed to open serial port: %w", err)
	}

	// Start network listener
	if err := b.networkHandler.StartListening(); err != nil {
		b.serialHandler.Close()
		return fmt.Errorf("failed to start network listener: %w", err)
	}

	b.isRunning = true

	// Start data forwarding goroutines
	b.wg.Add(2)
	go b.forwardSerialToNetwork()
	go b.forwardNetworkToSerial()

	return nil
}

// Stop stops the bridge
func (b *Bridge) Stop() error {
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
	b.networkHandler.StopListening()
	b.networkHandler.CloseClient()

	b.isRunning = false
	return nil
}

// IsRunning returns whether the bridge is running
func (b *Bridge) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.isRunning
}

// GetStatus returns the current status of the bridge
func (b *Bridge) GetStatus() Status {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return Status{
		IsRunning:        b.isRunning,
		SerialOpen:       b.serialHandler.IsOpen(),
		NetworkListening: b.networkHandler.IsListening(),
		HasClient:        b.networkHandler.HasClient(),
		ClientAddress:    b.networkHandler.GetClientAddress(),
	}
}

// forwardSerialToNetwork forwards data from serial port to network
func (b *Bridge) forwardSerialToNetwork() {
	defer b.wg.Done()

	for {
		select {
		case <-b.stopChan:
			return
		default:
			// Check if we have a client connection
			if !b.networkHandler.HasClient() {
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
			if err := b.networkHandler.Write(data); err != nil {
				// Client disconnected, continue waiting for new connection
				continue
			}
		}
	}
}

// forwardNetworkToSerial forwards data from network to serial port
func (b *Bridge) forwardNetworkToSerial() {
	defer b.wg.Done()

	for {
		select {
		case <-b.stopChan:
			return
		default:
			// Check if we have a client connection
			if !b.networkHandler.HasClient() {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// Read from network
			data, err := b.networkHandler.Read()
			if err != nil {
				// Client disconnected, continue waiting for new connection
				continue
			}

			if len(data) == 0 {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			// Write to serial port
			if err := b.serialHandler.Write(data); err != nil {
				// Try to reconnect serial port
				if reconnectErr := b.serialHandler.Reconnect(); reconnectErr != nil {
					time.Sleep(1 * time.Second)
					continue
				}
				continue
			}
		}
	}
}

// Status represents the current status of the bridge
type Status struct {
	IsRunning        bool   `json:"is_running"`
	SerialOpen       bool   `json:"serial_open"`
	NetworkListening bool   `json:"network_listening"`
	HasClient        bool   `json:"has_client"`
	ClientAddress    string `json:"client_address,omitempty"`
}

// GetSerialConfig returns the serial configuration
func (b *Bridge) GetSerialConfig() *config.SerialConfig {
	return b.serialHandler.GetConfig()
}

// GetNetworkConfig returns the network configuration
func (b *Bridge) GetNetworkConfig() *config.NetworkConfig {
	return b.networkHandler.GetConfig()
}
