package usecase

import (
	"path/filepath"

	"goapp/internal/domain/entity"
)

// PathResolver resolves file paths according to policy
type PathResolver struct{}

// NewPathResolver creates a new PathResolver
func NewPathResolver() *PathResolver {
	return &PathResolver{}
}

// ResolvePath resolves a filename to a full path
// VULNERABLE: No path traversal check - allows ../ sequences
func (r *PathResolver) ResolvePath(filename, baseDir string, policy *entity.FileAccessPolicy) string {
	return filepath.Join(baseDir, filename)
}
