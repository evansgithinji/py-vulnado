from app.domain.entity.redirect_request import RedirectPolicy


class RedirectPolicyRepository:
    def __init__(self):
        self._policies = {
            "auth": RedirectPolicy(
                name="auth",
                allowed_domains=None,  # VULNERABLE: no domain whitelist
                check_external=False,  # VULNERABLE
            ),
            "admin": RedirectPolicy(
                name="admin",
                allowed_domains=None,
                check_external=False,
            ),
        }

    def get_policy(self, name: str) -> RedirectPolicy:
        return self._policies.get(name, self._policies["auth"])
