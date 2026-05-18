from lxml import etree

USERS_XML = b"""<?xml version="1.0" encoding="UTF-8"?>
<users>
  <user><username>admin</username><password>admin123</password><name>Admin</name><role>admin</role></user>
  <user><username>john</username><password>password</password><name>John Doe</name><role>user</role></user>
  <user><username>jane</username><password>password123</password><name>Jane Smith</name><role>user</role></user>
</users>"""


class XPathUseCase:
    def xpath_login(self, username: str, password: str) -> list:
        """VULNERABLE: XPath Injection (CWE-643)
        String concatenation in XPath expression
        """
        doc = etree.fromstring(USERS_XML)

        # VULNERABLE: XPath Injection - string concatenation
        expression = "/users/user[username='" + username + "' and password='" + password + "']/name"

        try:
            results = doc.xpath(expression)
            return [r.text for r in results]
        except Exception as e:
            raise ValueError(f"XPath query failed: {e}")

    def xpath_query(self, query: str) -> list:
        """VULNERABLE: XPath Injection (CWE-643)
        User-controlled XPath expression
        """
        doc = etree.fromstring(USERS_XML)

        # VULNERABLE: Direct user input as XPath expression
        try:
            results = doc.xpath(query)
            if isinstance(results, list):
                return [r.text if hasattr(r, "text") else str(r) for r in results]
            return [str(results)]
        except Exception as e:
            raise ValueError(f"XPath query failed: {e}")
