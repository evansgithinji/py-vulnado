package repository

import "goapp/internal/domain/entity"

// DocumentCollectionRepository provides access to NoSQL document collections.
type DocumentCollectionRepository interface {
	FindCollectionByName(name string) (*entity.DocumentCollection, error)
	FindAllUsers() []entity.NoSqlUser
}

// InMemoryDocumentCollectionRepository holds document collections in memory.
type InMemoryDocumentCollectionRepository struct {
	collections map[string]*entity.DocumentCollection
}

// NewInMemoryDocumentCollectionRepository creates a new InMemoryDocumentCollectionRepository with default data.
func NewInMemoryDocumentCollectionRepository() *InMemoryDocumentCollectionRepository {
	users := []entity.NoSqlUser{
		{ID: "1", Username: "admin", Password: "admin123", Role: "admin", Email: "admin@example.com"},
		{ID: "2", Username: "john", Password: "john456", Role: "user", Email: "john@example.com"},
		{ID: "3", Username: "jane", Password: "jane789", Role: "user", Email: "jane@example.com"},
	}

	collections := map[string]*entity.DocumentCollection{
		"users": {
			Name:      "users",
			Documents: users,
		},
	}

	return &InMemoryDocumentCollectionRepository{collections: collections}
}

// FindCollectionByName returns the document collection matching the given name.
func (r *InMemoryDocumentCollectionRepository) FindCollectionByName(name string) (*entity.DocumentCollection, error) {
	if col, ok := r.collections[name]; ok {
		return col, nil
	}
	return nil, nil
}

// FindAllUsers returns all user documents from the users collection.
func (r *InMemoryDocumentCollectionRepository) FindAllUsers() []entity.NoSqlUser {
	if col, ok := r.collections["users"]; ok {
		result := make([]entity.NoSqlUser, len(col.Documents))
		copy(result, col.Documents)
		return result
	}
	return nil
}
