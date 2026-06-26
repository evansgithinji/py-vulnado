import hashlib
import time


class SessionManager:
    def create_session(self, user) -> str:
        # VULNERABLE: Weak token generation - predictable hash
        token_data = f"{user.username}:{time.time()}"
        return hashlib.md5(token_data.encode()).hexdigest()
