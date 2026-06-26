from app.domain.entity.content_submission import ContentPolicy


class ContentPolicyRepository:
    def __init__(self):
        self._policies = {
            "html": ContentPolicy(
                name="html",
                allow_html=True,  # VULNERABLE
                sanitize=False,  # VULNERABLE
                max_length=10000,
            ),
            "message": ContentPolicy(
                name="message",
                allow_html=True,
                sanitize=False,
                max_length=5000,
            ),
            "review": ContentPolicy(
                name="review",
                allow_html=True,
                sanitize=False,
                max_length=2000,
            ),
        }

    def get_policy(self, name: str) -> ContentPolicy:
        return self._policies.get(name, self._policies["html"])
