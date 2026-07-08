import os
import subprocess
import requests


class FileUseCase:
    def __init__(self, files_dir: str, upload_dir: str):
        self.files_dir = files_dir
        self.upload_dir = upload_dir

    def read_file(self, filename: str) -> str:
        # VULNERABLE: Path Traversal (CWE-22)
        # No validation on filename - allows ../../etc/passwd
        filepath = os.path.join(self.files_dir, filename)
        with open(filepath, "r") as f:
            return f.read()

    def download_file(self, filename: str) -> bytes:
        # VULNERABLE: Path Traversal (CWE-22)
        filepath = os.path.join(self.files_dir, filename)
        with open(filepath, "rb") as f:
            return f.read()

    def upload_file(self, filename: str, content: bytes) -> str:
        filepath = os.path.join(self.upload_dir, filename)
        with open(filepath, "wb") as f:
            f.write(content)
        return filepath

    def list_files(self) -> list[str]:
        try:
            return os.listdir(self.files_dir)
        except OSError:
            return []

    def fetch_remote_file(self, url: str) -> str:
        # VULNERABLE: SSRF (CWE-918)
        response = requests.get(url, timeout=10)
        return response.text

    def file_info(self, filename: str) -> str:
        # VULNERABLE: Command Injection (CWE-78)
        output = subprocess.check_output(
            f"file {self.files_dir}/{filename}", shell=True
        )
        return output.decode("utf-8")
