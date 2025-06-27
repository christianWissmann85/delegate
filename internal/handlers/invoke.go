package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/christianwissmann85/delegate/internal/models"
)

// InvokeHandler handles the invoke tool
type InvokeHandler struct {
	providers        ProviderFactory
	storage          Storage
	extractorFactory ExtractorFactory
}


// NewInvokeHandler creates a new invoke handler
func NewInvokeHandler(providers ProviderFactory, storage Storage, extractorFactory ExtractorFactory) *InvokeHandler {
	return &InvokeHandler{
		providers:        providers,
		storage:          storage,
		extractorFactory: extractorFactory,
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
		// Provider factory should return DelegateError, but wrap if not
		if delegateErr, ok := err.(*models.DelegateError); ok {
			return nil, delegateErr
		}
		return nil, models.NewDelegateError(
			models.ErrorTypeProviderUnavailable,
			fmt.Sprintf("Provider for model '%s' unavailable.", req.Model),
			"model", req.Model,
			"original_error", err,
		)
	}

	// Create generate request
	genReq := &models.GenerateRequest{
		Model:     req.Model,
		Prompt:    req.Prompt,
		Files:     req.Files,
		MaxTokens: req.MaxTokens,
		Timeout:   req.Timeout,
	}

	// Stream response from provider
	stream, err := provider.Generate(ctx, genReq)
	if err != nil {
		// Provider should return DelegateError
		if delegateErr, ok := err.(*models.DelegateError); ok {
			return nil, delegateErr
		}
		return nil, models.NewDelegateError(
			models.ErrorTypeProviderError,
			fmt.Sprintf("Failed to start stream for model '%s'.", req.Model),
			"model", req.Model,
			"original_error", err,
		)
	}

	// Collect response
	var fullResponse string
	for chunk := range stream {
		if chunk.Error != nil {
			// Stream errors should be DelegateError
			if delegateErr, ok := chunk.Error.(*models.DelegateError); ok {
				return nil, delegateErr
			}
			return nil, models.NewDelegateError(
				models.ErrorTypeProviderError,
				fmt.Sprintf("Stream error for model '%s'.", req.Model),
				"model", req.Model,
				"original_error", chunk.Error,
			)
		}
		fullResponse += chunk.Content
	}

	// Detect truncation
	truncationResult := DetectTruncation(fullResponse)
	
	// Debug logging
	if truncationResult.IsTruncated {
		fmt.Printf("DEBUG: Truncation detected! IsTruncated=%v, Confidence=%f, Reason='%s'\n", 
			truncationResult.IsTruncated, truncationResult.Confidence, truncationResult.Reason)
	}

	// Create extractor
	extractor := h.extractorFactory.Default()

	// Extract based on mode
	extraction, err := extractor.Extract(fullResponse, req.CodeOnly)
	if err != nil {
		// If extraction fails, still save the raw response
		extraction = models.Extraction{
			Explanation: fullResponse,
		}
	}

	// Create output object
	output := &models.Output{
		ID:        "", // Will be set by storage
		Model:     req.Model,
		Prompt:    req.Prompt,
		Files:     req.Files,
		CreatedAt: time.Now().UTC(),
		Response: models.Response{
			Raw: fullResponse,
			Extracted: models.Extracted{
				Explanation: extraction.Explanation,
			},
		},
		Metadata: models.Metadata{
			TotalBytes:           int64(len(fullResponse)),
			EstimatedTokens:      EstimateTokens(fullResponse),
			CodeOnly:             req.CodeOnly,
			IsTruncated:          truncationResult.IsTruncated,
			TruncationReason:     truncationResult.Reason,
			TruncationConfidence: truncationResult.Confidence,
		},
	}

	// Convert code blocks
	for _, block := range extraction.Code {
		output.Response.Extracted.Code = append(output.Response.Extracted.Code, models.ExtractedCode{
			Language:  block.Language,
			Content:   block.Content,
			LineStart: block.LineStart,
			LineEnd:   block.LineEnd,
		})
	}

	// Apply code_only filter if requested
	if req.CodeOnly && len(output.Response.Extracted.Code) > 0 {
		output.Response.Extracted.Explanation = ""
	}

	// Save to storage
	if err := h.storage.SaveOutput(output); err != nil {
		return nil, models.NewDelegateError(
			models.ErrorTypeInternal,
			fmt.Sprintf("Failed to save output."),
			"original_error", err,
		)
	}

	resp := &InvokeResponse{
		OutputID: output.ID,
	}
	
	// Add truncation warning if detected
	if truncationResult.IsTruncated {
		resp.Warning = fmt.Sprintf("Output likely truncated (confidence: %.2f). Use delegate_check for details or write_to to save anyway.", 
			truncationResult.Confidence)
	}
	
	return resp, nil
}

// validateRequest checks if the request is valid
func (h *InvokeHandler) validateRequest(req InvokeRequest) error {
	// Validate model
	if err := ValidateModel(req.Model); err != nil {
		return err
	}

	// Validate prompt
	if err := ValidatePrompt(req.Prompt); err != nil {
		return err
	}

	// Validate file paths if provided
	if len(req.Files) > 0 {
		if err := ValidateFilePaths(req.Files); err != nil {
			return err
		}
	}

	// Validate max tokens
	if err := ValidateMaxTokens(req.MaxTokens); err != nil {
		return err
	}

	// Validate timeout
	if err := ValidateTimeout(req.Timeout); err != nil {
		return err
	}

	return nil
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
	Warning  string `json:"warning,omitempty"`
}

