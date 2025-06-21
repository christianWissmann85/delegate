package google

import (
	"context"
	"encoding/json"
	"io"

	"github.com/christianwissmann85/delegate/internal/handlers"
)

// streamResponse handles streaming responses from Gemini
func (p *Provider) streamResponse(ctx context.Context, reader io.Reader) <-chan handlers.StreamChunk {
	ch := make(chan handlers.StreamChunk)

	go func() {
		defer close(ch)
		decoder := json.NewDecoder(reader)

		for {
			select {
			case <-ctx.Done():
				ch <- handlers.StreamChunk{Error: ctx.Err()}
				return
			default:
				// TODO: Parse streaming JSON responses
				var chunk interface{}
				if err := decoder.Decode(&chunk); err != nil {
					if err == io.EOF {
						return
					}
					ch <- handlers.StreamChunk{Error: err}
					return
				}
				// Process chunk
			}
		}
	}()

	return ch
}