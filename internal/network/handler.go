package network

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/config"
)

// Handler manages TCP network communication
type Handler struct {
	config      *config.NetworkConfig
	listener    net.Listener
	client      net.Conn
	mu          sync.RWMutex
	isListening bool
	hasClient   bool
	readBuffer  []byte
	writeBuffer []byte
}

// NewHandler creates a new network handler
func NewHandler(cfg *config.NetworkConfig) *Handler {
	return &Handler{
		config:      cfg,
		readBuffer:  make([]byte, 4096),
		writeBuffer: make([]byte, 4096),
	}
}

// StartListening starts listening for TCP connections
func (h *Handler) StartListening() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.isListening {
		return fmt.Errorf("already listening on port %d", h.config.ListenPort)
	}

	address := fmt.Sprintf("%s:%d", h.config.BindAddress, h.config.ListenPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start listening on %s: %w", address, err)
	}

	h.listener = listener
	h.isListening = true

	go h.acceptConnections()
	return nil
}

// StopListening stops listening for connections
func (h *Handler) StopListening() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.isListening {
		return nil
	}

	if h.listener != nil {
		err := h.listener.Close()
		h.listener = nil
		h.isListening = false
		return err
	}

	h.isListening = false
	return nil
}

// IsListening returns whether the handler is listening for connections
func (h *Handler) IsListening() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.isListening
}

// HasClient returns whether there is an active client connection
func (h *Handler) HasClient() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.hasClient
}

// Read reads data from the client connection
func (h *Handler) Read() ([]byte, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.hasClient || h.client == nil {
		return nil, fmt.Errorf("no client connection")
	}

	// Set read timeout
	h.client.SetReadDeadline(time.Now().Add(100 * time.Millisecond))

	n, err := h.client.Read(h.readBuffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, nil // Timeout is not an error
		}
		// Connection lost
		h.client.Close()
		h.client = nil
		h.hasClient = false
		return nil, fmt.Errorf("failed to read from client: %w", err)
	}

	if n == 0 {
		return nil, nil
	}

	// Return a copy of the data
	data := make([]byte, n)
	copy(data, h.readBuffer[:n])
	return data, nil
}

// Write writes data to the client connection
func (h *Handler) Write(data []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.hasClient || h.client == nil {
		return fmt.Errorf("no client connection")
	}

	_, err := h.client.Write(data)
	if err != nil {
		// Connection lost
		h.client.Close()
		h.client = nil
		h.hasClient = false
		return fmt.Errorf("failed to write to client: %w", err)
	}

	return nil
}

// CloseClient closes the current client connection
func (h *Handler) CloseClient() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.hasClient || h.client == nil {
		return nil
	}

	err := h.client.Close()
	h.client = nil
	h.hasClient = false
	return err
}

// acceptConnections accepts incoming TCP connections
func (h *Handler) acceptConnections() {
	for {
		h.mu.RLock()
		if !h.isListening || h.listener == nil {
			h.mu.RUnlock()
			break
		}
		listener := h.listener
		h.mu.RUnlock()

		conn, err := listener.Accept()
		if err != nil {
			h.mu.RLock()
			if !h.isListening {
				h.mu.RUnlock()
				break
			}
			h.mu.RUnlock()
			continue
		}

		// Close existing client if any
		h.mu.Lock()
		if h.hasClient && h.client != nil {
			h.client.Close()
		}
		h.client = conn
		h.hasClient = true
		h.mu.Unlock()
	}
}

// GetConfig returns the network configuration
func (h *Handler) GetConfig() *config.NetworkConfig {
	return h.config
}

// GetClientAddress returns the address of the current client
func (h *Handler) GetClientAddress() string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if !h.hasClient || h.client == nil {
		return ""
	}

	return h.client.RemoteAddr().String()
}
