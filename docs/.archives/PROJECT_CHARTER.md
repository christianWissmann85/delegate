# Delegate - AI Task Delegation via MCP

## Mission Statement
Enable Claude Code to delegate heavy tasks (code generation, document analysis, large file processing) to other LLMs without consuming context tokens, through a dead-simple, bulletproof MCP interface.

## Core Principles
1. **Simplicity over features** - If a feature adds complexity, we don't need it
2. **Reliability over flexibility** - Better to do 3 things perfectly than 10 things poorly
3. **Token efficiency is paramount** - Every operation designed to minimize Claude's token usage
4. **No scope creep allowed** - See NO_SCOPE_CREEP.md and stick to it

## Problem We're Solving
- Claude Code needs to handle large tasks without exhausting context tokens
- Reading multiple 4-5k line documents consumes context rapidly
- Analyzing entire codebases is impossible with limited tokens
- Current solutions (AAG) are overly complex and unreliable
- Timeouts, truncations, and parsing errors waste developer time
- We need industrial-strength delegation, not fancy features

## Success Metrics
- **Zero timeout failures** after first week of production use
- **<2s response time** for check operations
- **100% type safety** - Go compiler catches all errors before runtime
- **<500 lines of code** for core functionality
- **Single binary deployment** under 10MB

## Target User
Claude Code (and other AI assistants) that need to:
- Generate code without consuming their own tokens
- Check output sizes before reading
- Extract just the code without explanations

## Non-Goals
- Human-friendly CLI interface
- Complex orchestration
- Feature parity with AAG
- Supporting every LLM provider

## Timeline
- Week 1: Core implementation
- Week 2: Hardening and reliability
- Week 3: Polish and deployment
- Week 4: First production use with Claude Code

## Ownership
- **Architect**: Claude Code
- **Developer**: Chris
- **Primary User**: Claude Code
- **Maintenance**: Minimal - design for zero maintenance

## Decision Record
- Language: Go (for reliability and single binary)
- Protocol: MCP only (no CLI)
- Storage: Local filesystem (no databases)
- Models: Big 2 only (Gemini, Claude)