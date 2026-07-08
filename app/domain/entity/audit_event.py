from dataclasses import dataclass, field
from typing import Optional
from datetime import datetime, timezone


@dataclass
class AuditEvent:
    event_type: str
    actor: str
    details: str
    timestamp: str = ""
    source_ip: str = ""
    metadata: dict = field(default_factory=dict)

    def __post_init__(self):
        if not self.timestamp:
            self.timestamp = datetime.now(timezone.utc).isoformat()


@dataclass
class AuditPolicy:
    format_template: str
    level: str
    retention_days: int = 30
    include_source_ip: bool = True
