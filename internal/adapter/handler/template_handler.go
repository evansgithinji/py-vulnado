package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
)

// TemplateHandler handles template-related HTTP requests
type TemplateHandler struct {
	notificationService *service.NotificationService
}

// NewTemplateHandler creates a new TemplateHandler
func NewTemplateHandler(notificationService *service.NotificationService) *TemplateHandler {
	return &TemplateHandler{notificationService: notificationService}
}

// RegisterRoutes registers template routes
func (h *TemplateHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/api/greeting", h.Greeting)
	r.POST("/api/template/render", h.RenderTemplate)
	r.POST("/api/template/functions", h.RenderWithFunctions)
	r.GET("/api/invoice/preview", h.InvoicePreview)
	r.GET("/api/message/render", h.RenderMessage)

	// HTML pages
	r.GET("/greeting", h.GreetingPage)
	r.GET("/template", h.TemplatePage)
	r.GET("/template/functions", h.TemplateFunctionsPage)
	r.GET("/invoice", h.InvoicePage)
}

// Greeting renders a greeting message
// VULNERABLE: SSTI via message parameter (via NotificationService deep call graph)
func (h *TemplateHandler) Greeting(c *gin.Context) {
	message := c.DefaultQuery("message", "Welcome!")

	// VULNERABLE: SSTI via deep call graph
	request := &entity.RenderRequest{
		UserInput:    message,
		TemplateName: "greeting",
	}
	html, err := h.notificationService.RenderNotification(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// RenderTemplate renders a user-provided template
// VULNERABLE: SSTI - full template control (via NotificationService deep call graph)
func (h *TemplateHandler) RenderTemplate(c *gin.Context) {
	templateStr := c.PostForm("template")
	if templateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Template required"})
		return
	}

	// Build data from form
	variables := map[string]interface{}{
		"Name":      c.PostForm("name"),
		"Email":     c.PostForm("email"),
		"Secret":    "ADMIN_API_KEY_12345",
		"APIKey":    "sk-live-abcdef123456",
		"AdminPass": "super_secret_admin_pass",
	}

	// VULNERABLE: SSTI via deep call graph
	request := &entity.RenderRequest{
		UserInput:    templateStr,
		TemplateName: "custom",
		Variables:    variables,
	}
	result, err := h.notificationService.RenderNotification(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, result)
}

// InvoicePreview renders an invoice preview
// VULNERABLE: SSTI via template parameter (via NotificationService deep call graph)
func (h *TemplateHandler) InvoicePreview(c *gin.Context) {
	templateStr := c.Query("template")
	orderID := c.Query("order_id")
	totalStr := c.Query("total")

	if orderID == "" {
		orderID = "12345"
	}

	total := 99.99
	if totalStr != "" {
		if parsed, err := strconv.ParseFloat(totalStr, 64); err == nil {
			total = parsed
		}
	}

	variables := map[string]interface{}{
		"OrderID":   orderID,
		"Total":     total,
		"Secret":    "ADMIN_API_KEY_12345",
		"APIKey":    "sk-live-abcdef123456",
		"AdminPass": "super_secret_admin_pass",
	}

	// VULNERABLE: SSTI via deep call graph
	result, err := h.notificationService.RenderInvoice(templateStr, variables)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, result)
}

// RenderMessage renders a message template
// VULNERABLE: SSTI via message parameter
func (h *TemplateHandler) RenderMessage(c *gin.Context) {
	msg := c.Query("msg")
	if msg == "" {
		msg = "Hello, World!"
	}

	// VULNERABLE: SSTI via deep call graph
	request := &entity.RenderRequest{
		UserInput:    msg,
		TemplateName: "greeting",
	}
	result, err := h.notificationService.RenderNotification(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, result)
}

// RenderWithFunctions renders a template with exposed functions
// VULNERABLE: SSTI with dangerous function exposure
func (h *TemplateHandler) RenderWithFunctions(c *gin.Context) {
	templateStr := c.PostForm("template")
	if templateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Template required"})
		return
	}

	variables := map[string]interface{}{
		"Name":      "User",
		"Secret":    "ADMIN_API_KEY_12345",
		"APIKey":    "sk-live-abcdef123456",
		"AdminPass": "super_secret_admin_pass",
	}

	// VULNERABLE: SSTI via deep call graph
	request := &entity.RenderRequest{
		UserInput:    templateStr,
		TemplateName: "custom",
		Variables:    variables,
	}
	result, err := h.notificationService.RenderNotification(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, result)
}

// GreetingPage serves the greeting page
func (h *TemplateHandler) GreetingPage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Greeting - Vulnerable App</title></head>
<body>
	<h1>Greeting (SSTI Demo)</h1>

	<form method="GET" action="/api/greeting">
		<label>Message: <input type="text" name="message" size="50" placeholder="Welcome!"></label>
		<button type="submit">Render Greeting</button>
	</form>

	<p>Hints:</p>
	<ul>
		<li>Try: {{.}}</li>
		<li>Try: {{printf "%q" .}}</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// TemplatePage serves the template rendering page
func (h *TemplateHandler) TemplatePage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Template Render - Vulnerable App</title></head>
<body>
	<h1>Dynamic Template Rendering (SSTI Demo)</h1>

	<form method="POST" action="/api/template/render">
		<label>Name: <input type="text" name="name" placeholder="John"></label><br><br>
		<label>Email: <input type="text" name="email" placeholder="john@example.com"></label><br><br>
		<label>Template:</label><br>
		<textarea name="template" rows="5" cols="60">Hello {{.Name}}, your email is {{.Email}}</textarea><br><br>
		<button type="submit">Render Template</button>
	</form>

	<p>Hints:</p>
	<ul>
		<li>Access sensitive data: {{.Secret}}</li>
		<li>Access API key: {{.APIKey}}</li>
		<li>Access admin password: {{.AdminPass}}</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// InvoicePage serves the invoice preview page
func (h *TemplateHandler) InvoicePage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Invoice Preview - Vulnerable App</title></head>
<body>
	<h1>Invoice Preview (SSTI Demo)</h1>

	<form method="GET" action="/api/invoice/preview">
		<label>Order ID: <input type="text" name="order_id" placeholder="12345"></label><br><br>
		<label>Total: <input type="text" name="total" placeholder="99.99"></label><br><br>
		<label>Custom Template (optional):</label><br>
		<textarea name="template" rows="3" cols="60">Order #{{.OrderID}} - Total: ${{.Total}}</textarea><br><br>
		<button type="submit">Preview Invoice</button>
	</form>

	<p>Hints:</p>
	<ul>
		<li>Try custom template: {{.Secret}} - Expose sensitive data</li>
		<li>Default template exposes OrderID and Total</li>
		<li>Available fields: OrderID, Total, Secret, APIKey, AdminPass</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// TemplateFunctionsPage serves the template with functions page
func (h *TemplateHandler) TemplateFunctionsPage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Template Functions - Vulnerable App</title></head>
<body>
	<h1>Template with Functions (SSTI Demo)</h1>

	<form method="POST" action="/api/template/functions">
		<label>Template:</label><br>
		<textarea name="template" rows="5" cols="60">Hello {{.Name}}</textarea><br><br>
		<button type="submit">Render Template</button>
	</form>

	<p>Hints:</p>
	<ul>
		<li>Try: {{.Name}} - Access name variable</li>
		<li>Try: {{.Secret}} - Access sensitive data</li>
		<li>Try: {{.APIKey}} - Access API key</li>
		<li>Try: {{.AdminPass}} - Access admin password</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
