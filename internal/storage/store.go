package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/christianwissmann85/delegate/internal/logger"
	"github.com/christianwissmann85/delegate/internal/models"
)

// FileStore implements file-based storage
type FileStore struct {
	baseDir string
	logger  *logger.Logger
}

// NewFileStore creates a new file store
func NewFileStore(baseDir string) (*FileStore, error) {
	// Ensure base directory exists
	outputDir := filepath.Join(baseDir, "outputs")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("create output dir: %w", err)
	}

	// Ensure temp directory exists
	tempDir := filepath.Join(baseDir, "tmp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}

	return &FileStore{
		baseDir: baseDir,
		logger:  logger.New("storage", logger.InfoLevel),
	}, nil
}

// GenerateID creates a new output ID based on timestamp
func (s *FileStore) GenerateID() string {
	return fmt.Sprintf("out_%s", time.Now().UTC().Format("20060102_150405"))
}

// Save persists an output to disk atomically
func (s *FileStore) Save(output *models.Output) error {
	if output.ID == "" {
		output.ID = s.GenerateID()
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal output: %w", err)
	}

	// Write to temp file first (atomic write)
	tempPath := filepath.Join(s.TempDir(), output.ID+".tmp")
	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}

	// Move to final location atomically
	finalPath := s.GetOutputPath(output.ID)
	if err := os.Rename(tempPath, finalPath); err != nil {
		// Clean up temp file on error
		os.Remove(tempPath)
		return fmt.Errorf("move to final location: %w", err)
	}

	s.logger.Info("Saved output", map[string]interface{}{
		"id":   output.ID,
		"size": len(data),
		"path": finalPath,
	})

	return nil
}

// Get retrieves an output by ID
func (s *FileStore) Get(id string) (*models.Output, error) {
	// Sanitize ID to prevent path traversal
	if strings.Contains(id, "/") || strings.Contains(id, "..") {
		return nil, fmt.Errorf("invalid output ID: %s", id)
	}

	path := s.GetOutputPath(id)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("output not found: %s", id)
		}
		return nil, fmt.Errorf("read output file: %w", err)
	}

	var output models.Output
	if err := json.Unmarshal(data, &output); err != nil {
		return nil, fmt.Errorf("unmarshal output: %w", err)
	}

	return &output, nil
}

// Delete removes an output
func (s *FileStore) Delete(id string) error {
	// Sanitize ID
	if strings.Contains(id, "/") || strings.Contains(id, "..") {
		return fmt.Errorf("invalid output ID: %s", id)
	}

	path := s.GetOutputPath(id)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return nil // Already gone, that's fine
		}
		return fmt.Errorf("remove output file: %w", err)
	}

	s.logger.Info("Deleted output", map[string]interface{}{
		"id": id,
	})

	return nil
}

// ListOlderThan returns IDs of outputs older than the given age
func (s *FileStore) ListOlderThan(age time.Duration) ([]string, error) {
	outputDir := filepath.Join(s.baseDir, "outputs")
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return nil, fmt.Errorf("read output directory: %w", err)
	}

	cutoff := time.Now().Add(-age)
	var oldIDs []string

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			s.logger.Warn("Failed to get file info", map[string]interface{}{
				"file":  entry.Name(),
				"error": err.Error(),
			})
			continue
		}

		if info.ModTime().Before(cutoff) {
			// Extract ID from filename (remove .json suffix)
			id := strings.TrimSuffix(entry.Name(), ".json")
			oldIDs = append(oldIDs, id)
		}
	}

	return oldIDs, nil
}

// GetOutputPath returns the path for an output file
func (s *FileStore) GetOutputPath(id string) string {
	return filepath.Join(s.baseDir, "outputs", id+".json")
}

// TempDir returns the temporary directory for streaming
func (s *FileStore) TempDir() string {
	return filepath.Join(s.baseDir, "tmp")
}

// CreateTempFile creates a temporary file for streaming
func (s *FileStore) CreateTempFile(prefix string) (*os.File, error) {
	tempDir := s.TempDir()
	return os.CreateTemp(tempDir, prefix+"_*.tmp")
}

// SaveStream saves a stream to a temporary file and returns the path
func (s *FileStore) SaveStream(reader io.Reader, prefix string) (string, error) {
	tempFile, err := s.CreateTempFile(prefix)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer tempFile.Close()

	written, err := io.Copy(tempFile, reader)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("copy stream: %w", err)
	}

	s.logger.Debug("Saved stream to temp file", map[string]interface{}{
		"path": tempFile.Name(),
		"size": written,
	})

	return tempFile.Name(), nil
}