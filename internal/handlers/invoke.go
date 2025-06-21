package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/christianwissmann85/delegate/internal/models"
)

// InvokeHandler handles the invoke tool
type InvokeHandler struct {
	providers ProviderFactory
	storage   Storage
	extractor Extractor
}

// Provider generates content from LLMs
type Provider interface {
	GenerateStream(ctx context.Context, req GenerateRequest) (<-chan StreamChunk, error)
}

// ProviderFactory creates providers based on model
type ProviderFactory interface {
	GetProvider(model string) (Provider, error)
}

// Storage persists outputs
type Storage interface {
	Save(output *models.Output) error
	Get(id string) (*models.Output, error)
	Delete(id string) error
	ListOlderThan(age time.Duration) ([]string, error)
}

// Extractor extracts code and explanations
type Extractor interface {
	Extract(content string) (*Extraction, error)
	ExtractCode(content string) ([]CodeBlock, error)
	ExtractExplanation(content string) (string, error)
}

// NewInvokeHandler creates a new invoke handler
func NewInvokeHandler(providers ProviderFactory, storage Storage, extractor Extractor) *InvokeHandler {
	return &InvokeHandler{
		providers: providers,
		storage:   storage,
		extractor: extractor,
	}
}

// Handle processes an invoke request
func (h *InvokeHandler) Handle(ctx context.Context, req InvokeRequest) (*InvokeResponse, error) {
	// Validate request
	if err := h.validateRequest(req); err != nil {
		return nil, err
	}

	// Get provider for the model
	provider, err := h.providers.GetProvider(req.Model)
	if err != nil {
		return nil, fmt.Errorf("get provider: %w", err)
	}

	// Create generate request
	genReq := GenerateRequest{
		Model:     req.Model,
		Prompt:    req.Prompt,
		Files:     req.Files,
		MaxTokens: req.MaxTokens,
		Timeout:   req.Timeout,
	}

	// Stream response from provider
	stream, err := provider.GenerateStream(ctx, genReq)
	if err != nil {
		return nil, fmt.Errorf("start stream: %w", err)
	}

	// Collect response
	var fullResponse string
	for chunk := range stream {
		if chunk.Error != nil {
			return nil, fmt.Errorf("stream error: %w", chunk.Error)
		}
		fullResponse += chunk.Content
	}

	// Extract code and explanation
	extraction, err := h.extractor.Extract(fullResponse)
	if err != nil {
		// If extraction fails, still save the raw response
		extraction = &Extraction{
			Explanation: fullResponse,
		}
	}

	// Create output object
	output := &models.Output{
		Model:     req.Model,
		Prompt:    req.Prompt,
		CreatedAt: time.Now().UTC(),
		Response: models.Response{
			Raw: fullResponse,
			Extracted: models.Extracted{
				Explanation: extraction.Explanation,
			},
		},
		Metadata: models.Metadata{
			TotalBytes:      int64(len(fullResponse)),
			EstimatedTokens: estimateTokens(fullResponse),
		},
	}

	// Convert code blocks
	for _, block := range extraction.Code {
		output.Response.Extracted.Code = append(output.Response.Extracted.Code, models.ExtractedCode{
			Language: block.Language,
			Content:  block.Content,
		})
	}

	// Apply code_only filter if requested
	if req.CodeOnly && len(output.Response.Extracted.Code) > 0 {
		output.Response.Extracted.Explanation = ""
	}

	// Save to storage
	if err := h.storage.Save(output); err != nil {
		return nil, fmt.Errorf("save output: %w", err)
	}

	return &InvokeResponse{
		OutputID: output.ID,
	}, nil
}

// validateRequest checks if the request is valid
func (h *InvokeHandler) validateRequest(req InvokeRequest) error {
	if req.Model == "" {
		return fmt.Errorf("model is required")
	}
	if req.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	if req.MaxTokens < 0 {
		return fmt.Errorf("max_tokens cannot be negative")
	}
	if req.Timeout < 0 {
		return fmt.Errorf("timeout cannot be negative")
	}
	return nil
}

// estimateTokens provides a rough token count estimate
func estimateTokens(text string) int {
	// Simple estimation: ~4 characters per token
	return len(text) / 4
}

// InvokeRequest represents the invoke tool parameters
type InvokeRequest struct {
	Model        string   `json:"model"`
	Prompt       string   `json:"prompt"`
	Files        []string `json:"files,omitempty"`
	MaxTokens    int      `json:"max_tokens,omitempty"`
	CodeOnly     bool     `json:"code_only,omitempty"`
	LanguageHint string   `json:"language_hint,omitempty"`
	Timeout      int      `json:"timeout,omitempty"`
}

// InvokeResponse represents the invoke tool response
type InvokeResponse struct {
	OutputID string `json:"output_id"`
}

// GenerateRequest is sent to providers
type GenerateRequest struct {
	Model     string
	Prompt    string
	Files     []string
	MaxTokens int
	Timeout   int // Timeout in seconds
}

// StreamChunk represents a chunk of generated content
type StreamChunk struct {
	Content string
	Error   error
}

// Extraction contains extracted code and explanation
type Extraction struct {
	Code        []CodeBlock
	Explanation string
}

// CodeBlock represents an extracted code block
type CodeBlock struct {
	Language  string
	Content   string
	LineStart int
	LineEnd   int
}