from app.domain.entity.command_request import CommandPolicy


class CommandPolicyRepository:
    def __init__(self):
        self._policies = {
            "ping": CommandPolicy(
                name="ping",
                allowed_commands=["ping"],
                timeout=10,
                allow_shell=True,
            ),
            "backup": CommandPolicy(
                name="backup",
                allowed_commands=["tar"],
                timeout=60,
                allow_shell=True,
            ),
            "diagnostic": CommandPolicy(
                name="diagnostic",
                allowed_commands=["echo"],
                timeout=30,
                allow_shell=True,
            ),
            "port_check": CommandPolicy(
                name="port_check",
                allowed_commands=["nc"],
                timeout=10,
                allow_shell=True,
            ),
        }

    def get_policy(self, name: str) -> CommandPolicy:
        return self._policies.get(name, self._policies["ping"])
