from flask import Blueprint, request, jsonify
from app.domain.service.content_service import ContentService
from app.domain.entity.content_submission import ContentSubmission

bp = Blueprint("message", __name__)
_content_service: ContentService = None


def init(content_service: ContentService):
    global _content_service
    _content_service = content_service


@bp.route("/api/messages")
def list_messages():
    messages = _content_service.get_all_messages()
    return jsonify([{"id": m.id, "content": m.content, "author": m.author, "created_at": m.created_at} for m in messages])


@bp.route("/api/messages", methods=["POST"])
def create_message():
    data = request.get_json(silent=True) or {}
    content = data.get("content", "")
    author = data.get("author", "anonymous")

    if not content:
        return jsonify({"error": "Content required"}), 400

    submission = ContentSubmission(content=content, author=author, context="message")
    message = _content_service.submit_message(submission)
    return jsonify({"id": message.id, "content": message.content, "author": message.author}), 201


@bp.route("/messages")
def messages_board():
    # VULNERABLE: Stored XSS (CWE-79)
    # Raw HTML content rendered without escaping
    messages = _content_service.get_all_messages()
    html = "<h2>Message Board</h2>"
    for m in messages:
        html += f"<div class='message'><strong>{m.author}</strong>: {m.content}<br><small>{m.created_at}</small></div><hr>"
    html += """
    <h3>Post a Message</h3>
    <form method="POST" action="/api/messages" enctype="application/x-www-form-urlencoded">
        <input name="author" placeholder="Your name"><br>
        <textarea name="content" placeholder="Message"></textarea><br>
        <button type="submit">Post</button>
    </form>
    """
    return html
