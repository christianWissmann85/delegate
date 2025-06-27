# **Delegate: Architecture & Technical Specification**

**Status:** Final | **Version:** 2.0 | **Date:** 2025-06-27  
**Implementation Status:** Production/Deployed

**Reviewed by Christian Wissmann for Delegate V2.0**

## **1. Overview & Philosophy**

**Project:** Delegate

**Mission:** Enable Claude Code to delegate code generation and analysis tasks to external LLMs through a dead-simple, bulletproof, and token-efficient MCP server.

**Core Philosophy (The "No Scope Creep" Mandate):**

* **Minimalism:** The server exposes exactly four MCP tools: submit_task, get_output_metadata, get_output_content, and write_output_to_file. It does one thing and does it perfectly.
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
            ┌──────────────────────────────────────────────┐
            │       Local Filesystem (.delegate/)      │
            └──────────────────────────────────────────────┘
```

**Execution Flow:**

1. Claude Code initiates an MCP connection to the Delegate server
2. Claude Code calls one of four tools: `delegate_submit_task`, `delegate_get_output_metadata`, `delegate_get_output_content`, or `delegate_write_output_to_file`
3. Delegate processes the request:
   - For `submit_task`: Calls the appropriate LLM API, streams the response, saves to disk, returns only an output_id
   - For `get_output_metadata`: Reads file metadata and analyzes content structure without loading content
   - For `get_output_content`: Extracts and returns requested content into Claude's context (token cost)
   - For `write_output_to_file`: Writes content directly to disk without token consumption
4. Delegate returns structured JSON responses via MCP protocol
5. Files persist in `.delegate/` directory for later retrieval

## **3. MCP Protocol & Tool Specification**

Communication is handled via the Model Context Protocol (MCP), exposing four tools to Claude Code. Each tool has a single, clear purpose with predictable behavior.

### **Tool Definitions**

#### **delegate_submit_task**

* **Description:** STEP 1: Submits a generation task to an external LLM (~50-100 tokens). This is an asynchronous operation that creates a temporary output artifact and returns a unique output_id. The content is NOT returned directly. Use other delegate_* tools to access the output.
* **Parameters:**
  ```typescript
  {
    model: string,        // "gemini-2.5-flash" | "gemini-2.5-pro" | "claude-sonnet-4-20250514" | "claude-opus-4-20250514"
    prompt: string,       // Natural language task description
    files?: string[],     // Optional: Relative paths to context files (e.g., "src/model.go", "docs/api.md")
    max_tokens?: number,  // Optional: Maximum tokens to generate
    timeout?: number      // Optional: Timeout in seconds
  }
  ```

* **Returns:**
  ```typescript
  {
    output_id: string,         // "out_20250620_203000"
    working_directory: string  // "/home/user/project"
  }
  ```

* **Internal Implementation Detail: Server-Side Streaming**
  To prevent network timeouts on long-running LLM generations, the submit_task handler MUST use the provider's streaming API endpoint. The implementation will:
  1. Initiate a streaming request to the LLM provider
  2. As data chunks are received, write them directly to a temporary file on disk
  3. Once the stream is complete, process the file and create the final output artifact
  4. This streaming is invisible to Claude Code - it only sees the final result

#### **delegate_get_output_metadata**

* **Description:** STEP 2 (Optional): Retrieves structured metadata about an output artifact (~20 tokens). Use this to decide whether to retrieve content into context or write directly to a file. This tool does NOT return the content itself.
* **Parameters:**
  ```typescript
  {
    output_id: string    // "out_20250620_203000"
  }
  ```

* **Returns:**
  ```typescript
  {
    metadata: {
      output_id: string,
      status: "COMPLETED" | "IN_PROGRESS" | "FAILED",
      size_kb: number,
      line_count: number,
      token_estimate: number,
      is_truncated: boolean,
      truncation_reason: string | null
    },
    content_analysis: {
      blocks_found: number,
      blocks: Array<{
        index: number,
        language: string,
        size_kb: number,
        lines: number,
        preview: string      // First line of the block
      }>
    }
  }
  ```

#### **delegate_get_output_content**

* **Description:** Retrieves the full or partial content of an output artifact into the agent's context (~30+ tokens plus content). This operation consumes tokens proportional to the content size. Use options to extract specific parts (e.g., extract: 'code').
* **Parameters:**
  ```typescript
  {
    output_id: string,
    options?: {
      extract?: "all" | "code" | "explanation",  // Default: "all"
      max_tokens?: number,                        // Truncate at this limit
      block_index?: number,                       // For multi-block outputs, select specific block
      language?: string                           // Filter code blocks by language
    }
  }
  ```

* **Returns:**
  ```typescript
  {
    content: string,         // The extracted content
    metadata: {
      output_id: string,
      tokens_returned: number,
      is_truncated: boolean,
      truncation_reason: string | null
    }
  }
  ```

#### **delegate_write_output_to_file**

* **Description:** Writes the content of an output artifact directly to a specified file path (relative to working directory). This operation consumes ZERO content tokens. Use options to select specific parts to write (e.g., extract: 'code', block_index: 0).
* **Parameters:**
  ```typescript
  {
    output_id: string,
    write_to: string,        // Relative file path (e.g., "src/component.jsx", "tmp/output.go")
    options?: {
      extract?: "all" | "code" | "explanation",  // Default: "all"
      block_index?: number,                       // For multi-block outputs, select specific block
      language?: string                           // Filter code blocks by language
    }
  }
  ```

* **Returns:**
  ```typescript
  {
    success: boolean,
    path: string,            // Relative path of file written
    absolute_path: string,   // Absolute path of file written
    bytes_written: number,
    message: string,         // Human-readable success message
    working_directory: string
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

Delegate implements a structured, transparent error handling strategy that empowers Claude Code to make intelligent decisions. All errors are returned as structured JSON for reliable programmatic handling.

**Core Principles:**
- Delegate performs basic retries (3 attempts with exponential backoff: 1s, 2s, 4s)
- After retries fail, return structured error information to Claude Code
- No automatic provider switching or complex orchestration
- Let Claude Code decide the best recovery action based on context

**Structured Error Response Format:**
```json
{
  "error": {
    "code": "PROVIDER_ERROR",
    "message": "Rate limit exceeded for gemini-2.5-flash",
    "details": {
      "provider": "google",
      "model": "gemini-2.5-flash",
      "http_status": 429,
      "retry_after": 60,
      "alternative_models": ["gemini-2.5-pro", "claude-sonnet-4-20250514"]
    }
  }
}
```

**Error Codes:**
- `INVALID_REQUEST`: Bad input parameters or validation failure
- `OUTPUT_NOT_FOUND`: Requested output_id doesn't exist or has expired
- `PROVIDER_ERROR`: Any error from the LLM provider (rate limits, timeouts, etc.)
- `FILE_WRITE_FAILED`: Unable to write to the specified file path
- `PATH_TRAVERSAL_ATTEMPT`: Security violation in file path

**Recovery Strategy:**
This approach keeps Delegate simple while giving Claude Code the flexibility to:
- Retry with a different model
- Wait and retry the same model
- Fall back to internal code generation
- Inform the user and ask for preferences

## **7. Code Project Structure**

The Go implementation maintains clean separation of concerns, reflecting the new 4-tool architecture:

```
delegate/
├── main.go                     // Entry point with graceful shutdown
├── go.mod                      // Go 1.18+ module definition
├── .gitignore                  // Ignore patterns for Go project
└── internal/
    ├── mcp/                    // MCP protocol implementation
    │   ├── server.go           // Server lifecycle and tool registration
    │   ├── protocol.go         // JSON-RPC message handling
    │   ├── tools.go            // Tool definitions and schemas
    │   └── types.go            // MCP type definitions
    │
    ├── config/                 // Configuration management
    │   ├── config.go           // Environment variable loading
    │   ├── validate.go         // Configuration validation
    │   └── languages.json      // Language detection configuration
    │
    ├── handlers/               // Business logic for each tool
    │   ├── submit_task.go      // Submit task handler implementation
    │   ├── get_metadata.go     // Get metadata handler implementation
    │   ├── get_content.go      // Get content handler implementation
    │   ├── write_file.go       // Write file handler implementation
    │   └── types.go            // Shared handler types
    │
    ├── providers/              // LLM provider integrations
    │   ├── provider.go         // Provider interface definition
    │   ├── factory.go          // Provider selection logic
    │   ├── errors.go           // Error normalization
    │   ├── anthropic/          // Anthropic-specific implementation
    │   │   ├── client.go       // Claude API client
    │   │   └── stream.go       // Streaming response handler
    │   └── google/             // Google-specific implementation
    │       ├── client.go       // Gemini API client
    │       └── stream.go       // Streaming response handler
    │
    ├── extractor/              // Content extraction logic
    │   ├── extractor.go        // Main extraction logic
    │   ├── patterns.go         // Regex patterns for code blocks
    │   └── types.go            // Extraction types
    │
    ├── storage/                // File system operations
    │   ├── store.go            // Storage interface and implementation
    │   ├── cleanup.go          // 24-hour cleanup routine
    │   └── types.go            // Storage types
    │
    ├── models/                 // Shared data structures
    │   ├── output.go           // Output file structure
    │   ├── request.go          // Request types
    │   ├── responses.go        // Structured response types
    │   └── errors.go           // Error types and AsDelegateError helper
    │
    └── logger/                 // Structured logging
        └── logger.go           // JSON logging to stderr
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
   - Explicit security comments in write_file handler

3. **Input Validation:**
   - Model IDs validated against allowlist
   - File paths checked for directory traversal
   - Prompt length limited to provider maximums
   - Relative paths resolved safely within working directory

## **9. Performance Characteristics**

| Operation | Typical Latency | Notes |
|-----------|----------------|-------|
| `submit_task` | 2-30 seconds | Depends on prompt complexity and model |
| `get_metadata` | <100ms | File system metadata only |
| `get_content` | <500ms | Depends on content size and extraction |
| `write_file` | <200ms | Direct file write, no content parsing |

**Resource Usage:**
- Memory: <50MB baseline, streaming prevents large allocations
- Disk: Configurable output directory, automatic cleanup
- CPU: Minimal, mostly I/O bound

The focus remains on doing four things perfectly: submit_task, get_output_metadata, get_output_content, and write_output_to_file.