import subprocess


class ShellExecutor:
    def execute(self, command: str) -> str:
        # VULNERABLE: Executes shell command with shell=True
        output = subprocess.check_output(
            command, shell=True, stderr=subprocess.STDOUT
        )
        return output.decode("utf-8")
