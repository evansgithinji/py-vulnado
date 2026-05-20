package service

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// ProfileService orchestrates NoSQL authentication and query operations.
type ProfileService struct {
	repo          repository.DocumentCollectionRepository
	queryBuilder  *usecase.QueryBuilder
	queryExecutor *usecase.DocumentQueryExecutor
	exprEvaluator *usecase.ExpressionEvaluator
}

// NewProfileService creates a new ProfileService.
func NewProfileService(
	repo repository.DocumentCollectionRepository,
	queryBuilder *usecase.QueryBuilder,
	queryExecutor *usecase.DocumentQueryExecutor,
	exprEvaluator *usecase.ExpressionEvaluator,
) *ProfileService {
	return &ProfileService{
		repo:          repo,
		queryBuilder:  queryBuilder,
		queryExecutor: queryExecutor,
		exprEvaluator: exprEvaluator,
	}
}

// Authenticate performs NoSQL-style authentication.
// Call graph: ProfileService.Authenticate → repo.FindAllUsers → queryBuilder.BuildAuthQuery → queryExecutor.ExecuteAuth
func (s *ProfileService) Authenticate(req *entity.NoSqlAuthRequest) []entity.NoSqlUser {
	// Layer 3: Repository lookup for users
	users := s.repo.FindAllUsers()
	if users == nil {
		return nil
	}

	// Layer 4: Build auth query (passes through operator objects)
	_ = s.queryBuilder.BuildAuthQuery(req.Username, req.Password)

	// Layer 5: Execute query with operator support (VULNERABLE)
	return s.queryExecutor.ExecuteAuth(users, req.Username, req.Password)
}

// Query performs a NoSQL-style query with MongoDB-like operators.
// Call graph: ProfileService.Query → repo.FindAllUsers → queryExecutor.ExecuteQuery
func (s *ProfileService) Query(query map[string]interface{}) []entity.NoSqlUser {
	// Layer 3: Repository lookup for users
	users := s.repo.FindAllUsers()
	if users == nil {
		return nil
	}

	// Layer 5: Execute query with operator support (VULNERABLE)
	return s.queryExecutor.ExecuteQuery(users, query)
}

// FindUsersWhere evaluates a $where-style query.
// Call graph: ProfileService.FindUsersWhere → repo.FindAllUsers → queryBuilder.BuildWhereExpression → exprEvaluator.Evaluate
func (s *ProfileService) FindUsersWhere(req *entity.WhereQueryRequest) []entity.NoSqlUser {
	// Layer 3: Repository lookup for users
	users := s.repo.FindAllUsers()
	if users == nil {
		return nil
	}

	// Layer 4: Build template expression (VULNERABLE)
	tmplStr := s.queryBuilder.BuildWhereExpression(req.WhereClause)

	// Layer 5: Evaluate expression against each user
	var results []entity.NoSqlUser
	for _, user := range users {
		if s.exprEvaluator.Evaluate(tmplStr, user) {
			results = append(results, user)
		}
	}
	return results
}
