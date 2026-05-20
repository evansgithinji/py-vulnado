from lxml import etree
from app.domain.entity.xml_document import XmlDocument


class XPathExpressionBuilder:
    """Builds XPath expressions from user input.
    VULNERABLE: String concatenation without sanitization."""

    def build_auth_query(self, username: str, password: str) -> str:
        """VULNERABLE: XPath Injection (CWE-643)
        String concatenation in XPath expression."""
        return "/users/user[username='" + username + "' and password='" + password + "']/name"

    def build_search_query(self, field: str, value: str) -> str:
        """VULNERABLE: XPath Injection (CWE-643)"""
        return "/users/user[" + field + "='" + value + "']"


class XPathEvaluator:
    """Evaluates XPath expressions against XML documents.
    VULNERABLE: Executes potentially injected XPath expressions."""

    def evaluate(self, document: XmlDocument, expression: str) -> list:
        """VULNERABLE: XPath Injection (CWE-643)
        Executes XPath expression that may contain injected operators."""
        if document.parsed_doc is None:
            document.parsed_doc = etree.fromstring(document.content)

        try:
            results = document.parsed_doc.xpath(expression)
            if isinstance(results, list):
                return [r.text if hasattr(r, "text") else str(r) for r in results]
            return [str(results)]
        except Exception as e:
            raise ValueError(f"XPath query failed: {e}")

    def evaluate_raw(self, document: XmlDocument, query: str) -> list:
        """VULNERABLE: XPath Injection (CWE-643)
        Direct user input as XPath expression."""
        if document.parsed_doc is None:
            document.parsed_doc = etree.fromstring(document.content)

        try:
            results = document.parsed_doc.xpath(query)
            if isinstance(results, list):
                return [r.text if hasattr(r, "text") else str(r) for r in results]
            return [str(results)]
        except Exception as e:
            raise ValueError(f"XPath query failed: {e}")
