from app.domain.entity.header_policy import HeaderRequest, LocaleRequest
from app.domain.repository.header_policy_repository import (
    HeaderPolicyRepository,
    LocaleRepository,
)
from app.domain.usecase.header_processor import (
    HeaderValueProcessor,
    ResponseHeaderWriter,
    CookieManager,
    RedirectBuilder,
)


class ResponseCustomizationService:
    """Orchestrates custom header setting flow.
    Call graph: Handler → Service → PolicyRepository → ValueProcessor → HeaderWriter."""

    def __init__(
        self,
        policy_repo: HeaderPolicyRepository,
        value_processor: HeaderValueProcessor,
        header_writer: ResponseHeaderWriter,
    ):
        self._policy_repo = policy_repo
        self._value_processor = value_processor
        self._header_writer = header_writer

    def apply_custom_header(self, request: HeaderRequest) -> tuple:
        """Apply custom header from request.
        VULNERABLE: CRLF injection flows through DTO → Processor → Writer."""
        policy = self._policy_repo.get_policy(request.name)

        processed_name = self._value_processor.process_name(request.name, policy)
        processed_value = self._value_processor.process_value(request.value, policy)

        return self._header_writer.build_header_pair(processed_name, processed_value)


class LocalizationService:
    """Orchestrates locale-based redirect flow.
    Call graph: Handler → Service → LocaleRepository → CookieManager → RedirectBuilder."""

    def __init__(
        self,
        locale_repo: LocaleRepository,
        cookie_manager: CookieManager,
        redirect_builder: RedirectBuilder,
    ):
        self._locale_repo = locale_repo
        self._cookie_manager = cookie_manager
        self._redirect_builder = redirect_builder

    def build_localized_redirect(self, request: LocaleRequest) -> dict:
        """Build localized redirect with cookie.
        VULNERABLE: Header injection via unsanitized lang → locale → cookie/redirect."""
        locale = self._locale_repo.resolve_locale(request.lang)

        cookie_name, cookie_value = self._cookie_manager.build_locale_cookie(locale)
        redirect_url = self._redirect_builder.build_redirect_url(request.return_url, locale)

        return {
            "redirect_url": redirect_url,
            "cookie_name": cookie_name,
            "cookie_value": cookie_value,
            "locale": locale,
        }
