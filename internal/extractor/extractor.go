package extractor

import (
	"fmt"
	"regexp"

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
	usedRanges := make(map[string]bool)
	
	// Try each pattern
	for _, pattern := range e.patterns {
		if pattern.Name == "FencedCodeBlock" || pattern.Name == "AltFencedBlock" {
			matches := pattern.Regex.FindAllStringSubmatch(content, -1)
			indices := pattern.Regex.FindAllStringIndex(content, -1)
			
			for i, match := range matches {
				if len(match) < 3 || len(indices) <= i {
					continue
				}
				
				start, end := indices[i][0], indices[i][1]
				rangeKey := fmt.Sprintf("%d-%d", start, end)
				
				// Skip if this range was already captured
				if usedRanges[rangeKey] {
					continue
				}
				usedRanges[rangeKey] = true
				
				// Extract language and code
				lang := match[1]
				if lang == "" {
					lang = "plaintext"
				}
				code := match[2]
				
				// Calculate line numbers (approximate)
				linesBefore := countLines(content[:start])
				linesInBlock := countLines(code)
				
				blocks = append(blocks, handlers.CodeBlock{
					Language:  lang,
					Content:   code,
					LineStart: linesBefore + 1,
					LineEnd:   linesBefore + linesInBlock,
				})
			}
		}
	}
	
	return blocks, nil
}

// countLines counts newlines in text
func countLines(text string) int {
	count := 0
	for _, ch := range text {
		if ch == '\n' {
			count++
		}
	}
	return count
}

// ExtractExplanation extracts text without code blocks
func (e *Extractor) ExtractExplanation(content string) (string, error) {
	// Remove code blocks and return remaining text
	cleaned := e.removeCodeBlocks(content)
	return cleaned, nil
}

// removeCodeBlocks removes all code blocks from content
func (e *Extractor) removeCodeBlocks(content string) string {
	result := content
	
	// Remove fenced code blocks
	for _, pattern := range e.patterns {
		if pattern.Name == "FencedCodeBlock" || pattern.Name == "AltFencedBlock" {
			result = pattern.Regex.ReplaceAllString(result, "")
		}
	}
	
	// Clean up extra newlines (replace 3+ newlines with 2)
	result = cleanupNewlines(result)
	
	return result
}

// cleanupNewlines reduces excessive newlines
func cleanupNewlines(text string) string {
	// Simple approach: replace 3+ newlines with 2
	for {
		cleaned := regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")
		if cleaned == text {
			break
		}
		text = cleaned
	}
	return text
}