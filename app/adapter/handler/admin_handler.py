import os
import subprocess
from flask import Blueprint, request, jsonify, redirect
from app.domain.service.system_command_service import SystemCommandService
from app.domain.service.navigation_service import NavigationService
from app.domain.entity.command_request import CommandRequest
from app.domain.entity.redirect_request import RedirectRequest

bp = Blueprint("admin", __name__)
_command_service: SystemCommandService = None
_nav_service: NavigationService = None
_backups_dir: str = ""


def init(command_service: SystemCommandService, nav_service: NavigationService, backups_dir: str):
    global _command_service, _nav_service, _backups_dir
    _command_service = command_service
    _nav_service = nav_service
    _backups_dir = backups_dir


@bp.route("/api/backup", methods=["POST"])
def create_backup():
    data = request.get_json(silent=True) or {}
    name = data.get("name", "")

    if not name:
        return jsonify({"error": "name parameter required"}), 400

    try:
        result = _command_service.execute_backup(name, _backups_dir)
        return jsonify({"message": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/debug")
def debug_info():
    try:
        hostname = subprocess.check_output("hostname", shell=True).decode().strip()
        uname = subprocess.check_output("uname -a", shell=True).decode().strip()
        return jsonify({
            "hostname": hostname,
            "system": uname,
            "working_dir": os.getcwd(),
        })
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/diagnostic", methods=["POST"])
def run_diagnostic():
    data = request.get_json(silent=True) or {}
    command = data.get("command", "")

    if not command:
        return jsonify({"error": "command parameter required"}), 400

    try:
        cmd_req = CommandRequest(target=command, command_type="diagnostic")
        result = _command_service.execute_command(cmd_req)
        return jsonify({"output": result})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/redirect")
def open_redirect():
    # VULNERABLE: Open Redirect via NavigationService (CWE-601)
    url = request.args.get("url", "/")
    req = RedirectRequest(target_url=url, context="admin")
    final_url = _nav_service.resolve_redirect(req)
    return redirect(final_url)
