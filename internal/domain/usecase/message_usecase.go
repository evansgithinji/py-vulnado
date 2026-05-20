package usecase

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"time"
)

// MessageUseCase handles message operations (for XSS demo)
type MessageUseCase struct {
	messageRepo repository.MessageRepository
}

// NewMessageUseCase creates a new MessageUseCase
func NewMessageUseCase(messageRepo repository.MessageRepository) *MessageUseCase {
	return &MessageUseCase{messageRepo: messageRepo}
}

// GetAllMessages retrieves all messages
// VULNERABLE: Returns raw messages without sanitization (stored XSS)
func (uc *MessageUseCase) GetAllMessages() ([]entity.Message, error) {
	return uc.messageRepo.GetAll()
}

// PostMessage creates a new message
// VULNERABLE: No sanitization of message content (stored XSS)
func (uc *MessageUseCase) PostMessage(content, author string) error {
	message := &entity.Message{
		Content:   content, // VULNERABLE: Raw content stored
		Author:    author,
		CreatedAt: time.Now(),
	}
	return uc.messageRepo.Create(message)
}
