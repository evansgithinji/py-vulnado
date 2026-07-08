package service

import (
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// XmlProcessingService orchestrates XML processing with deep call graph
type XmlProcessingService struct {
	configRepo        *repository.XmlParserConfigRepository
	parserFactory     *usecase.XmlParserFactory
	documentProcessor *usecase.XmlDocumentProcessor
}

// NewXmlProcessingService creates a new XmlProcessingService
func NewXmlProcessingService(
	configRepo *repository.XmlParserConfigRepository,
	parserFactory *usecase.XmlParserFactory,
	documentProcessor *usecase.XmlDocumentProcessor,
) *XmlProcessingService {
	return &XmlProcessingService{
		configRepo:        configRepo,
		parserFactory:     parserFactory,
		documentProcessor: documentProcessor,
	}
}

// ProcessXml processes XML content through the deep call graph
// Deep call graph for XXE (depth 5):
// Handler -> XmlProcessingService.ProcessXml()
//   -> XmlParserConfigRepository.GetConfig() -> XmlParserConfig
//   -> XmlParserFactory.CreateParser() -> parser args (VULNERABLE: entities enabled)
//   -> XmlDocumentProcessor.Parse() -> result (VULNERABLE: XXE)
func (s *XmlProcessingService) ProcessXml(request *entity.XmlProcessingRequest) (string, error) {
	config := s.configRepo.GetConfig(request.ConfigName)
	parserArgs := s.parserFactory.CreateParser(config)
	return s.documentProcessor.Parse(parserArgs, request.XmlContent)
}

// ValidateXml validates XML content through the deep call graph
func (s *XmlProcessingService) ValidateXml(request *entity.XmlProcessingRequest) (string, error) {
	config := s.configRepo.GetConfig("validate")
	parserArgs := s.parserFactory.CreateParser(config)
	return s.documentProcessor.Validate(parserArgs, request.XmlContent)
}
