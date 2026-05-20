package service

import (
	"os"

	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// FileService orchestrates file operations with deep call graph
type FileService struct {
	policyRepo   *repository.FileAccessPolicyRepository
	pathResolver *usecase.PathResolver
	fileReader   *usecase.FileReader
}

// NewFileService creates a new FileService
func NewFileService(
	policyRepo *repository.FileAccessPolicyRepository,
	pathResolver *usecase.PathResolver,
	fileReader *usecase.FileReader,
) *FileService {
	return &FileService{
		policyRepo:   policyRepo,
		pathResolver: pathResolver,
		fileReader:   fileReader,
	}
}

// ProcessFileRequest processes a file read request through the deep call graph
// Deep call graph for Path Traversal (depth 5):
// Handler -> FileService.ProcessFileRequest()
//   -> FileAccessPolicyRepository.GetPolicy() -> FileAccessPolicy
//   -> PathResolver.ResolvePath() -> path (VULNERABLE: no .. check)
//   -> FileReader.ReadFile() -> content (VULNERABLE: reads any path)
func (s *FileService) ProcessFileRequest(request *entity.FileRequest) (string, error) {
	policy := s.policyRepo.GetPolicy(request.Operation)
	baseDir := request.BaseDir
	if baseDir == "" {
		baseDir = policy.BaseDirectory
	}
	path := s.pathResolver.ResolvePath(request.Filename, baseDir, policy)
	return s.fileReader.ReadFile(path)
}

// DownloadFile processes a file download request through the deep call graph
func (s *FileService) DownloadFile(request *entity.FileRequest) ([]byte, error) {
	policy := s.policyRepo.GetPolicy("download")
	baseDir := request.BaseDir
	if baseDir == "" {
		baseDir = policy.BaseDirectory
	}
	path := s.pathResolver.ResolvePath(request.Filename, baseDir, policy)
	return s.fileReader.ReadFileBytes(path)
}

// UploadFile saves an uploaded file
func (s *FileService) UploadFile(filename string, content []byte, uploadDir string) (string, error) {
	policy := s.policyRepo.GetPolicy("upload")
	path := s.pathResolver.ResolvePath(filename, uploadDir, policy)
	return path, os.WriteFile(path, content, 0644)
}

// ListFiles lists files in a directory
func (s *FileService) ListFiles(directory string) []string {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil
	}
	var files []string
	for _, entry := range entries {
		files = append(files, entry.Name())
	}
	return files
}
