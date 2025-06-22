package models

// DelegateError represents a structured error response
type DelegateError struct {
	Type         string   `json:"error"`
	Provider     string   `json:"provider"`
	Code         int      `json:"error_code,omitempty"`
	Message      string   `json:"message"`
	RetryAfter   int      `json:"retry_after,omitempty"`
	Alternatives []string `json:"alternative_models,omitempty"`
}

// Error implements the error interface
func (e *DelegateError) Error() string {
	return e.Message
}

// Common error types
const (
	ErrorTypeRateLimited         = "rate_limited"
	ErrorTypeProviderUnavailable = "provider_unavailable"
	ErrorTypeTimeout             = "timeout"
	ErrorTypeProviderError       = "provider_error"
	ErrorTypeNetworkError        = "network_error"
	ErrorTypeInvalidRequest      = "invalid_request"
	ErrorTypeAuthError           = "auth_error"
	ErrorTypeNotFound            = "not_found"
	ErrorTypeExtractionFailed    = "extraction_failed"
	ErrorTypeInternal            = "internal_error"
)

// NewDelegateError creates a new delegate error
func NewDelegateError(errorType, provider, message string) *DelegateError {
	return &DelegateError{
		Type:     errorType,
		Provider: provider,
		Message:  message,
	}
}
