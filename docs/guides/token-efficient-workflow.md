# Token-Efficient Development Workflow with Delegate

## Overview

This guide presents a revolutionary workflow that enables developers to work with massive codebases while using 95%+ fewer tokens. By leveraging the delegate MCP server, you can generate, fix, and iterate on large code projects without ever loading the full code into your primary Claude context.

## The Token Problem

Traditional AI-assisted development faces a critical limitation:
- Reading a 10,000 line codebase: ~50,000 tokens
- Making changes: ~50,000 tokens (to write it back)
- Fixing errors: Another ~50,000 tokens
- **Total: 150,000+ tokens for one iteration!**

## The Delegate Solution

With delegate, the same task uses:
- Orchestration: ~500 tokens
- Error analysis: ~500 tokens  
- Fix instructions: ~1,000 tokens
- **Total: ~2,000 tokens (98.7% reduction!)**

## Core Workflow

### 1. Initial Generation (Single File at a Time)

```yaml
# IMPORTANT: Generate ONE file per invoke call
delegate_invoke:
  model: gemini-2.5-flash    # Fast & cost-effective
  files:                     # Attach context files
    - "docs/api-spec.md"     
    - "src/interfaces.go"    
    - "examples/similar.go"  
  prompt: |
    Create models/user.go that:
    - Implements UserInterface from interfaces.go
    - Follows the patterns in similar.go
    - Uses GORM for database mapping
  timeout: 90                # 60-90s for code, 90-120s for docs
```

### 2. Check & Write Pattern

```bash
# 1. Check output size (instant, no tokens)
delegate_check(output_id) 
# → "15KB, ~3000 tokens"

# 2. Optional: Peek at structure (minimal tokens)
delegate_read(output_id, max_tokens: 200, extract: "code")

# 3. Write directly to project (ZERO tokens!)
# For source code files, use extract: "code" to strip markdown fences
delegate_read(output_id, write_to: "/absolute/path/to/src/services/user_service.go", extract: "code")
# → "Content written to src/services/user_service.go (15.2 KB, ~3800 tokens saved)"
```

### 2.1 NEW: Multi-Block Handling

When LLMs return multiple code blocks (common with "show me an example" prompts), delegate now provides intelligent block selection:

```bash
# Generate a React component with tests and styles
delegate_invoke(
  model: "gemini-2.5-flash",
  prompt: "Create a React TodoList component with tests and CSS styling"
)

# Attempt to write - delegate detects multiple blocks!
delegate_read(output_id, write_to: "/path/to/TodoList.jsx", extract: "code")
# → Warning: Multiple code blocks found (3 blocks). Use block_index option to select specific block.
#   
#   Block 0: jsx - "import React, { useState } from 'react'..." (4.3 KB, 150 lines)
#   Block 1: jsx - "import { render } from '@testing-library/react'..." (1.2 KB, 45 lines)  
#   Block 2: css - ".todo-container { ..." (892 bytes, 34 lines)

# Now you can intelligently select each block for the right file:
delegate_read(output_id, write_to: "/path/to/TodoList.jsx", extract: "code", block_index: 0)
# → "Content written to TodoList.jsx (4.3 KB, ~1075 tokens saved)"

delegate_read(output_id, write_to: "/path/to/TodoList.test.jsx", extract: "code", block_index: 1)
# → "Content written to TodoList.test.jsx (1.2 KB, ~300 tokens saved)"

delegate_read(output_id, write_to: "/path/to/TodoList.css", extract: "code", block_index: 2)
# → "Content written to TodoList.css (892 bytes, ~223 tokens saved)"
```

This feature eliminates the "merged code problem" where multiple files get concatenated into one!

### 3. Compile-Fix Loop

The magic happens here - let compilers do the analysis:

```bash
# Run compiler, capture errors to file
go build ./src/services/user_service.go 2> build_errors.txt

# Delegate reads its OWN output + errors
delegate_invoke:
  model: gemini-2.5-flash
  files:
    - "/absolute/path/to/src/services/user_service.go"  # The file it just wrote
    - "/absolute/path/to/build_errors.txt"              # Specific errors to fix
  prompt: "Fix only these specific compilation errors"
  code_only: true

# Write the fixed version (extract: "code" for clean source)
delegate_read(new_output_id, write_to: "/absolute/path/to/src/services/user_service.go", extract: "code")
```

### 4. Test-Fix Loop

Same pattern for tests:

```bash
# Run tests, capture failures
go test ./src/services/... > test_results.txt 2>&1

# Delegate fixes based on test output
delegate_invoke:
  model: gemini-2.5-flash
  files:
    - "/absolute/path/to/src/services/user_service.go"
    - "/absolute/path/to/src/services/user_service_test.go"
    - "/absolute/path/to/test_results.txt"
  prompt: "Fix the failing tests by updating the implementation"
```

## Token Saving Strategies

### 1. Never Read Generated Code
- ❌ Read full file → Edit → Write (uses full tokens)
- ✅ Let delegate read its own outputs (uses zero tokens)

### 2. Use Compilers as Analyzers
- Compilers provide precise error locations
- Linters identify style issues
- Test runners show exact failures
- All produce small, focused output files

### 3. Chain Small Operations
Instead of one large prompt, chain focused operations:
```
Generate core → Fix syntax → Add validation → Fix tests → Add docs
```

### 4. Model Selection Strategy
- **gemini-2.5-flash**: 95% of tasks (fast, cheap, capable)
- **gemini-2.5-pro**: Complex architectural decisions
- **claude-sonnet-4**: When you need Claude's specific strengths
- **claude-opus-4**: Only for the most complex reasoning tasks

## Real-World Example: Building a REST API

Traditional approach (150,000+ tokens):
1. Read existing codebase (50k tokens)
2. Generate new code in context (50k tokens)  
3. Fix errors through multiple iterations (50k+ tokens)

Delegate approach (~2,000 tokens):
```bash
# 1. Generate API with context (~500 tokens for orchestration)
delegate_invoke(model: "gemini-2.5-flash", 
                files: ["/absolute/path/to/api/openapi.yaml", "/absolute/path/to/internal/base_controller.go"],
                prompt: "Generate complete REST API from OpenAPI spec")

# 2. Write to project (0 tokens) - extract: "code" for clean source files
delegate_read(output_id, write_to: "/absolute/path/to/internal/api/controllers.go", extract: "code")

# 3. Fix compilation errors (~500 tokens)
go build ./internal/api/... 2> errors.txt
delegate_invoke(files: ["/absolute/path/to/internal/api/controllers.go", "/absolute/path/to/errors.txt"],
                prompt: "Fix compilation errors")

# 4. Fix failing tests (~1000 tokens)
go test ./internal/api/... > test_failures.txt 2>&1
delegate_invoke(files: ["/absolute/path/to/internal/api/controllers.go", "/absolute/path/to/test_failures.txt"],
                prompt: "Fix failing tests")

# Total: ~2,000 tokens vs 150,000+ tokens!
```

## Advanced Patterns

### Multi-File Generation
Generate entire packages by having delegate create file manifests:

```yaml
delegate_invoke:
  prompt: |
    Generate a complete user management system:
    Return a JSON with file paths and contents:
    {
      "src/models/user.go": "...",
      "src/services/user_service.go": "...",
      "src/handlers/user_handler.go": "..."
    }
```

Then use a script to extract and write each file.

### Incremental Enhancement
Add features without reading existing code:

```yaml
delegate_invoke:
  files: ["/absolute/path/to/src/service.go"]  # Delegate reads it, you don't
  prompt: "Add caching to all database queries in this service"
```

### Automated Testing Loop
```bash
while ! go test ./...; do
  go test ./... 2>&1 > test_output.txt
  delegate_invoke(
    files: ["/absolute/path/to/src/", "/absolute/path/to/test_output.txt"],
    prompt: "Fix the first failing test"
  )
  delegate_read(output_id, write_to: "/absolute/path/to/src/")
done
```

## Best Practices

1. **Always write to disk first**: Never read generated code into your context unless absolutely necessary

2. **Use precise prompts**: Since delegate can't ask clarifying questions, be specific

3. **Leverage file context**: Attach examples, interfaces, and specs rather than explaining

4. **Trust the process**: Resist the urge to "just quickly read" the generated code

5. **Think in pipelines**: Each step transforms output for the next, like Unix pipes

## Conclusion

This workflow transforms AI development from a token-intensive process to an efficient orchestration task. You become a conductor, not a code carrier, maintaining context for decision-making while delegating the heavy lifting.

The result: 95%+ token savings, faster development, and the ability to work with codebases of any size without context limitations.