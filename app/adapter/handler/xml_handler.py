from flask import Blueprint, request, jsonify
from app.domain.service.xml_processing_service import XmlProcessingService
from app.domain.entity.xml_processing_request import XmlProcessingRequest

bp = Blueprint("xml", __name__)
_xml_service: XmlProcessingService = None


def init(xml_service: XmlProcessingService):
    global _xml_service
    _xml_service = xml_service


@bp.route("/api/xml/parse", methods=["POST"])
def parse_xml():
    content_type = request.content_type or ""
    if "xml" in content_type:
        xml_content = request.data.decode("utf-8")
    else:
        data = request.get_json(silent=True) or {}
        xml_content = data.get("xml", "")

    if not xml_content:
        return jsonify({"error": "XML content required"}), 400

    try:
        req = XmlProcessingRequest(xml_content=xml_content, operation="parse", config_name="default")
        result = _xml_service.process_xml(req)
        return jsonify({"result": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/xml/validate", methods=["POST"])
def validate_xml():
    content_type = request.content_type or ""
    if "xml" in content_type:
        xml_content = request.data.decode("utf-8")
    else:
        data = request.get_json(silent=True) or {}
        xml_content = data.get("xml", "")

    if not xml_content:
        return jsonify({"error": "XML content required"}), 400

    try:
        req = XmlProcessingRequest(xml_content=xml_content, operation="validate", config_name="validate")
        result = _xml_service.validate_xml(req)
        return jsonify({"result": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/xml/transform", methods=["POST"])
def transform_xml():
    content_type = request.content_type or ""
    if "xml" in content_type:
        xml_content = request.data.decode("utf-8")
    else:
        data = request.get_json(silent=True) or {}
        xml_content = data.get("xml", "")

    if not xml_content:
        return jsonify({"error": "XML content required"}), 400

    try:
        req = XmlProcessingRequest(xml_content=xml_content, operation="transform", config_name="transform")
        result = _xml_service.transform_xml(req)
        return jsonify({"result": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500
