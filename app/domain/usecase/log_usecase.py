import logging
from datetime import datetime, timezone
from collections import deque

logger = logging.getLogger(__name__)


class LogUseCase:
    def __init__(self):
        self.entries = deque(maxlen=100)

    def _add_entry(self, level: str, message: str):
        entry = {
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "level": level,
            "message": message,
        }
        self.entries.append(entry)

    def log_login(self, username: str, status: str) -> str:
        """VULNERABLE: Log Injection (CWE-117)
        User input written to log without sanitization - newline injection allows log forging
        """
        # VULNERABLE: newline injection allows log forging
        message = f"Login attempt: user={username} status={status}"
        logger.info(message)
        self._add_entry("INFO", message)
        return message

    def log_search(self, query: str) -> str:
        """VULNERABLE: Log Injection (CWE-117)
        User input written to log without sanitization
        """
        # VULNERABLE: newline injection allows log forging
        message = f"Search query: q={query}"
        logger.info(message)
        self._add_entry("INFO", message)
        return message

    def get_logs(self) -> list:
        return list(self.entries)
