#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔨 Building Delegate...${NC}"
if go build -o delegate main.go; then
    echo -e "${GREEN}✅ Build successful!${NC}"
else
    echo -e "${YELLOW}❌ Build failed! Fix errors and try again.${NC}"
    exit 1
fi

echo -e "${BLUE}📦 Installing to /usr/local/bin...${NC}"
if sudo cp delegate /usr/local/bin/; then
    echo -e "${GREEN}✅ Delegate installed system-wide!${NC}"
else
    echo -e "${YELLOW}❌ Installation failed! Check permissions.${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 SUCCESS! Delegate is installed system-wide!${NC}"
echo ""
echo -e "${BLUE}=== WHAT IS DELEGATE? ===${NC}"
echo -e "Delegate is an MCP (Model Context Protocol) server that lets Claude Code"
echo -e "use OTHER AI models (Gemini/Claude) to generate code without consuming YOUR tokens!"
echo -e "Save 90%+ tokens on large file generation! 🚀"
echo ""
echo -e "${BLUE}=== WHAT IS MCP? ===${NC}"
echo -e "MCP lets Claude Code connect to external tools. Think of it as plugins for Claude!"
echo -e "Each project needs its own MCP setup (it's NOT global)."
echo ""
echo -e "${BLUE}=== HOW TO USE DELEGATE IN A PROJECT ===${NC}"
echo ""
echo -e "${YELLOW}STEP 1: Navigate to your project${NC}"
echo -e "   ${GREEN}cd your-project-directory${NC}"
echo ""
echo -e "${YELLOW}STEP 2: Add Delegate to THIS project (one-time setup)${NC}"
echo -e "   ${GREEN}claude mcp add delegate -s project -- delegate${NC}"
echo ""
echo -e "${YELLOW}STEP 3: Start Claude Code${NC}"
echo -e "   ${GREEN}claude${NC}"
echo ""
echo -e "${YELLOW}STEP 4: Use Delegate's 4 powerful tools!${NC}"
echo -e "   ${GREEN}1. submit_task(model: \"gemini-2.5-flash\", prompt: \"Create a REST API\")${NC}"
echo -e "   ${GREEN}2. get_output_metadata(output_id)  # Check size/blocks without tokens${NC}"
echo -e "   ${GREEN}3. write_output_to_file(output_id, write_to: \"api.go\")  # ZERO tokens!${NC}"
echo -e "   ${GREEN}4. get_output_content(output_id)  # Or read into context${NC}"
echo ""
echo -e "${BLUE}=== QUICK REFERENCE ===${NC}"
echo -e "List MCP servers:     ${GREEN}claude mcp list${NC}"
echo -e "Remove MCP server:    ${GREEN}claude mcp remove delegate${NC}"
echo -e "Check if installed:   ${GREEN}which delegate${NC}"
echo ""
echo -e "${BLUE}💡 REMEMBER:${NC}"
echo -e "• MCP setup is PER-PROJECT (not global)"
echo -e "• Run 'claude mcp add' in EACH new project"
echo -e "• Delegate binary is already installed system-wide"
echo -e "• Use relative paths like 'src/main.go', not absolute paths!"
echo ""
echo -e "${GREEN}🚀 Ready to save tokens! Check the comments in this script for MORE details! 🚀${NC}"

# ============================================================================
# COMPREHENSIVE MCP (Model Context Protocol) EXPLANATION
# ============================================================================
#
# What is MCP?
# ------------
# MCP (Model Context Protocol) is a protocol that allows Claude Code to communicate
# with external tools and services. Think of it as a bridge that lets Claude Code
# use special tools that aren't built into Claude itself. It's like giving Claude
# superpowers by connecting it to external programs!
#
# How does Claude Code use MCP?
# -----------------------------
# 1. Claude Code acts as an MCP CLIENT that can connect to MCP SERVERS
# 2. MCP servers (like Delegate) provide TOOLS that Claude can call
# 3. When you type "claude" in a directory, it looks for MCP configurations
# 4. Claude Code then connects to configured MCP servers and can use their tools
#
# The New Delegate MCP Server (4-Tool API v2.0)
# ---------------------------------------------
# Delegate now provides 4 clear, single-purpose tools to Claude:
#
# 1. submit_task: Generate code/content using other AI models
#    - Returns only an output_id (~50 tokens)
#    - Asynchronous operation, content stored temporarily
#
# 2. get_output_metadata: Check output details without reading content
#    - Returns size, line count, token estimate, block analysis
#    - Perfect for deciding next steps (~20-50 tokens)
#
# 3. write_output_to_file: Write content directly to disk
#    - ZERO token cost! Content never enters Claude's context
#    - Accepts relative paths like "src/component.jsx"
#
# 4. get_output_content: Read content into Claude's context
#    - High token cost (proportional to content size)
#    - Use when you need to review/modify the generated content
#
# How MCP Connections Work
# ------------------------
# IMPORTANT: MCP connections are PER-DIRECTORY, not global!
# - Each project directory needs its own MCP configuration
# - The configuration is stored in the project's root directory
# - When you run "claude" it only sees MCP servers configured for THAT directory
#
# Setting Up MCP for a New Project - Step by Step
# -----------------------------------------------
# 1. Navigate to your project directory:
#    cd your-awesome-project
#
# 2. Add Delegate as an MCP server for this project:
#    claude mcp add delegate -s project -- delegate
#
#    Breaking this down:
#    - "claude mcp add": Command to add an MCP server
#    - "delegate": The name you're giving this connection (can be anything)
#    - "-s project": Scope is "project" (local to this directory)
#    - "--": Separator between options and the command
#    - "delegate": The actual command to run (our installed binary)
#
# 3. Start Claude Code in that directory:
#    claude
#
# 4. Now Claude can use delegate tools IN THIS PROJECT!
#
# Do you need to copy the delegate executable?
# --------------------------------------------
# NO! After running this update.sh script, delegate is installed system-wide
# in /usr/local/bin. This means:
# - You can run "delegate" from ANY directory
# - No need to copy files around
# - All projects share the same delegate binary
# - Just run "claude mcp add" in each new project
#
# Common MCP Commands
# ------------------
# List MCP servers in current directory:
#   claude mcp list
#
# Remove an MCP server:
#   claude mcp remove delegate
#
# Add with custom name:
#   claude mcp add my-delegate -s project -- delegate
#
# Project Structure After MCP Setup
# ---------------------------------
# your-project/
# ├── .claude/              # Created by Claude Code
# │   └── mcp-settings.json # MCP configuration for this project
# ├── .delegate/            # Created by Delegate when used
# │   └── outputs/          # Where Delegate stores generated content
# └── your-code-files...
#
# New 4-Tool Workflow Examples
# ----------------------------
#
# WORKFLOW A: Zero-Token File Generation (Most Common)
# ---------------------------------------------------
# Goal: Generate a React component and save it directly to a file
#
# 1. submit_task({
#      model: "gemini-2.5-flash",
#      prompt: "Create a user profile component with hooks",
#      files: ["src/types.ts"]  // Relative paths for context
#    })
#    → Returns: { output_id: "gen-abc-123", working_directory: "/home/user/project" }
#
# 2. write_output_to_file({
#      output_id: "gen-abc-123",
#      write_to: "src/components/UserProfile.tsx",  // Relative path
#      options: { extract: "code" }
#    })
#    → File created with ZERO tokens consumed!
#    → Returns: { success: true, path: "src/components/UserProfile.tsx", bytes_written: 4182 }
#
# WORKFLOW B: Smart Multi-Block Handling
# --------------------------------------
# Goal: Generate code with documentation, decide what to save
#
# 1. submit_task({
#      model: "gemini-2.5-flash",
#      prompt: "Create an API with full documentation"
#    })
#    → Returns: { output_id: "gen-xyz-456" }
#
# 2. get_output_metadata({ output_id: "gen-xyz-456" })
#    → Returns structured JSON:
#    {
#      "metadata": { "size_kb": 25.3, "token_estimate": 6320 },
#      "content_analysis": {
#        "blocks_found": 2,
#        "blocks": [
#          { "index": 0, "language": "go", "size_kb": 18.7, "preview": "package main" },
#          { "index": 1, "language": "md", "size_kb": 6.6, "preview": "# API Documentation" }
#        ]
#      }
#    }
#
# 3. Based on structured data (no string parsing!), save each block:
#    write_output_to_file({
#      output_id: "gen-xyz-456",
#      write_to: "api/server.go",
#      options: { block_index: 0 }
#    })
#    write_output_to_file({
#      output_id: "gen-xyz-456", 
#      write_to: "docs/api.md",
#      options: { block_index: 1 }
#    })
#
# WORKFLOW C: Review Before Saving
# --------------------------------
# Goal: Check generated content before deciding where to save it
#
# 1. submit_task(...)  → { output_id: "gen-review-789" }
#
# 2. get_output_content({
#      output_id: "gen-review-789",
#      options: { max_tokens: 500 }  // Just a preview
#    })
#    → Returns: { content: "preview...", metadata: { tokens_returned: 500, is_truncated: true } }
#
# 3. After review, save the full content:
#    write_output_to_file({
#      output_id: "gen-review-789",
#      write_to: "final/output.go"
#    })
#
# Path Handling: Always Use Relative Paths
# ----------------------------------------
# ✅ CORRECT:   write_to: "src/main.go"
# ✅ CORRECT:   write_to: "docs/readme.md"  
# ✅ CORRECT:   files: ["config/settings.json", "src/types.ts"]
#
# ❌ WRONG:     write_to: "/home/user/project/src/main.go"
# ❌ WRONG:     files: ["/absolute/path/to/file.go"]
#
# Benefits:
# - Shorter, cleaner commands
# - Works across different environments
# - Reduces token usage
# - Less error-prone
#
# Why 4 Tools Instead of 3?
# -------------------------
# The old API had an ambiguous delegate_read that changed behavior based on parameters.
# Now each tool has ONE clear purpose:
#
# OLD (Confusing):
# - delegate_read(id, {write_to: "file.go"})     # Writes file, returns confirmation
# - delegate_read(id)                            # Returns content, costs tokens
#
# NEW (Clear):
# - write_output_to_file(id, write_to: "file.go")  # Always writes file
# - get_output_content(id)                          # Always returns content
#
# All responses are now structured JSON instead of strings, making them
# easier for AI agents to parse and act upon programmatically.
#
# Troubleshooting MCP
# ------------------
# 1. "Command not found: delegate"
#    → Run this update.sh script first!
#
# 2. "MCP server not found"
#    → Make sure you're in the right directory
#    → Run "claude mcp list" to check
#
# 3. Tools not showing in Claude
#    → Restart Claude Code after adding MCP
#    → Check "claude mcp list" shows delegate
#
# 4. "Path traversal" or "Invalid path" errors
#    → Use relative paths only: "src/file.go" not "/abs/path/file.go"
#    → Don't use ".." to go above project directory
#
# Why Use Delegate via MCP?
# ------------------------
# - Save 90%+ tokens on large code generation
# - Use Gemini's 2M token context window
# - Generate entire files without consuming Claude's context
# - Let Claude orchestrate while other models do heavy lifting
# - Clear, predictable API with structured responses
# - Zero-token file writing for maximum efficiency
#
# Token Cost Breakdown (New API)
# ------------------------------
# submit_task:           ~50-100 tokens (just the response)
# get_output_metadata:   ~20-50 tokens (structured data only)
# write_output_to_file:  ZERO tokens (content never enters context)
# get_output_content:    HIGH (proportional to content size)
#
# Compare to old workflow where delegate_read always consumed tokens!
#
# ============================================================================
