from app.domain.entity.redirect_request import RedirectPolicy


class UrlResolver:
    def resolve_url(self, target_url: str, policy: RedirectPolicy) -> str:
        # VULNERABLE: No validation of external URLs - allows open redirect
        if not target_url:
            return "/"
        return target_url
