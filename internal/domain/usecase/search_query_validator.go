package usecase

import "goapp/internal/domain/entity"

// SearchQueryValidator validates search queries against policy
type SearchQueryValidator struct{}

// NewSearchQueryValidator creates a new SearchQueryValidator
func NewSearchQueryValidator() *SearchQueryValidator {
	return &SearchQueryValidator{}
}

// Validate validates a search request against a policy
// VULNERABLE: "validates" but doesn't actually sanitize anything
func (v *SearchQueryValidator) Validate(request *entity.SearchRequest, policy *entity.SqlQueryPolicy) *entity.ValidatedQuery {
	query := request.Query
	// Fake validation - just checks length
	if len(query) > 1000 {
		query = query[:1000]
	}
	return &entity.ValidatedQuery{
		OriginalQuery: query,
		Table:         "products",
		Policy:        policy,
	}
}
