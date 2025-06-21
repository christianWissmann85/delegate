package handlers

import (
	"strings"
	"unicode"
)

// EstimateTokens provides a more accurate token count estimation
// This is still an approximation - actual tokenization varies by model
func EstimateTokens(text string) int {
	if text == "" {
		return 0
	}

	// More sophisticated estimation based on common tokenization patterns
	// Most LLMs use subword tokenization (like BPE)
	
	// Count words first
	words := 0
	inWord := false
	for _, r := range text {
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			inWord = false
		} else if !inWord {
			words++
			inWord = true
		}
	}

	// Count special tokens (code blocks, punctuation clusters)
	specialTokens := 0
	
	// Code fences are usually 2-3 tokens each
	specialTokens += strings.Count(text, "```") * 2
	
	// URLs tend to be multiple tokens
	specialTokens += strings.Count(text, "http://") * 3
	specialTokens += strings.Count(text, "https://") * 3
	
	// Common multi-character operators in code
	codeOperators := []string{"==", "!=", "<=", ">=", "&&", "||", "++", "--", "->", "=>", ":="}
	for _, op := range codeOperators {
		specialTokens += strings.Count(text, op)
	}

	// Estimate based on character count for very long strings (like base64)
	// These typically have high character-to-token ratios
	charBasedEstimate := len(text) / 4

	// Use the higher of word-based or character-based estimate
	wordBasedEstimate := words + specialTokens
	
	// For code-heavy content, lean toward character-based
	// For natural language, lean toward word-based
	codeIndicators := strings.Count(text, "{") + strings.Count(text, "}") + 
	                  strings.Count(text, "(") + strings.Count(text, ")") +
	                  strings.Count(text, ";")
	
	if codeIndicators > words/10 { // More than 10% code indicators
		// Likely code-heavy, use weighted average favoring characters
		return (charBasedEstimate*2 + wordBasedEstimate) / 3
	}
	
	// Likely natural language, use weighted average favoring words
	return (wordBasedEstimate*2 + charBasedEstimate) / 3
}

// EstimateTokensForJSON estimates tokens for JSON content
// JSON typically has more tokens due to structure characters
func EstimateTokensForJSON(text string) int {
	// JSON has additional overhead from quotes, braces, etc.
	baseEstimate := EstimateTokens(text)
	
	// Count structural elements that add tokens
	structureTokens := strings.Count(text, "\"") / 2 + // Each key/value pair
	                   strings.Count(text, ":") +      // Colons
	                   strings.Count(text, ",") +      // Commas
	                   strings.Count(text, "{") +      // Object starts
	                   strings.Count(text, "}")        // Object ends
	
	return baseEstimate + structureTokens/2 // Structural elements often combine into single tokens
}