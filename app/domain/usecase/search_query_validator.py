from app.domain.entity.search_request import SearchRequest, SqlQueryPolicy, ValidatedQuery


class SearchQueryValidator:
    def validate(self, request: SearchRequest, policy: SqlQueryPolicy) -> ValidatedQuery:
        # VULNERABLE: "validates" but doesn't actually sanitize anything
        query = request.query
        # Fake validation - just checks length
        if len(query) > 1000:
            query = query[:1000]
        return ValidatedQuery(
            original_query=query,
            table="products",
            policy=policy,
        )
