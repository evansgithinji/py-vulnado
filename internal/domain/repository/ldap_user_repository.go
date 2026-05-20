package repository

import "goapp/internal/domain/entity"

// LdapUserRepository provides access to LDAP user data and configuration.
type LdapUserRepository interface {
	FindConfigByDomain(domain string) (*entity.LdapConfig, error)
	FindAllUsers() []entity.LdapUser
}

// InMemoryLdapUserRepository is an in-memory implementation of LdapUserRepository.
type InMemoryLdapUserRepository struct {
	configs map[string]*entity.LdapConfig
	users   []entity.LdapUser
}

// NewInMemoryLdapUserRepository creates a new InMemoryLdapUserRepository with default data.
func NewInMemoryLdapUserRepository() *InMemoryLdapUserRepository {
	configs := map[string]*entity.LdapConfig{
		"default": {
			Domain:        "default",
			BaseDN:        "dc=example,dc=com",
			Host:          "localhost",
			Port:          389,
			BindDN:        "cn=admin,dc=example,dc=com",
			AdminPassword: "adminpass",
		},
	}

	users := []entity.LdapUser{
		{UID: "admin", CN: "Admin User", Mail: "admin@example.com", Password: "admin123", OU: "admins"},
		{UID: "john", CN: "John Doe", Mail: "john@example.com", Password: "john456", OU: "users"},
		{UID: "jane", CN: "Jane Smith", Mail: "jane@example.com", Password: "jane789", OU: "users"},
		{UID: "svc_backup", CN: "Backup Service", Mail: "backup@example.com", Password: "backup!@#", OU: "services"},
	}

	return &InMemoryLdapUserRepository{
		configs: configs,
		users:   users,
	}
}

// FindConfigByDomain returns the LDAP configuration for the given domain.
func (r *InMemoryLdapUserRepository) FindConfigByDomain(domain string) (*entity.LdapConfig, error) {
	if cfg, ok := r.configs[domain]; ok {
		return cfg, nil
	}
	// Fall back to default config
	if cfg, ok := r.configs["default"]; ok {
		return cfg, nil
	}
	return nil, nil
}

// FindAllUsers returns all users in the repository.
func (r *InMemoryLdapUserRepository) FindAllUsers() []entity.LdapUser {
	result := make([]entity.LdapUser, len(r.users))
	copy(result, r.users)
	return result
}
