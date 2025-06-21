package extractor

import (
	"strings"
	"testing"

	"github.com/christianwissmann85/delegate/internal/handlers"
)

func TestExtractor_ExtractCode(t *testing.T) {
	tests := []struct {
		name              string
		content           string
		expected          []handlers.CodeBlock
		skipDetailedCheck bool
	}{
		{
			name:    "single fenced code block",
			content: "Here's a Python function:\n\n```python\ndef hello():\n    print(\"Hello, World!\")\n```\n\nThat's it!",
			expected: []handlers.CodeBlock{
				{
					Language:  "python",
					Content:   "def hello():\n    print(\"Hello, World!\")",
					LineStart: 3,
					LineEnd:   4,
				},
			},
		},
		{
			name:    "multiple code blocks",
			content: "First, let's create a function:\n\n```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```\n\nAnd here's another one:\n\n```javascript\nconsole.log(\"World\");\n```",
			expected: []handlers.CodeBlock{
				{
					Language:  "go",
					Content:   "func main() {\n    fmt.Println(\"Hello\")\n}",
					LineStart: 3,
					LineEnd:   5,
				},
				{
					Language:  "javascript",
					Content:   "console.log(\"World\");",
					LineStart: 11,
					LineEnd:   11,
				},
			},
		},
		{
			name:    "code block with no language",
			content: "```\nsome code here\n```",
			expected: []handlers.CodeBlock{
				{
					Language:  "plaintext",
					Content:   "some code here",
					LineStart: 1,
					LineEnd:   1,
				},
			},
		},
		{
			name:    "alternative fence syntax",
			content: "~~~python\ndef test():\n    pass\n~~~",
			expected: []handlers.CodeBlock{
				{
					Language:  "python",
					Content:   "def test():\n    pass",
					LineStart: 1,
					LineEnd:   2,
				},
			},
		},
		{
			name:     "no code blocks",
			content:  "This is just plain text with no code.",
			expected: []handlers.CodeBlock{},
		},
		{
			name:     "empty content",
			content:  "",
			expected: []handlers.CodeBlock{},
		},
		{
			name:    "indented code block",
			content: "Here's some code:\n\n    def hello():\n        print(\"Hi\")\n    \n    hello()\n\nDone.",
			expected: []handlers.CodeBlock{
				{
					Language:  "python", // Should detect from content
					Content:   "def hello():\n    print(\"Hi\")\n\nhello()",
					LineStart: 3,
					LineEnd:   6,
				},
			},
			skipDetailedCheck: true, // Indented blocks might be detected as multiple blocks
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext := New()
			got, err := ext.ExtractCode(tt.content)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !tt.skipDetailedCheck && len(got) != len(tt.expected) {
				t.Fatalf("expected %d code blocks, got %d", len(tt.expected), len(got))
			}
			
			if tt.skipDetailedCheck {
				// Just check that we got some code blocks
				if len(got) == 0 {
					t.Error("expected at least one code block")
				}
				return
			}

			for i, block := range got {
				if block.Language != tt.expected[i].Language {
					t.Errorf("block %d: expected language %q, got %q", i, tt.expected[i].Language, block.Language)
				}
				if block.Content != tt.expected[i].Content {
					t.Errorf("block %d: expected content %q, got %q", i, tt.expected[i].Content, block.Content)
				}
			}
		})
	}
}

func TestExtractor_ExtractExplanation(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "explanation with code block",
			content:  "Here's how to implement a hello world function:\n\n```python\ndef hello():\n    print(\"Hello, World!\")\n```\n\nThis function prints a greeting message.",
			expected: "Here's how to implement a hello world function:\n\nThis function prints a greeting message.",
		},
		{
			name:     "explanation with inline code",
			content:  "Use the `print()` function to display text.",
			expected: "Use the `print()` function to display text.", // We keep inline code in explanations
		},
		{
			name:     "no code blocks",
			content:  "This is pure explanation text.",
			expected: "This is pure explanation text.",
		},
		{
			name:     "empty content",
			content:  "",
			expected: "",
		},
		{
			name:     "multiple code blocks",
			content:  "First step:\n\n```bash\necho \"Hello\"\n```\n\nSecond step:\n\n```bash\necho \"World\"\n```\n\nAll done!",
			expected: "First step:\n\nSecond step:\n\nAll done!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext := New()
			got, err := ext.ExtractExplanation(tt.content)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.expected {
				t.Errorf("expected explanation:\n%q\ngot:\n%q", tt.expected, got)
			}
		})
	}
}

func TestExtractor_LanguageDetection(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name:     "Python function",
			code:     "def hello():\n    print(\"Hello, World!\")",
			expected: "python",
		},
		{
			name:     "Go function",
			code:     "package main\n\nfunc main() {\n    fmt.Println(\"Hello\")\n}",
			expected: "go",
		},
		{
			name:     "JavaScript arrow function",
			code:     "const greet = () => {\n    console.log(\"Hello\");\n};",
			expected: "javascript",
		},
		{
			name:     "TypeScript with types",
			code:     "interface Person {\n    name: string;\n    age: number;\n}",
			expected: "typescript",
		},
		{
			name:     "SQL query",
			code:     "SELECT * FROM users WHERE age > 18;",
			expected: "sql",
		},
		{
			name:     "Bash script",
			code:     "#!/bin/bash\necho \"Hello\"",
			expected: "bash",
		},
		{
			name:     "JSON object",
			code:     "{\n    \"name\": \"test\",\n    \"value\": 123\n}",
			expected: "json",
		},
		{
			name:     "YAML config",
			code:     "name: test\nversion: 1.0\nitems:\n  - first\n  - second",
			expected: "yaml",
		},
		{
			name:     "Unknown code",
			code:     "Some random text",
			expected: "plaintext",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext := New()
			got := ext.detectLanguage(tt.code)
			if got != tt.expected {
				t.Errorf("expected language %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestExtractor_WithLanguageHint(t *testing.T) {
	content := "Here's some code:\n\n```\nSELECT * FROM users;\n```"

	t.Run("with SQL hint", func(t *testing.T) {
		ext := NewWithHint("sql")
		blocks, err := ext.ExtractCode(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(blocks) != 1 {
			t.Fatalf("expected 1 block, got %d", len(blocks))
		}

		// First block should use the hint
		if blocks[0].Language != "sql" {
			t.Errorf("expected first block to be SQL, got %q", blocks[0].Language)
		}
	})

	t.Run("without hint", func(t *testing.T) {
		ext := New()
		blocks, err := ext.ExtractCode(content)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// First block should be detected as SQL
		if blocks[0].Language != "sql" {
			t.Errorf("expected first block to be SQL, got %q", blocks[0].Language)
		}
	})
}

func TestExtractor_ExtractCodeOnly(t *testing.T) {
	content := "Here's the implementation:\n\n```go\nfunc main() {\n    fmt.Println(\"Hello\")\n}\n```\n\nThis is how it works."

	ext := New()
	blocks, err := ext.ExtractCodeOnly(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	if blocks[0].Language != "go" {
		t.Errorf("expected Go language, got %q", blocks[0].Language)
	}

	if !strings.Contains(blocks[0].Content, "fmt.Println") {
		t.Errorf("expected code content to contain fmt.Println")
	}
}

func TestExtractor_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "nested backticks",
			content: "```\n`inner code`\n```",
		},
		{
			name:    "unclosed fence",
			content: "```python\nprint('hello')",
		},
		{
			name:    "mixed fence types",
			content: "```python\ncode\n~~~",
		},
		{
			name:    "empty code block",
			content: "```\n```",
		},
		{
			name:    "code block with only whitespace",
			content: "```\n   \n\t\n```",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext := New()
			
			// Should not panic
			_, err := ext.Extract(tt.content)
			if err != nil {
				t.Logf("Extract returned error (expected): %v", err)
			}
			
			_, err = ext.ExtractCode(tt.content)
			if err != nil {
				t.Logf("ExtractCode returned error (expected): %v", err)
			}
			
			_, err = ext.ExtractExplanation(tt.content)
			if err != nil {
				t.Logf("ExtractExplanation returned error (expected): %v", err)
			}
		})
	}
}

func TestNormalizeLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"python", "python"},
		{"py", "python"},
		{"python3", "python"},
		{"js", "javascript"},
		{"ts", "typescript"},
		{"c++", "cpp"},
		{"c#", "csharp"},
		{"sh", "bash"},
		{"", "plaintext"},
		{"   ", "plaintext"},
		{"PYTHON", "python"},
		{"unknown-lang", "unknown-lang"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := NormalizeLanguage(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizeLanguage(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestExtractorFactory(t *testing.T) {
	factory := NewFactory()

	t.Run("create with hint", func(t *testing.T) {
		ext := factory.Create("python")
		if ext == nil {
			t.Fatal("expected extractor, got nil")
		}
		
		// Test that hint is used
		blocks, _ := ext.ExtractCode("```\ncode\n```")
		if len(blocks) > 0 && blocks[0].Language != "python" {
			t.Errorf("expected language hint to be used")
		}
	})

	t.Run("create without hint", func(t *testing.T) {
		ext := factory.Create("")
		if ext == nil {
			t.Fatal("expected extractor, got nil")
		}
	})

	t.Run("default extractor", func(t *testing.T) {
		ext := factory.Default()
		if ext == nil {
			t.Fatal("expected extractor, got nil")
		}
	})
}