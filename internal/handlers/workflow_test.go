package handlers_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/christianwissmann85/delegate/internal/extractor"
	"github.com/christianwissmann85/delegate/internal/handlers"
	"github.com/christianwissmann85/delegate/internal/providers/mock"
	"github.com/christianwissmann85/delegate/internal/storage"
)

// mockProviderFactory implements handlers.ProviderFactory
type mockProviderFactory struct{}

func (f *mockProviderFactory) GetProvider(model string) (handlers.Provider, error) {
	return mock.NewProvider(model), nil
}

func TestFullWorkflow_InvokeCheckRead(t *testing.T) {
	// Setup storage
	tempDir := t.TempDir()
	store, err := storage.NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Setup extractor factory
	extractFactory := extractor.NewFactory()

	// Setup provider factory
	providerFactory := &mockProviderFactory{}

	// Create handlers
	invokeHandler := handlers.NewInvokeHandler(providerFactory, store, extractFactory)
	checkHandler := handlers.NewCheckHandler(store)
	readHandler := handlers.NewReadHandler(store)

	ctx := context.Background()

	// Step 1: Invoke to generate content
	var outputID string
	t.Run("invoke", func(t *testing.T) {
		req := handlers.InvokeRequest{
			Model:  "mock-test",
			Prompt: "write code for hello world function in Python",
		}

		resp, err := invokeHandler.Handle(ctx, req)
		if err != nil {
			t.Fatalf("Invoke failed: %v", err)
		}

		if resp.OutputID == "" {
			t.Error("Expected output ID")
		}

		// Store output ID for next steps
		outputID = resp.OutputID
	})

	if outputID == "" {
		t.Fatal("No output ID from invoke step")
	}

	// Wait a moment for file write
	time.Sleep(10 * time.Millisecond)

	// Step 2: Check metadata
	t.Run("check", func(t *testing.T) {
		req := handlers.CheckRequest{
			OutputID: outputID,
		}

		resp, err := checkHandler.Handle(ctx, req)
		if err != nil {
			t.Fatalf("Check failed: %v", err)
		}

		// Verify metadata
		if resp.ID != outputID {
			t.Errorf("Expected ID %s, got %s", outputID, resp.ID)
		}

		if resp.Model != "mock-test" {
			t.Errorf("Expected model mock-test, got %s", resp.Model)
		}

		if !resp.HasCode {
			t.Error("Expected HasCode to be true")
		}

		if !resp.HasExplanation {
			t.Error("Expected HasExplanation to be true")
		}

		if resp.CodeBlocksCount == 0 {
			t.Error("Expected at least one code block")
		}

		if resp.EstimatedTokens == 0 {
			t.Error("Expected non-zero token estimate")
		}
	})

	// Step 3: Read different extraction modes
	t.Run("read_all", func(t *testing.T) {
		req := handlers.ReadRequest{
			OutputID: outputID,
			Options: handlers.ReadOptions{
				Extract: "all",
			},
		}

		resp, err := readHandler.Handle(ctx, req)
		if err != nil {
			t.Fatalf("Read all failed: %v", err)
		}

		if resp.Content == "" {
			t.Error("Expected content")
		}

		// Should contain both code and explanation
		if !strings.Contains(resp.Content, "```") {
			t.Error("Expected code blocks in 'all' mode")
		}
	})

	t.Run("read_code_only", func(t *testing.T) {
		req := handlers.ReadRequest{
			OutputID: outputID,
			Options: handlers.ReadOptions{
				Extract: "code",
			},
		}

		resp, err := readHandler.Handle(ctx, req)
		if err != nil {
			t.Fatalf("Read code failed: %v", err)
		}

		// Should only contain code blocks
		if !strings.HasPrefix(strings.TrimSpace(resp.Content), "```") {
			t.Error("Expected content to start with code block")
		}

		// Count code blocks
		codeBlocks := strings.Count(resp.Content, "```") / 2
		if codeBlocks == 0 {
			t.Error("Expected at least one code block")
		}
	})

	t.Run("read_explanation_only", func(t *testing.T) {
		req := handlers.ReadRequest{
			OutputID: outputID,
			Options: handlers.ReadOptions{
				Extract: "explanation",
			},
		}

		resp, err := readHandler.Handle(ctx, req)
		if err != nil {
			t.Fatalf("Read explanation failed: %v", err)
		}

		// Should not contain code blocks
		if strings.Contains(resp.Content, "```") {
			t.Error("Expected no code blocks in explanation mode")
		}
	})

	t.Run("read_with_truncation", func(t *testing.T) {
		req := handlers.ReadRequest{
			OutputID: outputID,
			Options: handlers.ReadOptions{
				Extract:   "all",
				MaxTokens: 10, // Very small limit
			},
		}

		resp, err := readHandler.Handle(ctx, req)
		if err != nil {
			t.Fatalf("Read with truncation failed: %v", err)
		}

		// Should be truncated
		if !strings.HasSuffix(resp.Content, "...") {
			t.Error("Expected truncated content to end with ...")
		}

		// Rough check that it's actually truncated
		if len(resp.Content) > 50 { // 10 tokens * 4 chars + some buffer
			t.Errorf("Content seems too long for 10 token limit: %d chars", len(resp.Content))
		}
	})
}

func TestCheckHandler_Errors(t *testing.T) {
	store, _ := storage.NewFileStore(t.TempDir())
	handler := handlers.NewCheckHandler(store)
	ctx := context.Background()

	t.Run("missing_output_id", func(t *testing.T) {
		req := handlers.CheckRequest{}
		_, err := handler.Handle(ctx, req)
		if err == nil || !strings.Contains(err.Error(), "output_id is required") {
			t.Errorf("Expected output_id required error, got: %v", err)
		}
	})

	t.Run("nonexistent_output", func(t *testing.T) {
		req := handlers.CheckRequest{
			OutputID: "nonexistent",
		}
		_, err := handler.Handle(ctx, req)
		if err == nil {
			t.Error("Expected error for nonexistent output")
		}
	})
}

func TestReadHandler_Errors(t *testing.T) {
	store, _ := storage.NewFileStore(t.TempDir())
	handler := handlers.NewReadHandler(store)
	ctx := context.Background()

	t.Run("missing_output_id", func(t *testing.T) {
		req := handlers.ReadRequest{}
		_, err := handler.Handle(ctx, req)
		if err == nil || !strings.Contains(err.Error(), "output_id is required") {
			t.Errorf("Expected output_id required error, got: %v", err)
		}
	})

	t.Run("invalid_extract_option", func(t *testing.T) {
		req := handlers.ReadRequest{
			OutputID: "test",
			Options: handlers.ReadOptions{
				Extract: "invalid",
			},
		}
		_, err := handler.Handle(ctx, req)
		if err == nil || !strings.Contains(err.Error(), "invalid extract option") {
			t.Errorf("Expected invalid extract option error, got: %v", err)
		}
	})
}

func TestTokenEstimation(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		minRatio float64 // minimum chars/token ratio
		maxRatio float64 // maximum chars/token ratio
	}{
		{
			name:     "simple text",
			text:     "Hello, this is a simple test.",
			minRatio: 3.0,
			maxRatio: 5.0,
		},
		{
			name: "code heavy",
			text: `function test() {
				const x = 10;
				return x * 2;
			}`,
			minRatio: 2.5,
			maxRatio: 5.0, // Allow higher ratio for short code
		},
		{
			name:     "json",
			text:     `{"name": "test", "value": 123, "nested": {"key": "value"}}`,
			minRatio: 2.0,
			maxRatio: 6.0, // Allow higher ratio for compact JSON
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := handlers.EstimateTokens(tt.text)
			ratio := float64(len(tt.text)) / float64(tokens)
			
			if ratio < tt.minRatio || ratio > tt.maxRatio {
				t.Errorf("Token estimation out of expected range. Text: %d chars, Tokens: %d, Ratio: %.2f (expected %.1f-%.1f)",
					len(tt.text), tokens, ratio, tt.minRatio, tt.maxRatio)
			}
		})
	}
}