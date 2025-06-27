# Delegate MCP Server Migration Guide

**From 3-Tool API to 4-Tool API | Version 2.0**

---

## Overview

This guide helps you migrate from the legacy 3-tool Delegate API to the new, more intuitive 4-tool architecture. The new API provides clearer separation of concerns, predictable token usage, and structured JSON responses.

**Migration Timeline:**
- **Phase 1 (Current):** Both APIs available
- **Phase 2 (v2.1):** Old API deprecated with warnings
- **Phase 3 (v3.0):** Old API removed

## Tool Mapping Table

| Old Tool | New Tools | Key Changes |
|----------|-----------|-------------|
| `delegate_invoke` | `delegate_submit_task` | ✅ Same functionality, cleaner response |
| `delegate_check` | `delegate_get_output_metadata` | ✅ Enhanced with structured content analysis |
| `delegate_read` | `delegate_get_output_content` OR `delegate_write_output_to_file` | ⚠️ **SPLIT**: One tool became two based on intent |

### Detailed Mapping

#### `delegate_invoke` → `delegate_submit_task`
- **Parameters:** Identical
- **Response:** Simplified to just `{"output_id": "...", "working_directory": "..."}`
- **Paths:** Now uses relative paths instead of absolute

#### `delegate_check` → `delegate_get_output_metadata`  
- **Parameters:** Identical
- **Response:** Enhanced with `content_analysis` block containing structured information about multi-block outputs
- **Benefit:** No more string parsing for multi-block warnings

#### `delegate_read` → Two Separate Tools

**For reading content into context:**
- Use `delegate_get_output_content`
- Same parameters as old `delegate_read` (minus `write_to`)

**For writing content to files:**
- Use `delegate_write_output_to_file` 
- Requires `write_to` parameter, same options as before

## Side-by-Side Examples

### Example 1: Code Generation to File

#### Old API (3-tool)
```javascript
// Step 1: Submit task
const task = await mcp.delegate_invoke({
  model: "gemini-2.5-flash",
  prompt: "Create a React component",
  files: ["/home/user/project/src/types.ts"]
});

// Step 2: Write to file (confusing dual-purpose tool)
const result = await mcp.delegate_read({
  output_id: task.output_id,
  options: {
    extract: "code",
    write_to: "/home/user/project/src/Component.jsx"
  }
});

// Result contains mixed success/content data
if (result.content) {
  // This shouldn't happen in write mode, but API was unclear
}
```

#### New API (4-tool)
```javascript
// Step 1: Submit task
const task = await mcp.delegate_submit_task({
  model: "gemini-2.5-flash", 
  prompt: "Create a React component",
  files: ["src/types.ts"]  // ✅ Relative paths
});

// Step 2: Write to file (clear, single-purpose tool)
const result = await mcp.delegate_write_output_to_file({
  output_id: task.output_id,
  write_to: "src/Component.jsx",  // ✅ Relative path
  options: {
    extract: "code"
  }
});

// ✅ Clear, structured response
console.log(`Wrote ${result.bytes_written} bytes to ${result.path}`);
```

### Example 2: Content Review Workflow

#### Old API (3-tool)
```javascript
// Submit task
const task = await mcp.delegate_invoke({
  model: "claude-sonnet-4-20250514",
  prompt: "Analyze this code",
  files: ["/home/user/project/src/main.go"]
});

// Check output - basic metadata only
const check = await mcp.delegate_check({
  output_id: task.output_id
});

// Read content and parse string warnings manually
const content = await mcp.delegate_read({
  output_id: task.output_id,
  options: { extract: "all" }
});

// ❌ Manual string parsing required
if (typeof content === 'string' && content.includes("WARNING: Multiple blocks")) {
  // Complex string parsing logic here...
  const blocks = parseMultiBlockWarning(content);
  // Handle each block...
}
```

#### New API (4-tool)
```javascript
// Submit task
const task = await mcp.delegate_submit_task({
  model: "claude-sonnet-4-20250514",
  prompt: "Analyze this code", 
  files: ["src/main.go"]  // ✅ Relative path
});

// Check output - rich structured metadata
const metadata = await mcp.delegate_get_output_metadata({
  output_id: task.output_id
});

// ✅ Structured data - no parsing needed
if (metadata.content_analysis.blocks_found > 1) {
  // Handle each block programmatically
  for (const block of metadata.content_analysis.blocks) {
    console.log(`Block ${block.index}: ${block.language}, ${block.lines} lines`);
    
    if (block.language === "go") {
      await mcp.delegate_write_output_to_file({
        output_id: task.output_id,
        write_to: `output/analysis_${block.index}.go`,
        options: { block_index: block.index }
      });
    }
  }
} else {
  // Single block - get content normally
  const content = await mcp.delegate_get_output_content({
    output_id: task.output_id
  });
  console.log(content.content);
}
```

### Example 3: Large Output Handling

#### Old API (3-tool)
```javascript
// Submit task
const task = await mcp.delegate_invoke({
  model: "gemini-2.5-pro",
  prompt: "Generate comprehensive documentation"
});

// Check basic info
const check = await mcp.delegate_check({
  output_id: task.output_id
});

// ❌ Limited metadata, hard to make smart decisions
if (check.status === "COMPLETED") {
  // Unclear what the token cost will be
  const content = await mcp.delegate_read({
    output_id: task.output_id
  });
  // Might consume thousands of unexpected tokens
}
```

#### New API (4-tool)
```javascript
// Submit task
const task = await mcp.delegate_submit_task({
  model: "gemini-2.5-pro",
  prompt: "Generate comprehensive documentation"
});

// Get rich metadata for smart decisions
const metadata = await mcp.delegate_get_output_metadata({
  output_id: task.output_id
});

// ✅ Make informed decisions based on structured data
console.log(`Output size: ${metadata.metadata.size_kb} KB`);
console.log(`Estimated tokens: ${metadata.metadata.token_estimate}`);

if (metadata.metadata.token_estimate > 5000) {
  // Large output - write directly to file (ZERO tokens)
  await mcp.delegate_write_output_to_file({
    output_id: task.output_id,
    write_to: "docs/generated.md"
  });
  console.log("Saved large output to file to preserve tokens");
} else {
  // Small output - safe to get content
  const content = await mcp.delegate_get_output_content({
    output_id: task.output_id,
    options: { max_tokens: 2000 }  // ✅ Explicit limit
  });
  console.log(content.content);
}
```

## Key Differences to Watch Out For

### 1. Path Handling: Absolute → Relative

**Old API:**
```javascript
// Required full absolute paths
files: ["/home/user/myproject/src/main.go", "/home/user/myproject/docs/api.md"]
write_to: "/home/user/myproject/output/result.txt"
```

**New API:**
```javascript
// Uses relative paths from working directory
files: ["src/main.go", "docs/api.md"]
write_to: "output/result.txt"

// Working directory provided in responses for context
// { "working_directory": "/home/user/myproject" }
```

### 2. Response Structure: Strings → JSON

**Old API:**
```javascript
// Mixed string/object responses, required parsing
const result = await mcp.delegate_read({...});

// Could be a string like:
// "WARNING: Multiple blocks detected. Block 0: JavaScript (245 lines)..."
// OR an object with content

if (typeof result === 'string' && result.includes('WARNING')) {
  // Manual parsing required
}
```

**New API:**
```javascript
// Always structured JSON
const metadata = await mcp.delegate_get_output_metadata({...});

// Reliable object structure:
metadata.content_analysis.blocks.forEach(block => {
  console.log(`${block.language}: ${block.lines} lines`);
});
```

### 3. Tool Purpose: Dual → Single

**Old API:**
```javascript
// delegate_read did two completely different things
const content = await mcp.delegate_read({
  output_id: "123",
  // No write_to = returns content (costs tokens)
});

const writeResult = await mcp.delegate_read({
  output_id: "123", 
  options: { write_to: "/path/file.txt" }  // Writes file (costs zero tokens)
});
```

**New API:**
```javascript
// Two separate tools with clear purposes
const content = await mcp.delegate_get_output_content({
  output_id: "123"
  // Always returns content, always costs tokens
});

const writeResult = await mcp.delegate_write_output_to_file({
  output_id: "123",
  write_to: "file.txt"
  // Always writes file, always costs zero tokens
});
```

### 4. Error Handling: Strings → Structured

**Old API:**
```javascript
// Errors were often mixed with content or in string format
try {
  const result = await mcp.delegate_read({...});
} catch (error) {
  // Error format was inconsistent
  console.log(error.message); // Varied format
}
```

**New API:**
```javascript
// Consistent structured error format
try {
  const result = await mcp.delegate_get_output_content({...});
} catch (error) {
  // Reliable structure:
  console.log(`Error ${error.code}: ${error.message}`);
  if (error.details) {
    console.log('Details:', error.details);
  }
}
```

## Migration Strategy

### Step 1: Update Path Handling
```javascript
// Change absolute paths to relative
const oldFiles = ["/home/user/project/src/main.go"];
const newFiles = ["src/main.go"];

const oldWritePath = "/home/user/project/output/result.txt";
const newWritePath = "output/result.txt";
```

### Step 2: Split delegate_read Usage
```javascript
// Identify your delegate_read calls
const oldRead = await mcp.delegate_read({
  output_id: "123",
  options: { write_to: "/path/file.txt" }  // ← Writing to file
});

// Becomes:
const newWrite = await mcp.delegate_write_output_to_file({
  output_id: "123",
  write_to: "path/file.txt"  // ← Use new tool
});

// ---

const oldRead2 = await mcp.delegate_read({
  output_id: "123",
  options: { extract: "code" }  // ← Getting content (no write_to)
});

// Becomes:
const newContent = await mcp.delegate_get_output_content({
  output_id: "123", 
  options: { extract: "code" }  // ← Use new tool
});
```

### Step 3: Replace String Parsing with Structured Data
```javascript
// Old manual parsing
if (typeof result === 'string' && result.includes('Multiple blocks')) {
  const blockMatches = result.match(/Block (\d+): (\w+) \((\d+) lines\)/g);
  // Complex regex parsing...
}

// New structured access
const metadata = await mcp.delegate_get_output_metadata({ output_id });
metadata.content_analysis.blocks.forEach(block => {
  console.log(`Block ${block.index}: ${block.language} (${block.lines} lines)`);
});
```

### Step 4: Update Tool Names
```javascript
// Simple renames
delegate_invoke → delegate_submit_task
delegate_check → delegate_get_output_metadata

// delegate_read splits based on usage:
// - Reading content → delegate_get_output_content  
// - Writing files → delegate_write_output_to_file
```

## FAQ

### Q: Can I use both APIs during migration?
**A:** Yes! Both APIs are available until v3.0. The old API will show deprecation warnings starting in v2.1.

### Q: Will my existing output_ids work with new tools?
**A:** Yes, output_ids are compatible across both APIs during the migration period.

### Q: How do I know if I'm using delegate_read for reading vs writing?
**A:** Check for the `write_to` parameter in options:
- **Has `write_to`** → Use `delegate_write_output_to_file`
- **No `write_to`** → Use `delegate_get_output_content`

### Q: What about working directory differences?
**A:** The new API shows `working_directory` in responses. Make sure your relative paths resolve correctly from that directory.

### Q: Are there performance benefits to the new API?
**A:** Yes! The structured metadata allows smarter token usage decisions, and the path handling is more efficient.

### Q: How do I handle the new error format?
**A:** Replace string-based error handling with structured error codes:

```javascript
// Old
if (error.message.includes("not found")) { ... }

// New  
if (error.code === "OUTPUT_NOT_FOUND") { ... }
```

### Q: What if I need both content AND file writing?
**A:** Use the metadata tool to make smart decisions:

```javascript
const metadata = await mcp.delegate_get_output_metadata({ output_id });

if (metadata.metadata.token_estimate < 1000) {
  // Small - get content AND write file
  const content = await mcp.delegate_get_output_content({ output_id });
  await mcp.delegate_write_output_to_file({ output_id, write_to: "backup.txt" });
} else {
  // Large - just write file to save tokens
  await mcp.delegate_write_output_to_file({ output_id, write_to: "output.txt" });
}
```

### Q: When should I use the metadata tool?
**A:** Use `delegate_get_output_metadata` when you need to:
- Make decisions about token usage
- Handle multi-block outputs intelligently  
- Check output size before processing
- Understand content structure before extraction

### Q: Can I still use absolute paths?
**A:** No, the new API only accepts relative paths. This reduces errors and improves portability.

## Deprecation Timeline

### Phase 1: Coexistence (Current - v2.0)
- ✅ Both old and new APIs available
- ✅ Full compatibility for existing code
- ✅ New projects should use 4-tool API
- ✅ Documentation focuses on new API

### Phase 2: Deprecation (v2.1)
- ⚠️ Old API returns deprecation warnings in responses
- ⚠️ Old API endpoints marked as deprecated in documentation
- ✅ Old API remains fully functional
- 📅 Migration deadline announced

### Phase 3: Removal (v3.0)
- ❌ Old API completely removed
- ❌ `delegate_invoke`, `delegate_check`, `delegate_read` no longer available
- ✅ Only 4-tool API supported
- 📅 Estimated timeline: 6+ months after v2.1 release

### Migration Checkpoints

**Before v2.1 (Recommended):**
- [ ] Update all file paths to relative format
- [ ] Replace `delegate_read` calls with appropriate new tools
- [ ] Update error handling to use structured format
- [ ] Test workflows with new API

**Before v3.0 (Required):**  
- [ ] Remove all old API usage
- [ ] Update any hardcoded tool names
- [ ] Verify no deprecation warnings in logs
- [ ] Complete migration testing

---

**Need Help?** 
- Check the [API Reference](api-reference.md) for detailed documentation
- Review the workflow examples for common patterns
- Test the new API alongside your existing code during Phase 1