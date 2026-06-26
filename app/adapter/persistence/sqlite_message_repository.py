import sqlite3
from app.domain.entity.message import Message
from app.domain.repository.message_repository import MessageRepository


class SQLiteMessageRepository(MessageRepository):
    def __init__(self, conn: sqlite3.Connection):
        self.conn = conn

    def find_all(self) -> list[Message]:
        cursor = self.conn.execute(
            "SELECT * FROM messages ORDER BY created_at DESC"
        )
        messages = []
        for row in cursor.fetchall():
            messages.append(
                Message(
                    id=row["id"],
                    content=row["content"],
                    author=row["author"],
                    created_at=row["created_at"] or "",
                )
            )
        return messages

    def create(self, message: Message) -> Message:
        cursor = self.conn.execute(
            "INSERT INTO messages (content, author) VALUES (?, ?)",
            (message.content, message.author),
        )
        self.conn.commit()
        message.id = cursor.lastrowid
        return message
