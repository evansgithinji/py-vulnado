from dataclasses import dataclass, field
from typing import Any, Optional


@dataclass
class XmlDocument:
    content: bytes
    source_id: str
    parsed_doc: Any = None


@dataclass
class XPathAuthRequest:
    username: str
    password: str
    document_source: str = "default"
