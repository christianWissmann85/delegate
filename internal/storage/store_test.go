package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/christianwissmann85/delegate/internal/models"
)

func TestFileStore_GenerateID(t *testing.T) {
	store := &FileStore{}
	
	id1 := store.GenerateID()
	if !strings.HasPrefix(id1, "out_") {
		t.Errorf("ID should start with 'out_', got: %s", id1)
	}
	
	// Wait at least 1 second to ensure different timestamp
	time.Sleep(1 * time.Second)
	
	id2 := store.GenerateID()
	if id1 == id2 {
		t.Error("Sequential IDs should be different")
	}
	
	// Verify format
	if len(id1) != 26 { // out_YYYYMMDD_HHMMSS_NNNNNN
		t.Errorf("ID should be 26 characters, got %d: %s", len(id1), id1)
	}
}

func TestFileStore_SaveAndGet(t *testing.T) {
	// Create temp directory for test
	tempDir := t.TempDir()
	
	store, err := NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Create test output
	output := &models.Output{
		Model:  "test-model",
		Prompt: "test prompt",
		Response: models.Response{
			Raw: "test response",
			Extracted: models.Extracted{
				Code: []models.ExtractedCode{
					{
						Language: "go",
						Content:  "func main() {}",
					},
				},
				Explanation: "test explanation",
			},
		},
		Metadata: models.Metadata{
			TotalBytes:      100,
			EstimatedTokens: 25,
		},
	}
	
	// Save output
	err = store.Save(output)
	if err != nil {
		t.Fatalf("Failed to save output: %v", err)
	}
	
	// Check that ID was generated
	if output.ID == "" {
		t.Error("Output ID should be generated")
	}
	
	// Get output back
	retrieved, err := store.Get(output.ID)
	if err != nil {
		t.Fatalf("Failed to get output: %v", err)
	}
	
	// Verify content
	if retrieved.Model != output.Model {
		t.Errorf("Model mismatch: got %s, want %s", retrieved.Model, output.Model)
	}
	if retrieved.Prompt != output.Prompt {
		t.Errorf("Prompt mismatch: got %s, want %s", retrieved.Prompt, output.Prompt)
	}
	if retrieved.Response.Raw != output.Response.Raw {
		t.Errorf("Response mismatch: got %s, want %s", retrieved.Response.Raw, output.Response.Raw)
	}
}

func TestFileStore_Delete(t *testing.T) {
	tempDir := t.TempDir()
	
	store, err := NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Create and save output
	output := &models.Output{
		Model:  "test-model",
		Prompt: "test prompt",
	}
	
	err = store.Save(output)
	if err != nil {
		t.Fatalf("Failed to save output: %v", err)
	}
	
	// Delete output
	err = store.Delete(output.ID)
	if err != nil {
		t.Fatalf("Failed to delete output: %v", err)
	}
	
	// Try to get deleted output
	_, err = store.Get(output.ID)
	if err == nil {
		t.Error("Should fail to get deleted output")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

func TestFileStore_ListOlderThan(t *testing.T) {
	tempDir := t.TempDir()
	
	store, err := NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Create outputs
	output1 := &models.Output{ID: "out_20250620_120000", Model: "old"}
	output2 := &models.Output{ID: "out_20250621_120000", Model: "new"}
	
	// Save outputs
	err = store.Save(output1)
	if err != nil {
		t.Fatalf("Failed to save output1: %v", err)
	}
	err = store.Save(output2)
	if err != nil {
		t.Fatalf("Failed to save output2: %v", err)
	}
	
	// Modify time of first output to be old
	oldPath := store.GetOutputPath(output1.ID)
	oldTime := time.Now().Add(-25 * time.Hour)
	_ = os.Chtimes(oldPath, oldTime, oldTime)
	
	// List old outputs
	oldIDs, err := store.ListOlderThan(24 * time.Hour)
	if err != nil {
		t.Fatalf("Failed to list old outputs: %v", err)
	}
	
	// Should find only the old output
	if len(oldIDs) != 1 {
		t.Errorf("Expected 1 old output, got %d", len(oldIDs))
	}
	if len(oldIDs) > 0 && oldIDs[0] != output1.ID {
		t.Errorf("Expected old output ID %s, got %s", output1.ID, oldIDs[0])
	}
}

func TestFileStore_PathTraversal(t *testing.T) {
	tempDir := t.TempDir()
	
	store, err := NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Try path traversal attacks
	badIDs := []string{
		"../../../etc/passwd",
		"out_../../secret",
		"out/../../secret",
		"/etc/passwd",
	}
	
	for _, badID := range badIDs {
		_, err := store.Get(badID)
		if err == nil {
			t.Errorf("Should reject bad ID: %s", badID)
		}
		if !strings.Contains(err.Error(), "invalid output ID") {
			t.Errorf("Expected 'invalid output ID' error for %s, got: %v", badID, err)
		}
		
		err = store.Delete(badID)
		if err == nil {
			t.Errorf("Should reject bad ID for delete: %s", badID)
		}
	}
}

func TestFileStore_AtomicWrite(t *testing.T) {
	tempDir := t.TempDir()
	
	store, err := NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	
	// Check that temp directory was created
	tempDirPath := filepath.Join(tempDir, "tmp")
	if _, err := os.Stat(tempDirPath); os.IsNotExist(err) {
		t.Error("Temp directory should be created")
	}
	
	// Save output and verify no temp files remain
	output := &models.Output{Model: "test"}
	err = store.Save(output)
	if err != nil {
		t.Fatalf("Failed to save output: %v", err)
	}
	
	// Check for leftover temp files
	entries, _ := os.ReadDir(tempDirPath)
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tmp") {
			t.Errorf("Found leftover temp file: %s", entry.Name())
		}
	}
}