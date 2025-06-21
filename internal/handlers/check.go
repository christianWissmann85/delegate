package handlers

import (
	"context"
	"fmt"
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
	// TODO: Implement check logic
	return nil, fmt.Errorf("not implemented")
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