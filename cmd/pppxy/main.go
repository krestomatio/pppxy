package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	pppxyLogger "github.com/krestomatio/pppxy/pkg/logger"
	pppxy "github.com/krestomatio/pppxy/pkg/pppxy"
)

// Configuration variables with default values
var (
	logLevel   = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	configPath = flag.String("config", "/etc/pppxy/config.yaml", "Path to configuration file")
)

func main() {
	flag.Parse()

	// Load the configuration from the specified YAML file
	config, err := pppxy.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Create and start pppxyGroup concurrently
	pppxyGroup := make([]*pppxy.PPPxy, len(config.PPPxyGroup))
	for i, pppxyConfig := range config.PPPxyGroup {
		pppxyID := pppxyConfig.ListenAddr
		logger := pppxyLogger.InitLogger(*logLevel, pppxyID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		newPPPxy := pppxy.NewPPPxy(pppxyConfig, logger)
		pppxyGroup[i] = newPPPxy
		go func(p *pppxy.PPPxy) {
			if err := p.Listen(); err != nil {
				p.Log.Error(err, "Failed to start TCP server")
				os.Exit(1)
			}
			p.Handle()
		}(newPPPxy)
	}

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Close all pppxy listeners
	for _, p := range pppxyGroup {
		if err := p.Close(); err != nil {
			p.Log.Error(err, "Failed to close pppxy listener", "address", p.Config.ListenAddr)
		}
	}
}
