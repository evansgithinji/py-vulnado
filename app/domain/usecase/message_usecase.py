from app.domain.repository.message_repository import MessageRepository
from app.domain.entity.message import Message


class MessageUseCase:
    def __init__(self, message_repo: MessageRepository):
        self.message_repo = message_repo

    def get_all_messages(self):
        return self.message_repo.find_all()

    def create_message(self, content: str, author: str):
        message = Message(content=content, author=author)
        return self.message_repo.create(message)
