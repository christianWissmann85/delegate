package mcp

import (
	"context"
	"encoding/json"
)

// Tool represents an MCP tool
type Tool interface {
	Name() string
	Description() string
	Schema() JSONSchema
	Handler(ctx context.Context, params json.RawMessage) (interface{}, error)
}

// JSONSchema represents a JSON Schema for tool parameters
type JSONSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]Property    `json:"properties"`
	Required   []string               `json:"required,omitempty"`
}

// Property represents a JSON Schema property
type Property struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Enum        []string    `json:"enum,omitempty"`
	Items       *Property   `json:"items,omitempty"`
}

// InvokeTool wraps the invoke handler as an MCP tool
type InvokeTool struct {
	handler interface{} // Will be *handlers.InvokeHandler
}

func (t *InvokeTool) Name() string {
	return "delegate.invoke"
}

func (t *InvokeTool) Description() string {
	return "Delegate heavy tasks (code generation, document analysis, large file processing) to other LLMs to save Claude Code's context tokens. Use this when: generating large amounts of code, analyzing multiple documents, processing entire codebases, or any task that would consume significant context. Supports Gemini models (1M token context) and Claude models. Returns an output_id for async retrieval."
}

func (t *InvokeTool) Schema() JSONSchema {
	// TODO: Return full schema
	return JSONSchema{
		Type: "object",
	}
}

func (t *InvokeTool) Handler(ctx context.Context, params json.RawMessage) (interface{}, error) {
	// TODO: Implement handler
	return nil, nil
}

// CheckTool wraps the check handler as an MCP tool
type CheckTool struct {
	handler interface{} // Will be *handlers.CheckHandler
}

func (t *CheckTool) Name() string {
	return "delegate.check"
}

func (t *CheckTool) Description() string {
	return "Get metadata about a delegated task output including size, token count, and creation time. Always use this before reading to avoid consuming unnecessary tokens. Returns file size in bytes and estimated token count."
}

func (t *CheckTool) Schema() JSONSchema {
	// TODO: Return full schema
	return JSONSchema{
		Type: "object",
	}
}

func (t *CheckTool) Handler(ctx context.Context, params json.RawMessage) (interface{}, error) {
	// TODO: Implement handler
	return nil, nil
}

// ReadTool wraps the read handler as an MCP tool
type ReadTool struct {
	handler interface{} // Will be *handlers.ReadHandler
}

func (t *ReadTool) Name() string {
	return "delegate.read"
}

func (t *ReadTool) Description() string {
	return "Retrieve results from a delegated task. Use 'extract' option to get only code or explanation. Use 'max_tokens' to limit response size. Best practice: always check() before read() to know what you're getting."
}

func (t *ReadTool) Schema() JSONSchema {
	// TODO: Return full schema
	return JSONSchema{
		Type: "object",
	}
}

func (t *ReadTool) Handler(ctx context.Context, params json.RawMessage) (interface{}, error) {
	// TODO: Implement handler
	return nil, nil
}