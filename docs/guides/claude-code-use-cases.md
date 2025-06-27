# Claude Code Use Cases

This guide shows practical examples of using Delegate for common coding tasks.

## 📚 Document & Knowledge Processing

### Analyze Multiple Documents
```
Use delegate_submit_task to analyze all RFC documents in docs/rfcs/ and summarize the key architectural decisions.
Then, use delegate_write_output_to_file to save the summary to docs/rfcs/summary.md.
```

### Research Across Codebases
```
Use delegate_submit_task to find all authentication patterns across these 5 microservice repositories.
Then, use delegate_write_output_to_file to save the findings to analysis/auth_patterns.md.
```

### Specification Compliance
```
Use delegate_submit_task to verify if our API implementation matches the OpenAPI spec in api-spec.yaml.
Then, use delegate_write_output_to_file to save the compliance report to reports/api_compliance.md.
```

## 🧪 Testing & Quality

### Comprehensive Test Generation
```
Use delegate_submit_task to generate complete test suites for the payment module including edge cases.
Then, use delegate_write_output_to_file to save the generated tests to tests/payment_module_tests.py.
```

### Test Coverage Analysis
```
Use delegate_submit_task to analyze existing tests and identify what functionality lacks coverage.
Then, use delegate_write_output_to_file to save the coverage analysis to reports/test_coverage.md.
```

### Performance Test Generation
```
Use delegate_submit_task to create load testing scenarios based on production usage patterns.
Then, use delegate_write_output_to_file to save the scenarios to tests/performance_scenarios.yaml.
```

## 🔄 Code Transformation

### Framework Migration
```
Use delegate_submit_task to convert this Express.js application to Fastify maintaining all endpoints.
Then, use delegate_write_output_to_file to save the converted code to src/fastify_app.js.
```

### Language Translation
```
Use delegate_submit_task to port this Python data processing pipeline to Go.
Then, use delegate_write_output_to_file to save the Go code to src/pipeline.go.
```

### Modernization
```
Use delegate_submit_task to refactor this callback-based code to use async/await patterns.
Then, use delegate_write_output_to_file to save the refactored code to src/modern_code.js.
```

## 📊 Analysis & Insights

### Log Analysis
```
Use delegate_submit_task to analyze these production logs and identify error patterns and bottlenecks.
Then, use delegate_write_output_to_file to save the analysis to reports/log_analysis.md.
```

### Code Quality Review
```
Use delegate_submit_task to review this codebase for security vulnerabilities and best practice violations.
Then, use delegate_write_output_to_file to save the review report to reports/code_review.md.
```

### Dependency Analysis
```
Use delegate_submit_task to analyze package.json files across all services and find version conflicts.
Then, use delegate_write_output_to_file to save the dependency report to reports/dependency_analysis.md.
```

## 🏗️ Infrastructure & DevOps

### IaC Generation
```
Use delegate_submit_task to generate Kubernetes manifests for this microservices architecture.
Then, use delegate_write_output_to_file to save the manifests to k8s/deployment.yaml.
```

### CI/CD Pipeline
```
Use delegate_submit_task to create GitHub Actions workflows based on this project structure.
Then, use delegate_write_output_to_file to save the workflow to .github/workflows/ci.yaml.
```

### Dockerfile Optimization
```
Use delegate_submit_task to optimize these Dockerfiles for smaller image sizes and better caching.
Then, use delegate_write_output_to_file to save the optimized Dockerfile to Dockerfile.optimized.
```

## 📝 Documentation

### API Documentation
```
Use delegate_submit_task to generate comprehensive API docs from these controller files.
Then, use delegate_write_output_to_file to save the docs to docs/api/index.md.
```

### Architecture Diagrams
```
Use delegate_submit_task to create PlantUML diagrams from this codebase structure.
Then, use delegate_write_output_to_file to save the diagram to docs/architecture.puml.
```

### README Generation
```
Use delegate_submit_task to create a professional README based on the codebase analysis.
Then, use delegate_write_output_to_file to save the README to README.md.
```

## 💡 Creative Uses

### Data Generation
```
Use delegate_submit_task to generate realistic test data matching our database schema.
Then, use delegate_write_output_to_file to save the data to data/test_data.json.
```

### Configuration Templates
```
Use delegate_submit_task to create environment-specific config files from this base template.
Then, use delegate_write_output_to_file to save the config to config/production.yaml.
```

### Migration Scripts
```
Use delegate_submit_task to generate database migration scripts from these schema changes.
Then, use delegate_write_output_to_file to save the script to db/migrations/001_initial.sql.
```

## 🔥 The Token-Free Workflow

### Generate Without Reading
```
Use delegate_submit_task to create src/auth/jwt.go implementing JWT authentication with refresh tokens.
Then use delegate_write_output_to_file to write it directly to disk without consuming tokens for content.
```

### Iterative Development
```
Use delegate_submit_task to create models/user.py with SQLAlchemy.
Then use delegate_write_output_to_file to write it to disk.

Use delegate_submit_task to create routes/auth.py using the user model.
Then use delegate_write_output_to_file to write it to disk.

Use delegate_submit_task to create tests/test_auth.py for the auth routes.
Then use delegate_write_output_to_file to write it to disk.
```

### The Compile-Fix Loop
```
Use delegate_submit_task to fix the TypeScript compilation errors in src/app.ts, providing src/app.ts as context.
Then use delegate_write_output_to_file to write the fixed version back to src/app.ts.
```

## Best Practices

1. **Generate ONE FILE at a time**
   - More reliable than multi-file generation
   - Allows iterative refinement
   - Prevents timeouts.

2. **Prioritize `delegate_write_output_to_file`**
   - Saves thousands of tokens per file by writing directly to disk.
   - Automatically handles code formatting.
   - Shows exactly how many tokens saved (ZERO content tokens).

3. **Always specify the model for `delegate_submit_task`**
   - Quick tasks: gemini-2.5-flash
   - Large documents: gemini-2.5-pro (1M context!)
   - Complex logic: claude-opus-4.

4. **Use `files` parameter liberally with `delegate_submit_task`**
   - Delegate handles large file inputs well.
   - More context = better results.

5. **Check metadata before retrieving content (if you must retrieve)**
   - Always use `delegate_get_output_metadata` to see output size and structure.
   - Prefer `delegate_write_output_to_file` over `delegate_get_output_content` to save tokens!