class CalculatorUseCase:
    def evaluate(self, expression: str) -> str:
        # VULNERABLE: eval injection (CWE-94)
        # Direct eval on user-supplied expression
        result = eval(expression)
        return str(result)

    def calculate_discount(self, price: float, formula: str) -> str:
        # VULNERABLE: exec injection (CWE-94)
        local_vars = {"price": price}
        exec(f"result = {formula}", {}, local_vars)
        return str(local_vars.get("result", "error"))

    def batch_calculate(self, expressions: list[str]) -> list[str]:
        results = []
        for expr in expressions:
            # VULNERABLE: eval injection (CWE-94)
            result = eval(expr)
            results.append(str(result))
        return results
