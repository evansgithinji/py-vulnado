from abc import ABC, abstractmethod
from app.domain.entity.review import Review


class ReviewRepository(ABC):
    @abstractmethod
    def find_by_product(self, product_id: int) -> list[Review]:
        pass

    @abstractmethod
    def create(self, review: Review) -> Review:
        pass
