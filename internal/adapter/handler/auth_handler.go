package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
	"goapp/internal/domain/usecase"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService       *service.AuthenticationService
	navigationService *service.NavigationService
	htmlBuilder       *usecase.HtmlResponseBuilder
	authUseCase       *usecase.AuthUseCase // kept for profile/search/getUser operations
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	authService *service.AuthenticationService,
	navigationService *service.NavigationService,
	htmlBuilder *usecase.HtmlResponseBuilder,
	authUseCase *usecase.AuthUseCase,
) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		navigationService: navigationService,
		htmlBuilder:       htmlBuilder,
		authUseCase:       authUseCase,
	}
}

// RegisterRoutes registers auth routes
func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/api/login", h.Login)
	r.GET("/api/users/search", h.SearchUsers)
	r.PUT("/api/users/:id/profile", h.UpdateProfile)
	r.GET("/api/users/:id", h.GetUser)

	// VULNERABLE: Login with redirect (Open Redirect)
	r.GET("/api/login/redirect", h.LoginWithRedirect)

	// HTML form pages
	r.GET("/login", h.LoginPage)
	r.GET("/users/search", h.SearchPage)
}

// Login handles user login
// VULNERABLE: SQL Injection via username and password (via AuthenticationService deep call graph)
func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password required"})
		return
	}

	// VULNERABLE: Credentials passed through deep call graph
	request := &entity.AuthRequest{Username: username, Password: password}
	result := h.authService.Authenticate(request)
	if result == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    result,
	})
}

// SearchUsers searches for users
// VULNERABLE: SQL Injection via search term
func (h *AuthHandler) SearchUsers(c *gin.Context) {
	term := c.Query("q")

	// VULNERABLE: Search term passed directly (SQLi)
	users, err := h.authUseCase.SearchUsers(term)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// UpdateProfile updates a user's profile
// VULNERABLE: Privilege escalation via is_admin parameter
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	email := c.PostForm("email")
	// VULNERABLE: is_admin can be set by any user
	isAdminStr := c.PostForm("is_admin")
	isAdmin := isAdminStr == "1" || isAdminStr == "true"

	err = h.authUseCase.UpdateProfile(userID, email, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated"})
}

// GetUser retrieves a user by ID
func (h *AuthHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.authUseCase.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// LoginPage serves the login form
func (h *AuthHandler) LoginPage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Login - Vulnerable App</title></head>
<body>
	<h1>Login</h1>
	<form method="POST" action="/api/login">
		<label>Username: <input type="text" name="username"></label><br><br>
		<label>Password: <input type="password" name="password"></label><br><br>
		<button type="submit">Login</button>
	</form>
	<p>Hint: Try admin' OR '1'='1 as username</p>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// LoginWithRedirect handles login with redirect parameter
// VULNERABLE: Open Redirect after login (via NavigationService deep call graph)
func (h *AuthHandler) LoginWithRedirect(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	redirect := c.Query("redirect")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password required"})
		return
	}

	// VULNERABLE: Credentials passed through deep call graph
	request := &entity.AuthRequest{Username: username, Password: password}
	result := h.authService.Authenticate(request)
	if result == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// VULNERABLE: Open redirect via NavigationService deep call graph
	if redirect != "" {
		redirectReq := &entity.RedirectRequest{TargetURL: redirect, Context: "auth"}
		finalURL := h.navigationService.ResolveRedirect(redirectReq)
		c.Redirect(http.StatusFound, finalURL)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    result,
	})
}

// SearchPage serves the user search page
// VULNERABLE: Reflected XSS via search term display (via HtmlResponseBuilder)
func (h *AuthHandler) SearchPage(c *gin.Context) {
	term := c.Query("q")

	users, _ := h.authUseCase.SearchUsers(term)

	if term != "" {
		// VULNERABLE: Reflected XSS via HtmlResponseBuilder
		htmlContent := h.htmlBuilder.BuildUserSearchPage(term, users)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, htmlContent)
		return
	}

	// VULNERABLE: Reflected XSS - term displayed without encoding
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><title>User Search - Vulnerable App</title></head>
<body>
	<h1>User Search</h1>
	<form method="GET" action="/users/search">
		<label>Search: <input type="text" name="q" value="%s"></label>
		<button type="submit">Search</button>
	</form>
	<h2>Results for: %s</h2>
	<ul>
`, term, term) // VULNERABLE: Reflected XSS

	for _, user := range users {
		html += fmt.Sprintf("<li>%s - %s</li>", user.Username, user.Email)
	}

	html += `
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
