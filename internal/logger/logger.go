package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Level represents log level
type Level string

const (
	DebugLevel Level = "debug"
	InfoLevel  Level = "info"
	WarnLevel  Level = "warn"
	ErrorLevel Level = "error"
)

// Logger provides structured logging to stderr
type Logger struct {
	component string
	level     Level
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     Level                  `json:"level"`
	Component string                 `json:"component"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// New creates a new logger for a component
func New(component string, level Level) *Logger {
	return &Logger{
		component: component,
		level:     level,
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, data ...map[string]interface{}) {
	if l.shouldLog(DebugLevel) {
		l.log(DebugLevel, message, data...)
	}
}

// Info logs an info message
func (l *Logger) Info(message string, data ...map[string]interface{}) {
	if l.shouldLog(InfoLevel) {
		l.log(InfoLevel, message, data...)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(message string, data ...map[string]interface{}) {
	if l.shouldLog(WarnLevel) {
		l.log(WarnLevel, message, data...)
	}
}

// Error logs an error message
func (l *Logger) Error(message string, data ...map[string]interface{}) {
	if l.shouldLog(ErrorLevel) {
		l.log(ErrorLevel, message, data...)
	}
}

// Fatal logs an error message and exits
func (l *Logger) Fatal(message string, data ...map[string]interface{}) {
	l.log(ErrorLevel, message, data...)
	os.Exit(1)
}

// log writes a log entry to stderr
func (l *Logger) log(level Level, message string, data ...map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Component: l.component,
		Message:   message,
	}

	if len(data) > 0 && data[0] != nil {
		entry.Data = data[0]
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, `{"level":"error","message":"Failed to marshal log entry: %v"}`+"\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "%s\n", jsonData)
}

// shouldLog determines if a message should be logged based on level
func (l *Logger) shouldLog(level Level) bool {
	levelOrder := map[Level]int{
		DebugLevel: 0,
		InfoLevel:  1,
		WarnLevel:  2,
		ErrorLevel: 3,
	}

	return levelOrder[level] >= levelOrder[l.level]
}

// ParseLevel converts a string to a log level
func ParseLevel(s string) Level {
	switch s {
	case "debug":
		return DebugLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}