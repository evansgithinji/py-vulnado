from abc import ABC, abstractmethod
from app.domain.entity.audit_event import AuditPolicy


class AuditPolicyRepository(ABC):
    @abstractmethod
    def get_policy(self, event_type: str) -> AuditPolicy:
        ...


class InMemoryAuditPolicyRepository(AuditPolicyRepository):
    def __init__(self):
        self._policies = {
            "login": AuditPolicy(
                format_template="[{level}] {timestamp} | AUTH | user={actor} status={details} ip={source_ip}",
                level="INFO",
                retention_days=90,
                include_source_ip=True,
            ),
            "search": AuditPolicy(
                format_template="[{level}] {timestamp} | SEARCH | query={details} actor={actor}",
                level="INFO",
                retention_days=30,
            ),
            "admin": AuditPolicy(
                format_template="[{level}] {timestamp} | ADMIN | action={details} actor={actor} ip={source_ip}",
                level="WARN",
                retention_days=365,
                include_source_ip=True,
            ),
        }
        self._default = AuditPolicy(
            format_template="[{level}] {timestamp} | {event_type} | actor={actor} details={details}",
            level="INFO",
        )

    def get_policy(self, event_type: str) -> AuditPolicy:
        return self._policies.get(event_type, self._default)
