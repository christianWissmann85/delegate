package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/christianwissmann85/delegate/internal/models"
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
	// Validate request
	if err := ValidateOutputID(req.OutputID); err != nil {
		return nil, err
	}

	// Set default extract option
	if req.Options.Extract == "" {
		req.Options.Extract = "all"
	}

	// Validate extract option
	if err := ValidateExtractOption(req.Options.Extract); err != nil {
		return nil, err
	}

	// Get output from storage
	output, err := h.storage.Get(req.OutputID)
	if err != nil {
		return nil, models.NewDelegateError(
			models.ErrorTypeNotFound,
			"",
			fmt.Sprintf("output not found: %v", err),
		)
	}

	// Extract requested content
	var content string
	var fileContent string // Content specifically prepared for file writing

	switch req.Options.Extract {
	case "all":
		content = output.Response.Raw
		fileContent = output.Response.Raw
	case "code":
		content = h.extractCodeContent(output)
		// For file writing, use the file-specific extraction
		if req.Options.WriteTo != "" {
			fileContent = h.extractCodeForFile(output, req.Options.WriteTo)
		} else {
			fileContent = content
		}
	case "explanation":
		content = output.Response.Extracted.Explanation
		fileContent = output.Response.Extracted.Explanation
	}

	// Track if content was truncated
	truncated := false
	originalLength := len(content)

	// Validate and apply token limit if specified
	if req.Options.MaxTokens > 0 {
		if err := ValidateMaxTokens(req.Options.MaxTokens); err != nil {
			return nil, err
		}
		content = h.truncateContent(content, req.Options.MaxTokens)
		fileContent = h.truncateContent(fileContent, req.Options.MaxTokens)
		truncated = len(content) < originalLength
	}

	// If WriteTo is specified, write to file instead of returning content
	if req.Options.WriteTo != "" {
		if err := h.writeToFile(req.Options.WriteTo, fileContent); err != nil {
			return nil, models.NewDelegateError(
				models.ErrorTypeInternal,
				"",
				fmt.Sprintf("failed to write to file: %v", err),
			)
		}

		// Calculate file size and tokens saved
		fileSize := len(fileContent)
		fileSizeKB := float64(fileSize) / 1024.0
		tokensSaved := fileSize / 4 // Approximate: 1 token ≈ 4 characters

		// Format file size nicely
		var sizeStr string
		if fileSizeKB < 1 {
			sizeStr = fmt.Sprintf("%d bytes", fileSize)
		} else if fileSizeKB < 1024 {
			sizeStr = fmt.Sprintf("%.1f KB", fileSizeKB)
		} else {
			sizeStr = fmt.Sprintf("%.1f MB", fileSizeKB/1024.0)
		}

		// Return success message with size and tokens saved
		return &ReadResponse{
			Content:     fmt.Sprintf("Content written to %s (%s, ~%d tokens saved)", req.Options.WriteTo, sizeStr, tokensSaved),
			Truncated:   truncated,
			Tokens:      0, // No tokens returned when writing to file
			Extraction:  req.Options.Extract,
			FileWritten: true,
		}, nil
	}

	// Calculate approximate token count for returned content
	// Using same approximation as truncateContent: 1 token ≈ 4 characters
	tokenCount := len(content) / 4

	return &ReadResponse{
		Content:    content,
		Truncated:  truncated,
		Tokens:     tokenCount,
		Extraction: req.Options.Extract,
	}, nil
}

// ReadRequest represents the read tool parameters
type ReadRequest struct {
	OutputID string      `json:"output_id"`
	Options  ReadOptions `json:"options,omitempty"`
}

// ReadOptions configures what to read
type ReadOptions struct {
	Extract   string `json:"extract,omitempty"`    // "all", "code", "explanation"
	MaxTokens int    `json:"max_tokens,omitempty"` // Limit response size
	WriteTo   string `json:"write_to,omitempty"`   // Write content to file instead of returning
}

// ReadResponse represents the read tool response
type ReadResponse struct {
	Content     string `json:"content"`
	Truncated   bool   `json:"truncated"`
	Tokens      int    `json:"tokens"`
	Extraction  string `json:"extraction"`
	Language    string `json:"language,omitempty"`
	FileWritten bool   `json:"file_written,omitempty"`
}

// extractCodeContent formats all code blocks into a single string
func (h *ReadHandler) extractCodeContent(output *models.Output) string {
	if len(output.Response.Extracted.Code) == 0 {
		return ""
	}

	var parts []string
	for i, block := range output.Response.Extracted.Code {
		// Format as fenced code block
		fence := fmt.Sprintf("```%s\n%s\n```", block.Language, block.Content)
		parts = append(parts, fence)

		// Add separator between blocks (except last)
		if i < len(output.Response.Extracted.Code)-1 {
			parts = append(parts, "")
		}
	}

	return strings.Join(parts, "\n")
}

// extractCodeForFile extracts code suitable for saving to a file
func (h *ReadHandler) extractCodeForFile(output *models.Output, filePath string) string {
	if len(output.Response.Extracted.Code) == 0 {
		// No code blocks found, return raw response as fallback
		return output.Response.Raw
	}

	// Detect if this is a source code file or documentation
	if h.isDocumentationFile(filePath) {
		// For documentation files, keep the fences
		return h.extractCodeContent(output)
	}

	// For source code files, concatenate raw code content without fences
	var parts []string
	for _, block := range output.Response.Extracted.Code {
		parts = append(parts, block.Content)
	}

	// Join blocks with single newline
	result := strings.Join(parts, "\n")

	// Clean up any stray markdown artifacts
	result = h.cleanupCodeArtifacts(result)

	return result
}

// isDocumentationFile checks if the file is a documentation file that should preserve markdown
func (h *ReadHandler) isDocumentationFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	docExtensions := []string{
		".md", ".markdown", ".rst", ".txt", ".adoc", ".asciidoc",
		".textile", ".rdoc", ".org", ".creole", ".mediawiki",
		".wiki", ".pod", ".rmd", ".mkd", ".mkdn", ".mdwn", ".mdown",
	}

	for _, docExt := range docExtensions {
		if ext == docExt {
			return true
		}
	}

	// Also check for common documentation filenames
	baseName := strings.ToLower(filepath.Base(filePath))
	docNames := []string{"readme", "changelog", "history", "license", "notice", "authors", "contributors", "todo"}
	for _, docName := range docNames {
		if strings.HasPrefix(baseName, docName) {
			return true
		}
	}

	return false
}

// cleanupCodeArtifacts removes common markdown artifacts from code
func (h *ReadHandler) cleanupCodeArtifacts(content string) string {
	// Remove leading/trailing backticks that might have slipped through
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")

	// Remove language identifiers at the start
	lines := strings.Split(content, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		// Common language identifiers that might appear
		langs := []string{
			"go", "golang", "python", "py", "python3", "javascript", "js", "typescript", "ts",
			"java", "cpp", "c++", "c", "csharp", "cs", "c#", "dart", "flutter", "rust", "rs",
			"ruby", "rb", "php", "swift", "kotlin", "kt", "scala", "r", "R", "julia", "jl",
			"bash", "sh", "shell", "zsh", "fish", "powershell", "ps1", "batch", "bat", "cmd",
			"sql", "mysql", "postgresql", "sqlite", "yaml", "yml", "json", "xml", "html",
			"css", "scss", "sass", "less", "vue", "jsx", "tsx", "svelte", "elm", "purescript",
			"haskell", "hs", "erlang", "erl", "elixir", "ex", "exs", "clojure", "clj", "cljs",
			"lisp", "scheme", "racket", "ocaml", "ml", "fsharp", "fs", "nim", "zig", "v", "vlang",
			"pascal", "delphi", "fortran", "f90", "cobol", "ada", "lua", "perl", "pl", "awk",
			"matlab", "octave", "mathematica", "maple", "sage", "prolog", "smalltalk", "forth",
			"assembly", "asm", "nasm", "masm", "gas", "llvm", "wasm", "webassembly",
			"terraform", "tf", "dockerfile", "docker", "makefile", "make", "cmake", "gradle",
			"maven", "ant", "bazel", "ninja", "meson", "scons", "rake", "gemfile", "cargo",
			"toml", "ini", "conf", "config", "properties", "env", "dotenv", "gitignore",
			"editorconfig", "eslintrc", "prettierrc", "babelrc", "webpack", "vite", "rollup",
			"proto", "protobuf", "graphql", "gql", "solidity", "sol", "vyper", "vy",
			"cuda", "cu", "opencl", "cl", "glsl", "hlsl", "metal", "msl", "wgsl",
		}

		lowerFirst := strings.ToLower(firstLine)
		for _, lang := range langs {
			if lowerFirst == lang || lowerFirst == "```"+lang || firstLine == lang || firstLine == "```"+lang {
				lines = lines[1:]
				break
			}
		}
	}

	// Also clean up any trailing ``` on its own line
	if len(lines) > 0 {
		lastIdx := len(lines) - 1
		if strings.TrimSpace(lines[lastIdx]) == "```" {
			lines = lines[:lastIdx]
		}
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// truncateContent truncates content to approximately maxTokens
func (h *ReadHandler) truncateContent(content string, maxTokens int) string {
	// Simple approximation: 1 token ≈ 4 characters
	// This is a rough estimate; actual tokenization varies by model
	maxChars := maxTokens * 4

	// Ensure we have a minimum reasonable size to avoid panic
	if maxChars < 10 {
		maxChars = 10
	}

	if len(content) <= maxChars {
		return content
	}

	// Ensure we have enough space for ellipsis
	if maxChars <= 3 {
		return "..."
	}

	// Truncate and add ellipsis
	truncated := content[:maxChars-3] + "..."

	// Try to break at a word boundary
	lastSpace := strings.LastIndexAny(truncated[:len(truncated)-3], " \n\t")
	if lastSpace > maxChars*3/4 { // Only break at word if we're keeping at least 75% of content
		truncated = truncated[:lastSpace] + "..."
	}

	return truncated
}

// writeToFile writes content to the specified file path with security checks
func (h *ReadHandler) writeToFile(filePath string, content string) error {
	// Clean and validate the path to prevent path traversal attacks
	cleanPath := filepath.Clean(filePath)

	// Reject paths that try to go outside current directory
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("invalid file path: path traversal detected")
	}

	// Convert to absolute path for validation
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Ensure the path is within the current working directory
	if !strings.HasPrefix(absPath, cwd) {
		return fmt.Errorf("invalid file path: must be within current directory")
	}

	// Ensure the directory exists
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the file
	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
