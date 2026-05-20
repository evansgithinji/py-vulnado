from dataclasses import dataclass


@dataclass
class FileRequest:
    filename: str = ""
    base_dir: str = ""
    operation: str = "read"


@dataclass
class FileAccessPolicy:
    name: str = ""
    base_directory: str = ""
    allowed_extensions: list = None
    check_traversal: bool = False  # VULNERABLE: traversal check disabled
