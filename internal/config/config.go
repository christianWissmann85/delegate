package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration values
type Config struct {
	// Logging
	LogLevel string

	// Timeouts
	TimeoutSeconds int

	// Storage
	OutputDir string

	// Provider API keys
	AnthropicKey string
	GoogleKey    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		// Defaults
		LogLevel:       getEnv("DELEGATE_LOG_LEVEL", "info"),
		TimeoutSeconds: getEnvInt("DELEGATE_TIMEOUT_SECONDS", 60),
		OutputDir:      getEnv("DELEGATE_OUTPUT_DIR", ".delegate"),

		// API keys (no defaults)
		AnthropicKey: os.Getenv("ANTHROPIC_API_KEY"),
		GoogleKey:    os.Getenv("GOOGLE_API_KEY"),
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// getEnv returns an environment variable or default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns an environment variable as int or default
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// HasProvider checks if at least one provider is configured
func (c *Config) HasProvider() bool {
	return c.AnthropicKey != "" || c.GoogleKey != ""
}

// SupportedModels returns models available based on configured API keys
func (c *Config) SupportedModels() []string {
	var models []string
	
	if c.GoogleKey != "" {
		models = append(models, "gemini-2.5-flash", "gemini-2.5-pro")
	}
	
	if c.AnthropicKey != "" {
		models = append(models, "claude-sonnet-4-20250514", "claude-opus-4-20250514")
	}
	
	return models
}