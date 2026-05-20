package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
)

type LdapHandler struct {
	directoryService *service.DirectoryService
}

func NewLdapHandler(directoryService *service.DirectoryService) *LdapHandler {
	return &LdapHandler{directoryService: directoryService}
}

func (h *LdapHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/ldap/login", h.Login)
	r.GET("/api/ldap/search", h.Search)
}

func (h *LdapHandler) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password required"})
		return
	}

	// Create DTO and call service
	req := &entity.LdapAuthRequest{
		Username: body.Username,
		Password: body.Password,
		Domain:   "default",
	}

	// VULNERABLE: LDAP Injection (CWE-90) - unsanitized user input
	results := h.directoryService.Authenticate(req)

	if len(results) > 0 {
		safeResults := make([]gin.H, len(results))
		for i, u := range results {
			safeResults[i] = gin.H{"uid": u.UID, "cn": u.CN, "mail": u.Mail, "ou": u.OU}
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login successful", "users": safeResults})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
	}
}

func (h *LdapHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "q parameter required"})
		return
	}

	// VULNERABLE: LDAP Injection (CWE-90) - unsanitized user input
	results := h.directoryService.Search(q)

	safeResults := make([]gin.H, len(results))
	for i, u := range results {
		safeResults[i] = gin.H{"uid": u.UID, "cn": u.CN, "mail": u.Mail, "ou": u.OU}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "count": len(results), "results": safeResults})
}
