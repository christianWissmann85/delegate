# **Delegate API Reference v1.0**

**Status:** Final | **Version:** 1.0 | **Date:** 2025-06-20

## **Introduction**

Welcome to the Delegate API Reference.

Delegate is a high-performance, minimalist MCP server designed to act as a token-efficient intermediary for Large Language Model (LLM) tasks. It allows Claude Code to delegate code generation and analysis to external models (Google's Gemini and Anthropic's Claude) without consuming its own context window.

As an MCP server, Delegate exposes three simple tools that Claude Code can call to manage delegated tasks efficiently.

The core workflow is a simple, three-step process:

1. **invoke**: Delegate a task to an LLM. This creates an output artifact but does not return its content.
2. **check**: Inspect the metadata of the output artifact to understand its size and structure.
3. **read**: Retrieve the content of the artifact, with options to extract only the necessary parts.

## **Authentication**

The Delegate MCP server requires API keys for the underlying LLM providers to be available as environment variables:

* `GOOGLE_API_KEY`: For Gemini models
* `ANTHROPIC_API_KEY`: For Claude models

These keys are never transmitted except to their respective providers over HTTPS.

## **MCP Tools**

### **delegate_invoke**

STEP 1: Delegates a generation task to a specified LLM. This is an asynchronous operation that creates a persistent output artifact and returns a unique ID for it. Does NOT write files directly - stores in temp storage. Each input file must be <1MB, but total can exceed.

#### **Tool Definition**

```typescript
{
  name: "delegate_invoke",
  description: "STEP 1: Delegate file generation to save tokens. Does NOT write files directly - stores in temp storage. Returns output_id for use with delegate_check then delegate_read(write_to). IMPORTANT: Use ABSOLUTE paths in 'files' parameter. Each file must be <1MB, but total can exceed.",
  parameters: {
    model: {
      type: "string",
      required: true,
      enum: ["gemini-2.5-flash", "gemini-2.5-pro", "claude-sonnet-4-20250514", "claude-opus-4-20250514"],
      description: "The LLM model to use for generation"
    },
    prompt: {
      type: "string", 
      required: true,
      description: "Natural language description of the task."
    },
    files: {
      type: "array",
      items: { type: "string" },
      required: false,
      description: "ABSOLUTE file paths to include as context (not relative!). Example: '/home/user/project/model.go' not 'model.go'. Each file must be <1MB."
    },
    max_tokens: {
      type: "number",
      required: false,
      description: "Maximum tokens to generate (defaults to model maximum)"
    },
    code_only: {
      type: "boolean",
      required: false,
      description: "Return only code without explanations (default: false)"
    },
    language_hint: {
      type: "string",
      required: false,
      description: "Expected programming language(s) for better extraction"
    },
    timeout: {
      type: "number",
      required: false,
      description: "Timeout in seconds (default: 180, max: 600). Suggested: 180s minimum for code, 400s minimum for creative tasks, 400-600s for very large file(s)/bundle(s) analysis."
    }
  }
}
```

#### **Example Usage**

```javascript
const result = await mcp.invoke({
  model: "gemini-2.5-flash",
  prompt: "Create a robust error handling middleware for Express.js",
  files: ["./src/server.js", "./src/routes.js"],
  max_tokens: 4000
});
```

#### **Success Response**

```javascript
{
  id: "out_20250620_204000",
  path: "/Users/you/project/.delegate/outputs/out_20250620_204000.json",
  size_kb: 2.1,
  model: "gemini-2.5-flash"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier for the output, used in `check` and `read` |
| `path` | string | Absolute path to the stored output artifact |
| `size_kb` | number | Size of the generated response in kilobytes |
| `model` | string | The model that was used, echoed back |

### **delegate_check**

STEP 2: Retrieves metadata about a previously generated output artifact without reading its content. Essential for token-efficient operations. Use this after delegate_invoke, before delegate_read to understand size and avoid consuming unnecessary tokens.

#### **Tool Definition**

```typescript
{
  name: "delegate_check",
  description: "STEP 2: Check delegated task status and size before retrieving. Shows token count and file size. Use this after delegate_invoke, before delegate_read. Helps avoid consuming unnecessary tokens.",
  parameters: {
    output_id: {
      type: "string",
      required: true,
      description: "The ID returned from invoke"
    }
  }
}
```

#### **Example Usage**

```javascript
const metadata = await mcp.check({
  output_id: "out_20250620_204000"
});
```

#### **Success Response**

```javascript
{
  bytes: 2150,
  size_kb: 2.1,
  estimated_tokens: 537,
  has_code: true,
  has_explanation: true,
  languages: ["javascript", "json"]
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `bytes` | number | Exact size in bytes |
| `size_kb` | number | Size in kilobytes |
| `estimated_tokens` | number | Rough token estimate (bytes / 4) |
| `has_code` | boolean | Whether code blocks were detected |
| `has_explanation` | boolean | Whether explanatory text was detected |
| `languages` | array | Programming languages found in code blocks |
| `is_truncated` | boolean | Whether the output appears to be truncated |
| `truncation_reason` | string | Reason for truncation detection (if applicable) |
| `truncation_confidence` | number | Confidence score (0.0-1.0) for truncation detection |

### **delegate_read**

STEP 3: Reads the content of an output artifact, with powerful options for extraction and truncation. The key feature is the 'write_to' option which saves files directly to disk WITHOUT consuming any tokens.

#### **Tool Definition**

```typescript
{
  name: "delegate_read",
  description: "STEP 3: Get delegated results. WORKFLOW: invoke -> check -> read. To save tokens: use 'write_to' with ABSOLUTE path to write file directly (no content returned). To get content: omit 'write_to'. Use 'extract' for code-only or explanation-only.",
  parameters: {
    output_id: {
      type: "string",
      required: true,
      description: "The ID returned from invoke"
    },
    options: {
      type: "object",
      required: false,
      properties: {
        extract: {
          type: "string",
          enum: ["all", "code", "explanation"],
          default: "all",
          description: "What to extract from the output"
        },
        max_tokens: {
          type: "number",
          description: "Truncate to this many tokens"
        },
        language: {
          type: "string",
          description: "When extracting code, filter to this language"
        },
        write_to: {
          type: "string",
          description: "ABSOLUTE file path to write content directly to disk (saves tokens - no content returned). Example: '/home/user/project/new_file.go' not 'new_file.go'"
        },
        block_index: {
          type: "number",
          description: "When multiple code blocks exist, specify which one to extract (0-based index). Use after getting block list warning."
        }
      }
    }
  }
}
```

#### **Example Usage**

```javascript
// Read everything
const full = await mcp.read({
  output_id: "out_20250620_204000"
});

// Read only code, max 2000 tokens
const code = await mcp.read({
  output_id: "out_20250620_204000",
  options: {
    extract: "code",
    max_tokens: 2000
  }
});

// Read only JavaScript code
const jsOnly = await mcp.read({
  output_id: "out_20250620_204000",
  options: {
    extract: "code",
    language: "javascript"
  }
});

// Write directly to file WITHOUT returning content (ZERO TOKENS!)
const result = await mcp.read({
  output_id: "out_20250620_204000",
  options: {
    extract: "code",
    write_to: "src/middleware/error-handler.js"
  }
});

// NEW: Handle multiple code blocks
const multiBlock = await mcp.read({
  output_id: "out_20250620_204000",
  options: {
    extract: "code",
    write_to: "src/components/TodoList.jsx"
  }
});
// Returns: Warning: Multiple code blocks found (3 blocks)...
// Block 0: jsx - "import React..." (4.3 KB, 150 lines)
// Block 1: jsx - "import { render }..." (1.2 KB, 45 lines)
// Block 2: css - ".todo-container..." (892 bytes, 34 lines)

// Select specific block
const component = await mcp.read({
  output_id: "out_20250620_204000",
  options: {
    extract: "code",
    write_to: "src/components/TodoList.jsx",
    block_index: 0
  }
});
```

#### **Success Response**

```javascript
// Standard response (when not using write_to)
{
  content: "const errorHandler = (err, req, res, next) => {\n  // ... code here\n}",
  truncated: false,
  tokens: 412,
  extraction: "code",
  language: "javascript"
}

// Response when using write_to (no content returned!)
{
  content: "Content written to src/middleware/error-handler.js (3.2 KB, ~800 tokens saved)",
  truncated: false,
  tokens: 0,
  extraction: "code",
  file_written: true
}

// Multi-block warning response (when extract: "code" with write_to finds multiple blocks)
{
  content: "Warning: Multiple code blocks found (3 blocks). Use block_index option to select specific block.\n\nBlock 0: jsx - \"import React, { useState } from 'react'...\" (4.3 KB, 150 lines)\nBlock 1: jsx - \"import { render } from '@testing-library/react'...\" (1.2 KB, 45 lines)\nBlock 2: css - \".todo-container { ...\" (892 bytes, 34 lines)",
  multiple_blocks: true,
  block_count: 3,
  extraction: "code"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `content` | string | The requested content (or success message when using write_to) |
| `truncated` | boolean | Whether content was cut off at max_tokens |
| `tokens` | number | Actual token count returned (0 when using write_to) |
| `extraction` | string | What was extracted ("all", "code", or "explanation") |
| `language` | string | Language filter applied (if any) |
| `file_written` | boolean | True when write_to was used successfully |
| `multiple_blocks` | boolean | True when multiple code blocks detected with extract:"code" + write_to |
| `block_count` | number | Number of code blocks found (when multiple_blocks is true) |

**Truncation Detection:**

When delegate detects that an LLM output was likely truncated mid-stream, it will:
1. Store truncation metadata (available via `delegate_check`)
2. Append a warning message with actionable hints when using `delegate_read`:
   ```
   [WARNING: Output was likely truncated (confidence: 0.85). Reason: unclosed brackets/braces]
   HINTS: Consider one of these actions:
   - Use write_to option to save the full output to disk (avoids token limits)
   - Retry with a more specific prompt asking for a smaller response
   - Use max_tokens with a higher value if you need more content
   - The response contains incomplete code/data structures
   ```

The truncation detection algorithm checks for:
- Unclosed quotes, brackets, braces, or code fences
- Content ending mid-word or mid-JSON structure
- Trailing operators or incomplete syntax
- Suspicious size boundaries (e.g., exactly 4096, 8192 bytes)

**For AI Agents:** The hints are designed to help you take appropriate action when truncation is detected. The most common solution is using the `write_to` option to save the full output directly to disk, bypassing token limits entirely.

## **Error Handling**

When a tool call fails, Delegate returns structured error information through the MCP protocol.

#### **Error Response Structure**

```javascript
{
  error: "Provider API call failed",
  code: "PROVIDER_ERROR",
  details: "API key for Gemini is invalid or has expired"
}
```

**Common Error Codes:**

| Code | Description | Action |
|------|-------------|--------|
| `INVALID_REQUEST` | Missing or invalid parameters | Check required fields |
| `INVALID_MODEL` | Model ID not recognized | Use supported model |
| `FILE_NOT_FOUND` | Context file doesn't exist | Verify file paths |
| `PROVIDER_ERROR` | LLM API returned error | Check API keys and limits |
| `OUTPUT_NOT_FOUND` | Output ID doesn't exist | Verify ID from invoke |
| `EXTRACTION_FAILED` | Could not parse LLM response | Try `extract: "all"` |
| `TIMEOUT` | Operation exceeded timeout | Increase timeout or simplify |

## **Supported Models**

| Model Identifier | Provider | Context Window | Recommended Use Case |
|------------------|----------|----------------|----------------------|
| `gemini-2.5-flash` | Google | 1M tokens | Fast, general-purpose code generation |
| `gemini-2.5-pro` | Google | 1M tokens | Complex reasoning and architecture |
| `claude-sonnet-4-20250514` | Anthropic | 200K tokens | Balanced quality and performance |
| `claude-opus-4-20250514` | Anthropic | 200K tokens | Highest quality for critical systems |

## **Best Practices**

### **1. Always Check Before Reading**
```javascript
// ❌ Bad - might consume thousands of tokens
const content = await mcp.read({ output_id });

// ✅ Good - make informed decisions
const info = await mcp.check({ output_id });
if (info.estimated_tokens > 5000) {
  // Extract only what you need
  const code = await mcp.read({ 
    output_id, 
    options: { extract: "code", max_tokens: 2000 }
  });
}
```

### **2. Use Context Files**
```javascript
// ❌ Bad - LLM lacks context
await mcp.invoke({
  model: "gemini-2.5-flash",
  prompt: "Update the API to handle the new requirements"
});

// ✅ Good - clear context improves quality
await mcp.invoke({
  model: "gemini-2.5-flash", 
  prompt: "Update the API to handle the new requirements",
  files: ["new_requirements.md", "current_api.js", "tests.js"]
});
```

### **3. Strategic Extraction**
```javascript
// First pass - get the code
const code = await mcp.read({
  output_id,
  options: { extract: "code" }
});

// Later, if needed - get explanation
const explanation = await mcp.read({
  output_id,
  options: { extract: "explanation", max_tokens: 500 }
});
```

## **Configuration**

Delegate behavior can be customized via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `DELEGATE_LOG_LEVEL` | `info` | Logging verbosity: debug, info, warn, error |
| `DELEGATE_TIMEOUT_SECONDS` | `60` | Maximum time for any invoke operation (can be overridden per request) |
| `DELEGATE_OUTPUT_DIR` | `./.delegate` | Directory for outputs and logs |

## **Output Lifecycle**

- Outputs are stored in `{DELEGATE_OUTPUT_DIR}/outputs/`
- Files are automatically cleaned up after 24 hours
- Output IDs are timestamp-based: `out_YYYYMMDD_HHMMSS`
- Each output is a complete JSON file containing the full LLM response

## **Error Handling**

Delegate uses structured errors to help Claude Code make intelligent recovery decisions:

### Error Response Format
```json
{
  "error": "provider_unavailable",
  "provider": "gemini-2.5-flash",
  "error_code": 429,
  "message": "Gemini is rate limited. Consider using Claude models or waiting 60s.",
  "retry_after": 60,
  "alternative_models": ["claude-sonnet-4-20250514", "gemini-2.5-pro"]
}
```

### Error Types
| Error Type | HTTP Code | Description | Claude Code Action |
|------------|-----------|-------------|-------------------|
| `rate_limited` | 429 | Too many requests | Try alternative model or wait |
| `provider_unavailable` | 503 | Service overloaded | Try alternative model or wait |
| `timeout` | 504 | Request took too long | Retry with simpler prompt |
| `provider_error` | 500 | Internal provider error | Retry or use alternative |
| `network_error` | - | Connection failed | Retry or inform user |

### Retry Behavior
- Automatic retry: 3 attempts with exponential backoff (1s, 2s, 4s)
- Retries happen within Delegate before returning error
- After 3 failures, structured error is returned to Claude Code

### Philosophy
Rather than complex automatic fallbacks, Delegate provides clear error information, letting Claude Code decide the best action based on context (retry, switch models, handle directly, or inform user).

---

**Next Steps:** Check out the Model Reference Card for detailed model selection guidance!