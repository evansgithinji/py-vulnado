from dataclasses import dataclass, field


@dataclass
class HttpRequestDto:
    url: str = ""
    method: str = "GET"
    headers: dict = field(default_factory=dict)
    body: str = ""
    timeout: int = 10


@dataclass
class UrlPolicy:
    name: str = ""
    blocked_hosts: list = None  # VULNERABLE: no blocked hosts configured
    allowed_schemes: list = None
    check_internal: bool = False  # VULNERABLE: doesn't block internal IPs


@dataclass
class PreparedRequest:
    url: str = ""
    method: str = "GET"
    headers: dict = field(default_factory=dict)
    body: str = ""
    timeout: int = 10
    policy: UrlPolicy = None
