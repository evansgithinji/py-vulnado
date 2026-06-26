import re


class NoSqlUseCase:
    def __init__(self):
        self.users = [
            {"_id": "1", "username": "admin", "password": "admin123", "role": "admin", "email": "admin@example.com"},
            {"_id": "2", "username": "john", "password": "john456", "role": "user", "email": "john@example.com"},
            {"_id": "3", "username": "jane", "password": "jane789", "role": "user", "email": "jane@example.com"},
        ]

    def _match_operator(self, field_val: str, condition) -> bool:
        if isinstance(condition, str):
            return field_val == condition
        if isinstance(condition, dict):
            for op, op_val in condition.items():
                if op == "$ne":
                    if field_val == str(op_val):
                        return False
                elif op == "$gt":
                    if not (field_val > str(op_val)):
                        return False
                elif op == "$regex":
                    if not re.search(str(op_val), field_val):
                        return False
            return True
        return False

    def authenticate(self, username, password) -> list:
        """VULNERABLE: NoSQL Injection (CWE-943)
        Operator injection via JSON objects
        """
        return [
            u for u in self.users
            if self._match_operator(u["username"], username)
            and self._match_operator(u["password"], password)
        ]

    def query(self, query_map: dict) -> list:
        """VULNERABLE: NoSQL Injection (CWE-943)"""
        results = []
        for user in self.users:
            match = True
            for field, condition in query_map.items():
                field_val = user.get(field, "")
                if not self._match_operator(field_val, condition):
                    match = False
                    break
            if match:
                results.append(user)
        return results

    def find_users_where(self, where_clause: str) -> list:
        """VULNERABLE: Code injection via eval() $where clause"""
        results = []
        for user in self.users:
            try:
                # VULNERABLE: User input evaluated via eval()
                if eval(where_clause, {"__builtins__": {}}, user):
                    results.append(user)
            except Exception:
                pass
        return results
