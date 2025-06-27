Of course. Here is the comprehensive refactoring plan for the Delegate MCP server, synthesizing all provided feedback and recommendations. This document is designed to be the single source of truth for the refactoring effort.

---

# **Delegate MCP Server Refactoring Plan v2.0**

**Document Status:** Final | **Version:** 1.0 | **Date:** 2024-10-27

## 1. Executive Summary

This document outlines a comprehensive refactoring plan for the Delegate MCP server. The primary goal is to enhance the clarity, reliability, and usability of the API for AI agents, directly addressing feedback synthesized from multiple analyses (`gemini_workflow_analysis.md`, `claude_workflow_analysis.md`, `ai-conclusion.md`).

**What We Are Changing:**
We are refactoring the existing three-tool API (`delegate_invoke`, `delegate_check`, `delegate_read`) into a new, more explicit four-tool architecture. The core of this change is splitting the ambiguous, dual-purpose `delegate_read` tool into two distinct, single-purpose tools: one for retrieving content into the agent's context and one for writing content directly to a file. Furthermore, all API responses, including errors and warnings, will be converted from human-readable strings to structured JSON.

**Why We Are Changing It:**
The current API forces a high cognitive load on AI agents. The `delegate_read` tool's behavior changes drastically based on a single parameter, and critical information (like multi-block warnings) is returned as unstructured strings, requiring brittle parsing. This ambiguity leads to complex decision-making logic, potential errors, and a confusing developer experience. The refactoring will:
*   **Eliminate Ambiguity:** Each tool will have one clear purpose and predictable behavior.
*   **Enable Programmatic Decisions:** Structured JSON responses will allow agents to reliably parse metadata, errors, and warnings without string manipulation.
*   **Increase Reliability:** A simpler, more intuitive workflow will reduce the likelihood of incorrect tool usage.
*   **Clarify Token Costs:** The cost of each operation will be transparent and explicit in the tool's function.

This refactoring strictly adheres to the **NO_SCOPE_CREEP** directive. It is a clarification and simplification of existing functionality, not an addition of new features.

### Path Handling Change: From Absolute to Relative

**What:** All file paths in the API will now use relative paths instead of absolute paths.

**Why:** 
- **Reduces cognitive load** - AI agents naturally think in project-relative terms ("src/main.go" not "/home/user/project/src/main.go")
- **Prevents errors** - No more constructing long absolute paths
- **Saves tokens** - Shorter paths mean fewer tokens in requests
- **Improves portability** - Same commands work across different environments

**How:**
- All `files` arrays and `write_to` parameters accept relative paths only
- Paths are resolved relative to the delegate server's working directory
- Responses include `working_directory` field where relevant for context
- The `write_output_to_file` response includes both relative `path` and `absolute_path` for clarity

## 2. New 4-Tool Architecture

The existing three tools will be replaced by the following four, which map directly to the original capabilities but with improved clarity.

### Tool Overview

| New Tool Name | Replaces | Core Function | Agent Token Cost |
| :--- | :--- | :--- | :--- |
| `delegate_submit_task` | `delegate_invoke` | Submits a generation task and returns an `output_id`. | **Low** (only the response) |
| `delegate_get_output_metadata` | `delegate_check` | Retrieves structured metadata about an output artifact. | **Low** (only the response) |
| `delegate_get_output_content` | `delegate_read` (read mode) | Retrieves content into the agent's context. | **High** (proportional to content size) |
| `delegate_write_output_to_file` | `delegate_read` (write mode) | Writes content directly to a file on disk. | **ZERO** |

---

### Tool Definitions

#### **1. `delegate_submit_task`**
*   **Description:** "STEP 1: Submits a generation task to an external LLM (~50-100 tokens). This is an asynchronous operation that creates a temporary output artifact and returns a unique `output_id`. The content is NOT returned directly. Use other `delegate_*` tools to access the output."
*   **Parameters:**
    *   `model` (string, required): The LLM model to use.
    *   `prompt` (string, required): The task description.
    *   `files` (array of strings, optional): List of relative file paths to include as context (e.g., "src/model.go", "docs/api.md").
    *   `max_tokens` (number, optional): Max tokens for the generation.
    *   `timeout` (number, optional): Timeout in seconds.
*   **Response:**
    ```json
    {
      "output_id": "out_20241027_103000",
      "working_directory": "/home/user/project"
    }
    ```

#### **2. `delegate_get_output_metadata`**
*   **Description:** "STEP 2 (Optional): Retrieves structured metadata about an output artifact (~20 tokens). Use this to decide whether to retrieve content into context or write directly to a file. This tool does NOT return the content itself."
*   **Parameters:**
    *   `output_id` (string, required): The ID from `delegate_submit_task`.
*   **Response:** See Section 3 for the full schema.
    ```json
    {
      "metadata": {
        "output_id": "out_20241027_103000",
        "status": "COMPLETED",
        "size_kb": 15.7,
        "line_count": 312,
        "token_estimate": 3925,
        "is_truncated": false,
        "truncation_reason": null
      },
      "content_analysis": {
        "blocks_found": 2,
        "blocks": [
          { "index": 0, "language": "jsx", "size_kb": 12.1, "lines": 250, "preview": "import React from 'react';" },
          { "index": 1, "language": "md", "size_kb": 3.6, "lines": 62, "preview": "# Explanation of the component" }
        ]
      }
    }
    ```

#### **3. `delegate_get_output_content`**
*   **Description:** "Retrieves the full or partial content of an output artifact into the agent's context (~30+ tokens plus content). This operation consumes tokens proportional to the content size. Use `options` to extract specific parts (e.g., `extract: 'code'`)."
*   **Parameters:**
    *   `output_id` (string, required): The ID from `delegate_submit_task`.
    *   `options` (object, optional):
        *   `extract` (string, enum: `all`, `code`, `explanation`): What part to extract.
        *   `max_tokens` (number): Truncate the returned content to this many tokens.
        *   `block_index` (number): For multi-block outputs, select a specific block.
        *   `language` (string): Filter code blocks by this language.
*   **Response:**
    ```json
    {
      "content": "import React from 'react';\n// ...rest of the code...",
      "metadata": {
        "output_id": "out_20241027_103000",
        "tokens_returned": 3925,
        "is_truncated": false,
        "truncation_reason": null
      }
    }
    ```

#### **4. `delegate_write_output_to_file`**
*   **Description:** "Writes the content of an output artifact directly to a specified file path (relative to working directory). This operation consumes ZERO content tokens. Use `options` to select specific parts to write (e.g., `extract: 'code'`, `block_index: 0`)."
*   **Parameters:**
    *   `output_id` (string, required): The ID from `delegate_submit_task`.
    *   `write_to` (string, required): The relative file path to write to (e.g., "src/component.jsx", "tmp/output.go").
    *   `options` (object, optional):
        *   `extract` (string, enum: `all`, `code`, `explanation`): What part to extract.
        *   `block_index` (number): For multi-block outputs, select a specific block.
        *   `language` (string): Filter code blocks by this language.
*   **Response:**
    ```json
    {
      "success": true,
      "path": "src/component.jsx",
      "absolute_path": "/home/user/project/src/component.jsx",
      "bytes_written": 12388,
      "message": "Successfully wrote 12.1 KB to src/component.jsx",
      "working_directory": "/home/user/project"
    }
    ```

---

### Improved Workflow Examples

#### Scenario A: Zero-Token Code Generation to File
**Goal:** Generate code and save it directly to a file.

1.  **Agent Call:** `delegate_submit_task`
    ```json
    { "tool_name": "delegate_submit_task", "parameters": { "prompt": "Refactor this to use hooks.", "files": ["src/Component.jsx"] } }
    ```
2.  **Server Response:** `{ "output_id": "gen-abc-123", "working_directory": "/home/user/project" }`
3.  **Agent Call:** The agent knows its goal is to write a file, so it calls `delegate_write_output_to_file` directly.
    ```json
    { "tool_name": "delegate_write_output_to_file", "parameters": { "output_id": "gen-abc-123", "write_to": "src/Component.hooks.jsx", "options": { "extract": "code" } } }
    ```
4.  **Server Response:**
    ```json
    { "success": true, "path": "src/Component.hooks.jsx", "absolute_path": "/home/user/project/src/Component.hooks.jsx", "bytes_written": 4182, "message": "...", "working_directory": "/home/user/project" }
    ```
**Outcome:** A clear, two-step, unambiguous workflow.

#### Scenario B: Reviewing Multi-Block Content Before Writing
**Goal:** Understand the output before deciding which part to save.

1.  **Agent Call:** `delegate_submit_task` -> **Server Response:** `{ "output_id": "gen-xyz-789" }`
2.  **Agent Call:** The agent is unsure of the output, so it checks the metadata first.
    ```json
    { "tool_name": "delegate_get_output_metadata", "parameters": { "output_id": "gen-xyz-789" } }
    ```
3.  **Server Response:** The agent receives structured data indicating multiple blocks.
    ```json
    { "metadata": { ... }, "content_analysis": { "blocks_found": 2, "blocks": [ { "index": 0, "language": "jsx", ... }, { "index": 1, "language": "md", ... } ] } }
    ```
4.  **Agent Decision:** The agent programmatically inspects the `content_analysis.blocks` array and decides to write the first block (`index: 0`) to a file.
5.  **Agent Call:**
    ```json
    { "tool_name": "delegate_write_output_to_file", "parameters": { "output_id": "gen-xyz-789", "write_to": "src/code.jsx", "options": { "block_index": 0 } } }
    ```
6.  **Server Response:** `{ "success": true, "path": "src/code.jsx", "absolute_path": "/home/user/project/src/code.jsx", "bytes_written": 12388, "message": "...", "working_directory": "/home/user/project" }`
**Outcome:** The agent handled a complex case without parsing strings, using a clear, logical flow.

## 3. Structured JSON Response Schemas

All responses will adhere to a strict JSON schema. String-based warnings are eliminated.

### Tool Success Schemas

*   **`delegate_submit_task`**
    ```json
    {
      "type": "object",
      "properties": {
        "output_id": { "type": "string", "description": "Unique identifier for the generated output." }
      },
      "required": ["output_id"]
    }
    ```
*   **`delegate_get_output_metadata`**
    ```json
    {
      "type": "object",
      "properties": {
        "metadata": {
          "type": "object",
          "properties": {
            "output_id": { "type": "string" },
            "status": { "type": "string", "enum": ["COMPLETED", "IN_PROGRESS", "FAILED"] },
            "size_kb": { "type": "number" },
            "line_count": { "type": "number" },
            "token_estimate": { "type": "number" },
            "is_truncated": { "type": "boolean" },
            "truncation_reason": { "type": ["string", "null"] }
          }
        },
        "content_analysis": {
          "type": "object",
          "properties": {
            "blocks_found": { "type": "number" },
            "blocks": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "index": { "type": "number" },
                  "language": { "type": "string" },
                  "size_kb": { "type": "number" },
                  "lines": { "type": "number" },
                  "preview": { "type": "string", "description": "The first line of the block." }
                },
                "required": ["index", "language", "size_kb", "lines", "preview"]
              }
            }
          }
        }
      },
      "required": ["metadata", "content_analysis"]
    }
    ```
*   **`delegate_get_output_content`**
    ```json
    {
      "type": "object",
      "properties": {
        "content": { "type": "string" },
        "metadata": {
          "type": "object",
          "properties": {
            "output_id": { "type": "string" },
            "tokens_returned": { "type": "number" },
            "is_truncated": { "type": "boolean", "description": "True if the content was truncated by the max_tokens parameter." },
            "truncation_reason": { "type": ["string", "null"], "description": "Reason for truncation, e.g., 'MAX_TOKENS_REACHED'." }
          }
        }
      },
      "required": ["content", "metadata"]
    }
    ```
*   **`delegate_write_output_to_file`**
    ```json
    {
      "type": "object",
      "properties": {
        "success": { "type": "boolean" },
        "path": { "type": "string", "description": "The relative path of the file written." },
        "absolute_path": { "type": "string", "description": "The absolute path of the file written." },
        "bytes_written": { "type": "number" },
        "message": { "type": "string", "description": "A human-readable success message." },
        "working_directory": { "type": "string", "description": "The current working directory." }
      },
      "required": ["success", "path", "absolute_path", "bytes_written", "message", "working_directory"]
    }
    ```

### Structured Error Response Schema

A single, consistent error format will be used for all failed tool calls.

```json
{
  "type": "object",
  "properties": {
    "error": {
      "type": "object",
      "properties": {
        "code": {
          "type": "string",
          "description": "A machine-readable error code.",
          "enum": ["INVALID_REQUEST", "OUTPUT_NOT_FOUND", "PROVIDER_ERROR", "FILE_WRITE_FAILED", "PATH_TRAVERSAL_ATTEMPT"]
        },
        "message": { "type": "string", "description": "A developer-friendly error message." },
        "details": {
          "type": "object",
          "description": "Optional object containing additional context about the error.",
          "additionalProperties": true
        }
      },
      "required": ["code", "message"]
    }
  }
}
```
**Example Error:**
```json
{
  "error": {
    "code": "OUTPUT_NOT_FOUND",
    "message": "The requested output ID does not exist or has expired.",
    "details": {
      "output_id_provided": "gen-this-is-fake"
    }
  }
}
```

## 4. Architecture Improvements (from Code Review)

The refactoring will incorporate the following improvements identified in the code review.

*   **Error Handling Enhancements:**
    *   Implement the new structured error format across all handlers.
    *   Create a helper function `models.AsDelegateError()` to reduce boilerplate when creating and wrapping custom errors, ensuring consistency.

*   **Security Improvements:**
    *   Add explicit, detailed comments to the `writeToFile` function in the new `delegate_write_output_to_file` handler, explaining the path traversal prevention and sandboxing mechanisms. This ensures future maintainers understand and preserve these critical security checks.

*   **Code Organization:**
    *   **Handler Separation:** The logic currently in `read.go` will be split. The file writing logic will move to a new `write_output_handler.go`, and the content retrieval logic will move to `get_content_handler.go`.
    *   **Shared Models:** Relocate shared data structures (`Extraction`, `CodeBlock`, `GenerateRequest`, etc.) from the `handlers` package to a new `internal/models` package to serve as a single source of truth for data contracts and improve the dependency graph.
    *   **Externalize Configuration:** Move hardcoded lists (e.g., language keywords in `extractor.go`, documentation file extensions in `read.go`) to a separate configuration file (e.g., `config/languages.json`) to be loaded at startup. This improves maintainability.
    *   **Centralize Defaults:** Move the default provider timeout from individual provider clients into the main application configuration (`config.go`).

*   **Testing Strategy:**
    *   **Unit Tests:** Each new handler (`submit_task`, `get_metadata`, `get_content`, `write_file`) will have comprehensive unit tests with mocked dependencies (storage, provider).
    *   **Integration Tests:** Create integration tests that cover the full, multi-step workflows (e.g., `submit -> metadata -> write`, `submit -> content`).
    *   **Regression Tests:** During the migration phase, maintain tests for the old, deprecated API to ensure it remains functional until removal.

*   **Documentation Updates:**
    *   The `api-reference.md` will be completely rewritten to reflect the new 4-tool architecture, including the new names, parameters, structured responses, and workflow examples.
    *   All internal code comments and `README.md` will be updated to reference the new tools.

## 5. Implementation Checklist

This checklist provides an ordered, file-by-file guide for implementation.

1.  **Project Structure & Models:**
    *   [ ] Create a new package `internal/models`.
    *   [ ] Move `Extraction`, `CodeBlock`, `GenerateRequest`, `StreamChunk` from `handlers/invoke.go` to `internal/models/models.go`.
    *   [ ] In `internal/models/error.go`, create the `AsDelegateError` helper function.
    *   [ ] Define the new structured response structs for all four tools in `internal/models/responses.go`.
    *   [ ] Define the new structured error struct in `internal/models/responses.go`.

2.  **Configuration:**
    *   [ ] Create `config/languages.json` and populate it with the keyword/extension lists from `extractor.go` and `read.go`.
    *   [ ] Update the application startup logic to load this new config file.
    *   [ ] Move the default provider timeout to the main `config.go` and remove it from provider clients.

3.  **Implement New Handlers:**
    *   [ ] **`delegate_submit_task`**: Rename `handlers/invoke.go` to `handlers/submit_task.go` and `InvokeHandler` to `SubmitTaskHandler`. Update its response to the simple `{"output_id": "..."}` format.
    *   [ ] **`delegate_get_output_metadata`**: Rename `handlers/check.go` to `handlers/get_metadata.go` and `CheckHandler` to `GetMetadataHandler`. Modify its `Handle` method to build and return the new, detailed structured response, including the `content_analysis` block.
    *   [ ] **`delegate_write_output_to_file`**: Create `handlers/write_file.go`. Move the `writeToFile` logic and the `if req.Options.WriteTo != ""` block from the old `read.go` into a new `WriteFileHandler`. Ensure it returns the new structured success/error response. Add security comments.
    *   [ ] **`delegate_get_output_content`**: Create `handlers/get_content.go`. Move the remaining logic from the old `read.go` (content retrieval) into a new `GetContentHandler`. Ensure it returns the new structured response.

4.  **Update Core Logic:**
    *   [ ] In `main.go`, unregister the old handlers and register the four new handlers with the router.
    *   [ ] Refactor all handlers to use the new structured error format and the `AsDelegateError` helper.
    *   [ ] Update the extractor and other relevant parts to use the new externalized language configuration.

5.  **Testing:**
    *   [ ] Write unit tests for `SubmitTaskHandler`.
    *   [ ] Write unit tests for `GetMetadataHandler`, covering cases with and without multiple blocks.
    *   [ ] Write unit tests for `WriteFileHandler`, including security path traversal checks.
    *   [ ] Write unit tests for `GetContentHandler`.
    *   [ ] Write new integration tests for the common workflows.

6.  **Migration & Deprecation (See Section 6):**
    *   [ ] Re-introduce the old handlers (`InvokeHandler`, `CheckHandler`, `ReadHandler`).
    *   [ ] Modify the old handlers to simply call the new handlers' logic.
    *   [ ] In the response of each old handler, add a `deprecation_warning` field.

7.  **Documentation:**
    *   [ ] Rewrite `api-reference.md` from scratch to document the new 4-tool API.
    *   [ ] Update `README.md` with the new workflow.
    *   [ ] Review and update all code comments.

## 6. Migration Strategy

To ensure zero downtime and a smooth transition for existing users, we will follow a three-phase migration path.

*   **Phase 1: Additive Introduction (v2.0)**
    1.  Implement and deploy the four new tools (`delegate_submit_task`, `delegate_get_output_metadata`, `delegate_get_output_content`, `delegate_write_output_to_file`).
    2.  The original three tools (`delegate_invoke`, `delegate_check`, `delegate_read`) will remain fully functional. Internally, they can be refactored to call the new handlers to reduce code duplication.
    3.  Update all documentation, examples, and agent-facing materials to **exclusively** feature the new 4-tool API. The old API will not be mentioned in any new documentation.

*   **Phase 2: Deprecation (v2.1)**
    1.  Modify the original three tools to include a `deprecation_warning` field in their JSON responses.
    2.  The warning message will state that the tool is deprecated, will be removed in a future version, and should point to its new replacement(s).
    3.  **Example `delegate_read` response:**
        ```json
        {
          "content": "...",
          "deprecation_warning": "The 'delegate_read' tool is deprecated and will be removed in v3.0. Please use 'delegate_get_output_content' or 'delegate_write_output_to_file' instead."
        }
        ```

*   **Phase 3: Removal (v3.0)**
    1.  In a future major version release, completely remove the old `delegate_invoke`, `delegate_check`, and `delegate_read` handlers and their routes from the codebase.

## 7. Success Metrics

We will measure the success of this refactoring based on the following metrics:

*   **Adoption Rate:**
    *   **Metric:** Percentage of API calls using the new toolset vs. the old, deprecated toolset.
    *   **Target:** 95% of traffic on the new API within 2 months of release.
*   **Error Rate Reduction:**
    *   **Metric:** A decrease in logs for `INVALID_REQUEST` errors, particularly those related to `delegate_read` misuse (e.g., parsing string warnings, incorrect parameter combinations).
    *   **Target:** 50% reduction in such errors within 1 month.
*   **Developer/Agent Experience (Qualitative):**
    *   **Metric:** Feedback from internal teams and AI agent trainers.
    *   **Target:** Feedback should indicate the new API is "easier to reason about," "more predictable," and "requires less complex logic to use correctly."
*   **Implementation Velocity:**
    *   **Metric:** Time required for new developers or agents to successfully implement a standard workflow.
    *   **Target:** A noticeable decrease in onboarding time and support questions related to the Delegate API.
*   **Deprecation Success:**
    *   **Metric:** Successful removal of the old API in a future major version without significant user disruption.
    *   **Target:** The v3.0 release (removing the old API) is a non-event for active users.