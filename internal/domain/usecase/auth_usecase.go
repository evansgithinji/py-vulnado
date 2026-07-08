package usecase

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
)

// AuthUseCase handles authentication operations
type AuthUseCase struct {
	userRepo repository.UserRepository
}

// NewAuthUseCase creates a new AuthUseCase
func NewAuthUseCase(userRepo repository.UserRepository) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo}
}

// Login authenticates a user
// VULNERABLE: Passes unsanitized credentials to repository (SQLi)
func (uc *AuthUseCase) Login(username, password string) (*entity.User, error) {
	return uc.userRepo.FindByCredentials(username, password)
}

// SearchUsers searches for users
// VULNERABLE: Passes unsanitized search term to repository (SQLi)
func (uc *AuthUseCase) SearchUsers(term string) ([]entity.User, error) {
	return uc.userRepo.Search(term)
}

// UpdateProfile updates a user's profile
// VULNERABLE: Allows privilege escalation via is_admin parameter
func (uc *AuthUseCase) UpdateProfile(userID int64, email string, isAdmin bool) error {
	return uc.userRepo.UpdateProfile(userID, email, isAdmin)
}

// GetUser retrieves a user by ID
func (uc *AuthUseCase) GetUser(userID int64) (*entity.User, error) {
	return uc.userRepo.FindByID(userID)
}
