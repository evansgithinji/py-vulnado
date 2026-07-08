package usecase

import (
	"fmt"

	"goapp/internal/domain/entity"
)

// SqlQueryBuilder builds SQL queries from validated queries
type SqlQueryBuilder struct{}

// NewSqlQueryBuilder creates a new SqlQueryBuilder
func NewSqlQueryBuilder() *SqlQueryBuilder {
	return &SqlQueryBuilder{}
}

// BuildSearchQuery builds a search SQL query
// VULNERABLE: String concatenation for SQL query building
func (b *SqlQueryBuilder) BuildSearchQuery(validated *entity.ValidatedQuery) string {
	query := validated.OriginalQuery
	table := validated.Table
	return fmt.Sprintf("SELECT * FROM %s WHERE name LIKE '%%%s%%'", table, query)
}

// BuildCategoryQuery builds a category search SQL query
// VULNERABLE: String concatenation for SQL query building
func (b *SqlQueryBuilder) BuildCategoryQuery(category string) string {
	return fmt.Sprintf("SELECT * FROM products WHERE category = '%s'", category)
}
