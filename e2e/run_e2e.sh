#!/bin/bash
# E2E Test Runner for Delegate

echo "🧪 Delegate E2E Test Suite"
echo "========================="
echo ""

# Load .env file if it exists
if [ -f ".env" ]; then
    echo "Loading .env file..."
    export $(cat .env | grep -v '^#' | xargs)
elif [ -f "../.env" ]; then
    echo "Loading ../.env file..."
    export $(cat ../.env | grep -v '^#' | xargs)
fi

# Check for API keys
if [ -z "$GOOGLE_API_KEY" ] && [ -z "$ANTHROPIC_API_KEY" ]; then
    echo "⚠️  WARNING: No API keys found!"
    echo ""
    echo "To run E2E tests with real LLM calls, set one of:"
    echo "  export GOOGLE_API_KEY=your-key-here"
    echo "  export ANTHROPIC_API_KEY=your-key-here"
    echo ""
    echo "Tests will be skipped without API keys."
    echo ""
fi

# Run E2E tests
echo "Running E2E tests..."
cd "$(dirname "$0")/.." || exit 1

# Build the main binary first to catch any compilation errors
echo "Building delegate..."
if ! /usr/local/go/bin/go build -o delegate main.go; then
    echo "❌ Build failed!"
    exit 1
fi
rm delegate

# Run the E2E tests
/usr/local/go/bin/go test -v --tags=e2e ./e2e/...

# Check exit code
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ All E2E tests passed!"
else
    echo ""
    echo "❌ Some E2E tests failed!"
    exit 1
fi