from dataclasses import dataclass
from typing import Optional


@dataclass
class HeaderPolicy:
    name: str
    allowed_pattern: Optional[str]
    max_length: int
    allow_raw: bool = True


@dataclass
class HeaderRequest:
    name: str
    value: str
    context: str = ""


@dataclass
class LocaleConfig:
    code: str
    name: str
    direction: str = "ltr"


@dataclass
class LocaleRequest:
    lang: str
    return_url: str = "/"
