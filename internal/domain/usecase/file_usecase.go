package usecase

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FileUseCase handles file operations
type FileUseCase struct {
	uploadDir string
	filesDir  string
	imagesDir string
}

// NewFileUseCase creates a new FileUseCase
func NewFileUseCase(uploadDir, filesDir, imagesDir string) *FileUseCase {
	return &FileUseCase{
		uploadDir: uploadDir,
		filesDir:  filesDir,
		imagesDir: imagesDir,
	}
}

// ReadFile reads a file from the files directory
// VULNERABLE: Path traversal - filepath.Join doesn't prevent ../
func (uc *FileUseCase) ReadFile(filename string) ([]byte, error) {
	// VULNERABLE: No path sanitization
	path := filepath.Join(uc.filesDir, filename)
	return os.ReadFile(path)
}

// ReadImage reads an image file
// VULNERABLE: Path traversal
func (uc *FileUseCase) ReadImage(filename string) ([]byte, error) {
	// VULNERABLE: Direct concatenation
	path := uc.imagesDir + "/" + filename
	return os.ReadFile(path)
}

// DownloadFile serves a file for download
// VULNERABLE: Path traversal via direct concatenation
func (uc *FileUseCase) DownloadFile(filename string) (*os.File, error) {
	// VULNERABLE: No sanitization
	return os.Open("./static/" + filename)
}

// SaveUploadedFile saves an uploaded file
// VULNERABLE: No file type validation, uses original filename
func (uc *FileUseCase) SaveUploadedFile(filename string, content io.Reader) (string, error) {
	// VULNERABLE: Uses original filename without sanitization
	destPath := filepath.Join(uc.uploadDir, filename)

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// VULNERABLE: No file size limit
	_, err = io.Copy(out, content)
	if err != nil {
		return "", err
	}

	return destPath, nil
}

// FetchRemoteFile fetches a file from a remote URL
// VULNERABLE: SSRF - no URL validation
func (uc *FileUseCase) FetchRemoteFile(url string) ([]byte, error) {
	// VULNERABLE: Direct HTTP request without validation
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// FetchRemoteImage fetches an image from URL and saves it
// VULNERABLE: SSRF + arbitrary file write
func (uc *FileUseCase) FetchRemoteImage(url, filename string) error {
	// VULNERABLE: SSRF
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// VULNERABLE: No content type validation
	destPath := filepath.Join(uc.imagesDir, filename)
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// IncludeRemoteContent fetches and includes remote content
// VULNERABLE: RFI - Remote File Inclusion
func (uc *FileUseCase) IncludeRemoteContent(url string) (string, error) {
	// VULNERABLE: Fetches arbitrary remote content
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

// ExportData exports database data to a file
// VULNERABLE: Command injection
func (uc *FileUseCase) ExportData(table, filename string) (string, error) {
	// VULNERABLE: Command injection via table and filename parameters
	cmd := fmt.Sprintf("sqlite3 /app/data/app.db 'SELECT * FROM %s' > /app/exports/%s.csv", table, filename)
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("export failed: %s - %v", string(output), err)
	}
	return string(output), nil
}

// ReadFileDirect reads a file directly without base path
// VULNERABLE: Path traversal - absolute path control
func (uc *FileUseCase) ReadFileDirect(filename string) ([]byte, error) {
	// VULNERABLE: Direct file read with user-controlled path
	return os.ReadFile(filename)
}
