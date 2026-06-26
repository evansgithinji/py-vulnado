from flask import Blueprint, request, jsonify
from app.domain.service.serialization_service import SerializationService
from app.domain.entity.deserialization_request import DeserializationRequest

bp = Blueprint("deserialize", __name__)
_serialization_service: SerializationService = None


def init(serialization_service: SerializationService):
    global _serialization_service
    _serialization_service = serialization_service


@bp.route("/api/cart/export", methods=["POST"])
def export_cart():
    data = request.get_json(silent=True) or {}
    items = data.get("items", [])

    if not items:
        return jsonify({"error": "items required"}), 400

    try:
        encoded = _serialization_service.serialize(items, "pickle")
        return jsonify({"format": "pickle", "data": encoded})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/cart/import", methods=["POST"])
def import_cart():
    data = request.get_json(silent=True) or {}
    encoded = data.get("data", "")
    fmt = data.get("format", "pickle")

    if not encoded:
        return jsonify({"error": "data required"}), 400

    try:
        req = DeserializationRequest(data=encoded, format=fmt)
        items = _serialization_service.deserialize(req)
        return jsonify({"items": items})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/cart/export/json", methods=["POST"])
def export_cart_json():
    data = request.get_json(silent=True) or {}
    items = data.get("items", [])

    if not items:
        return jsonify({"error": "items required"}), 400

    try:
        encoded = _serialization_service.serialize(items, "json")
        return jsonify({"format": "json", "data": encoded})
    except Exception as e:
        return jsonify({"error": str(e)}), 500
