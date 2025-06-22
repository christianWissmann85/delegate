package extractor

// ExtractionOptions configures extraction behavior
type ExtractionOptions struct {
	PreferLanguage string // Preferred language to extract
	MaxBlocks      int    // Maximum number of blocks to extract
}

// ExtractionStats provides statistics about extraction
type ExtractionStats struct {
	TotalBlocks      int
	ExtractedBlocks  int
	TotalLines       int
	CodeLines        int
	ExplanationLines int
}
