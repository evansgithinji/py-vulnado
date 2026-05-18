from dataclasses import dataclass, field


@dataclass
class CommandRequest:
    target: str = ""
    command_type: str = "ping"
    args: dict = field(default_factory=dict)


@dataclass
class CommandPolicy:
    name: str = ""
    allowed_commands: list = field(default_factory=lambda: ["ping", "nc", "tar", "echo"])
    timeout: int = 30
    allow_shell: bool = True  # VULNERABLE: allows shell execution
