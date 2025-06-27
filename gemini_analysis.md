The Delegate MCP server refactoring is encountering several compilation issues, primarily due to duplicate type declarations, inconsistent use of new `internal/models` types, and missing helper functions. The refactoring plan's goal of creating a "single source of truth" for models has not been fully realized, leading to conflicts.

Here's a comprehensive list of issues, organized by file, with proposed fixes to make the code compile:

---

### **1. Duplicate Type Declarations**

*   **`internal/models/output.go`**
    *   **Issue:** This file (`output.go`) defines `Output`, `Response`, `Extracted`, `ExtractedCode`, and `Metadata` which are also defined in `internal/models/models.go`. This causes "redeclared in this package" errors.
    *   **Fix:** Delete `internal/models/output.go`. Merge its unique fields (e.g., `Files` in `Output`, `LineStart`/`LineEnd` in `ExtractedCode`, `ProviderRequestID`/`ProcessingTimeMs` in `Metadata`) into `internal/models/models.go`.

*   **`handlers/invoke.go`**
    *   **Issue:** Lines 160-184: `GenerateRequest`, `StreamChunk`, `Extraction`, and `CodeBlock` are locally declared within the `handlers` package. The refactoring plan explicitly states these should be moved to `internal/models`.
    *   **Fix:** Remove these local declarations. Update all usages within `invoke.go` to use the fully qualified types from the `models` package (e.g., `models.GenerateRequest`, `models.StreamChunk`, `models.Extraction`, `models.CodeBlock`).

*   **`handlers/submit_task.go`, `handlers/get_metadata.go`, `handlers/get_content_handler.go`, `handlers/write_file_handler.go`, `handlers/invoke.go`, `handlers/check.go`, `handlers/read.go`**
    *   **Issue:** The interfaces `Provider`, `ProviderFactory`, `Storage`, `Extractor`, and `ExtractorFactory` are duplicated across multiple handler files.
    *   **Fix:** Create a new file, e.g., `handlers/interfaces.go`. Define all these interfaces there once. Then, in each handler file, remove the local interface definitions and ensure they use the qualified names (e.g., `handlers.Storage`).

---

### **2. Import Errors or Missing Imports**

*   No explicit "missing import" errors were found, but the resolution of duplicate types and helper functions will implicitly resolve some "undefined" errors.

---

### **3. Type Mismatches between Old and New Code**

*   **`handlers/invoke.go`**
    *   **Issue:** Line 40: `req GenerateRequest` should use the `models` package version.
    *   **Fix:** Change to `req models.GenerateRequest`.
    *   **Issue:** Line 78: `extraction *Extraction` should use the `models` package version.
    *   **Fix:** Change to `extraction *models.Extraction`.

---

### **4. Handler Constructor Issues**

*   No direct handler constructor issues were identified. `server.go` correctly initializes and registers the new handlers.

---

### **5. Any Other Compilation Errors**

*   **`internal/models/errors.go`**
    *   **Issue:** The constant `ErrorTypeNotFound` is used in `handlers/check.go` (line 34) and `handlers/read.go` (line 34) but is not defined in `errors.go`. The correct constant is `ErrorTypeOutputNotFound`.
    *   **Fix:** In `handlers/check.go` and `handlers/read.go`, change `models.ErrorTypeNotFound` to `models.ErrorTypeOutputNotFound`.

*   **`handlers/submit_task.go`**
    *   **Issue:** Lines 142-148: `TruncationResult`, `DetectTruncation`, `EstimateTokens`, `ValidateModel`, `ValidatePrompt`, `ValidateFilePaths`, `ValidateMaxTokens`, `ValidateTimeout` are declared as "NOTE: The following are assumed to be implemented elsewhere in the package." They are not defined.
    *   **Fix:**
        *   Move the `TruncationResult` struct definition to `internal/models/models.go`.
        *   Create a new file `handlers/validation.go` and define the `ValidateModel`, `ValidatePrompt`, `ValidateFilePaths`, `ValidateMaxTokens`, `ValidateTimeout`, `ValidateOutputID`, and `ValidateExtractOption` functions there. These functions should return `error` or `*models.DelegateError`.
        *   Create a new file `handlers/utils.go` and define `DetectTruncation` and `EstimateTokens` there.
        *   Update all calls in `submit_task.go` to use the `handlers.` prefix for these utility functions (e.g., `handlers.DetectTruncation`).

*   **`handlers/get_metadata.go`**
    *   **Issue:** Line 79: `ValidateOutputID` is not defined.
    *   **Fix:** Define `ValidateOutputID` in `handlers/validation.go` (as per the fix for `submit_task.go`). Update call to `handlers.ValidateOutputID`.

*   **`handlers/get_content_handler.go`**
    *   **Issue:** Line 30: `ValidateOutputID` is not defined.
    *   **Issue:** Line 37: `ValidateExtractOption` is not defined.
    *   **Issue:** Line 59: `ValidateMaxTokens` is not defined.
    *   **Fix:** Define these validation functions in `handlers/validation.go` and update calls to use `handlers.` prefix.
    *   **Issue:** Line 108: Incomplete `fmt.Sprintf("```")`. This is a syntax error.
    *   **Fix:** Complete the `fmt.Sprintf` to `fence := fmt.Sprintf("```%s\n%s\n```", block.Language, block.Content)`.
    *   **Issue:** `h.truncateContent` and `h.extractCodeContent` are called but not defined on `GetContentHandler`.
    *   **Fix:** Move `truncateContent` and `extractCodeContent` from `read.go` to `handlers/utils.go` and update calls to `handlers.TruncateContent` and `handlers.ExtractCodeContent`.

*   **`handlers/write_file_handler.go`**
    *   **Issue:** Line 29: `ValidateOutputID` is not defined.
    *   **Issue:** Line 39: `ValidateExtractOption` is not defined.
    *   **Fix:** Define these validation functions in `handlers/validation.go` and update calls to use `handlers.` prefix.
    *   **Issue:** Line 144: `h.getExtractedCodeForWriting` calls `h.isDocumentationFile` and `h.cleanupCodeArtifacts`, which are not defined on `WriteFileHandler`.
    *   **Fix:** Move `isDocumentationFile` and `cleanupCodeArtifacts` from `read.go` to `handlers/utils.go` and update calls to `handlers.IsDocumentationFile` and `handlers.CleanupCodeArtifacts`.
    *   **Issue:** The `writeToFile` method is called (line 90) with 4 arguments (`req.WriteTo`, `absolutePath`, `cwd`, `[]byte(contentToWrite)`), but its definition is missing in this file. The plan states this logic should be moved from `read.go`.
    *   **Fix:** Move the `writeToFile` function from `read.go` to `write_file_handler.go`. **Crucially, update its signature** to match the call: `func (h *WriteFileHandler) writeToFile(relativePath, absolutePath, cwd string, content []byte) (int64, error)`. Ensure the security comments are preserved as per the plan.
    *   **Issue:** Line 159: Incomplete `fmt.Sprintf("```")`.
    *   **Fix:** Complete the `fmt.Sprintf` to `fence := fmt.Sprintf("```%s\n%s\n```", block.Language, block.Content)`.

*   **`handlers/invoke.go` (Old Handler)**
    *   **Issue:** Lines 43, 55, 64, 131: Incorrect usage of `models.NewDelegateError`. It expects `key, value` pairs for `details` (e.g., `"model", req.Model`), not just raw strings or single values.
    *   **Fix:** Adjust calls to `models.NewDelegateError` to conform to its `(code ErrorType, message string, args ...interface{})` signature.
        *   Example: `return nil, models.NewDelegateError(models.ErrorTypeProviderUnavailable, fmt.Sprintf("Provider for model '%s' unavailable.", req.Model), "model", req.Model, "original_error", err)`
        *   For the `save output` error, use: `return nil, models.NewDelegateError(models.ErrorTypeInternal, "Failed to save output.", "original_error", err)`

*   **`handlers/read.go` (Old Handler)**
    *   **Issue:** This file still contains the `writeToFile` logic (lines 118-169, and the `writeToFile` method itself from line 260), `isDocumentationFile`, `cleanupCodeArtifacts`, and `truncateContent` methods. The refactoring plan dictates that the file writing logic and related helpers should move to `write_file_handler.go` and `handlers/utils.go`. While not a direct compilation error *yet*, it's a major deviation from the plan and will cause conflicts once the new handlers are fully functional and this file is intended for deprecation/removal.
    *   **Fix (for compilation and adherence to plan):**
        *   Remove the `writeToFile` method from `ReadHandler`.
        *   Remove the `if req.Options.WriteTo != ""` block (lines 118-169) from `Handle`.
        *   Remove `isDocumentationFile`, `cleanupCodeArtifacts`, `truncateContent`, `extractCodeContent`, `extractCodeForFile` methods from `ReadHandler`.
        *   Update calls to these functions to use `handlers.` prefix (e.g., `handlers.TruncateContent`, `handlers.ExtractCodeContent`).
        *   The `ReadResponse` struct (lines 208-219) still contains fields related to file writing (`FileWritten`) and multi-block warnings (`MultipleBlocks`, `BlockCount`, `BlockDescriptors`) which are now handled by `GetOutputMetadataResponse` and `WriteOutputToFileResponse`. These fields should be removed from `ReadResponse` if it's truly becoming `GetContentResponse`. For compilation, they can remain, but it's a design mismatch.

*   **`mcp/tools.go`**
    *   **Issue:** Line 44: `handlers.ValidModels` is used in the `SubmitTaskTool` schema but is not defined in the `handlers` package.
    *   **Fix:** Define `var ValidModels = []string{"gpt-4", "claude-3-opus-20240229" /* ... etc. */}` in `handlers/validation.go` (or a dedicated `handlers/constants.go` file).

---

**Summary of Required Actions (Ordered for Dependency Resolution):**

1.  **Consolidate Models:**
    *   Delete `internal/models/output.go`.
    *   Merge its fields into `internal/models/models.go`:
        *   `Output`: Add `Files []string`.
        *   `ExtractedCode`: Add `LineStart int`, `LineEnd int`.
        *   `Metadata`: Add `ProviderRequestID string`, `ProcessingTimeMs int64`.
        *   Move `TruncationResult` struct from `handlers/submit_task.go` to `internal/models/models.go`.

2.  **Create Shared Handler Components:**
    *   Create `handlers/interfaces.go` and move all duplicated interface definitions (`Provider`, `ProviderFactory`, `Storage`, `Extractor`, `ExtractorFactory`) there.
    *   Create `handlers/validation.go` and define all `ValidateX` functions (`ValidateModel`, `ValidatePrompt`, `ValidateFilePaths`, `ValidateMaxTokens`, `ValidateTimeout`, `ValidateOutputID`, `ValidateExtractOption`) there. Also define `var ValidModels = []string{...}` here.
    *   Create `handlers/utils.go` and define `DetectTruncation`, `EstimateTokens`, `TruncateContent`, `IsDocumentationFile`, `CleanupCodeArtifacts`, `ExtractCodeContent` there.

3.  **Update All Handler Files:**
    *   **Remove local interface definitions** from all handler files (`submit_task.go`, `get_metadata.go`, `get_content_handler.go`, `write_file_handler.go`, `invoke.go`, `check.go`, `read.go`).
    *   **`handlers/submit_task.go`**:
        *   Remove local `GenerateRequest`, `StreamChunk`, `Extraction`, `CodeBlock`, `TruncationResult` structs.
        *   Update all type usages to `models.X` and utility function calls to `handlers.X`.
    *   **`handlers/get_metadata.go`**:
        *   Change `models.ErrorTypeNotFound` to `models.ErrorTypeOutputNotFound`.
        *   Update `ValidateOutputID` call to `handlers.ValidateOutputID`.
    *   **`handlers/get_content_handler.go`**:
        *   Fix `fmt.Sprintf` syntax error.
        *   Update all validation and utility function calls to `handlers.X`.
    *   **`handlers/write_file_handler.go`**:
        *   Move `writeToFile` function from `read.go` to this file, updating its signature to `func (h *WriteFileHandler) writeToFile(relativePath, absolutePath, cwd string, content []byte) (int64, error)`.
        *   Update all validation and utility function calls to `handlers.X`.
        *   Fix `fmt.Sprintf` syntax error.
    *   **`handlers/invoke.go`**:
        *   Remove all local type declarations (`GenerateRequest`, `StreamChunk`, `Extraction`, `CodeBlock`, `TruncationResult`).
        *   Update all type usages to `models.X` and utility function calls to `handlers.X`.
        *   Correct `models.NewDelegateError` calls to pass details as key-value pairs.
    *   **`handlers/check.go`**:
        *   Change `models.ErrorTypeNotFound` to `models.ErrorTypeOutputNotFound`.
        *   Update `ValidateOutputID` call to `handlers.ValidateOutputID`.
    *   **`handlers/read.go`**:
        *   Remove `writeToFile`, `isDocumentationFile`, `cleanupCodeArtifacts`, `truncateContent`, `extractCodeContent`, `extractCodeForFile` methods from `ReadHandler`.
        *   Remove the entire `if req.Options.WriteTo != ""` block from `Handle`.
        *   Update all validation and utility function calls to `handlers.X`.
        *   Change `models.ErrorTypeNotFound` to `models.ErrorTypeOutputNotFound`.

4.  **Update `mcp/tools.go`**:
    *   Ensure `handlers.ValidModels` is correctly referenced (it will be defined in `handlers/validation.go`).

These steps will address all identified compilation errors and bring the codebase closer to the refactoring plan's architectural goals.