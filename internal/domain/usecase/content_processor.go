package usecase

import "goapp/internal/domain/entity"

// ContentProcessor processes user-submitted content
type ContentProcessor struct{}

// NewContentProcessor creates a new ContentProcessor
func NewContentProcessor() *ContentProcessor {
	return &ContentProcessor{}
}

// ProcessContent processes content according to policy
// VULNERABLE: No HTML sanitization - passes through raw HTML/script tags
func (p *ContentProcessor) ProcessContent(content string, policy *entity.ContentPolicy) string {
	if policy.MaxLength > 0 && len(content) > policy.MaxLength {
		content = content[:policy.MaxLength]
	}
	// Fake "processing" - doesn't actually sanitize
	return content
}
