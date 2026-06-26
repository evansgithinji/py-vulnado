from abc import ABC, abstractmethod
from app.domain.entity.ldap_config import LdapConfig


class LdapUserRepository(ABC):
    @abstractmethod
    def find_config_by_domain(self, domain: str) -> LdapConfig:
        ...

    @abstractmethod
    def get_users(self) -> list:
        ...


class InMemoryLdapUserRepository(LdapUserRepository):
    def __init__(self):
        self._configs = {
            "default": LdapConfig(
                base_dn="dc=example,dc=com",
                host="ldap.internal",
                port=389,
                bind_dn="cn=admin,dc=example,dc=com",
                admin_password="ldap_secret_123",
            ),
            "corp": LdapConfig(
                base_dn="dc=corp,dc=example,dc=com",
                host="ldap-corp.internal",
                port=636,
                bind_dn="cn=svc,dc=corp,dc=example,dc=com",
                admin_password="corp_ldap_pass!",
            ),
        }
        self._users = [
            {"uid": "admin", "cn": "Admin User", "mail": "admin@example.com", "userPassword": "admin123", "ou": "admins"},
            {"uid": "john", "cn": "John Doe", "mail": "john@example.com", "userPassword": "john456", "ou": "users"},
            {"uid": "jane", "cn": "Jane Smith", "mail": "jane@example.com", "userPassword": "jane789", "ou": "users"},
            {"uid": "svc_backup", "cn": "Backup Service", "mail": "backup@example.com", "userPassword": "backup!@#", "ou": "services"},
        ]

    def find_config_by_domain(self, domain: str) -> LdapConfig:
        return self._configs.get(domain, self._configs["default"])

    def get_users(self) -> list:
        return list(self._users)
