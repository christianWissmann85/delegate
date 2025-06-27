# **🚀 Delegate: Getting Started Guide**

**Reviewed by Christian Wissmann for Delegate V2.0**

Welcome to Delegate - the dead-simple way to save Claude Code's context tokens by delegating code generation to other LLMs!

## **What is Delegate?**

Delegate is an MCP server that gives Claude Code four simple tools:
- **submit_task** (STEP 1) - Send generation tasks to Gemini or Claude models  
- **get_output_metadata** (STEP 2) - Check what was generated before reading
- **get_output_content** (STEP 3) - Get generated content into Claude's context
- **write_output_to_file** (STEP 4) - Save content directly to disk (saves tokens!)

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

2. **Claude Code submits the task to Delegate**
   ```javascript
   // Behind the scenes:
   await delegate_submit_task({
     model: "gemini-2.5-flash",
     prompt: "Create a complete REST API...",
     files: ["docs/requirements.md", "src/models.go"]  // Now uses relative paths!
   })
   ```

3. **Claude Code checks what was generated**
   ```javascript
   // Smart token management:
   const info = await delegate_get_output_metadata({ output_id: "out_123" });
   // Returns structured JSON with size, blocks, token estimates
   ```

4. **Claude Code saves directly to disk (zero tokens!)**
   ```javascript
   // The magic zero-token workflow:
   await delegate_write_output_to_file({ 
     output_id: "out_123",
     write_to: "src/api/blog.go",
     options: { extract: "code", block_index: 0 }
   });
   ```

### **🎯 The Zero-Token Workflow**

The most powerful pattern is the **submit → write** workflow:

```javascript
// Step 1: Submit task (low tokens)
const { output_id } = await delegate_submit_task({
  prompt: "Create a React TodoList component with tests",
  files: ["src/types.ts"]
});

// Step 2: Write directly to file (ZERO tokens!)
await delegate_write_output_to_file({
  output_id,
  write_to: "src/components/TodoList.tsx",
  options: { extract: "code", block_index: 0 }
});
```

**Result:** You get a complete component file without consuming any of Claude Code's context tokens for the generated content!

## **🆕 Multi-Block Handling Made Simple**

When Delegate generates multiple files or sections, everything is now structured data:

```javascript
// Check what was generated
const metadata = await delegate_get_output_metadata({ output_id: "out_123" });

// You get structured JSON like this:
{
  "content_analysis": {
    "blocks_found": 3,
    "blocks": [
      { 
        "index": 0, 
        "language": "tsx", 
        "size_kb": 4.3, 
        "lines": 150, 
        "preview": "import React, { useState } from 'react';" 
      },
      { 
        "index": 1, 
        "language": "tsx", 
        "size_kb": 1.2, 
        "lines": 45, 
        "preview": "import { render } from '@testing-library/react';" 
      },
      { 
        "index": 2, 
        "language": "css", 
        "size_kb": 0.9, 
        "lines": 34, 
        "preview": ".todo-container {" 
      }
    ]
  }
}

// Now save each block to the right file:
await delegate_write_output_to_file({ 
  output_id: "out_123", 
  write_to: "src/TodoList.tsx", 
  options: { block_index: 0 } 
});

await delegate_write_output_to_file({ 
  output_id: "out_123", 
  write_to: "src/TodoList.test.tsx", 
  options: { block_index: 1 } 
});

await delegate_write_output_to_file({ 
  output_id: "out_123", 
  write_to: "src/TodoList.css", 
  options: { block_index: 2 } 
});
```

**No more string parsing!** Everything is structured, predictable JSON.

## **🎯 Essential Commands**

### **In Claude Code**

Just talk naturally! Claude Code knows when to use Delegate:
- "Generate a complex React component and save it to src/components/"
- "Create this with Gemini for speed"
- "Use Claude Opus for this critical payment logic"
- "Analyze these docs and write a summary to docs/analysis.md"

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

## **💡 Pro Tips**

### **1. Master the Zero-Token Pattern**
For any code generation task:
```
Generate [something] and save it to [relative/path/file.ext]
```
Claude Code will use the `submit_task` → `write_output_to_file` pattern automatically.

### **2. Use Relative Paths**
All paths are now relative to your project root:
- ✅ `src/components/Button.tsx`
- ✅ `docs/api.md`  
- ✅ `tests/integration/auth.test.js`
- ❌ `/home/user/project/src/Button.tsx` (old way)

### **3. Let Claude Code Choose Models**
Don't specify a model unless you have a preference. Claude Code is smart about picking the right one.

### **4. Include Context Files**
```
Use Delegate to update this API based on the requirements in docs/spec.md and existing code in src/api/
```

### **5. Perfect for Large Tasks**
- Generating entire applications
- Large refactoring tasks  
- Creating comprehensive test suites
- Writing extensive documentation
- **Analyzing multiple large documents**
- **Processing entire codebases**

### **6. Know Your Models**

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
Use Delegate to analyze architecture.md, api-spec.md, database-design.md and write an auth strategy summary to docs/auth-analysis.md
```
(Claude Code stays fresh for actual work!)

### **Real-World Scenarios**
- **"Analyze all test files and write a coverage report to docs/testing.md"**
- **"Review these RFC documents and extract key decisions to decisions.md"**  
- **"Read the entire codebase and write a refactoring plan to refactor-plan.md"**
- **"Process these log files and write an error analysis to logs/analysis.txt"**

## **🔧 New Structured Responses**

Everything is now predictable JSON (no more string parsing!):

### **Task Submission**
```json
{
  "output_id": "out_20250620_143022",
  "working_directory": "/home/user/project"
}
```

### **Metadata Check**
```json
{
  "metadata": {
    "status": "COMPLETED",
    "size_kb": 15.7,
    "token_estimate": 3925,
    "is_truncated": false
  },
  "content_analysis": {
    "blocks_found": 2,
    "blocks": [
      {
        "index": 0,
        "language": "jsx", 
        "size_kb": 12.1,
        "lines": 250,
        "preview": "import React from 'react';"
      }
    ]
  }
}
```

### **File Write Success**
```json
{
  "success": true,
  "path": "src/component.jsx",
  "absolute_path": "/home/user/project/src/component.jsx", 
  "bytes_written": 12388,
  "message": "Successfully wrote 12.1 KB to src/component.jsx",
  "working_directory": "/home/user/project"
}
```

### **Structured Errors**
```json
{
  "error": {
    "code": "OUTPUT_NOT_FOUND",
    "message": "The requested output ID does not exist or has expired.",
    "details": {
      "output_id_provided": "invalid-id"
    }
  }
}
```

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

**Q: What's the difference between the old and new API?**
- **Old (still works):** 3 tools with confusing dual-purpose `delegate_read`
- **New (recommended):** 4 clear tools with structured JSON responses
- The old API will show deprecation warnings

**Q: Can I read outputs later?**
- Yes! Use the output ID with `get_output_metadata` and other tools
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
- Check the exact ID from the submit_task response

### **"Path errors"** 
- Use relative paths only: `src/file.go` not `/abs/path/to/src/file.go`
- Paths are relative to your project working directory

## **📚 Next Steps**

1. **Try the zero-token workflow** - Generate code and save it directly to files
2. **Read the Model Reference** - Understand when to use each model  
3. **Check out the API Reference** - Deep dive into structured responses
4. **Explore multi-block handling** - Master complex generations

## **🎉 You're Ready!**

Start saving tokens and generating more code with Delegate. Remember the simple workflow:

1. **submit_task** → get an output_id
2. **write_output_to_file** → save directly (zero tokens!)

Or for complex cases:
1. **submit_task** → get an output_id
2. **get_output_metadata** → see what was generated  
3. **write_output_to_file** or **get_output_content** → save or read

Happy coding! 🚀

---

*Built with ❤️ and a strict No Scope Creep policy*