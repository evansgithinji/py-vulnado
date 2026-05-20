from dataclasses import dataclass


@dataclass
class Message:
    id: int = 0
    content: str = ""
    author: str = ""
    created_at: str = ""
