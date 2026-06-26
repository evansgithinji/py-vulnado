from dataclasses import dataclass, field
from typing import Any


@dataclass
class DocumentCollection:
    name: str
    documents: list = field(default_factory=list)
    indexes: list = field(default_factory=list)


@dataclass
class NoSqlQuery:
    field: str
    operator: str
    value: Any
    logical_op: str = "AND"


@dataclass
class NoSqlAuthRequest:
    username: Any
    password: Any
    collection: str = "users"


@dataclass
class WhereQueryRequest:
    where_expression: str
    collection: str = "users"
