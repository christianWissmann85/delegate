# **Delegate Model Reference Card v1.0**

**Status:** Final | **Version:** 1.0 | **Date:** 2025-06-20

## **Quick Reference**

| Model | ID | Speed | Quality | Cost | Best For |
|-------|----|----|---------|------|----------|
| **Gemini 2.5 Flash** | `gemini-2.5-flash` | ⚡⚡⚡ | ⭐⭐⭐ | 💰 | Everyday coding, fast iterations |
| **Gemini 2.5 Pro** | `gemini-2.5-pro` | ⚡⚡ | ⭐⭐⭐⭐ | 💰💰 | Complex architecture, deep reasoning |
| **Claude Sonnet 4** | `claude-sonnet-4-20250514` | ⚡⚡ | ⭐⭐⭐⭐ | 💰💰 | Balanced tasks, precise instructions |
| **Claude Opus 4** | `claude-opus-4-20250514` | ⚡ | ⭐⭐⭐⭐⭐ | 💰💰💰 | Critical code, highest quality |

## **Model Details**

### **🚀 Gemini 2.5 Flash**
- **Model ID:** `gemini-2.5-flash`
- **Context Window:** 1 million tokens
- **Strengths:** Lightning fast, huge context, cost-effective
- **Use When:** 
  - Generating boilerplate code
  - Quick refactoring tasks
  - API endpoints and CRUD operations
  - Data transformations
- **Example Prompt:** "Create REST endpoints for a todo app"

### **🧠 Gemini 2.5 Pro**
- **Model ID:** `gemini-2.5-pro`
- **Context Window:** 1 million tokens (2M coming soon)
- **Strengths:** Advanced reasoning, complex problem solving
- **Use When:**
  - System architecture design
  - Complex algorithm implementation
  - Multi-file refactoring
  - Performance optimization
- **Example Prompt:** "Design a scalable event-driven microservices architecture"

### **⚖️ Claude Sonnet 4**
- **Model ID:** `claude-sonnet-4-20250514`
- **Context Window:** 200,000 tokens
- **Strengths:** Precise instruction following, great at debugging
- **Use When:**
  - Following detailed specifications
  - Debugging complex issues
  - Writing tests and documentation
  - Code review and improvements
- **Example Prompt:** "Implement this feature following our strict coding standards"

### **👑 Claude Opus 4**
- **Model ID:** `claude-opus-4-20250514`
- **Context Window:** 200,000 tokens
- **Strengths:** World's best coding model, exceptional quality
- **Use When:**
  - Security-critical implementations
  - Complex business logic
  - Mission-critical systems
  - Code that needs to be perfect first time
- **Example Prompt:** "Implement a secure authentication system with OAuth2 and MFA"

## **Decision Matrix**

### **By Task Type**

| Task | Recommended Model | Why |
|------|-------------------|-----|
| Quick scripts | Gemini 2.5 Flash | Speed matters most |
| API development | Gemini 2.5 Flash | Standard patterns, fast delivery |
| Complex algorithms | Gemini 2.5 Pro | Needs deep reasoning |
| Security code | Claude Opus 4 | Cannot afford mistakes |
| Following specs | Claude Sonnet 4 | Best instruction adherence |
| Large codebase work | Gemini 2.5 Pro | 1M token context window |
| Production features | Claude Sonnet 4 | Good balance of quality/speed |
| Architectural design | Claude Opus 4 | Highest intelligence needed |

### **By Context Size**

| File Count | File Sizes | Best Model |
|------------|------------|------------|
| 1-5 files | Small (<1KB each) | Any model |
| 5-20 files | Medium (<10KB each) | Any model |
| 20+ files | Large | Gemini models (1M context) |
| Entire codebases | Any size | Gemini models only |

## **Cost Optimization Tips**

1. **Start with Flash:** Try Gemini 2.5 Flash first - it's often good enough
2. **Escalate when needed:** Only use Opus 4 for truly complex tasks
3. **Use `check` before `read`:** Always check output size to avoid token waste
4. **Extract strategically:** Use `extract: "code"` to skip explanations

## **Model Selection in Practice**

```javascript
// Let Claude Code choose (recommended)
"Generate a user authentication system"

// Explicitly request a model
"Use Gemini Flash to create a simple CRUD API"
"Use Claude Opus to implement the payment processing logic"

// Context-aware selection
"This needs to be production-perfect, use Opus 4"
"Need this fast for a prototype, use Flash"
```

## **Provider Limits**

| Model | Rate Limit | Max Output | Timeout |
|-------|------------|------------|---------|
| Gemini 2.5 Flash | 60 RPM | 8K tokens | 60s |
| Gemini 2.5 Pro | 30 RPM | 8K tokens | 60s |
| Claude Sonnet 4 | 50 RPM | 4K tokens | 60s |
| Claude Opus 4 | 20 RPM | 4K tokens | 60s |

*RPM = Requests Per Minute, varies by API tier*

---

**Remember:** Claude Code is smart about model selection. When in doubt, describe your needs and let it choose!