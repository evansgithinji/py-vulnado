package usecase

import (
	"fmt"

	"goapp/internal/domain/entity"
)

// TemplateCompiler compiles template definitions with user input
type TemplateCompiler struct{}

// NewTemplateCompiler creates a new TemplateCompiler
func NewTemplateCompiler() *TemplateCompiler {
	return &TemplateCompiler{}
}

// Compile compiles a template definition with user input
// VULNERABLE: Injects user input directly into template string
func (c *TemplateCompiler) Compile(definition *entity.TemplateDefinition, userInput string) string {
	if definition.Name == "greeting" {
		return fmt.Sprintf("Hello, %s! Welcome to our store.", userInput)
	} else if definition.Name == "custom" {
		return userInput // VULNERABLE: entire template from user
	}
	return definition.TemplateString
}
