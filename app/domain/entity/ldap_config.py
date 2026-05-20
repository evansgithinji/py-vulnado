from dataclasses import dataclass


@dataclass
class LdapConfig:
    base_dn: str
    host: str
    port: int
    bind_dn: str
    admin_password: str


@dataclass
class LdapAuthRequest:
    username: str
    password: str
    domain: str = "default"
