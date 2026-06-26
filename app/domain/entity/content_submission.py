from dataclasses import dataclass


@dataclass
class ContentSubmission:
    content: str = ""
    author: str = "anonymous"
    content_type: str = "html"
    context: str = "message"


@dataclass
class ContentPolicy:
    name: str = ""
    allow_html: bool = True  # VULNERABLE: allows HTML
    sanitize: bool = False  # VULNERABLE: no sanitization
    max_length: int = 10000
