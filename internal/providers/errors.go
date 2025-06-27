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

	// Build details map
	details := make([]interface{}, 0, 8)
	details = append(details, "provider", provider)
	details = append(details, "status_code", statusCode)
	
	// Add retry_after for rate limits
	if strings.Contains(errLower, "rate limit") {
		details = append(details, "retry_after", 60)
		details = append(details, "is_rate_limit", true)
	}
	
	// Add timeout flag
	if strings.Contains(errLower, "timeout") || strings.Contains(errLower, "deadline") {
		details = append(details, "is_timeout", true)
	}
	
	// Suggest alternatives for provider unavailable errors
	if errorType == models.ErrorTypeProviderUnavailable {
		details = append(details, "alternatives", suggestAlternatives(provider))
	}
	
	// Add original error at the end
	details = append(details, err)

	return models.NewDelegateError(errorType, errMsg, details...)
}

// determineErrorType maps status codes and error messages to error types
func determineErrorType(statusCode int, errMsg string) models.ErrorType {
	// Check status code first
	switch statusCode {
	case http.StatusUnauthorized:
		return models.ErrorTypeProviderError // Auth errors are provider errors
	case http.StatusNotFound:
		return models.ErrorTypeOutputNotFound
	case http.StatusTooManyRequests:
		return models.ErrorTypeProviderError // Rate limit is a provider error
	case http.StatusServiceUnavailable, http.StatusBadGateway, http.StatusGatewayTimeout:
		return models.ErrorTypeProviderUnavailable
	case http.StatusRequestTimeout:
		return models.ErrorTypeProviderError // Timeout is a provider error
	}

	// Check error message patterns
	switch {
	case strings.Contains(errMsg, "rate limit"):
		return models.ErrorTypeProviderError
	case strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded"):
		return models.ErrorTypeProviderError
	case strings.Contains(errMsg, "connection refused") || strings.Contains(errMsg, "no such host"):
		return models.ErrorTypeProviderError
	case strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "bad request"):
		return models.ErrorTypeInvalidRequest
	case strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "authentication"):
		return models.ErrorTypeProviderError
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
	// Check error code to determine if retryable
	switch err.Code {
	case models.ErrorTypeProviderUnavailable,
		models.ErrorTypeProviderError:
		// Check details for specific retryable conditions
		if details, ok := err.Details["is_rate_limit"].(bool); ok && details {
			return true
		}
		if details, ok := err.Details["is_timeout"].(bool); ok && details {
			return true
		}
		// Check if it's a connection error in the message
		if strings.Contains(strings.ToLower(err.Message), "connection") {
			return true
		}
		// Provider unavailable is always retryable
		if err.Code == models.ErrorTypeProviderUnavailable {
			return true
		}
		return false
	default:
		return false
	}
}