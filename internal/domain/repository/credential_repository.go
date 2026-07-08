package repository

import (
	"database/sql"
	"fmt"

	"goapp/internal/domain/entity"
)

// CredentialRepository wraps database for credential lookup operations
type CredentialRepository struct {
	db *sql.DB
}

// NewCredentialRepository creates a new CredentialRepository
func NewCredentialRepository(db *sql.DB) *CredentialRepository {
	return &CredentialRepository{db: db}
}

// FindUser looks up user credentials by username
// VULNERABLE: uses raw SQL string concatenation; returns full credential including password
func (r *CredentialRepository) FindUser(username string) *entity.UserCredential {
	// VULNERABLE: SQL injection via string concatenation
	query := fmt.Sprintf(
		"SELECT id, username, password, email, is_admin FROM users WHERE username = '%s'",
		username,
	)

	var cred entity.UserCredential
	err := r.db.QueryRow(query).Scan(
		&cred.UserID, &cred.Username, &cred.Password, &cred.Email, &cred.IsAdmin,
	)
	if err != nil {
		return nil
	}

	return &cred
}
