package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/christianwissmann85/delegate/internal/handlers"
	"github.com/christianwissmann85/delegate/internal/logger"
	"github.com/christianwissmann85/delegate/internal/models"
)

// Protocol handles JSON-RPC communication
type Protocol struct {
	server *Server
	reader *bufio.Reader
	writer io.Writer
	logger *logger.Logger
}

// NewProtocol creates a new protocol handler
func NewProtocol(server *Server, logLevel logger.Level) *Protocol {
	return &Protocol{
		server: server,
		reader: bufio.NewReader(os.Stdin),
		writer: os.Stdout,
		logger: logger.New("mcp.protocol", logLevel),
	}
}

// HandleMessages processes incoming JSON-RPC messages
func (p *Protocol) HandleMessages(ctx context.Context) error {
	p.logger.Info("MCP server started, waiting for messages")

	scanner := bufio.NewScanner(p.reader)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Check if there's input available
			if scanner.Scan() {
				line := scanner.Text()

				// Skip empty lines
				if strings.TrimSpace(line) == "" {
					continue
				}

				// Parse and handle message
				if err := p.handleMessage([]byte(line)); err != nil {
					p.logger.Error("Error handling message", map[string]interface{}{
						"error": err.Error(),
					})
					// Continue processing other messages
				}
			} else {
				// Check for scanner error
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("scanner error: %w", err)
				}
				// EOF reached - this is normal when the client disconnects
				return nil
			}
		}
	}
}

// handleMessage processes a single JSON-RPC message
func (p *Protocol) handleMessage(data []byte) error {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		_ = p.sendError(nil, &Error{
			Code:    ParseError,
			Message: "Parse error",
			Data:    err.Error(),
		})
		return fmt.Errorf("parse message: %w", err)
	}

	p.logger.Debug("Received method", map[string]interface{}{
		"method": msg.Method,
		"id":     msg.ID,
	})

	// Route based on method
	switch msg.Method {
	case "initialize":
		return p.handleInitialize(msg.ID, msg.Params)
	case "tools/list":
		return p.handleToolsList(msg.ID)
	case "tools/call":
		return p.handleToolCall(msg.ID, msg.Params)
	case "notifications/initialized":
		// This is a notification, not a request - no response needed
		p.logger.Info("Client initialization complete")
		return nil
	default:
		// Only send error response if this is a request (has an ID)
		if msg.ID != nil {
			_ = p.sendError(msg.ID, &Error{
				Code:    MethodNotFound,
				Message: fmt.Sprintf("Method not found: %s", msg.Method),
			})
		}
		return fmt.Errorf("unknown method: %s", msg.Method)
	}
}

// handleInitialize handles the initialize request
func (p *Protocol) handleInitialize(id interface{}, params json.RawMessage) error {
	var initParams InitializeParams
	if err := json.Unmarshal(params, &initParams); err != nil {
		_ = p.sendError(id, &Error{
			Code:    InvalidParams,
			Message: "Invalid params",
			Data:    err.Error(),
		})
		return err
	}

	p.logger.Info("Client connected", map[string]interface{}{
		"client_name":    initParams.ClientInfo.Name,
		"client_version": initParams.ClientInfo.Version,
	})

	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: Capabilities{
			Tools: &ToolsCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    "delegate",
			Version: "1.0.0",
		},
	}

	return p.sendResponse(id, result)
}

// handleToolsList handles the tools/list request
func (p *Protocol) handleToolsList(id interface{}) error {
	tools := make([]ToolInfo, 0, len(p.server.tools))
	for _, tool := range p.server.tools {
		tools = append(tools, ToolInfo{
			Name:        tool.Name(),
			Description: tool.Description(),
			InputSchema: tool.Schema(),
		})
	}

	// Wrap the tools array in an object with a "tools" field
	response := struct {
		Tools []ToolInfo `json:"tools"`
	}{
		Tools: tools,
	}

	return p.sendResponse(id, response)
}

// handleToolCall handles tool invocation
func (p *Protocol) handleToolCall(id interface{}, params json.RawMessage) error {
	var callParams ToolCallParams
	if err := json.Unmarshal(params, &callParams); err != nil {
		_ = p.sendError(id, &Error{
			Code:    InvalidParams,
			Message: "Invalid params",
			Data:    err.Error(),
		})
		return err
	}

	tool, exists := p.server.tools[callParams.Name]
	if !exists {
		_ = p.sendError(id, &Error{
			Code:    InvalidParams,
			Message: fmt.Sprintf("Tool not found: %s", callParams.Name),
		})
		return fmt.Errorf("tool not found: %s", callParams.Name)
	}

	// Call the tool handler
	result, err := tool.Handler(context.Background(), callParams.Arguments)
	if err != nil {
		// Convert DelegateError to JSON-RPC error with rich data
		if delegateErr, ok := err.(*models.DelegateError); ok {
			// Build error data from details
			errorData := map[string]interface{}{
				"error": string(delegateErr.Code),
			}
			// Copy relevant details
			if delegateErr.Details != nil {
				for k, v := range delegateErr.Details {
					errorData[k] = v
				}
			}
			_ = p.sendError(id, &Error{
				Code:    p.mapErrorTypeToCode(string(delegateErr.Code)),
				Message: delegateErr.Message,
				Data:    errorData,
			})
		} else {
			// Fallback for non-DelegateError
			_ = p.sendError(id, &Error{
				Code:    InternalError,
				Message: err.Error(),
			})
		}
		return err
	}

	// Wrap result in MCP content array format
	wrappedResult := p.wrapToolResult(callParams.Name, result)
	return p.sendResponse(id, wrappedResult)
}

// sendResponse sends a JSON-RPC response
func (p *Protocol) sendResponse(id interface{}, result interface{}) error {
	resp := Message{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal response: %w", err)
	}

	_, err = fmt.Fprintf(p.writer, "%s\n", data)
	return err
}

// sendError sends a JSON-RPC error response
func (p *Protocol) sendError(id interface{}, err *Error) error {
	resp := Message{
		JSONRPC: "2.0",
		ID:      id,
		Error:   err,
	}

	data, errMarshal := json.Marshal(resp)
	if errMarshal != nil {
		// Last resort - send a basic error
		_, _ = fmt.Fprintf(p.writer, `{"jsonrpc":"2.0","id":null,"error":{"code":-32603,"message":"Internal error"}}`+"\n")
		return errMarshal
	}

	_, errWrite := fmt.Fprintf(p.writer, "%s\n", data)
	return errWrite
}

// wrapToolResult wraps tool handler results in MCP content array format
func (p *Protocol) wrapToolResult(toolName string, result interface{}) map[string]interface{} {
	var message string

	switch toolName {
	case "delegate_invoke":
		if invokeResp, ok := result.(*handlers.InvokeResponse); ok {
			message = fmt.Sprintf("Task delegated successfully. Output ID: %s", invokeResp.OutputID)
		}
	case "delegate_check":
		if checkResp, ok := result.(*handlers.CheckResponse); ok {
			message = fmt.Sprintf("Output %s: %d bytes, ~%d tokens, created at %s",
				checkResp.ID, checkResp.FileSizeBytes, checkResp.EstimatedTokens, checkResp.CreatedAt)
		}
	case "delegate_read":
		if readResp, ok := result.(*handlers.ReadResponse); ok {
			// For read, we return the actual content
			message = readResp.Content
		}
	default:
		// Fallback: try to marshal as JSON
		if data, err := json.Marshal(result); err == nil {
			message = string(data)
		} else {
			message = "Operation completed successfully"
		}
	}

	return map[string]interface{}{
		"content": []map[string]interface{}{
			{
				"type": "text",
				"text": message,
			},
		},
	}
}

// mapErrorTypeToCode maps DelegateError types to JSON-RPC error codes
func (p *Protocol) mapErrorTypeToCode(errorType string) int {
	switch models.ErrorType(errorType) {
	case models.ErrorTypeInvalidRequest:
		return InvalidParams
	case models.ErrorTypeOutputNotFound, models.ErrorTypeFileNotFound:
		return InvalidParams
	case models.ErrorTypeProviderError, models.ErrorTypeProviderUnavailable:
		return InternalError
	case models.ErrorTypeFileWriteFailed, models.ErrorTypePathTraversalAttempt:
		return InternalError
	default:
		// All other errors map to InternalError
		return InternalError
	}
}
