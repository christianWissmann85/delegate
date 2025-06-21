package anthropic

import (
	"context"
	"fmt"

	"github.com/christianwissmann85/delegate/internal/handlers"
)

// Provider implements the Anthropic provider
type Provider struct {
	apiKey string
	model  string
}

// NewProvider creates a new Anthropic provider
func NewProvider(apiKey, model string) *Provider {
	return &Provider{
		apiKey: apiKey,
		model:  model,
	}
}

// GenerateStream generates content using Anthropic's API
func (p *Provider) GenerateStream(ctx context.Context, req handlers.GenerateRequest) (<-chan handlers.StreamChunk, error) {
	// TODO: Implement Anthropic API streaming
	return nil, fmt.Errorf("not implemented")
}

// buildRequest creates an Anthropic API request
func (p *Provider) buildRequest(req handlers.GenerateRequest) interface{} {
	// TODO: Build Anthropic-specific request
	return nil
}