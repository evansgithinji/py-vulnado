package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
)

type XPathHandler struct {
	xmlAuthService *service.XmlAuthService
}

func NewXPathHandler(xmlAuthService *service.XmlAuthService) *XPathHandler {
	return &XPathHandler{xmlAuthService: xmlAuthService}
}

func (h *XPathHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/xpath/login", h.Login)
	r.GET("/api/xpath/query", h.Query)
}

func (h *XPathHandler) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	// Create DTO and call service
	req := &entity.XPathAuthRequest{
		Username: body.Username,
		Password: body.Password,
	}

	// VULNERABLE: XPath Injection (CWE-643)
	result, err := h.xmlAuthService.Authenticate(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login successful", "name": result})
}

func (h *XPathHandler) Query(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q parameter required"})
		return
	}

	// VULNERABLE: XPath Injection (CWE-643)
	result, err := h.xmlAuthService.Query(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "result": result})
}
