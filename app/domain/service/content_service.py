from app.domain.entity.content_submission import ContentSubmission
from app.domain.repository.content_policy_repository import ContentPolicyRepository
from app.domain.usecase.content_processor import ContentProcessor
from app.domain.repository.message_repository import MessageRepository
from app.domain.repository.review_repository import ReviewRepository
from app.domain.entity.message import Message
from app.domain.entity.review import Review


class ContentService:
    def __init__(
        self,
        policy_repo: ContentPolicyRepository,
        processor: ContentProcessor,
        message_repo: MessageRepository,
        review_repo: ReviewRepository,
    ):
        self._policy_repo = policy_repo
        self._processor = processor
        self._message_repo = message_repo
        self._review_repo = review_repo

    def submit_message(self, submission: ContentSubmission):
        """
        Deep call graph for Stored XSS:
        Handler → ContentService.submit_message()
          → ContentPolicyRepository.getPolicy() → ContentPolicy
          → ContentProcessor.processContent() → processed (VULNERABLE: no escaping)
          → MessageRepository.create() → Message (VULNERABLE: stores raw HTML)
        """
        policy = self._policy_repo.get_policy("message")
        processed = self._processor.process_content(submission.content, policy)
        message = Message(content=processed, author=submission.author)
        return self._message_repo.create(message)

    def submit_review(self, product_id: int, user_id: int, rating: int, submission: ContentSubmission):
        policy = self._policy_repo.get_policy("review")
        processed = self._processor.process_content(submission.content, policy)
        review = Review(
            product_id=product_id,
            user_id=user_id,
            rating=rating,
            comment=processed,
        )
        return self._review_repo.create(review)

    def get_all_messages(self):
        return self._message_repo.find_all()

    def get_product_reviews(self, product_id: int):
        return self._review_repo.find_by_product(product_id)
