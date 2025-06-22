package mock

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/christianwissmann85/delegate/internal/handlers"
)

// Provider implements a mock provider for testing
type Provider struct {
	model        string
	responses    []string
	delay        time.Duration
	errorOnChunk int // Error on the nth chunk (0 = no error)
}

// NewProvider creates a new mock provider
func NewProvider(model string) *Provider {
	return &Provider{
		model: model,
		delay: 10 * time.Millisecond, // Small delay to simulate streaming
	}
}

// WithResponses sets the responses to stream
func (p *Provider) WithResponses(responses ...string) *Provider {
	p.responses = responses
	return p
}

// WithDelay sets the delay between chunks
func (p *Provider) WithDelay(delay time.Duration) *Provider {
	p.delay = delay
	return p
}

// WithError makes the provider error on the nth chunk
func (p *Provider) WithError(chunkNumber int) *Provider {
	p.errorOnChunk = chunkNumber
	return p
}

// GenerateStream generates mock content
func (p *Provider) GenerateStream(ctx context.Context, req handlers.GenerateRequest) (<-chan handlers.StreamChunk, error) {
	// Validate model matches
	if req.Model != p.model {
		return nil, fmt.Errorf("model mismatch: expected %s, got %s", p.model, req.Model)
	}

	ch := make(chan handlers.StreamChunk)

	go func() {
		defer close(ch)

		// Default response if none provided
		responses := p.responses
		if len(responses) == 0 {
			// Generate a default response based on the prompt
			if strings.Contains(strings.ToLower(req.Prompt), "code") {
				responses = []string{
					"Here's a simple example:\n\n",
					"```python\n",
					"def hello_world():\n",
					"    print('Hello, World!')\n",
					"```\n",
					"\nThis function prints a greeting message.",
				}
			} else {
				responses = []string{
					"This is a mock response. ",
					"The prompt was: ",
					req.Prompt[:min(50, len(req.Prompt))],
					"...",
				}
			}
		}

		// Stream chunks with delay
		for i, chunk := range responses {
			select {
			case <-ctx.Done():
				ch <- handlers.StreamChunk{Error: ctx.Err()}
				return
			case <-time.After(p.delay):
				// Simulate error if configured
				if p.errorOnChunk > 0 && i+1 == p.errorOnChunk {
					ch <- handlers.StreamChunk{Error: fmt.Errorf("mock error on chunk %d", i+1)}
					return
				}
				ch <- handlers.StreamChunk{Content: chunk}
			}
		}
	}()

	return ch, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
