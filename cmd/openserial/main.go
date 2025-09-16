package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/bridge"
	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/config"
	"github.com/Prototype-Cafe-LLC/OpenSerial/pkg/logger"
)

const (
	version = "1.0.0"
)

func main() {
	var (
		configPath  = flag.String("config", "", "Path to configuration file")
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("OpenSerial v%s\n", version)
		os.Exit(0)
	}

	if *showHelp {
		showUsage()
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Default.Fatal("Failed to load configuration: %v", err)
	}

	// Create bridge
	bridge := bridge.NewBridge(cfg)

	// Start bridge
	logger.Default.Info("Starting OpenSerial bridge...")
	logger.Default.Info("Serial port: %s", cfg.Serial.Port)
	logger.Default.Info("Network: %s:%d", cfg.Network.BindAddress, cfg.Network.ListenPort)

	if err := bridge.Start(); err != nil {
		logger.Default.Fatal("Failed to start bridge: %v", err)
	}

	logger.Default.Info("Bridge started successfully")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	logger.Default.Info("Shutdown signal received, stopping bridge...")

	// Stop bridge
	if err := bridge.Stop(); err != nil {
		logger.Default.Error("Error stopping bridge: %v", err)
		os.Exit(1)
	}

	logger.Default.Info("Bridge stopped successfully")
}

func showUsage() {
	fmt.Printf("OpenSerial v%s - UART-to-TCP Bridge\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  openserial [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -config string")
	fmt.Println("        Path to configuration file (YAML format)")
	fmt.Println("  -version")
	fmt.Println("        Show version information")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println("")
	fmt.Println("Configuration:")
	fmt.Println("  The application looks for configuration files in the following order:")
	fmt.Println("  1. File specified with -config flag")
	fmt.Println("  2. ./config.yaml")
	fmt.Println("  3. ./configs/config.yaml")
	fmt.Println("  4. /etc/openserial/config.yaml")
	fmt.Println("  5. $HOME/.openserial/config.yaml")
	fmt.Println("")
	fmt.Println("Example configuration file (config.yaml):")
	fmt.Println("  serial:")
	fmt.Println("    port: \"/dev/ttyUSB0\"")
	fmt.Println("    baud_rate: 115200")
	fmt.Println("    data_bits: 8")
	fmt.Println("    stop_bits: 1")
	fmt.Println("    parity: \"none\"")
	fmt.Println("    flow_control: \"none\"")
	fmt.Println("  network:")
	fmt.Println("    listen_port: 8080")
	fmt.Println("    bind_address: \"0.0.0.0\"")
}
