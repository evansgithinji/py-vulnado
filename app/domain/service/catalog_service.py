from app.domain.entity.search_request import SearchRequest
from app.domain.repository.sql_query_policy_repository import SqlQueryPolicyRepository
from app.domain.usecase.search_query_validator import SearchQueryValidator
from app.domain.usecase.sql_query_builder import SqlQueryBuilder
from app.domain.usecase.sql_query_executor import SqlQueryExecutor
from app.domain.usecase.search_result_mapper import SearchResultMapper


class CatalogService:
    def __init__(
        self,
        policy_repo: SqlQueryPolicyRepository,
        validator: SearchQueryValidator,
        query_builder: SqlQueryBuilder,
        executor: SqlQueryExecutor,
        mapper: SearchResultMapper,
    ):
        self._policy_repo = policy_repo
        self._validator = validator
        self._query_builder = query_builder
        self._executor = executor
        self._mapper = mapper

    def search(self, request: SearchRequest) -> list:
        """
        Deep call graph for SQL Injection:
        Handler → CatalogService.search()
          → SqlQueryPolicyRepository.getPolicy() → SqlQueryPolicy
          → SearchQueryValidator.validate() → ValidatedQuery
          → SqlQueryBuilder.buildSearchQuery() → raw SQL (VULNERABLE)
          → SqlQueryExecutor.execute() → rows (VULNERABLE)
          → SearchResultMapper.mapToProducts() → Product[]
        """
        policy = self._policy_repo.get_policy("product_search")
        validated = self._validator.validate(request, policy)
        sql = self._query_builder.build_search_query(validated)
        rows = self._executor.execute(sql)
        return self._mapper.map_to_products(rows)

    def search_by_category(self, category: str) -> list:
        policy = self._policy_repo.get_policy("category_search")
        sql = self._query_builder.build_category_query(category)
        rows = self._executor.execute(sql)
        return self._mapper.map_to_products(rows)
