package handlers

import (
	"context"
	
	"github.com/christianwissmann85/delegate/internal/models"
	"github.com/christianwissmann85/delegate/internal/extractor"
)

// Provider generates content using an LLM
type Provider interface {
	Generate(ctx context.Context, req *models.GenerateRequest) (<-chan models.StreamChunk, error)
}

// ProviderFactory creates Provider instances
type ProviderFactory interface {
	GetProvider(model string) (Provider, error)
	SupportedModels() []string
}

// Storage manages generated outputs
type Storage interface {
	GetOutput(id string) (*models.Output, error)
	SaveOutput(output *models.Output) error
}

// Extractor parses raw LLM output to extract code and explanations
type Extractor interface {
	Extract(raw string, codeOnly bool) (models.Extraction, error)
}

// ExtractorFactory creates Extractor instances
type ExtractorFactory interface {
	Default() Extractor
}

// ExtractorAdapter adapts the concrete extractor.Extractor to the handlers.Extractor interface
type ExtractorAdapter struct {
	e *extractor.Extractor
}

// NewExtractorAdapter creates a new ExtractorAdapter
func NewExtractorAdapter(e *extractor.Extractor) Extractor {
	return &ExtractorAdapter{e: e}
}

// Extract implements the Extractor interface
func (a *ExtractorAdapter) Extract(raw string, codeOnly bool) (models.Extraction, error) {
	extractionPtr, err := a.e.Extract(raw, codeOnly)
	if err != nil {
		return models.Extraction{}, err
	}
	if extractionPtr == nil {
		return models.Extraction{}, nil
	}
	return *extractionPtr, nil
}

// ExtractorFactoryAdapter adapts the concrete extractor.Factory to the handlers.ExtractorFactory interface  
type ExtractorFactoryAdapter struct {
	f *extractor.Factory
}

// NewExtractorFactoryAdapter creates a new ExtractorFactoryAdapter
func NewExtractorFactoryAdapter(f *extractor.Factory) ExtractorFactory {
	return &ExtractorFactoryAdapter{f: f}
}

// Default implements the ExtractorFactory interface
func (a *ExtractorFactoryAdapter) Default() Extractor {
	return NewExtractorAdapter(a.f.Default())
}