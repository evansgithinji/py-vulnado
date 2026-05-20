package entity

// XmlDocument represents an XML document stored in the repository.
type XmlDocument struct {
	Content  []byte
	SourceID string
	TempDir  string
}

// XPathAuthRequest is a DTO for XPath-based authentication requests.
type XPathAuthRequest struct {
	Username string
	Password string
}
