from app.domain.entity.header_policy import HeaderPolicy, LocaleConfig


class HeaderValueProcessor:
    """Processes header values according to policy.
    VULNERABLE: No CRLF sanitization."""

    def process_value(self, value: str, policy: HeaderPolicy) -> str:
        """VULNERABLE: HTTP Header Injection (CWE-113)
        No CRLF sanitization - user input passed through as-is."""
        if policy.max_length and len(value) > policy.max_length:
            value = value[:policy.max_length]
        # VULNERABLE: No newline/CR/LF stripping
        return value

    def process_name(self, name: str, policy: HeaderPolicy) -> str:
        """VULNERABLE: Header name passed through without sanitization."""
        return name


class ResponseHeaderWriter:
    """Writes headers to HTTP response.
    VULNERABLE: Writes unsanitized values to response headers."""

    def build_header_pair(self, name: str, value: str) -> tuple:
        """VULNERABLE: HTTP Header Injection (CWE-113)
        Returns header name/value pair without CRLF sanitization."""
        return name, value


class CookieManager:
    """Manages cookie operations.
    VULNERABLE: Cookie injection via unsanitized locale."""

    def build_locale_cookie(self, locale: LocaleConfig) -> tuple:
        """VULNERABLE: HTTP Header Injection (CWE-113)
        Cookie value from unsanitized locale code."""
        cookie_name = "locale"
        cookie_value = locale.code
        return cookie_name, cookie_value


class RedirectBuilder:
    """Builds redirect URLs.
    VULNERABLE: Open redirect via unsanitized input."""

    def build_redirect_url(self, return_url: str, locale: LocaleConfig) -> str:
        """VULNERABLE: Open redirect + CRLF injection (CWE-113)"""
        return f"{return_url}?lang={locale.code}"
