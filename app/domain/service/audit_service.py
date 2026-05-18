from app.domain.entity.audit_event import AuditEvent
from app.domain.repository.audit_policy_repository import AuditPolicyRepository
from app.domain.usecase.log_formatter import (
    AuditEventEnricher,
    LogFormatter,
    LogWriter,
    LogStorageAdapter,
)


class AuditService:
    """Orchestrates audit event recording flow.
    Call graph: Handler → Service → PolicyRepository → Enricher → Formatter → Writer → StorageAdapter."""

    def __init__(
        self,
        policy_repo: AuditPolicyRepository,
        enricher: AuditEventEnricher,
        formatter: LogFormatter,
        writer: LogWriter,
        storage: LogStorageAdapter,
    ):
        self._policy_repo = policy_repo
        self._enricher = enricher
        self._formatter = formatter
        self._writer = writer
        self._storage = storage

    def record_event(self, event: AuditEvent) -> str:
        """Record an audit event through the full pipeline.
        VULNERABLE: Log Injection flows through DTO → Enricher → Formatter → Writer → Storage."""
        policy = self._policy_repo.get_policy(event.event_type)

        enriched_event = self._enricher.enrich(event, policy)

        formatted_message = self._formatter.format(enriched_event, policy.format_template)

        self._writer.write(formatted_message, policy.level)

        self._storage.store(enriched_event, formatted_message)

        return formatted_message

    def get_logs(self) -> list:
        return self._storage.get_all()
