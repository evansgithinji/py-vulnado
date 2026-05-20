package repository

import "goapp/internal/domain/entity"

// TemplateConfigRepository stores template definitions
type TemplateConfigRepository struct {
	templates map[string]*entity.TemplateDefinition
}

// NewTemplateConfigRepository creates a new TemplateConfigRepository
func NewTemplateConfigRepository() *TemplateConfigRepository {
	return &TemplateConfigRepository{
		templates: map[string]*entity.TemplateDefinition{
			"greeting": {
				Name:           "greeting",
				TemplateString: "Hello, {user_input}! Welcome to our store.",
				Engine:         "go",
			},
			"invoice": {
				Name:           "invoice",
				TemplateString: "Invoice for {{.Customer}}: ${{.Amount}}",
				Engine:         "go",
			},
			"custom": {
				Name:           "custom",
				TemplateString: "",
				Engine:         "go",
			},
		},
	}
}

// LoadTemplate retrieves a template definition by name
func (r *TemplateConfigRepository) LoadTemplate(name string) *entity.TemplateDefinition {
	if tmpl, ok := r.templates[name]; ok {
		// Return a copy to allow mutation
		copy := *tmpl
		return &copy
	}
	def := r.templates["greeting"]
	copy := *def
	return &copy
}
