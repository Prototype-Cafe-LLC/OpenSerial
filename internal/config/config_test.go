package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	config, err := LoadConfig("")
	if err != nil {
		t.Fatalf("Failed to load default config: %v", err)
	}

	// Test default values
	if config.Serial.Port != "/dev/ttyUSB0" {
		t.Errorf("Expected default port /dev/ttyUSB0, got %s", config.Serial.Port)
	}

	if config.Serial.BaudRate != 115200 {
		t.Errorf("Expected default baud rate 115200, got %d", config.Serial.BaudRate)
	}

	if config.Serial.DataBits != 8 {
		t.Errorf("Expected default data bits 8, got %d", config.Serial.DataBits)
	}

	if config.Serial.StopBits != 1 {
		t.Errorf("Expected default stop bits 1, got %d", config.Serial.StopBits)
	}

	if config.Serial.Parity != "none" {
		t.Errorf("Expected default parity none, got %s", config.Serial.Parity)
	}

	if config.Serial.FlowControl != "none" {
		t.Errorf("Expected default flow control none, got %s", config.Serial.FlowControl)
	}

	if config.Network.ListenPort != 8080 {
		t.Errorf("Expected default listen port 8080, got %d", config.Network.ListenPort)
	}

	if config.Network.BindAddress != "0.0.0.0" {
		t.Errorf("Expected default bind address 0.0.0.0, got %s", config.Network.BindAddress)
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test_config.yaml")

	configContent := `
serial:
  port: "/dev/ttyACM0"
  baud_rate: 9600
  data_bits: 7
  stop_bits: 2
  parity: "even"
  flow_control: "rts_cts"

network:
  listen_port: 9090
  bind_address: "127.0.0.1"
`

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("Failed to load config from file: %v", err)
	}

	// Test loaded values
	if config.Serial.Port != "/dev/ttyACM0" {
		t.Errorf("Expected port /dev/ttyACM0, got %s", config.Serial.Port)
	}

	if config.Serial.BaudRate != 9600 {
		t.Errorf("Expected baud rate 9600, got %d", config.Serial.BaudRate)
	}

	if config.Serial.DataBits != 7 {
		t.Errorf("Expected data bits 7, got %d", config.Serial.DataBits)
	}

	if config.Serial.StopBits != 2 {
		t.Errorf("Expected stop bits 2, got %d", config.Serial.StopBits)
	}

	if config.Serial.Parity != "even" {
		t.Errorf("Expected parity even, got %s", config.Serial.Parity)
	}

	if config.Serial.FlowControl != "rts_cts" {
		t.Errorf("Expected flow control rts_cts, got %s", config.Serial.FlowControl)
	}

	if config.Network.ListenPort != 9090 {
		t.Errorf("Expected listen port 9090, got %d", config.Network.ListenPort)
	}

	if config.Network.BindAddress != "127.0.0.1" {
		t.Errorf("Expected bind address 127.0.0.1, got %s", config.Network.BindAddress)
	}
}

func TestValidateConfig_ValidConfig(t *testing.T) {
	config := &Config{
		Serial: SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    115200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	err := validateConfig(config)
	if err != nil {
		t.Errorf("Valid config should not produce error: %v", err)
	}
}

func TestValidateConfig_InvalidBaudRate(t *testing.T) {
	config := &Config{
		Serial: SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    999999, // Invalid baud rate
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Error("Invalid baud rate should produce error")
	}
}

func TestValidateConfig_InvalidDataBits(t *testing.T) {
	config := &Config{
		Serial: SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    115200,
			DataBits:    10, // Invalid data bits
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Error("Invalid data bits should produce error")
	}
}

func TestValidateConfig_InvalidPort(t *testing.T) {
	config := &Config{
		Serial: SerialConfig{
			Port:        "", // Empty port
			BaudRate:    115200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: NetworkConfig{
			ListenPort:  8080,
			BindAddress: "0.0.0.0",
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Error("Empty port should produce error")
	}
}

func TestValidateConfig_InvalidListenPort(t *testing.T) {
	config := &Config{
		Serial: SerialConfig{
			Port:        "/dev/ttyUSB0",
			BaudRate:    115200,
			DataBits:    8,
			StopBits:    1,
			Parity:      "none",
			FlowControl: "none",
		},
		Network: NetworkConfig{
			ListenPort:  99999, // Invalid port
			BindAddress: "0.0.0.0",
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Error("Invalid listen port should produce error")
	}
}
