package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	// Import the new handler structs from the handlers package
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

// SubmitTaskTool wraps the submit task handler as an MCP tool
type SubmitTaskTool struct {
	handler *handlers.SubmitTaskHandler // Renamed from InvokeHandler
}

func (t *SubmitTaskTool) Name() string {
	return "delegate_submit_task"
}

func (t *SubmitTaskTool) Description() string {
	return "STEP 1: Submits a generation task to an external LLM (~50-100 tokens). This is an asynchronous operation that creates a temporary output artifact and returns a unique `output_id`. The content is NOT returned directly. Use other `delegate_*` tools to access the output."
}

func (t *SubmitTaskTool) Schema() JSONSchema {
	return JSONSchema{
		Type: "object",
		Properties: map[string]Property{
			"model": {
				Type:        "string",
				Description: "The LLM model to use.",
				Enum:        handlers.ValidModels,
			},
			"prompt": {
				Type:        "string",
				Description: "The task description.",
			},
			"files": {
				Type:        "array",
				Description: "List of relative file paths to include as context (e.g., \"src/model.go\", \"docs/api.md\").",
				Items: &Property{
					Type: "string",
				},
			},
			"max_tokens": {
				Type:        "number",
				Description: "Max tokens for the generation.",
			},
			"timeout": {
				Type:        "number",
				Description: "Timeout in seconds. Default is 180 seconds. Advised to use at least 180 seconds for all tasks, up to 600 seconds maxium. Never use less than 180secs unless you know what you're doing.",
			},
		},
		Required: []string{"model", "prompt"},
	}
}

func (t *SubmitTaskTool) Handler(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req handlers.SubmitTaskRequest // Renamed from InvokeRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("unmarshal params: %w", err)
	}
	return t.handler.Handle(ctx, req)
}

// GetMetadataTool wraps the get metadata handler as an MCP tool
type GetMetadataTool struct {
	handler *handlers.GetMetadataHandler // Renamed from CheckHandler
}

func (t *GetMetadataTool) Name() string {
	return "delegate_get_output_metadata"
}

func (t *GetMetadataTool) Description() string {
	return "STEP 2 (Optional): Retrieves structured metadata about an output artifact (~20 tokens). Use this to decide whether to retrieve content into context or write directly to a file. This tool does NOT return the content itself."
}

func (t *GetMetadataTool) Schema() JSONSchema {
	return JSONSchema{
		Type: "object",
		Properties: map[string]Property{
			"output_id": {
				Type:        "string",
				Description: "The ID from `delegate_submit_task`.",
			},
		},
		Required: []string{"output_id"},
	}
}

func (t *GetMetadataTool) Handler(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req handlers.GetMetadataRequest // Renamed from CheckRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("unmarshal params: %w", err)
	}
	return t.handler.Handle(ctx, req)
}

// GetContentTool wraps the get content handler as an MCP tool
type GetContentTool struct {
	handler *handlers.GetContentHandler // Renamed from ReadHandler
}

func (t *GetContentTool) Name() string {
	return "delegate_get_output_content"
}

func (t *GetContentTool) Description() string {
	return "Retrieves the full or partial content of an output artifact into the agent's context (~30+ tokens plus content). This operation consumes tokens proportional to the content size. Use `options` to extract specific parts (e.g., `extract: 'code'`)."
}

func (t *GetContentTool) Schema() JSONSchema {
	return JSONSchema{
		Type: "object",
		Properties: map[string]Property{
			"output_id": {
				Type:        "string",
				Description: "The ID from `delegate_submit_task`.",
			},
			"options": {
				Type:        "object",
				Description: "Options for retrieving output content.",
				Properties: map[string]Property{
					"extract": {
						Type:        "string",
						Description: "What part to extract: 'all', 'code', 'explanation'.",
						Enum:        []string{"all", "code", "explanation"},
					},
					"max_tokens": {
						Type:        "number",
						Description: "Truncate the returned content to this many tokens.",
					},
					"block_index": {
						Type:        "number",
						Description: "For multi-block outputs, select a specific block (0-based index).",
					},
					"language": {
						Type:        "string",
						Description: "Filter code blocks by this language.",
					},
				},
			},
		},
		Required: []string{"output_id"},
	}
}

func (t *GetContentTool) Handler(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req handlers.GetContentRequest // Renamed from ReadRequest
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("unmarshal params: %w", err)
	}
	return t.handler.Handle(ctx, req)
}

// WriteFileTool wraps the write file handler as an MCP tool
type WriteFileTool struct {
	handler *handlers.WriteFileHandler // New handler
}

func (t *WriteFileTool) Name() string {
	return "delegate_write_output_to_file"
}

func (t *WriteFileTool) Description() string {
	return "Writes the content of an output artifact directly to a specified file path (relative to working directory). This operation consumes ZERO content tokens. Use `options` to select specific parts to write (e.g., `extract: 'code'`, `block_index: 0`)."
}

func (t *WriteFileTool) Schema() JSONSchema {
	return JSONSchema{
		Type: "object",
		Properties: map[string]Property{
			"output_id": {
				Type:        "string",
				Description: "The ID from `delegate_submit_task`.",
			},
			"write_to": {
				Type:        "string",
				Description: "The relative file path to write to (e.g., \"src/component.jsx\", \"tmp/output.go\").",
			},
			"options": {
				Type:        "object",
				Description: "Options for writing output content to a file.",
				Properties: map[string]Property{
					"extract": {
						Type:        "string",
						Description: "What part to extract: 'all', 'code', 'explanation'.",
						Enum:        []string{"all", "code", "explanation"},
					},
					"block_index": {
						Type:        "number",
						Description: "For multi-block outputs, select a specific block (0-based index).",
					},
					"language": {
						Type:        "string",
						Description: "Filter code blocks by this language.",
					},
				},
			},
		},
		Required: []string{"output_id", "write_to"},
	}
}

func (t *WriteFileTool) Handler(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var req handlers.WriteFileRequest // New request type
	if err := json.Unmarshal(params, &req); err != nil {
		return nil, fmt.Errorf("unmarshal params: %w", err)
	}
	return t.handler.Handle(ctx, req)
}