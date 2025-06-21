package google

import (
	"context"
	"fmt"
	"time"

	"github.com/christianwissmann85/delegate/internal/handlers"
	"github.com/christianwissmann85/delegate/internal/logger"
	"google.golang.org/genai"
)

// Provider implements the Google (Gemini) provider
type Provider struct {
	apiKey string
	model  string
	logger *logger.Logger
}

// NewProvider creates a new Google provider
func NewProvider(apiKey, model string) *Provider {
	return &Provider{
		apiKey: apiKey,
		model:  model,
		logger: logger.New("providers.google", logger.InfoLevel),
	}
}

// GenerateStream generates content using Google's Gemini API
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
		
		// Create client with API key
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  p.apiKey,
			Backend: genai.BackendGeminiAPI,
		})
		if err != nil {
			ch <- handlers.StreamChunk{Error: fmt.Errorf("create gemini client: %w", err)}
			return
		}
		
		// Build content with prompt and files
		promptText := req.Prompt
		if len(req.Files) > 0 {
			// Read files with memory limits
			fileContents, err := handlers.ReadFilesWithLimit(req.Files)
			if err != nil {
				ch <- handlers.StreamChunk{Error: err}
				return
			}
			promptText = handlers.BuildPromptWithFiles(promptText, fileContents)
		}
		
		// Create content for the request
		content := &genai.Content{
			Parts: []*genai.Part{
				{Text: promptText},
			},
		}
		
		// Configure generation settings
		config := &genai.GenerateContentConfig{
			Temperature: float32Ptr(0.3), // Lower for more deterministic output
			TopP:        float32Ptr(0.95),
		}
		if req.MaxTokens > 0 {
			config.MaxOutputTokens = int32(req.MaxTokens)
		}
		
		// Generate content with streaming
		stream := client.Models.GenerateContentStream(ctx, p.model, []*genai.Content{content}, config)
		
		// Stream the responses
		for result, err := range stream {
			if err != nil {
				ch <- handlers.StreamChunk{Error: fmt.Errorf("stream error: %w", err)}
				return
			}
			
			// Extract text from result
			for _, candidate := range result.Candidates {
				if candidate.Content != nil {
					for _, part := range candidate.Content.Parts {
						if part != nil && part.Text != "" {
							ch <- handlers.StreamChunk{Content: part.Text}
						}
					}
				}
			}
		}
		
		p.logger.Info("Streaming completed", map[string]interface{}{
			"model": p.model,
		})
	}()
	
	return ch, nil
}

// Helper function for pointer creation
func float32Ptr(f float32) *float32 {
	return &f
}