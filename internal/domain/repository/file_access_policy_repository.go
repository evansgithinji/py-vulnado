package repository

import "goapp/internal/domain/entity"

// FileAccessPolicyRepository stores file access policies
type FileAccessPolicyRepository struct {
	policies map[string]*entity.FileAccessPolicy
}

// NewFileAccessPolicyRepository creates a new FileAccessPolicyRepository
func NewFileAccessPolicyRepository(filesDir, uploadDir string) *FileAccessPolicyRepository {
	return &FileAccessPolicyRepository{
		policies: map[string]*entity.FileAccessPolicy{
			"read": {
				Name:              "read",
				BaseDirectory:     filesDir,
				AllowedExtensions: nil,
				CheckTraversal:    false, // VULNERABLE
			},
			"download": {
				Name:              "download",
				BaseDirectory:     filesDir,
				AllowedExtensions: nil,
				CheckTraversal:    false, // VULNERABLE
			},
			"upload": {
				Name:              "upload",
				BaseDirectory:     uploadDir,
				AllowedExtensions: nil,
				CheckTraversal:    false,
			},
			"info": {
				Name:              "info",
				BaseDirectory:     filesDir,
				AllowedExtensions: nil,
				CheckTraversal:    false,
			},
		},
	}
}

// GetPolicy retrieves a file access policy by name
func (r *FileAccessPolicyRepository) GetPolicy(name string) *entity.FileAccessPolicy {
	if policy, ok := r.policies[name]; ok {
		return policy
	}
	return r.policies["read"]
}
