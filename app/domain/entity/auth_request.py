from dataclasses import dataclass
from typing import Optional


@dataclass
class AuthRequest:
    username: str = ""
    password: str = ""


@dataclass
class UserCredential:
    user_id: int = 0
    username: str = ""
    password: str = ""  # VULNERABLE: stores plaintext password
    email: str = ""
    is_admin: bool = False
