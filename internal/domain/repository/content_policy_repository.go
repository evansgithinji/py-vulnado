package repository

import "goapp/internal/domain/entity"

// ContentPolicyRepository stores content processing policies
type ContentPolicyRepository struct {
	policies map[string]*entity.ContentPolicy
}

// NewContentPolicyRepository creates a new ContentPolicyRepository
func NewContentPolicyRepository() *ContentPolicyRepository {
	return &ContentPolicyRepository{
		policies: map[string]*entity.ContentPolicy{
			"html": {
				Name:      "html",
				AllowHTML: true,  // VULNERABLE
				Sanitize:  false, // VULNERABLE
				MaxLength: 10000,
			},
			"message": {
				Name:      "message",
				AllowHTML: true,
				Sanitize:  false,
				MaxLength: 5000,
			},
			"review": {
				Name:      "review",
				AllowHTML: true,
				Sanitize:  false,
				MaxLength: 2000,
			},
		},
	}
}

// GetPolicy retrieves a content policy by name
func (r *ContentPolicyRepository) GetPolicy(name string) *entity.ContentPolicy {
	if policy, ok := r.policies[name]; ok {
		return policy
	}
	return r.policies["html"]
}
