package usecase

import (
	"bytes"
	"fmt"
	"html/template"
)

// TemplateUseCase handles template rendering operations
type TemplateUseCase struct{}

// NewTemplateUseCase creates a new TemplateUseCase
func NewTemplateUseCase() *TemplateUseCase {
	return &TemplateUseCase{}
}

// TemplateData represents data passed to templates
type TemplateData struct {
	OrderID   string
	Total     float64
	Username  string
	Secret    string // VULNERABLE: Sensitive data exposed to templates
	APIKey    string // VULNERABLE: Sensitive data exposed to templates
	AdminPass string // VULNERABLE: Sensitive data exposed to templates
}

// RenderGreeting renders a greeting message
// VULNERABLE: SSTI - user input interpolated into template string
func (uc *TemplateUseCase) RenderGreeting(message string) (string, error) {
	// VULNERABLE: User input directly in template string
	tmpl := fmt.Sprintf(`
		<html><body>
			<h1>%s</h1>
			<p>Try: {{.}} or {{printf "pwned"}}</p>
		</body></html>
	`, message)

	t, err := template.New("page").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, message)
	return buf.String(), err
}

// RenderDynamicTemplate renders a user-provided template
// VULNERABLE: SSTI - entire template from user input
func (uc *TemplateUseCase) RenderDynamicTemplate(userTemplate string, data interface{}) (string, error) {
	// VULNERABLE: User-controlled template string
	tmpl, err := template.New("dynamic").Parse(userTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}

// RenderInvoice renders an invoice with a custom template
// VULNERABLE: SSTI with sensitive data exposure
func (uc *TemplateUseCase) RenderInvoice(templateStr, orderID string, total float64) (string, error) {
	if templateStr == "" {
		templateStr = "Order #{{.OrderID}} - Total: ${{.Total}}"
	}

	// VULNERABLE: User-controlled template with sensitive data
	tmpl, err := template.New("invoice").Parse(templateStr)
	if err != nil {
		return "", err
	}

	data := TemplateData{
		OrderID:   orderID,
		Total:     total,
		Secret:    "ADMIN_API_KEY_12345",
		APIKey:    "sk-live-abcdef123456",
		AdminPass: "super_secret_admin_pass",
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}

// RenderMessage renders a message template
// VULNERABLE: SSTI via message interpolation
func (uc *TemplateUseCase) RenderMessage(msg string) (string, error) {
	// VULNERABLE: User input in template
	tmpl := fmt.Sprintf(`
		<html><body>
			<div>%s</div>
		</body></html>
	`, msg)

	t := template.New("page")
	t, err := t.Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, msg)
	return buf.String(), err
}

// RenderWithFunctions renders a template with exposed functions
// VULNERABLE: SSTI with exposed dangerous functions
func (uc *TemplateUseCase) RenderWithFunctions(templateStr string) (string, error) {
	// VULNERABLE: Template with exposed functions
	data := map[string]interface{}{
		"Name": "User",
		// VULNERABLE: Exposing system-like functions
		"System": func(cmd string) string {
			return fmt.Sprintf("System command: %s", cmd)
		},
		"Env": func(key string) string {
			return fmt.Sprintf("ENV[%s]", key)
		},
		"Secret":    "ADMIN_API_KEY_12345",
		"APIKey":    "sk-live-abcdef123456",
		"AdminPass": "super_secret_admin_pass",
	}

	tmpl, err := template.New("funcs").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.String(), err
}
