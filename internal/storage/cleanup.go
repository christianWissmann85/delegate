package storage

import (
	"context"
	"time"

	"github.com/christianwissmann85/delegate/internal/logger"
)

// Cleaner handles periodic cleanup of old outputs
type Cleaner struct {
	store    *FileStore
	interval time.Duration
	maxAge   time.Duration
	logger   *logger.Logger
}

// NewCleaner creates a new cleaner
func NewCleaner(store *FileStore, interval, maxAge time.Duration) *Cleaner {
	return &Cleaner{
		store:    store,
		interval: interval,
		maxAge:   maxAge,
		logger:   logger.New("storage.cleaner", logger.InfoLevel),
	}
}

// Start begins the cleanup routine
func (c *Cleaner) Start(ctx context.Context) {
	c.logger.Info("Starting cleanup routine", map[string]interface{}{
		"interval": c.interval.String(),
		"max_age":  c.maxAge.String(),
	})

	// Run initial cleanup after a short delay
	time.AfterFunc(30*time.Second, c.cleanup)

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Stopping cleanup routine")
			return
		case <-ticker.C:
			c.cleanup()
		}
	}
}

// cleanup removes old outputs
func (c *Cleaner) cleanup() {
	start := time.Now()
	
	ids, err := c.store.ListOlderThan(c.maxAge)
	if err != nil {
		c.logger.Error("Failed to list old outputs", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if len(ids) == 0 {
		c.logger.Debug("No old outputs to clean up")
		return
	}

	deleted := 0
	failed := 0

	for _, id := range ids {
		if err := c.store.Delete(id); err != nil {
			c.logger.Warn("Failed to delete output", map[string]interface{}{
				"id":    id,
				"error": err.Error(),
			})
			failed++
		} else {
			deleted++
		}
	}

	c.logger.Info("Cleanup completed", map[string]interface{}{
		"found":    len(ids),
		"deleted":  deleted,
		"failed":   failed,
		"duration": time.Since(start).String(),
	})
}

// DefaultCleanupConfig returns the default cleanup configuration
func DefaultCleanupConfig() (interval, maxAge time.Duration) {
	return 1 * time.Hour, 24 * time.Hour
}