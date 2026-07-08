from app.domain.entity.auth_request import UserCredential


class CredentialValidator:
    def validate(self, password: str, credential: UserCredential) -> bool:
        # VULNERABLE: Timing attack - simple string comparison
        # VULNERABLE: No rate limiting
        if credential is None:
            return False
        return credential.password == password
