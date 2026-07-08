from flask import Blueprint, request, jsonify
from app.domain.service.system_command_service import SystemCommandService
from app.domain.service.external_request_service import ExternalRequestService
from app.domain.usecase.html_response_builder import HtmlResponseBuilder
from app.domain.entity.command_request import CommandRequest
from app.domain.entity.http_request_dto import HttpRequestDto

bp = Blueprint("network", __name__)
_command_service: SystemCommandService = None
_request_service: ExternalRequestService = None
_html_builder: HtmlResponseBuilder = None


def init(command_service: SystemCommandService, request_service: ExternalRequestService, html_builder: HtmlResponseBuilder):
    global _command_service, _request_service, _html_builder
    _command_service = command_service
    _request_service = request_service
    _html_builder = html_builder


@bp.route("/api/ping", methods=["POST"])
def ping():
    data = request.get_json(silent=True) or {}
    host = data.get("host", "")
    if not host:
        return jsonify({"error": "host parameter required"}), 400

    try:
        cmd_req = CommandRequest(target=host, command_type="ping")
        output = _command_service.execute_command(cmd_req)
        return jsonify({"host": host, "output": output})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/fetch")
def fetch_url():
    url = request.args.get("url", "")
    if not url:
        return jsonify({"error": "url parameter required"}), 400

    try:
        dto = HttpRequestDto(url=url, method="GET")
        result = _request_service.execute_request(dto)
        return jsonify(result)
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/webhook/test", methods=["POST"])
def test_webhook():
    data = request.get_json(silent=True) or {}
    url = data.get("url", "")
    payload = data.get("payload", {})

    if not url:
        return jsonify({"error": "url parameter required"}), 400

    try:
        result = _request_service.test_webhook(url, payload)
        return jsonify(result)
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/proxy")
def proxy():
    url = request.args.get("url", "")
    if not url:
        return jsonify({"error": "url parameter required"}), 400

    try:
        content = _request_service.proxy_request(url)
        return content, 200, {"Content-Type": "text/html"}
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/network/check", methods=["POST"])
def check_port():
    data = request.get_json(silent=True) or {}
    host = data.get("host", "")
    port = data.get("port", "")

    if not host or not port:
        return jsonify({"error": "host and port required"}), 400

    try:
        cmd_req = CommandRequest(target=f"{host} {port}", command_type="port_check")
        output = _command_service.execute_command(cmd_req)
        return jsonify({"host": host, "port": port, "output": output})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/network/ping")
def ping_page():
    host = request.args.get("host", "")
    if not host:
        return "<h2>Network Ping</h2><p>Provide ?host= parameter</p>"

    try:
        cmd_req = CommandRequest(target=host, command_type="ping")
        output = _command_service.execute_command(cmd_req)
        # VULNERABLE: Reflected XSS via HtmlResponseBuilder (CWE-79)
        return _html_builder.build_ping_page(host, output)
    except Exception as e:
        return _html_builder.build_ping_page(host, f"Error: {e}")
