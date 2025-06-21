package handlers

import (
	"context"
	"fmt"
	"time"
)

// CheckHandler handles the check tool
type CheckHandler struct {
	storage Storage
}

// NewCheckHandler creates a new check handler
func NewCheckHandler(storage Storage) *CheckHandler {
	return &CheckHandler{
		storage: storage,
	}
}

// Handle processes a check request
func (h *CheckHandler) Handle(ctx context.Context, req CheckRequest) (*CheckResponse, error) {
	// Validate request
	if req.OutputID == "" {
		return nil, fmt.Errorf("output_id is required")
	}

	// Get output from storage
	output, err := h.storage.Get(req.OutputID)
	if err != nil {
		return nil, fmt.Errorf("get output: %w", err)
	}

	// Build response with metadata
	resp := &CheckResponse{
		ID:               output.ID,
		CreatedAt:        output.CreatedAt.Format(time.RFC3339),
		Model:            output.Model,
		FileSizeBytes:    output.Metadata.TotalBytes,
		EstimatedTokens:  output.Metadata.EstimatedTokens,
		HasCode:          len(output.Response.Extracted.Code) > 0,
		HasExplanation:   output.Response.Extracted.Explanation != "",
		CodeBlocksCount:  len(output.Response.Extracted.Code),
	}

	return resp, nil
}

// CheckRequest represents the check tool parameters
type CheckRequest struct {
	OutputID string `json:"output_id"`
}

// CheckResponse represents the check tool response
type CheckResponse struct {
	ID               string `json:"id"`
	CreatedAt        string `json:"created_at"`
	Model            string `json:"model"`
	FileSizeBytes    int64  `json:"file_size_bytes"`
	EstimatedTokens  int    `json:"estimated_tokens"`
	HasCode          bool   `json:"has_code"`
	HasExplanation   bool   `json:"has_explanation"`
	CodeBlocksCount  int    `json:"code_blocks_count"`
}