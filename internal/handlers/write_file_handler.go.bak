package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/christianwissmann85/delegate/internal/models"
)

// WriteFileRequest represents the parameters for the delegate_write_output_to_file tool.
type WriteFileRequest struct {
	OutputID string          `json:"output_id"`
	WriteTo  string          `json:"write_to"` // The relative file path to write to.
	Options  WriteFileOptions `json:"options,omitempty"`
}

// WriteFileOptions configures what content to write.
type WriteFileOptions struct {
	Extract    string  `json:"extract,omitempty"`     // "all", "code", "explanation"
	BlockIndex *int    `json:"block_index,omitempty"` // For multi-block outputs, select a specific block.
	Language   *string `json:"language,omitempty"`    // Filter code blocks by this language.
}

// WriteFileHandler implements the delegate_write_output_to_file tool.
// It writes the content of an output artifact directly to a specified file path.
type WriteFileHandler struct {
	storage Storage
}

// NewWriteFileHandler creates a new WriteFileHandler.
func NewWriteFileHandler(storage Storage) *WriteFileHandler {
	return &WriteFileHandler{
		storage: storage,
	}
}

// Handle processes the request to write output content to a file.
func (h *WriteFileHandler) Handle(ctx context.Context, req WriteFileRequest) (*models.WriteOutputToFileResponse, error) {
	// 1. Validate Request Parameters
	if err := ValidateOutputID(req.OutputID); err != nil {
		return nil, models.AsDelegateError(err)
	}
	if req.WriteTo == "" {
		return nil, models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"The 'write_to' parameter is required.",
			"parameter", "write_to",
		)
	}

	// Set default extract option
	if req.Options.Extract == "" {
		req.Options.Extract = "all"
	}
	if err := ValidateExtractOption(req.Options.Extract); err != nil {
		return nil, models.AsDelegateError(err)
	}

	// 2. Retrieve Output from Storage
	output, err := h.storage.GetOutput(req.OutputID)
	if err != nil {
		return nil, models.NewDelegateError(
			models.ErrorTypeOutputNotFound,
			fmt.Sprintf("Output with ID '%s' not found.", req.OutputID),
			"output_id_provided", req.OutputID,
			err, // Wrap the underlying error for internal debugging
		)
	}

	// 3. Extract Content to Write based on Options
	var contentToWrite string
	switch req.Options.Extract {
	case "all":
		contentToWrite = output.Response.Raw
	case "explanation":
		contentToWrite = output.Response.Extracted.Explanation
	case "code":
			// This method handles filtering by language and block_index,
		// and formats content appropriately for file writing (e.g., no fences for code files).
		contentToWrite = h.getExtractedCodeForWriting(output.Response.Extracted.Code, req.Options, req.WriteTo)
	}

	// 4. Resolve Paths and Get Working Directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, models.NewDelegateError(
			models.ErrorTypeInternal,
			"Failed to get current working directory.",
			"original_error", err,
		)
	}

	// Resolve the relative path provided by the agent to an absolute path.
	// filepath.Join handles path cleaning (e.g., removing redundant slashes, resolving '..').
	// filepath.Abs then converts it to a canonical absolute path.
	absolutePath, err := filepath.Abs(filepath.Join(cwd, req.WriteTo))
	if err != nil {
		return nil, models.NewDelegateError(
			models.ErrorTypeFileWriteFailed,
			fmt.Sprintf("Failed to resolve absolute path for '%s'.", req.WriteTo),
			"path_provided", req.WriteTo,
			"original_error", err,
		)
	}

	// 5. Write Content to File with Security Checks
	bytesWritten, writeErr := h.writeToFile(req.WriteTo, absolutePath, cwd, []byte(contentToWrite))
	if writeErr != nil {
		return nil, writeErr
	}

	// 6. Construct Success Response
	fileSizeKB := float64(bytesWritten) / 1024.0
	message := fmt.Sprintf("Successfully wrote %.1f KB to %s", fileSizeKB, req.WriteTo)

	return &models.WriteOutputToFileResponse{
		Success:          true,
		Path:             req.WriteTo,
		AbsolutePath:     absolutePath,
		BytesWritten:     bytesWritten,
		Message:          message,
		WorkingDirectory: cwd,
	}, nil
}

// getExtractedCodeForWriting extracts code content based on options,
// suitable for writing to a file (e.g., removes markdown fences for source code files).
// It reuses logic from the original read.go's extractCodeForFile.
func (h *WriteFileHandler) getExtractedCodeForWriting(codeBlocks []models.ExtractedCode, options WriteFileOptions, filePath string) string {
	if len(codeBlocks) == 0 {
			// If no code blocks found, return empty string.
		return ""
	}

	var filteredBlocks []models.ExtractedCode
	// 1. Filter by language if specified
	if options.Language != nil && *options.Language != "" {
		langFilter := strings.ToLower(*options.Language)
		for _, block := range codeBlocks {
			if strings.ToLower(block.Language) == langFilter {
				filteredBlocks = append(filteredBlocks, block)
			}
		}
		if len(filteredBlocks) == 0 {
			// No blocks found for the specified language, return empty content.
			return ""
		}
	} else {
		// No language filter, consider all code blocks.
		filteredBlocks = codeBlocks
	}

	// 2. Select specific block if index is provided
	if options.BlockIndex != nil {
		idx := *options.BlockIndex
		if idx >= 0 && idx < len(filteredBlocks) {
			// If a specific block is requested, return its cleaned content.
			return CleanupCodeArtifacts(filteredBlocks[idx].Content, filteredBlocks[idx].Language)
		}
		// Index out of range or no blocks after filtering, return empty content.
		return ""
	}

	// 3. If no specific block/language filter, or multiple blocks remain,
	// concatenate all filtered blocks.
	// Detect if this is a source code file or documentation to decide on fences.
	if IsDocumentationFile(filePath) {
		// For documentation files, keep the fences around each block.
		var parts []string
		for _, block := range filteredBlocks {
				fence := fmt.Sprintf("```%s\n%s\n```", block.Language, block.Content)
			parts = append(parts, fence)
		}
		return strings.Join(parts, "\n\n")
	}

	// For source code files, concatenate raw code content without fences
	var parts []string
	for _, block := range filteredBlocks {
		parts = append(parts, CleanupCodeArtifacts(block.Content, block.Language))
	}

	return strings.Join(parts, "\n")
}

// writeToFile writes content to the specified file path with extensive security checks
// to prevent path traversal attacks and ensure files are written within the working directory.
func (h *WriteFileHandler) writeToFile(relativePath, absolutePath, cwd string, content []byte) (int64, error) {
	// SECURITY CHECK 1: Ensure the absolute path is within the current working directory
	// This prevents path traversal attacks where users might try to write to system files
	if !strings.HasPrefix(absolutePath, cwd) {
		return 0, models.NewDelegateError(
			models.ErrorTypePathTraversalAttempt,
			fmt.Sprintf("Security violation: attempted to write outside working directory"),
			"path_provided", relativePath,
			"resolved_path", absolutePath,
			"working_directory", cwd,
		)
	}

	// SECURITY CHECK 2: Additional validation for suspicious patterns
	// Even though filepath.Clean should handle these, we double-check for defense in depth
	if strings.Contains(relativePath, "..") || strings.Contains(absolutePath, "..") {
		return 0, models.NewDelegateError(
			models.ErrorTypePathTraversalAttempt,
			"Security violation: path contains directory traversal pattern '..'",
			"path_provided", relativePath,
		)
	}

	// SECURITY CHECK 3: Prevent writing to sensitive locations within the working directory
	// This is an additional safeguard against overwriting important project files
	sensitivePaths := []string{".git", ".ssh", ".env", "node_modules", ".delegate"}
	for _, sensitive := range sensitivePaths {
		if strings.Contains(absolutePath, string(os.PathSeparator)+sensitive+string(os.PathSeparator)) ||
			strings.HasSuffix(absolutePath, string(os.PathSeparator)+sensitive) {
			return 0, models.NewDelegateError(
				models.ErrorTypeFileWriteFailed,
				fmt.Sprintf("Cannot write to sensitive directory: %s", sensitive),
				"path_provided", relativePath,
			)
		}
	}

	// Ensure the parent directory exists
	dir := filepath.Dir(absolutePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return 0, models.NewDelegateError(
			models.ErrorTypeFileWriteFailed,
			fmt.Sprintf("Failed to create directory: %s", dir),
			"directory", dir,
			"original_error", err,
		)
	}

	// Write the file
	if err := os.WriteFile(absolutePath, content, 0644); err != nil {
		return 0, models.NewDelegateError(
			models.ErrorTypeFileWriteFailed,
			fmt.Sprintf("Failed to write file: %s", relativePath),
			"path_provided", relativePath,
			"original_error", err,
		)
	}

	return int64(len(content)), nil
}