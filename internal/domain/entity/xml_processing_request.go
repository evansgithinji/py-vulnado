package entity

// XmlProcessingRequest represents an XML processing request DTO
type XmlProcessingRequest struct {
	XmlContent string
	Operation  string
	ConfigName string
}

// XmlParserConfig defines the configuration for the XML parser
type XmlParserConfig struct {
	Name            string
	ResolveEntities bool // VULNERABLE: entities enabled
	LoadDTD         bool // VULNERABLE: DTD loading enabled
	NoNetwork       bool // VULNERABLE: allows network access
}
