package persistence

import (
	"database/sql"
	"time"

	"goapp/internal/domain/entity"
)

// SQLiteMessageRepository implements MessageRepository with SQLite
type SQLiteMessageRepository struct {
	db *sql.DB
}

// NewSQLiteMessageRepository creates a new SQLiteMessageRepository
func NewSQLiteMessageRepository(db *sql.DB) *SQLiteMessageRepository {
	return &SQLiteMessageRepository{db: db}
}

// GetAll returns all messages
// VULNERABLE: Returns raw messages without sanitization (for stored XSS display)
func (r *SQLiteMessageRepository) GetAll() ([]entity.Message, error) {
	rows, err := r.db.Query("SELECT id, content, author, created_at FROM messages ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message
	for rows.Next() {
		var msg entity.Message
		var createdAt string
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.Author, &createdAt); err != nil {
			return nil, err
		}
		msg.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		messages = append(messages, msg)
	}

	return messages, nil
}

// Create creates a new message
// VULNERABLE: Message content stored without sanitization (stored XSS)
func (r *SQLiteMessageRepository) Create(message *entity.Message) error {
	result, err := r.db.Exec(
		"INSERT INTO messages (content, author, created_at) VALUES (?, ?, ?)",
		message.Content, // VULNERABLE: Raw content stored
		message.Author,
		message.CreatedAt.Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	message.ID = id

	return nil
}
