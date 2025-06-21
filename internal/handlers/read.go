package handlers

import (
	"context"
	"fmt"
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
	// TODO: Implement read logic
	return nil, fmt.Errorf("not implemented")
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