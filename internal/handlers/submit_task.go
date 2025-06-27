// handlers/submit_task.go
package handlers

import (
	"context"
	"os"
	"time"

	"github.com/christianwissmann85/delegate/internal/models"
)

// SubmitTaskHandler handles the delegate_submit_task tool.
// Its purpose is to asynchronously submit a generation task to an external LLM.
// It saves the complete output and returns a unique `output_id` for subsequent operations.
// This handler does not return the content directly.
type SubmitTaskHandler struct {
	providers        ProviderFactory
	storage          Storage
	extractorFactory ExtractorFactory
}


// NewSubmitTaskHandler creates a new, configured SubmitTaskHandler.
func NewSubmitTaskHandler(providers ProviderFactory, storage Storage, extractorFactory ExtractorFactory) *SubmitTaskHandler {
	return &SubmitTaskHandler{
		providers:        providers,
		storage:          storage,
		extractorFactory: extractorFactory,
	}
}

// SubmitTaskRequest represents the parameters for the delegate_submit_task tool.
type SubmitTaskRequest struct {
	Model     string   `json:"model"`
	Prompt    string   `json:"prompt"`
	Files     []string `json:"files,omitempty"` // List of relative file paths to include as context.
	MaxTokens int      `json:"max_tokens,omitempty"`
	Timeout   int      `json:"timeout,omitempty"`
}

// Handle processes a request to submit a generation task.
func (h *SubmitTaskHandler) Handle(ctx context.Context, req SubmitTaskRequest) (*models.SubmitTaskResponse, error) {
	if err := h.validateRequest(req); err != nil {
		return nil, models.AsDelegateError(err)
	}

	provider, err := h.providers.GetProvider(req.Model)
	if err != nil {
		return nil, models.AsDelegateError(err)
	}

	genReq := &models.GenerateRequest{
		Model:     req.Model,
		Prompt:    req.Prompt,
		Files:     req.Files,
		MaxTokens: req.MaxTokens,
		Timeout:   req.Timeout,
	}

	stream, err := provider.Generate(ctx, genReq)
	if err != nil {
		return nil, models.AsDelegateError(err)
	}

	var fullResponse string
	for chunk := range stream {
		if chunk.Error != nil {
			return nil, models.AsDelegateError(chunk.Error)
		}
		fullResponse += chunk.Content
	}

	// The following logic is preserved from the original handler to ensure
	// the stored output is correctly processed for use by other tools.
	truncationResult := DetectTruncation(fullResponse)
	extractor := h.extractorFactory.Default()
	extraction, err := extractor.Extract(fullResponse, false)
	if err != nil {
		// If extraction fails, save the raw response as the explanation.
		extraction = models.Extraction{
			Explanation: fullResponse,
		}
	}

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
			TotalBytes:      int64(len(fullResponse)),
			EstimatedTokens: EstimateTokens(fullResponse),
			IsTruncated:     truncationResult.IsTruncated,
			TruncationReason:     truncationResult.Reason,
			TruncationConfidence: truncationResult.Confidence,
		},
	}

	for _, block := range extraction.Code {
		output.Response.Extracted.Code = append(output.Response.Extracted.Code, models.ExtractedCode{
			Language:  block.Language,
			Content:   block.Content,
			LineStart: block.LineStart,
			LineEnd:   block.LineEnd,
		})
	}

	if err := h.storage.SaveOutput(output); err != nil {
		return nil, models.NewDelegateError(models.ErrorTypeInternal, "Failed to save output", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		// This is a server-side issue, so we return an internal error.
		return nil, models.NewDelegateError(models.ErrorTypeInternal, "Failed to get working directory", err)
	}

	resp := &models.SubmitTaskResponse{
		OutputID:         output.ID,
		WorkingDirectory: wd,
	}

	return resp, nil
}

// validateRequest checks if the request parameters are valid.
func (h *SubmitTaskHandler) validateRequest(req SubmitTaskRequest) error {
	if err := ValidateModel(req.Model); err != nil {
		return err
	}
	if err := ValidatePrompt(req.Prompt); err != nil {
		return err
	}
	if len(req.Files) > 0 {
		if err := ValidateFilePaths(req.Files); err != nil {
			return err
		}
	}
	if err := ValidateMaxTokens(req.MaxTokens); err != nil {
		return err
	}
	if err := ValidateTimeout(req.Timeout); err != nil {
		return err
	}
	return nil
}

