from flask import Blueprint, request, jsonify, Response
from app.domain.service.file_service import FileService
from app.domain.service.external_request_service import ExternalRequestService
from app.domain.service.system_command_service import SystemCommandService
from app.domain.entity.file_request import FileRequest

bp = Blueprint("file", __name__)
_file_service: FileService = None
_request_service: ExternalRequestService = None
_command_service: SystemCommandService = None
_files_dir: str = ""
_upload_dir: str = ""


def init(file_service: FileService, request_service: ExternalRequestService, command_service: SystemCommandService, files_dir: str, upload_dir: str):
    global _file_service, _request_service, _command_service, _files_dir, _upload_dir
    _file_service = file_service
    _request_service = request_service
    _command_service = command_service
    _files_dir = files_dir
    _upload_dir = upload_dir


@bp.route("/api/files")
def read_file():
    filename = request.args.get("filename", "")
    if not filename:
        return jsonify({"error": "filename parameter required"}), 400

    try:
        req = FileRequest(filename=filename, base_dir=_files_dir, operation="read")
        content = _file_service.process_file_request(req)
        return jsonify({"filename": filename, "content": content})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/files/download")
def download_file():
    filename = request.args.get("filename", "")
    if not filename:
        return jsonify({"error": "filename parameter required"}), 400

    try:
        req = FileRequest(filename=filename, base_dir=_files_dir, operation="download")
        content = _file_service.download_file(req)
        return Response(
            content,
            mimetype="application/octet-stream",
            headers={"Content-Disposition": f"attachment; filename={filename}"},
        )
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/files/upload", methods=["POST"])
def upload_file():
    if "file" not in request.files:
        return jsonify({"error": "No file provided"}), 400

    f = request.files["file"]
    if not f.filename:
        return jsonify({"error": "No filename"}), 400

    try:
        filepath = _file_service.upload_file(f.filename, f.read(), _upload_dir)
        return jsonify({"message": "File uploaded", "path": filepath}), 201
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/files/list")
def list_files():
    files = _file_service.list_files(_files_dir)
    return jsonify({"files": files})


@bp.route("/api/files/fetch")
def fetch_remote():
    url = request.args.get("url", "")
    if not url:
        return jsonify({"error": "url parameter required"}), 400

    try:
        content = _request_service.fetch_remote_file(url)
        return jsonify({"url": url, "content": content})
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/api/files/info")
def file_info():
    filename = request.args.get("filename", "")
    if not filename:
        return jsonify({"error": "filename parameter required"}), 400

    try:
        info = _command_service.execute_file_info(_files_dir, filename)
        return jsonify({"filename": filename, "info": info})
    except Exception as e:
        return jsonify({"error": str(e)}), 500
