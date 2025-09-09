#!/bin/bash

# Test script for the get_repository_discussions function

# Ensure the script exits on any error
set -e

# Run the command and capture the output
echo '{"jsonrpc":"2.0","id":5,"params":{"name":"get_repository_discussions", "arguments":{"owner":"github", "repo":"engineering"}},"method":"tools/call"}' | go run ./cmd/github-mcp-server/main.go stdio | jq .

# Print a message indicating the test is complete
echo "Test for get_repository_discussions completed."
