package service

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// PricingEngine orchestrates calculation and expression evaluation operations.
type PricingEngine struct {
	ruleRepo       repository.RuleRepository
	preprocessor   *usecase.ExpressionPreprocessor
	formulaBuilder *usecase.FormulaBuilder
	evaluatorSink  *usecase.ExpressionEvaluatorSink
	formatter      *usecase.ResultFormatter
}

// NewPricingEngine creates a new PricingEngine.
func NewPricingEngine(
	ruleRepo repository.RuleRepository,
	preprocessor *usecase.ExpressionPreprocessor,
	formulaBuilder *usecase.FormulaBuilder,
	evaluatorSink *usecase.ExpressionEvaluatorSink,
	formatter *usecase.ResultFormatter,
) *PricingEngine {
	return &PricingEngine{
		ruleRepo:       ruleRepo,
		preprocessor:   preprocessor,
		formulaBuilder: formulaBuilder,
		evaluatorSink:  evaluatorSink,
		formatter:      formatter,
	}
}

// Evaluate evaluates a mathematical expression.
// Call graph: PricingEngine.Evaluate → ruleRepo.FindByName → preprocessor.Preprocess → formulaBuilder.BuildShellExpression → evaluatorSink.Execute → formatter.Format
func (s *PricingEngine) Evaluate(req *entity.CalculationRequest) (string, error) {
	// Layer 3: Repository lookup for rule
	rule, err := s.ruleRepo.FindByName(req.RuleName)
	if err != nil || rule == nil {
		rule = &entity.CalculationRule{Name: "default", MaxLength: 1024, AllowShell: true}
	}

	// Layer 4: Preprocess expression
	processed := s.preprocessor.Preprocess(rule, req.Expression)

	// Layer 4: Build shell expression (VULNERABLE)
	shellExpr := s.formulaBuilder.BuildShellExpression(processed)

	// Layer 5: Execute in shell (VULNERABLE)
	rawResult, err := s.evaluatorSink.Execute(shellExpr)
	if err != nil {
		return "", err
	}

	// Format result
	return s.formatter.Format(rawResult), nil
}

// CalculateDiscount calculates a discount using a user-provided formula.
// Call graph: PricingEngine.CalculateDiscount → ruleRepo.FindByName → preprocessor.PreprocessDiscount → formulaBuilder.BuildShellExpression → evaluatorSink.Execute → formatter.Format
func (s *PricingEngine) CalculateDiscount(req *entity.DiscountRequest) (string, error) {
	// Layer 3: Repository lookup for rule
	rule, err := s.ruleRepo.FindByName(req.RuleName)
	if err != nil || rule == nil {
		rule = &entity.CalculationRule{Name: "discount", MaxLength: 2048, AllowShell: true}
	}

	// Layer 4: Preprocess formula with price substitution
	processed := s.preprocessor.PreprocessDiscount(rule, req.Formula, req.Price)

	// Layer 4: Build shell expression (VULNERABLE)
	shellExpr := s.formulaBuilder.BuildShellExpression(processed)

	// Layer 5: Execute in shell (VULNERABLE)
	rawResult, err := s.evaluatorSink.Execute(shellExpr)
	if err != nil {
		return "", err
	}

	// Format result
	return s.formatter.Format(rawResult), nil
}
