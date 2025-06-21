package providers

import (
	"github.com/christianwissmann85/delegate/internal/models"
)

// NormalizeError converts provider-specific errors to DelegateError
func NormalizeError(provider string, err error, statusCode int) *models.DelegateError {
	// TODO: Implement error normalization logic
	return &models.DelegateError{
		Type:     "provider_error",
		Provider: provider,
		Code:     statusCode,
		Message:  err.Error(),
	}
}

// IsRetryable determines if an error should be retried
func IsRetryable(err *models.DelegateError) bool {
	switch err.Type {
	case "rate_limited", "provider_unavailable", "timeout":
		return true
	default:
		return false
	}
}