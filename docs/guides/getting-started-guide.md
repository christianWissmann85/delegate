# **🚀 Delegate: Getting Started Guide**

Welcome to Delegate - the dead-simple way to save Claude Code's context tokens by delegating code generation to other LLMs!

## **What is Delegate?**

Delegate is an MCP server that gives Claude Code three simple tools:
- **invoke** (STEP 1) - Generate code with Gemini or Claude models, each file <1MB
- **check** (STEP 2) - See how big the output is before reading
- **read** (STEP 3) - Get the generated content or write directly to disk (saves tokens!)

That's it. No complexity. Just industrial-strength delegation.

## **⚡ Quick Start (2 minutes)**

### **1. Get your API keys**

You'll need at least one:
- **Gemini**: Get it at [Google AI Studio](https://makersuite.google.com/app/apikey)
- **Claude**: Get it at [Anthropic Console](https://console.anthropic.com/)

### **2. Set your API keys**

Add to your shell profile (`~/.bashrc` or `~/.zshrc`):
```bash
export GOOGLE_API_KEY="your-gemini-key"
export ANTHROPIC_API_KEY="your-claude-key"
```

Reload your shell:
```bash
source ~/.bashrc  # or ~/.zshrc
```

### **3. Install Delegate**

One command:
```bash
claude mcp add delegate -s user -- npx -y @christianwissmann85/delegate
```

### **4. Test it!**

Start Claude Code:
```bash
claude
```

Try these examples:

**Code Generation:**
```
Use Delegate to generate a Python web server with Gemini
```

**Document Analysis:**
```
Use Delegate to analyze my project docs and summarize the API design patterns
```

**Large File Processing:**
```
Use Delegate to review this codebase and find all error handling patterns
```

Claude Code will automatically use Delegate without consuming its own precious tokens! 🎉

## **📖 How It Works**

### **The Token-Saving Workflow**

1. **You ask Claude Code for something big**
   ```
   Create a complete REST API for a blog system
   ```

2. **Claude Code uses Delegate to generate it**
   ```javascript
   // Behind the scenes:
   await delegate_invoke({
     model: "gemini-2.5-flash",
     prompt: "Create a complete REST API...",
     files: ["/absolute/path/to/requirements.md"]  // Must use absolute paths!
   })
   ```

3. **Claude Code checks the size first**
   ```javascript
   // Smart token management:
   const info = await delegate_check({ output_id: "out_123" });
   // Returns: { size_kb: 15.2, estimated_tokens: 3800, has_code: true }
   ```

4. **Claude Code reads strategically**
   ```javascript
   // Only read what's needed:
   const code = await delegate_read({ 
     output_id: "out_123",
     options: { extract: "code" }
   });
   ```

## **🎯 Essential Commands**

### **In Claude Code**

Just talk naturally! Claude Code knows when to use Delegate:
- "Generate a complex React component"
- "Create this with Gemini for speed"
- "Use Claude Opus for this critical payment logic"

### **Check what models are available**
```
/mcp
```
You should see `delegate: connected`

### **For debugging**
```bash
# See MCP logs
tail -f ~/.claude/logs/mcp-server-delegate.log

# Check your outputs
ls .delegate/outputs/
```

## **🆕 Multi-Block Handling**

When generating code that includes multiple files (like a component + tests + styles), Delegate now intelligently handles them:

```javascript
// You ask for:
"Create a React TodoList component with tests and CSS"

// Delegate generates multiple code blocks
// When you try to save:
await delegate_read({ 
  output_id: "out_123", 
  options: { 
    extract: "code", 
    write_to: "/path/to/TodoList.jsx" 
  }
});

// You'll see:
"Warning: Multiple code blocks found (3 blocks). Use block_index option to select specific block.

Block 0: jsx - "import React, { useState } from 'react'..." (4.3 KB, 150 lines)
Block 1: jsx - "import { render } from '@testing-library/react'..." (1.2 KB, 45 lines)  
Block 2: css - ".todo-container { ..." (892 bytes, 34 lines)"

// Now you can save each one properly:
await delegate_read({ output_id: "out_123", options: { extract: "code", write_to: "/path/to/TodoList.jsx", block_index: 0 }});
await delegate_read({ output_id: "out_123", options: { extract: "code", write_to: "/path/to/TodoList.test.jsx", block_index: 1 }});
await delegate_read({ output_id: "out_123", options: { extract: "code", write_to: "/path/to/TodoList.css", block_index: 2 }});
```

This prevents the common problem of multiple files being merged into one!

## **💡 Pro Tips**

### **1. Let Claude Code Choose**
Don't specify a model unless you have a preference. Claude Code is smart about picking the right one.

### **2. Use Context Files**
Delegate can include files for context:
```
Use Delegate to update this API based on the new requirements in requirements.md
```

### **3. Save Tokens on Large Tasks**
Perfect for:
- Generating entire applications
- Large refactoring tasks
- Creating comprehensive test suites
- Writing extensive documentation
- **Analyzing multiple large documents**
- **Processing entire codebases**

### **4. Know Your Models**

| Need | Use | Why |
|------|-----|-----|
| Speed | `gemini-2.5-flash` | Lightning fast, great for iteration |
| Complex logic | `gemini-2.5-pro` | Advanced reasoning |
| Best quality | `claude-opus-4-20250514` | When it must be perfect |
| Balanced | `claude-sonnet-4-20250514` | Great all-rounder |
| **Document Analysis** | `gemini-2.5-pro` | **1M token context window!** |

## **📚 Document Analysis - The Game Changer**

### **Why This Matters**
Claude Code has limited context tokens. Reading multiple large documents quickly exhausts them. Delegate solves this!

### **Example: Analyzing Multiple Docs**
Instead of:
```
Read architecture.md, api-spec.md, database-design.md and tell me about the auth strategy
```
(This would consume 50%+ of Claude Code's context!)

Do this:
```
Use Delegate to analyze architecture.md, api-spec.md, database-design.md and summarize the authentication strategy
```
(Claude Code stays fresh for actual work!)

### **Real-World Scenarios**
- **"Analyze all test files and tell me what's not covered"**
- **"Review these 10 RFC documents and extract the key decisions"**
- **"Read the entire codebase and find inconsistent naming patterns"**
- **"Process these log files and identify error patterns"**

## **🛠️ Configuration**

### **Optional Environment Variables**
```bash
# More detailed logging
export DELEGATE_LOG_LEVEL="debug"

# Longer timeout for huge generations
export DELEGATE_TIMEOUT_SECONDS="120"

# Custom output directory
export DELEGATE_OUTPUT_DIR="/tmp/delegate"
```

### **Project-specific Setup**
```bash
# Just for this project
cd /your/project
claude mcp add delegate -s project -- npx -y @christianwissmann85/delegate
```

## **❓ Common Questions**

**Q: How much context can I send?**
- Gemini models: Up to 1 million tokens!
- Claude models: Up to 200k tokens

**Q: Where are outputs stored?**
- In `.delegate/outputs/` in your current directory
- Auto-cleaned after 24 hours

**Q: Can I read outputs later?**
- Yes! Use the output ID with `check` and `read`
- IDs look like: `out_20250620_143022`

**Q: What if generation times out?**
- Increase timeout: `export DELEGATE_TIMEOUT_SECONDS=120`
- Or break into smaller tasks

## **🚨 Troubleshooting**

### **"API key error"**
```bash
# Check your keys are set:
echo $GOOGLE_API_KEY
echo $ANTHROPIC_API_KEY
```

### **"MCP server not connected"**
```bash
# Reinstall:
claude mcp remove delegate
claude mcp add delegate -s user -- npx -y @christianwissmann85/delegate
```

### **"Output not found"**
- Outputs expire after 24 hours
- Check the exact ID from the invoke response

## **📚 Next Steps**

1. **Read the Model Reference Card** - Understand when to use each model
2. **Check out claude-code-guide-updated.md** - Advanced patterns and examples
3. **Review NO_SCOPE_CREEP.md** - Understand why Delegate is so simple (and staying that way!)

## **🎉 You're Ready!**

Start saving tokens and generating more code with Delegate. Remember:
- **invoke** to generate
- **check** before reading  
- **read** strategically

Happy coding! 🚀

---

*Built with ❤️ and a strict No Scope Creep policy*