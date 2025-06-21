# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Delegate is an MCP (Model Context Protocol) server that allows Claude Code to delegate heavy tasks (code generation, document analysis, large file processing) to other LLMs (Gemini and Claude models) to save context tokens. 

**Current Status**: ✅ **SHIPPED and PRODUCTION READY!** All features implemented, tested, and working beautifully.

**Revolutionary Feature**: The `write_to` option in `delegate_read` lets you save massive outputs directly to disk WITHOUT consuming any tokens - achieving 95%+ token savings!

**Core Philosophy**: Read `docs/development/NO_SCOPE_CREEP.md` before making ANY changes. This project does exactly 3 things via MCP tools: invoke, check, and read.

## Installation and Usage

```bash
# Clone the repository
git clone https://github.com/christianwissmann85/delegate.git
cd delegate

# Build the project
go build -o delegate main.go

# Add to Claude Code
claude mcp add delegate -s project -- go run main.go
```

## Development Commands

```bash
# Run tests
go test ./...
go test -v -tags=e2e ./e2e/...  # E2E tests (now passing!)

# Build
go build -o delegate main.go

# Format and lint
go fmt ./...
go vet ./...
```

## Architecture

The codebase follows this structure:
```
delegate/
├── main.go               # MCP server entry point
├── go.mod
└── internal/
    ├── mcp/              # MCP protocol implementation
    │   └── server.go     # Handles MCP connections and tool routing
    ├── config/           # Configuration management
    ├── handlers/         # Implements invoke, check, read logic
    ├── providers/        # LLM provider implementations
    │   ├── interface.go  # Common provider interface
    │   ├── anthropic.go  # Claude integration
    │   └── google.go     # Gemini integration
    ├── extractor/        # Code/explanation extraction logic
    └── storage/          # File system operations
```

## Key Implementation Details

1. **MCP Tools** - Only 3 tools exist:
   - `invoke`: Generate code using specified model
   - `check`: Get metadata about generated output
   - `read`: Retrieve generated content
   - Tool descriptions for MCP registration: see `docs/mcp-tool-descriptions.md`

2. **Storage**: All outputs stored in `.delegate/` directory with format `out_YYYYMMDD_HHMMSS`

3. **Providers**: Only Gemini (Flash/Pro) and Claude (Sonnet/Opus 4) models

4. **Testing Strategy**:
   - Unit tests: `internal/*/[module]_test.go`
   - Integration tests: Same location
   - E2E tests: `e2e_test.go` with build tag
   - Test fixtures: `test/fixtures/`

5. **Performance Requirements**:
   - Check operation: <100ms
   - Read operation: <500ms
   - Streaming for invoke to prevent timeouts

## Critical Constraints

- NO new features beyond the 3 tools (see `docs/NO_SCOPE_CREEP.md`)
- NO session management, progress indicators, or complex routing
- NO web UI, CLI tools, or analytics
- Stateless - each operation is atomic
- Local filesystem only - no cloud storage
- Single prompt, single response - no conversations

## Current Implementation Status

✅ **ALL FEATURES COMPLETE AND TESTED!**

- ✅ MCP Server Foundation
  - JSON-RPC protocol handling
  - Tool registration with full schemas
  - Structured JSON logging
  - Configuration management
  
- ✅ Storage Layer
  - File-based storage implementation
  - Output ID generation with atomic counter
  - Atomic writes to prevent corruption
  - Hourly cleanup routine (deletes files >24h old)
  
- ✅ All 3 Tools Working
  - `delegate_invoke` - Delegate tasks to Gemini/Claude models
  - `delegate_check` - Get metadata without consuming tokens
  - `delegate_read` - Retrieve results (with revolutionary `write_to` feature!)
  
- ✅ Security Hardened
  - Path traversal prevention in `write_to`
  - Input validation on all parameters
  - Robust error handling
  
- ✅ Fully Tested
  - Unit tests passing
  - Integration tests passing
  - E2E tests passing (fixed MCP protocol parsing)
  - Real API tests with Gemini working

See `docs/development/implementation-roadmap-VICTORY.md` for the celebration roadmap!

## Documentation Structure

```
docs/
├── architecture/     # Technical specs and decisions
├── development/      # Roadmap, testing, philosophy
├── guides/          # User and developer guides
└── reference/       # API and model references
```

Key documents:
- Token-Efficient Workflow: `docs/guides/token-efficient-workflow.md` (MUST READ!)
- Architecture: `docs/architecture/architecture-spec.md`
- Victory Roadmap: `docs/development/implementation-roadmap-VICTORY.md`
- Philosophy: `docs/development/NO_SCOPE_CREEP.md`

## Production Usage Examples

```bash
# Iterative single-file generation workflow (RECOMMENDED)
delegate_invoke(model: "gemini-2.5-flash", prompt: "Create user.go model with GORM tags")
delegate_check(output_id)  # Check size first
delegate_read(output_id, options: {write_to: "models/user.go"})  # Output: "Content written to models/user.go (3.2 KB, ~800 tokens saved)"

# Fix compilation errors iteratively
go build models/user.go 2> errors.txt
delegate_invoke(model: "gemini-2.5-flash", files: ["models/user.go", "errors.txt"], prompt: "Fix these compilation errors")
delegate_read(output_id, options: {write_to: "models/user.go"})  # Output: "Content written to models/user.go (3.4 KB, ~850 tokens saved)"

# Build complex projects file by file
delegate_invoke(model: "gemini-2.5-flash", prompt: "Create README.md for BlogAPI project")
delegate_read(output_id, options: {write_to: "README.md"})
delegate_invoke(model: "gemini-2.5-flash", prompt: "Create auth middleware", files: ["README.md", "models/user.go"])
delegate_read(output_id, options: {write_to: "middleware/auth.go"})
```