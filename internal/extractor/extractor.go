package extractor

import (
	"fmt"

	"github.com/christianwissmann85/delegate/internal/handlers"
)

// Extractor extracts code and explanations from LLM responses
type Extractor struct {
	patterns []Pattern
}

// New creates a new extractor
func New() *Extractor {
	return &Extractor{
		patterns: GetPatterns(),
	}
}

// Extract extracts both code and explanation
func (e *Extractor) Extract(content string) (*handlers.Extraction, error) {
	code, err := e.ExtractCode(content)
	if err != nil {
		return nil, fmt.Errorf("extract code: %w", err)
	}

	explanation, err := e.ExtractExplanation(content)
	if err != nil {
		return nil, fmt.Errorf("extract explanation: %w", err)
	}

	return &handlers.Extraction{
		Code:        code,
		Explanation: explanation,
	}, nil
}

// ExtractCode extracts all code blocks
func (e *Extractor) ExtractCode(content string) ([]handlers.CodeBlock, error) {
	var blocks []handlers.CodeBlock
	
	// TODO: Implement code extraction using patterns
	
	return blocks, nil
}

// ExtractExplanation extracts text without code blocks
func (e *Extractor) ExtractExplanation(content string) (string, error) {
	// TODO: Remove code blocks and return remaining text
	return "", fmt.Errorf("not implemented")
}

// removeCodeBlocks removes all code blocks from content
func (e *Extractor) removeCodeBlocks(content string) string {
	// TODO: Implement code block removal
	return content
}