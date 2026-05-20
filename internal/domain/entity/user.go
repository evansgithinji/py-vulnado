package entity

// User represents a user in the system
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	IsAdmin  bool   `json:"is_admin"`
}

// UserCredentials for login
type UserCredentials struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
