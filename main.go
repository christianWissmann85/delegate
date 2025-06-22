package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/christianwissmann85/delegate/internal/config"
	"github.com/christianwissmann85/delegate/internal/logger"
	"github.com/christianwissmann85/delegate/internal/mcp"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// Use basic logger for early errors
		log := logger.New("main", logger.ErrorLevel)
		log.Fatal("Failed to load config", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Create logger with configured level
	log := logger.New("main", logger.ParseLevel(cfg.LogLevel))

	// Create MCP server
	server := mcp.NewServer(cfg)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("Received shutdown signal")
		cancel()
	}()

	// Start server
	log.Info("Starting Delegate", map[string]interface{}{
		"pid": os.Getpid(),
	})

	if err := server.Start(ctx); err != nil {
		log.Fatal("Server error", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
