package storage

import "time"

// StorageOptions configures storage behavior
type StorageOptions struct {
	MaxFileSize   int64         // Maximum file size in bytes
	CleanupAge    time.Duration // Age after which files are deleted
	CleanupInterval time.Duration // How often to run cleanup
}

// StorageStats provides storage statistics
type StorageStats struct {
	TotalFiles   int
	TotalSize    int64
	OldestFile   time.Time
	NewestFile   time.Time
}