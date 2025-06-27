package extractor

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/christianwissmann85/delegate/internal/models"
)

// Extractor extracts code and explanations from LLM responses
type Extractor struct {
	patterns     []Pattern
	languageHint string // optional hint for language detection
}

// New creates a new extractor
func New() *Extractor {
	return &Extractor{
		patterns: GetPatterns(),
	}
}

// NewWithHint creates a new extractor with language hint
func NewWithHint(languageHint string) *Extractor {
	return &Extractor{
		patterns:     GetPatterns(),
		languageHint: languageHint,
	}
}

// Extract extracts both code and explanation
func (e *Extractor) Extract(content string, codeOnly bool) (*models.Extraction, error) {
	// Handle edge case: empty content
	if strings.TrimSpace(content) == "" {
		return &models.Extraction{
			Code:        []models.CodeBlock{},
			Explanation: "",
		}, nil
	}

	code, err := e.ExtractCode(content)
	if err != nil {
		return nil, fmt.Errorf("extract code: %w", err)
	}

	explanation, err := e.ExtractExplanation(content)
	if err != nil {
		return nil, fmt.Errorf("extract explanation: %w", err)
	}

	return &models.Extraction{
		Code:        code,
		Explanation: explanation,
	}, nil
}

// ExtractCodeOnly extracts only code blocks, ignoring explanations
func (e *Extractor) ExtractCodeOnly(content string) ([]models.CodeBlock, error) {
	// Handle edge case: empty content
	if strings.TrimSpace(content) == "" {
		return []models.CodeBlock{}, nil
	}

	return e.ExtractCode(content)
}

// ExtractCode extracts all code blocks
func (e *Extractor) ExtractCode(content string) ([]models.CodeBlock, error) {
	var blocks []models.CodeBlock
	usedRanges := make(map[string]bool)

	// Process patterns in priority order
	for _, pattern := range e.patterns {
		switch pattern.Name {
		case "FencedCodeBlock", "AltFencedBlock":
			blocks = append(blocks, e.extractFencedBlocks(content, pattern, usedRanges)...)
		case "HTMLCodeBlock":
			blocks = append(blocks, e.extractHTMLBlocks(content, pattern, usedRanges)...)
		case "IndentedBlock":
			// Only extract indented blocks if no fenced blocks found
			if len(blocks) == 0 {
				blocks = append(blocks, e.extractIndentedBlocks(content, pattern, usedRanges)...)
			}
		}
	}

	return blocks, nil
}

// extractFencedBlocks extracts fenced code blocks (``` or ~~~)
func (e *Extractor) extractFencedBlocks(content string, pattern Pattern, usedRanges map[string]bool) []models.CodeBlock {
	var blocks []models.CodeBlock

	matches := pattern.Regex.FindAllStringSubmatch(content, -1)
	indices := pattern.Regex.FindAllStringIndex(content, -1)

	for i, match := range matches {
		if len(match) < 3 || len(indices) <= i {
			continue
		}

		start, end := indices[i][0], indices[i][1]
		rangeKey := fmt.Sprintf("%d-%d", start, end)

		// Skip if overlaps with existing block
		if e.overlapsWithUsedRange(start, end, usedRanges) {
			continue
		}
		usedRanges[rangeKey] = true

		// Extract and normalize language
		lang := NormalizeLanguage(match[1])
		code := strings.TrimRight(match[2], "\n") // Trim trailing newline from code

		// If no language specified, try to detect it
		if lang == "plaintext" && e.languageHint != "" {
			lang = NormalizeLanguage(e.languageHint)
		} else if lang == "plaintext" {
			lang = e.detectLanguage(code)
		}

		// Calculate line numbers
		linesBefore := countLines(content[:start])
		linesInBlock := countLines(code)

		blocks = append(blocks, models.CodeBlock{
			Language:  lang,
			Content:   code,
			LineStart: linesBefore + 1,
			LineEnd:   linesBefore + linesInBlock,
		})
	}

	return blocks
}

// extractHTMLBlocks extracts HTML <code> blocks
func (e *Extractor) extractHTMLBlocks(content string, pattern Pattern, usedRanges map[string]bool) []models.CodeBlock {
	var blocks []models.CodeBlock

	matches := pattern.Regex.FindAllStringSubmatch(content, -1)
	indices := pattern.Regex.FindAllStringIndex(content, -1)

	for i, match := range matches {
		if len(match) < 3 || len(indices) <= i {
			continue
		}

		start, end := indices[i][0], indices[i][1]
		rangeKey := fmt.Sprintf("%d-%d", start, end)

		if e.overlapsWithUsedRange(start, end, usedRanges) {
			continue
		}
		usedRanges[rangeKey] = true

		// Extract language from class attribute if present
		lang := "plaintext"
		if match[1] != "" {
			lang = NormalizeLanguage(match[1])
		}
		code := strings.TrimRight(match[2], "\n")

		// If no language specified, try to detect it
		if lang == "plaintext" && e.languageHint != "" {
			lang = NormalizeLanguage(e.languageHint)
		} else if lang == "plaintext" {
			lang = e.detectLanguage(code)
		}

		linesBefore := countLines(content[:start])
		linesInBlock := countLines(code)

		blocks = append(blocks, models.CodeBlock{
			Language:  lang,
			Content:   code,
			LineStart: linesBefore + 1,
			LineEnd:   linesBefore + linesInBlock,
		})
	}

	return blocks
}

// extractIndentedBlocks extracts indented code blocks
func (e *Extractor) extractIndentedBlocks(content string, pattern Pattern, usedRanges map[string]bool) []models.CodeBlock {
	var blocks []models.CodeBlock

	matches := pattern.Regex.FindAllString(content, -1)
	indices := pattern.Regex.FindAllStringIndex(content, -1)

	for i, match := range matches {
		if len(indices) <= i {
			continue
		}

		start, end := indices[i][0], indices[i][1]
		rangeKey := fmt.Sprintf("%d-%d", start, end)

		if e.overlapsWithUsedRange(start, end, usedRanges) {
			continue
		}
		usedRanges[rangeKey] = true

		// Remove common indentation
		code := e.removeCommonIndentation(match)

		linesBefore := countLines(content[:start])
		linesInBlock := countLines(code)

		// Try to detect language for indented blocks
		var lang string
		if e.languageHint != "" {
			lang = NormalizeLanguage(e.languageHint)
		} else {
			lang = e.detectLanguage(code)
		}

		blocks = append(blocks, models.CodeBlock{
			Language:  lang,
			Content:   code,
			LineStart: linesBefore + 1,
			LineEnd:   linesBefore + linesInBlock,
		})
	}

	return blocks
}

// overlapsWithUsedRange checks if a range overlaps with any used range
func (e *Extractor) overlapsWithUsedRange(start, end int, usedRanges map[string]bool) bool {
	for rangeKey := range usedRanges {
		var usedStart, usedEnd int
		if _, err := fmt.Sscanf(rangeKey, "%d-%d", &usedStart, &usedEnd); err != nil {
			continue
		}

		// Check for overlap
		if start < usedEnd && end > usedStart {
			return true
		}
	}
	return false
}

// removeCommonIndentation removes the common leading whitespace from lines
func (e *Extractor) removeCommonIndentation(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return text
	}

	// Find minimum indentation
	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return text
	}

	// Remove common indentation
	for i, line := range lines {
		if len(line) >= minIndent {
			lines[i] = line[minIndent:]
		}
	}

	return strings.Join(lines, "\n")
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
	// Get all code blocks to know what to exclude
	codeBlocks, err := e.ExtractCode(content)
	if err != nil {
		return "", fmt.Errorf("extract code blocks: %w", err)
	}

	// If no code blocks, return the entire content as explanation
	if len(codeBlocks) == 0 {
		return strings.TrimSpace(content), nil
	}

	// Build a list of ranges to exclude
	type exclusion struct {
		start, end int
	}
	var exclusions []exclusion

	// Find positions of all code blocks in original content
	for _, pattern := range e.patterns {
		if pattern.Name == "FencedCodeBlock" || pattern.Name == "AltFencedBlock" || pattern.Name == "HTMLCodeBlock" {
			indices := pattern.Regex.FindAllStringIndex(content, -1)
			for _, idx := range indices {
				exclusions = append(exclusions, exclusion{start: idx[0], end: idx[1]})
			}
		}
	}

	// Don't exclude inline code from explanations - we want to keep them

	// Sort exclusions by start position
	for i := 0; i < len(exclusions); i++ {
		for j := i + 1; j < len(exclusions); j++ {
			if exclusions[j].start < exclusions[i].start {
				exclusions[i], exclusions[j] = exclusions[j], exclusions[i]
			}
		}
	}

	// Build explanation by taking non-excluded parts
	var explanationParts []string
	lastEnd := 0

	for _, excl := range exclusions {
		if excl.start > lastEnd {
			part := content[lastEnd:excl.start]
			part = strings.TrimSpace(part)
			if part != "" {
				explanationParts = append(explanationParts, part)
			}
		}
		lastEnd = excl.end
	}

	// Add any remaining content after the last exclusion
	if lastEnd < len(content) {
		part := content[lastEnd:]
		part = strings.TrimSpace(part)
		if part != "" {
			explanationParts = append(explanationParts, part)
		}
	}

	// Join parts with proper spacing
	explanation := strings.Join(explanationParts, "\n\n")

	// Clean up extra newlines
	explanation = cleanupNewlines(explanation)

	return strings.TrimSpace(explanation), nil
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

// detectLanguage attempts to detect the language from code content
func (e *Extractor) detectLanguage(code string) string {
	// Common patterns for language detection - order matters!
	patterns := []struct {
		lang    string
		pattern *regexp.Regexp
	}{
		// Check for shebang first
		{"bash", regexp.MustCompile(`(?m)^#!/bin/(ba)?sh`)},
		{"python", regexp.MustCompile(`(?m)^#!/usr/bin/(env )?python`)},
		// Then check for specific language patterns
		{"python", regexp.MustCompile(`(?m)(^def |^class |^import |^from .+ import|if __name__ == ['"]__main__['"]:|print\()`)},
		{"javascript", regexp.MustCompile(`(?m)(^(const|let|var) |^function |=> |require\(|module\.exports|console\.log\()`)},
		{"typescript", regexp.MustCompile(`(?m)(^interface |^type |^enum |^namespace |^declare |: (string|number|boolean|any)\b)`)},
		{"go", regexp.MustCompile(`(?m)(^package |^import \(|^func |^type .+ struct|^var .+ = |fmt\.Print)`)},
		{"java", regexp.MustCompile(`(?m)(^public class |^private |^protected |^package |^import java\.|System\.out\.print)`)},
		{"cpp", regexp.MustCompile(`(?m)(^#include <|^using namespace |^int main\(|std::|cout <<|cin >>)`)},
		{"c", regexp.MustCompile(`(?m)(^#include <.*\.h>|^int main\(|printf\(|scanf\()`)},
		{"rust", regexp.MustCompile(`(?m)(^use |^fn |^let mut |^impl |^struct |^enum |println!\()`)},
		{"ruby", regexp.MustCompile(`(?m)(^def |^class |^module |^require |puts |p |attr_)`)},
		{"php", regexp.MustCompile(`(?m)(<\?php|\$[a-zA-Z_]|^function |echo |print_r\()`)},
		{"sql", regexp.MustCompile(`(?im)(^SELECT |^INSERT |^UPDATE |^DELETE |^CREATE TABLE |^ALTER |^DROP )`)},
		{"yaml", regexp.MustCompile(`(?m)^[a-zA-Z_]+:\s+(.*\n(  |\t))?`)},
		{"json", regexp.MustCompile(`^\s*\{[\s\S]*\}\s*$|^\s*\[[\s\S]*\]\s*$`)},
		{"xml", regexp.MustCompile(`^\s*<\?xml|^\s*<[a-zA-Z]+.*>[\s\S]*</[a-zA-Z]+>\s*$`)},
		{"html", regexp.MustCompile(`(?i)(<html|<head|<body|<div|<p|<h[1-6]|<!DOCTYPE html)`)},
		{"dockerfile", regexp.MustCompile(`(?m)(^FROM |^RUN |^CMD |^EXPOSE |^ENV |^WORKDIR )`)},
		{"terraform", regexp.MustCompile(`(?m)(^resource "|^provider "|^variable "|^output "|^module ")`)},
	}

	// Check each pattern
	for _, p := range patterns {
		if p.pattern.MatchString(code) {
			return p.lang
		}
	}

	// Check for shebang
	if strings.HasPrefix(strings.TrimSpace(code), "#!") {
		firstLine := strings.Split(code, "\n")[0]
		if strings.Contains(firstLine, "python") {
			return "python"
		} else if strings.Contains(firstLine, "node") {
			return "javascript"
		} else if strings.Contains(firstLine, "ruby") {
			return "ruby"
		} else if strings.Contains(firstLine, "sh") {
			return "bash"
		}
	}

	return "plaintext"
}
