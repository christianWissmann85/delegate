package extractor

import "github.com/christianwissmann85/delegate/internal/handlers"

// Factory creates extractors with optional configuration
type Factory struct{}

// NewFactory creates a new extractor factory
func NewFactory() *Factory {
	return &Factory{}
}

// Create creates an extractor with the given options
func (f *Factory) Create(languageHint string) handlers.Extractor {
	if languageHint != "" {
		return NewWithHint(languageHint)
	}
	return New()
}

// Default creates a default extractor with no hints
func (f *Factory) Default() handlers.Extractor {
	return New()
}