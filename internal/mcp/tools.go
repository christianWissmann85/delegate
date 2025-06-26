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
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

// Property represents a JSON Schema property
type Property struct {
	Type        string              `json:"type"`
	Description string              `json:"description"`
	Enum        []string            `json:"enum,omitempty"`
	Items       *Property           `json:"items,omitempty"`
	Properties  map[string]Property `json:"properties,omitempty"`
}

// InvokeTool wraps the invoke handler as an MCP tool
type InvokeTool struct {
	handler *handlers.InvokeHandler
}

func (t *InvokeTool) Name() string {
	return "delegate_invoke"
}

func (t *InvokeTool) Description() string {
	return "STEP 1: Delegate file generation to save tokens. Does NOT write files directly - stores in temp storage. Returns output_id for use with delegate_check then delegate_read(write_to). IMPORTANT: Use ABSOLUTE paths in 'files' parameter. Each file must be <1MB, but total can exceed."
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
				Description: "Task for ONE file. Good: 'Create user model with GORM tags'. Bad: 'Create entire REST API' (too many files).",
			},
			"files": {
				Type:        "array",
				Description: "ABSOLUTE file paths to include as context (not relative!). Example: '/home/user/project/model.go' not 'model.go'",
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
				Description: "Extract only code blocks from response, no explanations (default: false). Tip: Also ask for 'only code' in prompt.",
			},
			"language_hint": {
				Type:        "string",
				Description: "Expected programming language(s) for better extraction",
			},
			"timeout": {
				Type:        "number",
				Description: "Timeout in seconds (default: 180, max: 600). Suggested: 180s minimum for code, 400s minimum for creative tasks, 400-600s for very large file(s)/bundle(s) analysis.",
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
	return "delegate_check"
}

func (t *CheckTool) Description() string {
	return "STEP 2: Check delegated task status and size before retrieving. Shows token count and file size. Use this after delegate_invoke, before delegate_read. Helps avoid consuming unnecessary tokens."
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
	return "delegate_read"
}

func (t *ReadTool) Description() string {
	return "STEP 3: Get delegated results. WORKFLOW: invoke -> check -> read. To save tokens: use 'write_to' with ABSOLUTE path to write file directly (no content returned). To get content: omit 'write_to'. IMPORTANT: Use extract:'code' to strip markdown fences for source files, 'all' keeps formatting."
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
						Description: "What to extract: 'all' (keeps markdown), 'code' (strips fences for clean source), 'explanation' (docs only)",
						Enum:        []string{"all", "code", "explanation"},
					},
					"max_tokens": {
						Type:        "number",
						Description: "Limit response size in tokens",
					},
					"write_to": {
						Type:        "string",
						Description: "ABSOLUTE file path to write content directly to disk (saves tokens - no content returned). Example: '/home/user/project/new_file.go' not 'new_file.go'",
					},
					"block_index": {
						Type:        "number",
						Description: "When multiple code blocks exist, specify which one to extract (0-based index). Use after getting block list warning.",
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
