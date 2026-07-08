class FileReader:
    def read_file(self, path: str) -> str:
        # VULNERABLE: Reads any file path without restriction
        with open(path, "r") as f:
            return f.read()

    def read_file_bytes(self, path: str) -> bytes:
        # VULNERABLE: Reads any file path without restriction
        with open(path, "rb") as f:
            return f.read()
