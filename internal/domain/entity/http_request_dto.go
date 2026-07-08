package entity

// HttpRequestDto represents an outbound HTTP request DTO
type HttpRequestDto struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    string
	Timeout int
}

// UrlPolicy defines the policy for outbound URL access
type UrlPolicy struct {
	Name           string
	BlockedHosts   []string // VULNERABLE: no blocked hosts configured
	AllowedSchemes []string
	CheckInternal  bool // VULNERABLE: doesn't block internal IPs
}

// PreparedRequest represents a validated/prepared HTTP request
type PreparedRequest struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    string
	Timeout int
	Policy  *UrlPolicy
}
