package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

	// Track if content was truncated
	truncated := false
	originalLength := len(content)
	
	// Validate and apply token limit if specified
	if req.Options.MaxTokens > 0 {
		if err := ValidateMaxTokens(req.Options.MaxTokens); err != nil {
			return nil, err
		}
		content = h.truncateContent(content, req.Options.MaxTokens)
		truncated = len(content) < originalLength
	}

	// If WriteTo is specified, write to file instead of returning content
	if req.Options.WriteTo != "" {
		if err := h.writeToFile(req.Options.WriteTo, content); err != nil {
			return nil, models.NewDelegateError(
				models.ErrorTypeInternal,
				"",
				fmt.Sprintf("failed to write to file: %v", err),
			)
		}
		// Return success message instead of content
		return &ReadResponse{
			Content:     fmt.Sprintf("Content written to %s", req.Options.WriteTo),
			Truncated:   truncated,
			Tokens:      0, // No tokens returned when writing to file
			Extraction:  req.Options.Extract,
			FileWritten: true,
		}, nil
	}

	// Calculate approximate token count for returned content
	// Using same approximation as truncateContent: 1 token ≈ 4 characters
	tokenCount := len(content) / 4

	return &ReadResponse{
		Content:    content,
		Truncated:  truncated,
		Tokens:     tokenCount,
		Extraction: req.Options.Extract,
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
	WriteTo   string `json:"write_to,omitempty"`  // Write content to file instead of returning
}

// ReadResponse represents the read tool response
type ReadResponse struct {
	Content     string `json:"content"`
	Truncated   bool   `json:"truncated"`
	Tokens      int    `json:"tokens"`
	Extraction  string `json:"extraction"`
	Language    string `json:"language,omitempty"`
	FileWritten bool   `json:"file_written,omitempty"`
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
	
	// Ensure we have a minimum reasonable size to avoid panic
	if maxChars < 10 {
		maxChars = 10
	}
	
	if len(content) <= maxChars {
		return content
	}

	// Ensure we have enough space for ellipsis
	if maxChars <= 3 {
		return "..."
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

// writeToFile writes content to the specified file path with security checks
func (h *ReadHandler) writeToFile(filePath string, content string) error {
	// Clean and validate the path to prevent path traversal attacks
	cleanPath := filepath.Clean(filePath)
	
	// Reject paths that try to go outside current directory
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid file path: path traversal detected")
	}
	
	// Convert to absolute path for validation
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}
	
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	
	// Ensure the path is within the current working directory
	if !strings.HasPrefix(absPath, cwd) {
		return fmt.Errorf("invalid file path: must be within current directory")
	}
	
	// Ensure the directory exists
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the file
	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}