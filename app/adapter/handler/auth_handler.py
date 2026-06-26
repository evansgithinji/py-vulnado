from dataclasses import asdict
from flask import Blueprint, request, jsonify, redirect
from app.domain.service.authentication_service import AuthenticationService
from app.domain.service.navigation_service import NavigationService
from app.domain.usecase.html_response_builder import HtmlResponseBuilder
from app.domain.entity.auth_request import AuthRequest
from app.domain.entity.redirect_request import RedirectRequest
from app.domain.usecase.auth_usecase import AuthUseCase

bp = Blueprint("auth", __name__)
_auth_service: AuthenticationService = None
_nav_service: NavigationService = None
_html_builder: HtmlResponseBuilder = None
_auth_uc: AuthUseCase = None


def init(auth_service: AuthenticationService, nav_service: NavigationService, html_builder: HtmlResponseBuilder, auth_uc: AuthUseCase):
    global _auth_service, _nav_service, _html_builder, _auth_uc
    _auth_service = auth_service
    _nav_service = nav_service
    _html_builder = html_builder
    _auth_uc = auth_uc


@bp.route("/api/login", methods=["POST"])
def login():
    data = request.get_json(silent=True) or {}
    username = data.get("username", "")
    password = data.get("password", "")

    if not username or not password:
        return jsonify({"error": "Username and password required"}), 400

    try:
        auth_req = AuthRequest(username=username, password=password)
        result = _auth_service.authenticate(auth_req)
        if result:
            return jsonify({"message": "Login successful", "user": result})
        return jsonify({"error": "Invalid credentials"}), 401
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/users/search")
def search_users():
    q = request.args.get("q", "")
    if not q:
        return "<h2>User Search</h2><p>Provide ?q= parameter</p>"

    try:
        users = _auth_uc.search_users(q)
        # VULNERABLE: Reflected XSS via HtmlResponseBuilder (CWE-79)
        return _html_builder.build_user_search_page(q, users)
    except Exception as e:
        return f"<h2>Search results for: {q}</h2><p>Error: {e}</p>"


@bp.route("/api/users")
def list_users():
    users = _auth_uc.get_all_users()
    return jsonify([asdict(u) for u in users])


@bp.route("/login/callback")
def login_callback():
    # VULNERABLE: Open Redirect via NavigationService (CWE-601)
    redirect_url = request.args.get("redirect", "/")
    req = RedirectRequest(target_url=redirect_url, context="auth")
    final_url = _nav_service.resolve_redirect(req)
    return redirect(final_url)
