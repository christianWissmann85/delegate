# **Delegate: MCP Installation Guide for Claude Code v1.0**

**Status:** Final | **Version:** 1.0 | **Date:** 2025-06-20

## **1. Overview**

Delegate is an MCP (Model Context Protocol) server that integrates seamlessly with Claude Code CLI. No binaries, no system installation - just add it as an MCP server and start delegating code generation to save your context tokens!

**What you need:**
- Claude Code CLI installed and working
- API keys for Gemini and/or Claude models
- 2 minutes to set it up

## **2. Installation**

### **Step 1: Add Delegate MCP Server to Claude Code**

Run this command in your terminal:

```bash
claude mcp add delegate -s user -- npx -y @christianwissmann85/delegate
```

This command:
- Adds Delegate as an MCP server named "delegate"
- Uses user scope (`-s user`) so it's available in all your projects
- Runs via npx to always use the latest version

### **Step 2: Configure API Keys**

You need to set environment variables for the LLM providers you want to use.

**Option A: Add to your shell profile** (recommended)

Add these to your `~/.bashrc`, `~/.zshrc`, or equivalent:

```bash
export ANTHROPIC_API_KEY="sk-ant-..."
export GOOGLE_API_KEY="AIza..."
```

Then reload your shell:
```bash
source ~/.bashrc  # or source ~/.zshrc
```

**Option B: Project-specific configuration**

Create a `.env` file in your project directory:
```
ANTHROPIC_API_KEY=sk-ant-...
GOOGLE_API_KEY=AIza...
```

### **Step 3: Get Your API Keys**

**For Gemini models:**
1. Go to [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Click "Create API Key"
3. Copy your key

**For Claude models:**
1. Go to [Anthropic Console](https://console.anthropic.com/)
2. Navigate to API Keys
3. Create a new key
4. Copy your key

### **Step 4: Verify Installation**

Start Claude Code in any project:
```bash
claude
```

Then check if Delegate is available:
```
/mcp
```

You should see:
```
⎿ MCP Server Status
⎿ • delegate: connected
```

## **3. Alternative Installation Options**

### **Project-Scoped Installation**

To install Delegate only for a specific project:

```bash
# Navigate to your project
cd /path/to/your/project

# Add with project scope
claude mcp add delegate -s project -- npx -y @christianwissmann85/delegate
```

### **Custom Environment Variables**

You can pass environment variables directly in the MCP configuration:

```bash
claude mcp add delegate -s user -- env \
  ANTHROPIC_API_KEY="your-key" \
  GOOGLE_API_KEY="your-key" \
  DELEGATE_LOG_LEVEL="debug" \
  npx -y @christianwissmann85/delegate
```

## **4. Configuration Options**

### **Environment Variables**

Set these in your shell profile or project `.env` file:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ANTHROPIC_API_KEY` | Yes* | - | For Claude models |
| `GOOGLE_API_KEY` | Yes* | - | For Gemini models |
| `DELEGATE_LOG_LEVEL` | No | `info` | Options: `debug`, `info`, `warn`, `error` |
| `DELEGATE_TIMEOUT_SECONDS` | No | `60` | Max time for generation |
| `DELEGATE_OUTPUT_DIR` | No | `.delegate` | Where outputs are stored |

*At least one API key is required

## **5. Using Delegate in Claude Code**

Once installed, Claude Code will automatically use Delegate when appropriate. You can also explicitly request it:

```
"Use Delegate to generate a complex React component with Gemini"
```

The three tools available:
- `delegate_invoke` - Generate code with another LLM
- `delegate_check` - Check output size before reading
- `delegate_read` - Read generated content

## **6. Supported Models**

| Model ID | Provider | Best For |
|----------|----------|----------|
| `gemini-2.5-flash` | Google | Fast, everyday code generation |
| `gemini-2.5-pro` | Google | Complex reasoning and architecture |
| `claude-sonnet-4-20250514` | Anthropic | Balanced quality and speed |
| `claude-opus-4-20250514` | Anthropic | Highest quality, critical code |

## **7. Managing Your MCP Server**

### **View All MCP Servers**
```bash
claude mcp list
```

### **Remove Delegate**
```bash
claude mcp remove delegate
```

### **Disable Temporarily**
```bash
claude mcp disable delegate
```

### **Re-enable**
```bash
claude mcp enable delegate
```

## **8. Troubleshooting**

### **"MCP server not connected"**
1. Check your API keys are set correctly:
   ```bash
   echo $ANTHROPIC_API_KEY
   echo $GOOGLE_API_KEY
   ```
2. Try removing and re-adding:
   ```bash
   claude mcp remove delegate
   claude mcp add delegate -s user -- npx -y @christianwissmann85/delegate
   ```

### **"API key errors"**
- Ensure no quotes around your API keys in environment variables
- Verify keys are valid in their respective consoles
- Check you're using the correct key format

### **"Timeout errors"**
- Large generations might exceed 30 seconds
- Set a higher timeout:
  ```bash
  export DELEGATE_TIMEOUT_SECONDS=60
  ```

### **"Permission denied" on outputs**
- Delegate creates a `.delegate` folder in your current directory
- Ensure you have write permissions
- Or set a custom output directory:
  ```bash
  export DELEGATE_OUTPUT_DIR=/tmp/delegate
  ```

## **9. Debug Mode**

For troubleshooting, enable debug logging:

```bash
export DELEGATE_LOG_LEVEL=debug
claude mcp add delegate -s user -- npx -y @christianwissmann85/delegate
```

Then restart Claude Code and check the logs:
```bash
# Check MCP server logs
ls ~/.claude/logs/
tail -f ~/.claude/logs/mcp-server-delegate.log
```

## **10. Updating Delegate**

Delegate updates automatically via npx. To force an update:

1. Clear npm cache:
   ```bash
   npm cache clean --force
   ```

2. Remove and re-add the server:
   ```bash
   claude mcp remove delegate
   claude mcp add delegate -s user -- npx -y @christianwissmann85/delegate
   ```

## **11. Quick Start Example**

```bash
# 1. Set your API keys
export GOOGLE_API_KEY="your-google-key"
export ANTHROPIC_API_KEY="your-anthropic-key"

# 2. Add Delegate
claude mcp add delegate -s user -- npx -y @christianwissmann85/delegate

# 3. Start Claude Code
claude

# 4. Use it!
# Type: "Use Delegate with Gemini to create a REST API server"
```

## **12. Best Practices**

1. **Global vs Project Scope**: Use `-s user` for global access, `-s project` for project-specific needs
2. **API Key Security**: Never commit API keys to version control
3. **Output Cleanup**: Delegate automatically cleans up files older than 24 hours
4. **Model Selection**: Let Claude Code choose the model, or specify explicitly in your request

---

**Next Steps:** Check out `claude-code-guide-updated.md` for usage patterns and examples!