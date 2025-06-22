package google

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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

// apiKeyTransport is a custom HTTP transport that adds API key authentication
type apiKeyTransport struct {
	base   http.RoundTripper
	apiKey string
}

// RoundTrip adds the 'x-goog-api-key' header to every request
func (t *apiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original
	newReq := req.Clone(req.Context())
	
	// Remove any existing authorization headers that might interfere
	newReq.Header.Del("Authorization")
	newReq.Header.Del("x-goog-api-client")
	
	// Set our API key header
	newReq.Header.Set("x-goog-api-key", t.apiKey)
	
	return t.base.RoundTrip(newReq)
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

		// Debug logging
		p.logger.Info("Creating Gemini client with custom API key transport", map[string]interface{}{
			"model":       p.model,
			"apiKeyLen":   len(p.apiKey),
			"apiKeyStart": p.apiKey[:10] + "...",
		})

		// Create a custom HTTP client that uses our API key transport
		customHTTPClient := &http.Client{
			Transport: &apiKeyTransport{
				base:   http.DefaultTransport,
				apiKey: p.apiKey,
			},
		}

		// Create client with custom transport instead of API key
		// This bypasses any ADC (Application Default Credentials) detection
		// IMPORTANT: We explicitly set APIKey to empty string to prevent ADC fallback
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			Backend:    genai.BackendGeminiAPI,
			HTTPClient: customHTTPClient,
			APIKey:     "", // Explicitly empty to prevent ADC fallback
			// Authentication is handled via x-goog-api-key header in our custom transport
		})
		if err != nil {
			ch <- handlers.StreamChunk{Error: fmt.Errorf("create gemini client with custom transport: %w", err)}
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
		// Add "models/" prefix if not already present
		modelName := p.model
		if !strings.HasPrefix(modelName, "models/") {
			modelName = "models/" + modelName
		}

		p.logger.Info("Calling GenerateContentStream", map[string]interface{}{
			"originalModel": p.model,
			"modelName":     modelName,
		})

		stream := client.Models.GenerateContentStream(ctx, modelName, []*genai.Content{content}, config)

		// Stream the responses
		for result, err := range stream {
			if err != nil {
				// Check for authentication-related errors
				errMsg := err.Error()
				isAuthError := strings.Contains(errMsg, "API key") || 
					strings.Contains(errMsg, "authentication") || 
					strings.Contains(errMsg, "unauthorized") ||
					strings.Contains(errMsg, "403") ||
					strings.Contains(errMsg, "401")
				
				p.logger.Error("Stream error details", map[string]interface{}{
					"error":        errMsg,
					"type":         fmt.Sprintf("%T", err),
					"model":        modelName,
					"isAuthError":  isAuthError,
					"apiKeyLen":    len(p.apiKey),
					"apiKeyPrefix": p.apiKey[:10] + "...",
				})
				
				// Add hint about potential ADC interference
				if isAuthError {
					ch <- handlers.StreamChunk{Error: fmt.Errorf("authentication error (check if ADC is overriding API key): %w", err)}
				} else {
					ch <- handlers.StreamChunk{Error: fmt.Errorf("stream error: %w", err)}
				}
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
