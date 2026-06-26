class HeaderUseCase:
    def set_custom_header(self, name: str, value: str) -> tuple:
        """VULNERABLE: HTTP Header Injection / CRLF (CWE-113)
        No sanitization of user input before setting as header value
        """
        # VULNERABLE: CRLF injection - user controls header value without sanitization
        return name, value

    def get_redirect_url(self, lang: str) -> str:
        """VULNERABLE: HTTP Header Injection / CRLF (CWE-113)"""
        # VULNERABLE: CRLF injection via lang parameter
        return f"/?lang={lang}"
