package usecase

import "os"

// FileReader reads files from disk
type FileReader struct{}

// NewFileReader creates a new FileReader
func NewFileReader() *FileReader {
	return &FileReader{}
}

// ReadFile reads a file and returns its contents as a string
// VULNERABLE: Reads any file path without restriction
func (r *FileReader) ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadFileBytes reads a file and returns its contents as bytes
// VULNERABLE: Reads any file path without restriction
func (r *FileReader) ReadFileBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}
