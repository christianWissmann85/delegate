package analytics

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAnalyticsLogger(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "analytics_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create logger
	logger, err := New(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	// Log some events
	logger.Log(LogEntry{
		"event_type": "task_submitted",
		"model":      "gemini-2.5-flash",
		"duration_ms": 1234,
		"success":    true,
	})

	logger.Log(LogEntry{
		"event_type":  "file_written",
		"output_id":   "test_123",
		"bytes_written": 5000,
		"tokens_saved": 1250,
		"success":     true,
	})

	// Close the logger to flush all events
	logger.Close()

	// Check that a log file was created
	today := time.Now().UTC().Format("2006-01-02")
	logFile := filepath.Join(tmpDir, today+".jsonl")
	
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Expected log file %s to exist", logFile)
	}

	// Read the log file
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}

	// Check that both events were logged
	if len(content) == 0 {
		t.Error("Log file is empty")
	}

	// Simple content checks
	contentStr := string(content)
	if !contains(contentStr, "task_submitted") {
		t.Error("task_submitted event not found in log")
	}
	if !contains(contentStr, "file_written") {
		t.Error("file_written event not found in log")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || contains(s[1:], substr)))
}