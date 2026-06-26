from app.domain.repository.review_repository import ReviewRepository
from app.domain.entity.review import Review


class ReviewUseCase:
    def __init__(self, review_repo: ReviewRepository):
        self.review_repo = review_repo

    def get_product_reviews(self, product_id: int):
        return self.review_repo.find_by_product(product_id)

    def create_review(self, product_id: int, user_id: int, rating: int, comment: str):
        review = Review(
            product_id=product_id,
            user_id=user_id,
            rating=rating,
            comment=comment,
        )
        return self.review_repo.create(review)
