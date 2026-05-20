package service

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// AuthenticationService orchestrates authentication with deep call graph
type AuthenticationService struct {
	credentialRepo *repository.CredentialRepository
	validator      *usecase.CredentialValidator
	sessionManager *usecase.SessionManager
}

// NewAuthenticationService creates a new AuthenticationService
func NewAuthenticationService(
	credentialRepo *repository.CredentialRepository,
	validator *usecase.CredentialValidator,
	sessionManager *usecase.SessionManager,
) *AuthenticationService {
	return &AuthenticationService{
		credentialRepo: credentialRepo,
		validator:      validator,
		sessionManager: sessionManager,
	}
}

// Authenticate authenticates a user through the deep call graph
// Deep call graph for Broken Auth (depth 5):
// Handler -> AuthenticationService.Authenticate()
//   -> CredentialRepository.FindUser() -> UserCredential (VULNERABLE: returns password)
//   -> CredentialValidator.Validate() -> boolean (VULNERABLE: timing attack)
//   -> SessionManager.CreateSession() -> token (VULNERABLE: weak token)
func (s *AuthenticationService) Authenticate(request *entity.AuthRequest) map[string]interface{} {
	credential := s.credentialRepo.FindUser(request.Username)
	if credential == nil {
		return nil
	}

	isValid := s.validator.Validate(request.Password, credential)
	if !isValid {
		return nil
	}

	token := s.sessionManager.CreateSession(credential)
	return map[string]interface{}{
		"user_id":  credential.UserID,
		"username": credential.Username,
		"email":    credential.Email,
		"is_admin": credential.IsAdmin,
		"token":    token,
	}
}
