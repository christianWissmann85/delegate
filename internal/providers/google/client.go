package google

import (
	"context"
	"fmt"

	"github.com/christianwissmann85/delegate/internal/handlers"
)

// Provider implements the Google (Gemini) provider
type Provider struct {
	apiKey string
	model  string
}

// NewProvider creates a new Google provider
func NewProvider(apiKey, model string) *Provider {
	return &Provider{
		apiKey: apiKey,
		model:  model,
	}
}

// GenerateStream generates content using Google's Gemini API
func (p *Provider) GenerateStream(ctx context.Context, req handlers.GenerateRequest) (<-chan handlers.StreamChunk, error) {
	// TODO: Implement Gemini API streaming
	return nil, fmt.Errorf("not implemented")
}

// buildRequest creates a Gemini API request
func (p *Provider) buildRequest(req handlers.GenerateRequest) interface{} {
	// TODO: Build Gemini-specific request
	return nil
}