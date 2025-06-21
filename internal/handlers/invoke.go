package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/christianwissmann85/delegate/internal/models"
)

// InvokeHandler handles the invoke tool
type InvokeHandler struct {
	providers ProviderFactory
	storage   Storage
	extractor Extractor
}

// Provider generates content from LLMs
type Provider interface {
	GenerateStream(ctx context.Context, req GenerateRequest) (<-chan StreamChunk, error)
}

// ProviderFactory creates providers based on model
type ProviderFactory interface {
	GetProvider(model string) (Provider, error)
}

// Storage persists outputs
type Storage interface {
	Save(output *models.Output) error
	Get(id string) (*models.Output, error)
	Delete(id string) error
	ListOlderThan(age time.Duration) ([]string, error)
}

// Extractor extracts code and explanations
type Extractor interface {
	Extract(content string) (*Extraction, error)
	ExtractCode(content string) ([]CodeBlock, error)
	ExtractExplanation(content string) (string, error)
}

// NewInvokeHandler creates a new invoke handler
func NewInvokeHandler(providers ProviderFactory, storage Storage, extractor Extractor) *InvokeHandler {
	return &InvokeHandler{
		providers: providers,
		storage:   storage,
		extractor: extractor,
	}
}

// Handle processes an invoke request
func (h *InvokeHandler) Handle(ctx context.Context, req InvokeRequest) (*InvokeResponse, error) {
	// TODO: Implement invoke logic
	return nil, fmt.Errorf("not implemented")
}

// InvokeRequest represents the invoke tool parameters
type InvokeRequest struct {
	Model        string   `json:"model"`
	Prompt       string   `json:"prompt"`
	Files        []string `json:"files,omitempty"`
	MaxTokens    int      `json:"max_tokens,omitempty"`
	CodeOnly     bool     `json:"code_only,omitempty"`
	LanguageHint string   `json:"language_hint,omitempty"`
	Timeout      int      `json:"timeout,omitempty"`
}

// InvokeResponse represents the invoke tool response
type InvokeResponse struct {
	OutputID string `json:"output_id"`
}

// GenerateRequest is sent to providers
type GenerateRequest struct {
	Model     string
	Prompt    string
	Files     []string
	MaxTokens int
}

// StreamChunk represents a chunk of generated content
type StreamChunk struct {
	Content string
	Error   error
}

// Extraction contains extracted code and explanation
type Extraction struct {
	Code        []CodeBlock
	Explanation string
}

// CodeBlock represents an extracted code block
type CodeBlock struct {
	Language  string
	Content   string
	LineStart int
	LineEnd   int
}