package extractor

import (
	"regexp"
	"strings"
)

// Pattern represents a regex pattern for code extraction
type Pattern struct {
	Name  string
	Regex *regexp.Regexp
}

// GetPatterns returns all code extraction patterns
func GetPatterns() []Pattern {
	return []Pattern{
		{
			Name:  "FencedCodeBlock",
			Regex: regexp.MustCompile(`(?s)` + "```" + `(?P<lang>[\w+#-]*)\s*\n(?P<code>.*?)` + "```"),
		},
		{
			Name:  "AltFencedBlock",
			Regex: regexp.MustCompile(`(?s)~~~(?P<lang>[\w+#-]*)\s*\n(?P<code>.*?)~~~`),
		},
		{
			Name:  "IndentedBlock",
			Regex: regexp.MustCompile(`(?m)^((?:    |\t).+(?:\n|$))+`),
		},
		{
			Name:  "HTMLCodeBlock",
			Regex: regexp.MustCompile(`(?s)<code(?:\s+class="language-(\w+)")?>(.+?)</code>`),
		},
		{
			Name:  "MarkdownInlineCode",
			Regex: regexp.MustCompile("`([^`]+)`"),
		},
	}
}

// GetLanguageHints returns common language identifiers and their variants
func GetLanguageHints() map[string][]string {
	return map[string][]string{
		"python":     {"python", "py", "python3", "py3"},
		"javascript": {"javascript", "js", "node", "nodejs"},
		"typescript": {"typescript", "ts"},
		"go":         {"go", "golang"},
		"java":       {"java"},
		"cpp":        {"cpp", "c++", "cc", "cxx"},
		"c":          {"c"},
		"csharp":     {"csharp", "c#", "cs"},
		"rust":       {"rust", "rs"},
		"ruby":       {"ruby", "rb"},
		"php":        {"php"},
		"swift":      {"swift"},
		"kotlin":     {"kotlin", "kt"},
		"sql":        {"sql", "mysql", "postgres", "postgresql"},
		"bash":       {"bash", "sh", "shell", "zsh"},
		"powershell": {"powershell", "ps1"},
		"json":       {"json"},
		"yaml":       {"yaml", "yml"},
		"xml":        {"xml"},
		"html":       {"html", "htm"},
		"css":        {"css", "scss", "sass"},
		"markdown":   {"markdown", "md"},
		"dockerfile": {"dockerfile", "docker"},
		"terraform":  {"terraform", "tf", "hcl"},
		"plaintext":  {"text", "txt", "plaintext", "plain"},
	}
}

// NormalizeLanguage converts various language identifiers to standard names
func NormalizeLanguage(lang string) string {
	trimmed := strings.TrimSpace(lang)
	if trimmed == "" {
		return "plaintext"
	}
	
	langLower := strings.ToLower(trimmed)
	hints := GetLanguageHints()
	
	for standard, variants := range hints {
		for _, variant := range variants {
			if langLower == variant {
				return standard
			}
		}
	}
	
	// If not found in hints, return the original but cleaned
	return langLower
}

// MatchResult contains information about a pattern match
type MatchResult struct {
	Language  string
	Content   string
	StartPos  int
	EndPos    int
}