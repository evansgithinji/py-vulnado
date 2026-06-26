from dataclasses import asdict
from flask import Blueprint, request, jsonify
from app.domain.service.catalog_service import CatalogService
from app.domain.usecase.html_response_builder import HtmlResponseBuilder
from app.domain.entity.search_request import SearchRequest

bp = Blueprint("product", __name__)
_catalog_service: CatalogService = None
_html_builder: HtmlResponseBuilder = None


def init(catalog_service: CatalogService, html_builder: HtmlResponseBuilder):
    global _catalog_service, _html_builder
    _catalog_service = catalog_service
    _html_builder = html_builder


@bp.route("/api/products")
def list_products():
    req = SearchRequest(query="")
    products = _catalog_service.search(req)
    return jsonify([asdict(p) for p in products])


@bp.route("/api/products/<int:product_id>")
def get_product(product_id):
    req = SearchRequest(query=str(product_id))
    products = _catalog_service.search(req)
    for p in products:
        if p.id == product_id:
            return jsonify(asdict(p))
    return jsonify({"error": "Product not found"}), 404


@bp.route("/api/products/search")
def search_products():
    q = request.args.get("q", "")
    category = request.args.get("category", "")

    try:
        if category:
            products = _catalog_service.search_by_category(category)
        elif q:
            req = SearchRequest(query=q)
            products = _catalog_service.search(req)
        else:
            req = SearchRequest(query="")
            products = _catalog_service.search(req)
        return jsonify([asdict(p) for p in products])
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/products/search")
def search_products_html():
    q = request.args.get("q", "")
    if not q:
        return "<h2>Product Search</h2><p>Provide ?q= parameter</p>"

    try:
        req = SearchRequest(query=q)
        products = _catalog_service.search(req)
        # VULNERABLE: Reflected XSS via HtmlResponseBuilder (CWE-79)
        return _html_builder.build_search_page(q, products)
    except Exception as e:
        return _html_builder.build_error_page(q, str(e))
