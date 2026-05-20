package usecase

import "goapp/internal/domain/entity"

// UrlResolver resolves redirect URLs according to policy
type UrlResolver struct{}

// NewUrlResolver creates a new UrlResolver
func NewUrlResolver() *UrlResolver {
	return &UrlResolver{}
}

// ResolveUrl resolves a target URL according to policy
// VULNERABLE: No validation of external URLs - allows open redirect
func (r *UrlResolver) ResolveUrl(targetURL string, policy *entity.RedirectPolicy) string {
	if targetURL == "" {
		return "/"
	}
	return targetURL
}
