# MCP Tool Descriptions for Delegate

This document defines the exact descriptions that should be used when registering Delegate's tools with the MCP protocol. These descriptions are what Claude Code sees and uses to understand when/how to use each tool.

## Tool Descriptions

### delegate_invoke
```
description: "Delegate heavy tasks (code generation, document analysis, large file processing) to other LLMs to save Claude Code's context tokens. Use this when: generating large amounts of code, analyzing multiple documents, processing entire codebases, or any task that would consume significant context. Supports Gemini models (1M token context) and Claude models. Returns an output_id for async retrieval."
```

### delegate_check
```
description: "Get metadata about a delegated task output including size, token count, and creation time. Always use this before reading to avoid consuming unnecessary tokens. Returns file size in bytes and estimated token count."
```

### delegate_read
```
description: "Retrieve results from a delegated task. Use 'extract' option to get only code or explanation. Use 'max_tokens' to limit response size. **KEY FEATURE**: Use 'write_to' to save output directly to a file WITHOUT consuming any tokens! Best practice: always check() before read() to know what you're getting."
```

## Why These Descriptions Work

1. **Clear use cases** - I immediately know when to use invoke (heavy tasks, multiple docs, large generation)
2. **Token awareness** - Emphasizes the token-saving benefit, especially the write_to feature
3. **Practical guidance** - Tells me to check before reading
4. **Context hints** - Mentions Gemini's 1M token advantage
5. **Revolutionary feature highlight** - The write_to option in delegate_read is prominently featured as a KEY FEATURE

## Implementation Notes

When registering tools in the MCP server (likely in `internal/mcp/server.go`):

```go
tools := []MCPTool{
    {
        Name: "delegate_invoke",
        Description: "Delegate heavy tasks (code generation, document analysis, large file processing) to other LLMs...",
        // ... parameters
    },
    // ... other tools
}
```

These descriptions should be:
- Under 255 characters (MCP best practice)
- Action-oriented ("Delegate heavy tasks...")
- Include concrete examples
- Mention key benefits (token saving, 1M context)