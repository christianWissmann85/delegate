package mcp

import (
	"context"
	"fmt"

	"github.com/christianwissmann85/delegate/internal/config"
	"github.com/christianwissmann85/delegate/internal/extractor"
	"github.com/christianwissmann85/delegate/internal/handlers"
	"github.com/christianwissmann85/delegate/internal/logger"
	"github.com/christianwissmann85/delegate/internal/providers"
	"github.com/christianwissmann85/delegate/internal/storage"
)

// Server handles MCP protocol communication
type Server struct {
	config   *config.Config
	protocol *Protocol
	tools    map[string]Tool
	logger   *logger.Logger
	storage  handlers.Storage
}

// NewServer creates a new MCP server
func NewServer(cfg *config.Config) *Server {
	logLevel := logger.ParseLevel(cfg.LogLevel)
	
	s := &Server{
		config: cfg,
		tools:  make(map[string]Tool),
		logger: logger.New("mcp.server", logLevel),
	}

	// Initialize protocol handler
	s.protocol = NewProtocol(s, logLevel)

	// Register tools
	if err := s.registerTools(); err != nil {
		s.logger.Fatal("Failed to register tools", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return s
}

// Start begins serving MCP requests
func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("Starting Delegate MCP server", map[string]interface{}{
		"version":          "1.0.0",
		"supported_models": s.config.SupportedModels(),
		"output_dir":       s.config.OutputDir,
		"timeout_seconds":  s.config.TimeoutSeconds,
	})
	
	// Start cleanup routine
	if store, ok := s.storage.(*storage.FileStore); ok {
		interval, maxAge := storage.DefaultCleanupConfig()
		cleaner := storage.NewCleaner(store, interval, maxAge)
		go cleaner.Start(ctx)
	}
	
	// Start protocol handler
	return s.protocol.HandleMessages(ctx)
}

// RegisterTool registers a tool with the server
func (s *Server) RegisterTool(tool Tool) error {
	if _, exists := s.tools[tool.Name()]; exists {
		return fmt.Errorf("tool %s already registered", tool.Name())
	}
	s.tools[tool.Name()] = tool
	s.logger.Info("Registered tool", map[string]interface{}{
		"tool": tool.Name(),
	})
	return nil
}

// registerTools sets up all available tools
func (s *Server) registerTools() error {
	// Initialize storage
	store, err := storage.NewFileStore(s.config.OutputDir)
	if err != nil {
		return fmt.Errorf("create storage: %w", err)
	}
	s.storage = store

	// Initialize provider factory
	providerFactory := providers.NewFactory(s.config)

	// Initialize extractor factory
	extractFactory := extractor.NewFactory()

	// Create handlers
	invokeHandler := handlers.NewInvokeHandler(providerFactory, store, extractFactory)
	checkHandler := handlers.NewCheckHandler(store)
	readHandler := handlers.NewReadHandler(store)

	// Register tools
	tools := []Tool{
		&InvokeTool{handler: invokeHandler},
		&CheckTool{handler: checkHandler},
		&ReadTool{handler: readHandler},
	}

	for _, tool := range tools {
		if err := s.RegisterTool(tool); err != nil {
			return err
		}
	}

	return nil
}