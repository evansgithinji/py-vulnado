from flask import Blueprint, request, jsonify
from app.domain.entity.audit_event import AuditEvent
from app.domain.service.audit_service import AuditService

bp = Blueprint("log", __name__)
_service: AuditService = None


def init(service: AuditService):
    global _service
    _service = service


@bp.route("/api/log/login", methods=["POST"])
def log_login():
    data = request.get_json(silent=True) or {}
    username = data.get("username", "")
    password = data.get("password", "")

    if not username:
        return jsonify({"error": "username required"}), 400

    status = "success" if username == "admin" and password == "admin123" else "failed"

    # Create DTO with tainted user input
    event = AuditEvent(
        event_type="login",
        actor=username,
        details=status,
        source_ip=request.remote_addr or "",
    )

    # VULNERABLE: Log Injection (CWE-117) - flows through DTO → Service → Enricher → Formatter → Writer → Storage
    msg = _service.record_event(event)
    return jsonify({"success": True, "logged": msg})


@bp.route("/api/log/search", methods=["POST"])
def log_search():
    data = request.get_json(silent=True) or {}
    query = data.get("query", "")

    if not query:
        return jsonify({"error": "query required"}), 400

    # Create DTO with tainted user input
    event = AuditEvent(
        event_type="search",
        actor="anonymous",
        details=query,
        source_ip=request.remote_addr or "",
    )

    # VULNERABLE: Log Injection (CWE-117) - flows through DTO
    msg = _service.record_event(event)
    return jsonify({"success": True, "logged": msg})


@bp.route("/api/logs")
def get_logs():
    logs = _service.get_logs()
    return jsonify({"success": True, "count": len(logs), "logs": logs})
