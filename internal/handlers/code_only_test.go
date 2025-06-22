package handlers_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/christianwissmann85/delegate/internal/handlers"
	"github.com/christianwissmann85/delegate/internal/models"
	"github.com/christianwissmann85/delegate/internal/storage"
)

// createTestOutputWithCodeOnly is a helper to create a models.Output for testing
// It allows specifying the CodeOnly metadata and multiple code blocks.
func createTestOutputWithCodeOnly(id string, codeOnly bool, codeBlocks []models.ExtractedCode, explanation string) *models.Output {
	// The Raw content is primarily used by "extract: all" and not directly by "extract: code"
	// However, it's good practice to make it consistent if possible.
	// For these tests, we're focusing on the Extracted.Code part.
	rawContent := ""
	if explanation != "" {
		rawContent += explanation + "\n\n"
	}
	for _, block := range codeBlocks {
		rawContent += fmt.Sprintf("```%s\n%s\n```\n", block.Language, block.Content)
	}
	rawContent = strings.TrimSpace(rawContent)

	return &models.Output{
		ID:        id,
		Model:     "mock-test",
		Prompt:    "Test prompt for code_only functionality",
		CreatedAt: time.Now(),
		Response: models.Response{
			Raw: rawContent,
			Extracted: models.Extracted{
				Code:        codeBlocks,
				Explanation: explanation,
			},
		},
		Metadata: models.Metadata{
			TotalBytes:       int64(len(rawContent)),
			EstimatedTokens:  handlers.EstimateTokens(rawContent),
			ProcessingTimeMs: 10,
			CodeOnly:         codeOnly, // This is the key metadata field for this test
		},
	}
}

func TestReadHandler_CodeOnlyFunctionality(t *testing.T) {
	ctx := context.Background()
	tempDir := t.TempDir() // Create a temporary directory for storage
	store, err := storage.NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	readHandler := handlers.NewReadHandler(store)

	// Define common code blocks for reuse across tests
	singlePythonCode := []models.ExtractedCode{
		{Language: "python", Content: "def hello():\n    print('Hello, world!')"},
	}
	multipleCodeBlocks := []models.ExtractedCode{
		{Language: "go", Content: "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello from Go!\")\n}"},
		{Language: "javascript", Content: "function greet() {\n    console.log('Hello from JS!');\n}"},
	}
	explanationText := "This is a simple explanation for the code."

	tests := []struct {
		name                 string
		outputID             string
		codeOnlyMetadata     bool
		extractedCodeBlocks  []models.ExtractedCode
		explanation          string
		writeToPath          string // Optional: if specified, content will be written to this file
		expectedContentCheck func(t *testing.T, content string)
		expectedFileCheck    func(t *testing.T, filePath string) // Only used if writeToPath is set
	}{
		// --- Test cases for ReadResponse.Content (no write_to) ---
		{
			name:                "1. CodeOnly=true: Single code block, no fences in response",
			outputID:            "out_codeonly_true_single",
			codeOnlyMetadata:    true,
			extractedCodeBlocks: singlePythonCode,
			explanation:         explanationText,
			expectedContentCheck: func(t *testing.T, content string) {
				if strings.Contains(content, "```") {
					t.Errorf("Expected no markdown fences in content, got:\n%s", content)
				}
				if !strings.Contains(content, singlePythonCode[0].Content) {
					t.Errorf("Expected content to contain raw code, got:\n%s", content)
				}
				if content != singlePythonCode[0].Content {
					t.Errorf("Content mismatch. Expected:\n%q\nGot:\n%q", singlePythonCode[0].Content, content)
				}
			},
		},
		{
			name:                "2. CodeOnly=false: Single code block, with fences in response",
			outputID:            "out_codeonly_false_single",
			codeOnlyMetadata:    false,
			extractedCodeBlocks: singlePythonCode,
			explanation:         explanationText,
			expectedContentCheck: func(t *testing.T, content string) {
				expectedFencedCode := "```python\n" + singlePythonCode[0].Content + "\n```"
				if !strings.Contains(content, expectedFencedCode) {
					t.Errorf("Expected content with markdown fences and language, got:\n%s", content)
				}
				if content != expectedFencedCode {
					t.Errorf("Content mismatch. Expected:\n%q\nGot:\n%q", expectedFencedCode, content)
				}
			},
		},
		{
			name:                "3. CodeOnly=true: Multiple code blocks, no fences in response",
			outputID:            "out_codeonly_true_multi",
			codeOnlyMetadata:    true,
			extractedCodeBlocks: multipleCodeBlocks,
			explanation:         explanationText,
			expectedContentCheck: func(t *testing.T, content string) {
				if strings.Contains(content, "```") {
					t.Errorf("Expected no markdown fences in content, got:\n%s", content)
				}
				// Expect blocks joined by a newline, with an empty line between them as per extractCodeContent logic
				expected := strings.Join([]string{multipleCodeBlocks[0].Content, "", multipleCodeBlocks[1].Content}, "\n")
				if content != expected {
					t.Errorf("Content mismatch for multiple blocks. Expected:\n%q\nGot:\n%q", expected, content)
				}
			},
		},
		{
			name:                "4. CodeOnly=false: Multiple code blocks, with fences in response",
			outputID:            "out_codeonly_false_multi",
			codeOnlyMetadata:    false,
			extractedCodeBlocks: multipleCodeBlocks,
			explanation:         explanationText,
			expectedContentCheck: func(t *testing.T, content string) {
				expectedGo := fmt.Sprintf("```%s\n%s\n```", multipleCodeBlocks[0].Language, multipleCodeBlocks[0].Content)
				expectedJS := fmt.Sprintf("```%s\n%s\n```", multipleCodeBlocks[1].Language, multipleCodeBlocks[1].Content)
				// Expect blocks joined by a newline, with an empty line between them as per extractCodeContent logic
				expected := strings.Join([]string{expectedGo, "", expectedJS}, "\n")
				if !strings.Contains(content, expectedGo) || !strings.Contains(content, expectedJS) {
					t.Errorf("Expected content with markdown fences for multiple blocks, got:\n%s", content)
				}
				if content != expected {
					t.Errorf("Content mismatch for multiple blocks with fences. Expected:\n%q\nGot:\n%q", expected, content)
				}
			},
		},
		// --- WriteTo Integration Tests ---
		{
			name:                "5. WriteTo (doc file) with CodeOnly=true: No fences in file",
			outputID:            "out_write_doc_codeonly_true",
			codeOnlyMetadata:    true,
			extractedCodeBlocks: singlePythonCode,
			explanation:         explanationText,
			writeToPath:         "output_codeonly_true.md", // Documentation file - relative path
			expectedContentCheck: func(t *testing.T, content string) {
				if !strings.Contains(content, "Content written to") {
					t.Errorf("Expected success message for write_to, got: %s", content)
				}
				if !strings.Contains(content, "output_codeonly_true.md") {
					t.Errorf("Success message missing file path: %s", content)
				}
			},
			expectedFileCheck: func(t *testing.T, filePath string) {
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read file %s: %v", filePath, err)
				}
				contentStr := string(fileContent)
				if strings.Contains(contentStr, "```") {
					t.Errorf("Expected no markdown fences in documentation file %s (CodeOnly=true), got:\n%s", filePath, contentStr)
				}
				if contentStr != singlePythonCode[0].Content {
					t.Errorf("File content mismatch. Expected:\n%q\nGot:\n%q", singlePythonCode[0].Content, contentStr)
				}
			},
		},
		{
			name:                "6. WriteTo (doc file) with CodeOnly=false: Fences in file",
			outputID:            "out_write_doc_codeonly_false",
			codeOnlyMetadata:    false,
			extractedCodeBlocks: singlePythonCode,
			explanation:         explanationText,
			writeToPath:         "output_codeonly_false.md", // Documentation file - relative path
			expectedContentCheck: func(t *testing.T, content string) {
				if !strings.Contains(content, "Content written to") {
					t.Errorf("Expected success message for write_to, got: %s", content)
				}
			},
			expectedFileCheck: func(t *testing.T, filePath string) {
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read file %s: %v", filePath, err)
				}
				contentStr := string(fileContent)
				expectedFencedCode := "```python\n" + singlePythonCode[0].Content + "\n```"
				if !strings.Contains(contentStr, expectedFencedCode) {
					t.Errorf("Expected markdown fences in documentation file %s (CodeOnly=false), got:\n%s", filePath, contentStr)
				}
				if contentStr != expectedFencedCode {
					t.Errorf("File content mismatch. Expected:\n%q\nGot:\n%q", expectedFencedCode, contentStr)
				}
			},
		},
		{
			name:                "7. WriteTo (source file) with CodeOnly=true: No fences in file (extractCodeForFile overrides)",
			outputID:            "out_write_src_codeonly_true",
			codeOnlyMetadata:    true,
			extractedCodeBlocks: singlePythonCode,
			explanation:         explanationText,
			writeToPath:         "output_codeonly_true.py", // Source code file - relative path
			expectedContentCheck: func(t *testing.T, content string) {
				if !strings.Contains(content, "Content written to") {
					t.Errorf("Expected success message for write_to, got: %s", content)
				}
			},
			expectedFileCheck: func(t *testing.T, filePath string) {
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read file %s: %v", filePath, err)
				}
				contentStr := string(fileContent)
				if strings.Contains(contentStr, "```") {
					t.Errorf("Expected no markdown fences in source file %s, got:\n%s", filePath, contentStr)
				}
				if contentStr != singlePythonCode[0].Content {
					t.Errorf("File content mismatch. Expected:\n%q\nGot:\n%q", singlePythonCode[0].Content, contentStr)
				}
			},
		},
		{
			name:                "8. WriteTo (source file) with CodeOnly=false: No fences in file (extractCodeForFile overrides)",
			outputID:            "out_write_src_codeonly_false",
			codeOnlyMetadata:    false,
			extractedCodeBlocks: singlePythonCode,
			explanation:         explanationText,
			writeToPath:         "output_codeonly_false.py", // Source code file - relative path
			expectedContentCheck: func(t *testing.T, content string) {
				if !strings.Contains(content, "Content written to") {
					t.Errorf("Expected success message for write_to, got: %s", content)
				}
			},
			expectedFileCheck: func(t *testing.T, filePath string) {
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("Failed to read file %s: %v", filePath, err)
				}
				contentStr := string(fileContent)
				if strings.Contains(contentStr, "```") {
					t.Errorf("Expected no markdown fences in source file %s, got:\n%s", filePath, contentStr)
				}
				if contentStr != singlePythonCode[0].Content {
					t.Errorf("File content mismatch. Expected:\n%q\nGot:\n%q", singlePythonCode[0].Content, contentStr)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create and store the test output in the mock storage
			output := createTestOutputWithCodeOnly(tt.outputID, tt.codeOnlyMetadata, tt.extractedCodeBlocks, tt.explanation)
			err := store.Save(output)
			if err != nil {
				t.Fatalf("Failed to save test output %s: %v", tt.outputID, err)
			}

			req := handlers.ReadRequest{
				OutputID: tt.outputID,
				Options: handlers.ReadOptions{
					Extract: "code", // Always extract code for these tests
				},
			}

			// If writeToPath is specified, add it to the request options
			if tt.writeToPath != "" {
				req.Options.WriteTo = tt.writeToPath
			}

			resp, err := readHandler.Handle(ctx, req)
			if err != nil {
				t.Fatalf("ReadHandler.Handle failed: %v", err)
			}

			// Check the ReadResponse.Content
			tt.expectedContentCheck(t, resp.Content)

			// If writeToPath was specified, perform file-specific checks
			if tt.writeToPath != "" {
				if !resp.FileWritten {
					t.Error("Expected FileWritten to be true when WriteTo is specified")
				}
				tt.expectedFileCheck(t, tt.writeToPath)
				// Clean up the written file
				os.Remove(tt.writeToPath)
			} else {
				// If no writeToPath, FileWritten should be false
				if resp.FileWritten {
					t.Error("Expected FileWritten to be false when WriteTo is not specified")
				}
			}
		})
	}
}

