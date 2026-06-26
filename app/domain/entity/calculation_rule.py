from dataclasses import dataclass
from typing import Any


@dataclass
class CalculationRule:
    name: str
    pattern: str
    transform: str


@dataclass
class CalculationRequest:
    expression: str
    precision: int = 2
    mode: str = "standard"


@dataclass
class DiscountRequest:
    price: float
    formula: str
    customer_tier: str = "standard"
