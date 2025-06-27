# NO SCOPE CREEP - The Sacred Document

**Reviewed by Christian Wissmann for Delegate V2.0**

## The Prime Directive
**If it's not one of the 4 tools (`submit_task`, `get_output_metadata`, `get_output_content`, `write_output_to_file`), it doesn't exist.**

## Things We Will NOT Build (No Matter How "Easy")

### ❌ Session Management
- "But what about tracking usage across..." - **NO**
- "It would be simple to add session..." - **NO**
- "Other tools have session management..." - **We are not other tools**

### ❌ Token Counting
- "We could add accurate token counting..." - **NO**
- "Just a simple tokenizer..." - **NO**
- Token estimation (bytes/4) is enough

### ❌ Progress Indicators
- "Users want to see progress..." - **NO**
- "Just a simple percentage..." - **NO**
- "Streaming updates..." - **Absolutely NO**

### ❌ Web UI / CLI
- "A simple web interface..." - **NO**
- "Just a basic CLI for testing..." - **NO**
- MCP only. Period.

### ❌ Multiple Storage Backends
- "Support for S3..." - **NO**
- "Database storage..." - **NO**
- "Network file systems..." - **NO**
- Local filesystem only.

### ❌ Complex Routing/Orchestration
- "Route based on prompt type..." - **NO**
- "Automatic model selection..." - **NO**
- "Load balancing..." - **NO**
- Explicit model parameter only.

### ❌ Command System
- "Create/Edit/Analyze commands..." - **NO**
- "It's just organizing the prompts..." - **NO**
- That's how AAG died. Learn from history.

### ❌ Analytics/Metrics Dashboard
- "Track success rates..." - **NO**
- "Usage analytics..." - **NO**
- "Performance metrics..." - **NO**
- Basic logs are enough.

### ❌ Conversation Management
- "Multi-turn conversations..." - **NO**
- "Discussion features..." - **NO**
- "Context management..." - **NO**
- Single prompt, single response.

### ❌ Batch Operations (v1.0)
- "Process multiple files..." - **NO**
- "Parallel execution..." - **NO**

### ❌ Middleware/Plugins
- "Extensibility is important..." - **NO**
- "Plugin architecture..." - **NO**
- "Hooks for customization..." - **NO**

### ❌ Advanced Error Recovery
- "Automatic fallback models..." - **NO**
- "Smart retry strategies..." - **NO**
- Simple 3-retry with backoff. That's it.

## The Slippery Slope Examples

### Example 1: "Just Add Request IDs"
- Day 1: "We need request IDs for debugging"
- Day 3: "Now we need request tracking"
- Day 5: "Let's add request history"
- Day 10: "We need a database for the history"
- Day 20: You've built AAG again

**Answer: NO**. Use timestamps in logs.

### Example 2: "Simple Progress Updates"
- "Just emit a 'started' event"
- "Now add 'progress' percentage"
- "Stream the LLM responses"
- "Buffer management for streaming"
- "Backpressure handling"

**Answer: NO**. Request -> Response. Nothing in between.

### Example 3: "Basic Validation"
- "Validate the prompt isn't empty"
- "Check for prohibited content"
- "Add prompt templates"
- "Template variables"
- "Template management system"

**Answer: NO**. Minimal validation only (required fields).

## When Someone Asks for a Feature

### The Decision Tree
```
Is it submit_task, get_output_metadata, get_output_content, or write_output_to_file?
├─ Yes: Consider it
│   └─ Does it add complexity?
│       ├─ Yes: NO
│       └─ No: Maybe
└─ No: NO
```

### Standard Responses
- "That's a great idea for v3" (Translation: Never)
- "Let's see how v2 performs first" (Translation: No)
- "That would complicate the core design" (Translation: Obviously no)
- "Check NO_SCOPE_CREEP.md" (Translation: Read this document)

## Features That Seem Harmless But Aren't

1. **"Just cache the responses"**
   - Cache invalidation
   - Storage management
   - Cache configuration
   - Before you know it: Redis

2. **"Add JSON Schema validation"**
   - Schema versioning
   - Migration strategies
   - Validation error messages
   - Suddenly you're building a framework

3. **"Support YAML/TOML config"**
   - Config validation
   - Config reloading
   - Config migration
   - Now you have a config management system

## The Mantra

When in doubt, chant:
- **Four tools**
- **One purpose each**
- **Zero complexity**
- **No scope creep**

## Remember Why AAG Failed

AAG started simple and became:
- 3,283 lines in orchestrator.py alone
- Session management
- Token tracking
- Progress reporting  
- Command routing
- Discussion coordination
- Batch processing
- Analysis frameworks
- Review systems

**Delegate will not repeat history.**

## The Enforcement

This document is **sacred**. 
- Print it
- Frame it
- Look at it when tempted
- Say NO to scope creep

**Every feature that isn't `submit_task`/`get_output_metadata`/`get_output_content`/`write_output_to_file` is a step toward another failed refactor.**

Stay strong. Ship simple. Win.