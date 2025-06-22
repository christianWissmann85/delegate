package providers

import (
	"context"
	"strings"
	"testing"

	"github.com/christianwissmann85/delegate/internal/config"
	"github.com/christianwissmann85/delegate/internal/handlers"
)

func TestFactory_MockProvider(t *testing.T) {
	cfg := &config.Config{}
	factory := NewFactory(cfg)

	// Test getting mock provider
	provider, err := factory.GetProvider("mock-test")
	if err != nil {
		t.Fatalf("Failed to get mock provider: %v", err)
	}

	if provider == nil {
		t.Fatal("Expected provider, got nil")
	}

	// Test streaming
	req := handlers.GenerateRequest{
		Model:  "mock-test",
		Prompt: "Write some code",
	}

	stream, err := provider.GenerateStream(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to start stream: %v", err)
	}

	var response strings.Builder
	for chunk := range stream {
		if chunk.Error != nil {
			t.Fatalf("Stream error: %v", chunk.Error)
		}
		response.WriteString(chunk.Content)
	}

	result := response.String()
	if result == "" {
		t.Error("Expected response content")
	}

	// Should contain code since prompt mentions "code"
	if !strings.Contains(result, "```") {
		t.Error("Expected code block in response")
	}
}

func TestFactory_RetryLogic(t *testing.T) {
	cfg := &config.Config{}
	factory := NewFactory(cfg)

	// Get a provider (will be wrapped with retry)
	provider, err := factory.GetProvider("mock-retry")
	if err != nil {
		t.Fatalf("Failed to get provider: %v", err)
	}

	// Verify it's a RetryableProvider
	_, ok := provider.(*RetryableProvider)
	if !ok {
		t.Error("Expected provider to be wrapped with RetryableProvider")
	}
}

func TestFactory_UnsupportedModel(t *testing.T) {
	cfg := &config.Config{}
	factory := NewFactory(cfg)

	_, err := factory.GetProvider("unsupported-model")
	if err == nil {
		t.Error("Expected error for unsupported model")
	}

	if !strings.Contains(err.Error(), "unsupported model") {
		t.Errorf("Expected 'unsupported model' error, got: %v", err)
	}
}

func TestProviderIntegration_GeminiModels(t *testing.T) {
	cfg := &config.Config{
		GoogleKey: "test-key",
	}
	factory := NewFactory(cfg)

	geminiModels := []string{"gemini-2.5-flash", "gemini-2.5-pro"}

	for _, model := range geminiModels {
		t.Run(model, func(t *testing.T) {
			provider, err := factory.GetProvider(model)
			if err != nil {
				t.Fatalf("Failed to get provider for %s: %v", model, err)
			}

			if provider == nil {
				t.Errorf("Expected provider for %s, got nil", model)
			}
		})
	}
}

func TestProviderIntegration_ClaudeModels(t *testing.T) {
	cfg := &config.Config{
		AnthropicKey: "test-key",
	}
	factory := NewFactory(cfg)

	claudeModels := []string{"claude-sonnet-4-20250514", "claude-opus-4-20250514"}

	for _, model := range claudeModels {
		t.Run(model, func(t *testing.T) {
			provider, err := factory.GetProvider(model)
			if err != nil {
				t.Fatalf("Failed to get provider for %s: %v", model, err)
			}

			if provider == nil {
				t.Errorf("Expected provider for %s, got nil", model)
			}
		})
	}
}
