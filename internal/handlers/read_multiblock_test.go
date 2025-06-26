package handlers

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/christianwissmann85/delegate/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadHandler_MultiBlockHandling(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "delegate_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Change to temp directory for tests
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	err = os.Chdir(tmpDir)
	require.NoError(t, err)
	defer os.Chdir(originalWd)

	tests := []struct {
		name           string
		output         *models.Output
		request        ReadRequest
		wantContent    string
		wantMultiBlock bool
		wantBlockCount int
		wantError      bool
		errorContains  string
	}{
		{
			name: "multiple blocks without block_index shows warning",
			output: &models.Output{
				Response: models.Response{
					Extracted: models.Extracted{
						Code: []models.ExtractedCode{
							{
								Language: "javascript",
								Content:  "const Component = () => {\n  return <div>Hello</div>;\n}",
							},
							{
								Language: "css",
								Content:  ".container {\n  display: flex;\n}",
							},
						},
					},
				},
			},
			request: ReadRequest{
				OutputID: "test_output_001",
				Options: ReadOptions{
					Extract: "code",
					WriteTo: "test.js",
				},
			},
			wantContent:    "Warning: Multiple code blocks found (2 blocks). Use block_index option to select specific block.\n\nBlock 0: javascript - \"const Component = () => {\" (54 bytes, 3 lines)\nBlock 1: css - \".container {\" (31 bytes, 3 lines)\n",
			wantMultiBlock: true,
			wantBlockCount: 2,
		},
		{
			name: "multiple blocks with valid block_index writes specific block",
			output: &models.Output{
				Response: models.Response{
					Extracted: models.Extracted{
						Code: []models.ExtractedCode{
							{
								Language: "go",
								Content:  "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}",
							},
							{
								Language: "go",
								Content:  "package main_test\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {}",
							},
						},
					},
				},
			},
			request: ReadRequest{
				OutputID: "test_output_002",
				Options: ReadOptions{
					Extract:    "code",
					WriteTo:    "main.go",
					BlockIndex: intPtr(0),
				},
			},
			wantContent: "Content written to main.go",
			wantError:   false,
		},
		{
			name: "block_index out of range returns error",
			output: &models.Output{
				Response: models.Response{
					Extracted: models.Extracted{
						Code: []models.ExtractedCode{
							{
								Language: "python",
								Content:  "def hello():\n    print(\"Hello\")",
							},
							{
								Language: "python",
								Content:  "def test():\n    pass",
							},
						},
					},
				},
			},
			request: ReadRequest{
				OutputID: "test_output_003",
				Options: ReadOptions{
					Extract:    "code",
					WriteTo:    "test.py",
					BlockIndex: intPtr(5),
				},
			},
			wantError:     true,
			errorContains: "block_index 5 out of range (0-1)",
		},
		{
			name: "single block proceeds without warning",
			output: &models.Output{
				Response: models.Response{
					Extracted: models.Extracted{
						Code: []models.ExtractedCode{
							{
								Language: "ruby",
								Content:  "puts 'Hello World'",
							},
						},
					},
				},
			},
			request: ReadRequest{
				OutputID: "test_output_004",
				Options: ReadOptions{
					Extract: "code",
					WriteTo: "test.rb",
				},
			},
			wantContent: "Content written to test.rb",
			wantError:   false,
		},
		{
			name: "extract all bypasses multi-block logic",
			output: &models.Output{
				Response: models.Response{
					Raw: "# Documentation\n\n```js\ncode1\n```\n\n```css\ncode2\n```",
					Extracted: models.Extracted{
						Code: []models.ExtractedCode{
							{Language: "js", Content: "code1"},
							{Language: "css", Content: "code2"},
						},
					},
				},
			},
			request: ReadRequest{
				OutputID: "test_output_005",
				Options: ReadOptions{
					Extract: "all",
					WriteTo: "doc.md",
				},
			},
			wantContent: "Content written to doc.md",
			wantError:   false,
		},
		{
			name: "long first line gets truncated",
			output: &models.Output{
				Response: models.Response{
					Extracted: models.Extracted{
						Code: []models.ExtractedCode{
							{
								Language: "javascript",
								Content:  "const veryLongVariableNameThatExceedsSixtyCharactersForTestingTruncation = 'This is a very long line that should be truncated';",
							},
							{
								Language: "javascript",
								Content:  "short line",
							},
						},
					},
				},
			},
			request: ReadRequest{
				OutputID: "test_output_006",
				Options: ReadOptions{
					Extract: "code",
					WriteTo: "test.js",
				},
			},
			wantContent:    "Warning: Multiple code blocks found (2 blocks). Use block_index option to select specific block.\n\nBlock 0: javascript - \"const veryLongVariableNameThatExceedsSixtyCharactersForTe...\" (127 bytes, 1 lines)\nBlock 1: javascript - \"short line\" (10 bytes, 1 lines)\n",
			wantMultiBlock: true,
			wantBlockCount: 2,
		},
		{
			name: "size formatting for large blocks",
			output: &models.Output{
				Response: models.Response{
					Extracted: models.Extracted{
						Code: []models.ExtractedCode{
							{
								Language: "text",
								Content:  "// This is a large file with lots of content...\n" + string(make([]byte, 2000)),
							},
							{
								Language: "text",
								Content:  "// Smaller file\n" + string(make([]byte, 496)),
							},
						},
					},
				},
			},
			request: ReadRequest{
				OutputID: "test_output_007",
				Options: ReadOptions{
					Extract: "code",
					WriteTo: "test.txt",
				},
			},
			wantContent:    "Warning: Multiple code blocks found (2 blocks). Use block_index option to select specific block.\n\nBlock 0: text - \"// This is a large file with lots of content...\" (2.0 KB, 2 lines)\nBlock 1: text - \"// Smaller file\" (512 bytes, 2 lines)\n",
			wantMultiBlock: true,
			wantBlockCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock storage
			mockStorage := &mockReadStorage{
				outputs: map[string]*models.Output{
					tt.request.OutputID: tt.output,
				},
			}

			handler := NewReadHandler(mockStorage)
			readResp, err := handler.Handle(nil, tt.request)

			if tt.wantError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, readResp)
			
			// For write_to operations, check the content message
			if tt.request.Options.WriteTo != "" && !tt.wantMultiBlock {
				assert.Contains(t, readResp.Content, "Content written to")
				assert.Contains(t, readResp.Content, tt.request.Options.WriteTo)
				
				// Verify file was actually written
				filePath := filepath.Join(tmpDir, tt.request.Options.WriteTo)
				if _, err := os.Stat(filePath); err == nil {
					// File exists, check it was written correctly
					assert.True(t, readResp.FileWritten)
				}
			}

			// For multi-block warnings
			if tt.wantMultiBlock {
				assert.Equal(t, tt.wantContent, readResp.Content)
				assert.Equal(t, tt.wantMultiBlock, readResp.MultipleBlocks)
				assert.Equal(t, tt.wantBlockCount, readResp.BlockCount)
			}
		})
	}
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}

// Mock storage for testing
type mockReadStorage struct {
	outputs map[string]*models.Output
}

func (m *mockReadStorage) Save(output *models.Output) error {
	return nil
}

func (m *mockReadStorage) Get(id string) (*models.Output, error) {
	output, ok := m.outputs[id]
	if !ok {
		return nil, models.NewDelegateError(models.ErrorTypeNotFound, "", "not found")
	}
	return output, nil
}

func (m *mockReadStorage) Delete(id string) error {
	return nil
}

func (m *mockReadStorage) ListOlderThan(age time.Duration) ([]string, error) {
	return nil, nil
}