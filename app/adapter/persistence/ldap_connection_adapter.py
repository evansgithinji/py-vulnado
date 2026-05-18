from app.domain.entity.ldap_config import LdapConfig


class LdapConnection:
    """Simulated LDAP connection that executes filters against in-memory data."""

    def __init__(self, config: LdapConfig, users: list):
        self._config = config
        self._users = users

    def search(self, filter_str: str) -> list:
        """VULNERABLE: Executes LDAP filter string against user store.
        The filter_str may contain injected LDAP operators."""
        return [u for u in self._users if self._match_filter(u, filter_str)]

    def _match_filter(self, user: dict, filt: str) -> bool:
        filt = filt.strip()
        while filt.startswith("(") and filt.endswith(")"):
            filt = filt[1:-1].strip()

        if filt.startswith("&"):
            conditions = self._split_conditions(filt[1:])
            return all(self._match_filter(user, c) for c in conditions)

        if filt.startswith("|"):
            conditions = self._split_conditions(filt[1:])
            return any(self._match_filter(user, c) for c in conditions)

        eq_idx = filt.find("=")
        if eq_idx == -1:
            return True

        attr = filt[:eq_idx].strip()
        val = filt[eq_idx + 1:].strip()

        if val == "*":
            return True

        user_val = user.get(attr, "") or user.get(attr.lower(), "")
        return str(user_val).lower() == val.lower()

    def _split_conditions(self, s: str) -> list:
        conditions = []
        depth = 0
        start = -1
        for i, ch in enumerate(s):
            if ch == "(":
                if depth == 0:
                    start = i
                depth += 1
            elif ch == ")":
                depth -= 1
                if depth == 0 and start >= 0:
                    conditions.append(s[start:i + 1])
                    start = -1
        return conditions


class LdapConnectionAdapter:
    """Creates LDAP connections from config."""

    def connect(self, config: LdapConfig, users: list) -> LdapConnection:
        return LdapConnection(config, users)


class LdapFilterBuilder:
    """Builds LDAP filter strings from user input.
    VULNERABLE: String concatenation without sanitization."""

    def build_auth_filter(self, username: str, password: str, base_dn: str) -> str:
        """VULNERABLE: LDAP Injection (CWE-90)
        Direct string concatenation into LDAP filter."""
        return f"(&(uid={username})(userPassword={password}))"

    def build_search_filter(self, query: str) -> str:
        """VULNERABLE: LDAP Injection (CWE-90)"""
        return f"(|(uid={query})(cn={query})(mail={query}))"
