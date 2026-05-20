package entity

// FileRequest represents a file access request DTO
type FileRequest struct {
	Filename  string
	BaseDir   string
	Operation string
}

// FileAccessPolicy defines the policy for file access
type FileAccessPolicy struct {
	Name              string
	BaseDirectory     string
	AllowedExtensions []string
	CheckTraversal    bool // VULNERABLE: traversal check disabled
}
