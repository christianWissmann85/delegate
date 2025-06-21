package providers

import (
	"fmt"
	"strings"

	"github.com/christianwissmann85/delegate/internal/config"
	"github.com/christianwissmann85/delegate/internal/handlers"
	"github.com/christianwissmann85/delegate/internal/providers/anthropic"
	"github.com/christianwissmann85/delegate/internal/providers/google"
	"github.com/christianwissmann85/delegate/internal/providers/mock"
)

// Factory creates providers based on model
type Factory struct {
	config *config.Config
}

// NewFactory creates a new provider factory
func NewFactory(cfg *config.Config) *Factory {
	return &Factory{
		config: cfg,
	}
}

// GetProvider returns a provider for the given model
func (f *Factory) GetProvider(model string) (handlers.Provider, error) {
	// Support mock providers for testing
	if strings.HasPrefix(model, "mock-") {
		return mock.NewProvider(model), nil
	}
	
	switch model {
	case "gemini-2.5-flash", "gemini-2.5-pro":
		return google.NewProvider(f.config.GoogleKey, model), nil
	case "claude-sonnet-4-20250514", "claude-opus-4-20250514":
		return anthropic.NewProvider(f.config.AnthropicKey, model), nil
	default:
		return nil, fmt.Errorf("unsupported model: %s", model)
	}
}