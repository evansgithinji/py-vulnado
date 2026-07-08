package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
	"goapp/internal/domain/usecase"
)

// NetworkHandler handles network-related HTTP requests
type NetworkHandler struct {
	commandService        *service.SystemCommandService
	externalRequestService *service.ExternalRequestService
	htmlBuilder           *usecase.HtmlResponseBuilder
}

// NewNetworkHandler creates a new NetworkHandler
func NewNetworkHandler(
	commandService *service.SystemCommandService,
	externalRequestService *service.ExternalRequestService,
	htmlBuilder *usecase.HtmlResponseBuilder,
) *NetworkHandler {
	return &NetworkHandler{
		commandService:        commandService,
		externalRequestService: externalRequestService,
		htmlBuilder:           htmlBuilder,
	}
}

// RegisterRoutes registers network routes
func (h *NetworkHandler) RegisterRoutes(r *gin.Engine) {
	// Command Injection
	r.POST("/api/ping", h.PingHost)
	r.GET("/api/ping", h.PingHostGet)

	// SSRF
	r.GET("/api/proxy", h.ProxyRequest)
	r.POST("/api/webhook/test", h.TestWebhook)
	r.Any("/api/fetch/cached", h.FetchAndCacheURL)
	r.GET("/api/fetch/history", h.GetCachedURLs)

	// HTML pages
	r.GET("/ping", h.PingPage)
	r.GET("/webhook", h.WebhookPage)
	r.GET("/fetch/cached", h.CachedFetchPage)
}

// PingHost pings a host
// VULNERABLE: Command Injection via SystemCommandService deep call graph
func (h *NetworkHandler) PingHost(c *gin.Context) {
	host := c.PostForm("host")
	if host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Host required"})
		return
	}

	// VULNERABLE: Command injection via deep call graph
	request := &entity.CommandRequest{Target: host, CommandType: "ping"}
	output, err := h.commandService.ExecuteCommand(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Ping failed",
			"output": output,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"host":   host,
		"output": output,
	})
}

// PingHostGet pings a host (GET variant)
// VULNERABLE: Command Injection
func (h *NetworkHandler) PingHostGet(c *gin.Context) {
	host := c.Query("host")
	if host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Host required"})
		return
	}

	// VULNERABLE: Command injection via deep call graph
	request := &entity.CommandRequest{Target: host, CommandType: "ping"}
	output, err := h.commandService.ExecuteCommand(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Ping failed",
			"output": output,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"host":   host,
		"output": output,
	})
}

// ProxyRequest proxies a request to another URL
// VULNERABLE: SSRF via ExternalRequestService deep call graph
func (h *NetworkHandler) ProxyRequest(c *gin.Context) {
	targetURL := c.Query("url")
	if targetURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL required"})
		return
	}

	// VULNERABLE: SSRF via deep call graph
	resp, err := h.externalRequestService.ProxyRequest(targetURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Proxy failed", "detail": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Forward response
	c.DataFromReader(resp.StatusCode, resp.ContentLength, resp.Header.Get("Content-Type"), resp.Body, nil)
}

// TestWebhook tests a webhook URL
// VULNERABLE: SSRF via POST
func (h *NetworkHandler) TestWebhook(c *gin.Context) {
	webhookURL := c.PostForm("url")
	payload := c.PostForm("payload")

	if webhookURL == "" {
		webhookURL = c.Query("url")
	}
	if payload == "" {
		payload = c.Query("payload")
	}

	if webhookURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Webhook URL required"})
		return
	}

	if payload == "" {
		payload = `{"test": true}`
	}

	// VULNERABLE: SSRF via deep call graph
	response, err := h.externalRequestService.TestWebhook(webhookURL, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Webhook test failed", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Webhook tested",
		"url":      webhookURL,
		"response": response,
	})
}

// PingPage serves the ping page
// VULNERABLE: Reflected XSS via HtmlResponseBuilder
func (h *NetworkHandler) PingPage(c *gin.Context) {
	host := c.Query("host")
	var output string

	if host != "" {
		request := &entity.CommandRequest{Target: host, CommandType: "ping"}
		var err error
		output, err = h.commandService.ExecuteCommand(request)
		if err != nil {
			output = "Error: " + err.Error() + "\n" + output
		}

		// VULNERABLE: Reflected XSS via HtmlResponseBuilder
		htmlContent := h.htmlBuilder.BuildPingPage(host, output)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlContent)
		return
	}

	// VULNERABLE: Reflected XSS via host parameter
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><title>Ping - Vulnerable App</title></head>
<body>
	<h1>Network Ping (Command Injection Demo)</h1>

	<form method="GET" action="/ping">
		<label>Host: <input type="text" name="host" value="%s"></label>
		<button type="submit">Ping (GET)</button>
	</form>

	<form method="POST" action="/api/ping">
		<label>Host: <input type="text" name="host" placeholder="127.0.0.1"></label>
		<button type="submit">Ping (POST)</button>
	</form>
`, host) // VULNERABLE: Reflected XSS

	html += `
	<p>Hints:</p>
	<ul>
		<li>Try: 127.0.0.1; id</li>
		<li>Try: 127.0.0.1 && cat /etc/passwd</li>
		<li>Try: 127.0.0.1 | whoami</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// WebhookPage serves the webhook testing page
func (h *NetworkHandler) WebhookPage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Webhook Test - Vulnerable App</title></head>
<body>
	<h1>Webhook Test (SSRF Demo)</h1>

	<form method="POST" action="/api/webhook/test">
		<label>Webhook URL: <input type="text" name="url" size="50" placeholder="http://example.com/webhook"></label><br><br>
		<label>Payload: <textarea name="payload" rows="3" cols="50">{"test": true}</textarea></label><br><br>
		<button type="submit">Test Webhook</button>
	</form>

	<p>Hints:</p>
	<ul>
		<li>Try: http://localhost:8080/api/debug</li>
		<li>Try: http://internal-service:8080/admin</li>
		<li>Try cloud metadata endpoints</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// FetchAndCacheURL fetches a URL and caches it
// VULNERABLE: SSRF + cache information leakage
func (h *NetworkHandler) FetchAndCacheURL(c *gin.Context) {
	url := c.PostForm("url")
	if url == "" {
		url = c.Query("url")
	}

	if url == "" {
		c.Redirect(http.StatusFound, "/fetch/cached")
		return
	}

	// VULNERABLE: SSRF via deep call graph
	content, err := h.externalRequestService.FetchRemoteFile(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.String(http.StatusOK, "Content from %s:\n\n%s", url, content)
}

// GetCachedURLs returns all cached URLs
// VULNERABLE: Information disclosure
func (h *NetworkHandler) GetCachedURLs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Cache endpoint"})
}

// CachedFetchPage serves the cached URL fetch page
func (h *NetworkHandler) CachedFetchPage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>URL Fetcher with Cache - Vulnerable App</title></head>
<body>
	<h1>URL Content Fetcher (SSRF with Cache)</h1>
	<form method="POST" action="/api/fetch/cached">
		<input type="text" name="url" placeholder="Enter URL to fetch" style="width:300px">
		<input type="submit" value="Fetch">
	</form>
	<p>Hints:</p>
	<ul>
		<li>Try: http://localhost:8080/api/debug</li>
		<li>Try: http://internal-api/admin/credentials</li>
		<li>Try: http://metadata-service/</li>
		<li>Check /api/fetch/history to see all cached URLs</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
