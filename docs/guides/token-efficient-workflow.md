# Token-Efficient Development Workflow with Delegate

**Reviewed by Christian Wissmann for Delegate V2.0**

## Overview

This guide presents a revolutionary workflow that enables developers to work with massive codebases while using 95%+ fewer tokens. By leveraging the delegate MCP server's new 4-tool API, you can generate, fix, and iterate on large code projects without ever loading the full code into your primary Claude context.

## The Token Problem

Traditional AI-assisted development faces a critical limitation:
- Reading a 10,000 line codebase: ~50,000 tokens
- Making changes: ~50,000 tokens (to write it back)
- Fixing errors: Another ~50,000 tokens
- **Total: 150,000+ tokens for one iteration!**

## The Delegate Solution

With delegate's new 4-tool API, the same task uses:
- Task submission: ~500 tokens
- Error analysis: ~500 tokens  
- Fix instructions: ~1,000 tokens
- **Total: ~2,000 tokens (98.7% reduction!)**

## New 4-Tool API Architecture

The delegate server now provides four distinct, single-purpose tools that eliminate ambiguity and make token costs explicit:

| Tool | Purpose | Token Cost | When to Use |
|------|---------|------------|-------------|
| `delegate_submit_task` | Submit generation task | **Low** (~50-100 tokens) | Always first step |
| `delegate_get_output_metadata` | Get structured info about output | **Low** (~20 tokens) | When unsure about output structure |
| `delegate_get_output_content` | Retrieve content into context | **High** (proportional to size) | Only when you need to read the content |
| `delegate_write_output_to_file` | Write content directly to file | **ZERO** | Most common - direct file operations |

## Core Workflow Patterns

### Pattern 1: Zero-Token Generation (Most Common)

Generate code and write it directly to disk without ever reading it:

```json
// Step 1: Submit task
{
  "tool_name": "delegate_submit_task",
  "parameters": {
    "model": "gemini-2.5-flash",
    "prompt": "Create a React TodoList component with TypeScript",
    "files": ["src/types.ts", "examples/similar-component.tsx"],
    "timeout": 90
  }
}
```

**Response:**
```json
{
  "output_id": "out_20241027_103000",
  "working_directory": "/home/user/project"
}
```

```json
// Step 2: Write directly to file (ZERO tokens!)
{
  "tool_name": "delegate_write_output_to_file", 
  "parameters": {
    "output_id": "out_20241027_103000",
    "write_to": "src/components/TodoList.tsx",
    "options": {
      "extract": "code"
    }
  }
}
```

**Response:**
```json
{
  "success": true,
  "path": "src/components/TodoList.tsx",
  "absolute_path": "/home/user/project/src/components/TodoList.tsx",
  "bytes_written": 4182,
  "message": "Successfully wrote 4.1 KB to src/components/TodoList.tsx",
  "working_directory": "/home/user/project"
}
```

### Pattern 2: Smart Multi-Block Handling

When LLMs generate multiple code blocks, delegate provides structured metadata for intelligent selection:

```json
// Step 1: Generate component with tests and styles
{
  "tool_name": "delegate_submit_task",
  "parameters": {
    "prompt": "Create a React TodoList component with unit tests and CSS styling"
  }
}
```

```json
// Step 2: Check what was generated
{
  "tool_name": "delegate_get_output_metadata",
  "parameters": {
    "output_id": "out_20241027_103000"
  }
}
```

**Response shows structured block information:**
```json
{
  "metadata": {
    "output_id": "out_20241027_103000",
    "status": "COMPLETED",
    "size_kb": 18.5,
    "line_count": 487,
    "token_estimate": 4625,
    "is_truncated": false,
    "truncation_reason": null
  },
  "content_analysis": {
    "blocks_found": 3,
    "blocks": [
      {
        "index": 0,
        "language": "tsx",
        "size_kb": 12.1,
        "lines": 250,
        "preview": "import React, { useState } from 'react';"
      },
      {
        "index": 1,
        "language": "tsx", 
        "size_kb": 3.8,
        "lines": 95,
        "preview": "import { render, screen } from '@testing-library/react';"
      },
      {
        "index": 2,
        "language": "css",
        "size_kb": 2.6,
        "lines": 142,
        "preview": ".todo-container {"
      }
    ]
  }
}
```

```json
// Step 3: Write each block to appropriate files
{
  "tool_name": "delegate_write_output_to_file",
  "parameters": {
    "output_id": "out_20241027_103000",
    "write_to": "src/components/TodoList.tsx",
    "options": {
      "block_index": 0
    }
  }
}

// Repeat for test file
{
  "tool_name": "delegate_write_output_to_file", 
  "parameters": {
    "output_id": "out_20241027_103000",
    "write_to": "src/components/TodoList.test.tsx",
    "options": {
      "block_index": 1
    }
  }
}

// And CSS file
{
  "tool_name": "delegate_write_output_to_file",
  "parameters": {
    "output_id": "out_20241027_103000", 
    "write_to": "src/components/TodoList.css",
    "options": {
      "block_index": 2
    }
  }
}
```

### Pattern 3: Compile-Fix Loop

Let compilers do the error analysis:

```bash
# Compile and capture errors
go build ./src/services/user_service.go 2> build_errors.txt
```

```json
// Fix based on compiler output
{
  "tool_name": "delegate_submit_task",
  "parameters": {
    "model": "gemini-2.5-flash",
    "files": ["src/services/user_service.go", "build_errors.txt"],
    "prompt": "Fix only these specific compilation errors"
  }
}
```

```json
// Write the fixed version directly
{
  "tool_name": "delegate_write_output_to_file",
  "parameters": {
    "output_id": "new_output_id",
    "write_to": "src/services/user_service.go", 
    "options": {
      "extract": "code"
    }
  }
}
```

## Advanced Decision-Making with Structured Responses

The new API enables sophisticated programmatic logic:

```json
// Get metadata to make informed decisions
{
  "tool_name": "delegate_get_output_metadata",
  "parameters": {"output_id": "out_123"}
}
```

**Agent can now programmatically decide:**
```javascript
// Pseudo-code for agent decision logic
if (response.content_analysis.blocks_found > 1) {
  // Handle multi-block output
  response.content_analysis.blocks.forEach((block, index) => {
    if (block.language === 'jsx') {
      writeBlock(index, `src/components/${componentName}.jsx`);
    } else if (block.language === 'css') {
      writeBlock(index, `src/styles/${componentName}.css`);
    }
  });
} else if (response.metadata.token_estimate > 10000) {
  // Large output - write directly to avoid token cost
  writeToFile();
} else {
  // Small output - safe to retrieve into context
  getContent();
}
```

## Real-World Example: Building a REST API

**Traditional approach (150,000+ tokens):**
1. Read existing codebase (50k tokens)
2. Generate new code in context (50k tokens)  
3. Fix errors through multiple iterations (50k+ tokens)

**Delegate approach (~2,000 tokens):**

```json
// 1. Generate API with context (~500 tokens for orchestration)
{
  "tool_name": "delegate_submit_task",
  "parameters": {
    "model": "gemini-2.5-flash",
    "files": ["api/openapi.yaml", "internal/base_controller.go"],
    "prompt": "Generate complete REST API handlers from OpenAPI spec"
  }
}
```

```json
// 2. Write to project (0 tokens)
{
  "tool_name": "delegate_write_output_to_file",
  "parameters": {
    "output_id": "api_gen_001",
    "write_to": "internal/api/handlers.go",
    "options": {"extract": "code"}
  }
}
```

```bash
# 3. Fix compilation errors (~500 tokens)
go build ./internal/api/... 2> errors.txt
```

```json
{
  "tool_name": "delegate_submit_task",
  "parameters": {
    "files": ["internal/api/handlers.go", "errors.txt"],
    "prompt": "Fix compilation errors"
  }
}
```

```json
{
  "tool_name": "delegate_write_output_to_file",
  "parameters": {
    "output_id": "fix_001",
    "write_to": "internal/api/handlers.go",
    "options": {"extract": "code"}
  }
}
```

```bash
# 4. Fix failing tests (~1000 tokens)
go test ./internal/api/... > test_failures.txt 2>&1
```

```json
{
  "tool_name": "delegate_submit_task", 
  "parameters": {
    "files": ["internal/api/handlers.go", "test_failures.txt"],
    "prompt": "Fix failing tests"
  }
}
```

**Total: ~2,000 tokens vs 150,000+ tokens!**

## Token Saving Strategies

### 1. Never Read Generated Code Unnecessarily
- ❌ Retrieve content just to write it elsewhere
- ✅ Use `delegate_write_output_to_file` for direct file operations

### 2. Use Metadata for Smart Decisions
- Check `delegate_get_output_metadata` first for complex outputs
- Make programmatic decisions based on structured `content_analysis`
- Avoid parsing string warnings (now eliminated!)

### 3. Leverage External Tools as Analyzers
- Compilers provide precise error locations
- Linters identify style issues  
- Test runners show exact failures
- All produce small, focused input files

### 4. Chain Focused Operations
Instead of one large prompt, chain small operations:
```
Generate core → Fix syntax → Add validation → Fix tests → Add docs
```

### 5. Model Selection Strategy
- **gemini-2.5-flash**: 95% of tasks (fast, cheap, capable)
- **gemini-2.5-pro**: Complex architectural decisions
- **claude-sonnet-4**: When you need Claude's specific strengths

## Path Handling Simplification

All paths are now relative to the delegate server's working directory:

```json
// Old way (absolute paths)
{
  "files": ["/home/user/project/src/models/user.go"],
  "write_to": "/home/user/project/src/services/user_service.go"
}

// New way (relative paths)
{
  "files": ["src/models/user.go"], 
  "write_to": "src/services/user_service.go"
}
```

Benefits:
- **Reduces cognitive load** - think in project terms
- **Prevents path errors** - no more absolute path construction
- **Saves tokens** - shorter paths in requests
- **Improves portability** - same commands work everywhere

## Best Practices

1. **Default to `delegate_write_output_to_file`**: Most operations should write directly to disk

2. **Use metadata for complex decisions**: Check `delegate_get_output_metadata` when unsure about output structure

3. **Leverage structured responses**: No more string parsing - all data is programmatically accessible

4. **Think in pipelines**: Each step transforms output for the next, like Unix pipes

5. **Trust the process**: Resist the urge to read generated code unless necessary

6. **Use relative paths**: Simpler, more portable, fewer tokens

## Error Handling

The new API provides structured error responses:

```json
{
  "error": {
    "code": "OUTPUT_NOT_FOUND",
    "message": "The requested output ID does not exist or has expired.",
    "details": {
      "output_id_provided": "invalid_id_123"
    }
  }
}
```

Common error codes:
- `INVALID_REQUEST`: Malformed parameters
- `OUTPUT_NOT_FOUND`: Invalid or expired output_id
- `PROVIDER_ERROR`: External LLM service issues
- `FILE_WRITE_FAILED`: Disk I/O problems
- `PATH_TRAVERSAL_ATTEMPT`: Security violation

## Conclusion

The new 4-tool delegate API transforms AI development from a token-intensive process to an efficient, predictable workflow. With explicit tool purposes, structured responses, and zero-token file operations, you can:

- **Work with any codebase size** without context limitations
- **Save 95%+ tokens** through direct file operations  
- **Make reliable programmatic decisions** using structured metadata
- **Handle complex multi-block outputs** intelligently
- **Build robust error handling** with structured error responses

The result: You become a conductor orchestrating development, not a code carrier burning tokens. The new API's clarity and predictability enable sophisticated automation while maintaining the speed and cost benefits that make delegate revolutionary.