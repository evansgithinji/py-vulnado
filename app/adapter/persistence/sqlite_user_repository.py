import sqlite3
from app.domain.entity.user import User
from app.domain.repository.user_repository import UserRepository


class SQLiteUserRepository(UserRepository):
    def __init__(self, conn: sqlite3.Connection):
        self.conn = conn

    def find_by_credentials(self, username: str, password: str) -> User | None:
        # VULNERABLE: SQL Injection (CWE-89)
        query = f"SELECT * FROM users WHERE username='{username}' AND password='{password}'"
        cursor = self.conn.execute(query)
        row = cursor.fetchone()
        if row:
            return User(
                id=row["id"],
                username=row["username"],
                password=row["password"],
                email=row["email"],
                role=row["role"],
                is_admin=bool(row["is_admin"]),
            )
        return None

    def find_by_username(self, username: str) -> User | None:
        # VULNERABLE: SQL Injection (CWE-89)
        query = f"SELECT * FROM users WHERE username='{username}'"
        cursor = self.conn.execute(query)
        row = cursor.fetchone()
        if row:
            return User(
                id=row["id"],
                username=row["username"],
                password=row["password"],
                email=row["email"],
                role=row["role"],
                is_admin=bool(row["is_admin"]),
            )
        return None

    def search(self, query: str) -> list[User]:
        # VULNERABLE: SQL Injection (CWE-89)
        sql = f"SELECT * FROM users WHERE username LIKE '%{query}%' OR email LIKE '%{query}%'"
        cursor = self.conn.execute(sql)
        users = []
        for row in cursor.fetchall():
            users.append(
                User(
                    id=row["id"],
                    username=row["username"],
                    password=row["password"],
                    email=row["email"],
                    role=row["role"],
                    is_admin=bool(row["is_admin"]),
                )
            )
        return users

    def find_all(self) -> list[User]:
        cursor = self.conn.execute("SELECT * FROM users")
        users = []
        for row in cursor.fetchall():
            users.append(
                User(
                    id=row["id"],
                    username=row["username"],
                    password=row["password"],
                    email=row["email"],
                    role=row["role"],
                    is_admin=bool(row["is_admin"]),
                )
            )
        return users
