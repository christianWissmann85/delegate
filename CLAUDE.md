# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Delegate is an MCP (Model Context Protocol) server that allows Claude Code to delegate code generation tasks to other LLMs (Gemini and Claude models) to save context tokens. The project is currently in the documentation phase, with implementation starting per the 3-week roadmap.

**Core Philosophy**: Read `docs/NO_SCOPE_CREEP.md` before making ANY changes. This project does exactly 3 things via MCP tools: invoke, check, and read.

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

## Current Status

Project is in documentation phase. Implementation follows the roadmap in `docs/implementation-roadmap.md`. All technical specifications are in `docs/Delegate Architecture & Technical Specification.md`.

When implementing, follow the daily tasks in the roadmap strictly. The goal is a production-ready MCP server available via `npx @christianwissmann85/delegate` in 3 weeks.