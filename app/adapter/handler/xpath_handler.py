from flask import Blueprint, request, jsonify
from app.domain.entity.xml_document import XPathAuthRequest
from app.domain.service.xml_auth_service import XmlAuthService

bp = Blueprint("xpath", __name__)
_service: XmlAuthService = None


def init(service: XmlAuthService):
    global _service
    _service = service


@bp.route("/api/xpath/login", methods=["POST"])
def xpath_login():
    data = request.get_json(silent=True) or {}
    username = data.get("username", "")
    password = data.get("password", "")

    if not username or not password:
        return jsonify({"error": "username and password required"}), 400

    try:
        # Create DTO with tainted user input
        auth_request = XPathAuthRequest(
            username=username,
            password=password,
            document_source=data.get("source", "default"),
        )

        # VULNERABLE: XPath Injection (CWE-643) - flows through DTO
        results = _service.authenticate_via_xpath(auth_request)

        if results:
            return jsonify({"success": True, "message": "Login successful", "name": results})
        else:
            return jsonify({"success": False, "message": "Invalid credentials"}), 401
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/xpath/query")
def xpath_query():
    q = request.args.get("q", "")
    if not q:
        return jsonify({"error": "q parameter required"}), 400

    try:
        # VULNERABLE: XPath Injection (CWE-643) - user-controlled XPath expression
        results = _service.query(q)
        return jsonify({"success": True, "result": results})
    except Exception as e:
        return jsonify({"error": str(e)}), 500
