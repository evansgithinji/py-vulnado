class LdapUseCase:
    def __init__(self):
        self.users = [
            {"uid": "admin", "cn": "Admin User", "mail": "admin@example.com", "userPassword": "admin123", "ou": "admins"},
            {"uid": "john", "cn": "John Doe", "mail": "john@example.com", "userPassword": "john456", "ou": "users"},
            {"uid": "jane", "cn": "Jane Smith", "mail": "jane@example.com", "userPassword": "jane789", "ou": "users"},
            {"uid": "svc_backup", "cn": "Backup Service", "mail": "backup@example.com", "userPassword": "backup!@#", "ou": "services"},
        ]

    def _match_filter(self, user: dict, filt: str) -> bool:
        filt = filt.strip()
        while filt.startswith("(") and filt.endswith(")"):
            filt = filt[1:-1].strip()

        # AND filter
        if filt.startswith("&"):
            conditions = self._split_conditions(filt[1:])
            return all(self._match_filter(user, c) for c in conditions)

        # OR filter
        if filt.startswith("|"):
            conditions = self._split_conditions(filt[1:])
            return any(self._match_filter(user, c) for c in conditions)

        # Simple attribute=value
        eq_idx = filt.find("=")
        if eq_idx == -1:
            return True  # malformed = match all (vulnerable)

        attr = filt[:eq_idx].strip()
        val = filt[eq_idx + 1:].strip()

        if val == "*":
            return True

        user_val = user.get(attr, "") or user.get(attr.lower(), "")
        return user_val.lower() == val.lower()

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

    def authenticate(self, username: str, password: str) -> list:
        """VULNERABLE: LDAP Injection (CWE-90)
        String concatenation for LDAP filter construction
        """
        # VULNERABLE: Direct string concatenation into LDAP filter
        filt = f"(&(uid={username})(userPassword={password}))"

        return [u for u in self.users if self._match_filter(u, filt)]

    def search(self, query: str) -> list:
        """VULNERABLE: LDAP Injection (CWE-90)"""
        # VULNERABLE: Direct string concatenation into LDAP filter
        filt = f"(|(uid={query})(cn={query})(mail={query}))"

        return [u for u in self.users if self._match_filter(u, filt)]
