package models

import (
	"fmt"
)

// ErrorType defines machine-readable error codes for Delegate API errors.
type ErrorType string

const (
	// ErrorTypeInvalidRequest indicates issues with the request parameters (e.g., missing required fields, invalid values).
	ErrorTypeInvalidRequest ErrorType = "INVALID_REQUEST"
	// ErrorTypeOutputNotFound indicates that the requested output ID does not exist or has expired.
	ErrorTypeOutputNotFound ErrorType = "OUTPUT_NOT_FOUND"
	// ErrorTypeProviderError indicates an error originating from the external LLM provider (e.g., API error, rate limit).
	ErrorTypeProviderError ErrorType = "PROVIDER_ERROR"
	// ErrorTypeFileWriteFailed indicates a failure during a file writing operation (e.g., permissions, disk full).
	ErrorTypeFileWriteFailed ErrorType = "FILE_WRITE_FAILED"
	// ErrorTypePathTraversalAttempt indicates a security violation where a path traversal was attempted.
	ErrorTypePathTraversalAttempt ErrorType = "PATH_TRAVERSAL_ATTEMPT"
	// ErrorTypeInternal indicates an unexpected internal server error not covered by other specific codes.
	ErrorTypeInternal ErrorType = "INTERNAL_ERROR"
	// ErrorTypeProviderUnavailable indicates that the requested LLM model provider is not available or configured.
	ErrorTypeProviderUnavailable ErrorType = "PROVIDER_UNAVAILABLE"
	// ErrorTypeFileNotFound indicates that a requested file could not be found.
	ErrorTypeFileNotFound ErrorType = "FILE_NOT_FOUND"
)

// DelegateError is a custom error type for Delegate MCP server.
// It encapsulates a machine-readable code, a developer-friendly message, and optional details.
// It also allows wrapping an underlying error for internal debugging and chaining.
type DelegateError struct {
	Code    ErrorType              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Err     error                  `json:"-"` // Original error for internal logging/chaining, not exposed in JSON
}

// NewDelegateError creates a new DelegateError.
// It takes an ErrorType, a message, and optional key-value pairs for details.
// The last argument can optionally be an error to wrap.
//
// Example usage:
// NewDelegateError(models.ErrorTypeInvalidRequest, "Invalid ID", "id_provided", "123", "reason", "format error")
// NewDelegateError(models.ErrorTypeProviderError, "LLM call failed", "model", "gpt-4", fmt.Errorf("API timeout"))
func NewDelegateError(code ErrorType, message string, args ...interface{}) *DelegateError {
	de := &DelegateError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}

	for i := 0; i < len(args); i++ {
		if err, ok := args[i].(error); ok {
			// If it's an error and it's the last argument, treat it as the wrapped error.
			if i == len(args)-1 {
				de.Err = err
				continue
			}
		}
		// Treat as key-value pair for details
		if i+1 < len(args) {
			if key, ok := args[i].(string); ok {
				de.Details[key] = args[i+1]
				i++ // Consume the value
			}
			// If key is not a string, it's an invalid argument, just skip it.
		}
	}
	return de
}

// Error implements the error interface for DelegateError.
// It returns the developer-friendly message, optionally including the wrapped error's message.
func (e *DelegateError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error, allowing for error chaining with errors.Is and errors.As.
func (e *DelegateError) Unwrap() error {
	return e.Err
}

// AsDelegateError attempts to convert a given error to a *DelegateError.
// If the error is already a *DelegateError, it returns it directly.
// If it's a standard error, it wraps it as an ErrorTypeInternal.
// If the input error is nil, it returns nil.
func AsDelegateError(err error) *DelegateError {
	if err == nil {
		return nil
	}
	if de, ok := err.(*DelegateError); ok {
		return de
	}
	// Default to internal error if it's not a known DelegateError type
	return NewDelegateError(ErrorTypeInternal, "An unexpected internal error occurred.", "original_error", err, err)
}