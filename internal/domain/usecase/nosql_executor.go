package usecase

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
	"text/template"

	"goapp/internal/domain/entity"
)

// QueryBuilder builds NoSQL queries from request parameters.
type QueryBuilder struct{}

// NewQueryBuilder creates a new QueryBuilder.
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// BuildAuthQuery constructs an authentication query from username and password.
func (b *QueryBuilder) BuildAuthQuery(username, password interface{}) map[string]interface{} {
	return map[string]interface{}{
		"username": username,
		"password": password,
	}
}

// BuildWhereExpression constructs a Go template expression from a $where clause.
// VULNERABLE: Code injection via template evaluation
func (b *QueryBuilder) BuildWhereExpression(whereClause string) string {
	return "{{if " + whereClause + "}}true{{end}}"
}

// DocumentQueryExecutor executes queries against document collections.
type DocumentQueryExecutor struct{}

// NewDocumentQueryExecutor creates a new DocumentQueryExecutor.
func NewDocumentQueryExecutor() *DocumentQueryExecutor {
	return &DocumentQueryExecutor{}
}

// ExecuteAuth performs authentication by matching username and password with operator support.
// VULNERABLE: NoSQL Injection (CWE-943) - operator injection via JSON
func (e *DocumentQueryExecutor) ExecuteAuth(users []entity.NoSqlUser, username, password interface{}) []entity.NoSqlUser {
	var results []entity.NoSqlUser
	for _, user := range users {
		if e.matchOperator(user.Username, username) && e.matchOperator(user.Password, password) {
			results = append(results, user)
		}
	}
	return results
}

// ExecuteQuery performs a query by matching fields with operator support.
// VULNERABLE: NoSQL Injection (CWE-943)
func (e *DocumentQueryExecutor) ExecuteQuery(users []entity.NoSqlUser, query map[string]interface{}) []entity.NoSqlUser {
	var results []entity.NoSqlUser
	for _, user := range users {
		match := true
		for field, condition := range query {
			fieldVal := e.getField(user, field)
			if !e.matchOperator(fieldVal, condition) {
				match = false
				break
			}
		}
		if match {
			results = append(results, user)
		}
	}
	return results
}

// matchOperator evaluates a MongoDB-like operator against a field value.
// VULNERABLE: NoSQL Injection (CWE-943) - operator injection
func (e *DocumentQueryExecutor) matchOperator(fieldVal string, condition interface{}) bool {
	switch v := condition.(type) {
	case string:
		return fieldVal == v
	case map[string]interface{}:
		for op, opVal := range v {
			strVal, _ := opVal.(string)
			switch op {
			case "$ne":
				if fieldVal == strVal {
					return false
				}
			case "$gt":
				if !(fieldVal > strVal) {
					return false
				}
			case "$regex":
				matched, _ := regexp.MatchString(strVal, fieldVal)
				if !matched {
					return false
				}
			}
		}
		return true
	}
	return false
}

func (e *DocumentQueryExecutor) getField(user entity.NoSqlUser, field string) string {
	switch strings.ToLower(field) {
	case "username":
		return user.Username
	case "password":
		return user.Password
	case "role":
		return user.Role
	case "email":
		return user.Email
	case "_id":
		return user.ID
	}
	return ""
}

// ExpressionEvaluator evaluates Go template expressions for $where-style queries.
type ExpressionEvaluator struct{}

// NewExpressionEvaluator creates a new ExpressionEvaluator.
func NewExpressionEvaluator() *ExpressionEvaluator {
	return &ExpressionEvaluator{}
}

// Evaluate evaluates a Go template expression against a user document.
// VULNERABLE: Code injection via template evaluation
func (e *ExpressionEvaluator) Evaluate(tmplStr string, user entity.NoSqlUser) bool {
	tmpl, err := template.New("where").Parse(tmplStr)
	if err != nil {
		return false
	}

	data := map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
		"role":     user.Role,
		"email":    user.Email,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return false
	}
	return buf.String() == "true"
}

// ParseNoSqlBody parses JSON body and returns username/password as interface{} to allow operator objects.
func ParseNoSqlBody(body []byte) (interface{}, interface{}, error) {
	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, nil, err
	}
	return parsed["username"], parsed["password"], nil
}
