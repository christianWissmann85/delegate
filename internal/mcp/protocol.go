package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

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
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Read next message
			line, err := p.reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("read error: %w", err)
			}

			// Parse and handle message
			if err := p.handleMessage([]byte(line)); err != nil {
				p.logger.Error("Error handling message", map[string]interface{}{
					"error": err.Error(),
				})
				// Continue processing other messages
			}
		}
	}
}

// handleMessage processes a single JSON-RPC message
func (p *Protocol) handleMessage(data []byte) error {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		p.sendError(nil, &Error{
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
	default:
		p.sendError(msg.ID, &Error{
			Code:    MethodNotFound,
			Message: fmt.Sprintf("Method not found: %s", msg.Method),
		})
		return fmt.Errorf("unknown method: %s", msg.Method)
	}
}

// handleInitialize handles the initialize request
func (p *Protocol) handleInitialize(id interface{}, params json.RawMessage) error {
	var initParams InitializeParams
	if err := json.Unmarshal(params, &initParams); err != nil {
		p.sendError(id, &Error{
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

	// Build tool list
	tools := make([]ToolInfo, 0, len(p.server.tools))
	for _, tool := range p.server.tools {
		tools = append(tools, ToolInfo{
			Name:        tool.Name(),
			Description: tool.Description(),
			InputSchema: tool.Schema(),
		})
	}

	result := InitializeResult{
		ServerInfo: ServerInfo{
			Name:    "delegate",
			Version: "1.0.0",
		},
		Tools: tools,
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
	return p.sendResponse(id, tools)
}

// handleToolCall handles tool invocation
func (p *Protocol) handleToolCall(id interface{}, params json.RawMessage) error {
	var callParams ToolCallParams
	if err := json.Unmarshal(params, &callParams); err != nil {
		p.sendError(id, &Error{
			Code:    InvalidParams,
			Message: "Invalid params",
			Data:    err.Error(),
		})
		return err
	}

	tool, exists := p.server.tools[callParams.Name]
	if !exists {
		p.sendError(id, &Error{
			Code:    InvalidParams,
			Message: fmt.Sprintf("Tool not found: %s", callParams.Name),
		})
		return fmt.Errorf("tool not found: %s", callParams.Name)
	}

	// Call the tool handler
	result, err := tool.Handler(context.Background(), callParams.Params)
	if err != nil {
		// Convert DelegateError to JSON-RPC error with rich data
		if delegateErr, ok := err.(*models.DelegateError); ok {
			p.sendError(id, &Error{
				Code:    p.mapErrorTypeToCode(delegateErr.Type),
				Message: delegateErr.Message,
				Data: map[string]interface{}{
					"error":               delegateErr.Type,
					"provider":            delegateErr.Provider,
					"retry_after":         delegateErr.RetryAfter,
					"alternative_models":  delegateErr.Alternatives,
				},
			})
		} else {
			// Fallback for non-DelegateError
			p.sendError(id, &Error{
				Code:    InternalError,
				Message: err.Error(),
			})
		}
		return err
	}

	return p.sendResponse(id, result)
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
		fmt.Fprintf(p.writer, `{"jsonrpc":"2.0","id":null,"error":{"code":-32603,"message":"Internal error"}}`+"\n")
		return errMarshal
	}

	_, errWrite := fmt.Fprintf(p.writer, "%s\n", data)
	return errWrite
}

// mapErrorTypeToCode maps DelegateError types to JSON-RPC error codes
func (p *Protocol) mapErrorTypeToCode(errorType string) int {
	switch errorType {
	case models.ErrorTypeInvalidRequest:
		return InvalidParams
	case models.ErrorTypeAuthError:
		return InvalidRequest
	case models.ErrorTypeNotFound:
		return InvalidParams
	default:
		// All provider errors map to InternalError
		return InternalError
	}
}