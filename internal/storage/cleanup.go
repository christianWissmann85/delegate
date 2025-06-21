package storage

import (
	"context"
	"log"
	"time"
)

// Cleaner handles periodic cleanup of old outputs
type Cleaner struct {
	store    *FileStore
	interval time.Duration
	maxAge   time.Duration
}

// NewCleaner creates a new cleaner
func NewCleaner(store *FileStore, interval, maxAge time.Duration) *Cleaner {
	return &Cleaner{
		store:    store,
		interval: interval,
		maxAge:   maxAge,
	}
}

// Start begins the cleanup routine
func (c *Cleaner) Start(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	// Run initial cleanup
	c.cleanup()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.cleanup()
		}
	}
}

// cleanup removes old outputs
func (c *Cleaner) cleanup() {
	ids, err := c.store.ListOlderThan(c.maxAge)
	if err != nil {
		log.Printf("Failed to list old outputs: %v", err)
		return
	}

	for _, id := range ids {
		if err := c.store.Delete(id); err != nil {
			log.Printf("Failed to delete output %s: %v", id, err)
		} else {
			log.Printf("Deleted old output: %s", id)
		}
	}
}