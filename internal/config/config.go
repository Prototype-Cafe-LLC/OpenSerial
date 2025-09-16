package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config represents the complete configuration structure
type Config struct {
	Serial  SerialConfig  `mapstructure:"serial"`
	Network NetworkConfig `mapstructure:"network"`
}

// SerialConfig represents serial port configuration
type SerialConfig struct {
	Port        string `mapstructure:"port"`
	BaudRate    int    `mapstructure:"baud_rate"`
	DataBits    int    `mapstructure:"data_bits"`
	StopBits    int    `mapstructure:"stop_bits"`
	Parity      string `mapstructure:"parity"`
	FlowControl string `mapstructure:"flow_control"`
}

// NetworkConfig represents network configuration
type NetworkConfig struct {
	ListenPort  int    `mapstructure:"listen_port"`
	BindAddress string `mapstructure:"bind_address"`
}

// LoadConfig loads configuration from file or creates default
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/etc/openserial")
	viper.AddConfigPath("$HOME/.openserial")

	// Set default values
	setDefaults()

	// Override with config file if provided
	if configPath != "" {
		viper.SetConfigFile(configPath)
	}

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, use defaults
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Serial defaults
	viper.SetDefault("serial.port", "/dev/ttyUSB0")
	viper.SetDefault("serial.baud_rate", 115200)
	viper.SetDefault("serial.data_bits", 8)
	viper.SetDefault("serial.stop_bits", 1)
	viper.SetDefault("serial.parity", "none")
	viper.SetDefault("serial.flow_control", "none")

	// Network defaults
	viper.SetDefault("network.listen_port", 8080)
	viper.SetDefault("network.bind_address", "0.0.0.0")
}

// validateConfig validates the configuration values
func validateConfig(config *Config) error {
	// Validate serial configuration
	if config.Serial.Port == "" {
		return fmt.Errorf("serial port cannot be empty")
	}

	validBaudRates := map[int]bool{
		300: true, 600: true, 1200: true, 2400: true, 4800: true,
		9600: true, 19200: true, 38400: true, 57600: true, 115200: true,
		230400: true, 460800: true, 921600: true,
	}
	if !validBaudRates[config.Serial.BaudRate] {
		return fmt.Errorf("invalid baud rate: %d", config.Serial.BaudRate)
	}

	if config.Serial.DataBits < 5 || config.Serial.DataBits > 8 {
		return fmt.Errorf("data bits must be between 5 and 8, got: %d", config.Serial.DataBits)
	}

	if config.Serial.StopBits < 1 || config.Serial.StopBits > 2 {
		return fmt.Errorf("stop bits must be 1 or 2, got: %d", config.Serial.StopBits)
	}

	validParity := map[string]bool{
		"none": true, "odd": true, "even": true, "mark": true, "space": true,
	}
	if !validParity[config.Serial.Parity] {
		return fmt.Errorf("invalid parity: %s", config.Serial.Parity)
	}

	validFlowControl := map[string]bool{
		"none": true, "rts_cts": true, "xon_xoff": true,
	}
	if !validFlowControl[config.Serial.FlowControl] {
		return fmt.Errorf("invalid flow control: %s", config.Serial.FlowControl)
	}

	// Validate network configuration
	if config.Network.ListenPort < 1 || config.Network.ListenPort > 65535 {
		return fmt.Errorf("listen port must be between 1 and 65535, got: %d", config.Network.ListenPort)
	}

	if config.Network.BindAddress == "" {
		return fmt.Errorf("bind address cannot be empty")
	}

	return nil
}

// GetSerialTimeout returns the serial port timeout
func (c *SerialConfig) GetSerialTimeout() time.Duration {
	return 1 * time.Second
}

// GetReadTimeout returns the read timeout
func (c *SerialConfig) GetReadTimeout() time.Duration {
	return 100 * time.Millisecond
}
