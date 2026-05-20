package entity

import "time"

// Message represents a message/comment in the system (for XSS demo)
type Message struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}
