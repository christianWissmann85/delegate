# Module Architecture - Clear Boundaries, No Surprises

## Core Principles
1. Each module has ONE responsibility
2. Dependencies flow inward (handlers → providers/storage, not vice versa)
3. Interfaces defined in the consuming module
4. No circular dependencies

## Module Dependency Graph
```
main.go
    ↓
mcp/server.go
    ↓
handlers/{invoke,check,read}.go
    ↓               ↓
providers/*    storage/*
    ↓               ↓
models/*       models/*
```

## Module Interfaces & Contracts

### 1. MCP Module
**Purpose:** Handle MCP protocol communication ONLY

**Exports:**
```go
type Server interface {
    Start(ctx context.Context) error
    RegisterTool(tool Tool) error
}

type Tool interface {
    Name() string
    Description() string
    Schema() JSONSchema
    Handler(ctx context.Context, params json.RawMessage) (interface{}, error)
}
```

**Depends on:** Nothing (except handlers for tool registration)

### 2. Handlers Module
**Purpose:** Orchestrate business logic for each tool

**Exports:**
```go
type InvokeHandler struct {
    providers ProviderFactory
    storage   Storage
    extractor Extractor
}

func (h *InvokeHandler) Handle(ctx context.Context, req InvokeRequest) (*InvokeResponse, error)
```

**Depends on:** providers, storage, extractor interfaces

### 3. Providers Module
**Purpose:** Abstract LLM communication

**Interface (defined in handlers):**
```go
type Provider interface {
    GenerateStream(ctx context.Context, req GenerateRequest) (<-chan StreamChunk, error)
}

type ProviderFactory interface {
    GetProvider(model string) (Provider, error)
}
```

**Internal Structure:**
- `anthropic/` - Anthropic-specific implementation
- `google/` - Google-specific implementation
- `factory.go` - Provider selection logic
- `errors.go` - Error normalization

### 4. Storage Module
**Purpose:** Persist and retrieve outputs

**Interface (defined in handlers):**
```go
type Storage interface {
    Save(output *Output) error
    Get(id string) (*Output, error)
    Delete(id string) error
    ListOlderThan(age time.Duration) ([]string, error)
}
```

**Depends on:** models.Output only

### 5. Extractor Module
**Purpose:** Extract code and explanations from LLM responses

**Interface (defined in handlers):**
```go
type Extractor interface {
    Extract(content string) (*Extraction, error)
    ExtractCode(content string) ([]CodeBlock, error)
    ExtractExplanation(content string) (string, error)
}

type Extraction struct {
    Code        []CodeBlock
    Explanation string
}
```

**Depends on:** Nothing

### 6. Config Module
**Purpose:** Load and validate configuration

**Exports:**
```go
type Config struct {
    LogLevel       string
    TimeoutSeconds int
    OutputDir      string
    // Provider configs
    AnthropicKey string
    GoogleKey    string
}

func Load() (*Config, error)
func (c *Config) Validate() error
```

**Depends on:** Nothing

### 7. Models Module
**Purpose:** Shared data structures (no logic!)

**Exports:**
```go
// output.go
type Output struct {
    ID         string
    CreatedAt  time.Time
    Model      string
    Prompt     string
    Files      []string
    Response   Response
    Metadata   Metadata
}

// errors.go
type DelegateError struct {
    Type         string
    Provider     string
    Code         int
    Message      string
    RetryAfter   int
    Alternatives []string
}
```

## Anti-Patterns to Avoid

### ❌ DON'T: Import across layers
```go
// Bad: storage importing providers
import "delegate/internal/providers"
```

### ❌ DON'T: Business logic in models
```go
// Bad: methods on data structures
func (o *Output) Generate() error { ... }
```

### ❌ DON'T: Direct provider access from MCP
```go
// Bad: MCP calling providers directly
provider.GenerateStream(...)
```

### ✅ DO: Keep interfaces small
```go
// Good: focused interfaces
type Extractor interface {
    Extract(content string) (*Extraction, error)
}
```

### ✅ DO: Mock at interface boundaries
```go
// Good: easy to test
type mockProvider struct{}
func (m *mockProvider) GenerateStream(...) { ... }
```

## Module Size Limits

| Module | Max Files | Max LOC/File | Notes |
|--------|-----------|--------------|-------|
| mcp | 4 | 200 | Protocol complexity contained |
| handlers | 4 | 200 | Business logic stays simple |
| providers | 7 | 200 | Split by provider |
| storage | 3 | 200 | File I/O focused |
| extractor | 3 | 200 | Regex complexity isolated |
| config | 2 | 150 | Simple validation |
| models | 3 | 100 | Data only, no logic |

## Testing Strategy

Each module has clear test boundaries:

```
internal/
├── mcp/
│   └── protocol_test.go      # Test JSON-RPC parsing
├── handlers/
│   ├── invoke_test.go        # Test with mock providers/storage
│   └── testdata/             # Sample requests
├── providers/
│   └── factory_test.go       # Test provider selection
├── extractor/
│   ├── extractor_test.go     # Test extraction patterns
│   └── testdata/             # Sample LLM outputs
└── storage/
    └── store_test.go         # Test with temp directories
```

## Summary

This architecture:
1. **Prevents refactoring** - Clear boundaries from day 1
2. **Enables parallel development** - Each module can be built independently
3. **Simplifies testing** - Mock at interface boundaries
4. **Maintains simplicity** - Each module does ONE thing

Ready to build without fear of refactoring! 🏗️