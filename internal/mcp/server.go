package mcp

import (
	"context"
	"fmt"

	"github.com/christianwissmann85/delegate/internal/config"
	"github.com/christianwissmann85/delegate/internal/handlers"
)

// Server handles MCP protocol communication
type Server struct {
	config   *config.Config
	protocol *Protocol
	tools    map[string]Tool
}

// NewServer creates a new MCP server
func NewServer(cfg *config.Config) *Server {
	s := &Server{
		config: cfg,
		tools:  make(map[string]Tool),
	}

	// Initialize protocol handler
	s.protocol = NewProtocol(s)

	// Register tools
	s.registerTools()

	return s
}

// Start begins serving MCP requests
func (s *Server) Start(ctx context.Context) error {
	// TODO: Implement server start logic
	return fmt.Errorf("not implemented")
}

// RegisterTool registers a tool with the server
func (s *Server) RegisterTool(tool Tool) error {
	if _, exists := s.tools[tool.Name()]; exists {
		return fmt.Errorf("tool %s already registered", tool.Name())
	}
	s.tools[tool.Name()] = tool
	return nil
}

// registerTools sets up all available tools
func (s *Server) registerTools() {
	// TODO: Initialize handlers and register tools
}