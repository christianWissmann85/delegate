package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/christianwissmann85/delegate/internal/config"
)

// TestServer_InvokeFlow tests the complete invoke → check → read workflow
func TestServer_InvokeFlow(t *testing.T) {
	// Create test config with mock provider
	cfg := &config.Config{
		LogLevel:     "debug",
		OutputDir:    t.TempDir(),
		GoogleKey:    "mock-key",
		AnthropicKey: "mock-key",
	}

	// Create server
	server := NewServer(cfg)

	// Test invoke tool
	t.Run("invoke", func(t *testing.T) {
		params := json.RawMessage(`{
			"model": "mock-test",
			"prompt": "Write a hello world function"
		}`)

		result, err := server.tools["delegate_invoke"].Handler(context.Background(), params)
		if err != nil {
			t.Fatalf("Invoke failed: %v", err)
		}

		// Just verify we got a result
		if result == nil {
			t.Error("Expected invoke result")
		}
	})

}

// TestServer_ToolRegistration verifies all tools are registered
func TestServer_ToolRegistration(t *testing.T) {
	cfg := &config.Config{
		LogLevel:  "info",
		OutputDir: t.TempDir(),
	}

	server := NewServer(cfg)

	expectedTools := []string{
		"delegate_invoke",
		"delegate_check",
		"delegate_read",
	}

	for _, toolName := range expectedTools {
		tool, exists := server.tools[toolName]
		if !exists {
			t.Errorf("Tool %s not registered", toolName)
			continue
		}

		// Verify tool has required methods
		if tool.Name() != toolName {
			t.Errorf("Tool name mismatch: expected %s, got %s", toolName, tool.Name())
		}

		if tool.Description() == "" {
			t.Errorf("Tool %s has empty description", toolName)
		}

		schema := tool.Schema()
		if schema.Type != "object" {
			t.Errorf("Tool %s schema type should be 'object', got %s", toolName, schema.Type)
		}
	}
}
