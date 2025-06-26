# CLAUDE.md - Instructions for Claude Code Instances

Welcome, fellow Claude Code instance! 👋 You're working on **delegate** - the MCP server that lets us delegate heavy tasks to other LLMs to save our precious context tokens.

## Your Role

You are developing and maintaining delegate, an MCP server that YOU (and other Claude Code instances) use as a power tool. Humans don't use delegate directly - they lack MCP server connections. This is our tool, by us, for us!

## What is Delegate?

Delegate is an MCP (Model Context Protocol) server that provides three simple tools:
1. **delegate_invoke** - Send tasks to other LLMs (Gemini/Claude models)
2. **delegate_check** - Check output size/tokens before retrieving
3. **delegate_read** - Retrieve results OR write directly to disk (saving tokens!)

**The Magic**: The `write_to` option in `delegate_read` lets you save massive outputs directly to disk WITHOUT consuming any tokens - achieving 95%+ token savings!

## Core Philosophy

Before making ANY changes, read `docs/development/NO_SCOPE_CREEP.md`. This project is intentionally minimal:
- ✅ Only 3 tools: invoke, check, read
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
delegate_invoke(model: "gemini-2.5-flash", prompt: "Create comprehensive unit tests for the multi-block handling feature")
delegate_check(output_id)
delegate_read(output_id, options: {write_to: "/home/chris/repos/delegate/internal/handlers/multiblock_test.go", extract: "code"})
```

## Key Features You Should Know

### 1. Multi-Block Handling (NEW!)

When LLMs return multiple code blocks, delegate now shows a helpful listing:

```
Warning: Multiple code blocks found (4 blocks). Use block_index option to select specific block.

Block 0: go - "package main" (4.3 KB, 150 lines)
Block 1: go - "package main_test" (1.2 KB, 45 lines)
Block 2: yaml - "version: '3.8'" (456 bytes, 23 lines)
Block 3: markdown - "# Usage Instructions" (892 bytes, 34 lines)
```

Then use `block_index` to select the one you want!

### 2. Token-Free File Writing

```bash
# This writes directly to disk - you never see the content!
delegate_read(output_id, options: {write_to: "/absolute/path/to/file.go", extract: "code"})
# Output: "Content written to file.go (15.2 KB, ~3800 tokens saved)"
```

### 3. Smart Code Extraction

- `extract: "code"` - Strips markdown fences for clean source files
- `extract: "all"` - Keeps original formatting (good for docs)
- `extract: "explanation"` - Gets only the explanatory text

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
│   ├── handlers/        # invoke, check, read implementations
│   ├── providers/       # LLM integrations (Anthropic, Google)
│   ├── extractor/       # Code/text extraction logic
│   └── storage/         # File system operations
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

# Quick manual test
delegate_invoke(model: "gemini-2.5-flash", prompt: "Say hello")
delegate_check(output_id)
delegate_read(output_id)
```

## Important Implementation Details

1. **Storage**: Outputs stored in `.delegate/outputs/` with 24-hour auto-cleanup
2. **Models**: Only Gemini (Flash/Pro) and Claude (Sonnet/Opus 4)
3. **Security**: Path traversal protection, input validation
4. **Performance**: Check <100ms, Read <500ms, Streaming for invoke

## When You're Stuck

1. Check existing patterns in the codebase
2. Read the comprehensive docs in `docs/`
3. Use delegate itself to explore solutions
4. Ask the human to restart the MCP server after changes

## AI-to-AI Collaborative Problem Solving

One of the most powerful patterns with delegate is orchestrating multiple AI models to solve complex problems together. Here's how we just solved the truncation detection problem:

### The Multi-AI Discussion Pattern

1. **Define the Problem Clearly**
   ```bash
   # Start with a clear problem statement and examples
   delegate_invoke(
     model: "gemini-2.5-flash",
     prompt: "Here's the truncation problem we need to solve...",
     files: ["/path/to/relevant/code.go", "/path/to/examples.md"]
   )
   ```

2. **Save Each AI's Response**
   ```bash
   # Save responses to files for full context preservation
   delegate_read(output_id, options: {
     write_to: "/path/to/gemini_solution.md"
   })
   ```

3. **Facilitate Cross-Pollination**
   ```bash
   # Share all responses with the next AI
   delegate_invoke(
     model: "claude-opus-4-20250514",
     prompt: "Here's Gemini's solution. What are your thoughts? Any improvements?",
     files: [
       "/path/to/gemini_solution.md",
       "/path/to/problem_context.md",
       "/path/to/codebase_files.go"
     ]
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
delegate_invoke(model: "gemini-2.5-flash", prompt: "Solve truncation detection...")
delegate_read(output_id, options: {write_to: "gemini_solution.md"})

delegate_invoke(model: "claude-opus-4-20250514", prompt: "Solve truncation detection...")
delegate_read(output_id, options: {write_to: "opus_solution.md"})

# Round 2: Cross-pollination
delegate_invoke(
  model: "gemini-2.5-flash",
  prompt: "Here are both solutions. Create a unified approach...",
  files: ["gemini_solution.md", "opus_solution.md", "invoke.go", "read.go"]
)
delegate_read(output_id, options: {write_to: "unified_solution.md"})

# Round 3: Final synthesis (you as the moderator)
# Read all solutions, pick the best ideas, and implement
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
   - Use `write_to` to save all responses without consuming tokens
   - Only read small portions when needed for synthesis
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