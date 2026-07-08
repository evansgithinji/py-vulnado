from dataclasses import dataclass, field
from typing import Optional


@dataclass
class SearchRequest:
    query: str = ""
    category: str = ""
    sort_by: str = ""
    order: str = "asc"


@dataclass
class SqlQueryPolicy:
    name: str = ""
    allowed_tables: list = field(default_factory=lambda: ["products"])
    max_results: int = 100
    allow_wildcards: bool = True
    allow_union: bool = True  # VULNERABLE: allows UNION queries


@dataclass
class ValidatedQuery:
    original_query: str = ""
    table: str = "products"
    policy: Optional[SqlQueryPolicy] = None
