from flask import Blueprint, request, jsonify
from app.domain.entity.ldap_config import LdapAuthRequest
from app.domain.service.directory_service import DirectoryService

bp = Blueprint("ldap", __name__)
_service: DirectoryService = None


def init(service: DirectoryService):
    global _service
    _service = service


def _safe_user(u: dict) -> dict:
    return {"uid": u["uid"], "cn": u["cn"], "mail": u["mail"], "ou": u["ou"]}


@bp.route("/api/ldap/login", methods=["POST"])
def login():
    data = request.get_json(silent=True) or {}
    username = data.get("username", "")
    password = data.get("password", "")

    if not username or not password:
        return jsonify({"error": "username and password required"}), 400

    # Create DTO with tainted user input
    auth_request = LdapAuthRequest(
        username=username,
        password=password,
        domain=data.get("domain", "default"),
    )

    # VULNERABLE: LDAP Injection (CWE-90) - unsanitized user input flows through DTO
    results = _service.authenticate(auth_request)

    if results:
        return jsonify({
            "success": True,
            "message": "Login successful",
            "users": [_safe_user(u) for u in results],
        })
    else:
        return jsonify({"success": False, "message": "Invalid credentials"}), 401


@bp.route("/api/ldap/search")
def search():
    q = request.args.get("q", "")
    if not q:
        return jsonify({"error": "q parameter required"}), 400

    # VULNERABLE: LDAP Injection (CWE-90) - unsanitized user input
    results = _service.search(q)

    return jsonify({
        "success": True,
        "count": len(results),
        "results": [_safe_user(u) for u in results],
    })
