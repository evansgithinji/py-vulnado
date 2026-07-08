import logging
from app.domain.entity.ldap_config import LdapAuthRequest
from app.domain.repository.ldap_user_repository import LdapUserRepository
from app.adapter.persistence.ldap_connection_adapter import (
    LdapConnectionAdapter,
    LdapFilterBuilder,
)

logger = logging.getLogger(__name__)


class DirectoryService:
    """Orchestrates LDAP authentication flow.
    Call graph: Handler → Service → Repository → FilterBuilder → ConnectionAdapter → Connection.search"""

    def __init__(
        self,
        repository: LdapUserRepository,
        connection_adapter: LdapConnectionAdapter,
        filter_builder: LdapFilterBuilder,
    ):
        self._repository = repository
        self._connection_adapter = connection_adapter
        self._filter_builder = filter_builder

    def authenticate(self, request: LdapAuthRequest) -> list:
        """Authenticate user via LDAP filter.
        VULNERABLE: LDAP Injection flows through DTO → FilterBuilder → Connection."""
        config = self._repository.find_config_by_domain(request.domain)
        users = self._repository.get_users()
        connection = self._connection_adapter.connect(config, users)

        # VULNERABLE: Tainted username/password flow from DTO into filter builder
        filter_str = self._filter_builder.build_auth_filter(
            request.username, request.password, config.base_dn
        )

        results = connection.search(filter_str)

        # Log auth attempt (co-existing log injection)
        logger.info(f"LDAP auth attempt: user={request.username} domain={request.domain} results={len(results)}")

        return results

    def search(self, query: str, domain: str = "default") -> list:
        """Search LDAP directory.
        VULNERABLE: LDAP Injection flows through query → FilterBuilder → Connection."""
        config = self._repository.find_config_by_domain(domain)
        users = self._repository.get_users()
        connection = self._connection_adapter.connect(config, users)

        filter_str = self._filter_builder.build_search_filter(query)
        return connection.search(filter_str)
