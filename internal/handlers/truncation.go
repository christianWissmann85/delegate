package handlers

import (
	"strings"
	"unicode"
)

// TruncationResult holds the outcome of truncation detection
type TruncationResult struct {
	IsTruncated bool
	Confidence  float64
	Reason      string
}

// DetectTruncation analyzes content for signs of truncation
// This is the 80/20 solution - catches most cases with minimal complexity
func DetectTruncation(content string) TruncationResult {
	if len(content) == 0 {
		return TruncationResult{IsTruncated: false, Confidence: 0.0}
	}

	// Only examine the last 500 characters for efficiency
	tailSize := 500
	tail := content
	if len(content) > tailSize {
		tail = content[len(content)-tailSize:]
	}

	// Trim whitespace for analysis
	trimmed := strings.TrimSpace(tail)
	if trimmed == "" {
		return TruncationResult{IsTruncated: false, Confidence: 0.0}
	}

	// 1. Check for complete endings (high confidence not truncated)
	completeEndings := []string{
		".", "!", "?", ";", "</", "/>", "}", "]", ")",
		"\n\n", "---", "* * *", "[DONE]", "[END]", "EOF",
	}
	
	// Special handling for code fences - check if content before ``` looks truncated
	if strings.HasSuffix(trimmed, "```") {
		// Get content before the closing ```
		beforeFence := trimmed[:len(trimmed)-3]
		beforeFence = strings.TrimSpace(beforeFence)
		
		// If the content before ``` ends with incomplete patterns, it's still truncated
		if strings.HasSuffix(beforeFence, ",") || 
		   strings.HasSuffix(beforeFence, ":") ||
		   strings.HasSuffix(beforeFence, "\"") ||
		   !strings.HasSuffix(beforeFence, "}") && !strings.HasSuffix(beforeFence, "]") {
			// Don't return early - let the rest of the detection run
		} else {
			return TruncationResult{IsTruncated: false, Confidence: 0.9}
		}
	}
	
	for _, ending := range completeEndings {
		if strings.HasSuffix(trimmed, ending) {
			return TruncationResult{IsTruncated: false, Confidence: 0.9}
		}
	}

	// 2. Check for obvious truncation patterns
	truncationPatterns := []struct {
		suffix     string
		confidence float64
		reason     string
	}{
		{`","`, 0.95, "mid-JSON array/object"},
		{`":`, 0.9, "mid-JSON key-value"},
		{`": "`, 0.95, "mid-JSON string value"},
		{`\"`, 0.85, "escaped quote at end"},
		{",", 0.7, "trailing comma"},
		{":", 0.7, "trailing colon"},
		{"=", 0.8, "trailing assignment"},
		{"&&", 0.85, "incomplete logical operator"},
		{"||", 0.85, "incomplete logical operator"},
	}

	for _, pattern := range truncationPatterns {
		if strings.HasSuffix(trimmed, pattern.suffix) {
			return TruncationResult{
				IsTruncated: true,
				Confidence:  pattern.confidence,
				Reason:      pattern.reason,
			}
		}
	}

	// 3. Count unclosed structures in the tail
	openBrackets := 0
	openParens := 0
	openBraces := 0
	quotes := 0
	backticks := 0
	inString := false
	lastWasBackslash := false

	for _, r := range tail {
		// Handle escape sequences
		if lastWasBackslash {
			lastWasBackslash = false
			continue
		}
		if r == '\\' && inString {
			lastWasBackslash = true
			continue
		}

		// Track quotes (considering escape sequences)
		if r == '"' && !lastWasBackslash {
			quotes++
			inString = quotes%2 == 1
		}

		// Only count brackets outside of strings
		if !inString {
			switch r {
			case '{':
				openBraces++
			case '}':
				openBraces--
			case '[':
				openBrackets++
			case ']':
				openBrackets--
			case '(':
				openParens++
			case ')':
				openParens--
			case '`':
				backticks++
			}
		}
	}

	// Check for unclosed structures
	if openBraces > 0 || openBrackets > 0 || openParens > 0 {
		return TruncationResult{
			IsTruncated: true,
			Confidence:  0.85,
			Reason:      "unclosed brackets/braces",
		}
	}

	if quotes%2 != 0 {
		return TruncationResult{
			IsTruncated: true,
			Confidence:  0.9,
			Reason:      "unclosed quote",
		}
	}

	if backticks%3 == 1 || backticks%3 == 2 {
		return TruncationResult{
			IsTruncated: true,
			Confidence:  0.8,
			Reason:      "incomplete code fence",
		}
	}

	// 4. Check if ends mid-word
	lastChar := rune(trimmed[len(trimmed)-1])
	if unicode.IsLetter(lastChar) || unicode.IsDigit(lastChar) {
		// Look for common word endings that suggest completion
		wordEndings := []string{"ing", "ed", "ly", "er", "est", "ion", "tion", "ment"}
		lowerTail := strings.ToLower(trimmed)
		for _, ending := range wordEndings {
			if strings.HasSuffix(lowerTail, ending) {
				return TruncationResult{IsTruncated: false, Confidence: 0.6}
			}
		}
		// If no common ending, likely truncated
		return TruncationResult{
			IsTruncated: true,
			Confidence:  0.75,
			Reason:      "ends mid-word",
		}
	}

	// 5. Check response size for suspicious boundaries
	contentSize := len(content)
	suspiciousSizes := []int{4096, 8192, 16384, 32768, 65536}
	for _, size := range suspiciousSizes {
		if contentSize >= size-50 && contentSize <= size {
			return TruncationResult{
				IsTruncated: true,
				Confidence:  0.6,
				Reason:      "suspicious size boundary",
			}
		}
	}

	// Default: probably not truncated
	return TruncationResult{IsTruncated: false, Confidence: 0.7}
}
