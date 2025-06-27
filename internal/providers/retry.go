package providers

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/christianwissmann85/delegate/internal/handlers"
	"github.com/christianwissmann85/delegate/internal/logger"
	"github.com/christianwissmann85/delegate/internal/models"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultRetryConfig returns default retry settings
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Second,
		MaxDelay:    30 * time.Second,
	}
}

// RetryableProvider wraps a provider with retry logic
type RetryableProvider struct {
	provider handlers.Provider
	config   RetryConfig
	logger   *logger.Logger
	name     string
}

// NewRetryableProvider creates a provider with retry capabilities
func NewRetryableProvider(provider handlers.Provider, name string) *RetryableProvider {
	return &RetryableProvider{
		provider: provider,
		config:   DefaultRetryConfig(),
		logger:   logger.New("providers.retry", logger.InfoLevel),
		name:     name,
	}
}

// Generate implements handlers.Provider with retry logic
func (r *RetryableProvider) Generate(ctx context.Context, req *models.GenerateRequest) (<-chan models.StreamChunk, error) {
	ch := make(chan models.StreamChunk)

	go func() {
		defer close(ch)

		var lastErr *models.DelegateError

		for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
			r.logger.Debug("Attempting request", map[string]interface{}{
				"provider": r.name,
				"attempt":  attempt,
				"model":    req.Model,
			})

			// Create a new context for this attempt
			attemptCtx, cancel := context.WithCancel(ctx)

			// Try to generate
			stream, err := r.provider.Generate(attemptCtx, req)
			if err != nil {
				cancel()
				// Initial error before streaming starts
				lastErr = NormalizeError(r.name, err, 0)
				r.handleRetry(ch, lastErr, attempt)
				continue
			}

			// Stream the response
			success := true
			for chunk := range stream {
				if chunk.Error != nil {
					// Error during streaming
					lastErr = NormalizeError(r.name, chunk.Error, 0)
					success = false
					break
				}
				ch <- chunk
			}

			cancel()

			if success {
				// Success! No need to retry
				return
			}

			// Handle retry for streaming error
			r.handleRetry(ch, lastErr, attempt)
		}

		// All attempts failed
		if lastErr != nil {
			ch <- models.StreamChunk{
				Error: fmt.Errorf("all retry attempts failed: %w", lastErr),
			}
		}
	}()

	return ch, nil
}

// handleRetry decides whether to retry and calculates delay
func (r *RetryableProvider) handleRetry(ch chan<- models.StreamChunk, err *models.DelegateError, attempt int) {
	if attempt >= r.config.MaxAttempts {
		// No more retries
		return
	}

	if !IsRetryable(err) {
		// Error is not retryable
		ch <- models.StreamChunk{Error: err}
		return
	}

	// Calculate backoff delay
	retryAfter := 0
	if val, ok := err.Details["retry_after"].(int); ok {
		retryAfter = val
	}
	delay := r.calculateBackoff(attempt, retryAfter)

	r.logger.Info("Retrying after error", map[string]interface{}{
		"provider":    r.name,
		"attempt":     attempt,
		"error_type":  err.Code,
		"retry_after": delay.Seconds(),
	})

	// Wait before retry
	time.Sleep(delay)
}

// calculateBackoff computes the exponential backoff delay
func (r *RetryableProvider) calculateBackoff(attempt int, retryAfter int) time.Duration {
	// If provider specified retry_after, use that
	if retryAfter > 0 {
		return time.Duration(retryAfter) * time.Second
	}

	// Exponential backoff: base * 2^(attempt-1)
	delay := float64(r.config.BaseDelay) * math.Pow(2, float64(attempt-1))

	// Add jitter (±10%)
	jitter := delay * 0.1 * (2*rand() - 1)
	delay += jitter

	// Cap at max delay
	if delay > float64(r.config.MaxDelay) {
		delay = float64(r.config.MaxDelay)
	}

	return time.Duration(delay)
}

// Simple random float between 0 and 1
func rand() float64 {
	return float64(time.Now().UnixNano()%1000) / 1000.0
}
