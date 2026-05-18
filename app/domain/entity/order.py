from dataclasses import dataclass


@dataclass
class Order:
    id: int = 0
    user_id: int = 0
    product_id: int = 0
    quantity: int = 1
    total: float = 0.0
    status: str = "pending"
