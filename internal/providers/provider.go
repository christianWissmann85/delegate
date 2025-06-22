package providers

import (
	"context"

	"github.com/christianwissmann85/delegate/internal/handlers"
)

// Provider generates content from an LLM
type Provider interface {
	GenerateStream(ctx context.Context, req handlers.GenerateRequest) (<-chan handlers.StreamChunk, error)
}

// ProviderCapabilities describes what a provider can do
type ProviderCapabilities struct {
	MaxTokens        int
	SupportsFiles    bool
	StreamingSupport bool
}
