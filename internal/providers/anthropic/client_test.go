package anthropic

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/christianwissmann85/delegate/internal/handlers"
)

func TestProvider_GenerateStream(t *testing.T) {
	// Skip if no API key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set")
	}

	provider := NewProvider(apiKey, "claude-sonnet-4-20250514")

	ctx := context.Background()
	req := handlers.GenerateRequest{
		Model:     "claude-sonnet-4-20250514",
		Prompt:    "Write a simple hello world function in Python. Keep it very short.",
		MaxTokens: 100,
		Timeout:   30,
	}

	stream, err := provider.GenerateStream(ctx, req)
	if err != nil {
		t.Fatalf("Failed to start stream: %v", err)
	}

	var response strings.Builder
	chunkCount := 0

	for chunk := range stream {
		if chunk.Error != nil {
			t.Fatalf("Stream error: %v", chunk.Error)
		}
		response.WriteString(chunk.Content)
		chunkCount++
	}

	result := response.String()

	// Verify we got a response
	if result == "" {
		t.Error("Got empty response")
	}

	// Verify we got multiple chunks (streaming worked)
	if chunkCount < 2 {
		t.Errorf("Expected multiple chunks for streaming, got %d", chunkCount)
	}

	// Verify content makes sense
	if !strings.Contains(strings.ToLower(result), "hello") {
		t.Errorf("Expected 'hello' in response, got: %s", result)
	}

	t.Logf("Got %d chunks, total response: %s", chunkCount, result)
}

func TestProvider_Timeout(t *testing.T) {
	// Skip if no API key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("ANTHROPIC_API_KEY not set")
	}

	provider := NewProvider(apiKey, "claude-sonnet-4-20250514")

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := handlers.GenerateRequest{
		Model:   "claude-sonnet-4-20250514",
		Prompt:  "This should timeout immediately",
		Timeout: 1, // 1 second
	}

	stream, err := provider.GenerateStream(ctx, req)
	if err != nil {
		t.Fatalf("Failed to start stream: %v", err)
	}

	// Should get an error due to cancelled context
	errorReceived := false
	for chunk := range stream {
		if chunk.Error != nil {
			errorReceived = true
			t.Logf("Got expected error: %v", chunk.Error)
			break
		}
	}

	if !errorReceived {
		t.Error("Expected timeout error but got none")
	}
}

func TestProvider_InvalidModel(t *testing.T) {
	provider := NewProvider("fake-api-key", "invalid-model")

	ctx := context.Background()
	req := handlers.GenerateRequest{
		Model:  "invalid-model",
		Prompt: "Test",
	}

	stream, err := provider.GenerateStream(ctx, req)
	if err != nil {
		t.Fatalf("Failed to start stream: %v", err)
	}

	// Should get an error about invalid model or API key
	errorReceived := false
	for chunk := range stream {
		if chunk.Error != nil {
			errorReceived = true
			t.Logf("Got expected error: %v", chunk.Error)
			break
		}
	}

	if !errorReceived {
		t.Error("Expected error for invalid model/key but got none")
	}
}
