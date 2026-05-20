from dataclasses import dataclass


@dataclass
class Product:
    id: int = 0
    name: str = ""
    price: float = 0.0
    stock: int = 0
    category: str = ""
