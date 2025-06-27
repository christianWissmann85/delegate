# MCP Tool Descriptions for Delegate

**Reviewed by Christian Wissmann for Delegate V2.0**

This document defines the exact descriptions that should be used when registering Delegate's tools with the MCP protocol. These descriptions are what Claude Code sees and uses to understand when/how to use each tool.

## Tool Descriptions

### delegate_submit_task
```
description: "STEP 1: Submits a generation task to an external LLM (~50-100 tokens). This is an asynchronous operation that creates a temporary output artifact and returns a unique `output_id`. The content is NOT returned directly. Use other `delegate_*` tools to access the output."
```

### delegate_get_output_metadata
```
description: "STEP 2 (Optional): Retrieves structured metadata about an output artifact (~20 tokens). Use this to decide whether to retrieve content into context or write directly to a file. This tool does NOT return the content itself."
```

### delegate_get_output_content
```
description: "Retrieves the full or partial content of an output artifact into the agent's context (~30+ tokens plus content). This operation consumes tokens proportional to the content size. Use `options` to extract specific parts (e.g., `extract: 'code'`)."
```

### delegate_write_output_to_file
```
description: "Writes the content of an output artifact directly to a specified file path (relative to working directory). This operation consumes ZERO content tokens. Use `options` to select specific parts to write (e.g., `extract: 'code'`, `block_index: 0`)."
```

## Why These Descriptions Work

1. **Clear workflow steps** - STEP 1 → STEP 2 (optional) → final action makes the process obvious
2. **Relative path clarity** - Specifies "relative to working directory" to prevent path confusion
3. **Token cost transparency** - Each tool clearly states its token consumption impact
4. **Purpose-driven design** - Each tool has one clear, unambiguous function
5. **Workflow guidance** - Natural progression from submit → check → retrieve/write
6. **Zero-token emphasis** - Highlights the key benefit of write_output_to_file

## Implementation Notes

When registering tools in the MCP server (`internal/mcp/tools.go`):

```go
func (t *SubmitTaskTool) Description() string {
    return "STEP 1: Submits a generation task to an external LLM (~50-100 tokens). This is an asynchronous operation that creates a temporary output artifact and returns a unique `output_id`. The content is NOT returned directly. Use other `delegate_*` tools to access the output."
}

func (t *GetMetadataTool) Description() string {
    return "STEP 2 (Optional): Retrieves structured metadata about an output artifact (~20 tokens). Use this to decide whether to retrieve content into context or write directly to a file. This tool does NOT return the content itself."
}

func (t *GetContentTool) Description() string {
    return "Retrieves the full or partial content of an output artifact into the agent's context (~30+ tokens plus content). This operation consumes tokens proportional to the content size. Use `options` to extract specific parts (e.g., `extract: 'code'`)."
}

func (t *WriteFileTool) Description() string {
    return "Writes the content of an output artifact directly to a specified file path (relative to working directory). This operation consumes ZERO content tokens. Use `options` to select specific parts to write (e.g., `extract: 'code'`, `block_index: 0`)."
}
```

These descriptions are designed to:
- **Guide proper workflow** - Clear progression from task submission to final action
- **Prevent token waste** - Emphasize when operations consume tokens vs. when they don't
- **Enable smart decisions** - Metadata step helps agents choose the right final action
- **Use relative paths** - Simpler, more natural path handling
- **Be concise and actionable** - Each description fits MCP best practices while being maximally useful

## Common Workflows

**Zero-Token File Generation:**
1. `delegate_submit_task` → get output_id
2. `delegate_write_output_to_file` → write directly to disk (0 tokens)

**Content Review Before Writing:**
1. `delegate_submit_task` → get output_id  
2. `delegate_get_output_metadata` → check structure/size
3. `delegate_write_output_to_file` → write selected content (0 tokens)

**Content Integration:**
1. `delegate_submit_task` → get output_id
2. `delegate_get_output_content` → bring into context (high tokens)