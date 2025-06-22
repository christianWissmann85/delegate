# 🎯 Task Agent + Delegate Workflow Examples

*Real-world scenarios with step-by-step execution guides*

## Table of Contents
1. [Workflow 1: Building a REST API from Scratch](#workflow-1-building-a-rest-api-from-scratch)
2. [Workflow 2: Adding Authentication to Existing App](#workflow-2-adding-authentication-to-existing-app)
3. [Workflow 3: Codebase Modernization Project](#workflow-3-codebase-modernization-project)
4. [Workflow 4: Emergency Bug Fix](#workflow-4-emergency-bug-fix)
5. [Workflow 5: Feature Development with TDD](#workflow-5-feature-development-with-tdd)
6. [Workflow 6: Multi-Service Integration](#workflow-6-multi-service-integration)
7. [Workflow 7: Performance Optimization Sprint](#workflow-7-performance-optimization-sprint)
8. [Workflow 8: Documentation Overhaul](#workflow-8-documentation-overhaul)

---

## Workflow 1: Building a REST API from Scratch

**Scenario**: You need a complete REST API for a blog platform with posts, comments, and users.

### Step 1: Initial Planning
```markdown
You: "I need a REST API for a blog platform with posts, comments, and users. 
     Use TypeScript, Express, and PostgreSQL."

Claude: "I'll orchestrate a team of agents to build this. Let me break it down..."
```

### Step 2: Architecture Agent
```markdown
Claude: "First, let me deploy an architecture agent to design the structure."

Use a Task agent to design the blog API architecture:
1. Research best practices for Express + TypeScript + PostgreSQL apps
2. Use delegate to create:
   - Project structure documentation
   - Database schema (SQL)
   - API endpoint specification
   - Type definitions for all entities
3. Save all docs in docs/architecture/
4. Report the proposed structure
```

**Agent Output**:
```
✅ Created docs/architecture/structure.md
✅ Created docs/architecture/schema.sql
✅ Created docs/architecture/api-spec.md
✅ Created types/entities.ts

Architecture uses layered approach:
- Routes → Controllers → Services → Repositories
- Strong typing throughout
- JWT authentication ready
```

### Step 3: Parallel Implementation
```markdown
Claude: "Now I'll deploy three specialized agents in parallel."

Launch three Task agents simultaneously:

Agent 1 - Database Layer:
- Create all model files using delegate
- Generate migration files
- Create repository classes
- Set up database connection

Agent 2 - Business Logic:
- Create service classes using delegate
- Implement validation logic
- Create utility functions
- Handle business rules

Agent 3 - API Layer:
- Create route files using delegate
- Generate controllers
- Set up middleware
- Create error handlers

Each agent should use the architecture docs as context.
```

### Step 4: Integration Agent
```markdown
Claude: "Time to wire everything together."

Deploy integration agent to:
1. Create the main server.ts file
2. Wire up all routes
3. Configure middleware pipeline
4. Set up error handling
5. Create package.json with all dependencies
6. Generate .env.example
7. Create docker-compose.yml for PostgreSQL
8. Verify everything compiles
```

### Step 5: Testing Agent
```markdown
Claude: "Let's ensure quality with comprehensive tests."

Task agent: Create test suite for the blog API:
1. Use delegate to generate tests for each endpoint
2. Include:
   - Unit tests for services
   - Integration tests for API routes
   - Database repository tests
3. Create test fixtures and factories
4. Set up test database configuration
5. Ensure all tests pass
```

### Step 6: Final Polish
```markdown
Claude: "Final touches from a DevOps agent."

Deploy DevOps agent to:
1. Create README.md with setup instructions
2. Add GitHub Actions workflow
3. Create Dockerfile
4. Generate API documentation
5. Add logging configuration
6. Create deployment guide
```

### Final Result
```
📁 blog-api/
├── 📁 src/
│   ├── 📁 routes/      (5 files)
│   ├── 📁 controllers/ (5 files)
│   ├── 📁 services/    (5 files)
│   ├── 📁 repositories/(5 files)
│   ├── 📁 models/      (3 files)
│   ├── 📁 middleware/  (4 files)
│   ├── 📁 utils/       (3 files)
│   └── server.ts
├── 📁 tests/           (15 files)
├── 📁 docs/            (6 files)
├── 📁 types/           (4 files)
├── package.json
├── tsconfig.json
├── Dockerfile
├── docker-compose.yml
└── README.md

Total: ~50 files generated in ~5 minutes
Token cost: <1000 (vs ~100,000 traditional)
```

---

## Workflow 2: Adding Authentication to Existing App

**Scenario**: Your Express app needs JWT authentication with refresh tokens.

### Step 1: Analysis
```markdown
You: "Add JWT authentication with refresh tokens to our Express app."

Claude: "I'll analyze your current setup first."

Task agent: Analyze the existing app for auth integration:
1. Scan routes to identify protected endpoints
2. Check existing user model/database
3. Identify middleware pipeline
4. Use delegate with findings to create integration plan
5. Report current state and proposed changes
```

### Step 2: Implementation
```markdown
Claude: "Based on analysis, I'll implement auth systematically."

Deploy auth implementation agent:
1. Create auth service with delegate:
   - JWT token generation
   - Refresh token logic
   - Password hashing (bcrypt)
2. Create auth middleware
3. Create auth routes:
   - POST /auth/register
   - POST /auth/login
   - POST /auth/refresh
   - POST /auth/logout
4. Update user model if needed
5. Add auth to existing routes
```

### Step 3: Security Hardening
```markdown
Claude: "Let's ensure this is production-ready."

Security agent tasks:
1. Add rate limiting to auth endpoints
2. Implement refresh token rotation
3. Add CORS configuration
4. Set secure cookie options
5. Create security headers middleware
6. Add input validation with Joi/Zod
7. Document security considerations
```

### Step 4: Testing & Documentation
```markdown
Claude: "Comprehensive tests and docs for the team."

Final agent: Complete the auth feature:
1. Generate auth tests with delegate
2. Create Postman collection
3. Update API documentation
4. Create auth flow diagrams
5. Write migration guide for existing users
```

---

## Workflow 3: Codebase Modernization Project

**Scenario**: Migrate a legacy JavaScript/callback-based app to modern TypeScript/async-await.

### Step 1: Assessment
```markdown
You: "We need to modernize our old Node.js app from callbacks to TypeScript + async/await."

Claude: "This needs careful orchestration. Starting with assessment."

Task agent: Assess the legacy codebase:
1. Count JS files and total lines
2. Identify callback patterns (cb, callback, done)
3. Find dependencies that need updating
4. Check for tests (migration safety net)
5. Use delegate to create migration report
6. Prioritize files by complexity/importance
```

### Step 2: Preparation
```markdown
Claude: "Setting up the modernization infrastructure."

Prep agent: Set up TypeScript environment:
1. Create tsconfig.json with strict settings
2. Update package.json with TS dependencies
3. Create type definition files for legacy code
4. Set up build scripts
5. Configure ESLint for TS
6. Create migration tracking document
```

### Step 3: Parallel Migration Waves
```markdown
Claude: "I'll migrate in waves to maintain stability."

Wave 1 - Utility Files (3 agents):
- Agent 1: Convert utility/helper files
- Agent 2: Convert configuration files
- Agent 3: Convert constants/types

Each agent:
1. Read JS file
2. Use delegate to convert to TS + async/await
3. Fix any type errors
4. Update imports in other files
5. Verify functionality unchanged

[After Wave 1 succeeds]

Wave 2 - Core Business Logic (4 agents):
- Agent 1: Database layer
- Agent 2: Service layer
- Agent 3: API routes
- Agent 4: Middleware

Wave 3 - Tests and Documentation:
- Convert all tests to TS
- Update all documentation
- Create migration guide
```

### Step 4: Validation
```markdown
Claude: "Ensuring nothing broke during migration."

Validation agent:
1. Run all tests (old vs new)
2. Compare API responses
3. Check for any runtime errors
4. Verify performance didn't degrade
5. Create comparison report
6. List any breaking changes
```

---

## Workflow 4: Emergency Bug Fix

**Scenario**: Production is down! Users can't login - 500 errors everywhere.

### Step 1: Rapid Diagnosis
```markdown
You: "URGENT: Users can't login, getting 500 errors!"

Claude: "I'll immediately deploy a diagnostic agent."

Emergency diagnostic agent:
1. Check auth-related files for recent changes
2. Look for login/auth error handling
3. Search for database connection code
4. Check for environment variable usage
5. Scan logs patterns if available
6. Report findings immediately
```

### Step 2: Root Cause Analysis
```markdown
Claude: "Found potential issue in auth service. Investigating deeper."

Deep dive agent:
1. Read the problematic file(s)
2. Trace the login flow step-by-step
3. Identify where 500 error originates
4. Check related files for context
5. Use delegate to analyze and suggest fix
```

**Agent finds**: "JWT secret environment variable is undefined in production!"

### Step 3: Fix Implementation
```markdown
Claude: "Implementing emergency fix with safeguards."

Fix implementation agent:
1. Add fallback for missing JWT secret
2. Improve error messages
3. Add environment variable validation
4. Create health check endpoint
5. Add monitoring hooks
6. Test the fix thoroughly
```

### Step 4: Prevention
```markdown
Claude: "Let's prevent this from happening again."

Prevention agent:
1. Create env variable validator
2. Add startup checks
3. Generate .env.example with all vars
4. Create deployment checklist
5. Add integration tests for auth flow
6. Document the incident and fix
```

---

## Workflow 5: Feature Development with TDD

**Scenario**: Build a shopping cart feature using Test-Driven Development.

### Step 1: Test First
```markdown
You: "Build a shopping cart feature using TDD approach."

Claude: "Perfect! I'll start with tests, then implementation."

Test design agent:
1. Use delegate to create comprehensive cart tests:
   - Add item to cart
   - Remove item from cart
   - Update quantities
   - Calculate totals
   - Apply discounts
   - Handle out-of-stock
2. Create test fixtures
3. All tests should initially fail
```

### Step 2: Implementation Cycles
```markdown
Claude: "Now I'll implement features to make tests pass, one by one."

For each test group, deploy an agent to:
1. Read the failing test
2. Use delegate to implement minimal code to pass
3. Run the test to verify it passes
4. Refactor if needed
5. Move to next test
6. Report progress

Cycle 1: Basic cart operations (add/remove)
Cycle 2: Quantity management
Cycle 3: Price calculations
Cycle 4: Discount system
Cycle 5: Inventory integration
```

### Step 3: Integration
```markdown
Claude: "Integrating cart with the rest of the system."

Integration agent:
1. Create cart API endpoints
2. Add cart to user session/database
3. Integrate with product service
4. Add cart persistence
5. Create cart UI components (if needed)
6. Wire up to frontend
```

### Step 4: Polish
```markdown
Claude: "Final touches for production readiness."

Polish agent:
1. Add cart analytics events
2. Create admin cart viewer
3. Add cart abandonment tracking
4. Performance optimize large carts
5. Add cart-related emails
6. Document cart API
```

---

## Workflow 6: Multi-Service Integration

**Scenario**: Integrate Stripe, SendGrid, and Discord webhooks into your app.

### Step 1: Parallel Service Wrappers
```markdown
You: "Integrate Stripe payments, SendGrid emails, and Discord notifications."

Claude: "I'll create dedicated agents for each service."

Deploy 3 specialist agents in parallel:

Stripe Agent:
1. Create Stripe service wrapper with delegate
2. Implement: payments, subscriptions, webhooks
3. Add strong typing for all Stripe objects
4. Create test mode toggle
5. Add webhook signature verification

SendGrid Agent:
1. Create email service with delegate
2. Implement: transactional, bulk, templates
3. Create email queue system
4. Add bounce handling
5. Create test email viewer

Discord Agent:
1. Create Discord notification service
2. Implement: webhooks, embeds, mentions
3. Add rate limiting
4. Create notification templates
5. Add error channel
```

### Step 2: Orchestration Layer
```markdown
Claude: "Now to orchestrate these services."

Orchestration agent:
1. Create unified event system
2. Wire services to business events:
   - Purchase → Stripe + SendGrid + Discord
   - Refund → Stripe + SendGrid
   - User signup → SendGrid + Discord
3. Add circuit breakers
4. Create fallback strategies
5. Add comprehensive logging
```

### Step 3: Testing & Monitoring
```markdown
Claude: "Ensuring reliability of integrations."

Integration test agent:
1. Create integration tests with mocked services
2. Add webhook endpoint tests
3. Create service health checks
4. Add monitoring dashboards
5. Create runbooks for failures
6. Document all integrations
```

---

## Workflow 7: Performance Optimization Sprint

**Scenario**: App is slow, needs comprehensive performance optimization.

### Step 1: Profiling
```markdown
You: "The app is running slowly, optimize performance."

Claude: "I'll run a comprehensive performance audit."

Performance audit agent:
1. Identify all database queries
2. Find N+1 query problems
3. Check for missing indexes
4. Identify large JSON operations
5. Find synchronous operations that could be async
6. Use delegate to create optimization report
```

### Step 2: Database Optimization
```markdown
Claude: "Starting with the biggest impact - database."

DB optimization agent:
1. Add missing indexes based on query patterns
2. Implement query result caching
3. Convert N+1 queries to joins
4. Add database connection pooling
5. Implement read replicas if needed
6. Create query performance tests
```

### Step 3: Application Layer
```markdown
Claude: "Now optimizing application code."

App optimization agent:
1. Implement response caching
2. Add pagination to large lists
3. Optimize JSON serialization
4. Implement lazy loading
5. Add request compression
6. Optimize bundle sizes
```

### Step 4: Infrastructure
```markdown
Claude: "Final layer - infrastructure optimization."

Infra agent:
1. Configure CDN for static assets
2. Implement Redis for caching
3. Add job queues for heavy tasks
4. Configure auto-scaling rules
5. Optimize Docker images
6. Create performance monitoring
```

---

## Workflow 8: Documentation Overhaul

**Scenario**: Your project has grown but documentation is outdated and incomplete.

### Step 1: Documentation Audit
```markdown
You: "Our docs are a mess. Need complete overhaul."

Claude: "I'll audit and rebuild systematically."

Audit agent:
1. Catalog all existing documentation
2. Check which docs are outdated
3. Identify undocumented features
4. Survey code comments quality
5. List missing documentation types
6. Create documentation roadmap
```

### Step 2: Parallel Documentation Generation
```markdown
Claude: "Deploying specialized documentation agents."

Deploy 5 documentation agents:

API Docs Agent:
- Generate OpenAPI spec from routes
- Create interactive API explorer
- Generate client SDKs
- Create webhook documentation

Code Docs Agent:
- Generate JSDoc/TSDoc for all functions
- Create architecture diagrams
- Document design patterns used
- Create code style guide

User Docs Agent:
- Write getting started guide
- Create feature tutorials
- Build troubleshooting guide
- Generate FAQ

Dev Docs Agent:
- Create development setup guide
- Document deployment process
- Write testing guide
- Create contribution guidelines

Ops Docs Agent:
- Document infrastructure
- Create runbooks
- Write monitoring guide
- Document backup procedures
```

### Step 3: Documentation Site
```markdown
Claude: "Building a proper documentation site."

Doc site agent:
1. Set up documentation framework (Docusaurus/VitePress)
2. Organize all documentation
3. Add search functionality
4. Create versioning system
5. Add feedback widgets
6. Deploy to GitHub Pages/Netlify
```

### Step 4: Maintenance System
```markdown
Claude: "Ensuring docs stay updated."

Maintenance agent:
1. Create documentation tests
2. Add docs to PR checklist
3. Create update reminders
4. Set up broken link checker
5. Add analytics to track usage
6. Create doc review process
```

---

## Execution Tips

### 1. Trust the Process
Let agents complete their tasks. Resist the urge to micromanage.

### 2. Parallel When Possible
Independent tasks = parallel agents = faster completion

### 3. Sequential When Necessary
Some tasks need previous outputs = sequential execution

### 4. Verify at Milestones
Check in after major phases, not every step

### 5. Iterate on Failures
If an agent fails, refine the prompt and retry

### 6. Document Patterns
Save successful workflows for future reuse

---

## Common Patterns Across Workflows

### The Scout Pattern
Always start with an analysis agent to understand the landscape

### The Specialist Pattern
Deploy focused agents for specific technologies/layers

### The Validator Pattern
End with testing/verification agents

### The Documentation Pattern
Always include documentation as part of the workflow

### The Safety Pattern
Include rollback/recovery strategies in critical workflows

---

*"A journey of a thousand files begins with a single agent." - Ancient Proverb (circa 2024)*

**Your workflows are limited only by your imagination. Dream big, execute bigger! 🚀**