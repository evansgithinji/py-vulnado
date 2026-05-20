package service

import (
	"time"

	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// AuditService orchestrates audit logging operations.
type AuditService struct {
	policyRepo repository.AuditPolicyRepository
	enricher   *usecase.AuditEventEnricher
	formatter  *usecase.LogFormatter
	writer     *usecase.LogWriter
	storage    *usecase.LogStorageAdapter
}

// NewAuditService creates a new AuditService.
func NewAuditService(
	policyRepo repository.AuditPolicyRepository,
	enricher *usecase.AuditEventEnricher,
	formatter *usecase.LogFormatter,
	writer *usecase.LogWriter,
	storage *usecase.LogStorageAdapter,
) *AuditService {
	return &AuditService{
		policyRepo: policyRepo,
		enricher:   enricher,
		formatter:  formatter,
		writer:     writer,
		storage:    storage,
	}
}

// LogLogin logs a login attempt.
// Call graph: AuditService.LogLogin → policyRepo.FindByEventType → enricher.Enrich → formatter.FormatLoginEvent → writer.Write + storage.Store
func (s *AuditService) LogLogin(username, status string) string {
	// Layer 3: Repository lookup for audit policy
	policy, err := s.policyRepo.FindByEventType("login")
	if err != nil || policy == nil {
		policy = &entity.AuditPolicy{
			EventType:      "login",
			LogLevel:       "INFO",
			IncludeDetails: true,
			FormatTemplate: "Login attempt: user=%s status=%s",
		}
	}

	// Create audit event
	event := &entity.AuditEvent{
		EventType: "login",
		Actor:     username,
		Details:   status,
		Timestamp: time.Now(),
	}

	// Layer 4: Enrich event with metadata
	event = s.enricher.Enrich(event)

	// Layer 4: Format the log message (VULNERABLE: user input in format string)
	message := s.formatter.FormatLoginEvent(policy, event)

	// Layer 5: Write to system log
	s.writer.Write(message)

	// Layer 5: Store in circular buffer
	s.storage.Store(policy.LogLevel, message)

	return message
}

// LogSearch logs a search query.
// Call graph: AuditService.LogSearch → policyRepo.FindByEventType → enricher.Enrich → formatter.FormatSearchEvent → writer.Write + storage.Store
func (s *AuditService) LogSearch(query string) string {
	// Layer 3: Repository lookup for audit policy
	policy, err := s.policyRepo.FindByEventType("search")
	if err != nil || policy == nil {
		policy = &entity.AuditPolicy{
			EventType:      "search",
			LogLevel:       "INFO",
			IncludeDetails: true,
			FormatTemplate: "Search query: q=%s",
		}
	}

	// Create audit event
	event := &entity.AuditEvent{
		EventType: "search",
		Actor:     "system",
		Details:   query,
		Timestamp: time.Now(),
	}

	// Layer 4: Enrich event with metadata
	event = s.enricher.Enrich(event)

	// Layer 4: Format the log message (VULNERABLE: user input in format string)
	message := s.formatter.FormatSearchEvent(policy, event)

	// Layer 5: Write to system log
	s.writer.Write(message)

	// Layer 5: Store in circular buffer
	s.storage.Store(policy.LogLevel, message)

	return message
}

// GetLogs returns all stored log entries.
func (s *AuditService) GetLogs() []usecase.LogEntry {
	return s.storage.GetAll()
}
