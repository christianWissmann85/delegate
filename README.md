# Delegate

**Reviewed by Christian Wissmann for Delegate V2.0**

> Transform Claude Code from a coder into a Team Lead. Delegate heavy lifting to other LLMs. Save 90%+ tokens. No fluff. Just works.

## What is Delegate?

Delegate is an MCP server that transforms how Claude Code works. Instead of burning through context tokens on large file generation, Claude becomes a technical lead - orchestrating work while Gemini and other models do the heavy lifting. Generate entire codebases, analyze massive documents, process huge files - all while keeping Claude's context pristine for the important architectural decisions.

**Think of it this way**: Why have Claude write 1000 lines of boilerplate when it can delegate that to Gemini and focus on the architecture?

## Installation

### 🚀 Quick Install (Recommended)

We provide an installer script that handles everything and teaches you about MCP:

```bash
# Clone the repository
git clone https://github.com/christianwissmann85/delegate.git
cd delegate

# Run the installer (builds & installs system-wide)
./update.sh
```

**What update.sh does:**
- ✅ Builds the delegate binary
- ✅ Installs it to `/usr/local/bin` (system-wide)
- ✅ Shows you EXACTLY how to use MCP with any project
- ✅ Contains extensive comments explaining MCP in detail
- ✅ No more confusion about per-project setup!

### Manual Installation

```bash
# Build manually
go build -o delegate main.go

# Add to a specific project
cd /your/project
claude mcp add delegate -s project -- /path/to/delegate
```

## Using Delegate in Your Projects

After running `update.sh`, here's how to use Delegate in ANY project:

```bash
# 1. Go to your project
cd /path/to/your/awesome-project

# 2. Add Delegate to THIS project (one-time setup)
claude mcp add delegate -s project -- delegate

# 3. Start Claude Code
claude

# 4. Now Claude can delegate work!
```

**Inside Claude Code:**
```python
# The new 2-step workflow - cleaner, simpler, no confusion
response = delegate_submit_task(model="gemini-2.5-flash", prompt="Create a complete REST API server with auth")
delegate_write_output_to_file(output_id=response["output_id"], write_to="server.go")  # ZERO tokens!

# Or review before saving
response = delegate_submit_task(model="gemini-2.5-flash", prompt="Create comprehensive test suite")
metadata = delegate_get_output_metadata(output_id=response["output_id"])  # Check size/structure
content = delegate_get_output_content(output_id=response["output_id"], options={"extract": "code"})  # Read what you need
```

## Why Delegate Makes Claude a Team Lead

### Without Delegate (Claude as a Coder):
- 🔥 Burns thousands of tokens writing boilerplate
- 😓 Context fills up with generated code
- 🐌 Slower responses as context grows
- 💸 Higher costs from token usage

### With Delegate (Claude as a Team Lead):
- 🧠 Claude focuses on architecture & decisions
- 🚀 Gemini handles the heavy lifting (1M context!)
- 💾 Generated files saved directly (0 tokens!)
- 🎯 Claude's context stays clean for important work
- 💰 90%+ token savings on large generations

## Real-World Example

```python
# Traditional approach: Claude writes everything (5000+ tokens)
❌ "Write a complete user authentication system with JWT"

# Delegate approach: Claude orchestrates, Gemini codes (500 tokens)
✅ response = delegate_submit_task(
    model="gemini-2.5-flash", 
    prompt="""
    Create a complete user authentication system:
    - JWT token generation and validation
    - Password hashing with bcrypt
    - User registration/login endpoints
    - Middleware for protected routes
    - Proper error handling
    """)
✅ delegate_write_output_to_file(output_id=response["output_id"], write_to="auth/auth.go")
```

## Quick Start

1. Install Delegate with our helpful script:
   ```bash
   git clone https://github.com/christianwissmann85/delegate.git
   cd delegate
   ./update.sh  # Read the output - it explains EVERYTHING!
   ```

2. Set your API keys:
   ```bash
   export GOOGLE_API_KEY="your-key"      # For Gemini
   export ANTHROPIC_API_KEY="your-key"  # For Claude models
   ```

3. In each project where you want to use Delegate:
   ```bash
   claude mcp add delegate -s project -- delegate
   claude
   ```

[Getting Started Guide →](docs/Getting%20Started%20Guide.md)

## Key Features

- ✅ **4 Clear Tools**: Each tool does ONE thing well - no ambiguity
- ✅ **Structured JSON Responses**: Every response is parseable, predictable JSON
- ✅ **Zero-Token File Writes**: `delegate_write_output_to_file` saves directly to disk
- ✅ **Smart Metadata**: Check output size/structure before deciding what to do
- ✅ **Relative Paths**: Use simple paths like "src/main.go" - no more absolute path headaches
- ✅ **4 Powerful Models**: Gemini 2.0 Flash Experimental/1206 (1M tokens!), Claude Sonnet/Opus
- ✅ **No Complexity**: Read [NO_SCOPE_CREEP.md](docs/development/NO_SCOPE_CREEP.md)

## Documentation

📚 **[View Full Documentation](docs/README.md)**

### Quick Links
- [Getting Started](docs/guides/getting-started-guide.md) - Start here!
- [API Reference](docs/reference/api-reference.md) - Tool specifications
- [Claude Code Guide](docs/guides/claude-code-guide.md) - Usage patterns
- [Architecture](docs/architecture/architecture-spec.md) - Technical details

## Project Status

✅ **Ready for Use** - All core features implemented and tested. The new 4-tool API makes workflows clearer and more reliable than ever. The `delegate_write_output_to_file` feature lets you save massive outputs directly to disk without consuming any tokens!

## Philosophy

This project has one sacred document: [NO_SCOPE_CREEP.md](docs/development/NO_SCOPE_CREEP.md). We do four things. We do them well. That's it.

## Requirements

- Go 1.21+ (for building from source)
- Claude Code CLI
- API key for at least one provider (Gemini or Claude)
- That's it

## Contributing

Want to add a feature? Read [NO_SCOPE_CREEP.md](docs/NO_SCOPE_CREEP.md) first. The answer is probably no. 

Found a bug? That's different. Please open an issue!

## License

MIT - Because complexity bad, simplicity good.

---

Built with ❤️ and an iron-clad commitment to simplicity.

