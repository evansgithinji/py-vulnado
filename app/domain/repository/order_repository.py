from abc import ABC, abstractmethod
from app.domain.entity.order import Order


class OrderRepository(ABC):
    @abstractmethod
    def find_by_user(self, user_id: str) -> list[Order]:
        pass

    @abstractmethod
    def create(self, order: Order) -> Order:
        pass

    @abstractmethod
    def find_all(self) -> list[Order]:
        pass
