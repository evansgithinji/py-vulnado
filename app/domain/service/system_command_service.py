from app.domain.entity.command_request import CommandRequest
from app.domain.repository.command_policy_repository import CommandPolicyRepository
from app.domain.usecase.command_builder import CommandBuilder
from app.domain.usecase.shell_executor import ShellExecutor


class SystemCommandService:
    def __init__(
        self,
        policy_repo: CommandPolicyRepository,
        command_builder: CommandBuilder,
        shell_executor: ShellExecutor,
    ):
        self._policy_repo = policy_repo
        self._command_builder = command_builder
        self._shell_executor = shell_executor

    def execute_command(self, request: CommandRequest) -> str:
        """
        Deep call graph for Command Injection:
        Handler → SystemCommandService.execute_command()
          → CommandPolicyRepository.getPolicy() → CommandPolicy
          → CommandBuilder.buildCommand() → command string (VULNERABLE)
          → ShellExecutor.execute() → output (VULNERABLE)
        """
        policy = self._policy_repo.get_policy(request.command_type)
        command = self._command_builder.build_command(request.target, policy)
        return self._shell_executor.execute(command)

    def execute_backup(self, name: str, backups_dir: str) -> str:
        policy = self._policy_repo.get_policy("backup")
        command = self._command_builder.build_backup_command(name, backups_dir)
        self._shell_executor.execute(command)
        return f"Backup '{name}' created successfully"

    def execute_file_info(self, files_dir: str, filename: str) -> str:
        policy = self._policy_repo.get_policy("diagnostic")
        command = self._command_builder.build_file_info_command(files_dir, filename)
        return self._shell_executor.execute(command)
