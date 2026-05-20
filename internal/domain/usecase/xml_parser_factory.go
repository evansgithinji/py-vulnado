package usecase

import "goapp/internal/domain/entity"

// XmlParserFactory creates XML parser configurations
type XmlParserFactory struct{}

// NewXmlParserFactory creates a new XmlParserFactory
func NewXmlParserFactory() *XmlParserFactory {
	return &XmlParserFactory{}
}

// CreateParser creates xmllint arguments based on config
// VULNERABLE: Creates parser with entities enabled
func (f *XmlParserFactory) CreateParser(config *entity.XmlParserConfig) []string {
	args := []string{}
	if config.ResolveEntities {
		args = append(args, "--noent")
	}
	if !config.NoNetwork {
		// Network access allowed - vulnerable to SSRF via XXE
	}
	return args
}
