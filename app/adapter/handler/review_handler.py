from flask import Blueprint, request, jsonify
from app.domain.service.content_service import ContentService
from app.domain.entity.content_submission import ContentSubmission

bp = Blueprint("review", __name__)
_content_service: ContentService = None


def init(content_service: ContentService):
    global _content_service
    _content_service = content_service


@bp.route("/api/products/<int:product_id>/reviews")
def get_reviews(product_id):
    reviews = _content_service.get_product_reviews(product_id)
    return jsonify([{"id": r.id, "product_id": r.product_id, "user_id": r.user_id, "rating": r.rating, "comment": r.comment} for r in reviews])


@bp.route("/api/products/<int:product_id>/reviews", methods=["POST"])
def create_review(product_id):
    data = request.get_json(silent=True) or {}
    user_id = data.get("user_id", 0)
    rating = data.get("rating", 0)
    comment = data.get("comment", "")

    if not user_id or not rating:
        return jsonify({"error": "user_id and rating required"}), 400

    submission = ContentSubmission(content=comment, context="review")
    review = _content_service.submit_review(product_id, user_id, rating, submission)
    return jsonify({"id": review.id, "product_id": review.product_id, "rating": review.rating, "comment": review.comment}), 201


@bp.route("/products/<int:product_id>/reviews")
def reviews_page(product_id):
    # VULNERABLE: Stored XSS (CWE-79)
    # Review comments rendered as raw HTML
    reviews = _content_service.get_product_reviews(product_id)
    html = f"<h2>Reviews for Product #{product_id}</h2>"
    for r in reviews:
        html += f"<div class='review'><strong>User #{r.user_id}</strong> - Rating: {r.rating}/5<p>{r.comment}</p></div><hr>"
    return html
