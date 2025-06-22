package models

import "time"

// Output represents a stored generation result
type Output struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Model     string    `json:"model"`
	Prompt    string    `json:"prompt"`
	Files     []string  `json:"files,omitempty"`
	Response  Response  `json:"response"`
	Metadata  Metadata  `json:"metadata"`
}

// Response contains the LLM response
type Response struct {
	Raw       string    `json:"raw"`
	Extracted Extracted `json:"extracted"`
}

// Extracted contains extracted code and explanation
type Extracted struct {
	Code        []ExtractedCode `json:"code"`
	Explanation string          `json:"explanation"`
}

// ExtractedCode represents a code block
type ExtractedCode struct {
	Language  string `json:"language"`
	Content   string `json:"content"`
	LineStart int    `json:"line_start"`
	LineEnd   int    `json:"line_end"`
}

// Metadata contains output metadata
type Metadata struct {
	TotalBytes        int64  `json:"total_bytes"`
	EstimatedTokens   int    `json:"estimated_tokens"`
	ProviderRequestID string `json:"provider_request_id,omitempty"`
	ProcessingTimeMs  int64  `json:"processing_time_ms"`
	CodeOnly          bool   `json:"code_only,omitempty"`
}
