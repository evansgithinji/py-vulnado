from abc import ABC, abstractmethod
from app.domain.entity.document_collection import DocumentCollection


class DocumentCollectionRepository(ABC):
    @abstractmethod
    def get_collection(self, name: str) -> DocumentCollection:
        ...


class InMemoryDocumentCollectionRepository(DocumentCollectionRepository):
    def __init__(self):
        self._collections = {
            "users": DocumentCollection(
                name="users",
                documents=[
                    {"_id": "1", "username": "admin", "password": "admin123", "role": "admin", "email": "admin@example.com"},
                    {"_id": "2", "username": "john", "password": "john456", "role": "user", "email": "john@example.com"},
                    {"_id": "3", "username": "jane", "password": "jane789", "role": "user", "email": "jane@example.com"},
                ],
                indexes=["username", "email"],
            ),
        }

    def get_collection(self, name: str) -> DocumentCollection:
        return self._collections.get(name, DocumentCollection(name=name))
