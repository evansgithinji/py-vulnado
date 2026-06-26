from app.domain.entity.command_request import CommandPolicy


class CommandBuilder:
    def build_command(self, target: str, policy: CommandPolicy) -> str:
        # VULNERABLE: String concatenation for shell command building
        command_type = policy.name
        if command_type == "ping":
            return f"ping -c 1 {target}"
        elif command_type == "backup":
            return f"tar czf {target}"
        elif command_type == "diagnostic":
            return f"echo 'Running diagnostic: {target}'"
        elif command_type == "port_check":
            return f"nc -zv -w 3 {target} 2>&1 || true"
        else:
            return f"echo {target}"

    def build_backup_command(self, name: str, backups_dir: str) -> str:
        # VULNERABLE: String concatenation
        return f"tar czf {backups_dir}/{name}.tar.gz -C /app/data ."

    def build_file_info_command(self, files_dir: str, filename: str) -> str:
        # VULNERABLE: String concatenation
        return f"file {files_dir}/{filename}"
