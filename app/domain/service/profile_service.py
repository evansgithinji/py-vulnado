from app.domain.entity.document_collection import NoSqlAuthRequest, WhereQueryRequest
from app.domain.repository.document_collection_repository import DocumentCollectionRepository
from app.domain.usecase.nosql_executor import (
    QueryBuilder,
    DocumentQueryExecutor,
    ExpressionEvaluator,
)


class ProfileService:
    """Orchestrates NoSQL authentication and query flow.
    Call graph: Handler → Service → CollectionRepository → QueryBuilder → Executor."""

    def __init__(
        self,
        collection_repo: DocumentCollectionRepository,
        query_builder: QueryBuilder,
        executor: DocumentQueryExecutor,
        expression_evaluator: ExpressionEvaluator,
    ):
        self._collection_repo = collection_repo
        self._query_builder = query_builder
        self._executor = executor
        self._expression_evaluator = expression_evaluator

    def authenticate_user(self, request: NoSqlAuthRequest) -> list:
        """Authenticate user via NoSQL query.
        VULNERABLE: NoSQL Injection flows through DTO → QueryBuilder → Executor."""
        collection = self._collection_repo.get_collection(request.collection)

        queries = self._query_builder.build_match_query(
            request.username, request.password
        )

        return self._executor.execute(collection, queries)

    def query(self, query_map: dict, collection_name: str = "users") -> list:
        """Execute NoSQL query.
        VULNERABLE: NoSQL Injection via operator objects."""
        collection = self._collection_repo.get_collection(collection_name)
        return self._executor.execute_generic(collection, query_map)

    def find_by_expression(self, request: WhereQueryRequest) -> list:
        """Find documents by $where expression.
        VULNERABLE: Code injection via eval()."""
        collection = self._collection_repo.get_collection(request.collection)
        return self._expression_evaluator.evaluate(
            request.where_expression, collection.documents
        )
