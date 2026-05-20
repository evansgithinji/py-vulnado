package service

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// NotificationService orchestrates template rendering with deep call graph
type NotificationService struct {
	templateRepo *repository.TemplateConfigRepository
	compiler     *usecase.TemplateCompiler
	engine       *usecase.TemplateEngine
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(
	templateRepo *repository.TemplateConfigRepository,
	compiler *usecase.TemplateCompiler,
	engine *usecase.TemplateEngine,
) *NotificationService {
	return &NotificationService{
		templateRepo: templateRepo,
		compiler:     compiler,
		engine:       engine,
	}
}

// RenderNotification renders a notification through the deep call graph
// Deep call graph for SSTI (depth 5):
// Handler -> NotificationService.RenderNotification()
//   -> TemplateConfigRepository.LoadTemplate() -> TemplateDefinition
//   -> TemplateCompiler.Compile() -> compiled (VULNERABLE: user input in template)
//   -> TemplateEngine.Render() -> output (VULNERABLE: SSTI)
func (s *NotificationService) RenderNotification(request *entity.RenderRequest) (string, error) {
	definition := s.templateRepo.LoadTemplate(request.TemplateName)
	compiled := s.compiler.Compile(definition, request.UserInput)
	return s.engine.Render(compiled, request.Variables)
}

// RenderInvoice renders an invoice through the deep call graph
func (s *NotificationService) RenderInvoice(templateStr string, variables map[string]interface{}) (string, error) {
	definition := s.templateRepo.LoadTemplate("invoice")
	if templateStr != "" {
		definition.TemplateString = templateStr
	}
	compiled := s.compiler.Compile(definition, templateStr)
	return s.engine.Render(compiled, variables)
}
