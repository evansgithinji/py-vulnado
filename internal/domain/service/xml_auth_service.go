package service

import (
	"fmt"

	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// XmlAuthService orchestrates XPath-based authentication and queries.
type XmlAuthService struct {
	repo              repository.XmlDocumentRepository
	expressionBuilder *usecase.XPathExpressionBuilder
	evaluator         *usecase.XPathEvaluator
}

// NewXmlAuthService creates a new XmlAuthService.
func NewXmlAuthService(
	repo repository.XmlDocumentRepository,
	expressionBuilder *usecase.XPathExpressionBuilder,
	evaluator *usecase.XPathEvaluator,
) *XmlAuthService {
	return &XmlAuthService{
		repo:              repo,
		expressionBuilder: expressionBuilder,
		evaluator:         evaluator,
	}
}

// Authenticate performs XPath-based login.
// Call graph: XmlAuthService.Authenticate → repo.FindBySourceID → expressionBuilder.BuildAuthQuery → evaluator.Evaluate
func (s *XmlAuthService) Authenticate(req *entity.XPathAuthRequest) (string, error) {
	// Layer 3: Repository lookup for XML document
	doc, err := s.repo.FindBySourceID("users")
	if err != nil {
		return "", fmt.Errorf("failed to load user document: %v", err)
	}
	if doc == nil {
		return "", fmt.Errorf("user document not found")
	}

	// Layer 4: Build vulnerable XPath expression
	xpathExpr := s.expressionBuilder.BuildAuthQuery(req.Username, req.Password)

	// Layer 5: Evaluate XPath against document
	return s.evaluator.Evaluate(doc, xpathExpr)
}

// Query performs a raw XPath query.
// Call graph: XmlAuthService.Query → repo.FindBySourceID → expressionBuilder.BuildRawQuery → evaluator.Evaluate
func (s *XmlAuthService) Query(query string) (string, error) {
	// Layer 3: Repository lookup for XML document
	doc, err := s.repo.FindBySourceID("users")
	if err != nil {
		return "", fmt.Errorf("failed to load user document: %v", err)
	}
	if doc == nil {
		return "", fmt.Errorf("user document not found")
	}

	// Layer 4: Build raw query expression (vulnerable - direct user input)
	xpathExpr := s.expressionBuilder.BuildRawQuery(query)

	// Layer 5: Evaluate XPath against document
	return s.evaluator.Evaluate(doc, xpathExpr)
}
