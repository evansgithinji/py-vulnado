package usecase

import "goapp/internal/domain/entity"

// CredentialValidator validates user credentials
type CredentialValidator struct{}

// NewCredentialValidator creates a new CredentialValidator
func NewCredentialValidator() *CredentialValidator {
	return &CredentialValidator{}
}

// Validate validates a password against stored credentials
// VULNERABLE: Timing attack - simple string comparison
// VULNERABLE: No rate limiting
func (v *CredentialValidator) Validate(password string, credential *entity.UserCredential) bool {
	if credential == nil {
		return false
	}
	return credential.Password == password
}
