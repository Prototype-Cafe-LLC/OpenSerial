package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/config"
	"github.com/Prototype-Cafe-LLC/OpenSerial/internal/tcpbridge"
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
		fmt.Printf("OpenSerial TCP Bridge v%s\n", version)
		os.Exit(0)
	}

	if *showHelp {
		showUsage()
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.LoadTCPBridgeConfig(*configPath)
	if err != nil {
		logger.Default.Fatal("Failed to load configuration: %v", err)
	}

	// Create TCP bridge server
	server := tcpbridge.NewServer(cfg)

	// Start server
	logger.Default.Info("Starting OpenSerial TCP Bridge server...")
	logger.Default.Info("Listening on: %s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Default.Info("Target server: %s:%d", cfg.Target.Host, cfg.Target.Port)

	if err := server.Start(); err != nil {
		logger.Default.Fatal("Failed to start server: %v", err)
	}

	logger.Default.Info("TCP Bridge server started successfully")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	logger.Default.Info("Shutdown signal received, stopping server...")

	// Stop server
	if err := server.Stop(); err != nil {
		logger.Default.Error("Error stopping server: %v", err)
		os.Exit(1)
	}

	logger.Default.Info("TCP Bridge server stopped successfully")
}

func showUsage() {
	fmt.Printf("OpenSerial TCP Bridge v%s - TCP Relay Server\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  tcpbridge [options]")
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
	fmt.Println("  2. ./tcpbridge.yaml")
	fmt.Println("  3. ./configs/tcpbridge.yaml")
	fmt.Println("  4. /etc/openserial/tcpbridge.yaml")
	fmt.Println("  5. $HOME/.openserial/tcpbridge.yaml")
	fmt.Println("")
	fmt.Println("Example configuration file (tcpbridge.yaml):")
	fmt.Println("  server:")
	fmt.Println("    port: 8080")
	fmt.Println("    host: \"0.0.0.0\"")
	fmt.Println("  target:")
	fmt.Println("    host: \"localhost\"")
	fmt.Println("    port: 8081")
	fmt.Println("  clients:")
	fmt.Println("    max_connections: 10")
	fmt.Println("    connection_timeout: \"5m\"")
}
