from abc import ABC, abstractmethod
from app.domain.entity.product import Product


class ProductRepository(ABC):
    @abstractmethod
    def find_all(self) -> list[Product]:
        pass

    @abstractmethod
    def find_by_id(self, product_id: int) -> Product | None:
        pass

    @abstractmethod
    def search(self, query: str) -> list[Product]:
        pass

    @abstractmethod
    def search_by_category(self, category: str) -> list[Product]:
        pass
