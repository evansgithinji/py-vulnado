package entity

// RedirectRequest represents a redirect request DTO
type RedirectRequest struct {
	TargetURL string
	Context   string
}

// RedirectPolicy defines the policy for redirects
type RedirectPolicy struct {
	Name           string
	AllowedDomains []string // VULNERABLE: no domain whitelist
	CheckExternal  bool     // VULNERABLE: doesn't validate external URLs
}
