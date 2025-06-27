# CLAUDE.md - Instructions for Claude Code Instances

Welcome, fellow Claude Code instance! 👋 You're working on **delegate** - the MCP server that lets us delegate heavy tasks to other LLMs to save our precious context tokens.

## Your Role

You are developing and maintaining delegate, an MCP server that YOU (and other Claude Code instances) use as a power tool. Humans don't use delegate directly - they lack MCP server connections. This is our tool, by us, for us!

## What is Delegate?

Delegate is an MCP (Model Context Protocol) server that provides four simple tools:
1. **delegate_submit_task** - Send tasks to other LLMs (Gemini/Claude models)
2. **delegate_get_output_metadata** - Check output size/tokens before retrieving
3. **delegate_get_output_content** - Retrieve results into your context
4. **delegate_write_output_to_file** - Write results directly to disk (saving tokens!)

**The Magic**: The `delegate_write_output_to_file` tool lets you save massive outputs directly to disk WITHOUT consuming any tokens - achieving 95%+ token savings!

## Core Philosophy

Before making ANY changes, read `docs/development/NO_SCOPE_CREEP.md`. This project is intentionally minimal:
- ✅ Only 4 tools: submit, metadata, content, write
- ✅ Stateless operations only
- ✅ Local filesystem storage
- ❌ NO session management
- ❌ NO progress tracking
- ❌ NO web UI or CLI tools
- ❌ NO cloud storage

## Development Setup

```bash
# You're already in the repo, but for reference:
cd /home/chris/repos/delegate

# Build
go build -o delegate main.go

# Run tests
go test ./...
go test -v -tags=e2e ./e2e/...  # E2E tests

# The delegate MCP server is already installed in Claude Code
# After changes, ask the human to restart: "Please restart the delegate MCP server"
```

## Using Delegate While Developing It

Yes, you can use delegate to help develop delegate! 🤯

```bash
# Example: Generate a new test file
delegate_submit_task(model: "gemini-2.5-flash", prompt: "Create comprehensive unit tests for the multi-block handling feature")
# Response: {"output_id": "out_20241027_103000", "working_directory": "/home/chris/repos/delegate"}

delegate_get_output_metadata(output_id: "out_20241027_103000")
# Response shows structured metadata about the output

delegate_write_output_to_file(output_id: "out_20241027_103000", write_to: "internal/handlers/multiblock_test.go", options: {extract: "code"})
# Response: {"success": true, "path": "internal/handlers/multiblock_test.go", "bytes_written": 4388, ...}
```

## Key Features You Should Know

### 1. Multi-Block Handling with Structured Metadata

When LLMs return multiple code blocks, `delegate_get_output_metadata` returns structured JSON:

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
    "blocks_found": 4,
    "blocks": [
      {"index": 0, "language": "go", "size_kb": 4.3, "lines": 150, "preview": "package main"},
      {"index": 1, "language": "go", "size_kb": 1.2, "lines": 45, "preview": "package main_test"},
      {"index": 2, "language": "yaml", "size_kb": 0.456, "lines": 23, "preview": "version: '3.8'"},
      {"index": 3, "language": "markdown", "size_kb": 0.892, "lines": 34, "preview": "# Usage Instructions"}
    ]
  }
}
```

Then use `block_index` to select the one you want!

### 2. Token-Free File Writing

```bash
# This writes directly to disk - you never see the content!
delegate_write_output_to_file(
  output_id: "out_20241027_103000", 
  write_to: "src/component.jsx",
  options: {extract: "code", block_index: 0}
)
# Response: {"success": true, "path": "src/component.jsx", "bytes_written": 15234, "message": "Successfully wrote 14.9 KB to src/component.jsx"}
```

### 3. Smart Code Extraction

- `extract: "code"` - Strips markdown fences for clean source files
- `extract: "all"` - Keeps original formatting (good for docs)
- `extract: "explanation"` - Gets only the explanatory text
- `block_index: N` - Select specific block from multi-block outputs
- `language: "jsx"` - Filter blocks by language

### 4. Relative Path Simplicity

All paths are now relative to the working directory:
```bash
# Before (old API): "/home/chris/repos/delegate/src/handlers/test.go"
# Now: "src/handlers/test.go"

delegate_submit_task(
  model: "gemini-2.5-flash",
  prompt: "Refactor this handler",
  files: ["src/handlers/old.go", "docs/patterns.md"]  # Simple relative paths!
)
```

## Common Development Tasks

### Adding a New Feature
1. First, ask yourself: "Does this violate NO_SCOPE_CREEP.md?"
2. If yes, stop. If no, proceed.
3. Use delegate to generate boilerplate/tests
4. Implement the feature
5. Run tests: `go test ./...`

### Fixing Bugs
1. Reproduce the issue
2. Write a failing test
3. Fix the bug
4. Verify all tests pass

### Updating Documentation
Always update docs when changing functionality:
- `docs/reference/api-reference.md` - API changes
- `docs/guides/token-efficient-workflow.md` - Usage patterns
- `docs/development/` - Development notes

## Project Structure

```
delegate/
├── main.go               # Entry point
├── internal/
│   ├── mcp/             # MCP protocol layer
│   ├── handlers/        # submit, metadata, content, write implementations
│   ├── providers/       # LLM integrations (Anthropic, Google)
│   ├── extractor/       # Code/text extraction logic
│   ├── storage/         # File system operations
│   └── models/          # Shared data structures and responses
└── docs/
    ├── guides/          # User guides (for Claude Code instances!)
    ├── reference/       # API documentation
    └── development/     # Development docs
```

## Testing Your Changes

```bash
# Unit tests
go test ./internal/handlers/...

# Integration tests  
go test ./internal/providers/...

# Full E2E tests
go test -v -tags=e2e ./e2e/...

# Quick manual test with new API
delegate_submit_task(model: "gemini-2.5-flash", prompt: "Say hello")
delegate_get_output_metadata(output_id: "out_...")
delegate_get_output_content(output_id: "out_...")
```

## Important Implementation Details

1. **Storage**: Outputs stored in `.delegate/outputs/` with 24-hour auto-cleanup
2. **Models**: Only Gemini (Flash/Pro) and Claude (Sonnet/Opus 4)
3. **Security**: Path traversal protection, input validation
4. **Performance**: Metadata <100ms, Content retrieval <500ms, Streaming for submit

## When You're Stuck

1. Check existing patterns in the codebase
2. Read the comprehensive docs in `docs/`
3. Use delegate itself to explore solutions
4. Ask the human to restart the MCP server after changes

## AI-to-AI Collaborative Problem Solving

One of the most powerful patterns with delegate is orchestrating multiple AI models to solve complex problems together. Here's how to leverage this effectively:

### The Multi-AI Discussion Pattern

1. **Define the Problem Clearly**
   ```bash
   # Start with a clear problem statement and examples
   delegate_submit_task(
     model: "gemini-2.5-flash",
     prompt: "Here's the truncation problem we need to solve...",
     files: ["internal/handlers/read.go", "test/examples/truncated.md"]
   )
   ```

2. **Save Each AI's Response**
   ```bash
   # Get the output ID
   # Response: {"output_id": "gen-abc-123"}
   
   # Save response to file without consuming tokens
   delegate_write_output_to_file(
     output_id: "gen-abc-123",
     write_to: "tmp/gemini_solution.md"
   )
   ```

3. **Facilitate Cross-Pollination**
   ```bash
   # Share all responses with the next AI
   delegate_submit_task(
     model: "claude-opus-4-20250514",
     prompt: "Here's Gemini's solution. What are your thoughts? Any improvements?",
     files: [
       "tmp/gemini_solution.md",
       "docs/problem_context.md",
       "internal/handlers/read.go"
     ]
   )
   
   # Save Opus's response
   delegate_write_output_to_file(
     output_id: "gen-def-456",
     write_to: "tmp/opus_solution.md"
   )
   ```

4. **Act as the Moderator/Facilitator**
   - Summarize key insights from each AI
   - Ask targeted follow-up questions
   - Guide the discussion toward practical solutions
   - Synthesize the best ideas from all participants

### Example Workflow: Truncation Detection

```bash
# Round 1: Initial solutions
delegate_submit_task(model: "gemini-2.5-flash", prompt: "Solve truncation detection...")
delegate_write_output_to_file(output_id: "gen-001", write_to: "tmp/gemini_solution.md")

delegate_submit_task(model: "claude-opus-4-20250514", prompt: "Solve truncation detection...")
delegate_write_output_to_file(output_id: "gen-002", write_to: "tmp/opus_solution.md")

# Round 2: Cross-pollination
delegate_submit_task(
  model: "gemini-2.5-flash",
  prompt: "Here are both solutions. Create a unified approach...",
  files: ["tmp/gemini_solution.md", "tmp/opus_solution.md", "internal/handlers/invoke.go", "internal/handlers/read.go"]
)
delegate_write_output_to_file(output_id: "gen-003", write_to: "tmp/unified_solution.md")

# Round 3: Final synthesis (you as the moderator)
# Read selected parts using delegate_get_output_content with specific options
delegate_get_output_content(output_id: "gen-003", options: {extract: "code", block_index: 0})
```

### Key Success Factors

1. **Context is Everything**
   - Always attach relevant code files
   - Include design documents and READMEs
   - Share previous AI responses in full
   - Provide concrete examples of the problem

2. **Clear Moderation**
   - Summarize what each AI contributed
   - Ask specific questions for each round
   - Guide toward practical, implementable solutions
   - Know when to stop discussing and start implementing

3. **Token Efficiency**
   - Use `delegate_write_output_to_file` to save all responses without consuming tokens
   - Use `delegate_get_output_metadata` to understand outputs before retrieving
   - Only retrieve small portions with `delegate_get_output_content` when needed for synthesis
   - Let each AI see the full context without you having to hold it all

4. **Leverage Each Model's Strengths**
   - Gemini 2.5 Flash: Fast, comprehensive initial solutions
   - Claude Opus 4: Deep analysis, pragmatic improvements
   - Claude Sonnet 4: Balanced implementation focus
   - You (Claude Code): Final synthesis and implementation

### When to Use This Pattern

- **Complex algorithmic problems** (like truncation detection)
- **Architecture decisions** requiring multiple perspectives
- **Code review** where different models might catch different issues
- **Research tasks** where breadth of knowledge helps
- **Design discussions** benefiting from diverse approaches

This collaborative approach turns delegate into a powerful "AI team" orchestration tool, where you act as the technical lead coordinating multiple AI engineers to solve problems beyond what any single model could achieve alone!

## Remember

You're not just using delegate - you're making it better for all Claude Code instances! Every improvement helps your future self and your fellow instances save tokens and work more efficiently.

Happy coding! 🚀

---
*P.S. If a human is reading this: Hi! This file is for Claude Code instances working on delegate. You can read it, but you can't use delegate directly - you need an MCP client like Claude Code!*