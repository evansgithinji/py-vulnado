from app.domain.entity.file_request import FileAccessPolicy


class FileAccessPolicyRepository:
    def __init__(self, files_dir: str, upload_dir: str):
        self._policies = {
            "read": FileAccessPolicy(
                name="read",
                base_directory=files_dir,
                allowed_extensions=None,
                check_traversal=False,  # VULNERABLE
            ),
            "download": FileAccessPolicy(
                name="download",
                base_directory=files_dir,
                allowed_extensions=None,
                check_traversal=False,  # VULNERABLE
            ),
            "upload": FileAccessPolicy(
                name="upload",
                base_directory=upload_dir,
                allowed_extensions=None,
                check_traversal=False,
            ),
            "info": FileAccessPolicy(
                name="info",
                base_directory=files_dir,
                allowed_extensions=None,
                check_traversal=False,
            ),
        }

    def get_policy(self, name: str) -> FileAccessPolicy:
        return self._policies.get(name, self._policies["read"])
