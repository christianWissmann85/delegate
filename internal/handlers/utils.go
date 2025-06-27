package handlers

import (
	"fmt"
	"path/filepath"
	"strings"
	
	"github.com/christianwissmann85/delegate/internal/models"
)



// TruncateContent truncates content to approximately the specified number of tokens
func TruncateContent(content string, maxTokens int) string {
	// Rough estimation: ~4 characters per token
	maxChars := maxTokens * 4
	
	if len(content) <= maxChars {
		return content
	}
	
	// Try to truncate at a natural boundary (newline)
	truncated := content[:maxChars]
	lastNewline := strings.LastIndex(truncated, "\n")
	if lastNewline > maxChars/2 { // Only use newline if it's not too far back
		return content[:lastNewline] + "\n\n[Content truncated]"
	}
	
	return truncated + "\n\n[Content truncated]"
}

// IsDocumentationFile checks if a file extension indicates a documentation file
func IsDocumentationFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	docExtensions := []string{
		".md", ".markdown", ".rst", ".txt", ".doc", ".docx",
		".adoc", ".asciidoc", ".org", ".tex", ".latex",
	}
	
	for _, docExt := range docExtensions {
		if ext == docExt {
			return true
		}
	}
	
	return false
}

// CleanupCodeArtifacts removes common artifacts from code content
func CleanupCodeArtifacts(content, language string) string {
	// Remove trailing whitespace from each line
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	content = strings.Join(lines, "\n")
	
	// Remove multiple consecutive blank lines
	for strings.Contains(content, "\n\n\n") {
		content = strings.Replace(content, "\n\n\n", "\n\n", -1)
	}
	
	// Trim leading and trailing whitespace
	content = strings.TrimSpace(content)
	
	return content
}

// ExtractCodeContent extracts only code blocks from content
func ExtractCodeContent(blocks []models.CodeBlock) string {
	if len(blocks) == 0 {
		return ""
	}
	
	var result strings.Builder
	for i, block := range blocks {
		if i > 0 {
			result.WriteString("\n\n")
		}
		
		// Add fence with language
		result.WriteString(fmt.Sprintf("```%s\n", block.Language))
		result.WriteString(block.Content)
		if !strings.HasSuffix(block.Content, "\n") {
			result.WriteString("\n")
		}
		result.WriteString("```")
	}
	
	return result.String()
}

// ExtractCodeForFile extracts code suitable for writing to a file
func ExtractCodeForFile(blocks []models.CodeBlock, language string) string {
	// If language is specified, only include blocks of that language
	var relevantBlocks []models.CodeBlock
	if language != "" {
		for _, block := range blocks {
			if block.Language == language {
				relevantBlocks = append(relevantBlocks, block)
			}
		}
	} else {
		relevantBlocks = blocks
	}
	
	if len(relevantBlocks) == 0 {
		return ""
	}
	
	// For file writing, we concatenate code blocks without fences
	var result strings.Builder
	for i, block := range relevantBlocks {
		if i > 0 {
			result.WriteString("\n\n")
		}
		result.WriteString(block.Content)
	}
	
	return result.String()
}