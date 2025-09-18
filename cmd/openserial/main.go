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
		clientMode  = flag.Bool("client", false, "Run in client mode (connect to server)")
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

	// Create bridge based on mode
	if *clientMode {
		// Client mode - connect to server
		clientBridge := bridge.NewClientBridge(cfg)

		logger.Default.Info("Starting OpenSerial client bridge...")
		logger.Default.Info("Serial port: %s", cfg.Serial.Port)
		logger.Default.Info("Connecting to: %s:%d", cfg.Network.BindAddress, cfg.Network.ListenPort)

		if err := clientBridge.Start(); err != nil {
			logger.Default.Fatal("Failed to start client bridge: %v", err)
		}

		logger.Default.Info("Client bridge started successfully")

		// Set up signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Wait for shutdown signal
		<-sigChan
		logger.Default.Info("Shutdown signal received, stopping client bridge...")

		// Stop client bridge
		if err := clientBridge.Stop(); err != nil {
			logger.Default.Error("Error stopping client bridge: %v", err)
			os.Exit(1)
		}

		logger.Default.Info("Client bridge stopped successfully")
		return
	} else {
		// Server mode - listen for connections
		serverBridge := bridge.NewBridge(cfg)

		logger.Default.Info("Starting OpenSerial server bridge...")
		logger.Default.Info("Serial port: %s", cfg.Serial.Port)
		logger.Default.Info("Listening on: %s:%d", cfg.Network.BindAddress, cfg.Network.ListenPort)

		if err := serverBridge.Start(); err != nil {
			logger.Default.Fatal("Failed to start server bridge: %v", err)
		}

		logger.Default.Info("Server bridge started successfully")

		// Set up signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Wait for shutdown signal
		<-sigChan
		logger.Default.Info("Shutdown signal received, stopping server bridge...")

		// Stop server bridge
		if err := serverBridge.Stop(); err != nil {
			logger.Default.Error("Error stopping server bridge: %v", err)
			os.Exit(1)
		}

		logger.Default.Info("Server bridge stopped successfully")
		return
	}
}

func showUsage() {
	fmt.Printf("OpenSerial v%s - UART-to-TCP Bridge\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  openserial [options]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -config string")
	fmt.Println("        Path to configuration file (YAML format)")
	fmt.Println("  -client")
	fmt.Println("        Run in client mode (connect to server)")
	fmt.Println("  -version")
	fmt.Println("        Show version information")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println("")
	fmt.Println("Modes:")
	fmt.Println("  Server mode (default):")
	fmt.Println("    - Listens for incoming TCP connections")
	fmt.Println("    - Bridges serial port to TCP clients")
	fmt.Println("    - Use on Mac with admin privileges")
	fmt.Println("")
	fmt.Println("  Client mode (-client flag):")
	fmt.Println("    - Connects to a TCP server")
	fmt.Println("    - Bridges serial port to server connection")
	fmt.Println("    - Use on Windows without admin privileges")
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
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  # Server mode (Mac)")
	fmt.Println("  ./openserial -config server.yaml")
	fmt.Println("")
	fmt.Println("  # Client mode (Windows)")
	fmt.Println("  ./openserial -client -config client.yaml")
}
