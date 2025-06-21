# Implementation Questions to Resolve

## Critical Decisions Needed Before Starting

### 1. MCP Implementation
- [ ] Which Go MCP library to use? (or implement from scratch?)
- [ ] How to handle stdio vs TCP transport?
- [ ] Streaming response strategy over MCP?

### 2. File Handling
- [ ] How to resolve file paths in `files` parameter?
  - Relative to what directory?
  - Support for glob patterns?
  - How to handle missing files?
- [ ] Security boundaries for file access?

### 3. Code Extraction
- [ ] Regex patterns for code block detection?
- [ ] Handle alternative fence styles (~~~, ```, etc)?
- [ ] Extract language hints from code blocks?
- [ ] What if no code blocks found?
- [ ] How to handle mixed content (code + explanation)?

### 4. Provider Details
- [ ] Exact Anthropic API endpoint and version?
- [ ] Gemini API: Vertex AI or Google AI Studio?
- [ ] How to handle provider-specific errors?
- [ ] Streaming implementation differences?

### 5. Storage Details
- [ ] File permissions for .delegate directory?
- [ ] JSON structure for stored outputs?
- [ ] Include raw response or just extracted content?
- [ ] How to handle concurrent writes?

### 6. Error Scenarios
- [ ] What if extraction completely fails?
- [ ] Partial streaming failures?
- [ ] Cleanup failures?
- [ ] Invalid model names?

### 7. Testing Strategy
- [ ] Mock provider responses format?
- [ ] Integration test with real MCP client?
- [ ] How to test streaming behavior?

## Proposed Solutions

### MCP Library
Recommend: Start with basic stdio implementation, no external libraries needed:
```go
// Read from stdin, write to stdout
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() {
    message := scanner.Text()
    // Parse JSON-RPC message
}
```

### File Resolution
```go
// Always resolve relative to current working directory
absPath := filepath.Abs(filepath.Join(cwd, requestedPath))
// Verify it's under allowed directory
if !strings.HasPrefix(absPath, cwd) {
    return error("path traversal attempted")
}
```

### Code Extraction
```go
// Primary pattern
codeBlockRegex := regexp.MustCompile("```(?:(\w+)\n)?(.*?)```")
// Fallback pattern  
altBlockRegex := regexp.MustCompile("~~~(?:(\w+)\n)?(.*?)~~~")
```

These decisions need to be made before Day 1 of implementation!