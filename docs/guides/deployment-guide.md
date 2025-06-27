# Delegate MCP Server - Deployment Guide

**Reviewed by Christian Wissmann for Delegate V2.0**

## 1. Prerequisites

Before installing Delegate, ensure you have the following:

### Go Installation
- **Go 1.21 or higher** is required for building from source
- Check your Go version: `go version`
- Install Go from [golang.org](https://golang.org/download/) if needed

### Claude Code CLI
- Install the Claude Code CLI following the [official installation guide](https://docs.anthropic.com/claude/docs/claude-code)
- Verify installation: `claude --version`

### API Keys
You need at least one API key from the supported providers:

- **Google AI Studio API Key** (for Gemini models)
  - Get your key at [ai.google.dev](https://ai.google.dev/)
  - Set environment variable: `export GOOGLE_API_KEY="your-key-here"`

- **Anthropic API Key** (for Claude models)
  - Get your key at [console.anthropic.com](https://console.anthropic.com/)
  - Set environment variable: `export ANTHROPIC_API_KEY="your-key-here"`

## 2. Installation Methods

### Option A: Quick Install Script (Recommended)

The fastest way to get started:

```bash
# Clone the repository
git clone https://github.com/christianwissmann85/delegate.git
cd delegate

# Run the installer (builds & installs system-wide)
./update.sh
```

**What the installer does:**
- ✅ Builds the delegate binary
- ✅ Installs it to `/usr/local/bin` (system-wide access)
- ✅ Shows you exactly how to use MCP with any project
- ✅ Contains extensive comments explaining MCP setup
- ✅ No more confusion about per-project setup!

### Option B: Manual Installation

If you prefer to build and install manually:

```bash
# Clone the repository
git clone https://github.com/christianwissmann85/delegate.git
cd delegate

# Build the binary
go build -o delegate main.go

# Install system-wide (optional)
sudo mv delegate /usr/local/bin/

# Or install to a specific location
cp delegate /path/to/your/binaries/
```

### Option C: Binary Download

Download pre-built binaries from the [releases page](https://github.com/christianwissmann85/delegate/releases):

```bash
# Download for your platform
curl -L https://github.com/christianwissmann85/delegate/releases/latest/download/delegate-linux-amd64 -o delegate

# Make executable
chmod +x delegate

# Install system-wide
sudo mv delegate /usr/local/bin/
```

## 3. Configuration

### Environment Variables

Create a `.env` file in your project directory or set environment variables:

```bash
# Required: At least one API key
export GOOGLE_API_KEY="your-google-api-key"
export ANTHROPIC_API_KEY="your-anthropic-api-key"

# Optional: Server configuration
export DELEGATE_PORT="8080"                    # Default port
export DELEGATE_TIMEOUT="120"                  # Default timeout in seconds
export DELEGATE_MAX_TOKENS="100000"           # Default max tokens
export DELEGATE_WORKING_DIR="/path/to/project" # Default working directory
```

### Configuration File

Alternatively, create a `config.json` file:

```json
{
  "port": 8080,
  "timeout": 120,
  "max_tokens": 100000,
  "working_directory": "/path/to/your/project",
  "providers": {
    "google": {
      "api_key": "your-google-api-key"
    },
    "anthropic": {
      "api_key": "your-anthropic-api-key"
    }
  }
}
```

## 4. Running the Server

### Standalone Mode

For testing or development:

```bash
# Run with environment variables
export GOOGLE_API_KEY="your-key"
delegate

# Or run with explicit configuration
delegate --config config.json --port 8080
```

### With Claude Code (Recommended)

Add Delegate to your project and run with Claude Code:

```bash
# Navigate to your project
cd /path/to/your/project

# Add Delegate to this project (one-time setup)
claude mcp add delegate -s project -- delegate

# Start Claude Code
claude
```

### System Service (Production)

For production deployments, create a systemd service:

```ini
# /etc/systemd/system/delegate.service
[Unit]
Description=Delegate MCP Server
After=network.target

[Service]
Type=simple
User=delegate
WorkingDirectory=/opt/delegate
ExecStart=/usr/local/bin/delegate
Environment=GOOGLE_API_KEY=your-key
Environment=ANTHROPIC_API_KEY=your-key
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start the service
sudo systemctl enable delegate
sudo systemctl start delegate
sudo systemctl status delegate
```

## 5. Using in Claude Code - The New 4-Tool API

Once Delegate is running with Claude Code, you have access to four powerful tools:

### Tool Overview

| Tool | Purpose | Token Cost |
|------|---------|------------|
| `delegate_submit_task` | Submit generation task, get output_id | Low (~50-100 tokens) |
| `delegate_get_output_metadata` | Check output size/structure | Low (~20 tokens) |
| `delegate_get_output_content` | Read content into context | High (proportional to content) |
| `delegate_write_output_to_file` | Write content directly to file | **ZERO tokens** |

## 6. Troubleshooting

### Common Issues

#### "delegate: command not found"

```bash
# Check if delegate is in PATH
which delegate

# If not found, install system-wide
sudo cp delegate /usr/local/bin/

# Or add to PATH
export PATH=$PATH:/path/to/delegate/directory
```

#### "API key not found" errors

```bash
# Verify environment variables are set
echo $GOOGLE_API_KEY
echo $ANTHROPIC_API_KEY

# Set them if missing
export GOOGLE_API_KEY="your-key-here"
export ANTHROPIC_API_KEY="your-key-here"

# Or create .env file in project directory
cat > .env << EOF
GOOGLE_API_KEY=your-key-here
ANTHROPIC_API_KEY=your-key-here
EOF
```

#### "Connection refused" or server not starting

```bash
# Check if port is already in use
lsof -i :8080

# Use a different port
delegate --port 8081

# Or set environment variable
export DELEGATE_PORT=8081
```

#### "File write failed" errors
```bash
# Check permissions in target directory
ls -la /path/to/target/directory

# Create directory if it doesn't exist
mkdir -p src/generated

# Check working directory
pwd
# Delegate uses relative paths from its working directory
```

#### Model timeout errors

```bash
# Increase timeout for large generations
export DELEGATE_TIMEOUT=300  # 5 minutes

# Or specify in the request
delegate_submit_task(
    model="gemini-2.0-flash-experimental",
    prompt="your prompt",
    timeout=300
)
```

### Debug Mode

Enable debug logging for troubleshooting:

```bash
# Run with debug logging
export DELEGATE_DEBUG=true
delegate

# Or with specific log level
export DELEGATE_LOG_LEVEL=debug
delegate
```

### Health Check

Verify Delegate is working correctly:

```bash
# Check server status
curl http://localhost:8080/health

# Test basic functionality
curl -X POST http://localhost:8080/delegate_submit_task \
  -H "Content-Type: application/json" \
  -d '{"model": "gemini-2.5-flash", "prompt": "Hello, world!"}'
```

### Performance Optimization

For better performance:

```bash
# Increase memory limit for Go
export GOGC=100
export GOMEMLIMIT=4GiB

# Use faster models for development
# gemini-2.5-flash instead of gemini-2.5-pro

# Reduce max_tokens for faster responses
delegate_submit_task(
    max_tokens=50000,  # Instead of default 100000
    # ... other parameters
)
```

### Getting Help

If you're still having issues:

1. **Check the logs**: Look for error messages in the console output
2. **Verify configuration**: Double-check API keys and environment variables
3. **Test connectivity**: Ensure you can reach the AI provider APIs
4. **Check file permissions**: Verify write access to target directories
5. **Review working directory**: Ensure relative paths resolve correctly

For additional support, check the [GitHub Issues](https://github.com/christianwissmann85/delegate/issues) or create a new issue with:
- Your operating system
- Go version (`go version`)
- Claude Code version (`claude --version`)
- Error messages and logs
- Steps to reproduce the issue

---

**Ready to transform Claude Code into a Team Lead?** Start with the quick install script and begin delegating heavy lifting to other LLMs while keeping Claude's context clean for the important architectural decisions.