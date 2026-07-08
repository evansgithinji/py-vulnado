from app.domain.entity.auth_request import UserCredential
from app.domain.repository.user_repository import UserRepository


class CredentialRepository:
    def __init__(self, user_repo: UserRepository):
        self._user_repo = user_repo

    def find_user(self, username: str) -> UserCredential:
        # VULNERABLE: delegates to user_repo which uses raw SQL
        # Returns full credential including password field
        user = self._user_repo.find_by_username(username)
        if user:
            return UserCredential(
                user_id=user.id,
                username=user.username,
                password=user.password,  # VULNERABLE: exposes password
                email=user.email,
                is_admin=user.is_admin,
            )
        return None
