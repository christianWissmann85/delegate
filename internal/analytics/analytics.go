// internal/analytics/analytics.go
package analytics

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	// logFileDateFormat defines the format for the daily log file names (e.g., "2023-10-27").
	logFileDateFormat = "2006-01-02"
	// logRetentionDays is the number of days to keep log files before deleting them.
	logRetentionDays = 30
	// logChannelBuffer is the size of the buffered channel for log entries.
	// If the buffer is full, new log entries will be dropped.
	logChannelBuffer = 1000
)

// LogEntry represents a single analytics event. It's a flexible map to accommodate different event types.
type LogEntry map[string]interface{}

// Logger handles writing analytics events to daily rotating log files.
// It is thread-safe and designed for high performance with non-blocking writes.
type Logger struct {
	logDir      string
	logChan     chan LogEntry
	wg          sync.WaitGroup
	mu          sync.Mutex
	currentFile *os.File
	currentDate string
	closed      bool
}

// New creates and starts a new analytics logger.
// It initializes the writer goroutine and performs an initial cleanup of old logs.
func New(logDir string) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %s: %w", logDir, err)
	}

	l := &Logger{
		logDir:  logDir,
		logChan: make(chan LogEntry, logChannelBuffer),
	}

	// Start the background writer goroutine.
	l.wg.Add(1)
	go l.writer()

	// Perform an initial cleanup of old logs in the background.
	go l.cleanupOldLogs()

	return l, nil
}

// Log sends an event to the analytics logger. This is a non-blocking call.
// If the logger's buffer is full, the entry is dropped to prevent blocking
// the calling application thread.
func (l *Logger) Log(entry LogEntry) {
	// Add a timestamp in UTC if one isn't already present.
	if _, ok := entry["timestamp"]; !ok {
		entry["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}

	// Non-blocking send to the channel.
	select {
	case l.logChan <- entry:
	default:
		// This message is printed to stderr for visibility into dropped logs.
		fmt.Fprintf(os.Stderr, "Analytics channel buffer is full. Dropping log entry: %v\n", entry)
	}
}

// writer is the single goroutine that consumes log entries from the channel and handles all file I/O.
func (l *Logger) writer() {
	defer l.wg.Done()

	for entry := range l.logChan {
		l.writeToFile(entry)
	}

	// After the channel is closed and all entries are processed,
	// ensure the current file handle is closed.
	l.mu.Lock()
	if l.currentFile != nil {
		l.currentFile.Close()
	}
	l.mu.Unlock()
}

// writeToFile handles the actual file writing and rotation logic.
// It is only called by the single writer goroutine.
func (l *Logger) writeToFile(entry LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	today := time.Now().UTC().Format(logFileDateFormat)
	if l.currentFile == nil || l.currentDate != today {
		if err := l.rotate(today); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to rotate analytics log: %v\n", err)
			// If rotation fails, we cannot write the log.
			return
		}
	}

	line, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal analytics entry: %v\n", err)
		return
	}

	// Append a newline to conform to the JSON Lines format.
	if _, err := l.currentFile.Write(append(line, '\n')); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to analytics log: %v\n", err)
	}
}

// rotate closes the current log file (if any) and opens a new one for the given date string.
// This method must be called with the mutex held.
func (l *Logger) rotate(dateStr string) error {
	if l.currentFile != nil {
		l.currentFile.Close()
	}

	filename := fmt.Sprintf("%s.jsonl", dateStr)
	path := filepath.Join(l.logDir, filename)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		l.currentFile = nil // Ensure we don't hold a reference to a failed open
		return fmt.Errorf("could not open log file %s: %w", path, err)
	}

	l.currentFile = file
	l.currentDate = dateStr

	// Schedule a cleanup to run in the background. This is efficient
	// as it only runs once per day upon the first log of that day.
	go l.cleanupOldLogs()

	return nil
}

// cleanupOldLogs scans the log directory and removes any log files
// that are older than the configured retention period.
func (l *Logger) cleanupOldLogs() {
	cutoff := time.Now().UTC().AddDate(0, 0, -logRetentionDays)

	files, err := os.ReadDir(l.logDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read log directory for cleanup: %v\n", err)
		return
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".jsonl") {
			continue
		}

		dateStr := strings.TrimSuffix(file.Name(), ".jsonl")
		fileDate, err := time.Parse(logFileDateFormat, dateStr)
		if err != nil {
			// Skip files with names that don't match the date format.
			continue
		}

		if fileDate.Before(cutoff) {
			path := filepath.Join(l.logDir, file.Name())
			if err := os.Remove(path); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to remove old log file %s: %v\n", path, err)
			}
		}
	}
}

// Close gracefully shuts down the logger. It closes the log channel and waits
// for the writer goroutine to process all buffered entries before returning.
func (l *Logger) Close() {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return
	}
	l.closed = true
	l.mu.Unlock()
	
	close(l.logChan)
	l.wg.Wait()
}