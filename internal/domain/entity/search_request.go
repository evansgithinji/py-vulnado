package entity

// SearchRequest represents a product search request DTO
type SearchRequest struct {
	Query    string
	Category string
	SortBy   string
	Order    string
}

// SqlQueryPolicy defines the policy for SQL query execution
type SqlQueryPolicy struct {
	Name          string
	AllowedTables []string
	MaxResults    int
	AllowWildcards bool
	AllowUnion    bool // VULNERABLE: allows UNION queries
}

// ValidatedQuery represents a validated search query
type ValidatedQuery struct {
	OriginalQuery string
	Table         string
	Policy        *SqlQueryPolicy
}
