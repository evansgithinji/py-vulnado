import sqlite3
from app.domain.entity.review import Review
from app.domain.repository.review_repository import ReviewRepository


class SQLiteReviewRepository(ReviewRepository):
    def __init__(self, conn: sqlite3.Connection):
        self.conn = conn

    def find_by_product(self, product_id: int) -> list[Review]:
        cursor = self.conn.execute(
            "SELECT * FROM reviews WHERE product_id = ?", (product_id,)
        )
        reviews = []
        for row in cursor.fetchall():
            reviews.append(
                Review(
                    id=row["id"],
                    product_id=row["product_id"],
                    user_id=row["user_id"],
                    rating=row["rating"],
                    comment=row["comment"],
                )
            )
        return reviews

    def create(self, review: Review) -> Review:
        cursor = self.conn.execute(
            "INSERT INTO reviews (product_id, user_id, rating, comment) VALUES (?, ?, ?, ?)",
            (review.product_id, review.user_id, review.rating, review.comment),
        )
        self.conn.commit()
        review.id = cursor.lastrowid
        return review
