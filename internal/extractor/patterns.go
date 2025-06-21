package extractor

import "regexp"

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
			Regex: regexp.MustCompile("(?s)```(?P<lang>\\w+)?\\n(?P<code>.*?)```"),
		},
		{
			Name:  "AltFencedBlock",
			Regex: regexp.MustCompile("(?s)~~~(?P<lang>\\w+)?\\n(?P<code>.*?)~~~"),
		},
		{
			Name:  "IndentedBlock",
			Regex: regexp.MustCompile("(?m)^(    .+\\n)+"),
		},
	}
}

// MatchResult contains information about a pattern match
type MatchResult struct {
	Language  string
	Content   string
	StartPos  int
	EndPos    int
}