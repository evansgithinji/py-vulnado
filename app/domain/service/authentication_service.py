from app.domain.entity.auth_request import AuthRequest
from app.domain.repository.credential_repository import CredentialRepository
from app.domain.usecase.credential_validator import CredentialValidator
from app.domain.usecase.session_manager import SessionManager


class AuthenticationService:
    def __init__(
        self,
        credential_repo: CredentialRepository,
        validator: CredentialValidator,
        session_manager: SessionManager,
    ):
        self._credential_repo = credential_repo
        self._validator = validator
        self._session_manager = session_manager

    def authenticate(self, request: AuthRequest):
        """
        Deep call graph for Broken Auth / Info Disclosure:
        Handler → AuthenticationService.authenticate()
          → CredentialRepository.findUser() → UserCredential (VULNERABLE: returns password)
          → CredentialValidator.validate() → boolean (VULNERABLE: timing attack)
          → SessionManager.createSession() → token (VULNERABLE: weak token)
        """
        credential = self._credential_repo.find_user(request.username)
        if credential is None:
            return None

        is_valid = self._validator.validate(request.password, credential)
        if not is_valid:
            return None

        token = self._session_manager.create_session(credential)
        return {
            "user_id": credential.user_id,
            "username": credential.username,
            "email": credential.email,
            "is_admin": credential.is_admin,
            "token": token,
        }
