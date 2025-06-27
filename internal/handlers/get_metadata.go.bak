// handlers/get_metadata.go
package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/christianwissmann85/delegate/internal/models"
)

// GetMetadataHandler handles the delegate_get_output_metadata tool.
// It retrieves structured metadata and a content analysis for a given output_id,
// allowing an agent to inspect an output's properties without fetching its full content.
type GetMetadataHandler struct {
	storage Storage
}

// NewGetMetadataHandler creates a new, configured GetMetadataHandler.
func NewGetMetadataHandler(storage Storage) *GetMetadataHandler {
	return &GetMetadataHandler{
		storage: storage,
	}
}

// GetMetadataRequest represents the parameters for the delegate_get_output_metadata tool.
type GetMetadataRequest struct {
	OutputID string `json:"output_id"`
}

// Handle processes a request to get output metadata.
func (h *GetMetadataHandler) Handle(ctx context.Context, req GetMetadataRequest) (*models.GetOutputMetadataResponse, error) {
	if err := ValidateOutputID(req.OutputID); err != nil {
		return nil, models.AsDelegateError(err)
	}

	output, err := h.storage.GetOutput(req.OutputID)
	if err != nil {
		// Wrap the storage error into a structured DelegateError with details.
		return nil, models.NewDelegateError(
			models.ErrorTypeOutputNotFound,
			fmt.Sprintf("The requested output ID does not exist or has expired."),
			"output_id_provided", req.OutputID,
			err,
		)
	}

	// Build the structured metadata block.
	var truncationReason *string
	if output.Metadata.IsTruncated && output.Metadata.TruncationReason != "" {
		reason := output.Metadata.TruncationReason
		truncationReason = &reason
	}

	metadataBlock := models.MetadataBlock{
		OutputID:         output.ID,
		Status:           "COMPLETED", // If the output is retrieved, it's considered completed.
		SizeKB:           float64(output.Metadata.TotalBytes) / 1024.0,
		LineCount:        strings.Count(output.Response.Raw, "\n") + 1,
		TokenEstimate:    output.Metadata.EstimatedTokens,
		IsTruncated:      output.Metadata.IsTruncated,
		TruncationReason: truncationReason,
	}

	// Build the content analysis block with details about each code block.
	analysisBlocks := make([]models.CodeBlockMetadata, 0, len(output.Response.Extracted.Code))
	for i, block := range output.Response.Extracted.Code {
		preview := ""
		if block.Content != "" {
			preview = strings.SplitN(block.Content, "\n", 2)[0]
		}

		analysisBlocks = append(analysisBlocks, models.CodeBlockMetadata{
			Index:    i,
			Language: block.Language,
			SizeKB:   float64(len(block.Content)) / 1024.0,
			Lines:    strings.Count(block.Content, "\n") + 1,
			Preview:  preview,
		})
	}

	contentAnalysisBlock := models.ContentAnalysisBlock{
		BlocksFound: len(analysisBlocks),
		Blocks:      analysisBlocks,
	}

	// Assemble the final structured response.
	resp := &models.GetOutputMetadataResponse{
		Metadata:        metadataBlock,
		ContentAnalysis: contentAnalysisBlock,
	}

	return resp, nil
}

