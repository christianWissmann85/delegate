package mock

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/christianwissmann85/delegate/internal/handlers"
)

func TestMockProvider_DefaultResponse(t *testing.T) {
	provider := NewProvider("test-model")

	req := handlers.GenerateRequest{
		Model:  "test-model",
		Prompt: "Hello, world!",
	}

	ch, err := provider.GenerateStream(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}

	var response strings.Builder
	for chunk := range ch {
		if chunk.Error != nil {
			t.Fatalf("Unexpected error: %v", chunk.Error)
		}
		response.WriteString(chunk.Content)
	}

	result := response.String()
	if !strings.Contains(result, "mock response") {
		t.Errorf("Expected mock response, got: %s", result)
	}
	if !strings.Contains(result, "Hello, world!") {
		t.Errorf("Expected prompt echo, got: %s", result)
	}
}

func TestMockProvider_CodeResponse(t *testing.T) {
	provider := NewProvider("test-model")

	req := handlers.GenerateRequest{
		Model:  "test-model",
		Prompt: "Write some code for me",
	}

	ch, err := provider.GenerateStream(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}

	var response strings.Builder
	for chunk := range ch {
		if chunk.Error != nil {
			t.Fatalf("Unexpected error: %v", chunk.Error)
		}
		response.WriteString(chunk.Content)
	}

	result := response.String()
	if !strings.Contains(result, "```python") {
		t.Errorf("Expected code block, got: %s", result)
	}
	if !strings.Contains(result, "def hello_world()") {
		t.Errorf("Expected function definition, got: %s", result)
	}
}

func TestMockProvider_CustomResponses(t *testing.T) {
	provider := NewProvider("test-model").
		WithResponses("First chunk", " Second chunk", " Third chunk")

	req := handlers.GenerateRequest{
		Model:  "test-model",
		Prompt: "Test",
	}

	ch, err := provider.GenerateStream(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}

	var response strings.Builder
	chunkCount := 0
	for chunk := range ch {
		if chunk.Error != nil {
			t.Fatalf("Unexpected error: %v", chunk.Error)
		}
		response.WriteString(chunk.Content)
		chunkCount++
	}

	if chunkCount != 3 {
		t.Errorf("Expected 3 chunks, got %d", chunkCount)
	}

	result := response.String()
	if result != "First chunk Second chunk Third chunk" {
		t.Errorf("Expected custom response, got: %s", result)
	}
}

func TestMockProvider_WithError(t *testing.T) {
	provider := NewProvider("test-model").
		WithResponses("Chunk 1", "Chunk 2", "Chunk 3").
		WithError(2) // Error on second chunk

	req := handlers.GenerateRequest{
		Model:  "test-model",
		Prompt: "Test",
	}

	ch, err := provider.GenerateStream(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}

	chunkCount := 0
	errorReceived := false
	for chunk := range ch {
		chunkCount++
		if chunk.Error != nil {
			errorReceived = true
			if !strings.Contains(chunk.Error.Error(), "mock error on chunk 2") {
				t.Errorf("Expected mock error, got: %v", chunk.Error)
			}
			break
		}
	}

	if !errorReceived {
		t.Error("Expected error but none received")
	}
	if chunkCount != 2 {
		t.Errorf("Expected 2 chunks before error, got %d", chunkCount)
	}
}

func TestMockProvider_ContextCancellation(t *testing.T) {
	provider := NewProvider("test-model").
		WithResponses("Chunk 1", "Chunk 2", "Chunk 3").
		WithDelay(50 * time.Millisecond)

	req := handlers.GenerateRequest{
		Model:  "test-model",
		Prompt: "Test",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is always called

	ch, err := provider.GenerateStream(ctx, req)
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}

	// Cancel after first chunk
	chunkCount := 0
	for chunk := range ch {
		chunkCount++
		if chunkCount == 1 {
			cancel()
		}
		if chunk.Error != nil {
			// Should get context cancelled error
			if chunk.Error != context.Canceled {
				t.Errorf("Expected context.Canceled, got: %v", chunk.Error)
			}
			break
		}
	}

	if chunkCount > 2 {
		t.Errorf("Expected at most 2 chunks before cancellation, got %d", chunkCount)
	}
}

func TestMockProvider_ModelMismatch(t *testing.T) {
	provider := NewProvider("test-model")

	req := handlers.GenerateRequest{
		Model:  "wrong-model",
		Prompt: "Test",
	}

	_, err := provider.GenerateStream(context.Background(), req)
	if err == nil {
		t.Error("Expected error for model mismatch")
	}
	if !strings.Contains(err.Error(), "model mismatch") {
		t.Errorf("Expected model mismatch error, got: %v", err)
	}
}
