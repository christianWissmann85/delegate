//go:build e2e
// +build e2e

package e2e_test

import (
	"testing"
	"os"
)

// TestMockE2E tests the MCP server with mock providers (no API keys needed)
func TestMockE2E(t *testing.T) {
	// Set environment to use mock provider
	os.Setenv("DELEGATE_MOCK_MODE", "true")
	defer os.Unsetenv("DELEGATE_MOCK_MODE")

	// Start MCP server
	client := StartMCPServer(t)
	defer client.Stop(t)

	// Initialize connection
	t.Log("Initializing MCP connection...")
	initResp, err := client.SendMessage("initialize", map[string]interface{}{
		"clientInfo": map[string]string{
			"name":    "test-client",
			"version": "1.0.0",
		},
	})
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	// List tools to verify server is working
	t.Log("Listing available tools...")
	toolsResp, err := client.SendMessage("tools/list", nil)
	if err != nil {
		t.Fatalf("Failed to list tools: %v", err)
	}

	t.Logf("Server initialized successfully!")
	t.Logf("Init response: %+v", initResp)
	t.Logf("Tools response: %+v", toolsResp)
}