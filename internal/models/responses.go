package models

// SubmitTaskResponse represents the successful response for the delegate_submit_task tool.
type SubmitTaskResponse struct {
	OutputID         string `json:"output_id"`         // Unique identifier for the generated output.
	WorkingDirectory string `json:"working_directory"` // The current working directory of the server.
}

// GetOutputMetadataResponse represents the successful response for the delegate_get_output_metadata tool.
type GetOutputMetadataResponse struct {
	Metadata       MetadataBlock      `json:"metadata"`
	ContentAnalysis ContentAnalysisBlock `json:"content_analysis"`
}

// MetadataBlock contains general metadata about the output artifact.
type MetadataBlock struct {
	OutputID         string  `json:"output_id"`
	Status           string  `json:"status"` // Enum: "COMPLETED", "IN_PROGRESS", "FAILED"
	SizeKB           float64 `json:"size_kb"`
	LineCount        int     `json:"line_count"`
	TokenEstimate    int     `json:"token_estimate"`
	IsTruncated      bool    `json:"is_truncated"`
	TruncationReason *string `json:"truncation_reason"` // Reason for truncation, e.g., 'MAX_TOKENS_REACHED'. Null if not truncated.
}

// ContentAnalysisBlock provides details about the structure of the output content,
// especially for multi-block outputs.
type ContentAnalysisBlock struct {
	BlocksFound int               `json:"blocks_found"`
	Blocks      []CodeBlockMetadata `json:"blocks"`
}

// CodeBlockMetadata describes a single code block within the output content,
// providing summary information without the full content.
type CodeBlockMetadata struct {
	Index    int     `json:"index"`
	Language string  `json:"language"`
	SizeKB   float64 `json:"size_kb"`
	Lines    int     `json:"lines"`
	Preview  string  `json:"preview"` // The first line of the block.
}

// GetOutputContentResponse represents the successful response for the delegate_get_output_content tool.
type GetOutputContentResponse struct {
	Content  string             `json:"content"`
	Metadata ContentMetadataBlock `json:"metadata"`
}

// GetContentResponse is an alias for GetOutputContentResponse for backward compatibility
type GetContentResponse = GetOutputContentResponse

// ContentMetadataBlock contains metadata specific to the retrieved content,
// indicating if it was truncated for the agent's context.
type ContentMetadataBlock struct {
	OutputID         string  `json:"output_id"`
	TokensReturned   int     `json:"tokens_returned"`
	IsTruncated      bool    `json:"is_truncated"`     // True if the content was truncated by the max_tokens parameter.
	TruncationReason *string `json:"truncation_reason"` // Reason for truncation, e.g., 'MAX_TOKENS_REACHED'. Null if not truncated.
}

// WriteOutputToFileResponse represents the successful response for the delegate_write_output_to_file tool.
type WriteOutputToFileResponse struct {
	Success          bool   `json:"success"`
	Path             string `json:"path"`              // The relative path of the file written.
	AbsolutePath     string `json:"absolute_path"`     // The absolute path of the file written.
	BytesWritten     int64  `json:"bytes_written"`
	Message          string `json:"message"`           // A human-readable success message.
	WorkingDirectory string `json:"working_directory"` // The current working directory.
}

// ErrorResponse represents the structured error format for all failed tool calls.
type ErrorResponse struct {
	Error DelegateErrorDetails `json:"error"`
}

// DelegateErrorDetails contains the machine-readable code and developer-friendly message for an error,
// along with optional additional context.
type DelegateErrorDetails struct {
	Code    ErrorType              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"` // Optional object containing additional context.
}