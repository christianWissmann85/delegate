# **Implementation Roadmap - Delegate v1.0**

**Status:** Final | **Version:** 1.0 | **Date:** 2025-06-20

## **3-Week Sprint to Production**

### **Week 1: Core MCP Implementation (Days 1-7)**

#### **Day 1-2: Project Setup & MCP Server Foundation**
- [x] Initialize Go module: `go mod init github.com/christianwissmann85/delegate`
- [x] Set up MCP server framework that can handle tool calls
- [x] Create directory structure as specified in architecture doc
- [x] Implement basic MCP protocol handling (connect, initialize, tool registration)
- [x] Add structured logging to stderr (JSON format)

#### **Day 3-4: Storage Layer**
- [x] Implement `Storage` interface for file operations
- [x] Output ID generation (timestamp-based: `out_YYYYMMDD_HHMMSS`)
- [x] Atomic file writing to prevent corruption
- [x] 24-hour cleanup goroutine for old outputs
- [x] Unit tests for all storage operations

#### **Day 5-6: First Provider Integration (Gemini)**
- [x] Define `Provider` interface with `GenerateStream` method
- [x] Implement Gemini provider using official Google SDK
- [x] Add streaming response handling (write to temp file during stream)
- [x] Implement timeout handling (60s default, configurable per request)
- [x] Create mock provider for testing

#### **Day 7: Wire First Tool (invoke)**
- [x] Implement `invoke` tool handler
- [x] Connect to MCP server tool registry
- [x] Basic code extraction using regex
- [x] Manual test with Claude Code: invoke → file created
- [x] Verify streaming prevents timeouts on long generations

**Week 1 Deliverable**: Can invoke Gemini from Claude Code and save outputs

### **Week 2: Complete Implementation (Days 8-14)**

#### **Day 8-9: Remaining Providers**
- [x] Implement Anthropic provider (Claude models)
- [x] Normalize error handling across providers
- [x] Add retry logic (3 attempts, exponential backoff)
- [x] Provider selection based on model parameter
- [x] Integration tests with mock providers

#### **Day 10-11: Robust Code Extraction**
- [x] Improve extraction to handle multiple code blocks
- [x] Language detection for each code block
- [x] Separate code from explanation text
- [x] Handle edge cases (no code, malformed blocks)
- [x] Implement code_only mode to return just code without explanations
- [x] Add language_hint parameter for better extraction accuracy
- [x] Unit tests for extractor module

#### **Day 12-13: Check & Read Tools**
- [x] Implement `check` tool (fast metadata inspection)
- [x] Implement `read` tool with extraction options
- [x] Token estimation and counting
- [x] Truncation logic for `max_tokens` parameter
- [x] Test full workflow: invoke → check → read

#### **Day 14: Error Handling & Hardening**
- [ ] Implement structured error responses (DelegateError type)
- [ ] Map provider errors to normalized error types
- [ ] Add retry_after and alternative_models to error responses
- [ ] Input validation for all tool parameters
- [ ] Path traversal prevention
- [ ] Memory limits for file operations
- [ ] Load testing with concurrent tool calls

**Week 2 Deliverable**: All 3 tools working reliably with both providers

### **Week 3: Polish & Production Ready (Days 15-21)**

#### **Day 15-16: MCP Package & Distribution**
- [ ] Create npm package wrapper for easy `npx` execution
- [ ] Package Go binary for multiple platforms (darwin-arm64, darwin-amd64, linux-amd64)
- [ ] Add package.json with proper bin configuration
- [ ] Test `npx @christianwissmann85/delegate` execution
- [ ] Publish to npm registry

#### **Day 17-18: Claude Code Integration Testing**
- [ ] Full integration test with Claude Code CLI
- [ ] Test all supported models
- [ ] Performance profiling and optimization
- [ ] Fix any integration issues
- [ ] Create troubleshooting guide

#### **Day 19-20: Documentation & Examples**
- [ ] Finalize all documentation
- [ ] Create example workflows in claude-code-guide-updated.md
- [ ] Add debugging tips
- [ ] Record demo video/GIF
- [ ] Final documentation review

#### **Day 21: Launch**
- [ ] Tag v1.0.0 release
- [ ] Build and upload release binaries
- [ ] Publish npm package
- [ ] Create GitHub release with changelog
- [ ] Announce to Chris!

**Week 3 Deliverable**: Production-ready MCP server available via npx

## **Development Guidelines**

### **Code Quality Standards**
- Every public function has a doc comment
- Every error includes context: `fmt.Errorf("invoke failed: %w", err)`
- No function longer than 50 lines
- No file longer than 300 lines
- Test coverage >80% for core logic

### **Testing Strategy**
- Unit tests for extractor and config
- Integration tests with mock providers
- E2E tests with mock MCP client
- Real API tests only in CI
- No flaky tests allowed

### **Daily Practices**
- Commit at end of each day
- Run all tests before commits
- Update this roadmap with ✅ when complete
- Keep a DECISIONS.md log
- Refer to NO_SCOPE_CREEP.md daily

## **Risk Mitigation**

### **Technical Risks**

1. **MCP Protocol Complexity**
   - Mitigation: Start simple, only implement required methods
   - Use existing MCP libraries, don't reinvent

2. **Provider API Changes**
   - Mitigation: Version lock SDKs
   - Test with real APIs daily in development

3. **Streaming Timeouts**
   - Mitigation: Implement streaming early (Day 6)
   - Test with large generation tasks

### **Schedule Risks**

1. **Scope Creep**
   - Mitigation: NO_SCOPE_CREEP.md is the bible
   - Three tools only. No exceptions.

2. **Integration Issues**
   - Mitigation: Test with Claude Code from Day 7
   - Keep Chris in the loop for early feedback

## **Success Criteria**

### **Week 1**
- [ ] Basic invoke working with Gemini
- [ ] Files saved and retrievable
- [ ] No panics or crashes
- [ ] Works with Claude Code

### **Week 2**  
- [ ] All 3 tools working
- [ ] Both providers integrated
- [ ] Code extraction >90% accurate
- [ ] <2s response time for check/read

### **Week 3**
- [ ] Available via `npx`
- [ ] Claude Code using it successfully
- [ ] Zero maintenance required
- [ ] Documentation complete

## **Post-Launch (v1.1+)**

Only consider after 1 month of stable operation:
- Batch operations (invoke multiple prompts)
- Additional providers (only if Chris asks)
- Caching layer (only if performance demands)

Remember: **Ship v1.0 first, enhance later!**

## **Quick Command Reference**

```bash
# Development
go mod init github.com/christianwissmann85/delegate
go test ./...
go test -v --tags=e2e .

# Building
go build -o delegate main.go

# Testing with Claude Code
claude mcp add delegate-dev -s project -- go run main.go

# NPM Publishing
npm init -y
npm publish

# Usage
npx @christianwissmann85/delegate
```

---

**The Mantra:** Three tools. Pure MCP. No scope creep. Ship it! 🚀