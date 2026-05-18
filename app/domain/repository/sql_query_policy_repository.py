from app.domain.entity.search_request import SqlQueryPolicy


class SqlQueryPolicyRepository:
    def __init__(self):
        self._policies = {
            "product_search": SqlQueryPolicy(
                name="product_search",
                allowed_tables=["products"],
                max_results=100,
                allow_wildcards=True,
                allow_union=True,
            ),
            "category_search": SqlQueryPolicy(
                name="category_search",
                allowed_tables=["products"],
                max_results=50,
                allow_wildcards=True,
                allow_union=True,
            ),
        }

    def get_policy(self, name: str) -> SqlQueryPolicy:
        return self._policies.get(name, self._policies["product_search"])
