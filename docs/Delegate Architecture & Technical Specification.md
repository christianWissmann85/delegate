# **Delegate: Architecture & Technical Specification v1.0**

**Status:** Final | **Version:** 1.0 | **Date:** 2025-06-20

## **1. Overview & Philosophy**

**Project:** Delegate

**Mission:** Enable Claude Code to delegate code generation and analysis tasks to external LLMs through a dead-simple, bulletproof, and token-efficient MCP server.

**Core Philosophy (The "No Scope Creep" Mandate):**

* **Minimalism:** The server exposes exactly three MCP tools: invoke, check, and read. It does one thing and does it perfectly.
* **Reliability:** The system is designed for zero maintenance. With a Go-based MCP server, static typing, and no external runtime dependencies, it is built to be unbreakable.
* **Statelessness:** Each tool call is an atomic, independent transaction. There are no sessions, no conversation history, and no state management.
* **Robustness:** The server is resilient to network latency and provider timeouts through smart internal stream handling.

This system is the antithesis of the overly complex "AAG" project. It is an industrial-strength tool, not a development framework.

## **2. System Architecture**

Delegate runs as an MCP server that Claude Code connects to via the Model Context Protocol.

```
┌─────────────────┐      ┌────────────────────────┐      ┌───────────────────┐
│  Claude Code    │      │   Delegate MCP Server  │      │   LLM APIs        │
│  (MCP Client)   │<====>│   (Go Implementation)  │=====>│ (Claude, Gemini)  │
└─────────────────┘ MCP  └────────────────────────┘ HTTPS└───────────────────┘
                              |               ^
                              | (read/write)  |
                              v               |
            ┌──────────────────────────────────────────┐
            │       Local Filesystem (.delegate/)      │
            └──────────────────────────────────────────┘
```

**Execution Flow:**

1. Claude Code initiates an MCP connection to the Delegate server
2. Claude Code calls one of three tools: `invoke`, `check`, or `read`
3. Delegate processes the request:
   - For `invoke`: Calls the appropriate LLM API, streams the response, saves to disk
   - For `check`: Reads file metadata, analyzes content structure
   - For `read`: Extracts and returns requested content
4. Delegate returns the result via MCP protocol
5. Files persist in `.delegate/` directory for later retrieval

## **3. MCP Protocol & Tool Specification**

Communication is handled via the Model Context Protocol (MCP), exposing three tools to Claude Code.

### **Tool Definitions**

#### **delegate.invoke**

* **Description:** The primary tool to delegate a task to an external LLM. It generates an output file and returns metadata about it. It does **not** return the content itself.
* **Parameters:**
  ```typescript
  {
    model: string,        // "gemini-2.5-flash" | "gemini-2.5-pro" | "claude-sonnet-4-20250514" | "claude-opus-4-20250514"
    prompt: string,       // Natural language task description
    files?: string[],     // Optional: Paths to context files
    max_tokens?: number   // Optional: Maximum tokens to generate
  }
  ```

* **Returns:**
  ```typescript
  {
    id: string,          // "out_20250620_203000"
    path: string,        // "/path/to/repo/.delegate/outputs/out_20250620_203000.json"
    size_kb: number,     // 1.8
    model: string        // Echo of the model used
  }
  ```

* **Internal Implementation Detail: Server-Side Streaming**
  To prevent network timeouts on long-running LLM generations, the invoke handler MUST use the provider's streaming API endpoint. The implementation will:
  1. Initiate a streaming request to the LLM provider
  2. As data chunks are received, write them directly to a temporary file on disk
  3. Once the stream is complete, process the file and create the final output artifact
  4. This streaming is invisible to Claude Code - it only sees the final result

#### **delegate.check**

* **Description:** Inspects a previously generated output file without reading its content. Provides crucial metadata for Claude Code to decide *if* and *how* to read the file.
* **Parameters:**
  ```typescript
  {
    output_id: string    // "out_20250620_203000"
  }
  ```

* **Returns:**
  ```typescript
  {
    bytes: number,           // 1843
    size_kb: number,         // 1.8
    estimated_tokens: number,// 460
    has_code: boolean,       // true
    has_explanation: boolean,// true
    languages: string[]      // ["javascript"]
  }
  ```

#### **delegate.read**

* **Description:** Retrieves the actual content from an output file, with options to extract specific parts and limit token count.
* **Parameters:**
  ```typescript
  {
    output_id: string,
    options?: {
      extract?: "all" | "code" | "explanation",  // Default: "all"
      max_tokens?: number,                        // Truncate at this limit
      language?: string                           // Filter to specific language
    }
  }
  ```

* **Returns:**
  ```typescript
  {
    content: string,         // The extracted content
    truncated: boolean,      // true if max_tokens was hit
    tokens: number,          // Actual token count returned
    extraction: string,      // What was extracted
    language?: string        // If language filter was applied
  }
  ```

## **4. Configuration Management**

Configuration is handled via environment variables, with an optional config file for defaults.

1. **Environment Variables (Highest Priority):**
   * `ANTHROPIC_API_KEY` - Required for Claude models
   * `GOOGLE_API_KEY` - Required for Gemini models
   * `DELEGATE_LOG_LEVEL` - Options: `debug`, `info`, `warn`, `error`
   * `DELEGATE_TIMEOUT_SECONDS` - Default: 60
   * `DELEGATE_OUTPUT_DIR` - Default: `./.delegate`

2. **Config File (Optional):**
   * **Path:** `.delegate/config.json` within the project root
   * **Structure:**
     ```json
     {
       "default_model": "gemini-2.5-flash",
       "timeout_seconds": 60,
       "log_level": "info"
     }
     ```
   * Environment variables always override config file values

## **5. Logging & Storage**

All generated artifacts and logs are stored within a `.delegate` directory to avoid polluting the user's project.

**Directory Structure:**
```
.delegate/
├── outputs/
│   └── out_20250620_203000.json  // Stored output objects
├── logs/
│   └── delegate_20250620.log      // Daily log file
└── tmp/
    └── stream_xyz.tmp             // Temporary files during streaming
```

**Structured Logging:**
- Logs are written in JSON format for easy parsing
- Log levels: debug, info, warn, error
- Each log entry includes timestamp, level, message, and context

## **6. Provider Integration**

### **Supported Models**

| Model ID | Provider | Context Window | Use Case |
|----------|----------|----------------|----------|
| `gemini-2.5-flash` | Google | 1M tokens | Fast, everyday code generation |
| `gemini-2.5-pro` | Google | 1M tokens | Complex reasoning tasks |
| `claude-sonnet-4-20250514` | Anthropic | 200K tokens | Balanced quality and speed |
| `claude-opus-4-20250514` | Anthropic | 200K tokens | Highest quality output |

### **Provider Interface**

Each provider implements a common interface:

```go
type LLMProvider interface {
    GenerateStream(ctx context.Context, req GenerateRequest) (<-chan StreamChunk, error)
    GetCapabilities() ProviderCapabilities
}
```

### **Error Handling**

Delegate implements a simple, transparent error handling strategy that empowers Claude Code to make intelligent decisions:

**Core Principles:**
- Delegate performs basic retries (3 attempts with exponential backoff: 1s, 2s, 4s)
- After retries fail, return structured error information to Claude Code
- No automatic provider switching or complex orchestration
- Let Claude Code decide the best recovery action based on context

**Error Response Structure:**
```go
type DelegateError struct {
    Type            string   `json:"error"`           // rate_limited, provider_unavailable, timeout, etc.
    Provider        string   `json:"provider"`        // Which provider failed
    Code            int      `json:"error_code"`      // HTTP status code if applicable
    Message         string   `json:"message"`         // Human-readable description
    RetryAfter      int      `json:"retry_after"`     // Seconds to wait (if provided by API)
    Alternatives    []string `json:"alternative_models"` // Other models that could be tried
}
```

**Error Types:**
- `rate_limited` (429): Provider rate limit exceeded
- `provider_unavailable` (503): Service temporarily unavailable
- `timeout` (504): Request exceeded timeout limit
- `provider_error` (500): Internal provider error
- `network_error`: Connection or network failure

**Recovery Strategy:**
This approach keeps Delegate simple while giving Claude Code the flexibility to:
- Retry with a different model
- Wait and retry the same model
- Fall back to internal code generation
- Inform the user and ask for preferences

## **7. Code Project Structure**

The Go implementation maintains clean separation of concerns:

```
delegate/
├── main.go               // MCP server entry point
├── go.mod
└── internal/
    ├── mcp/              // MCP protocol implementation
    │   └── server.go     // Handles MCP connections and tool routing
    ├── config/           // Configuration management
    ├── handlers/         // Implements invoke, check, read logic
    ├── providers/        // LLM provider implementations
    │   ├── interface.go  // Common provider interface
    │   ├── anthropic.go  // Claude integration
    │   └── google.go     // Gemini integration
    ├── extractor/        // Code/explanation extraction logic
    └── storage/          // File system operations
```

## **8. Security Considerations**

1. **API Key Protection:**
   - Never logged or exposed in responses
   - Only transmitted to respective LLM providers over HTTPS
   - Stored in environment variables, not config files

2. **File System Security:**
   - Path traversal prevention on all file operations
   - Output files use sanitized IDs
   - Automatic cleanup of files older than 24 hours

3. **Input Validation:**
   - Model IDs validated against allowlist
   - File paths checked for directory traversal
   - Prompt length limited to provider maximums

## **9. Performance Characteristics**

| Operation | Typical Latency | Notes |
|-----------|----------------|-------|
| `invoke` | 2-30 seconds | Depends on prompt complexity and model |
| `check` | <100ms | File system metadata only |
| `read` | <500ms | Depends on content size and extraction |

**Resource Usage:**
- Memory: <50MB baseline, streaming prevents large allocations
- Disk: Configurable output directory, automatic cleanup
- CPU: Minimal, mostly I/O bound

## **10. Future Considerations (Post v1.0)**

Per the "No Scope Creep" mandate, these are **not** in v1.0:
- Batch operations (invoke multiple prompts)
- Caching of responses
- Additional providers
- Progress indicators
- Analytics or metrics

The focus remains on doing three things perfectly: invoke, check, and read.