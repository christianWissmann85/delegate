# Claude Code Use Cases for Delegate

A quick reference for creative ways to use Delegate beyond basic code generation.

## 📚 Document & Knowledge Processing

### Analyze Multiple Documents
```
Use Delegate to analyze all RFC documents in docs/rfcs/ and summarize the key architectural decisions
```

### Research Across Codebases
```
Use Delegate to find all authentication patterns across these 5 microservice repositories
```

### Specification Compliance
```
Use Delegate to verify if our API implementation matches the OpenAPI spec in api-spec.yaml
```

## 🧪 Testing & Quality

### Comprehensive Test Generation
```
Use Delegate to generate complete test suites for the payment module including edge cases
```

### Test Coverage Analysis
```
Use Delegate to analyze existing tests and identify what functionality lacks coverage
```

### Performance Test Generation
```
Use Delegate to create load testing scenarios based on production usage patterns
```

## 🔄 Code Transformation

### Framework Migration
```
Use Delegate to convert this Express.js application to Fastify maintaining all endpoints
```

### Language Translation
```
Use Delegate to port this Python data processing pipeline to Go
```

### Modernization
```
Use Delegate to refactor this callback-based code to use async/await patterns
```

## 📊 Analysis & Insights

### Log Analysis
```
Use Delegate to analyze these production logs and identify error patterns and bottlenecks
```

### Code Quality Review
```
Use Delegate to review this codebase for security vulnerabilities and best practice violations
```

### Dependency Analysis
```
Use Delegate to analyze package.json files across all services and find version conflicts
```

## 🏗️ Infrastructure & DevOps

### IaC Generation
```
Use Delegate to generate Kubernetes manifests for this microservices architecture
```

### CI/CD Pipeline
```
Use Delegate to create GitHub Actions workflows based on this project structure
```

### Dockerfile Optimization
```
Use Delegate to optimize these Dockerfiles for smaller image sizes and better caching
```

## 📝 Documentation

### API Documentation
```
Use Delegate to generate comprehensive API docs from these controller files
```

### Architecture Diagrams
```
Use Delegate to create PlantUML diagrams from this codebase structure
```

### README Generation
```
Use Delegate to create a professional README based on the codebase analysis
```

## 💡 Creative Uses

### Data Generation
```
Use Delegate to generate realistic test data matching our database schema
```

### Configuration Templates
```
Use Delegate to create environment-specific config files from this base template
```

### Migration Scripts
```
Use Delegate to generate database migration scripts from these schema changes
```

## 🔥 The Token-Free Workflow

### Generate Without Reading
```
Use Delegate to create src/auth/jwt.go implementing JWT authentication with refresh tokens
Then write it directly to disk without reading
```

### Iterative Development
```
Use Delegate to create models/user.py with SQLAlchemy
Write to disk
Use Delegate to create routes/auth.py using the user model
Write to disk
Use Delegate to create tests/test_auth.py for the auth routes
Write to disk
```

### The Compile-Fix Loop
```
Use Delegate to fix the TypeScript compilation errors in src/app.ts
Write the fixed version back to src/app.ts
```

## Best Practices

1. **Generate ONE FILE at a time**
   - More reliable than multi-file generation
   - Allows iterative refinement
   - Prevents timeouts

2. **Use write_to for everything**
   - Saves thousands of tokens per file
   - Automatically handles code formatting
   - Shows exactly how many tokens saved

3. **Always specify the model based on task size**
   - Quick tasks: gemini-2.5-flash
   - Large documents: gemini-2.5-pro (1M context!)
   - Complex logic: claude-opus-4

4. **Use file context liberally**
   - Delegate handles large file inputs well
   - More context = better results

5. **Check before reading (if you must read)**
   - Always use check() to see output size
   - But prefer write_to over reading!