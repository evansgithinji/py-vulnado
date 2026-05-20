from app.domain.entity.search_request import ValidatedQuery


class SqlQueryBuilder:
    def build_search_query(self, validated: ValidatedQuery) -> str:
        # VULNERABLE: String concatenation for SQL query building
        query = validated.original_query
        table = validated.table
        return f"SELECT * FROM {table} WHERE name LIKE '%{query}%'"

    def build_category_query(self, category: str) -> str:
        # VULNERABLE: String concatenation for SQL query building
        return f"SELECT * FROM products WHERE category = '{category}'"
