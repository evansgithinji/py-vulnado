from app.domain.repository.user_repository import UserRepository


class AuthUseCase:
    def __init__(self, user_repo: UserRepository):
        self.user_repo = user_repo

    def login(self, username: str, password: str):
        return self.user_repo.find_by_credentials(username, password)

    def search_users(self, query: str):
        return self.user_repo.search(query)

    def get_all_users(self):
        return self.user_repo.find_all()
