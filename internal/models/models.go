package models

import "time"

// GenerateRequest represents the parameters sent to an LLM provider for content generation.
// This struct is used internally by handlers to communicate with LLM providers.
type GenerateRequest struct {
	Model     string   `json:"model"`
	Prompt    string   `json:"prompt"`
	Files     []string `json:"files,omitempty"`
	MaxTokens int      `json:"max_tokens,omitempty"`
	Timeout   int      `json:"timeout,omitempty"` // Timeout in seconds
}

// StreamChunk represents a single chunk of content from a streaming LLM response.
// This struct is used internally by LLM providers to stream responses back to handlers.
type StreamChunk struct {
	Content string `json:"content"`
	Error   error  `json:"error,omitempty"` // Error encountered during streaming
}

// Extraction contains extracted code blocks and explanation text from an LLM response.
// This type is used internally by the extractor component for parsing raw LLM output.
type Extraction struct {
	Code        []CodeBlock `json:"code"`
	Explanation string      `json:"explanation"`
}

// CodeBlock represents a single extracted code block, including its language and content.
// This type is used internally by the extractor component for parsing raw LLM output.
type CodeBlock struct {
	Language  string `json:"language"`
	Content   string `json:"content"`
	LineStart int    `json:"line_start"` // Starting line number in the original raw response (1-based)
	LineEnd   int    `json:"line_end"`   // Ending line number in the original raw response (1-based)
}

// Output represents a stored LLM generation result.
// This is the primary data structure for persisting LLM responses in storage.
type Output struct {
	ID        string    `json:"id"`
	Model     string    `json:"model"`
	Prompt    string    `json:"prompt"`
	Files     []string  `json:"files,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Response  Response  `json:"response"`
	Metadata  Metadata  `json:"metadata"`
}

// Response holds the raw and extracted content of an LLM generation.
// It's part of the Output struct.
type Response struct {
	Raw       string    `json:"raw"`
	Extracted Extracted `json:"extracted"`
}

// Extracted holds the parsed code blocks and explanation from the LLM response.
// It's part of the Response struct.
type Extracted struct {
	Code        []ExtractedCode `json:"code"`
	Explanation string          `json:"explanation"`
}

// ExtractedCode represents a single code block extracted from the LLM response,
// used within the Extracted struct.
type ExtractedCode struct {
	Language  string `json:"language"`
	Content   string `json:"content"`
	LineStart int    `json:"line_start"`
	LineEnd   int    `json:"line_end"`
}

// Metadata holds various metrics and flags about the LLM generation.
// It's part of the Output struct.
type Metadata struct {
	TotalBytes           int64   `json:"total_bytes"`
	EstimatedTokens      int     `json:"estimated_tokens"`
	ProviderRequestID    string  `json:"provider_request_id,omitempty"`
	ProcessingTimeMs     int64   `json:"processing_time_ms"`
	CodeOnly             bool    `json:"code_only"`
	IsTruncated          bool    `json:"is_truncated"`
	TruncationReason     string  `json:"truncation_reason,omitempty"`
	TruncationConfidence float64 `json:"truncation_confidence,omitempty"`
}

// TruncationResult represents the result of truncation detection
type TruncationResult struct {
	IsTruncated bool    `json:"is_truncated"`
	Reason      string  `json:"reason"`
	Confidence  float64 `json:"confidence"`
}