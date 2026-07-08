package repository

import "goapp/internal/domain/entity"

// SqlQueryPolicyRepository stores SQL query execution policies
type SqlQueryPolicyRepository struct {
	policies map[string]*entity.SqlQueryPolicy
}

// NewSqlQueryPolicyRepository creates a new SqlQueryPolicyRepository
func NewSqlQueryPolicyRepository() *SqlQueryPolicyRepository {
	return &SqlQueryPolicyRepository{
		policies: map[string]*entity.SqlQueryPolicy{
			"product_search": {
				Name:           "product_search",
				AllowedTables:  []string{"products"},
				MaxResults:     100,
				AllowWildcards: true,
				AllowUnion:     true, // VULNERABLE: allows UNION queries
			},
			"category_search": {
				Name:           "category_search",
				AllowedTables:  []string{"products"},
				MaxResults:     50,
				AllowWildcards: true,
				AllowUnion:     true, // VULNERABLE
			},
		},
	}
}

// GetPolicy retrieves a SQL query policy by name
func (r *SqlQueryPolicyRepository) GetPolicy(name string) *entity.SqlQueryPolicy {
	if policy, ok := r.policies[name]; ok {
		return policy
	}
	return r.policies["product_search"]
}
