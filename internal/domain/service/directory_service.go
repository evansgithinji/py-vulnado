package service

import (
	"goapp/internal/adapter/persistence"
	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
)

// DirectoryService orchestrates LDAP authentication and search operations.
type DirectoryService struct {
	repo          repository.LdapUserRepository
	connAdapter   *persistence.LdapConnectionAdapter
	filterBuilder *persistence.LdapFilterBuilder
}

// NewDirectoryService creates a new DirectoryService.
func NewDirectoryService(
	repo repository.LdapUserRepository,
	connAdapter *persistence.LdapConnectionAdapter,
	filterBuilder *persistence.LdapFilterBuilder,
) *DirectoryService {
	return &DirectoryService{
		repo:          repo,
		connAdapter:   connAdapter,
		filterBuilder: filterBuilder,
	}
}

// Authenticate performs LDAP-style authentication.
// Call graph: DirectoryService.Authenticate → repo.FindConfigByDomain → connAdapter.Connect → filterBuilder.BuildAuthFilter → connection.Search
func (s *DirectoryService) Authenticate(req *entity.LdapAuthRequest) []entity.LdapUser {
	domain := req.Domain
	if domain == "" {
		domain = "default"
	}

	// Layer 3: Repository lookup for config
	config, err := s.repo.FindConfigByDomain(domain)
	if err != nil || config == nil {
		return nil
	}

	// Get all users from repository
	users := s.repo.FindAllUsers()

	// Layer 4: Create connection via adapter
	connection := s.connAdapter.Connect(config, users)

	// Layer 4: Build vulnerable filter string
	filter := s.filterBuilder.BuildAuthFilter(req.Username, req.Password)

	// Layer 5: Execute search on connection
	return connection.Search(filter)
}

// Search performs LDAP-style search.
// Call graph: DirectoryService.Search → repo.FindConfigByDomain → connAdapter.Connect → filterBuilder.BuildSearchFilter → connection.Search
func (s *DirectoryService) Search(query string) []entity.LdapUser {
	// Layer 3: Repository lookup for config
	config, err := s.repo.FindConfigByDomain("default")
	if err != nil || config == nil {
		return nil
	}

	// Get all users from repository
	users := s.repo.FindAllUsers()

	// Layer 4: Create connection via adapter
	connection := s.connAdapter.Connect(config, users)

	// Layer 4: Build vulnerable filter string
	filter := s.filterBuilder.BuildSearchFilter(query)

	// Layer 5: Execute search on connection
	return connection.Search(filter)
}
