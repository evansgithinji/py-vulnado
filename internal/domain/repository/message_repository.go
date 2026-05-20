package repository

import "goapp/internal/domain/entity"

// MessageRepository defines the interface for message data operations
// VULNERABLE: Used for XSS demonstration
type MessageRepository interface {
	// GetAll returns all messages (for stored XSS display)
	GetAll() ([]entity.Message, error)

	// Create creates a new message (stored XSS vulnerable)
	Create(message *entity.Message) error
}
