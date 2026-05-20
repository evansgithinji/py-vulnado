package entity

import "time"

// AuditEvent represents an auditable event in the system.
type AuditEvent struct {
	EventType string
	Actor     string
	Details   string
	Timestamp time.Time
	SourceIP  string
	Metadata  map[string]string
}

// AuditPolicy defines how audit events should be logged.
type AuditPolicy struct {
	EventType      string
	LogLevel       string
	IncludeDetails bool
	FormatTemplate string
}
