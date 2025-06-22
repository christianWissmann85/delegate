package providers

import (
	"net/http"
	"strings"

	"github.com/christianwissmann85/delegate/internal/models"
)

// NormalizeError converts provider-specific errors to DelegateError
func NormalizeError(provider string, err error, statusCode int) *models.DelegateError {
	errMsg := err.Error()
	errLower := strings.ToLower(errMsg)

	// Determine error type based on status code and message
	errorType := determineErrorType(statusCode, errLower)

	// Create base error
	delegateErr := &models.DelegateError{
		Type:     errorType,
		Provider: provider,
		Code:     statusCode,
		Message:  errMsg,
	}

	// Add retry_after for rate limits
	if errorType == models.ErrorTypeRateLimited {
		// Default to 60 seconds if not specified
		delegateErr.RetryAfter = 60

		// Try to extract retry_after from error message
		// TODO: Add provider-specific parsing for retry_after values when needed
	}

	// Suggest alternatives for certain errors
	if errorType == models.ErrorTypeRateLimited || errorType == models.ErrorTypeProviderUnavailable {
		delegateErr.Alternatives = suggestAlternatives(provider)
	}

	return delegateErr
}

// determineErrorType maps status codes and error messages to error types
func determineErrorType(statusCode int, errMsg string) string {
	// Check status code first
	switch statusCode {
	case http.StatusUnauthorized:
		return models.ErrorTypeAuthError
	case http.StatusNotFound:
		return models.ErrorTypeNotFound
	case http.StatusTooManyRequests:
		return models.ErrorTypeRateLimited
	case http.StatusServiceUnavailable, http.StatusBadGateway, http.StatusGatewayTimeout:
		return models.ErrorTypeProviderUnavailable
	case http.StatusRequestTimeout:
		return models.ErrorTypeTimeout
	}

	// Check error message patterns
	switch {
	case strings.Contains(errMsg, "rate limit"):
		return models.ErrorTypeRateLimited
	case strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded"):
		return models.ErrorTypeTimeout
	case strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "no such host"):
		return models.ErrorTypeNetworkError
	case strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "bad request"):
		return models.ErrorTypeInvalidRequest
	case strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "authentication"):
		return models.ErrorTypeAuthError
	case strings.Contains(errMsg, "service unavailable") || strings.Contains(errMsg, "503"):
		return models.ErrorTypeProviderUnavailable
	default:
		return models.ErrorTypeProviderError
	}
}

// suggestAlternatives returns alternative models when a provider fails
func suggestAlternatives(failedProvider string) []string {
	switch failedProvider {
	case "google":
		return []string{"claude-sonnet-4-20250514", "claude-opus-4-20250514"}
	case "anthropic":
		return []string{"gemini-2.5-flash", "gemini-2.5-pro"}
	default:
		// Return all available models except the failed provider's
		alternatives := []string{}
		if failedProvider != "google" {
			alternatives = append(alternatives, "gemini-2.5-flash", "gemini-2.5-pro")
		}
		if failedProvider != "anthropic" {
			alternatives = append(alternatives, "claude-sonnet-4-20250514")
		}
		return alternatives
	}
}

// IsRetryable determines if an error should be retried
func IsRetryable(err *models.DelegateError) bool {
	switch err.Type {
	case models.ErrorTypeRateLimited,
		models.ErrorTypeProviderUnavailable,
		models.ErrorTypeTimeout,
		models.ErrorTypeNetworkError:
		return true
	default:
		return false
	}
}
