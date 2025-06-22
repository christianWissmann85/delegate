package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/christianwissmann85/delegate/internal/models"
)

func TestCleaner_Cleanup(t *testing.T) {
	tempDir := t.TempDir()

	store, err := NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create outputs with different ages
	oldOutput := &models.Output{ID: "out_old", Model: "old"}
	newOutput := &models.Output{ID: "out_new", Model: "new"}

	// Save outputs
	err = store.Save(oldOutput)
	if err != nil {
		t.Fatalf("Failed to save old output: %v", err)
	}
	err = store.Save(newOutput)
	if err != nil {
		t.Fatalf("Failed to save new output: %v", err)
	}

	// Make first output old
	oldPath := store.GetOutputPath(oldOutput.ID)
	oldTime := time.Now().Add(-25 * time.Hour)
	_ = os.Chtimes(oldPath, oldTime, oldTime)

	// Create cleaner with short intervals for testing
	cleaner := NewCleaner(store, 100*time.Millisecond, 24*time.Hour)

	// Run cleanup once
	cleaner.cleanup()

	// Verify old output is deleted
	_, err = store.Get(oldOutput.ID)
	if err == nil {
		t.Error("Old output should be deleted")
	}

	// Verify new output still exists
	_, err = store.Get(newOutput.ID)
	if err != nil {
		t.Error("New output should still exist")
	}
}

func TestCleaner_Start(t *testing.T) {
	tempDir := t.TempDir()

	store, err := NewFileStore(tempDir)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	// Create old output
	oldOutput := &models.Output{ID: "out_old", Model: "old"}
	err = store.Save(oldOutput)
	if err != nil {
		t.Fatalf("Failed to save output: %v", err)
	}

	// Make it old
	oldPath := store.GetOutputPath(oldOutput.ID)
	oldTime := time.Now().Add(-25 * time.Hour)
	_ = os.Chtimes(oldPath, oldTime, oldTime)

	// Create cleaner with very short intervals for testing
	cleaner := NewCleaner(store, 50*time.Millisecond, 24*time.Hour)

	// Start cleaner in background
	ctx, cancel := context.WithCancel(context.Background())
	go cleaner.Start(ctx)

	// Wait for at least one cleanup cycle
	time.Sleep(100 * time.Millisecond)

	// Verify old output is deleted
	_, err = store.Get(oldOutput.ID)
	if err == nil {
		t.Error("Old output should be deleted by background cleaner")
	}

	// Stop cleaner
	cancel()

	// Give it time to stop
	time.Sleep(50 * time.Millisecond)
}

func TestDefaultCleanupConfig(t *testing.T) {
	interval, maxAge := DefaultCleanupConfig()

	if interval != 1*time.Hour {
		t.Errorf("Expected 1 hour interval, got %v", interval)
	}

	if maxAge != 24*time.Hour {
		t.Errorf("Expected 24 hour max age, got %v", maxAge)
	}
}
