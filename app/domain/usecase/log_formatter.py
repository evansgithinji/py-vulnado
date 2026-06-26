import logging
from collections import deque
from app.domain.entity.audit_event import AuditEvent, AuditPolicy

logger = logging.getLogger(__name__)


class AuditEventEnricher:
    """Enriches audit events with metadata from policy.
    Taint flows through object fields."""

    def enrich(self, event: AuditEvent, policy: AuditPolicy) -> AuditEvent:
        """Adds metadata fields to event. Tainted data preserved in object fields."""
        event.metadata["policy_level"] = policy.level
        event.metadata["retention_days"] = str(policy.retention_days)
        if policy.include_source_ip and event.source_ip:
            event.metadata["client_ip"] = event.source_ip
        return event


class LogFormatter:
    """Formats audit events into log strings.
    VULNERABLE: Format string with user data."""

    def format(self, event: AuditEvent, format_template: str) -> str:
        """VULNERABLE: Log Injection (CWE-117)
        User data interpolated into format string without sanitization.
        Newlines in actor/details fields enable log forging."""
        return format_template.format(
            level=event.metadata.get("policy_level", "INFO"),
            timestamp=event.timestamp,
            event_type=event.event_type,
            actor=event.actor,
            details=event.details,
            source_ip=event.source_ip,
        )


class LogWriter:
    """Writes formatted log messages.
    VULNERABLE: Writes unsanitized messages to log."""

    def write(self, message: str, level: str) -> None:
        """VULNERABLE: Log Injection (CWE-117)
        Newline injection in sink - writes tainted message to log."""
        if level == "WARN":
            logger.warning(message)
        elif level == "ERROR":
            logger.error(message)
        else:
            logger.info(message)


class LogStorageAdapter:
    """Persists log events to in-memory store for GET /api/logs.
    VULNERABLE: Stored taint."""

    def __init__(self, max_size: int = 100):
        self._entries = deque(maxlen=max_size)

    def store(self, event: AuditEvent, formatted_message: str) -> None:
        """VULNERABLE: Stored taint - user-controlled data persisted."""
        self._entries.append({
            "timestamp": event.timestamp,
            "level": event.metadata.get("policy_level", "INFO"),
            "message": formatted_message,
        })

    def get_all(self) -> list:
        return list(self._entries)
