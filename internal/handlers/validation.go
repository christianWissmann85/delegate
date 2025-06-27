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
	// Minimum timeout: 10 seconds (but 120+ seconds recommended)
	MinTimeout = 10
	// Recommended timeout for stable LLM responses
	RecommendedMinTimeout = 120
)

// ValidateOutputID validates an output ID is safe
func ValidateOutputID(id string) error {
	if id == "" {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"Output ID is required", // Moved from unnamed argument
		)
	}

	// Check for path traversal attempts
	if strings.Contains(id, "/") || strings.Contains(id, "\\") || strings.Contains(id, "..") {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"Invalid output ID: contains path separators", // Moved from unnamed argument
		)
	}

	// Check length
	if len(id) > 100 {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"Invalid output ID: too long (max 100 characters)", // Moved from unnamed argument
		)
	}

	// Check format (should be out_YYYYMMDD_HHMMSS or test_output_XXX for tests)
	if !strings.HasPrefix(id, "out_") && !strings.HasPrefix(id, "test_output_") {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"Invalid output ID: must start with 'out_' or 'test_output_'", // Moved from unnamed argument
		)
	}

	return nil
}

// ValidateFilePaths validates file paths are safe and accessible
func ValidateFilePaths(files []string) error {
	if len(files) > MaxFileCount {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			fmt.Sprintf("Too many files: %d (max %d)", len(files), MaxFileCount), // Moved from unnamed argument
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
			fmt.Sprintf("File path too long: %s (max %d characters)", path, MaxFilePathLength), // Moved from unnamed argument
		)
	}

	// Must be relative path (changed from absolute)
	if filepath.IsAbs(path) {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"Path must be relative, not absolute", // Already correct
			"path", path,
			"hint", "Use relative paths like 'src/main.go' instead of '/home/user/project/src/main.go'",
		)
	}

	// Clean the path to resolve any .. or . elements
	cleanPath := filepath.Clean(path)

	// Convert relative path to absolute for file system checks
	absPath := cleanPath
	if !filepath.IsAbs(cleanPath) {
		cwd, err := os.Getwd()
		if err != nil {
			return models.NewDelegateError(
				models.ErrorTypeInvalidRequest,
				"Failed to get working directory", // Already correct
				"error", err.Error(),
			)
		}
		absPath = filepath.Join(cwd, cleanPath)
	}

	// Check file exists and is readable
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return models.NewDelegateError(
				models.ErrorTypeFileNotFound,
				fmt.Sprintf("File not found: %s", path), // Moved from unnamed argument
			)
		}
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			fmt.Sprintf("Cannot access file: %s", path), // Moved from unnamed argument
		)
	}

	// Must be a regular file
	if !info.Mode().IsRegular() {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			fmt.Sprintf("Not a regular file: %s", path), // Moved from unnamed argument
		)
	}

	// Check file size (max 1MB per file)
	if info.Size() > 1024*1024 {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			fmt.Sprintf("File too large: %s (max 1MB)", cleanPath), // Moved from unnamed argument
		)
	}

	return nil
}

// ValidatePrompt validates the prompt
func ValidatePrompt(prompt string) error {
	if prompt == "" {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"Prompt is required", // Moved from unnamed argument
		)
	}

	if len(prompt) > MaxPromptSize {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			fmt.Sprintf("Prompt too large: %d bytes (max %d)", len(prompt), MaxPromptSize), // Moved from unnamed argument
		)
	}

	return nil
}

// ValidateModel validates the model name
func ValidateModel(model string) error {
	if model == "" {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"Model is required", // Moved from unnamed argument
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
			"Timeout cannot be negative", // Moved from unnamed argument
		)
	}

	if timeout > 0 && timeout < MinTimeout {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			fmt.Sprintf("Timeout too short: %d seconds (minimum: %d seconds, recommended: %d+ seconds for stable LLM responses)", 
				timeout, MinTimeout, RecommendedMinTimeout),
		)
	}

	if timeout > MaxTimeout {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			fmt.Sprintf("Timeout too long: %d seconds (max %d)", timeout, MaxTimeout), // Moved from unnamed argument
		)
	}

	return nil
}

// ValidateMaxTokens validates max_tokens parameter
func ValidateMaxTokens(maxTokens int) error {
	if maxTokens < 0 {
		return models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			"Max tokens cannot be negative", // Moved from unnamed argument
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
			fmt.Sprintf("Invalid extract option: %s (must be 'all', 'code', or 'explanation')", extract), // Moved from unnamed argument
		)
	}
}