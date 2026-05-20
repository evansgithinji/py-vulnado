from flask import Blueprint, request, jsonify, redirect, make_response
from app.domain.entity.header_policy import HeaderRequest, LocaleRequest
from app.domain.service.response_customization_service import (
    ResponseCustomizationService,
    LocalizationService,
)

bp = Blueprint("header", __name__)
_customization_service: ResponseCustomizationService = None
_localization_service: LocalizationService = None


def init(customization_service: ResponseCustomizationService, localization_service: LocalizationService):
    global _customization_service, _localization_service
    _customization_service = customization_service
    _localization_service = localization_service


@bp.route("/api/header/set")
def set_header():
    name = request.args.get("name", "")
    value = request.args.get("value", "")

    if not name or not value:
        return jsonify({"error": "name and value parameters required"}), 400

    # Create DTO with tainted user input
    header_request = HeaderRequest(name=name, value=value, context=request.url)

    # VULNERABLE: HTTP Header Injection (CWE-113) - flows through DTO → Service → Processor → Writer
    header_name, header_value = _customization_service.apply_custom_header(header_request)

    resp = make_response(jsonify({"success": True, "header": header_name, "value": header_value}))
    resp.headers[header_name] = header_value
    return resp


@bp.route("/api/header/redirect")
def redirect_with_lang():
    lang = request.args.get("lang", "")

    if not lang:
        return jsonify({"error": "lang parameter required"}), 400

    # Create DTO with tainted user input
    locale_request = LocaleRequest(lang=lang, return_url="/")

    # VULNERABLE: HTTP Header Injection (CWE-113) - flows through DTO → Service → Repository → Cookie/Redirect
    result = _localization_service.build_localized_redirect(locale_request)

    resp = make_response(redirect(result["redirect_url"], code=302))
    resp.headers["X-Language"] = lang
    resp.set_cookie(result["cookie_name"], result["cookie_value"])
    return resp
