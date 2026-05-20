package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goapp/internal/domain/service"
)

type LogHandler struct {
	auditService *service.AuditService
}

func NewLogHandler(auditService *service.AuditService) *LogHandler {
	return &LogHandler{auditService: auditService}
}

func (h *LogHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/log/login", h.LogLogin)
	r.POST("/api/log/search", h.LogSearch)
	r.GET("/api/logs", h.GetLogs)
}

func (h *LogHandler) LogLogin(c *gin.Context) {
	var body struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username required"})
		return
	}

	status := "failed"
	if body.Username == "admin" && body.Password == "admin123" {
		status = "success"
	}

	// VULNERABLE: Log Injection (CWE-117) - unsanitized user input in log
	msg := h.auditService.LogLogin(body.Username, status)
	c.JSON(http.StatusOK, gin.H{"success": true, "logged": msg})
}

func (h *LogHandler) LogSearch(c *gin.Context) {
	var body struct {
		Query string `json:"query" form:"query"`
	}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query required"})
		return
	}

	// VULNERABLE: Log Injection (CWE-117) - unsanitized user input in log
	msg := h.auditService.LogSearch(body.Query)
	c.JSON(http.StatusOK, gin.H{"success": true, "logged": msg})
}

func (h *LogHandler) GetLogs(c *gin.Context) {
	logs := h.auditService.GetLogs()
	c.JSON(http.StatusOK, gin.H{"success": true, "count": len(logs), "logs": logs})
}
