package anthropic

import (
	"bufio"
	"context"
	"io"

	"github.com/christianwissmann85/delegate/internal/handlers"
)

// streamResponse handles streaming responses from Anthropic
func (p *Provider) streamResponse(ctx context.Context, reader io.Reader) <-chan handlers.StreamChunk {
	ch := make(chan handlers.StreamChunk)

	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(reader)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				ch <- handlers.StreamChunk{Error: ctx.Err()}
				return
			default:
				// TODO: Parse SSE events and extract content
				line := scanner.Text()
				_ = line // Process line
			}
		}

		if err := scanner.Err(); err != nil {
			ch <- handlers.StreamChunk{Error: err}
		}
	}()

	return ch
}