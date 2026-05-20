package repository

import "goapp/internal/domain/entity"

// UserRepository defines the interface for user data operations
// VULNERABLE: Implementations intentionally use unsafe SQL patterns
type UserRepository interface {
	// FindByCredentials authenticates a user (SQLi vulnerable)
	FindByCredentials(username, password string) (*entity.User, error)

	// Search searches users by term (SQLi vulnerable)
	Search(term string) ([]entity.User, error)

	// FindByID finds user by ID
	FindByID(id int64) (*entity.User, error)

	// UpdateProfile updates user profile (SQLi vulnerable, privilege escalation)
	UpdateProfile(userID int64, email string, isAdmin bool) error

	// Create creates a new user
	Create(user *entity.User) error
}
