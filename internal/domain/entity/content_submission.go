package entity

// ContentSubmission represents a user content submission DTO
type ContentSubmission struct {
	Content     string
	Author      string
	ContentType string
	Context     string
}

// ContentPolicy defines the policy for content processing
type ContentPolicy struct {
	Name      string
	AllowHTML bool // VULNERABLE: allows HTML
	Sanitize  bool // VULNERABLE: no sanitization
	MaxLength int
}
