package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/christianwissmann85/delegate/internal/handlers"
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
	Type        string               `json:"type"`
	Description string               `json:"description"`
	Enum        []string             `json:"enum,omitempty"`
	Items       *Property            `json:"items,omitempty"`
	Properties  map[string]Property  `json:"properties,omitempty"`
}

// InvokeTool wraps the invoke handler as an MCP tool
type InvokeTool struct {
	handler *handlers.InvokeHandler
}

func (t *InvokeTool) Name() string {
	return "delegate.invoke"
}

func (t *InvokeTool) Description() string {
	return "Delegate heavy tasks (code generation, document analysis, large file processing) to other LLMs to save Claude Code's context tokens. Use this when: generating large amounts of code, analyzing multiple documents, processing entire codebases, or any task that would consume significant context. Supports Gemini models (1M token context) and Claude models. Returns an output_id for async retrieval."
}

func (t *InvokeTool) Schema() JSONSchema {
	return JSONSchema{
		Type: "object",
		Properties: map[string]Property{
			"model": {
				Type:        "string",
				Description: "The LLM model to use",
				Enum:        handlers.ValidModels,
			},
			"prompt": {
				Type:        "string",
				Description: "Natural language description of the task.",
			},
			"files": {
				Type:        "array",
				Description: "File paths to include as context.",
				Items: &Property{
					Type: "string",
				},
			},
			"max_tokens": {
				Type:        "number",
				Description: "Maximum tokens to generate (defaults to model maximum)",
			},
			"code_only": {
				Type:        "boolean",
				Description: "Return only code without explanations (default: false)",
			},
			"language_hint": {
				Type:        "string",
				Description: "Expected programming language(s) for better extraction",
			},
			"timeout": {
				Type:        "number",
				Description: "Request-specific timeout in seconds (overrides DELEGATE_TIMEOUT_SECONDS)",
			},
		},
		Required: []string{"model", "prompt"},
	}
}

func (t *InvokeTool) Handler(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req handlers.InvokeRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("unmarshal params: %w", err)
	}
	return t.handler.Handle(ctx, req)
}

// CheckTool wraps the check handler as an MCP tool
type CheckTool struct {
	handler *handlers.CheckHandler
}

func (t *CheckTool) Name() string {
	return "delegate.check"
}

func (t *CheckTool) Description() string {
	return "Get metadata about a delegated task output including size, token count, and creation time. Always use this before reading to avoid consuming unnecessary tokens. Returns file size in bytes and estimated token count."
}

func (t *CheckTool) Schema() JSONSchema {
	return JSONSchema{
		Type: "object",
		Properties: map[string]Property{
			"output_id": {
				Type:        "string",
				Description: "The output ID returned from invoke",
			},
		},
		Required: []string{"output_id"},
	}
}

func (t *CheckTool) Handler(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req handlers.CheckRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("unmarshal params: %w", err)
	}
	return t.handler.Handle(ctx, req)
}

// ReadTool wraps the read handler as an MCP tool
type ReadTool struct {
	handler *handlers.ReadHandler
}

func (t *ReadTool) Name() string {
	return "delegate.read"
}

func (t *ReadTool) Description() string {
	return "Retrieve results from a delegated task. Use 'extract' option to get only code or explanation. Use 'max_tokens' to limit response size. Best practice: always check() before read() to know what you're getting."
}

func (t *ReadTool) Schema() JSONSchema {
	return JSONSchema{
		Type: "object",
		Properties: map[string]Property{
			"output_id": {
				Type:        "string",
				Description: "The output ID returned from invoke",
			},
			"options": {
				Type:        "object",
				Description: "Options for reading output",
				Properties: map[string]Property{
					"extract": {
						Type:        "string",
						Description: "What to extract: 'all', 'code', 'explanation'",
						Enum:        []string{"all", "code", "explanation"},
					},
					"max_tokens": {
						Type:        "number",
						Description: "Limit response size in tokens",
					},
				},
			},
		},
		Required: []string{"output_id"},
	}
}

func (t *ReadTool) Handler(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req handlers.ReadRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("unmarshal params: %w", err)
	}
	return t.handler.Handle(ctx, req)
}