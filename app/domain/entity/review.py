from dataclasses import dataclass


@dataclass
class Review:
    id: int = 0
    product_id: int = 0
    user_id: int = 0
    rating: int = 0
    comment: str = ""
