from flask import Blueprint, request, jsonify
from app.domain.entity.calculation_rule import CalculationRequest, DiscountRequest
from app.domain.service.pricing_engine import PricingEngine

bp = Blueprint("calculator", __name__)
_engine: PricingEngine = None


def init(engine: PricingEngine):
    global _engine
    _engine = engine


@bp.route("/api/calculate")
def calculate():
    expr = request.args.get("expr", "")
    if not expr:
        return jsonify({"error": "expr parameter required"}), 400

    try:
        # Create DTO with tainted user input
        calc_request = CalculationRequest(
            expression=expr,
            precision=int(request.args.get("precision", "2")),
            mode=request.args.get("mode", "standard"),
        )

        # VULNERABLE: Code Injection (CWE-94) - flows through DTO → Engine → Preprocessor → Evaluator
        result = _engine.calculate(calc_request)
        return jsonify({"expression": expr, "result": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/calculate/discount", methods=["POST"])
def calculate_discount():
    data = request.get_json(silent=True) or {}
    price = data.get("price", 0)
    formula = data.get("formula", "")

    if not formula:
        return jsonify({"error": "formula parameter required"}), 400

    try:
        # Create DTO with tainted user input
        discount_request = DiscountRequest(
            price=float(price),
            formula=formula,
            customer_tier=data.get("tier", "standard"),
        )

        # VULNERABLE: Code Injection (CWE-94) - flows through DTO → Engine → FormulaBuilder → Evaluator
        result = _engine.calculate_discount(discount_request)
        return jsonify({"price": price, "formula": formula, "result": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/calculate/batch", methods=["POST"])
def batch_calculate():
    data = request.get_json(silent=True) or {}
    expressions = data.get("expressions", [])

    if not expressions:
        return jsonify({"error": "expressions required"}), 400

    try:
        # VULNERABLE: Code Injection (CWE-94) - batch eval on user expressions
        results = _engine.batch_calculate(expressions)
        return jsonify({"results": results})
    except Exception as e:
        return jsonify({"error": str(e)}), 500
