from abc import ABC, abstractmethod
from app.domain.entity.message import Message


class MessageRepository(ABC):
    @abstractmethod
    def find_all(self) -> list[Message]:
        pass

    @abstractmethod
    def create(self, message: Message) -> Message:
        pass
