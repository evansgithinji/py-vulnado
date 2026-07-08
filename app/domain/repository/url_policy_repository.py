from app.domain.entity.http_request_dto import UrlPolicy


class UrlPolicyRepository:
    def __init__(self):
        self._policies = {
            "fetch": UrlPolicy(
                name="fetch",
                blocked_hosts=None,  # VULNERABLE: no blocked hosts
                allowed_schemes=["http", "https"],
                check_internal=False,  # VULNERABLE
            ),
            "webhook": UrlPolicy(
                name="webhook",
                blocked_hosts=None,
                allowed_schemes=["http", "https"],
                check_internal=False,
            ),
            "proxy": UrlPolicy(
                name="proxy",
                blocked_hosts=None,
                allowed_schemes=["http", "https"],
                check_internal=False,
            ),
            "remote_file": UrlPolicy(
                name="remote_file",
                blocked_hosts=None,
                allowed_schemes=["http", "https", "file"],
                check_internal=False,
            ),
        }

    def get_policy(self, name: str) -> UrlPolicy:
        return self._policies.get(name, self._policies["fetch"])
