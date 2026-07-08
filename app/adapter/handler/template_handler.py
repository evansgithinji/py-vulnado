from flask import Blueprint, request, jsonify
from app.domain.service.notification_service import NotificationService
from app.domain.entity.render_request import RenderRequest

bp = Blueprint("template", __name__)
_notification_service: NotificationService = None


def init(notification_service: NotificationService):
    global _notification_service
    _notification_service = notification_service


@bp.route("/api/greeting")
def greeting():
    name = request.args.get("name", "World")

    try:
        req = RenderRequest(user_input=name, template_name="greeting")
        result = _notification_service.render_notification(req)
        return jsonify({"greeting": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/template/render", methods=["POST"])
def render_template():
    data = request.get_json(silent=True) or {}
    template_str = data.get("template", "")

    if not template_str:
        return jsonify({"error": "template parameter required"}), 400

    try:
        req = RenderRequest(user_input=template_str, template_name="custom")
        result = _notification_service.render_notification(req)
        return jsonify({"result": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/invoice/preview", methods=["POST"])
def invoice_preview():
    data = request.get_json(silent=True) or {}
    template_str = data.get("template", "Invoice for {{ customer }}: ${{ amount }}")
    customer = data.get("customer", "")
    amount = data.get("amount", "0")

    try:
        result = _notification_service.render_invoice(
            template_str, {"customer": customer, "amount": amount}
        )
        return jsonify({"preview": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500
