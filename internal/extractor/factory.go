package extractor

// Factory creates extractors with optional configuration
type Factory struct{}

// NewFactory creates a new extractor factory
func NewFactory() *Factory {
	return &Factory{}
}

// Create creates an extractor with the given options
func (f *Factory) Create(languageHint string) *Extractor {
	if languageHint != "" {
		return NewWithHint(languageHint)
	}
	return New()
}

// Default creates a default extractor with no hints
func (f *Factory) Default() *Extractor {
	return New()
}
