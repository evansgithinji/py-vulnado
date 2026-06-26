from dataclasses import asdict
from flask import Blueprint, request, jsonify
from app.domain.usecase.order_usecase import OrderUseCase

bp = Blueprint("order", __name__)
_usecase: OrderUseCase = None


def init(usecase: OrderUseCase):
    global _usecase
    _usecase = usecase


@bp.route("/api/orders")
def list_orders():
    user_id = request.args.get("user_id", "")
    try:
        if user_id:
            orders = _usecase.get_user_orders(user_id)
        else:
            orders = _usecase.get_all_orders()
        return jsonify([asdict(o) for o in orders])
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/orders", methods=["POST"])
def create_order():
    data = request.get_json(silent=True) or {}
    user_id = data.get("user_id", 0)
    product_id = data.get("product_id", 0)
    quantity = data.get("quantity", 1)
    total = data.get("total", 0.0)

    if not user_id or not product_id:
        return jsonify({"error": "user_id and product_id required"}), 400

    try:
        order = _usecase.create_order(user_id, product_id, quantity, total)
        return jsonify(asdict(order)), 201
    except Exception as e:
        return jsonify({"error": str(e)}), 500
