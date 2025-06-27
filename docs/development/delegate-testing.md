# **Delegate: Testing Strategy v1.0**

**Status:** Final | **Version:** 1.0 | **Date:** 2025-06-20

## **1. Guiding Philosophy**

Our testing strategy follows the same core principles as Delegate itself: **simplicity, reliability, and pragmatism.**

The goal of testing is not to achieve 100% code coverage but to gain **maximum confidence** that the system works as expected in real-world MCP scenarios. We focus on tests that verify the interactions between components and the MCP protocol implementation. We avoid writing brittle tests or testing trivial logic.

This is a "No Scope Creep" approach to testing.

## **2. The Testing Pyramid (Simplified)**

We use a simplified testing pyramid that emphasizes the most valuable tests for an MCP server.

```
      /\
     /E2E\      <-- Smallest number, highest value. Real MCP integration
    /------\
   /Integration\  <-- The core of our strategy. Provider & storage tests
  /------------\
 / Unit Tests   \ <-- Used sparingly for complex, pure logic
/----------------\
```

### **Level 1: Unit Tests (The Foundation)**

* **Purpose:** To verify small, isolated pieces of pure, critical logic. We will **not** write unit tests for simple functions.
* **Scope:**
  * **Code Extraction:** Test the logic that parses LLM responses to separate code blocks from explanations
  * **Configuration:** Test that environment variables are loaded and prioritized correctly
  * **MCP Protocol:** Test the tool parameter validation and response formatting
  * **Path Resolution:** Test relative path handling and security validation
  * **Structured Error Handling:** Test the new `AsDelegateError` helper function
* **Tooling:** Go's standard testing package
* **Execution:** `go test ./internal/...`

### **Level 2: Integration Tests (The Core)**

* **Purpose:** To verify that our application modules correctly interact with external systems (filesystem and mocked LLM APIs)
* **Scope:**
  * **Storage Integration:**
    * Test that the storage package correctly creates the `.delegate` directory structure
    * Verify that `delegate_submit_task` correctly writes output files and `delegate_get_output_content` correctly reads them
    * Test the 24-hour cleanup logic
  * **LLM Provider Integration (with Mocks):**
    * Create a **mock LLM provider** that implements the Provider interface
    * Return predictable, pre-defined streaming responses
    * Test the submit task handler's stream handling and error normalization
    * Verify timeout handling and retry logic
  * **MCP Tool Integration:**
    * Test each of the four new tools with various parameter combinations
    * Verify structured JSON error responses for invalid inputs
    * Test the full workflow: submit → metadata → content/write
* **Tooling:** Go's standard testing package with test tables and mock structs
* **Execution:** Part of standard `go test` run - fast and reliable

### **Level 3: End-to-End (E2E) Tests (The Golden Path)**

* **Purpose:** To provide ultimate confirmation that the MCP server works with Claude Code
* **Scope:**
  * Test the MCP server with a mock MCP client that simulates Claude Code
  * Use **real API keys** for actual LLM calls (in CI only)
  * Primary E2E test flows:
    1. **Zero-token workflow:** Submit task → write directly to file
    2. **Review workflow:** Submit task → get metadata → get content → make decision
    3. **Multi-block workflow:** Handle outputs with multiple code blocks
* **Tooling:** Go's testing package with MCP client library
* **Execution:**
  * Separate build-tagged file (`//go:build e2e`)
  * Run manually or in dedicated CI step
  * **Command:** `go test -v --tags=e2e .`

## **3. Test Environment and Automation**

* **Local Development:** 
  * Unit and integration tests run instantly with `go test ./...`
  * E2E tests available with `--tags=e2e` flag and API keys

* **Continuous Integration (CI) with GitHub Actions:**
  1. **On Every Push/PR:** Run all unit and integration tests (fast feedback)
  2. **On Merge to main:** Run E2E tests with real API calls
  3. **API Keys:** Stored securely as GitHub Actions Secrets

## **4. Test File Structure**

Following Go conventions, test files are placed alongside the code they test with a `_test.go` suffix:

```
delegate/
├── testdata/                               # Test fixtures (Go convention for test data)
│   ├── mock_llm_response_code.json        # Mock LLM response with code
│   ├── mock_llm_response_mixed.json       # Mock response with code + explanation
│   ├── mock_stream_chunks.json            # Mock streaming response chunks
│   └── mock_multiblock_response.json      # Mock response with multiple code blocks
│
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   ├── config_test.go                 # Unit: Environment variable loading
│   │   └── languages_test.go              # Unit: Language config loading
│   ├── models/
│   │   ├── models.go
│   │   ├── responses.go
│   │   ├── error.go
│   │   ├── models_test.go                 # Unit: Data structures
│   │   └── error_test.go                  # Unit: AsDelegateError helper
│   ├── extractor/
│   │   ├── extractor.go
│   │   ├── extractor_test.go              # Unit: Code/explanation extraction
│   │   └── patterns_test.go               # Unit: Regex patterns
│   ├── handlers/
│   │   ├── submit_task.go
│   │   ├── submit_task_test.go            # Integration: Submit task handler
│   │   ├── get_metadata.go
│   │   ├── get_metadata_test.go           # Integration: Get metadata handler
│   │   ├── get_content.go
│   │   ├── get_content_test.go            # Integration: Get content handler
│   │   ├── write_file.go
│   │   └── write_file_test.go             # Integration: Write file handler
│   ├── providers/
│   │   ├── mock/                          # Mock provider package
│   │   │   ├── provider.go                # Mock LLM provider implementation
│   │   │   └── provider_test.go           # Unit: Mock provider tests
│   │   ├── google/
│   │   │   ├── client.go
│   │   │   └── client_test.go             # Integration: Gemini provider
│   │   ├── anthropic/
│   │   │   ├── client.go
│   │   │   └── client_test.go             # Integration: Claude provider
│   │   └── factory_test.go                # Unit: Provider factory
│   ├── mcp/
│   │   ├── server.go
│   │   ├── server_test.go                 # Integration: MCP server
│   │   ├── protocol_test.go               # Unit: Protocol handling
│   │   └── tools_test.go                  # Unit: Tool definitions
│   ├── storage/
│   │   ├── store.go
│   │   ├── store_test.go                  # Integration: File operations
│   │   ├── cleanup.go
│   │   └── cleanup_test.go                # Integration: Cleanup routine
│   └── logger/
│       ├── logger.go
│       └── logger_test.go                 # Unit: Logging functionality
│
├── cmd/
│   └── delegate/
│       └── main_test.go                   # Integration: CLI startup
│
└── e2e/
    ├── golden_path_test.go                # E2E: Full MCP workflow (//go:build e2e)
    └── migration_test.go                  # E2E: Old and new API compatibility
```

**Key Points:**
- Test files use `_test.go` suffix and live next to the code they test
- `testdata/` directory contains test fixtures (Go convention)
- Mock provider is properly placed in `internal/providers/mock/`
- E2E tests are in a separate `e2e/` package with build tags
- New `internal/models/` package has dedicated test files

## **5. Testing Workflow**

Our testing strategy focuses on the new 4-tool API workflow while maintaining compatibility testing during the migration period.

### **New API Workflow Testing**

The primary test scenarios cover the new 4-tool architecture:

1. **Submit Task** → **Get Metadata** → **Get Content** (review workflow)
2. **Submit Task** → **Write to File** (zero-token workflow)
3. **Submit Task** → **Get Metadata** → **Write Specific Block** (multi-block workflow)

### **Key Test Scenarios**

#### **Unit Test Examples**

```go
// extractor_test.go
func TestExtractCode(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        wantCode []models.CodeBlock
        wantErr  bool
    }{
        {
            name:  "extracts JavaScript code",
            input: "Here's the code:\n```javascript\nconst x = 1;\n```",
            wantCode: []models.CodeBlock{{Language: "javascript", Content: "const x = 1;"}},
        },
        {
            name:  "handles multiple code blocks",
            input: "```go\nfunc main() {}\n```\n\nExplanation\n\n```md\n# README\n```",
            wantCode: []models.CodeBlock{
                {Language: "go", Content: "func main() {}"},
                {Language: "md", Content: "# README"},
            },
        },
        // More test cases...
    }
}

// models/error_test.go
func TestAsDelegateError(t *testing.T) {
    tests := []struct {
        name     string
        code     string
        message  string
        details  map[string]interface{}
        wantJSON string
    }{
        {
            name:    "basic error",
            code:    "OUTPUT_NOT_FOUND",
            message: "Output not found",
            wantJSON: `{"error":{"code":"OUTPUT_NOT_FOUND","message":"Output not found"}}`,
        },
        {
            name:    "error with details",
            code:    "FILE_WRITE_FAILED",
            message: "Cannot write file",
            details: map[string]interface{}{"path": "src/test.go"},
            wantJSON: `{"error":{"code":"FILE_WRITE_FAILED","message":"Cannot write file","details":{"path":"src/test.go"}}}`,
        },
    }
}
```

#### **Integration Test Examples**

```go
// handlers/submit_task_test.go
func TestSubmitTaskHandler(t *testing.T) {
    // Setup mock provider
    mockProvider := &MockProvider{
        Response: "```go\nfunc main() {}\n```",
    }
    
    // Test submit task tool
    result := SubmitTaskWithProvider(mockProvider, models.SubmitTaskParams{
        Model:  "gemini-2.5-flash",
        Prompt: "Create a hello world",
        Files:  []string{"src/main.go"}, // Relative paths
    })
    
    // Verify structured response
    assert.NotEmpty(t, result.OutputID)
    assert.Contains(t, result.WorkingDirectory, "/")
    
    // Verify file was created in storage
    assert.FileExists(t, filepath.Join(storageDir, result.OutputID))
}

// handlers/get_metadata_test.go
func TestGetMetadataHandler(t *testing.T) {
    // Setup: Create test output with multiple blocks
    outputID := createTestOutputWithMultipleBlocks(t)
    
    // Test get metadata tool
    result := GetMetadataHandler(models.GetMetadataParams{
        OutputID: outputID,
    })
    
    // Verify structured metadata response
    assert.Equal(t, "COMPLETED", result.Metadata.Status)
    assert.Equal(t, 2, result.ContentAnalysis.BlocksFound)
    assert.Equal(t, "go", result.ContentAnalysis.Blocks[0].Language)
    assert.Equal(t, "md", result.ContentAnalysis.Blocks[1].Language)
    assert.Greater(t, result.Metadata.TokenEstimate, 0)
}

// handlers/write_file_test.go
func TestWriteFileHandler(t *testing.T) {
    // Setup: Create test output
    outputID := createTestOutput(t, "```go\nfunc test() {}\n```")
    tempDir := t.TempDir()
    
    // Test write output to file tool
    result := WriteFileHandler(models.WriteFileParams{
        OutputID: outputID,
        WriteTo:  "src/test.go", // Relative path
        Options: &models.ExtractionOptions{
            Extract: "code",
        },
    }, tempDir)
    
    // Verify structured success response
    assert.True(t, result.Success)
    assert.Equal(t, "src/test.go", result.Path)
    assert.Contains(t, result.AbsolutePath, "/src/test.go")
    assert.Greater(t, result.BytesWritten, 0)
    assert.Equal(t, tempDir, result.WorkingDirectory)
    
    // Verify file was actually written
    content, err := os.ReadFile(filepath.Join(tempDir, "src/test.go"))
    assert.NoError(t, err)
    assert.Contains(t, string(content), "func test()")
}
```

#### **E2E Test Examples**

```go
// e2e/golden_path_test.go
//go:build e2e

func TestNewAPIGoldenPath(t *testing.T) {
    // Start MCP server
    server := StartTestMCPServer(t)
    defer server.Stop()
    
    // Connect mock client
    client := mcp.NewClient(server.Address())
    
    // Test new 4-tool workflow
    
    // Step 1: Submit task
    submitResp := client.CallTool("delegate_submit_task", map[string]any{
        "model": "gemini-2.5-flash",
        "prompt": "Create a function to calculate fibonacci",
        "files": []string{"docs/requirements.md"}, // Relative paths
    })
    
    outputID := submitResp["output_id"].(string)
    assert.NotEmpty(t, outputID)
    
    // Step 2: Get metadata
    metadataResp := client.CallTool("delegate_get_output_metadata", map[string]any{
        "output_id": outputID,
    })
    
    metadata := metadataResp["metadata"].(map[string]any)
    assert.Equal(t, "COMPLETED", metadata["status"])
    
    contentAnalysis := metadataResp["content_analysis"].(map[string]any)
    blocksFound := int(contentAnalysis["blocks_found"].(float64))
    assert.Greater(t, blocksFound, 0)
    
    // Step 3a: Get content (high token cost)
    contentResp := client.CallTool("delegate_get_output_content", map[string]any{
        "output_id": outputID,
        "options": map[string]any{"extract": "code"},
    })
    
    content := contentResp["content"].(string)
    assert.Contains(t, content, "fibonacci")
    
    // Step 3b: Write to file (zero token cost)
    writeResp := client.CallTool("delegate_write_output_to_file", map[string]any{
        "output_id": outputID,
        "write_to": "src/fibonacci.go", // Relative path
        "options": map[string]any{"extract": "code"},
    })
    
    assert.True(t, writeResp["success"].(bool))
    assert.Equal(t, "src/fibonacci.go", writeResp["path"])
    assert.Greater(t, int(writeResp["bytes_written"].(float64)), 0)
}

func TestMultiBlockWorkflow(t *testing.T) {
    server := StartTestMCPServer(t)
    defer server.Stop()
    client := mcp.NewClient(server.Address())
    
    // Submit task that generates multiple blocks
    submitResp := client.CallTool("delegate_submit_task", map[string]any{
        "model": "gemini-2.5-flash",
        "prompt": "Create a React component with documentation",
    })
    
    outputID := submitResp["output_id"].(string)
    
    // Get metadata to understand structure
    metadataResp := client.CallTool("delegate_get_output_metadata", map[string]any{
        "output_id": outputID,
    })
    
    contentAnalysis := metadataResp["content_analysis"].(map[string]any)
    blocks := contentAnalysis["blocks"].([]interface{})
    
    // Write specific blocks to different files
    for i, block := range blocks {
        blockMap := block.(map[string]any)
        language := blockMap["language"].(string)
        
        var filename string
        switch language {
        case "jsx", "js":
            filename = "src/Component.jsx"
        case "md":
            filename = "docs/component.md"
        default:
            filename = fmt.Sprintf("output_%d.txt", i)
        }
        
        writeResp := client.CallTool("delegate_write_output_to_file", map[string]any{
            "output_id": outputID,
            "write_to": filename,
            "options": map[string]any{"block_index": i},
        })
        
        assert.True(t, writeResp["success"].(bool))
    }
}
```

## **6. Migration Testing**

During the migration period, we need to test both the old 3-tool API and the new 4-tool API to ensure compatibility and smooth transition.

### **Migration Test Strategy**

```go
// e2e/migration_test.go
//go:build e2e

func TestAPIMigrationCompatibility(t *testing.T) {
    server := StartTestMCPServer(t)
    defer server.Stop()
    client := mcp.NewClient(server.Address())
    
    t.Run("old_api_still_works", func(t *testing.T) {
        // Test the deprecated 3-tool workflow
        
        // Step 1: Old invoke
        invokeResp := client.CallTool("delegate_invoke", map[string]any{
            "model": "gemini-2.5-flash",
            "prompt": "Create a hello world function",
        })
        
        outputID := invokeResp["id"].(string)
        
        // Check for deprecation warning
        if warning, exists := invokeResp["deprecation_warning"]; exists {
            assert.Contains(t, warning.(string), "deprecated")
            assert.Contains(t, warning.(string), "delegate_submit_task")
        }
        
        // Step 2: Old check
        checkResp := client.CallTool("delegate_check", map[string]any{
            "output_id": outputID,
        })
        
        assert.True(t, checkResp["has_code"].(bool))
        
        // Check for deprecation warning
        if warning, exists := checkResp["deprecation_warning"]; exists {
            assert.Contains(t, warning.(string), "deprecated")
            assert.Contains(t, warning.(string), "delegate_get_output_metadata")
        }
        
        // Step 3: Old read (content mode)
        readResp := client.CallTool("delegate_read", map[string]any{
            "output_id": outputID,
            "options": map[string]any{"extract": "code"},
        })
        
        assert.Contains(t, readResp["content"].(string), "hello")
        
        // Check for deprecation warning
        if warning, exists := readResp["deprecation_warning"]; exists {
            assert.Contains(t, warning.(string), "deprecated")
            assert.Contains(t, warning.(string), "delegate_get_output_content")
        }
        
        // Step 4: Old read (write mode)
        writeResp := client.CallTool("delegate_read", map[string]any{
            "output_id": outputID,
            "options": map[string]any{
                "extract": "code",
                "write_to": "test_old_api.go",
            },
        })
        
        assert.Contains(t, writeResp["message"].(string), "written")
        
        // Check for deprecation warning
        if warning, exists := writeResp["deprecation_warning"]; exists {
            assert.Contains(t, warning.(string), "deprecated")
            assert.Contains(t, warning.(string), "delegate_write_output_to_file")
        }
    })
    
    t.Run("new_api_works", func(t *testing.T) {
        // Test the new 4-tool workflow (same as golden path test)
        // This ensures both APIs work simultaneously during migration
        
        submitResp := client.CallTool("delegate_submit_task", map[string]any{
            "model": "gemini-2.5-flash",
            "prompt": "Create a hello world function",
        })
        
        outputID := submitResp["output_id"].(string)
        
        // Ensure no deprecation warnings in new API
        assert.NotContains(t, submitResp, "deprecation_warning")
        
        metadataResp := client.CallTool("delegate_get_output_metadata", map[string]any{
            "output_id": outputID,
        })
        
        assert.NotContains(t, metadataResp, "deprecation_warning")
        assert.Equal(t, "COMPLETED", metadataResp["metadata"].(map[string]any)["status"])
        
        writeResp := client.CallTool("delegate_write_output_to_file", map[string]any{
            "output_id": outputID,
            "write_to": "test_new_api.go",
            "options": map[string]any{"extract": "code"},
        })
        
        assert.NotContains(t, writeResp, "deprecation_warning")
        assert.True(t, writeResp["success"].(bool))
    })
    
    t.Run("same_output_different_apis", func(t *testing.T) {
        // Create output with new API
        submitResp := client.CallTool("delegate_submit_task", map[string]any{
            "model": "gemini-2.5-flash",
            "prompt": "Create a test function",
        })
        
        outputID := submitResp["output_id"].(string)
        
        // Read with old API
        oldReadResp := client.CallTool("delegate_read", map[string]any{
            "output_id": outputID,
            "options": map[string]any{"extract": "all"},
        })
        
        // Read with new API
        newReadResp := client.CallTool("delegate_get_output_content", map[string]any{
            "output_id": outputID,
            "options": map[string]any{"extract": "all"},
        })
        
        // Content should be the same (ignoring structure differences)
        oldContent := oldReadResp["content"].(string)
        newContent := newReadResp["content"].(string)
        assert.Equal(t, oldContent, newContent)
    })