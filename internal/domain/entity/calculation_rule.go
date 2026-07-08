package entity

// CalculationRule defines a rule for expression evaluation.
type CalculationRule struct {
	Name       string
	MaxLength  int
	AllowShell bool
}

// CalculationRequest is a DTO for calculation requests.
type CalculationRequest struct {
	Expression string
	RuleName   string
}

// DiscountRequest is a DTO for discount calculation requests.
type DiscountRequest struct {
	Price    float64
	Formula  string
	RuleName string
}
