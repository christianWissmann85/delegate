package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/christianwissmann85/delegate/internal/models"
)

const (
	// Maximum file size: 1MB per file
	MaxFileSize = 1024 * 1024
	// Maximum total size for all files: 5MB
	MaxTotalFileSize = 5 * 1024 * 1024
)

// FileContent represents the content of a file
type FileContent struct {
	Path    string
	Content string
}

// ReadFilesWithLimit reads files with memory limits
func ReadFilesWithLimit(filePaths []string) ([]FileContent, error) {
	if len(filePaths) == 0 {
		return nil, nil
	}

	// Validate paths first
	if err := ValidateFilePaths(filePaths); err != nil {
		return nil, err
	}

	var totalSize int64
	var files []FileContent

	for _, path := range filePaths {
		// Clean the path
		cleanPath := filepath.Clean(path)
		
		// Convert relative path to absolute for file operations
		absPath := cleanPath
		if !filepath.IsAbs(cleanPath) {
			cwd, err := os.Getwd()
			if err != nil {
				return nil, models.NewDelegateError(
					models.ErrorTypeInvalidRequest,
					"Failed to get working directory",
					"error", err.Error(),
				)
			}
			absPath = filepath.Join(cwd, cleanPath)
		}

		// Open file with size limit
		// Pass both absolute path (for file operations) and original path (for error messages)
		content, size, err := readFileWithLimit(absPath, path, MaxFileSize)
		if err != nil {
			return nil, err
		}

		// Check total size limit
		totalSize += size
		if totalSize > MaxTotalFileSize {
			return nil, models.NewDelegateError(
				models.ErrorTypeInvalidRequest,
				fmt.Sprintf("Total file size exceeds limit: %d bytes (max %d)", totalSize, MaxTotalFileSize),
			)
		}

		files = append(files, FileContent{
			Path:    path, // Store original relative path
			Content: content,
		})
	}

	return files, nil
}

// readFileWithLimit reads a file with size limit
// absPath is used for file operations, originalPath is used for error messages
func readFileWithLimit(absPath string, originalPath string, maxSize int64) (string, int64, error) {
	file, err := os.Open(absPath)
	if err != nil {
		return "", 0, models.NewDelegateError(
			models.ErrorTypeFileNotFound,
			fmt.Sprintf("Cannot open file: %s", originalPath),
		)
	}
	defer func() {
		_ = file.Close()
	}()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		return "", 0, models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			fmt.Sprintf("Cannot stat file: %s", originalPath),
		)
	}

	// Check size before reading
	if info.Size() > maxSize {
		return "", 0, models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			fmt.Sprintf("File too large: %s (%d bytes, max %d)", originalPath, info.Size(), maxSize),
		)
	}

	// Use a limited reader to prevent memory exhaustion
	limitedReader := io.LimitReader(file, maxSize+1) // +1 to detect if file grew
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", 0, models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			fmt.Sprintf("Error reading file: %s", originalPath),
		)
	}

	// Check if we hit the limit (file grew while reading)
	if int64(len(content)) > maxSize {
		return "", 0, models.NewDelegateError(
			models.ErrorTypeInvalidRequest,
			fmt.Sprintf("File too large: %s (exceeded %d bytes while reading)", originalPath, maxSize),
		)
	}

	return string(content), int64(len(content)), nil
}

// BuildPromptWithFiles builds a prompt with file contents
func BuildPromptWithFiles(prompt string, files []FileContent) string {
	if len(files) == 0 {
		return prompt
	}

	var builder strings.Builder
	builder.WriteString(prompt)

	for _, file := range files {
		// Get just the filename for display
		filename := filepath.Base(file.Path)

		// Detect file type for better formatting
		ext := strings.ToLower(filepath.Ext(filename))

		// Add file content with appropriate formatting
		builder.WriteString(fmt.Sprintf("\n\n--- File: %s ---\n", filename))

		// For code files, wrap in code blocks
		if isCodeFile(ext) {
			lang := getLanguageFromExt(ext)
			builder.WriteString(fmt.Sprintf("```%s\n%s\n```", lang, file.Content))
		} else {
			builder.WriteString(file.Content)
		}
	}

	return builder.String()
}

// isCodeFile checks if the file extension indicates a code file
func isCodeFile(ext string) bool {
	codeExts := map[string]bool{
		".go":    true,
		".py":    true,
		".js":    true,
		".ts":    true,
		".java":  true,
		".c":     true,
		".cpp":   true,
		".h":     true,
		".hpp":   true,
		".rs":    true,
		".rb":    true,
		".php":   true,
		".swift": true,
		".kt":    true,
		".cs":    true,
		".sh":    true,
		".bash":  true,
		".sql":   true,
		".r":     true,
		".m":     true,
		".scala": true,
		".lua":   true,
		".pl":    true,
		".json":  true,
		".xml":   true,
		".yaml":  true,
		".yml":   true,
		".toml":  true,
	}
	return codeExts[ext]
}

// getLanguageFromExt maps file extension to language name
func getLanguageFromExt(ext string) string {
	langMap := map[string]string{
		".go":    "go",
		".py":    "python",
		".js":    "javascript",
		".ts":    "typescript",
		".java":  "java",
		".c":     "c",
		".cpp":   "cpp",
		".h":     "c",
		".hpp":   "cpp",
		".rs":    "rust",
		".rb":    "ruby",
		".php":   "php",
		".swift": "swift",
		".kt":    "kotlin",
		".cs":    "csharp",
		".sh":    "bash",
		".bash":  "bash",
		".sql":   "sql",
		".r":     "r",
		".m":     "matlab",
		".scala": "scala",
		".lua":   "lua",
		".pl":    "perl",
		".json":  "json",
		".xml":   "xml",
		".yaml":  "yaml",
		".yml":   "yaml",
		".toml":  "toml",
	}
	if lang, ok := langMap[ext]; ok {
		return lang
	}
	return "text"
}
