package usecase

import (
	"crypto/md5"
	"fmt"
	"time"

	"goapp/internal/domain/entity"
)

// SessionManager manages user sessions
type SessionManager struct{}

// NewSessionManager creates a new SessionManager
func NewSessionManager() *SessionManager {
	return &SessionManager{}
}

// CreateSession creates a new session token for a user
// VULNERABLE: Weak token generation - predictable hash
func (m *SessionManager) CreateSession(user *entity.UserCredential) string {
	tokenData := fmt.Sprintf("%s:%d", user.Username, time.Now().UnixNano())
	hash := md5.Sum([]byte(tokenData))
	return fmt.Sprintf("%x", hash)
}
