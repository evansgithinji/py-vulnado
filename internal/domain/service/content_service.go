package service

import (
	"time"

	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// ContentService orchestrates content submission with deep call graph
type ContentService struct {
	policyRepo  *repository.ContentPolicyRepository
	processor   *usecase.ContentProcessor
	messageRepo repository.MessageRepository
	reviewRepo  repository.ReviewRepository
}

// NewContentService creates a new ContentService
func NewContentService(
	policyRepo *repository.ContentPolicyRepository,
	processor *usecase.ContentProcessor,
	messageRepo repository.MessageRepository,
	reviewRepo repository.ReviewRepository,
) *ContentService {
	return &ContentService{
		policyRepo:  policyRepo,
		processor:   processor,
		messageRepo: messageRepo,
		reviewRepo:  reviewRepo,
	}
}

// SubmitMessage submits a message through the deep call graph
// Deep call graph for Stored XSS (depth 5):
// Handler -> ContentService.SubmitMessage()
//   -> ContentPolicyRepository.GetPolicy() -> ContentPolicy
//   -> ContentProcessor.ProcessContent() -> processed (VULNERABLE: no escaping)
//   -> MessageRepository.Create() -> Message (VULNERABLE: stores raw HTML)
func (s *ContentService) SubmitMessage(submission *entity.ContentSubmission) error {
	policy := s.policyRepo.GetPolicy("message")
	processed := s.processor.ProcessContent(submission.Content, policy)
	message := &entity.Message{
		Content:   processed,
		Author:    submission.Author,
		CreatedAt: time.Now(),
	}
	return s.messageRepo.Create(message)
}

// SubmitReview submits a review through the deep call graph
func (s *ContentService) SubmitReview(productID, userID int64, rating int, submission *entity.ContentSubmission) error {
	policy := s.policyRepo.GetPolicy("review")
	processed := s.processor.ProcessContent(submission.Content, policy)
	review := &entity.Review{
		ProductID: productID,
		UserID:    userID,
		Rating:    rating,
		Comment:   processed,
	}
	return s.reviewRepo.Create(review)
}

// GetAllMessages retrieves all messages
func (s *ContentService) GetAllMessages() ([]entity.Message, error) {
	return s.messageRepo.GetAll()
}

// GetProductReviews retrieves reviews for a product
func (s *ContentService) GetProductReviews(productID int64) ([]entity.Review, error) {
	return s.reviewRepo.FindByProductID(productID)
}
