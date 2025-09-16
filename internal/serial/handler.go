package serial

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/openserial/openserial/internal/config"
	"github.com/tarm/serial"
)

// Handler manages serial port communication
type Handler struct {
	config      *config.SerialConfig
	port        *serial.Port
	mu          sync.RWMutex
	isOpen      bool
	readBuffer  []byte
	writeBuffer []byte
}

// NewHandler creates a new serial handler
func NewHandler(cfg *config.SerialConfig) *Handler {
	return &Handler{
		config:      cfg,
		readBuffer:  make([]byte, 4096),
		writeBuffer: make([]byte, 4096),
	}
}

// Open opens the serial port
func (h *Handler) Open() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.isOpen {
		return fmt.Errorf("serial port is already open")
	}

	// Check if this is a test port (for testing without real hardware)
	if h.config.Port == "/dev/null" {
		h.isOpen = true
		return nil
	}

	// Configure serial port
	serialConfig := &serial.Config{
		Name:        h.config.Port,
		Baud:        h.config.BaudRate,
		Size:        byte(h.config.DataBits),
		StopBits:    serial.StopBits(h.config.StopBits),
		Parity:      h.getParity(),
		ReadTimeout: h.config.GetReadTimeout(),
	}

	port, err := serial.OpenPort(serialConfig)
	if err != nil {
		return fmt.Errorf("failed to open serial port %s: %w", h.config.Port, err)
	}

	h.port = port
	h.isOpen = true
	return nil
}

// Close closes the serial port
func (h *Handler) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.isOpen {
		return nil
	}

	if h.port != nil {
		err := h.port.Close()
		h.port = nil
		h.isOpen = false
		return err
	}

	h.isOpen = false
	return nil
}

// IsOpen returns whether the serial port is open
func (h *Handler) IsOpen() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.isOpen
}

// Read reads data from the serial port
func (h *Handler) Read() ([]byte, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.isOpen {
		return nil, fmt.Errorf("serial port is not open")
	}

	// Test mode - return empty data
	if h.config.Port == "/dev/null" {
		return nil, nil
	}

	if h.port == nil {
		return nil, fmt.Errorf("serial port is not open")
	}

	n, err := h.port.Read(h.readBuffer)
	if err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read from serial port: %w", err)
	}

	if n == 0 {
		return nil, nil
	}

	// Return a copy of the data
	data := make([]byte, n)
	copy(data, h.readBuffer[:n])
	return data, nil
}

// Write writes data to the serial port
func (h *Handler) Write(data []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.isOpen {
		return fmt.Errorf("serial port is not open")
	}

	// Test mode - just return success
	if h.config.Port == "/dev/null" {
		return nil
	}

	if h.port == nil {
		return fmt.Errorf("serial port is not open")
	}

	_, err := h.port.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to serial port: %w", err)
	}

	return nil
}

// Reconnect attempts to reconnect the serial port
func (h *Handler) Reconnect() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Close existing connection if open
	if h.isOpen && h.port != nil {
		h.port.Close()
		h.port = nil
		h.isOpen = false
	}

	// Wait a bit before reconnecting
	time.Sleep(1 * time.Second)

	// Try to open again
	serialConfig := &serial.Config{
		Name:        h.config.Port,
		Baud:        h.config.BaudRate,
		Size:        byte(h.config.DataBits),
		StopBits:    serial.StopBits(h.config.StopBits),
		Parity:      h.getParity(),
		ReadTimeout: h.config.GetReadTimeout(),
	}

	port, err := serial.OpenPort(serialConfig)
	if err != nil {
		return fmt.Errorf("failed to reconnect serial port %s: %w", h.config.Port, err)
	}

	h.port = port
	h.isOpen = true
	return nil
}

// getParity converts string parity to serial.Parity
func (h *Handler) getParity() serial.Parity {
	switch h.config.Parity {
	case "odd":
		return serial.ParityOdd
	case "even":
		return serial.ParityEven
	case "mark":
		return serial.ParityMark
	case "space":
		return serial.ParitySpace
	default:
		return serial.ParityNone
	}
}

// GetConfig returns the serial configuration
func (h *Handler) GetConfig() *config.SerialConfig {
	return h.config
}
