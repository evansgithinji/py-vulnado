package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
)

// MessageHandler handles message-related HTTP requests (XSS demo)
type MessageHandler struct {
	contentService *service.ContentService
}

// NewMessageHandler creates a new MessageHandler
func NewMessageHandler(contentService *service.ContentService) *MessageHandler {
	return &MessageHandler{contentService: contentService}
}

// RegisterRoutes registers message routes
func (h *MessageHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/messages", h.GetMessages)
	r.POST("/api/messages", h.PostMessage)

	// HTML pages
	r.GET("/messages", h.MessagesPage)
	r.POST("/messages", h.PostMessagePage)
}

// GetMessages retrieves all messages
// VULNERABLE: Returns raw messages (stored XSS when displayed)
func (h *MessageHandler) GetMessages(c *gin.Context) {
	messages, err := h.contentService.GetAllMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// PostMessage creates a new message
// VULNERABLE: No sanitization (stored XSS) via ContentService deep call graph
func (h *MessageHandler) PostMessage(c *gin.Context) {
	content := c.PostForm("content")
	author := c.PostForm("author")

	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content required"})
		return
	}

	if author == "" {
		author = "Anonymous"
	}

	// VULNERABLE: Raw content stored via deep call graph
	submission := &entity.ContentSubmission{
		Content:     content,
		Author:      author,
		ContentType: "html",
		Context:     "message",
	}
	err := h.contentService.SubmitMessage(submission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Message posted"})
}

// MessagesPage serves the message board page
// VULNERABLE: Stored XSS - displays raw message content
func (h *MessageHandler) MessagesPage(c *gin.Context) {
	messages, _ := h.contentService.GetAllMessages()

	html := `
<!DOCTYPE html>
<html>
<head><title>Message Board - Vulnerable App</title></head>
<body>
	<h1>Message Board (Stored XSS Demo)</h1>
	<form method="POST" action="/messages">
		<label>Author: <input type="text" name="author"></label><br><br>
		<label>Message: <textarea name="content" rows="3" cols="50"></textarea></label><br><br>
		<button type="submit">Post Message</button>
	</form>
	<h2>Messages</h2>
	<div id="messages">
`

	for _, msg := range messages {
		// VULNERABLE: Stored XSS - content displayed without encoding
		html += fmt.Sprintf(`
		<div style="border: 1px solid #ccc; padding: 10px; margin: 5px;">
			<b>%s</b> wrote:
			<p>%s</p>
			<small>%s</small>
		</div>
		`, msg.Author, msg.Content, msg.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	html += `
	</div>
	<p>Hint: Try posting: <code>&lt;script&gt;alert('XSS')&lt;/script&gt;</code></p>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// PostMessagePage handles form submission for messages
func (h *MessageHandler) PostMessagePage(c *gin.Context) {
	content := c.PostForm("content")
	author := c.PostForm("author")

	if author == "" {
		author = "Anonymous"
	}

	if content != "" {
		submission := &entity.ContentSubmission{
			Content:     content,
			Author:      author,
			ContentType: "html",
			Context:     "message",
		}
		h.contentService.SubmitMessage(submission)
	}

	c.Redirect(http.StatusFound, "/messages")
}
