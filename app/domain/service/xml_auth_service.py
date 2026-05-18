from app.domain.entity.xml_document import XPathAuthRequest
from app.domain.repository.xml_document_repository import XmlDocumentRepository
from app.domain.usecase.xpath_evaluator import XPathExpressionBuilder, XPathEvaluator


class XmlAuthService:
    """Orchestrates XPath-based authentication flow.
    Call graph: Handler → Service → Repository → ExpressionBuilder → Evaluator."""

    def __init__(
        self,
        repository: XmlDocumentRepository,
        expression_builder: XPathExpressionBuilder,
        evaluator: XPathEvaluator,
    ):
        self._repository = repository
        self._expression_builder = expression_builder
        self._evaluator = evaluator

    def authenticate_via_xpath(self, request: XPathAuthRequest) -> list:
        """Authenticate user via XPath query.
        VULNERABLE: XPath Injection flows through DTO → ExpressionBuilder → Evaluator."""
        document = self._repository.load_document(request.document_source)

        expression = self._expression_builder.build_auth_query(
            request.username, request.password
        )

        results = self._evaluator.evaluate(document, expression)
        return results

    def query(self, query: str, source: str = "default") -> list:
        """Execute raw XPath query.
        VULNERABLE: XPath Injection - user-controlled expression."""
        document = self._repository.load_document(source)
        return self._evaluator.evaluate_raw(document, query)
