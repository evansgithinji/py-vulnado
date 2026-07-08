package usecase

import (
	"fmt"
	"os/exec"
	"strings"

	"goapp/internal/domain/entity"
)

// ExpressionPreprocessor preprocesses expressions before evaluation.
type ExpressionPreprocessor struct{}

// NewExpressionPreprocessor creates a new ExpressionPreprocessor.
func NewExpressionPreprocessor() *ExpressionPreprocessor {
	return &ExpressionPreprocessor{}
}

// Preprocess prepares an expression for evaluation according to the rule.
func (p *ExpressionPreprocessor) Preprocess(rule *entity.CalculationRule, expression string) string {
	// Apply rule constraints (but AllowShell is always true in vulnerable mode)
	if rule.MaxLength > 0 && len(expression) > rule.MaxLength {
		expression = expression[:rule.MaxLength]
	}
	return expression
}

// PreprocessDiscount substitutes the price variable in a formula.
func (p *ExpressionPreprocessor) PreprocessDiscount(rule *entity.CalculationRule, formula string, price float64) string {
	expr := strings.ReplaceAll(formula, "price", fmt.Sprintf("%f", price))
	return p.Preprocess(rule, expr)
}

// FormulaBuilder constructs shell commands from expressions.
type FormulaBuilder struct{}

// NewFormulaBuilder creates a new FormulaBuilder.
func NewFormulaBuilder() *FormulaBuilder {
	return &FormulaBuilder{}
}

// BuildShellExpression builds a shell expression string for bc evaluation.
// VULNERABLE: Code Injection (CWE-94) - shell command injection via expression
func (b *FormulaBuilder) BuildShellExpression(expression string) string {
	// VULNERABLE: User input passed to shell via string concatenation
	return "echo '" + expression + "' | bc"
}

// ExpressionEvaluatorSink executes shell expressions.
type ExpressionEvaluatorSink struct{}

// NewExpressionEvaluatorSink creates a new ExpressionEvaluatorSink.
func NewExpressionEvaluatorSink() *ExpressionEvaluatorSink {
	return &ExpressionEvaluatorSink{}
}

// Execute runs a shell expression and returns the output.
// VULNERABLE: Code Injection (CWE-94) - shell command injection
func (e *ExpressionEvaluatorSink) Execute(shellExpr string) (string, error) {
	cmd := exec.Command("sh", "-c", shellExpr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("evaluation failed: %s", strings.TrimSpace(string(output)))
	}
	return strings.TrimSpace(string(output)), nil
}

// ResultFormatter formats calculation results.
type ResultFormatter struct{}

// NewResultFormatter creates a new ResultFormatter.
func NewResultFormatter() *ResultFormatter {
	return &ResultFormatter{}
}

// Format formats a raw calculation result string.
func (f *ResultFormatter) Format(rawResult string) string {
	return strings.TrimSpace(rawResult)
}
