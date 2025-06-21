# Delegate

> The industrial-strength MCP server for delegating LLM tasks. Save Claude Code's context tokens. No fluff. Just works.

## What is Delegate?

Delegate lets Claude Code save context tokens by delegating heavy tasks to other LLMs (Gemini & Claude). Generate code, analyze documents, process large files - anything that would eat up Claude Code's context. Three tools. Zero complexity.

## Installation

```bash
claude mcp add delegate -s user -- npx -y @christianwissmann85/delegate
```

That's it. [Full installation guide →](docs/Delegate%20MCP%20Installation%20Guide%20v1.0.md)

## Key Use Cases

### 🚀 Code Generation
```
Use Delegate to generate a complete authentication system with Gemini
```

### 📚 Document Analysis
```
Use Delegate to analyze these 5 architecture documents and find all API patterns
```

### 🔍 Large File Processing
```
Use Delegate to review this 10k line codebase and identify security issues
```

## Quick Start

1. Set your API keys:
   ```bash
   export GOOGLE_API_KEY="your-key"
   export ANTHROPIC_API_KEY="your-key"
   ```

2. Start Claude Code and delegate heavy tasks:
   ```
   Use Delegate to analyze all documentation files and summarize the testing strategy
   ```

[Getting Started Guide →](docs/Getting%20Started%20Guide.md)

## Features

- ✅ **3 Simple Tools**: invoke, check, read
- ✅ **4 Powerful Models**: Gemini 2.5 Flash/Pro (1M tokens!), Claude Sonnet/Opus 4
- ✅ **Token Efficient**: Delegate document analysis, code generation, any heavy lifting
- ✅ **Context Preservation**: Keep Claude Code's context clean for actual work
- ✅ **No Complexity**: Read [NO_SCOPE_CREEP.md](docs/NO_SCOPE_CREEP.md)

## Documentation

- [Getting Started](docs/Getting%20Started%20Guide.md) - Start here!
- [API Reference](docs/Delegate%20API%20Reference%20v1.0.md) - Tool specifications
- [Model Reference](docs/Model%20Reference%20Card%20v1.0.md) - Which model when?
- [Architecture](docs/Delegate%20Architecture%20&%20Technical%20Specification%20v1.0.md) - How it works
- [Claude Code Guide](docs/claude-code-guide-updated.md) - Usage patterns

## Project Status

🚧 **In Development** - Following the [Implementation Roadmap](docs/Implementation%20Roadmap%20-%20Delegate%20v1.0.md)

## Philosophy

This project has one sacred document: [NO_SCOPE_CREEP.md](docs/NO_SCOPE_CREEP.md). We do three things. We do them well. That's it.

## Requirements

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