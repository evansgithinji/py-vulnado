from dataclasses import dataclass


@dataclass
class User:
    id: int = 0
    username: str = ""
    password: str = ""
    email: str = ""
    role: str = "user"
    is_admin: bool = False
