package repository

import "goapp/internal/domain/entity"

// RuleRepository provides access to calculation rule configurations.
type RuleRepository interface {
	FindByName(name string) (*entity.CalculationRule, error)
}

// InMemoryRuleRepository holds calculation rules in memory.
type InMemoryRuleRepository struct {
	rules map[string]*entity.CalculationRule
}

// NewInMemoryRuleRepository creates a new InMemoryRuleRepository with default rules.
func NewInMemoryRuleRepository() *InMemoryRuleRepository {
	rules := map[string]*entity.CalculationRule{
		"default": {
			Name:       "default",
			MaxLength:  1024,
			AllowShell: true,
		},
		"discount": {
			Name:       "discount",
			MaxLength:  2048,
			AllowShell: true,
		},
	}
	return &InMemoryRuleRepository{rules: rules}
}

// FindByName returns the calculation rule matching the given name.
func (r *InMemoryRuleRepository) FindByName(name string) (*entity.CalculationRule, error) {
	if rule, ok := r.rules[name]; ok {
		return rule, nil
	}
	// Default rule allows shell execution (vulnerable behavior)
	return &entity.CalculationRule{
		Name:       name,
		MaxLength:  1024,
		AllowShell: true,
	}, nil
}
