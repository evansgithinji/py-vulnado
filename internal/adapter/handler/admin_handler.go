package handler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
)

// AdminHandler handles admin-related HTTP requests
type AdminHandler struct {
	commandService    *service.SystemCommandService
	navigationService *service.NavigationService
	backupsDir        string
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(
	commandService *service.SystemCommandService,
	navigationService *service.NavigationService,
	backupsDir string,
) *AdminHandler {
	return &AdminHandler{
		commandService:    commandService,
		navigationService: navigationService,
		backupsDir:        backupsDir,
	}
}

// RegisterRoutes registers admin routes
func (h *AdminHandler) RegisterRoutes(r *gin.Engine) {
	// Information Disclosure
	r.GET("/api/debug", h.DebugInfo)

	// Command Injection
	r.POST("/api/invoice/generate", h.GenerateInvoice)
	r.POST("/api/config/apply", h.ApplyConfig)
	r.POST("/api/backup", h.BackupDatabase)
	r.POST("/api/restore", h.RestoreDatabase)

	// Open Redirect
	r.GET("/redirect", h.Redirect)
	r.GET("/api/redirect", h.RedirectAPI)

	// HTML page
	r.GET("/admin", h.AdminPage)
}

// DebugInfo returns debug information
// VULNERABLE: Information Disclosure
func (h *AdminHandler) DebugInfo(c *gin.Context) {
	// VULNERABLE: Exposes sensitive information via diagnostic command
	request := &entity.CommandRequest{Target: "env", CommandType: "diagnostic"}
	info, _ := h.commandService.ExecuteCommand(request)
	c.String(http.StatusOK, info)
}

// GenerateInvoice generates an invoice PDF
// VULNERABLE: Command Injection via SystemCommandService deep call graph
func (h *AdminHandler) GenerateInvoice(c *gin.Context) {
	orderID := c.PostForm("order_id")
	format := c.PostForm("format")

	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID required"})
		return
	}

	if format == "" {
		format = "pdf"
	}

	// VULNERABLE: Command injection via deep call graph
	target := fmt.Sprintf("echo 'Invoice for Order #%s' | wkhtmltopdf - invoice_%s.%s", orderID, orderID, format)
	request := &entity.CommandRequest{Target: target, CommandType: "diagnostic"}
	output, err := h.commandService.ExecuteCommand(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"output": output,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Invoice generated",
		"order_id": orderID,
		"format":   format,
	})
}

// ApplyConfig applies a configuration
// VULNERABLE: Command Injection via config parameter
func (h *AdminHandler) ApplyConfig(c *gin.Context) {
	config := c.PostForm("config")
	if config == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Config required"})
		return
	}

	// VULNERABLE: Command injection via deep call graph
	request := &entity.CommandRequest{Target: config, CommandType: "diagnostic"}
	output, err := h.commandService.ExecuteCommand(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Config applied",
		"output":  output,
	})
}

// BackupDatabase creates a database backup
// VULNERABLE: Command Injection via SystemCommandService deep call graph
func (h *AdminHandler) BackupDatabase(c *gin.Context) {
	backupName := c.PostForm("name")
	if backupName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Backup name required"})
		return
	}

	// VULNERABLE: Command injection via deep call graph
	output, err := h.commandService.ExecuteBackup(backupName, h.backupsDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": output})
}

// RestoreDatabase restores a database backup
// VULNERABLE: Command Injection
func (h *AdminHandler) RestoreDatabase(c *gin.Context) {
	backupName := c.PostForm("name")
	if backupName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Backup name required"})
		return
	}

	// VULNERABLE: Command injection via deep call graph
	target := fmt.Sprintf("cp %s/%s.db /app/data/app.db", h.backupsDir, backupName)
	request := &entity.CommandRequest{Target: target, CommandType: "diagnostic"}
	output, err := h.commandService.ExecuteCommand(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database restored from: " + backupName, "output": output})
}

// Redirect redirects to a URL
// VULNERABLE: Open Redirect via NavigationService deep call graph
func (h *AdminHandler) Redirect(c *gin.Context) {
	redirectURL := c.Query("url")
	if redirectURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL required"})
		return
	}

	// VULNERABLE: URL format check doesn't prevent open redirect
	if _, err := url.ParseRequestURI(redirectURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	// VULNERABLE: Open redirect via deep call graph
	request := &entity.RedirectRequest{TargetURL: redirectURL, Context: "admin"}
	finalURL := h.navigationService.ResolveRedirect(request)
	c.Redirect(http.StatusFound, finalURL)
}

// RedirectAPI redirects to a URL (no validation)
// VULNERABLE: Open Redirect via NavigationService deep call graph
func (h *AdminHandler) RedirectAPI(c *gin.Context) {
	redirectURL := c.Query("url")
	if redirectURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL required"})
		return
	}

	// VULNERABLE: Direct redirect via deep call graph
	request := &entity.RedirectRequest{TargetURL: redirectURL, Context: "admin"}
	finalURL := h.navigationService.ResolveRedirect(request)
	c.Redirect(http.StatusFound, finalURL)
}

// AdminPage serves the admin page
func (h *AdminHandler) AdminPage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Admin Panel - Vulnerable App</title></head>
<body>
	<h1>Admin Panel</h1>

	<h2>Debug Information</h2>
	<a href="/api/debug" target="_blank">View Debug Info</a> (Information Disclosure)

	<h2>Generate Invoice (Command Injection)</h2>
	<form method="POST" action="/api/invoice/generate">
		<label>Order ID: <input type="text" name="order_id" placeholder="12345"></label><br><br>
		<label>Format: <input type="text" name="format" placeholder="pdf"></label><br><br>
		<button type="submit">Generate</button>
	</form>
	<p>Hint: Try order_id=12345; id</p>

	<h2>Apply Config (Command Injection)</h2>
	<form method="POST" action="/api/config/apply">
		<label>Config: <input type="text" name="config" size="50" placeholder="cmd:whoami"></label>
		<button type="submit">Apply</button>
	</form>
	<p>Hint: Try config=cmd:cat /etc/passwd</p>

	<h2>Database Backup (Command Injection)</h2>
	<form method="POST" action="/api/backup">
		<label>Backup Name: <input type="text" name="name" placeholder="backup_2024"></label>
		<button type="submit">Backup</button>
	</form>
	<p>Hint: Try name=test; id</p>

	<h2>Open Redirect</h2>
	<form method="GET" action="/redirect">
		<label>URL: <input type="text" name="url" size="50" placeholder="https://evil.com"></label>
		<button type="submit">Redirect (with URL parse)</button>
	</form>
	<form method="GET" action="/api/redirect">
		<label>URL: <input type="text" name="url" size="50" placeholder="https://evil.com"></label>
		<button type="submit">Redirect (no validation)</button>
	</form>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
