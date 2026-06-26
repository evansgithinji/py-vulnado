from abc import ABC, abstractmethod
from app.domain.entity.user import User


class UserRepository(ABC):
    @abstractmethod
    def find_by_credentials(self, username: str, password: str) -> User | None:
        pass

    @abstractmethod
    def find_by_username(self, username: str) -> User | None:
        pass

    @abstractmethod
    def search(self, query: str) -> list[User]:
        pass

    @abstractmethod
    def find_all(self) -> list[User]:
        pass
