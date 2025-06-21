package config

import (
	"fmt"
	"strings"
)

// Validate ensures the configuration is valid
func (c *Config) Validate() error {
	// Check log level
	validLogLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLogLevels, c.LogLevel) {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	// Check timeout
	if c.TimeoutSeconds < 1 || c.TimeoutSeconds > 600 {
		return fmt.Errorf("timeout must be between 1 and 600 seconds")
	}

	// Check at least one provider is configured
	if !c.HasProvider() {
		return fmt.Errorf("at least one API key must be set (ANTHROPIC_API_KEY or GOOGLE_API_KEY)")
	}

	// Validate API keys format (basic check)
	if c.AnthropicKey != "" && !strings.HasPrefix(c.AnthropicKey, "sk-") {
		return fmt.Errorf("invalid Anthropic API key format")
	}

	return nil
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}