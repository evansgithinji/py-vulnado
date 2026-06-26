from app.domain.entity.calculation_rule import CalculationRequest, DiscountRequest
from app.domain.repository.rule_repository import RuleRepository
from app.domain.usecase.expression_evaluator import (
    ExpressionPreprocessor,
    ExpressionEvaluatorSink,
    FormulaBuilder,
    ResultFormatter,
)


class PricingEngine:
    """Orchestrates calculation and pricing flow.
    Call graph: Handler → Engine → RuleRepository → Preprocessor → Evaluator."""

    def __init__(
        self,
        rule_repo: RuleRepository,
        preprocessor: ExpressionPreprocessor,
        evaluator: ExpressionEvaluatorSink,
        formula_builder: FormulaBuilder,
        result_formatter: ResultFormatter,
    ):
        self._rule_repo = rule_repo
        self._preprocessor = preprocessor
        self._evaluator = evaluator
        self._formula_builder = formula_builder
        self._result_formatter = result_formatter

    def calculate(self, request: CalculationRequest) -> str:
        """Calculate expression.
        VULNERABLE: Code Injection flows through DTO → Preprocessor → Evaluator."""
        rules = self._rule_repo.get_rules(request.mode)

        processed_expr = self._preprocessor.preprocess(request.expression, rules)

        raw_result = self._evaluator.evaluate(processed_expr)

        return self._result_formatter.format(raw_result, request.precision)

    def calculate_discount(self, request: DiscountRequest) -> str:
        """Calculate discount.
        VULNERABLE: Code Injection flows through DTO → FormulaBuilder → Evaluator."""
        rules = self._rule_repo.get_rules("discount")

        formula = self._formula_builder.build_formula(
            request.formula, request.price, rules
        )

        raw_result = self._evaluator.evaluate_with_exec(formula, {"price": request.price})

        return self._result_formatter.format(raw_result, 2)

    def batch_calculate(self, expressions: list, mode: str = "standard") -> list:
        """Batch calculate expressions.
        VULNERABLE: Code Injection via eval() on each expression."""
        rules = self._rule_repo.get_rules(mode)
        results = []
        for expr in expressions:
            processed = self._preprocessor.preprocess(expr, rules)
            try:
                result = self._evaluator.evaluate(processed)
                results.append(result)
            except Exception as e:
                results.append(f"error: {e}")
        return results
