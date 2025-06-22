# 📚 Task Agent + Delegate Prompt Patterns Cookbook

*Copy → Paste → Customize → Execute → Marvel at the results*

## Quick Navigation
- [Single File Generation](#single-file-generation)
- [Multi-File Features](#multi-file-features)
- [Analysis & Review](#analysis--review)
- [Refactoring & Migration](#refactoring--migration)
- [Testing & Quality](#testing--quality)
- [Documentation](#documentation)
- [Debugging & Fixes](#debugging--fixes)
- [Advanced Orchestration](#advanced-orchestration)

---

## Single File Generation

### Pattern: Create New File from Scratch
```
Use a Task agent to create [FILENAME]. The agent should:
1. Research similar files in the codebase for patterns and conventions
2. Use delegate_invoke with model "gemini-2.5-flash" to generate the file
3. Include proper imports, types, and follow existing code style
4. Save using delegate_read with write_to option
5. Verify the file compiles/has no syntax errors
6. Report back with confirmation and any issues fixed
```

**Example:**
```
Use a Task agent to create middleware/rateLimiter.ts. The agent should:
1. Research similar files in the codebase for patterns and conventions
2. Use delegate_invoke with model "gemini-2.5-flash" to generate the file
3. Include proper imports, types, and follow existing code style
4. Save using delegate_read with write_to option
5. Verify the file compiles/has no syntax errors
6. Report back with confirmation and any issues fixed
```

### Pattern: Enhance Existing File
```
Task agent: Enhance [FILEPATH] by adding [FEATURE]. Steps:
1. Read the current file to understand its structure
2. Research related files that might be affected
3. Use delegate with the file as context to generate enhanced version
4. The enhancement should include [SPECIFIC REQUIREMENTS]
5. Save the updated file and verify it still works
6. Update any imports in other files if needed
```

---

## Multi-File Features

### Pattern: Complete Feature Implementation
```
Use Task agents to implement [FEATURE NAME]:
1. First, research the codebase structure and identify where each part belongs
2. Create a plan listing all files that need to be created/modified
3. For each new file, use delegate_invoke to generate it
4. For modifications, read existing file then use delegate to update
5. Ensure all files work together (imports, types match, etc.)
6. Run any build commands to verify no errors
7. Report: List all files created/modified with brief description of each
```

**Example:**
```
Use Task agents to implement user notifications system:
1. First, research the codebase structure and identify where each part belongs
2. Create a plan listing all files that need to be created/modified
3. For each new file, use delegate_invoke to generate it
4. For modifications, read existing file then use delegate to update
5. Ensure all files work together (imports, types match, etc.)
6. Run any build commands to verify no errors
7. Report: List all files created/modified with brief description of each
```

### Pattern: Parallel Component Generation
```
Launch 3 Task agents in parallel to create [SYSTEM]:
- Agent 1: Create the data layer (models, database queries, types)
- Agent 2: Create the business logic layer (services, utilities) 
- Agent 3: Create the presentation layer (UI components, styles)

Each agent should:
- Use delegate for all file generation
- Follow existing patterns in their layer
- Ensure their outputs will integrate with the other layers
- Save all files directly with write_to
- Report what they created
```

---

## Analysis & Review

### Pattern: Comprehensive Code Review
```
Use a Task agent to perform a deep code review:
1. Identify all source files in [DIRECTORY/FEATURE]
2. Create a bundle of these files (or use existing bundle.md)
3. Use delegate_invoke with model "gemini-2.5-pro" to analyze:
   - Code quality and best practices
   - Security vulnerabilities
   - Performance concerns
   - Architectural issues
   - Suggested improvements
4. Save the review as docs/reviews/[FEATURE]-review-[DATE].md
5. Create a prioritized action items list
```

### Pattern: Architecture Analysis
```
Task agent mission - Analyze and document our architecture:
1. Scan the entire codebase structure
2. Identify key patterns, frameworks, and dependencies
3. Use delegate with Gemini Pro to create:
   - Architecture overview document
   - Component interaction diagrams (mermaid/ASCII)
   - Technology stack assessment
   - Recommendations for improvements
4. Save all documents in docs/architecture/
5. Report key findings and concerns
```

---

## Refactoring & Migration

### Pattern: Technology Migration
```
Deploy a Task agent to migrate from [OLD TECH] to [NEW TECH]:
1. Find all files using [OLD TECH]
2. Create a migration plan with phases
3. For each file:
   - Read current implementation
   - Use delegate to generate migrated version
   - Ensure new version maintains same functionality
   - Update imports and dependencies
4. Update package.json/config files as needed
5. Run tests to ensure nothing broke
6. Report: Migration complete with X files updated
```

**Example:**
```
Deploy a Task agent to migrate from JavaScript to TypeScript:
1. Find all .js files in the src/ directory
2. Create a migration plan with phases
3. For each file:
   - Read current implementation  
   - Use delegate to generate TypeScript version
   - Ensure types are properly defined
   - Update imports to use .ts extensions
4. Create/update tsconfig.json
5. Run tsc to check for type errors
6. Report: Migration complete with X files converted
```

### Pattern: Pattern Refactoring
```
Task agent: Refactor codebase to use [PATTERN]:
1. Identify files that could benefit from [PATTERN]
2. Create examples of the pattern applied to our code
3. Use delegate to refactor each identified file
4. Ensure backward compatibility
5. Update tests to cover new pattern
6. Document the pattern for team reference
```

---

## Testing & Quality

### Pattern: Test Suite Generation
```
Use a Task agent to create comprehensive tests for [FEATURE/FILE]:
1. Analyze the code to understand all functions/methods
2. Identify edge cases and error scenarios
3. Use delegate to generate test file(s) with:
   - Unit tests for each function
   - Integration tests for feature flows
   - Edge case coverage
   - Mock data and fixtures
4. Ensure tests actually run and pass
5. Aim for >80% code coverage
6. Report coverage stats and any uncovered lines
```

### Pattern: E2E Test Creation
```
Task agent: Create end-to-end tests for [USER FLOW]:
1. Understand the complete user journey
2. Identify all interaction points and validations
3. Use delegate to create E2E tests using [FRAMEWORK]
4. Include:
   - Happy path scenarios
   - Error scenarios  
   - Data validation checks
   - UI state verifications
5. Run tests and fix any failures
6. Save tests and report results
```

---

## Documentation

### Pattern: API Documentation
```
Deploy documentation agent:
1. Scan all API endpoints/routes
2. Extract method, path, parameters, responses
3. Use delegate with Gemini Pro to generate:
   - OpenAPI/Swagger specification
   - Markdown API reference
   - Postman collection
   - Example curl commands
4. Include error codes and rate limits
5. Save in docs/api/
6. Generate a simple HTML preview
```

### Pattern: README Generation
```
Task agent: Create/update README.md for [PROJECT/FEATURE]:
1. Analyze the codebase structure and purpose
2. Use delegate to generate README with:
   - Clear description and purpose
   - Installation instructions
   - Usage examples with code
   - API reference (if applicable)
   - Configuration options
   - Contributing guidelines
   - License information
3. Ensure all examples actually work
4. Include badges for build status, coverage, etc.
```

---

## Debugging & Fixes

### Pattern: Error Investigation and Fix
```
Task agent: Debug and fix [ERROR MESSAGE/ISSUE]:
1. Search for the error message in the codebase
2. Identify the root cause and affected files
3. Research similar issues and solutions
4. Use delegate to generate the fix
5. Apply the fix and verify it resolves the issue
6. Check for any side effects
7. Update tests if needed
8. Report: Root cause, fix applied, files changed
```

### Pattern: Performance Optimization
```
Use a Task agent to optimize [FEATURE/FILE] performance:
1. Analyze current implementation for bottlenecks
2. Research best practices for optimization
3. Use delegate to generate optimized version with:
   - Algorithm improvements
   - Caching strategies
   - Query optimization
   - Lazy loading
   - Memoization where appropriate
4. Benchmark before/after if possible
5. Ensure functionality unchanged
6. Document optimization choices
```

---

## Advanced Orchestration

### Pattern: Multi-Agent Project Bootstrap
```
Initialize a new [PROJECT TYPE] project using agent team:

Launch 4 specialized agents:
- Agent 1 (Architect): Design structure and create boilerplate
- Agent 2 (Backend): Implement API and database layer
- Agent 3 (Frontend): Create UI components and pages
- Agent 4 (DevOps): Setup build, test, deployment configs

Each agent uses delegate for their domain and coordinates through 
shared files. Agents should read each other's outputs for context.

Final report: Complete project ready to run with all layers integrated.
```

### Pattern: Continuous Improvement Loop
```
Deploy an improvement agent that:
1. Analyzes codebase for quality metrics
2. Identifies top 3 areas for improvement
3. For each area:
   - Use delegate to generate improved version
   - Run tests to ensure nothing breaks
   - Measure improvement (performance/readability/etc)
4. If improvement is significant, keep changes
5. Document what was improved and why
6. Repeat until no significant improvements found
```

### Pattern: Cross-Codebase Integration
```
Task agent: Integrate [EXTERNAL SERVICE/API] into our codebase:
1. Research the external service documentation
2. Identify integration points in our code
3. Use delegate to generate:
   - Client/SDK wrapper
   - Type definitions
   - Service integration layer
   - Error handling
   - Tests with mocked responses
4. Create example usage in our app
5. Document configuration and setup
6. Ensure graceful fallbacks for service unavailability
```

---

## Meta Patterns

### Pattern: Agent That Creates Agents
```
Use a Task agent to analyze [LARGE TASK] and create a multi-agent plan:
1. Break down the task into parallel workstreams
2. Define what each agent should do
3. Identify dependencies between agents
4. Generate specific prompts for each agent
5. Create a coordination strategy
6. Return the complete execution plan
Then execute the plan with: "Run the multi-agent plan above"
```

### Pattern: Learning From Execution
```
Task agent: Implement [FEATURE] and document the process:
1. Keep a log of all decisions made
2. Track what worked and what didn't
3. Use delegate as normal for implementation
4. After completion, use delegate to create:
   - Implementation guide for similar features
   - Lessons learned document
   - Reusable patterns discovered
   - Template for future use
```

---

## Usage Tips

### Variables to Replace:
- `[FILENAME]` - Specific file path like `services/emailService.ts`
- `[FEATURE]` - Feature name like "password reset flow"
- `[DIRECTORY]` - Directory path like `src/components/`
- `[OLD TECH]` / `[NEW TECH]` - Technologies like "Redux" / "Zustand"
- `[PATTERN]` - Design pattern like "Repository Pattern"
- `[PROJECT TYPE]` - Project type like "Next.js blog"
- `[FRAMEWORK]` - Test framework like "Playwright" or "Cypress"

### Model Selection Guide:
- **gemini-2.5-flash**: Quick iterations, standard code generation
- **gemini-2.5-pro**: Complex analysis, architecture decisions
- **claude-sonnet-4**: Nuanced logic, complex algorithms
- **claude-opus-4**: When you need the absolute best quality

### Success Multipliers:
1. Always include "research existing patterns first"
2. Specify "ensure it compiles/builds/runs"
3. Add "fix any issues before reporting"
4. Include "follow our coding conventions"
5. Request "report what was done and why"

---

*"The best code is the code you didn't have to write yourself." - Every developer using this cookbook*

**Copy. Paste. Customize. Conquer. 🚀**