from app.domain.entity.content_submission import ContentPolicy


class ContentProcessor:
    def process_content(self, content: str, policy: ContentPolicy) -> str:
        # VULNERABLE: No HTML sanitization - passes through raw HTML/script tags
        if policy.max_length and len(content) > policy.max_length:
            content = content[:policy.max_length]
        # Fake "processing" - doesn't actually sanitize
        return content
