import subprocess
import os


class AdminUseCase:
    def __init__(self, backups_dir: str):
        self.backups_dir = backups_dir

    def create_backup(self, name: str) -> str:
        # VULNERABLE: Command Injection (CWE-78)
        output = subprocess.check_output(
            f"tar czf {self.backups_dir}/{name}.tar.gz -C /app/data .",
            shell=True,
            stderr=subprocess.STDOUT,
        )
        return f"Backup '{name}' created successfully"

    def get_system_info(self) -> dict:
        hostname = subprocess.check_output("hostname", shell=True).decode().strip()
        uname = subprocess.check_output("uname -a", shell=True).decode().strip()
        return {
            "hostname": hostname,
            "system": uname,
            "working_dir": os.getcwd(),
        }

    def run_diagnostic(self, command: str) -> str:
        # VULNERABLE: Command Injection (CWE-78)
        output = subprocess.check_output(
            f"echo 'Running diagnostic: {command}'",
            shell=True,
            stderr=subprocess.STDOUT,
        )
        return output.decode("utf-8")
