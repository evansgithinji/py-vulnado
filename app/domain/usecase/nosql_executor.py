import re
from app.domain.entity.document_collection import DocumentCollection, NoSqlQuery


class QueryBuilder:
    """Builds NoSQL query objects from user input.
    VULNERABLE: Operator injection preserved in Query objects."""

    def build_match_query(self, username, password) -> list:
        """VULNERABLE: NoSQL Injection (CWE-943)
        Operator injection - if username/password are dicts with $ne/$regex/$gt,
        the operators are preserved in the query objects."""
        queries = []
        queries.append({"field": "username", "condition": username})
        queries.append({"field": "password", "condition": password})
        return queries


class DocumentQueryExecutor:
    """Executes queries against document collections.
    VULNERABLE: Processes MongoDB-like operators."""

    def execute(self, collection: DocumentCollection, queries: list) -> list:
        """VULNERABLE: NoSQL Injection (CWE-943)
        $ne/$regex/$gt operators evaluated against documents."""
        results = []
        for doc in collection.documents:
            if self._matches_all(doc, queries):
                results.append(doc)
        return results

    def execute_generic(self, collection: DocumentCollection, query_map: dict) -> list:
        """VULNERABLE: NoSQL Injection (CWE-943)"""
        results = []
        for doc in collection.documents:
            match = True
            for field, condition in query_map.items():
                field_val = doc.get(field, "")
                if not self._match_operator(field_val, condition):
                    match = False
                    break
            if match:
                results.append(doc)
        return results

    def _matches_all(self, doc: dict, queries: list) -> bool:
        for q in queries:
            field_val = doc.get(q["field"], "")
            if not self._match_operator(field_val, q["condition"]):
                return False
        return True

    def _match_operator(self, field_val: str, condition) -> bool:
        if isinstance(condition, str):
            return field_val == condition
        if isinstance(condition, dict):
            for op, op_val in condition.items():
                if op == "$ne":
                    if field_val == str(op_val):
                        return False
                elif op == "$gt":
                    if not (field_val > str(op_val)):
                        return False
                elif op == "$regex":
                    if not re.search(str(op_val), field_val):
                        return False
            return True
        return False


class ExpressionEvaluator:
    """Evaluates $where expressions against documents.
    VULNERABLE: eval() with user input."""

    def evaluate(self, where_expression: str, documents: list) -> list:
        """VULNERABLE: Code injection via eval() $where clause (CWE-94)"""
        results = []
        for doc in documents:
            try:
                if eval(where_expression, {"__builtins__": {}}, doc):
                    results.append(doc)
            except Exception:
                pass
        return results
