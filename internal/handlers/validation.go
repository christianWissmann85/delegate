package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/christianwissmann85/delegate/internal/models"
)

const (
	// Maximum prompt size: 100KB
	MaxPromptSize = 100 * 1024
	// Maximum file path length
	MaxFilePathLength = 1024
	// Maximum number of files
	MaxFileCount = 50
	// Maximum timeout: 10 minutes
	MaxTimeout = 600
	// Minimum timeout: 10 seconds
	MinTimeout = 10
)

// ValidateOutputID validates an output ID is safe
func ValidateOutputID(id string) error {
	if id == "" {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			"output_id is required",
		)
	}

	// Check for path traversal attempts
	if strings.Contains(id, "/") || strings.Contains(id, "\\") || strings.Contains(id, "..") {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			"invalid output_id: contains path separators",
		)
	}

	// Check length
	if len(id) > 100 {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			"invalid output_id: too long (max 100 characters)",
		)
	}

	// Check format (should be out_YYYYMMDD_HHMMSS or test_output_XXX for tests)
	if !strings.HasPrefix(id, "out_") && !strings.HasPrefix(id, "test_output_") {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			"invalid output_id: must start with 'out_' or 'test_output_'",
		)
	}

	return nil
}

// ValidateFilePaths validates file paths are safe and accessible
func ValidateFilePaths(files []string) error {
	if len(files) > MaxFileCount {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			fmt.Sprintf("too many files: %d (max %d)", len(files), MaxFileCount),
		)
	}

	for _, file := range files {
		if err := validateSingleFilePath(file); err != nil {
			return err
		}
	}

	return nil
}

func validateSingleFilePath(path string) error {
	// Check length
	if len(path) > MaxFilePathLength {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			fmt.Sprintf("file path too long: %s (max %d characters)", path, MaxFilePathLength),
		)
	}

	// Must be absolute path
	if !filepath.IsAbs(path) {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			fmt.Sprintf("file path must be absolute: %s", path),
		)
	}

	// Clean the path to resolve any .. or . elements
	cleanPath := filepath.Clean(path)

	// Check file exists and is readable
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return models.NewDelegateError(
				models.ErrorTypeNotFound,
				"",
				fmt.Sprintf("file not found: %s", cleanPath),
			)
		}
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			fmt.Sprintf("cannot access file: %s", cleanPath),
		)
	}

	// Must be a regular file
	if !info.Mode().IsRegular() {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			fmt.Sprintf("not a regular file: %s", cleanPath),
		)
	}

	// Check file size (max 1MB per file)
	if info.Size() > 1024*1024 {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			fmt.Sprintf("file too large: %s (max 1MB)", cleanPath),
		)
	}

	return nil
}

// ValidatePrompt validates the prompt
func ValidatePrompt(prompt string) error {
	if prompt == "" {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			"prompt is required",
		)
	}

	if len(prompt) > MaxPromptSize {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			fmt.Sprintf("prompt too large: %d bytes (max %d)", len(prompt), MaxPromptSize),
		)
	}

	return nil
}

// ValidateModel validates the model name
func ValidateModel(model string) error {
	if model == "" {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			"model is required",
		)
	}

	// Valid models are checked by provider factory
	return nil
}

// ValidateTimeout validates timeout value
func ValidateTimeout(timeout int) error {
	if timeout < 0 {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			"timeout cannot be negative",
		)
	}

	if timeout > 0 && timeout < MinTimeout {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			fmt.Sprintf("timeout too short: %d seconds (min %d)", timeout, MinTimeout),
		)
	}

	if timeout > MaxTimeout {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			fmt.Sprintf("timeout too long: %d seconds (max %d)", timeout, MaxTimeout),
		)
	}

	return nil
}

// ValidateMaxTokens validates max_tokens parameter
func ValidateMaxTokens(maxTokens int) error {
	if maxTokens < 0 {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			"max_tokens cannot be negative",
		)
	}

	// Provider-specific limits are handled by the providers
	return nil
}

// ValidateExtractOption validates the extract option for read
func ValidateExtractOption(extract string) error {
	if extract == "" {
		return nil // Default is "all"
	}

	switch extract {
	case "all", "code", "explanation":
		return nil
	default:
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"",
			fmt.Sprintf("invalid extract option: %s (must be 'all', 'code', or 'explanation')", extract),
		)
	}
}