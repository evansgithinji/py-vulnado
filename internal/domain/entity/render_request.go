package entity

// RenderRequest represents a template render request DTO
type RenderRequest struct {
	UserInput    string
	TemplateName string
	Variables    map[string]interface{}
}

// TemplateDefinition defines a stored template
type TemplateDefinition struct {
	Name           string
	TemplateString string
	Engine         string
}
