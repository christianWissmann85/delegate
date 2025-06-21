package anthropic

import (
	"context"
	"fmt"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/christianwissmann85/delegate/internal/handlers"
	"github.com/christianwissmann85/delegate/internal/logger"
)

// Provider implements the Anthropic provider
type Provider struct {
	apiKey string
	model  string
	logger *logger.Logger
}

// NewProvider creates a new Anthropic provider
func NewProvider(apiKey, model string) *Provider {
	return &Provider{
		apiKey: apiKey,
		model:  model,
		logger: logger.New("providers.anthropic", logger.InfoLevel),
	}
}

// GenerateStream generates content using Anthropic's API
func (p *Provider) GenerateStream(ctx context.Context, req handlers.GenerateRequest) (<-chan handlers.StreamChunk, error) {
	// Create output channel
	ch := make(chan handlers.StreamChunk)
	
	// Start streaming in a goroutine
	go func() {
		defer close(ch)
		
		// Set timeout from request or use default
		timeout := 60 * time.Second
		if req.Timeout > 0 {
			timeout = time.Duration(req.Timeout) * time.Second
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		
		// Create Anthropic client
		client := anthropic.NewClient(
			option.WithAPIKey(p.apiKey),
		)
		
		// Build prompt with files
		promptText := req.Prompt
		for _, filePath := range req.Files {
			// For now, we'll just add file paths as text
			// In a full implementation, we'd read and process the files
			promptText += fmt.Sprintf("\n\nFile: %s\n<file content would go here>", filePath)
		}
		
		// Configure message parameters
		params := anthropic.MessageNewParams{
			Model:     anthropic.Model(p.model),
			MaxTokens: 4096, // Default max
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock(promptText)),
			},
		}
		
		// Override max tokens if specified
		if req.MaxTokens > 0 {
			params.MaxTokens = int64(req.MaxTokens)
		}
		
		// Create streaming request
		stream := client.Messages.NewStreaming(ctx, params)
		
		// Process stream events
		for stream.Next() {
			event := stream.Current()
			
			switch eventVariant := event.AsAny().(type) {
			case anthropic.ContentBlockDeltaEvent:
				switch deltaVariant := eventVariant.Delta.AsAny().(type) {
				case anthropic.TextDelta:
					ch <- handlers.StreamChunk{Content: deltaVariant.Text}
				}
			case anthropic.MessageStopEvent:
				// Stream completed successfully
				p.logger.Info("Streaming completed", map[string]interface{}{
					"model": p.model,
				})
			}
		}
		
		// Check for stream error
		if err := stream.Err(); err != nil {
			ch <- handlers.StreamChunk{Error: fmt.Errorf("stream error: %w", err)}
			return
		}
	}()
	
	return ch, nil
}

