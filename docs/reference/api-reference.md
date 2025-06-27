# **Delegate API Reference v2.0**

**Status:** Final | **Version:** 2.0 | **Date:** 2024-10-27

**Reviewed by Christian Wissmann for Delegate V2.0**

## **Introduction**

Welcome to the Delegate API Reference.

Delegate is a high-performance, minimalist MCP server designed to act as a token-efficient intermediary for Large Language Model (LLM) tasks. It allows Claude and other AI agents to delegate code generation and analysis to external models (Google's Gemini and Anthropic's Claude) without consuming their own context window.

As an MCP server, Delegate exposes four simple tools that AI agents can call to manage delegated tasks efficiently.

The core workflow is a simple, four-step process:

1. **submit_task**: Delegate a task to an LLM. This creates an output artifact and returns an ID.
2. **get_output_metadata**: (Optional) Inspect the metadata of the output artifact to understand its size and structure.
3. **get_output_content**: Retrieve the content into your context (consumes tokens proportional to content size).
4. **write_output_to_file**: Write content directly to a file (consumes ZERO tokens).

## **Authentication**

The Delegate MCP server requires API keys for the underlying LLM providers to be available as environment variables:

* `GOOGLE_API_KEY`: For Gemini models
* `ANTHROPIC_API_KEY`: For Claude models

These keys are never transmitted except to their respective providers over HTTPS.

## **MCP Tools**

### **delegate_submit_task**

STEP 1: Submits a generation task to an external LLM. This is an asynchronous operation that creates a temporary output artifact and returns a unique `output_id`. The content is NOT returned directly - use other `delegate_*` tools to access the output.

#### **Tool Definition**

```typescript
{
  name: "delegate_submit_task",
  description: "STEP 1: Submits a generation task to an external LLM (~50-100 tokens). Creates a temporary output artifact and returns a unique output_id. The content is NOT returned directly. Use other delegate_* tools to access the output.",
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
      description: "Relative file paths to include as context (e.g., 'src/model.go', 'docs/api.md'). Each file must be <1MB."
    },
    max_tokens: {
      type: "number",
      required: false,
      description: "Maximum tokens to generate (defaults to model maximum)"
    },
    timeout: {
      type: "number",
      required: false,
      description: "Timeout in seconds (default: 180, max: 600). Suggested: 180s minimum for code, 400s minimum for creative tasks."
    }
  }
}
```

#### **Example Usage**

```javascript
const result = await mcp.submit_task({
  model: "gemini-2.5-flash",
  prompt: "Create a robust error handling middleware for Express.js",
  files: ["src/server.js", "src/routes.js"],
  max_tokens: 4000
});
```

#### **Success Response**

```json
{
  "output_id": "out_20241027_103000",
  "working_directory": "/home/user/project"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `output_id` | string | Unique identifier for the output, used in subsequent tools |
| `working_directory` | string | The server's working directory for relative path resolution |

### **delegate_get_output_metadata**

STEP 2 (Optional): Retrieves structured metadata about an output artifact without reading its content. Essential for token-efficient operations. Use this to decide whether to retrieve content into context or write directly to a file.

#### **Tool Definition**

```typescript
{
  name: "delegate_get_output_metadata",
  description: "STEP 2 (Optional): Retrieves structured metadata about an output artifact (~20 tokens). Use this to decide whether to retrieve content into context or write directly to a file. This tool does NOT return the content itself.",
  parameters: {
    output_id: {
      type: "string",
      required: true,
      description: "The ID returned from delegate_submit_task"
    }
  }
}
```

#### **Example Usage**

```javascript
const metadata = await mcp.get_output_metadata({
  output_id: "out_20241027_103000"
});
```

#### **Success Response**

```json
{
  "metadata": {
    "output_id": "out_20241027_103000",
    "status": "COMPLETED",
    "size_kb": 15.7,
    "line_count": 312,
    "token_estimate": 3925,
    "is_truncated": false,
    "truncation_reason": null
  },
  "content_analysis": {
    "blocks_found": 2,
    "blocks": [
      {
        "index": 0,
        "language": "jsx",
        "size_kb": 12.1,
        "lines": 250,
        "preview": "import React from 'react';"
      },
      {
        "index": 1,
        "language": "md",
        "size_kb": 3.6,
        "lines": 62,
        "preview": "# Explanation of the component"
      }
    ]
  }
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `metadata.output_id` | string | The output identifier |
| `metadata.status` | string | Generation status: "COMPLETED", "IN_PROGRESS", or "FAILED" |
| `metadata.size_kb` | number | Size of the output in kilobytes |
| `metadata.line_count` | number | Total number of lines in the output |
| `metadata.token_estimate` | number | Estimated token count for the full content |
| `metadata.is_truncated` | boolean | Whether the LLM output was truncated |
| `metadata.truncation_reason` | string/null | Reason for truncation if applicable |
| `content_analysis.blocks_found` | number | Number of distinct content blocks detected |
| `content_analysis.blocks` | array | Array of block metadata objects |
| `content_analysis.blocks[].index` | number | Zero-based index for referencing this block |
| `content_analysis.blocks[].language` | string | Programming language or content type |
| `content_analysis.blocks[].size_kb` | number | Size of this block in kilobytes |
| `content_analysis.blocks[].lines` | number | Number of lines in this block |
| `content_analysis.blocks[].preview` | string | First line of the block content |

### **delegate_get_output_content**

STEP 3A: Retrieves the full or partial content of an output artifact into the agent's context. This operation consumes tokens proportional to the content size. Use `options` to extract specific parts.

#### **Tool Definition**

```typescript
{
  name: "delegate_get_output_content",
  description: "Retrieves the full or partial content of an output artifact into the agent's context (~30+ tokens plus content). This operation consumes tokens proportional to the content size. Use options to extract specific parts (e.g., extract: 'code').",
  parameters: {
    output_id: {
      type: "string",
      required: true,
      description: "The ID returned from delegate_submit_task"
    },
    options: {
      type: "object",
      required: false,
      properties: {
        extract: {
          type: "string",
          enum: ["all", "code", "explanation"],
          default: "all",
          description: "What part to extract from the output"
        },
        max_tokens: {
          type: "number",
          description: "Truncate the returned content to this many tokens"
        },
        block_index: {
          type: "number",
          description: "For multi-block outputs, select a specific block (0-based index)"
        },
        language: {
          type: "string",
          description: "Filter code blocks by this programming language"
        }
      }
    }
  }
}
```

#### **Example Usage**

```javascript
// Get all content
const full = await mcp.get_output_content({
  output_id: "out_20241027_103000"
});

// Get only code, max 2000 tokens
const code = await mcp.get_output_content({
  output_id: "out_20241027_103000",
  options: {
    extract: "code",
    max_tokens: 2000
  }
});

// Get specific block from multi-block output
const component = await mcp.get_output_content({
  output_id: "out_20241027_103000",
  options: {
    extract: "code",
    block_index: 0
  }
});

// Get only JavaScript code
const jsOnly = await mcp.get_output_content({
  output_id: "out_20241027_103000",
  options: {
    extract: "code",
    language: "javascript"
  }
});
```

#### **Success Response**

```json
{
  "content": "import React from 'react';\n\nconst ErrorHandler = ({ error, retry }) => {\n  // Component implementation here\n};\n\nexport default ErrorHandler;",
  "metadata": {
    "output_id": "out_20241027_103000",
    "tokens_returned": 3925,
    "is_truncated": false,
    "truncation_reason": null
  }
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `content` | string | The requested content |
| `metadata.output_id` | string | The output identifier |
| `metadata.tokens_returned` | number | Number of tokens in the returned content |
| `metadata.is_truncated` | boolean | Whether content was truncated by max_tokens parameter |
| `metadata.truncation_reason` | string/null | Reason for truncation (e.g., "MAX_TOKENS_REACHED") |

### **delegate_write_output_to_file**

STEP 3B: Writes the content of an output artifact directly to a specified file path (relative to working directory). This operation consumes ZERO content tokens. Use `options` to select specific parts to write.

#### **Tool Definition**

```typescript
{
  name: "delegate_write_output_to_file",
  description: "Writes the content of an output artifact directly to a specified file path (relative to working directory). This operation consumes ZERO content tokens. Use options to select specific parts to write (e.g., extract: 'code', block_index: 0).",
  parameters: {
    output_id: {
      type: "string",
      required: true,
      description: "The ID returned from delegate_submit_task"
    },
    write_to: {
      type: "string",
      required: true,
      description: "Relative file path to write to (e.g., 'src/component.jsx', 'tmp/output.go')"
    },
    options: {
      type: "object",
      required: false,
      properties: {
        extract: {
          type: "string",
          enum: ["all", "code", "explanation"],
          default: "all",
          description: "What part to extract from the output"
        },
        block_index: {
          type: "number",
          description: "For multi-block outputs, select a specific block (0-based index)"
        },
        language: {
          type: "string",
          description: "Filter code blocks by this programming language"
        }
      }
    }
  }
}
```

#### **Example Usage**

```javascript
// Write all content to file
const result = await mcp.write_output_to_file({
  output_id: "out_20241027_103000",
  write_to: "src/ErrorHandler.jsx"
});

// Write only code to file
const codeResult = await mcp.write_output_to_file({
  output_id: "out_20241027_103000",
  write_to: "src/components/ErrorHandler.jsx",
  options: {
    extract: "code"
  }
});

// Write specific block from multi-block output
const blockResult = await mcp.write_output_to_file({
  output_id: "out_20241027_103000",
  write_to: "src/components/TodoList.jsx",
  options: {
    extract: "code",
    block_index: 0
  }
});

// Write only JavaScript code
const jsResult = await mcp.write_output_to_file({
  output_id: "out_20241027_103000",
  write_to: "src/utils/helpers.js",
  options: {
    extract: "code",
    language: "javascript"
  }
});
```

#### **Success Response**

```json
{
  "success": true,
  "path": "src/components/ErrorHandler.jsx",
  "absolute_path": "/home/user/project/src/components/ErrorHandler.jsx",
  "bytes_written": 12388,
  "message": "Successfully wrote 12.1 KB to src/components/ErrorHandler.jsx",
  "working_directory": "/home/user/project"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `success` | boolean | Whether the write operation succeeded |
| `path` | string | The relative path of the file written |
| `absolute_path` | string | The absolute path of the file written |
| `bytes_written` | number | Number of bytes written to the file |
| `message` | string | Human-readable success message |
| `working_directory` | string | The server's working directory |

## **Error Handling**

When a tool call fails, Delegate returns structured error information through the MCP protocol.

#### **Error Response Structure**

```json
{
  "error": {
    "code": "OUTPUT_NOT_FOUND",
    "message": "The requested output ID does not exist or has expired.",
    "details": {
      "output_id_provided": "gen-this-is-fake"
    }
  }
}
```

**Common Error Codes:**

| Code | Description | Action |
|------|-------------|--------|
| `INVALID_REQUEST` | Missing or invalid parameters | Check required fields |
| `INVALID_MODEL` | Model ID not recognized | Use supported model |
| `FILE_NOT_FOUND` | Context file doesn't exist | Verify file paths |
| `PROVIDER_ERROR` | LLM API returned error | Check API keys and limits |
| `OUTPUT_NOT_FOUND` | Output ID doesn't exist | Verify ID from submit_task |
| `FILE_WRITE_FAILED` | Could not write to specified path | Check file permissions |
| `PATH_TRAVERSAL_ATTEMPT` | Illegal path detected | Use relative paths within project |
| `TIMEOUT` | Operation exceeded timeout | Increase timeout or simplify |

## **Supported Models**

| Model Identifier | Provider | Context Window | Recommended Use Case |
|------------------|----------|----------------|----------------------|
| `gemini-2.5-flash` | Google | 1M tokens | Fast, general-purpose code generation |
| `gemini-2.5-pro` | Google | 1M tokens | Complex reasoning and architecture |
| `claude-sonnet-4-20250514` | Anthropic | 200K tokens | Balanced quality and performance |
| `claude-opus-4-20250514` | Anthropic | 200K tokens | Highest quality for critical systems |

## **Workflow Examples**

### **Scenario A: Zero-Token Code Generation**
**Goal:** Generate code and save it directly to a file without consuming any tokens.

```javascript
// Step 1: Submit the task
const task = await mcp.submit_task({
  model: "gemini-2.5-flash",
  prompt: "Create a React component for handling user authentication",
  files: ["src/types.ts", "src/api.js"]
});

// Step 2: Write directly to file (ZERO tokens consumed)
const result = await mcp.write_output_to_file({
  output_id: task.output_id,
  write_to: "src/components/AuthHandler.jsx",
  options: {
    extract: "code"
  }
});

console.log(result.message); // "Successfully wrote 8.3 KB to src/components/AuthHandler.jsx"
```

### **Scenario B: Intelligent Multi-Block Handling**
**Goal:** Understand the output structure before deciding which parts to save.

```javascript
// Step 1: Submit the task
const task = await mcp.submit_task({
  model: "gemini-2.5-pro",
  prompt: "Create a complete React component with tests and documentation",
  files: ["src/existing-component.jsx"]
});

// Step 2: Check what was generated
const metadata = await mcp.get_output_metadata({
  output_id: task.output_id
});

// Step 3: Programmatically handle different blocks
for (const block of metadata.content_analysis.blocks) {
  if (block.language === "jsx") {
    // Save the component
    await mcp.write_output_to_file({
      output_id: task.output_id,
      write_to: `src/components/NewComponent.jsx`,
      options: {
        block_index: block.index
      }
    });
  } else if (block.language === "javascript" && block.preview.includes("test")) {
    // Save the tests
    await mcp.write_output_to_file({
      output_id: task.output_id,
      write_to: `src/components/__tests__/NewComponent.test.js`,
      options: {
        block_index: block.index
      }
    });
  } else if (block.language === "md") {
    // Save the documentation
    await mcp.write_output_to_file({
      output_id: task.output_id,
      write_to: `docs/NewComponent.md`,
      options: {
        block_index: block.index
      }
    });
  }
}
```

### **Scenario C: Token-Efficient Content Review**
**Goal:** Review small parts of the output before deciding what to do with the full content.

```javascript
// Step 1: Submit the task
const task = await mcp.submit_task({
  model: "claude-sonnet-4-20250514",
  prompt: "Refactor this legacy code to use modern patterns",
  files: ["src/legacy-module.js"]
});

// Step 2: Check the size first
const metadata = await mcp.get_output_metadata({
  output_id: task.output_id
});

if (metadata.metadata.token_estimate > 5000) {
  // Large output - get just the explanation first
  const explanation = await mcp.get_output_content({
    output_id: task.output_id,
    options: {
      extract: "explanation",
      max_tokens: 500
    }
  });
  
  console.log("Summary:", explanation.content);
  
  // If satisfied, write the code directly to file
  await mcp.write_output_to_file({
    output_id: task.output_id,
    write_to: "src/modern-module.js",
    options: {
      extract: "code"
    }
  });
} else {
  // Small output - safe to get everything
  const fullContent = await mcp.get_output_content({
    output_id: task.output_id
  });
  
  console.log("Full output:", fullContent.content);
}
```

## **Best Practices**

### **1. Use Metadata for Smart Decisions**
```javascript
// ❌ Bad - might consume thousands of tokens
const content = await mcp.get_output_content({ output_id });

// ✅ Good - make informed decisions
const info = await mcp.get_output_metadata({ output_id });
if (info.metadata.token_estimate > 5000) {
  // Write directly to file to save tokens
  await mcp.write_output_to_file({
    output_id,
    write_to: "output.txt",
    options: { extract: "code" }
  });
} else {
  // Safe to get content
  const content = await mcp.get_output_content({ output_id });
}
```

### **2. Leverage Context Files**
```javascript
// ❌ Bad - LLM lacks context
await mcp.submit_task({
  model: "gemini-2.5-flash",
  prompt: "Update the API to handle the new requirements"
});

// ✅ Good - clear context improves quality
await mcp.submit_task({
  model: "gemini-2.5-flash", 
  prompt: "Update the API to handle the new requirements",
  files: ["docs/new-requirements.md", "src/current-api.js", "tests/api.test.js"]
});
```

### **3. Handle Multi-Block Outputs Gracefully**
```javascript
// Get metadata first to understand the structure
const metadata = await mcp.get_output_metadata({ output_id });

if (metadata.content_analysis.blocks_found > 1) {
  // Handle each block appropriately
  for (const block of metadata.content_analysis.blocks) {
    const filename = `output_${block.index}.${block.language}`;
    await mcp.write_output_to_file({
      output_id,
      write_to: filename,
      options: { block_index: block.index }
    });
  }
} else {
  // Single block - write directly
  await mcp.write_output_to_file({
    output_id,
    write_to: "output.txt"
  });
}
```

### **4. Use Relative Paths**
```javascript
// ✅ Good - relative paths from project root
const files = [
  "src/models/user.go",
  "docs/api-spec.md",
  "config/database.json"
];

// ✅ Good - relative output paths
await mcp.write_output_to_file({
  output_id,
  write_to: "src/handlers/new-endpoint.go"
});
```

## **Configuration**

Delegate behavior can be customized via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `DELEGATE_LOG_LEVEL` | `info` | Logging verbosity: debug, info, warn, error |
| `DELEGATE_TIMEOUT_SECONDS` | `180` | Default timeout for submit_task operations |
| `DELEGATE_OUTPUT_DIR` | `./.delegate` | Directory for outputs and logs |
| `DELEGATE_WORKING_DIR` | `./` | Working directory for relative path resolution |

## **Output Lifecycle**

- Outputs are stored in `{DELEGATE_OUTPUT_DIR}/outputs/`
- Files are automatically cleaned up after 24 hours
- Output IDs are timestamp-based: `out_YYYYMMDD_HHMMSS`
- Each output is a complete JSON file containing the full LLM response and metadata

## **Security Considerations**

- All file paths are resolved relative to the configured working directory
- Path traversal attempts (e.g., `../../../etc/passwd`) are blocked
- File writes are sandboxed to the project directory
- Context files are validated for size limits (<1MB each)
- API keys are never logged or transmitted except to their respective providers

---

**Migration Note:** This API replaces the previous 3-tool architecture. The new 4-tool design provides clearer separation of concerns and more predictable token usage patterns.