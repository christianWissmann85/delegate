#!/bin/bash

# Test script to verify MCP server functionality

echo "Testing Delegate MCP Server..."

# Test 1: Initialize
echo -e "\n1. Testing initialize..."
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"test","version":"1.0"}}}' | ./delegate | jq .

# Test 2: List tools
echo -e "\n2. Testing tools/list..."
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./delegate | jq .

echo -e "\nMCP server is working! You can now use it with Claude Code:"
echo "claude mcp add delegate -s project -- $(pwd)/delegate"