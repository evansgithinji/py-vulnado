package usecase

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
)

// ReviewUseCase handles review operations
type ReviewUseCase struct {
	reviewRepo repository.ReviewRepository
}

// NewReviewUseCase creates a new ReviewUseCase
func NewReviewUseCase(reviewRepo repository.ReviewRepository) *ReviewUseCase {
	return &ReviewUseCase{reviewRepo: reviewRepo}
}

// GetProductReviews retrieves reviews for a product
func (uc *ReviewUseCase) GetProductReviews(productID int64) ([]entity.Review, error) {
	return uc.reviewRepo.FindByProductID(productID)
}

// CreateReview creates a new review
// VULNERABLE: No sanitization of review content (SQLi, stored XSS)
func (uc *ReviewUseCase) CreateReview(review *entity.Review) error {
	return uc.reviewRepo.Create(review)
}
