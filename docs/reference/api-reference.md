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

### **delegate.invoke**

Delegates a generation task to a specified LLM. This is an asynchronous operation that creates a persistent output artifact and returns a unique ID for it.

#### **Tool Definition**

```typescript
{
  name: "delegate.invoke",
  description: "Generate code or content using an external LLM",
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
      description: "File paths to include as context."
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
      description: "Request-specific timeout in seconds (overrides DELEGATE_TIMEOUT_SECONDS)"
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

### **delegate.check**

Retrieves metadata about a previously generated output artifact without reading its content. Essential for token-efficient operations.

#### **Tool Definition**

```typescript
{
  name: "delegate.check",
  description: "Get metadata about a generated output without reading content",
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

### **delegate.read**

Reads the content of an output artifact, with powerful options for extraction and truncation.

#### **Tool Definition**

```typescript
{
  name: "delegate.read",
  description: "Read content from a generated output",
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
```

#### **Success Response**

```javascript
{
  content: "const errorHandler = (err, req, res, next) => {\n  // ... code here\n}",
  truncated: false,
  tokens: 412,
  extraction: "code",
  language: "javascript"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `content` | string | The requested content |
| `truncated` | boolean | Whether content was cut off at max_tokens |
| `tokens` | number | Actual token count returned |
| `extraction` | string | What was extracted ("all", "code", or "explanation") |
| `language` | string | Language filter applied (if any) |

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