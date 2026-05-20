package repository

import "goapp/internal/domain/entity"

// UrlPolicyRepository stores URL access policies
type UrlPolicyRepository struct {
	policies map[string]*entity.UrlPolicy
}

// NewUrlPolicyRepository creates a new UrlPolicyRepository
func NewUrlPolicyRepository() *UrlPolicyRepository {
	return &UrlPolicyRepository{
		policies: map[string]*entity.UrlPolicy{
			"fetch": {
				Name:           "fetch",
				BlockedHosts:   nil, // VULNERABLE: no blocked hosts
				AllowedSchemes: []string{"http", "https"},
				CheckInternal:  false, // VULNERABLE
			},
			"webhook": {
				Name:           "webhook",
				BlockedHosts:   nil,
				AllowedSchemes: []string{"http", "https"},
				CheckInternal:  false,
			},
			"proxy": {
				Name:           "proxy",
				BlockedHosts:   nil,
				AllowedSchemes: []string{"http", "https"},
				CheckInternal:  false,
			},
			"remote_file": {
				Name:           "remote_file",
				BlockedHosts:   nil,
				AllowedSchemes: []string{"http", "https", "file"},
				CheckInternal:  false,
			},
		},
	}
}

// GetPolicy retrieves a URL policy by name
func (r *UrlPolicyRepository) GetPolicy(name string) *entity.UrlPolicy {
	if policy, ok := r.policies[name]; ok {
		return policy
	}
	return r.policies["fetch"]
}
