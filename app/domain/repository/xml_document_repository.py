from abc import ABC, abstractmethod
from app.domain.entity.xml_document import XmlDocument


class XmlDocumentRepository(ABC):
    @abstractmethod
    def load_document(self, source: str) -> XmlDocument:
        ...


class InMemoryXmlDocumentRepository(XmlDocumentRepository):
    def __init__(self):
        self._documents = {
            "default": XmlDocument(
                content=b"""<?xml version="1.0" encoding="UTF-8"?>
<users>
  <user><username>admin</username><password>admin123</password><name>Admin</name><role>admin</role></user>
  <user><username>john</username><password>password</password><name>John Doe</name><role>user</role></user>
  <user><username>jane</username><password>password123</password><name>Jane Smith</name><role>user</role></user>
</users>""",
                source_id="default",
            ),
        }

    def load_document(self, source: str) -> XmlDocument:
        doc = self._documents.get(source, self._documents["default"])
        return doc
