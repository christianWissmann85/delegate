# Claude Code Integration Guide - Delegate

## What is Delegate?

Delegate is YOUR tool - built specifically to help you (Claude Code) generate code without eating up your context window. It's the industrial-strength replacement for AAG, with just 3 simple commands.

## Quick Start

### Your MCP Tools
```javascript
// That's it. Three tools. Nothing else.
delegate_invoke(params)  // Generate code with another LLM
delegate_check(params)   // Check output size before reading
delegate_read(params)    // Read the output (or parts of it)
```

## Core Workflow Pattern

```javascript
// 1. Generate code without using your tokens
const output = await delegate_invoke({
    model: "gemini-2.5-flash",
    prompt: "Create a complete Express.js REST API for a todo app",
    files: ["requirements.md", "database_schema.sql"]
});

// 2. Always check size before reading!
const info = await delegate_check({
    output_id: output.id
});
console.log(`Output is ${info.size_kb}KB (≈${info.estimated_tokens} tokens)`);

// 3. Read strategically
if (info.estimated_tokens < 1000) {
    // Small enough - read everything
    const result = await delegate_read({
        output_id: output.id,
        options: { extract: "all" }
    });
} else {
    // Too big - just get the code
    const result = await delegate_read({
        output_id: output.id,
        options: { extract: "code", max_tokens: 2000 }
    });
}
```

## Model Selection Guide

### gemini-2.5-flash (Your Workhorse)
- **When**: Most code generation tasks
- **Why**: Fast, cheap, huge context window (1M tokens)
- **Example**: API endpoints, data models, utility functions

### gemini-2.5-pro (Heavy Lifting)
- **When**: Complex architectural decisions, system design
- **Why**: Advanced reasoning, still with 1M context
- **Example**: Microservice architecture, complex algorithms

### claude-sonnet-4-20250514 (Precision Work)
- **When**: Need precise instruction following
- **Why**: Best at following detailed specifications
- **Example**: Implementing to strict standards, refactoring

### claude-opus-4-20250514 (Crown Jewel)
- **When**: Security-critical code, complex business logic
- **Why**: Highest quality output, best understanding
- **Example**: Authentication systems, payment processing

## Common Patterns

### Pattern 1: Large Feature Implementation
```javascript
// User wants a complete feature
const output = await delegate_invoke({
    model: "gemini-2.5-flash",
    prompt: `Create a complete user authentication system with:
    - JWT tokens
    - Password reset via email
    - Role-based access control
    - Rate limiting
    Include all models, routes, middleware, and tests.`,
    files: ["tech_stack.md", "existing_auth_code.js"]
});

// Check what we got
const info = await delegate_check({ output_id: output.id });
// Likely 10-20KB of code

// Read in chunks
const models = await delegate_read({
    output_id: output.id,
    options: { extract: "code", language: "javascript", max_tokens: 1000 }
});
// Review, then get more...
```

### Pattern 2: Code Analysis/Refactoring
```javascript
// Attach existing code for context
const output = await delegate_invoke({
    model: "claude-sonnet-4-20250514",  // Use Claude for analysis
    prompt: "Analyze this code for security vulnerabilities and suggest fixes",
    files: ["src/auth/login.js", "src/auth/session.js"]
});

// Read the analysis
const analysis = await delegate_read({
    output_id: output.id,
    options: { extract: "explanation" }  // Just the analysis, no code
});
```

### Pattern 3: Iterative Development
```javascript
// First pass - basic structure
const v1 = await delegate_invoke({
    model: "gemini-2.5-flash",
    prompt: "Create the basic structure for a GraphQL API server"
});

// Check and read
const structure = await delegate_read({ 
    output_id: v1.id, 
    options: { extract: "code" }
});

// Second pass - add specific features
const v2 = await delegate_invoke({
    model: "gemini-2.5-pro",  // Upgrade model for complex logic
    prompt: "Add user authentication to this GraphQL server",
    files: ["generated_structure.js"]  // Feed back the previous output
});
```

### Pattern 4: Document Analysis (Your Context Saver!)
```javascript
// Scenario: Need to analyze multiple large documents
const analysis = await delegate_invoke({
    model: "gemini-2.5-pro",  // 1M token context window!
    prompt: `Analyze these architecture documents and extract:
    1. All API endpoint patterns
    2. Authentication strategies used
    3. Database schema decisions
    4. Testing approaches
    
    Provide a structured summary with examples.`,
    files: ["arch-doc-1.md", "arch-doc-2.md", "arch-doc-3.md", "api-spec.md"]
});

// I get a focused summary instead of reading 20k lines
const insights = await delegate_read({ output_id: analysis.id });
```

### Pattern 5: Multi-Document Research
```javascript
// Research across massive documentation
const research = await delegate_invoke({
    model: "gemini-2.5-pro",
    prompt: `Read all these docs and answer:
    - How is error handling implemented across the codebase?
    - What patterns are used for data validation?
    - Are there any security best practices mentioned?
    
    Cite specific examples with file names.`,
    files: ["docs/**/*.md"]  // Could be 100+ files!
});

// Extract just the findings
const findings = await delegate_read({
    output_id: research.id,
    options: { max_tokens: 2000 }  // Get concise results
});
```

### Pattern 6: Codebase Analysis
```javascript
// Analyze entire codebases without filling my context
const review = await delegate_invoke({
    model: "gemini-2.5-flash",  // Fast for large volume
    prompt: `Review this codebase for:
    - Potential security vulnerabilities
    - Performance bottlenecks
    - Code quality issues
    - Missing error handling
    
    Focus on critical issues only.`,
    files: ["src/**/*.js", "lib/**/*.js"]  // Thousands of lines!
});

// Get actionable insights
const issues = await delegate_read({
    output_id: review.id,
    options: { extract: "all" }
});
```

## Pro Tips

### 1. Always Check Before Reading
```javascript
// ❌ BAD - Might consume 10k tokens unexpectedly
const result = await delegate_read({ output_id: output.id });

// ✅ GOOD - Know what you're getting into
const info = await delegate_check({ output_id: output.id });
if (info.estimated_tokens > 5000) {
    // Too big! Extract just what you need
    const code = await delegate_read({
        output_id: output.id,
        options: { extract: "code", max_tokens: 2000 }
    });
}
```

### 2. Use File Attachments Liberally
```javascript
// ❌ BAD - LLM has no context
await delegate_invoke({
    prompt: "Update the API to handle the new requirements"
});

// ✅ GOOD - Clear context
await delegate_invoke({
    prompt: "Update the API to handle the new requirements",
    files: ["new_requirements.md", "current_api.js", "test_cases.js"]
});
```

### 3. Extract Strategically
```javascript
// If output has both code and explanation:
// - First read just the code to implement
// - Then read explanation if user asks questions

const code = await delegate_read({
    output_id: output.id,
    options: { extract: "code" }
});
// Implement the code...

// Later, if user asks "why did you do X?"
const explanation = await delegate_read({
    output_id: output.id,
    options: { extract: "explanation", max_tokens: 500 }
});
```

### 4. Handle Errors Gracefully
```javascript
try {
    const output = await delegate_invoke({
        model: "gemini-2.5-flash",
        prompt: "Generate code",
        max_tokens: 8000
    });
} catch (error) {
    if (error.code === 'TIMEOUT') {
        // Try with smaller scope or different model
    } else if (error.code === 'PROVIDER_ERROR') {
        // API key issue or rate limit
    }
}
```

## What NOT to Do

### ❌ Don't Chain Too Many Calls
Each invoke is 2-30 seconds. Users get impatient after 3-4 calls.

### ❌ Don't Read Without Checking
You might consume your entire context on one file.

### ❌ Don't Ignore File Attachments
LLMs perform much better with context.

### ❌ Don't Use Wrong Models
- Don't use Opus for simple boilerplate (expensive, slow)
- Don't use Flash for security-critical code (fast but less thorough)

## Model Decision Tree

```
Is it security/payment related?
├─ Yes → claude-opus-4-20250514
└─ No → Is it architecturally complex?
    ├─ Yes → gemini-2.5-pro
    └─ No → Does it need strict spec adherence?
        ├─ Yes → claude-sonnet-4-20250514
        └─ No → gemini-2.5-flash (default)
```

## Advanced Features

### Code-Only Mode
```javascript
// When you just need the code without explanations
const output = await delegate_invoke({
    model: "gemini-2.5-flash",
    prompt: "Create a Python function to calculate fibonacci numbers",
    code_only: true  // Returns only code blocks
});

// Reading will return just the code
const code = await delegate_read({ output_id: output.id });
```

### Language Hints for Better Extraction
```javascript
// Help the extractor by specifying expected language
const output = await delegate_invoke({
    model: "gemini-2.5-pro",
    prompt: "Create a REST API with TypeScript and tests",
    language_hint: "typescript"  // Improves extraction accuracy
});
```

### Custom Timeouts for Long Tasks
```javascript
// Override default timeout for complex generations
const output = await delegate_invoke({
    model: "claude-opus-4-20250514",
    prompt: "Generate a complete microservices architecture...",
    timeout: 120  // 2 minutes instead of default 60s
});
```

## Error Handling Examples

### Handling Provider Errors
When Delegate encounters provider issues, it returns structured errors that help me make smart decisions:

```javascript
try {
    const output = await delegate_invoke({
        model: "gemini-2.5-flash",
        prompt: "Generate a complex React dashboard"
    });
} catch (error) {
    if (error.error === "rate_limited") {
        // Option 1: Wait and retry
        console.log(`Rate limited. Waiting ${error.retry_after}s...`);
        
        // Option 2: Try alternative model
        const output = await delegate_invoke({
            model: error.alternative_models[0], // e.g., "claude-sonnet-4-20250514"
            prompt: "Generate a complex React dashboard"
        });
        
        // Option 3: I'll handle it directly
        console.log("Both providers are busy. I'll generate this code directly.");
    }
}
```

### Smart Recovery Patterns
```javascript
// Pattern 1: Try fast model first, fall back to powerful model
async function generateWithFallback(prompt) {
    try {
        return await delegate_invoke({ model: "gemini-2.5-flash", prompt });
    } catch (error) {
        if (error.error === "provider_unavailable") {
            console.log("Gemini unavailable, trying Claude...");
            return await delegate_invoke({ model: "claude-opus-4-20250514", prompt });
        }
        throw error;
    }
}

// Pattern 2: Inform user and let them decide
async function generateWithUserChoice(prompt) {
    try {
        return await delegate_invoke({ model: "gemini-2.5-pro", prompt });
    } catch (error) {
        if (error.error === "rate_limited") {
            console.log(`Gemini is rate limited (retry in ${error.retry_after}s).`);
            console.log("Should I: 1) Wait and retry, 2) Use Claude, or 3) Generate directly?");
            // Handle based on user preference
        }
    }
}
```

## Troubleshooting

### "Output not found"
- Output files expire after 24 hours
- Check the output_id is correct

### "Timeout error"
- Default 60-second limit
- Simplify prompt or break into smaller tasks

### "Extraction failed"
- LLM didn't format code properly
- Try `extract: "all"` and parse manually

### "Provider error"
- Check API key is set correctly
- May be rate limited - wait and retry

### "File too large"
- Max 100KB per attached file
- Split large files or summarize first

## Example: Complete Feature Flow

```javascript
// User: "Create a real-time chat application with rooms"

// 1. Generate the data models
const models = await delegate_invoke({
    model: "gemini-2.5-flash",
    prompt: `Create data models for a real-time chat app:
    - Users (with online status)
    - Rooms (public/private)
    - Messages (with read receipts)
    Include Mongoose schemas with all validations.`
});

// 2. Check size
const modelsInfo = await delegate_check({ output_id: models.id });
console.log(`Models: ${modelsInfo.size_kb}KB`);

// 3. Read and implement models
const modelCode = await delegate_read({
    output_id: models.id,
    options: { extract: "code" }
});

// 4. Generate WebSocket handlers
const websocket = await delegate_invoke({
    model: "gemini-2.5-pro",  // More complex, upgrade model
    prompt: "Create Socket.io handlers for real-time chat with rooms",
    files: ["generated_models.js"]  // Pass the models as context
});

// 5. Generate frontend components
const frontend = await delegate_invoke({
    model: "gemini-2.5-flash",
    prompt: "Create React components for the chat interface",
    files: ["socket_events.js", "ui_mockup.png"]
});

// Continue pattern...
```

## Remember

Delegate is YOUR tool. It's designed to:
- Save your tokens for thinking, not generating
- Be boringly reliable (no fancy features to break)
- Get out of your way

When in doubt: `invoke` → `check` → `read`. That's it!

## Quick Reference

| Task | Model | Why |
|------|-------|-----|
| Boilerplate | `gemini-2.5-flash` | Fast & cheap |
| Complex logic | `gemini-2.5-pro` | Better reasoning |
| Following specs | `claude-sonnet-4-20250514` | Precise adherence |
| Critical code | `claude-opus-4-20250514` | Highest quality |

**The Golden Rule**: Always `check` before you `read`!