package service

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// CatalogService orchestrates product search operations with deep call graph
type CatalogService struct {
	policyRepo *repository.SqlQueryPolicyRepository
	validator  *usecase.SearchQueryValidator
	builder    *usecase.SqlQueryBuilder
	executor   *usecase.SqlQueryExecutor
	mapper     *usecase.SearchResultMapper
}

// NewCatalogService creates a new CatalogService
func NewCatalogService(
	policyRepo *repository.SqlQueryPolicyRepository,
	validator *usecase.SearchQueryValidator,
	builder *usecase.SqlQueryBuilder,
	executor *usecase.SqlQueryExecutor,
	mapper *usecase.SearchResultMapper,
) *CatalogService {
	return &CatalogService{
		policyRepo: policyRepo,
		validator:  validator,
		builder:    builder,
		executor:   executor,
		mapper:     mapper,
	}
}

// Search executes a product search through the deep call graph
// Deep call graph for SQL Injection (depth 6):
// Handler -> CatalogService.Search()
//   -> SqlQueryPolicyRepository.GetPolicy() -> SqlQueryPolicy
//   -> SearchQueryValidator.Validate() -> ValidatedQuery
//   -> SqlQueryBuilder.BuildSearchQuery() -> raw SQL (VULNERABLE)
//   -> SqlQueryExecutor.Execute() -> rows (VULNERABLE)
//   -> SearchResultMapper.MapToProducts() -> Product[]
func (s *CatalogService) Search(request *entity.SearchRequest) ([]entity.Product, error) {
	policy := s.policyRepo.GetPolicy("product_search")
	validated := s.validator.Validate(request, policy)
	sql := s.builder.BuildSearchQuery(validated)
	rows, err := s.executor.Execute(sql)
	if err != nil {
		return nil, err
	}
	return s.mapper.MapToProducts(rows), nil
}

// SearchByCategory executes a category search through the deep call graph
func (s *CatalogService) SearchByCategory(category string) ([]entity.Product, error) {
	_ = s.policyRepo.GetPolicy("category_search")
	sql := s.builder.BuildCategoryQuery(category)
	rows, err := s.executor.Execute(sql)
	if err != nil {
		return nil, err
	}
	return s.mapper.MapToProducts(rows), nil
}
