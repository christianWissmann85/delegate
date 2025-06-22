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
echo -e "   ${GREEN}cd /path/to/your/project${NC}"
echo ""
echo -e "${YELLOW}STEP 2: Add Delegate to THIS project (one-time setup)${NC}"
echo -e "   ${GREEN}claude mcp add delegate -s project -- delegate${NC}"
echo ""
echo -e "${YELLOW}STEP 3: Start Claude Code${NC}"
echo -e "   ${GREEN}claude${NC}"
echo ""
echo -e "${YELLOW}STEP 4: Use Delegate tools in Claude!${NC}"
echo -e "   ${GREEN}delegate_invoke(model: \"gemini-2.5-flash\", prompt: \"Create a REST API\")${NC}"
echo -e "   ${GREEN}delegate_check(output_id)  # See size without reading${NC}"
echo -e "   ${GREEN}delegate_read(output_id, options: {write_to: \"api.go\"})  # Save without tokens!${NC}"
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
echo -e "• No need to copy files - just add MCP connection!"
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
# The Delegate MCP Server
# -----------------------
# Delegate is an MCP server that provides 3 tools to Claude:
# - delegate_invoke: Generate code/content using other AI models (saves tokens!)
# - delegate_check: Check output size without reading it (saves tokens!)
# - delegate_read: Read generated content (with optional write_to for saving tokens!)
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
#    cd /path/to/your/project
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
# Why Use Delegate via MCP?
# ------------------------
# - Save 90%+ tokens on large code generation
# - Use Gemini's 2M token context window
# - Generate entire files without consuming Claude's context
# - Let Claude orchestrate while Gemini does heavy lifting
#
# Example Workflow
# ---------------
# 1. cd my-awesome-project
# 2. claude mcp add delegate -s project -- delegate
# 3. claude
# 4. In Claude: delegate_invoke(model: "gemini-2.5-flash", prompt: "Create user.go")
# 5. In Claude: delegate_read(output_id, options: {write_to: "user.go"})
# 6. File created with ~1000 tokens saved!
#
# ============================================================================