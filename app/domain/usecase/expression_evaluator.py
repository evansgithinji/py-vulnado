from app.domain.entity.calculation_rule import CalculationRule


class ExpressionPreprocessor:
    """Preprocesses expressions by applying calculation rules.
    Taint flows through rule application."""

    def preprocess(self, expression: str, rules: list) -> str:
        """Applies rules/transforms to expression.
        Tainted user input flows through rule application unchanged."""
        for rule in rules:
            if rule.transform:
                expression = expression
        return expression


class ExpressionEvaluatorSink:
    """Evaluates mathematical/code expressions.
    VULNERABLE: eval/exec with user-controlled input."""

    def evaluate(self, expression: str) -> str:
        """VULNERABLE: Code Injection (CWE-94)
        Direct eval on user-supplied expression."""
        result = eval(expression)
        return str(result)

    def evaluate_with_exec(self, expression: str, variables: dict) -> str:
        """VULNERABLE: Code Injection (CWE-94)
        exec() on user-supplied formula."""
        local_vars = dict(variables)
        exec(f"result = {expression}", {}, local_vars)
        return str(local_vars.get("result", "error"))


class FormulaBuilder:
    """Constructs formula expressions from user input and rules.
    VULNERABLE: Formula injection."""

    def build_formula(self, formula: str, price: float, rules: list) -> str:
        """VULNERABLE: Code Injection (CWE-94)
        User-controlled formula with variable substitution."""
        import re
        result = re.sub(r'\bprice\b', str(price), formula)
        return result


class ResultFormatter:
    """Formats calculation results."""

    def format(self, result: str, precision: int) -> str:
        try:
            numeric = float(result)
            return str(round(numeric, precision))
        except (ValueError, TypeError):
            return result
