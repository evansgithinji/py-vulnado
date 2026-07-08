package entity

// DocumentCollection represents a NoSQL document collection.
type DocumentCollection struct {
	Name      string
	Documents []NoSqlUser
}

// NoSqlUser represents a user document in the NoSQL store.
type NoSqlUser struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

// NoSqlQuery represents a parsed NoSQL query.
type NoSqlQuery struct {
	Collection string
	Filter     map[string]interface{}
}

// NoSqlAuthRequest is a DTO for NoSQL authentication requests.
type NoSqlAuthRequest struct {
	Username interface{}
	Password interface{}
}

// WhereQueryRequest is a DTO for $where-style query requests.
type WhereQueryRequest struct {
	WhereClause string
	Collection  string
}
