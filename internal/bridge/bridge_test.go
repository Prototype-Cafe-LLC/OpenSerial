package bridge

import (
	"testing"
	"time"

	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/config"
)

func TestNewBridge(t *testing.T) {
	cfg := &config.Config{
		Serial: config.SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    115200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: config.NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	bridge := NewBridge(cfg)

	if bridge == nil {
		t.Fatal("NewBridge should not return nil")
	}

	if bridge.config != cfg {
		t.Error("Bridge config should match provided config")
	}

	if bridge.serialHandler == nil {
		t.Error("Serial handler should not be nil")
	}

	if bridge.networkHandler == nil {
		t.Error("Network handler should not be nil")
	}

	if bridge.isRunning {
		t.Error("Bridge should not be running initially")
	}
}

func TestBridge_IsRunning(t *testing.T) {
	cfg := &config.Config{
		Serial: config.SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    115200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: config.NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	bridge := NewBridge(cfg)

	if bridge.IsRunning() {
		t.Error("Bridge should not be running initially")
	}

	bridge.mu.Lock()
	bridge.isRunning = true
	bridge.mu.Unlock()

	if !bridge.IsRunning() {
		t.Error("Bridge should be running after setting isRunning to true")
	}
}

func TestBridge_GetStatus(t *testing.T) {
	cfg := &config.Config{
		Serial: config.SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    115200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: config.NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	bridge := NewBridge(cfg)
	status := bridge.GetStatus()

	if status.IsRunning {
		t.Error("Status should show bridge not running initially")
	}

	if status.SerialOpen {
		t.Error("Status should show serial not open initially")
	}

	if status.NetworkListening {
		t.Error("Status should show network not listening initially")
	}

	if status.HasClient {
		t.Error("Status should show no client initially")
	}
}

func TestBridge_GetSerialConfig(t *testing.T) {
	cfg := &config.Config{
		Serial: config.SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    115200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: config.NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	bridge := NewBridge(cfg)
	serialConfig := bridge.GetSerialConfig()

	if serialConfig == nil {
		t.Fatal("GetSerialConfig should not return nil")
	}

	if serialConfig.Port != cfg.Serial.Port {
		t.Errorf("Expected port %s, got %s", cfg.Serial.Port, serialConfig.Port)
	}

	if serialConfig.BaudRate != cfg.Serial.BaudRate {
		t.Errorf("Expected baud rate %d, got %d", cfg.Serial.BaudRate, serialConfig.BaudRate)
	}
}

func TestBridge_GetNetworkConfig(t *testing.T) {
	cfg := &config.Config{
		Serial: config.SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    115200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: config.NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	bridge := NewBridge(cfg)
	networkConfig := bridge.GetNetworkConfig()

	if networkConfig == nil {
		t.Fatal("GetNetworkConfig should not return nil")
	}

	if networkConfig.ListenPort != cfg.Network.ListenPort {
		t.Errorf("Expected listen port %d, got %d", cfg.Network.ListenPort, networkConfig.ListenPort)
	}

	if networkConfig.BindAddress != cfg.Network.BindAddress {
		t.Errorf("Expected bind address %s, got %s", cfg.Network.BindAddress, networkConfig.BindAddress)
	}
}

func TestBridge_Start_AlreadyRunning(t *testing.T) {
	cfg := &config.Config{
		Serial: config.SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    115200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: config.NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	bridge := NewBridge(cfg)
	bridge.mu.Lock()
	bridge.isRunning = true
	bridge.mu.Unlock()

	err := bridge.Start()
	if err == nil {
		t.Error("Start should return error when already running")
	}
}

func TestBridge_Stop_NotRunning(t *testing.T) {
	cfg := &config.Config{
		Serial: config.SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    115200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: config.NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	bridge := NewBridge(cfg)

	err := bridge.Stop()
	if err != nil {
		t.Errorf("Stop should not return error when not running: %v", err)
	}
}

func TestBridge_Stop_Running(t *testing.T) {
	cfg := &config.Config{
		Serial: config.SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    115200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: config.NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	bridge := NewBridge(cfg)
	bridge.mu.Lock()
	bridge.isRunning = true
	bridge.mu.Unlock()

	// Create a new stop channel to avoid panic
	bridge.stopChan = make(chan struct{})

	err := bridge.Stop()
	if err != nil {
		t.Errorf("Stop should not return error when running: %v", err)
	}

	// Give goroutines time to finish
	time.Sleep(100 * time.Millisecond)

	if bridge.IsRunning() {
		t.Error("Bridge should not be running after stop")
	}
}
