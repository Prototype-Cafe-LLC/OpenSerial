package tcpbridge

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/config"
	"github.com/Prototype-Cafe-LLC/OpenSerial/pkg/logger"
)

// Server represents the TCP bridge server
type Server struct {
	config      *config.TCPBridgeConfig
	listener    net.Listener
	connections map[net.Conn]*Connection
	connMutex   sync.RWMutex
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

// Connection represents a client connection
type Connection struct {
	ClientConn net.Conn
	TargetConn net.Conn
	StopChan   chan struct{}
	mu         sync.Mutex
}

// NewServer creates a new TCP bridge server
func NewServer(cfg *config.TCPBridgeConfig) *Server {
	return &Server{
		config:      cfg,
		connections: make(map[net.Conn]*Connection),
		stopChan:    make(chan struct{}),
	}
}

// Start starts the TCP bridge server
func (s *Server) Start() error {
	// Start listening
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port))
	if err != nil {
		return fmt.Errorf("failed to start listening: %w", err)
	}

	s.listener = listener
	logger.Default.Info("TCP Bridge listening on %s:%d", s.config.Server.Host, s.config.Server.Port)

	// Start accepting connections
	s.wg.Add(1)
	go s.acceptConnections()

	// Start cleanup routine
	s.wg.Add(1)
	go s.cleanupRoutine()

	return nil
}

// Stop stops the TCP bridge server
func (s *Server) Stop() error {
	// Signal stop
	close(s.stopChan)

	// Close listener
	if s.listener != nil {
		s.listener.Close()
	}

	// Close all connections
	s.connMutex.Lock()
	for _, conn := range s.connections {
		conn.Close()
	}
	s.connections = make(map[net.Conn]*Connection)
	s.connMutex.Unlock()

	// Wait for goroutines to finish
	s.wg.Wait()

	return nil
}

// acceptConnections accepts incoming client connections
func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.stopChan:
			return
		default:
			// Accept connection
			clientConn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.stopChan:
					return
				default:
					logger.Default.Error("Failed to accept connection: %v", err)
					continue
				}
			}

			logger.Default.Info("Client connected from: %s", clientConn.RemoteAddr())

			// Check max connections
			s.connMutex.RLock()
			if len(s.connections) >= s.config.Clients.MaxConnections {
				s.connMutex.RUnlock()
				logger.Default.Error("Max connections reached, rejecting client: %s", clientConn.RemoteAddr())
				clientConn.Close()
				continue
			}
			s.connMutex.RUnlock()

			// Connect to target server
			targetConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.config.Target.Host, s.config.Target.Port))
			if err != nil {
				logger.Default.Error("Failed to connect to target server: %v", err)
				clientConn.Close()
				continue
			}

			logger.Default.Info("Connected to target server: %s:%d", s.config.Target.Host, s.config.Target.Port)

			// Create connection
			conn := &Connection{
				ClientConn: clientConn,
				TargetConn: targetConn,
				StopChan:   make(chan struct{}),
			}

			// Add to connections
			s.connMutex.Lock()
			s.connections[clientConn] = conn
			s.connMutex.Unlock()

			// Start data forwarding
			s.wg.Add(2)
			go s.forwardData(conn, clientConn, targetConn, "client->target")
			go s.forwardData(conn, targetConn, clientConn, "target->client")
		}
	}
}

// forwardData forwards data between two connections
func (s *Server) forwardData(conn *Connection, from, to net.Conn, direction string) {
	defer s.wg.Done()
	defer conn.Close()

	buffer := make([]byte, 4096)
	for {
		select {
		case <-conn.StopChan:
			return
		default:
			// Set read timeout
			from.SetReadDeadline(time.Now().Add(1 * time.Second))

			// Read data
			n, err := from.Read(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				logger.Default.Debug("Read error in %s: %v", direction, err)
				return
			}

			if n == 0 {
				continue
			}

			// Write data
			_, err = to.Write(buffer[:n])
			if err != nil {
				logger.Default.Debug("Write error in %s: %v", direction, err)
				return
			}

			logger.Default.Debug("Forwarded %d bytes: %s", n, direction)
		}
	}
}

// cleanupRoutine periodically cleans up closed connections
func (s *Server) cleanupRoutine() {
	defer s.wg.Done()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.cleanupClosedConnections()
		}
	}
}

// cleanupClosedConnections removes closed connections
func (s *Server) cleanupClosedConnections() {
	s.connMutex.Lock()
	defer s.connMutex.Unlock()

	for clientConn, conn := range s.connections {
		// Check if connection is still alive
		clientConn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
		_, err := clientConn.Read(make([]byte, 1))
		if err != nil {
			// Connection is closed
			logger.Default.Info("Cleaning up closed connection: %s", clientConn.RemoteAddr())
			conn.Close()
			delete(s.connections, clientConn)
		} else {
			// Connection is alive, reset deadline
			clientConn.SetReadDeadline(time.Time{})
		}
	}
}

// Close closes a connection
func (c *Connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.StopChan:
		return
	default:
		close(c.StopChan)
	}

	if c.ClientConn != nil {
		c.ClientConn.Close()
	}
	if c.TargetConn != nil {
		c.TargetConn.Close()
	}
}

// GetStatus returns the current status of the server
func (s *Server) GetStatus() ServerStatus {
	s.connMutex.RLock()
	defer s.connMutex.RUnlock()

	return ServerStatus{
		IsRunning:         s.listener != nil,
		ActiveConnections: len(s.connections),
		MaxConnections:    s.config.Clients.MaxConnections,
		TargetHost:        s.config.Target.Host,
		TargetPort:        s.config.Target.Port,
	}
}

// ServerStatus represents the current status of the server
type ServerStatus struct {
	IsRunning         bool   `json:"is_running"`
	ActiveConnections int    `json:"active_connections"`
	MaxConnections    int    `json:"max_connections"`
	TargetHost        string `json:"target_host"`
	TargetPort        int    `json:"target_port"`
}
