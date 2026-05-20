package repository

import "goapp/internal/domain/entity"

// AuditPolicyRepository provides access to audit policy configurations.
type AuditPolicyRepository interface {
	FindByEventType(eventType string) (*entity.AuditPolicy, error)
}

// InMemoryAuditPolicyRepository holds audit policies in memory.
type InMemoryAuditPolicyRepository struct {
	policies map[string]*entity.AuditPolicy
}

// NewInMemoryAuditPolicyRepository creates a new InMemoryAuditPolicyRepository with default policies.
func NewInMemoryAuditPolicyRepository() *InMemoryAuditPolicyRepository {
	policies := map[string]*entity.AuditPolicy{
		"login": {
			EventType:      "login",
			LogLevel:       "INFO",
			IncludeDetails: true,
			FormatTemplate: "Login attempt: user=%s status=%s",
		},
		"search": {
			EventType:      "search",
			LogLevel:       "INFO",
			IncludeDetails: true,
			FormatTemplate: "Search query: q=%s",
		},
	}
	return &InMemoryAuditPolicyRepository{policies: policies}
}

// FindByEventType returns the audit policy matching the given event type.
func (r *InMemoryAuditPolicyRepository) FindByEventType(eventType string) (*entity.AuditPolicy, error) {
	if policy, ok := r.policies[eventType]; ok {
		return policy, nil
	}
	// Default policy
	return &entity.AuditPolicy{
		EventType:      eventType,
		LogLevel:       "INFO",
		IncludeDetails: true,
		FormatTemplate: "%s",
	}, nil
}
