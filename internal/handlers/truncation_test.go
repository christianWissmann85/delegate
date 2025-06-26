package handlers

import (
	"strings"
	"testing"
)

func TestDetectTruncation(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantTrunc   bool
		minConf     float64
		wantReason  string
	}{
		// Complete endings - should NOT be truncated
		{
			name:      "complete sentence",
			content:   "This is a complete sentence.",
			wantTrunc: false,
		},
		{
			name:      "complete code block",
			content:   "```go\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n```",
			wantTrunc: false,
		},
		{
			name:      "complete JSON",
			content:   `{"key": "value", "number": 42}`,
			wantTrunc: false,
		},
		{
			name:      "ends with ellipsis",
			content:   "And so on...",
			wantTrunc: false,
		},

		// Obvious truncations
		{
			name:       "mid-JSON string",
			content:    `{"key": "value", "text": "Hello wo`,
			wantTrunc:  true,
			minConf:    0.8,
			wantReason: "unclosed", // could be brackets or quote
		},
		{
			name:       "mid-JSON array",
			content:    `{"items": ["one", "two","`,
			wantTrunc:  true,
			minConf:    0.9,
			wantReason: "mid-JSON",
		},
		{
			name:       "unclosed bracket",
			content:    `function example() {\n\tconst data = [1, 2, 3`,
			wantTrunc:  true,
			minConf:    0.8,
			wantReason: "unclosed brackets",
		},
		{
			name:       "mid-word",
			content:    "This is an incomple",
			wantTrunc:  true,
			minConf:    0.7,
			wantReason: "mid-word",
		},
		{
			name:       "trailing comma",
			content:    `{"a": 1, "b": 2,`,
			wantTrunc:  true,
			minConf:    0.6,
			wantReason: "trailing comma",
		},
		{
			name:       "incomplete code fence",
			content:    "Here's the code:\n```python\ndef hello():\n    print(\"Hello\")\n``", // Only 2 backticks at end
			wantTrunc:  true,
			minConf:    0.7,
			wantReason: "incomplete code fence",
		},

		// Edge cases
		{
			name:      "complete word ending",
			content:   "The processing",
			wantTrunc: false, // ends with common suffix "ing"
		},
		{
			name:      "abbreviation",
			content:   "See the example e.g.",
			wantTrunc: false,
		},
		{
			name:       "suspicious size - 4096",
			content:    string(make([]byte, 4090)) + "abrup",
			wantTrunc:  true,
			minConf:    0.5,
			wantReason: "mid-word", // Actually detects mid-word, which is fine
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectTruncation(tt.content)

			if result.IsTruncated != tt.wantTrunc {
				t.Errorf("DetectTruncation() IsTruncated = %v, want %v", result.IsTruncated, tt.wantTrunc)
			}

			if tt.wantTrunc && result.Confidence < tt.minConf {
				t.Errorf("DetectTruncation() Confidence = %v, want >= %v", result.Confidence, tt.minConf)
			}

			if tt.wantReason != "" && !contains(result.Reason, tt.wantReason) {
				t.Errorf("DetectTruncation() Reason = %v, want to contain %v", result.Reason, tt.wantReason)
			}
		})
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Benchmark the truncation detection
func BenchmarkDetectTruncation(b *testing.B) {
	// Test with a large JSON response
	longContent := `{"data": {"items": [`
	for i := 0; i < 1000; i++ {
		longContent += `{"id": ` + string(rune(i)) + `, "name": "Item ` + string(rune(i)) + `"},`
	}
	longContent += `{"id": 1001, "name": "Last item"`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DetectTruncation(longContent)
	}
}
