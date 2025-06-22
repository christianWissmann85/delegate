package handlers

// Common types shared across handlers

// ExtractOption specifies what to extract
type ExtractOption string

const (
	ExtractAll         ExtractOption = "all"
	ExtractCode        ExtractOption = "code"
	ExtractExplanation ExtractOption = "explanation"
)

// ValidModels contains all supported models
var ValidModels = []string{
	"gemini-2.5-flash",
	"gemini-2.5-pro",
	"claude-sonnet-4-20250514",
	"claude-opus-4-20250514",
}
