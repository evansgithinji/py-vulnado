import os
from app.domain.entity.file_request import FileAccessPolicy


class PathResolver:
    def resolve_path(self, filename: str, base_dir: str, policy: FileAccessPolicy) -> str:
        # VULNERABLE: No path traversal check - allows ../ sequences
        return os.path.join(base_dir, filename)
