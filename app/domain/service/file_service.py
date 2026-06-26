from app.domain.entity.file_request import FileRequest
from app.domain.repository.file_access_policy_repository import FileAccessPolicyRepository
from app.domain.usecase.path_resolver import PathResolver
from app.domain.usecase.file_reader import FileReader


class FileService:
    def __init__(
        self,
        policy_repo: FileAccessPolicyRepository,
        path_resolver: PathResolver,
        file_reader: FileReader,
    ):
        self._policy_repo = policy_repo
        self._path_resolver = path_resolver
        self._file_reader = file_reader

    def process_file_request(self, request: FileRequest) -> str:
        """
        Deep call graph for Path Traversal:
        Handler → FileService.process_file_request()
          → FileAccessPolicyRepository.getPolicy() → FileAccessPolicy
          → PathResolver.resolvePath() → path (VULNERABLE: no .. check)
          → FileReader.readFile() → content (VULNERABLE: reads any path)
        """
        policy = self._policy_repo.get_policy(request.operation)
        path = self._path_resolver.resolve_path(
            request.filename, request.base_dir or policy.base_directory, policy
        )
        return self._file_reader.read_file(path)

    def download_file(self, request: FileRequest) -> bytes:
        policy = self._policy_repo.get_policy("download")
        path = self._path_resolver.resolve_path(
            request.filename, request.base_dir or policy.base_directory, policy
        )
        return self._file_reader.read_file_bytes(path)

    def upload_file(self, filename: str, content: bytes, upload_dir: str) -> str:
        import os
        policy = self._policy_repo.get_policy("upload")
        path = self._path_resolver.resolve_path(filename, upload_dir, policy)
        with open(path, "wb") as f:
            f.write(content)
        return path

    def list_files(self, directory: str) -> list:
        import os
        try:
            return os.listdir(directory)
        except OSError:
            return []
