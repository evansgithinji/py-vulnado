class SqlQueryExecutor:
    def __init__(self, connection):
        self._connection = connection

    def execute(self, sql: str) -> list:
        # VULNERABLE: Executes raw SQL without parameterization
        cursor = self._connection.cursor()
        cursor.execute(sql)
        columns = [desc[0] for desc in cursor.description] if cursor.description else []
        rows = cursor.fetchall()
        return [dict(zip(columns, row)) for row in rows]
