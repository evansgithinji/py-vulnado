package usecase

import (
	"bytes"
	"fmt"
	"html/template"
)

// TemplateEngine renders compiled templates with variables
type TemplateEngine struct{}

// NewTemplateEngine creates a new TemplateEngine
func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{}
}

// Render renders a compiled template string with variables
// VULNERABLE: Renders Go template with user-controlled content (SSTI)
func (e *TemplateEngine) Render(compiledTemplate string, variables map[string]interface{}) (string, error) {
	tmpl, err := template.New("dynamic").Parse(compiledTemplate)
	if err != nil {
		// If template parsing fails, return the raw string (still vulnerable to XSS)
		return fmt.Sprintf("<html><body>%s</body></html>", compiledTemplate), nil
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, variables)
	if err != nil {
		return fmt.Sprintf("<html><body>%s</body></html>", compiledTemplate), nil
	}

	return buf.String(), nil
}
