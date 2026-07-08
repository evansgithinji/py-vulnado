from dataclasses import dataclass


@dataclass
class RedirectRequest:
    target_url: str = "/"
    context: str = "auth"


@dataclass
class RedirectPolicy:
    name: str = ""
    allowed_domains: list = None  # VULNERABLE: no domain whitelist
    check_external: bool = False  # VULNERABLE: doesn't validate external URLs
