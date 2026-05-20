package usecase

import (
	"fmt"
	"log"
	"sync"
	"time"

	"goapp/internal/domain/entity"
)

// LogEntry represents a single log entry in the log store.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// AuditEventEnricher enriches audit events with additional metadata.
type AuditEventEnricher struct{}

// NewAuditEventEnricher creates a new AuditEventEnricher.
func NewAuditEventEnricher() *AuditEventEnricher {
	return &AuditEventEnricher{}
}

// Enrich adds metadata to an audit event.
func (e *AuditEventEnricher) Enrich(event *entity.AuditEvent) *entity.AuditEvent {
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	if event.Metadata == nil {
		event.Metadata = make(map[string]string)
	}
	event.Metadata["enriched"] = "true"
	return event
}

// LogFormatter formats audit events into log messages.
type LogFormatter struct{}

// NewLogFormatter creates a new LogFormatter.
func NewLogFormatter() *LogFormatter {
	return &LogFormatter{}
}

// FormatLoginEvent formats a login audit event.
// VULNERABLE: Log Injection (CWE-117) - user input written to log format string without sanitization
func (f *LogFormatter) FormatLoginEvent(policy *entity.AuditPolicy, event *entity.AuditEvent) string {
	// VULNERABLE: newline injection allows log forging
	return fmt.Sprintf(policy.FormatTemplate, event.Actor, event.Details)
}

// FormatSearchEvent formats a search audit event.
// VULNERABLE: Log Injection (CWE-117) - user input written to log format string without sanitization
func (f *LogFormatter) FormatSearchEvent(policy *entity.AuditPolicy, event *entity.AuditEvent) string {
	// VULNERABLE: newline injection allows log forging
	return fmt.Sprintf(policy.FormatTemplate, event.Details)
}

// LogWriter writes formatted messages to the system log.
type LogWriter struct{}

// NewLogWriter creates a new LogWriter.
func NewLogWriter() *LogWriter {
	return &LogWriter{}
}

// Write writes a message to the system log.
func (w *LogWriter) Write(message string) {
	log.Println(message)
}

// LogStorageAdapter stores log entries in a circular buffer.
type LogStorageAdapter struct {
	mu      sync.Mutex
	entries []LogEntry
	maxSize int
}

// NewLogStorageAdapter creates a new LogStorageAdapter with the specified buffer size.
func NewLogStorageAdapter(maxSize int) *LogStorageAdapter {
	return &LogStorageAdapter{
		entries: make([]LogEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Store adds a log entry to the circular buffer.
func (s *LogStorageAdapter) Store(level, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
	}

	if len(s.entries) >= s.maxSize {
		s.entries = s.entries[1:]
	}
	s.entries = append(s.entries, entry)
}

// GetAll returns a copy of all stored log entries.
func (s *LogStorageAdapter) GetAll() []LogEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]LogEntry, len(s.entries))
	copy(result, s.entries)
	return result
}
