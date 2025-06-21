package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/christianwissmann85/delegate/internal/models"
)

// ReadHandler handles the read tool
type ReadHandler struct {
	storage Storage
}

// NewReadHandler creates a new read handler
func NewReadHandler(storage Storage) *ReadHandler {
	return &ReadHandler{
		storage: storage,
	}
}

// Handle processes a read request
func (h *ReadHandler) Handle(ctx context.Context, req ReadRequest) (*ReadResponse, error) {
	// Validate request
	if err := ValidateOutputID(req.OutputID); err != nil {
		return nil, err
	}

	// Set default extract option
	if req.Options.Extract == "" {
		req.Options.Extract = "all"
	}

	// Validate extract option
	if err := ValidateExtractOption(req.Options.Extract); err != nil {
		return nil, err
	}

	// Get output from storage
	output, err := h.storage.Get(req.OutputID)
	if err != nil {
		return nil, models.NewDelegateError(
			models.ErrorTypeNotFound,
			"",
			fmt.Sprintf("output not found: %v", err),
		)
	}

	// Extract requested content
	var content string
	switch req.Options.Extract {
	case "all":
		content = output.Response.Raw
	case "code":
		content = h.extractCodeContent(output)
	case "explanation":
		content = output.Response.Extracted.Explanation
	}

	// Validate and apply token limit if specified
	if req.Options.MaxTokens > 0 {
		if err := ValidateMaxTokens(req.Options.MaxTokens); err != nil {
			return nil, err
		}
		content = h.truncateContent(content, req.Options.MaxTokens)
	}

	return &ReadResponse{
		Content: content,
	}, nil
}

// ReadRequest represents the read tool parameters
type ReadRequest struct {
	OutputID string      `json:"output_id"`
	Options  ReadOptions `json:"options,omitempty"`
}

// ReadOptions configures what to read
type ReadOptions struct {
	Extract   string `json:"extract,omitempty"`   // "all", "code", "explanation"
	MaxTokens int    `json:"max_tokens,omitempty"` // Limit response size
}

// ReadResponse represents the read tool response
type ReadResponse struct {
	Content string `json:"content"`
}

// extractCodeContent formats all code blocks into a single string
func (h *ReadHandler) extractCodeContent(output *models.Output) string {
	if len(output.Response.Extracted.Code) == 0 {
		return ""
	}

	var parts []string
	for i, block := range output.Response.Extracted.Code {
		// Format as fenced code block
		fence := fmt.Sprintf("```%s\n%s\n```", block.Language, block.Content)
		parts = append(parts, fence)
		
		// Add separator between blocks (except last)
		if i < len(output.Response.Extracted.Code)-1 {
			parts = append(parts, "")
		}
	}

	return strings.Join(parts, "\n")
}

// truncateContent truncates content to approximately maxTokens
func (h *ReadHandler) truncateContent(content string, maxTokens int) string {
	// Simple approximation: 1 token ≈ 4 characters
	// This is a rough estimate; actual tokenization varies by model
	maxChars := maxTokens * 4
	
	if len(content) <= maxChars {
		return content
	}

	// Truncate and add ellipsis
	truncated := content[:maxChars-3] + "..."
	
	// Try to break at a word boundary
	lastSpace := strings.LastIndexAny(truncated[:len(truncated)-3], " \n\t")
	if lastSpace > maxChars*3/4 { // Only break at word if we're keeping at least 75% of content
		truncated = truncated[:lastSpace] + "..."
	}

	return truncated
}