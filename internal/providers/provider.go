package providers

import (
	"context"

	"github.com/christianwissmann85/delegate/internal/models"
)

// Provider generates content from an LLM
type Provider interface {
	GenerateStream(ctx context.Context, req models.GenerateRequest) (<-chan models.StreamChunk, error)
}

// ProviderCapabilities describes what a provider can do
type ProviderCapabilities struct {
	MaxTokens        int
	SupportsFiles    bool
	StreamingSupport bool
}
