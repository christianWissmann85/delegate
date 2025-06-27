// get_content_handler.go
package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/christianwissmann85/delegate/internal/models"
)

// GetContentRequest represents the parameters for the delegate_get_output_content tool.
type GetContentRequest struct {
	OutputID string            `json:"output_id"`
	Options  GetContentOptions `json:"options,omitempty"`
}

// GetContentOptions configures what content to retrieve.
type GetContentOptions struct {
	Extract    string  `json:"extract,omitempty"`     // "all", "code", "explanation"
	MaxTokens  *int    `json:"max_tokens,omitempty"`  // Limit response size
	BlockIndex *int    `json:"block_index,omitempty"` // Select a specific code block
	Language   *string `json:"language,omitempty"`    // Filter code blocks by language
}

// GetContentHandler implements the delegate_get_output_content tool.
// It retrieves the content of an output artifact into the agent's context.
type GetContentHandler struct {
	storage Storage
}

// NewGetContentHandler creates a new GetContentHandler.
func NewGetContentHandler(storage Storage) *GetContentHandler {
	return &GetContentHandler{
		storage: storage,
	}
}

// Handle processes the request to get output content.
func (h *GetContentHandler) Handle(ctx context.Context, req GetContentRequest) (*models.GetContentResponse, error) {
	if err := ValidateOutputID(req.OutputID); err != nil {
		return nil, models.AsDelegateError(err)
	}

	if req.Options.Extract == "" {
		req.Options.Extract = "all"
	}
	if err := ValidateExtractOption(req.Options.Extract); err != nil {
		return nil, models.AsDelegateError(err)
	}

	output, err := h.storage.GetOutput(req.OutputID)
	if err != nil {
		return nil, models.NewDelegateError(models.ErrorTypeOutputNotFound,
			fmt.Sprintf("Output with ID '%s' not found.", req.OutputID),
			"output_id_provided", req.OutputID,
			err,
		)
	}

	var content string
	switch req.Options.Extract {
	case "all":
		content = output.Response.Raw
	case "explanation":
		content = output.Response.Extracted.Explanation
	case "code":
		content = h.extractCodeContent(output.Response.Extracted.Code, req.Options)
	}

	isTruncated := false
	var truncationReason *string
	if req.Options.MaxTokens != nil {
		if err := ValidateMaxTokens(*req.Options.MaxTokens); err != nil {
			return nil, models.AsDelegateError(err)
		}
		originalLength := len(content)
		content = TruncateContent(content, *req.Options.MaxTokens)
		if len(content) < originalLength {
			isTruncated = true
			reason := "MAX_TOKENS_REACHED"
			truncationReason = &reason
		}
	}

	// Approximate token count for the returned content
	tokenCount := len(content) / 4

	resp := &models.GetContentResponse{
		Content: content,
		Metadata: models.ContentMetadataBlock{
			OutputID:         req.OutputID,
			TokensReturned:   tokenCount,
			IsTruncated:      isTruncated,
			TruncationReason: truncationReason,
		},
	}

	return resp, nil
}

// extractCodeContent formats code blocks into a single string based on options.
func (h *GetContentHandler) extractCodeContent(codeBlocks []models.ExtractedCode, options GetContentOptions) string {
	if len(codeBlocks) == 0 {
		return ""
	}

	// Filter blocks by language if specified
	var filteredBlocks []models.ExtractedCode
	if options.Language != nil {
		langFilter := strings.ToLower(*options.Language)
		for _, block := range codeBlocks {
			if strings.ToLower(block.Language) == langFilter {
				filteredBlocks = append(filteredBlocks, block)
			}
		}
	} else {
		filteredBlocks = codeBlocks
	}

	// If a specific block is requested, return only that one's content
	if options.BlockIndex != nil {
		idx := *options.BlockIndex
		if idx >= 0 && idx < len(filteredBlocks) {
			// Return just the content of the single block, unfenced
			return filteredBlocks[idx].Content
		}
		// If index is out of range, return empty string.
		return ""
	}

	// Concatenate all filtered blocks, adding markdown fences
	var parts []string
	for _, block := range filteredBlocks {
			fence := fmt.Sprintf("```%s\n%s\n```", block.Language, block.Content)
		parts = append(parts, fence)
	}

	return strings.Join(parts, "\n\n")
}