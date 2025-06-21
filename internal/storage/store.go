package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/christianwissmann85/delegate/internal/models"
)

// FileStore implements file-based storage
type FileStore struct {
	baseDir string
}

// NewFileStore creates a new file store
func NewFileStore(baseDir string) (*FileStore, error) {
	// Ensure base directory exists
	outputDir := filepath.Join(baseDir, "outputs")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	return &FileStore{
		baseDir: baseDir,
	}, nil
}

// Save persists an output to disk
func (s *FileStore) Save(output *models.Output) error {
	// TODO: Implement save logic
	return fmt.Errorf("not implemented")
}

// Get retrieves an output by ID
func (s *FileStore) Get(id string) (*models.Output, error) {
	// TODO: Implement get logic
	return nil, fmt.Errorf("not implemented")
}

// Delete removes an output
func (s *FileStore) Delete(id string) error {
	// TODO: Implement delete logic
	return fmt.Errorf("not implemented")
}

// ListOlderThan returns IDs of outputs older than the given age
func (s *FileStore) ListOlderThan(age time.Duration) ([]string, error) {
	// TODO: Implement list logic
	return nil, fmt.Errorf("not implemented")
}

// GetOutputPath returns the path for an output file
func (s *FileStore) GetOutputPath(id string) string {
	return filepath.Join(s.baseDir, "outputs", id+".json")
}

// TempDir returns the temporary directory for streaming
func (s *FileStore) TempDir() string {
	return filepath.Join(s.baseDir, "tmp")
}