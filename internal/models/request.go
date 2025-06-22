package models

// Request types for internal use

// InvokeParams represents internal invoke parameters
type InvokeParams struct {
	Model        string
	Prompt       string
	Files        []string
	MaxTokens    int
	CodeOnly     bool
	LanguageHint string
	Timeout      int
}

// CheckParams represents internal check parameters
type CheckParams struct {
	OutputID string
}

// ReadParams represents internal read parameters
type ReadParams struct {
	OutputID  string
	Extract   string
	MaxTokens int
}
