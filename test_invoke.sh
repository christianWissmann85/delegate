#!/bin/bash
# Test script for the invoke tool using mock provider

# Set environment variables for testing
export GOOGLE_API_KEY="test-key"
export DELEGATE_LOG_LEVEL="debug"

# Start the delegate server in the background
echo "Starting delegate server..."
./delegate &
SERVER_PID=$!

# Give server time to start
sleep 2

# Test MCP initialize
echo "Testing MCP connection..."
echo '{"jsonrpc": "2.0", "method": "initialize", "params": {"protocolVersion": "0.1.0", "capabilities": {"roots": {"listChanged": true}, "sampling": {}}}, "id": 1}' | nc -N localhost 3000

# Kill the server
kill $SERVER_PID