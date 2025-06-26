# MCP Tool Descriptions for Delegate

This document defines the exact descriptions that should be used when registering Delegate's tools with the MCP protocol. These descriptions are what Claude Code sees and uses to understand when/how to use each tool.

## Tool Descriptions

### delegate_invoke
```
description: "STEP 1: Delegate file generation to save tokens. Does NOT write files directly - stores in temp storage. Returns output_id for use with delegate_check then delegate_read(write_to). IMPORTANT: Use ABSOLUTE paths in 'files' parameter. Each file must be <1MB, but total can exceed."
```

### delegate_check
```
description: "STEP 2: Check delegated task status and size before retrieving. Shows token count and file size. Use this after delegate_invoke, before delegate_read. Helps avoid consuming unnecessary tokens."
```

### delegate_read
```
description: "STEP 3: Get delegated results. WORKFLOW: invoke -> check -> read. To save tokens: use 'write_to' with ABSOLUTE path to write file directly (no content returned). To get content: omit 'write_to'. IMPORTANT: Use extract:'code' to strip markdown fences for source files, 'all' keeps formatting."
```

## Why These Descriptions Work

1. **Clear workflow steps** - STEP 1 → STEP 2 → STEP 3 makes the process obvious
2. **Absolute path emphasis** - Multiple warnings prevent relative path confusion  
3. **File writing clarity** - Makes it clear that only delegate_read actually writes files
4. **Token awareness** - Emphasizes the token-saving benefit of write_to
5. **Workflow guidance** - "invoke -> check -> read" pattern is explicit
6. **Practical warnings** - Warns against common mistakes (relative paths, wrong tool for file writing)

## Implementation Notes

When registering tools in the MCP server (`internal/mcp/tools.go`):

```go
func (t *InvokeTool) Description() string {
    return "STEP 1: Delegate file generation to save tokens. Does NOT write files directly - stores in temp storage. Returns output_id for use with delegate_check then delegate_read(write_to). IMPORTANT: Use ABSOLUTE paths in 'files' parameter. Each file must be <1MB, but total can exceed."
}
```

These descriptions are designed to:
- **Prevent common AI agent mistakes** (relative paths, trying to write with invoke)
- **Guide proper workflow** (numbered steps, explicit sequence)
- **Emphasize key features** (write_to for token saving)
- **Be under 255 characters** (MCP best practice)
- **Include practical examples** (absolute vs relative paths)