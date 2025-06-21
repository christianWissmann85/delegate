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
* **Tooling:** Go's standard testing package
* **Execution:** `go test ./internal/...`

### **Level 2: Integration Tests (The Core)**

* **Purpose:** To verify that our application modules correctly interact with external systems (filesystem and mocked LLM APIs)
* **Scope:**
  * **Storage Integration:**
    * Test that the storage package correctly creates the `.delegate` directory structure
    * Verify that `invoke` correctly writes output files and `read` correctly reads them
    * Test the 24-hour cleanup logic
  * **LLM Provider Integration (with Mocks):**
    * Create a **mock LLM provider** that implements the Provider interface
    * Return predictable, pre-defined streaming responses
    * Test the invoke handler's stream handling and error normalization
    * Verify timeout handling and retry logic
  * **MCP Tool Integration:**
    * Test each tool with various parameter combinations
    * Verify error responses for invalid inputs
    * Test the full invoke → check → read workflow
* **Tooling:** Go's standard testing package with test tables and mock structs
* **Execution:** Part of standard `go test` run - fast and reliable

### **Level 3: End-to-End (E2E) Tests (The Golden Path)**

* **Purpose:** To provide ultimate confirmation that the MCP server works with Claude Code
* **Scope:**
  * Test the MCP server with a mock MCP client that simulates Claude Code
  * Use **real API keys** for actual LLM calls (in CI only)
  * Primary E2E test flow:
    1. Connect mock client to MCP server
    2. Call `invoke` tool with a simple prompt
    3. Verify response has valid ID and metadata
    4. Call `check` tool with the ID
    5. Verify metadata matches expectations
    6. Call `read` tool to get content
    7. Verify content quality and structure
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
│   └── mock_stream_chunks.json            # Mock streaming response chunks
│
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   ├── config_test.go                 # Unit: Environment variable loading
│   │   └── validate_test.go                # Unit: Config validation
│   ├── extractor/
│   │   ├── extractor.go
│   │   ├── extractor_test.go              # Unit: Code/explanation extraction
│   │   └── patterns_test.go                # Unit: Regex patterns
│   ├── handlers/
│   │   ├── invoke.go
│   │   ├── invoke_test.go                 # Integration: Invoke handler
│   │   ├── check_test.go                  # Integration: Check handler
│   │   └── read_test.go                   # Integration: Read handler
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
    └── golden_path_test.go                # E2E: Full MCP workflow (//go:build e2e)
```

**Key Points:**
- Test files use `_test.go` suffix and live next to the code they test
- `testdata/` directory contains test fixtures (Go convention)
- Mock provider is properly placed in `internal/providers/mock/`
- E2E tests are in a separate `e2e/` package with build tags

## **5. Key Test Scenarios**

### **Unit Test Examples**

```go
// extractor_test.go
func TestExtractCode(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        wantCode []CodeBlock
        wantErr  bool
    }{
        {
            name:  "extracts JavaScript code",
            input: "Here's the code:\n```javascript\nconst x = 1;\n```",
            wantCode: []CodeBlock{{Language: "javascript", Content: "const x = 1;"}},
        },
        // More test cases...
    }
}
```

### **Integration Test Examples**

```go
// handlers_test.go
func TestInvokeHandler(t *testing.T) {
    // Setup mock provider
    mockProvider := &MockProvider{
        Response: "```go\nfunc main() {}\n```",
    }
    
    // Test invoke tool
    result := InvokeWithProvider(mockProvider, InvokeParams{
        Model:  "gemini-2.5-flash",
        Prompt: "Create a hello world",
    })
    
    // Verify file was created
    assert.FileExists(t, result.Path)
    assert.Equal(t, "gemini-2.5-flash", result.Model)
}
```

### **E2E Test Example**

```go
// e2e_test.go
//go:build e2e

func TestGoldenPath(t *testing.T) {
    // Start MCP server
    server := StartTestMCPServer(t)
    defer server.Stop()
    
    // Connect mock client
    client := mcp.NewClient(server.Address())
    
    // Test full workflow
    invokeResp := client.CallTool("delegate.invoke", map[string]any{
        "model": "gemini-2.5-flash",
        "prompt": "Create a function to calculate fibonacci",
    })
    
    checkResp := client.CallTool("delegate.check", map[string]any{
        "output_id": invokeResp["id"],
    })
    
    assert.True(t, checkResp["has_code"].(bool))
    
    readResp := client.CallTool("delegate.read", map[string]any{
        "output_id": invokeResp["id"],
        "options": map[string]any{"extract": "code"},
    })
    
    assert.Contains(t, readResp["content"].(string), "fibonacci")
}
```

## **6. Test Data Management**

* **Mock Responses:** Stored in `test/fixtures/` as JSON files
* **Test Outputs:** Use temporary directories that are cleaned up after tests
* **API Keys for E2E:** Only used in CI, never committed to repository

## **7. Performance Testing**

While not part of regular test runs, we include basic performance benchmarks:

```go
func BenchmarkInvoke(b *testing.B) {
    // Benchmark with mock provider to test overhead
}

func BenchmarkExtractor(b *testing.B) {
    // Benchmark code extraction on large responses
}
```

## **8. Testing Guidelines**

1. **Keep tests simple and focused** - One test, one assertion
2. **Use table-driven tests** - Easy to add new cases
3. **Mock external dependencies** - Tests should be fast and deterministic
4. **Test edge cases** - Empty inputs, large files, network errors
5. **No flaky tests** - If it fails intermittently, fix it or remove it

---

Remember: We test to gain confidence, not to achieve metrics. Quality over quantity!