package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
	"goapp/internal/domain/usecase"
)

type NoSqlHandler struct {
	profileService *service.ProfileService
}

func NewNoSqlHandler(profileService *service.ProfileService) *NoSqlHandler {
	return &NoSqlHandler{profileService: profileService}
}

func (h *NoSqlHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/nosql/login", h.Login)
	r.POST("/api/nosql/query", h.Query)
	r.GET("/api/nosql/users", h.FindWhere)
}

func safeNoSqlUser(u entity.NoSqlUser) gin.H {
	return gin.H{"_id": u.ID, "username": u.Username, "role": u.Role, "email": u.Email}
}

func (h *NoSqlHandler) Login(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	username, password, err := usecase.ParseNoSqlBody(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	// Create DTO and call service
	req := &entity.NoSqlAuthRequest{
		Username: username,
		Password: password,
	}

	// VULNERABLE: NoSQL Injection (CWE-943) - operator injection
	results := h.profileService.Authenticate(req)

	if len(results) > 0 {
		safe := make([]gin.H, len(results))
		for i, u := range results {
			safe[i] = safeNoSqlUser(u)
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login successful", "users": safe})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials"})
	}
}

func (h *NoSqlHandler) Query(c *gin.Context) {
	var body struct {
		Query map[string]interface{} `json:"query"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query object required"})
		return
	}

	// VULNERABLE: NoSQL Injection (CWE-943)
	results := h.profileService.Query(body.Query)

	safe := make([]gin.H, len(results))
	for i, u := range results {
		safe[i] = safeNoSqlUser(u)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "count": len(results), "results": safe})
}

func (h *NoSqlHandler) FindWhere(c *gin.Context) {
	where := c.Query("where")
	if where == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "where parameter required"})
		return
	}

	// Create DTO and call service
	req := &entity.WhereQueryRequest{
		WhereClause: where,
		Collection:  "users",
	}

	// VULNERABLE: Code injection via $where clause
	results := h.profileService.FindUsersWhere(req)

	safe := make([]gin.H, len(results))
	for i, u := range results {
		safe[i] = safeNoSqlUser(u)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "count": len(results), "results": safe})
}
