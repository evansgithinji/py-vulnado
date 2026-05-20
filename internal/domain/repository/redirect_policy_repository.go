package repository

import "goapp/internal/domain/entity"

// RedirectPolicyRepository stores redirect policies
type RedirectPolicyRepository struct {
	policies map[string]*entity.RedirectPolicy
}

// NewRedirectPolicyRepository creates a new RedirectPolicyRepository
func NewRedirectPolicyRepository() *RedirectPolicyRepository {
	return &RedirectPolicyRepository{
		policies: map[string]*entity.RedirectPolicy{
			"auth": {
				Name:           "auth",
				AllowedDomains: nil,   // VULNERABLE: no domain whitelist
				CheckExternal:  false, // VULNERABLE
			},
			"admin": {
				Name:           "admin",
				AllowedDomains: nil,
				CheckExternal:  false,
			},
		},
	}
}

// GetPolicy retrieves a redirect policy by name
func (r *RedirectPolicyRepository) GetPolicy(name string) *entity.RedirectPolicy {
	if policy, ok := r.policies[name]; ok {
		return policy
	}
	return r.policies["auth"]
}
