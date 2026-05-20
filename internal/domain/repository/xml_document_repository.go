package repository

import "goapp/internal/domain/entity"

// XmlDocumentRepository provides access to XML document data.
type XmlDocumentRepository interface {
	FindBySourceID(sourceID string) (*entity.XmlDocument, error)
}

// InMemoryXmlDocumentRepository holds XML documents in memory.
type InMemoryXmlDocumentRepository struct {
	documents map[string]*entity.XmlDocument
}

// NewInMemoryXmlDocumentRepository creates a new InMemoryXmlDocumentRepository with default user XML.
func NewInMemoryXmlDocumentRepository(tempDir string) *InMemoryXmlDocumentRepository {
	usersXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<users>
  <user><username>admin</username><password>admin123</password><name>Admin</name><role>admin</role></user>
  <user><username>john</username><password>password</password><name>John Doe</name><role>user</role></user>
  <user><username>jane</username><password>password123</password><name>Jane Smith</name><role>user</role></user>
</users>`)

	docs := map[string]*entity.XmlDocument{
		"users": {
			Content:  usersXML,
			SourceID: "users",
			TempDir:  tempDir,
		},
	}

	return &InMemoryXmlDocumentRepository{documents: docs}
}

// FindBySourceID returns the XML document matching the given source ID.
func (r *InMemoryXmlDocumentRepository) FindBySourceID(sourceID string) (*entity.XmlDocument, error) {
	if doc, ok := r.documents[sourceID]; ok {
		return doc, nil
	}
	return nil, nil
}
