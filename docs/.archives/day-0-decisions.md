# Day 0 Decisions - Lock These Down Before Coding

## 🎯 Purpose
Make all tactical decisions upfront to avoid refactoring later. Each decision here is final.

## 1. MCP Implementation Strategy

### Decision: Build Minimal MCP Handler from Scratch
```go
// No external dependencies, just standard library
// Simple JSON-RPC over stdio
```

**Rationale:**
- No mature Go MCP libraries exist
- We only need to support 3 tools
- Keeps binary size small
- Full control over implementation

### Implementation Details:
- Read JSON-RPC from stdin
- Write JSON-RPC to stdout
- Support only required MCP methods:
  - `initialize`
  - `tools/list`
  - `tools/call`
- Log to stderr only

## 2. File Structure (No Refactoring Needed)

```
delegate/
├── main.go                         # ~50 LOC - Just wire up and start server
├── go.mod
├── go.sum
└── internal/
    ├── mcp/
    │   ├── server.go               # ~200 LOC - MCP server coordination
    │   ├── protocol.go             # ~200 LOC - JSON-RPC message handling
    │   ├── tools.go                # ~150 LOC - Tool registration/dispatch
    │   └── types.go                # ~100 LOC - MCP type definitions
    │
    ├── handlers/
    │   ├── invoke.go               # ~200 LOC - Invoke orchestration
    │   ├── check.go                # ~80 LOC - Simple metadata lookup
    │   ├── read.go                 # ~150 LOC - Read with options
    │   └── types.go                # ~50 LOC - Handler types
    │
    ├── providers/
    │   ├── provider.go             # ~50 LOC - Provider interface
    │   ├── anthropic/
    │   │   ├── client.go           # ~200 LOC - Anthropic API client
    │   │   └── stream.go           # ~150 LOC - Streaming logic
    │   ├── google/
    │   │   ├── client.go           # ~200 LOC - Gemini API client
    │   │   └── stream.go           # ~150 LOC - Streaming logic
    │   ├── factory.go              # ~100 LOC - Provider selection
    │   └── errors.go               # ~100 LOC - Error normalization
    │
    ├── extractor/
    │   ├── extractor.go            # ~200 LOC - Main extraction logic
    │   ├── patterns.go             # ~100 LOC - Regex patterns
    │   └── types.go                # ~50 LOC - Extraction types
    │
    ├── storage/
    │   ├── store.go                # ~200 LOC - Storage interface & impl
    │   ├── cleanup.go              # ~100 LOC - 24-hour cleanup
    │   └── types.go                # ~50 LOC - Storage types
    │
    ├── config/
    │   ├── config.go               # ~150 LOC - Configuration loading
    │   └── validate.go             # ~80 LOC - Validation logic
    │
    └── models/
        ├── output.go               # ~80 LOC - Output data structure
        ├── request.go              # ~60 LOC - Request types
        └── errors.go               # ~100 LOC - Error types
```

**Key Decisions:**
- Split providers into subdirectories (anthropic/, google/)
- Separate types.go files to keep main logic clean
- No file exceeds 200 LOC
- Clear separation of concerns

## 3. Core Technical Decisions

### 3.1 File Path Handling
```go
// Decision: All paths relative to CWD where delegate is running
type FileResolver struct {
    workDir string  // Set at startup
}

func (f *FileResolver) Resolve(path string) (string, error) {
    // 1. Clean the path
    cleaned := filepath.Clean(path)
    
    // 2. Make absolute
    abs := filepath.Join(f.workDir, cleaned)
    
    // 3. Verify within workDir (security)
    if !strings.HasPrefix(abs, f.workDir) {
        return "", ErrPathTraversal
    }
    
    return abs, nil
}
```

**Glob Support:** NO - Keep it simple. Explicit file lists only.

### 3.2 Code Extraction Patterns
```go
// Primary patterns (in order of precedence)
var extractPatterns = []Pattern{
    {
        Name:  "FencedCodeBlock",
        Regex: regexp.MustCompile("```(?P<lang>\\w+)?\\n(?P<code>.*?)```"),
    },
    {
        Name:  "AltFencedBlock", 
        Regex: regexp.MustCompile("~~~(?P<lang>\\w+)?\\n(?P<code>.*?)~~~"),
    },
    {
        Name:  "IndentedBlock",
        Regex: regexp.MustCompile("(?m)^(    .+\\n)+"),
    },
}

// Decision: Extract ALL code blocks, preserve order
// Let Claude Code decide what to use
```

### 3.3 Storage Format
```json
{
    "id": "out_20250620_143022",
    "created_at": "2025-06-20T14:30:22Z",
    "model": "gemini-2.5-flash",
    "prompt": "Original prompt text...",
    "files": ["file1.js", "file2.js"],
    "response": {
        "raw": "Full LLM response with markdown...",
        "extracted": {
            "code": [
                {
                    "language": "javascript",
                    "content": "console.log('hello');",
                    "line_start": 5,
                    "line_end": 7
                }
            ],
            "explanation": "Text without code blocks..."
        }
    },
    "metadata": {
        "total_bytes": 4523,
        "estimated_tokens": 1130,
        "provider_request_id": "abc-123"
    }
}
```

### 3.4 Streaming Strategy
```go
// Decision: Stream to temp file, move when complete
func (h *InvokeHandler) streamToFile(stream <-chan Chunk) (string, error) {
    // 1. Create temp file
    tmp := filepath.Join(h.store.TempDir(), "stream_*")
    
    // 2. Stream chunks to temp
    for chunk := range stream {
        // Write chunk
        // Update progress
    }
    
    // 3. Move to final location atomically
    final := h.store.GetOutputPath(outputID)
    return os.Rename(tmp, final)
}
```

### 3.5 Provider API Details

**Anthropic:**
- Endpoint: `https://api.anthropic.com/v1/messages`
- API Version: `2023-06-01`
- Model IDs: Exactly as specified in docs

**Google (Gemini):**
- Use Google AI Studio API (simpler than Vertex)
- Endpoint: `https://generativelanguage.googleapis.com/v1/models/{model}:generateContent`
- Models: `gemini-2.5-flash`, `gemini-2.5-pro`

### 3.6 Error Handling Matrix

| Provider Error | HTTP Code | Our Error Type | Retry? |
|---------------|-----------|----------------|--------|
| Rate limit | 429 | `rate_limited` | Yes, with backoff |
| Server error | 500 | `provider_error` | Yes, 3 times |
| Timeout | - | `timeout` | Yes, once |
| Bad request | 400 | `invalid_request` | No |
| Auth error | 401/403 | `auth_error` | No |

### 3.7 MCP Tool Registration
```go
// Exact tool definitions
var tools = []MCPTool{
    {
        Name: "delegate_invoke",
        Description: "Delegate heavy tasks (code generation, document analysis, large file processing) to other LLMs to save Claude Code's context tokens. Use this when: generating large amounts of code, analyzing multiple documents, processing entire codebases, or any task that would consume significant context. Supports Gemini models (1M token context) and Claude models. Returns an output_id for async retrieval.",
        InputSchema: JSONSchema{
            Type: "object",
            Properties: map[string]Property{
                "model": {
                    Type: "string",
                    Enum: []string{"gemini-2.5-flash", "gemini-2.5-pro", "claude-sonnet-4-20250514", "claude-opus-4-20250514"},
                    Description: "The LLM model to use",
                },
                // ... etc
            },
        },
    },
    // ... check and read tools
}
```

## 4. Development Constraints

### 4.1 No External Dependencies Beyond:
- Official Google Gemini SDK
- Official Anthropic SDK
- Standard library only for everything else

### 4.2 Coding Standards
- Every exported function has a doc comment
- Every error includes context
- No naked returns
- No init() functions
- Contexts passed explicitly

### 4.3 Testing Requirements
- Mock providers for all tests
- No real API calls in unit tests
- Integration tests behind build tag
- Table-driven tests preferred

## 5. Day 1 Starting Points

### Morning: MCP Protocol
1. Create `internal/mcp/types.go` with all type definitions
2. Implement `internal/mcp/protocol.go` for JSON-RPC
3. Get basic message exchange working

### Afternoon: Storage Layer
1. Create `internal/storage/store.go` with interface
2. Implement file-based storage
3. Add cleanup goroutine

This plan eliminates refactoring risk. Every file has a clear purpose and size limit. Interfaces are defined upfront. Let's build it right the first time! 🚀