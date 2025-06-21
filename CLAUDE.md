# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Delegate is an MCP (Model Context Protocol) server that allows Claude Code to delegate heavy tasks (code generation, document analysis, large file processing) to other LLMs (Gemini and Claude models) to save context tokens. 

**Current Status**: Day 2 of 21 - MCP server foundation complete, starting storage layer.

**Core Philosophy**: Read `docs/development/NO_SCOPE_CREEP.md` before making ANY changes. This project does exactly 3 things via MCP tools: invoke, check, and read.

## Development Commands

```bash
# Initial setup (when starting implementation)
go mod init github.com/christianwissmann85/delegate

# Run tests
go test ./...
go test -v --tags=e2e .  # E2E tests

# Build
go build -o delegate main.go

# Test with Claude Code during development
claude mcp add delegate-dev -s project -- go run main.go

# Format and lint
go fmt ./...
go vet ./...

# NPM packaging (Week 3)
npm init -y
npm publish
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

- ✅ Day 1-2: MCP Server Foundation (COMPLETE)
  - JSON-RPC protocol handling
  - Tool registration with full schemas
  - Structured JSON logging
  - Configuration management
  
- 🚧 Day 3-4: Storage Layer (NEXT)
  - File-based storage implementation
  - Output ID generation
  - Atomic writes
  - Cleanup routine

See `docs/development/implementation-roadmap.md` for full schedule.

## Documentation Structure

```
docs/
├── architecture/     # Technical specs and decisions
├── development/      # Roadmap, testing, philosophy
├── guides/          # User and developer guides
└── reference/       # API and model references
```

Key documents:
- Architecture: `docs/architecture/architecture-spec.md`
- Current roadmap: `docs/development/implementation-roadmap.md`
- Philosophy: `docs/development/NO_SCOPE_CREEP.md`