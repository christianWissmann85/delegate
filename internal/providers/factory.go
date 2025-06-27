package providers

import (
	"fmt"
	"strings"

	"github.com/christianwissmann85/delegate/internal/config"
	"github.com/christianwissmann85/delegate/internal/handlers"
	"github.com/christianwissmann85/delegate/internal/models"
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
	var provider handlers.Provider
	var providerName string

	// Support mock providers for testing
	if strings.HasPrefix(model, "mock-") {
		provider = mock.NewProvider(model)
		providerName = "mock"
	} else {
		switch model {
		case "gemini-2.5-flash", "gemini-2.5-pro":
			provider = google.NewProvider(f.config.GoogleKey, model, f.config.TimeoutSeconds)
			providerName = "google"
		case "claude-sonnet-4-20250514", "claude-opus-4-20250514":
			provider = anthropic.NewProvider(f.config.AnthropicKey, model, f.config.TimeoutSeconds)
			providerName = "anthropic"
		default:
			return nil, models.NewDelegateError(
				models.ErrorTypeInvalidRequest,
				fmt.Sprintf("Unsupported model: %s. Supported models are: %s", 
					model, strings.Join(f.SupportedModels(), ", ")),
			)
		}
	}

	// Wrap all providers with retry logic
	return NewRetryableProvider(provider, providerName), nil
}

// SupportedModels returns the list of supported models
func (f *Factory) SupportedModels() []string {
	return []string{
		"gemini-2.5-flash",
		"gemini-2.5-pro",
		"claude-sonnet-4-20250514",
		"claude-opus-4-20250514",
	}
}
