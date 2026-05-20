from flask import Blueprint, request, jsonify
from app.domain.entity.document_collection import NoSqlAuthRequest, WhereQueryRequest
from app.domain.service.profile_service import ProfileService

bp = Blueprint("nosql", __name__)
_service: ProfileService = None


def init(service: ProfileService):
    global _service
    _service = service


def _safe_user(u: dict) -> dict:
    return {"_id": u["_id"], "username": u["username"], "role": u["role"], "email": u["email"]}


@bp.route("/api/nosql/login", methods=["POST"])
def login():
    data = request.get_json(silent=True) or {}
    username = data.get("username")
    password = data.get("password")

    if username is None or password is None:
        return jsonify({"error": "username and password required"}), 400

    # Create DTO with tainted user input (username/password may be objects with $ne/$gt/$regex)
    auth_request = NoSqlAuthRequest(
        username=username,
        password=password,
        collection=data.get("collection", "users"),
    )

    # VULNERABLE: NoSQL Injection (CWE-943) - operator injection flows through DTO
    results = _service.authenticate_user(auth_request)

    if results:
        return jsonify({
            "success": True,
            "message": "Login successful",
            "users": [_safe_user(u) for u in results],
        })
    else:
        return jsonify({"success": False, "message": "Invalid credentials"}), 401


@bp.route("/api/nosql/query", methods=["POST"])
def query():
    data = request.get_json(silent=True) or {}
    query_map = data.get("query", {})

    # VULNERABLE: NoSQL Injection (CWE-943) - operator objects in query map
    results = _service.query(query_map)

    return jsonify({
        "success": True,
        "count": len(results),
        "results": [_safe_user(u) for u in results],
    })


@bp.route("/api/nosql/users")
def find_where():
    where = request.args.get("where", "")
    if not where:
        return jsonify({"error": "where parameter required"}), 400

    # Create DTO with tainted user input
    where_request = WhereQueryRequest(where_expression=where, collection="users")

    # VULNERABLE: Code injection via $where clause flows through DTO
    results = _service.find_by_expression(where_request)

    return jsonify({
        "success": True,
        "count": len(results),
        "results": [_safe_user(u) for u in results],
    })
