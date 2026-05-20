package repository

import "goapp/internal/domain/entity"

// XmlParserConfigRepository stores XML parser configurations
type XmlParserConfigRepository struct {
	configs map[string]*entity.XmlParserConfig
}

// NewXmlParserConfigRepository creates a new XmlParserConfigRepository
func NewXmlParserConfigRepository() *XmlParserConfigRepository {
	return &XmlParserConfigRepository{
		configs: map[string]*entity.XmlParserConfig{
			"default": {
				Name:            "default",
				ResolveEntities: true,  // VULNERABLE
				LoadDTD:         true,  // VULNERABLE
				NoNetwork:       false, // VULNERABLE
			},
			"validate": {
				Name:            "validate",
				ResolveEntities: true,
				LoadDTD:         true,
				NoNetwork:       false,
			},
			"transform": {
				Name:            "transform",
				ResolveEntities: true,
				LoadDTD:         true,
				NoNetwork:       false,
			},
		},
	}
}

// GetConfig retrieves an XML parser configuration by name
func (r *XmlParserConfigRepository) GetConfig(name string) *entity.XmlParserConfig {
	if config, ok := r.configs[name]; ok {
		return config
	}
	return r.configs["default"]
}
