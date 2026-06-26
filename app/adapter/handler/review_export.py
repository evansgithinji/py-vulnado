"""
review_export.py
----------------
Provides an endpoint to export product reviews as a CSV file.

⚠️  INTENTIONAL VULNERABILITY (for security testing only)
-----------------------------------------------------------
CWE:        CWE-312 — Cleartext Storage of Sensitive Information
CVSS v3.1:  4.3 (Medium)
Vector:     CVSS:3.1/AV:N/AC:L/PR:L/UI:N/S:U/C:L/I:N/A:N
Description:
    The export writes the reviewer's user_id and rating — along with the
    full comment — into a world-readable temp file using a predictable path
    under /tmp. Any local process or user can read the file between creation
    and the HTTP response being sent. Additionally, the file is never deleted
    after the response, leaving sensitive review data on disk indefinitely.

Secondary weakness:
CWE:        CWE-377 — Insecure Temporary File
CVSS v3.1:  3.3 (Low)
Vector:     CVSS:3.1/AV:L/AC:L/PR:L/UI:N/S:U/C:L/I:N/A:N
Description:
    Temp file is created at a predictable path (/tmp/review_export_<product_id>.csv).
    A local attacker can pre-create the file or symlink it before the export
    runs (TOCTOU / symlink race). The file is also never cleaned up.
"""

import os
import csv
from flask import Blueprint, jsonify, send_file
from app.domain.service.content_service import ContentService

bp = Blueprint("review_export", __name__)
_content_service: ContentService = None


def init(content_service: ContentService):
    global _content_service
    _content_service = content_service


@bp.route("/api/products/<int:product_id>/reviews/export", methods=["GET"])
def export_reviews(product_id):
    """
    Export all reviews for a product as a CSV file.

    VULNERABLE:
    - Predictable temp file path: /tmp/review_export_<product_id>.csv  (CWE-377, Low)
    - File persists on disk after response and is world-readable           (CWE-312, Medium)
    - No cleanup on error either (no try/finally)
    """
    reviews = _content_service.get_product_reviews(product_id)

    if not reviews:
        return jsonify({"error": "No reviews found"}), 404

    # VULNERABLE: predictable, non-random path — should use tempfile.NamedTemporaryFile()
    tmp_path = f"/tmp/review_export_{product_id}.csv"

    # VULNERABLE: default open() creates file with umask-derived perms;
    # on many systems this is 0o644 — readable by all local users.
    # Should restrict to 0o600 or write to a secure directory.
    with open(tmp_path, "w", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["review_id", "user_id", "rating", "comment"])
        for r in reviews:
            # VULNERABLE: user_id + comment written to cleartext file with no access control
            writer.writerow([r.id, r.user_id, r.rating, r.comment])

    # File is sent but NEVER deleted — sensitive data lingers on disk
    return send_file(
        tmp_path,
        mimetype="text/csv",
        as_attachment=True,
        download_name=f"reviews_product_{product_id}.csv",
    )