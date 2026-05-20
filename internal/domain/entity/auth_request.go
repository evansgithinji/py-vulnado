package entity

// AuthRequest represents an authentication request DTO
type AuthRequest struct {
	Username string
	Password string
}

// UserCredential represents stored user credentials
type UserCredential struct {
	UserID   int64
	Username string
	Password string // VULNERABLE: stores plaintext password
	Email    string
	IsAdmin  bool
}
