package handler

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
)

// FileHandler handles file-related HTTP requests
type FileHandler struct {
	fileService            *service.FileService
	externalRequestService *service.ExternalRequestService
	commandService         *service.SystemCommandService
	filesDir               string
	uploadDir              string
}

// NewFileHandler creates a new FileHandler
func NewFileHandler(
	fileService *service.FileService,
	externalRequestService *service.ExternalRequestService,
	commandService *service.SystemCommandService,
	filesDir string,
	uploadDir string,
) *FileHandler {
	return &FileHandler{
		fileService:            fileService,
		externalRequestService: externalRequestService,
		commandService:         commandService,
		filesDir:               filesDir,
		uploadDir:              uploadDir,
	}
}

// RegisterRoutes registers file routes
func (h *FileHandler) RegisterRoutes(r *gin.Engine) {
	// Path Traversal endpoints
	r.GET("/api/files", h.ReadFile)
	r.GET("/api/files/direct", h.ReadFileDirect)
	r.GET("/api/images", h.ReadImage)
	r.GET("/api/download", h.DownloadFile)

	// File Upload endpoints
	r.POST("/api/upload", h.UploadFile)
	r.Static("/uploads", "./uploads")

	// SSRF endpoints
	r.GET("/api/fetch", h.FetchURL)
	r.POST("/api/fetch", h.FetchURLPost)
	r.GET("/api/import-image", h.ImportImage)

	// RFI endpoint
	r.GET("/api/include", h.IncludeRemote)
	r.POST("/api/include", h.IncludeRemotePost)

	// Export endpoint (Command Injection)
	r.GET("/api/export", h.ExportData)

	// HTML pages
	r.GET("/files", h.FilesPage)
	r.GET("/upload", h.UploadPage)
	r.GET("/fetch", h.FetchPage)
	r.GET("/include", h.IncludePage)
}

// ReadFile reads a file
// VULNERABLE: Path Traversal via FileService deep call graph
func (h *FileHandler) ReadFile(c *gin.Context) {
	filename := c.Query("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename required"})
		return
	}

	// VULNERABLE: Path traversal via deep call graph
	request := &entity.FileRequest{Filename: filename, Operation: "read"}
	content, err := h.fileService.ProcessFileRequest(request)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found", "detail": err.Error()})
		return
	}

	c.String(http.StatusOK, content)
}

// ReadFileDirect reads any file directly
// VULNERABLE: Path Traversal - direct file path control
func (h *FileHandler) ReadFileDirect(c *gin.Context) {
	filename := c.Query("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename required"})
		return
	}

	// VULNERABLE: Direct path traversal via deep call graph with empty base dir
	request := &entity.FileRequest{Filename: filename, BaseDir: "/", Operation: "read"}
	content, err := h.fileService.ProcessFileRequest(request)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found", "detail": err.Error()})
		return
	}

	c.String(http.StatusOK, content)
}

// ReadImage reads an image file
// VULNERABLE: Path Traversal
func (h *FileHandler) ReadImage(c *gin.Context) {
	filename := c.Query("file")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File parameter required"})
		return
	}

	// VULNERABLE: Path traversal via deep call graph
	request := &entity.FileRequest{Filename: filename, Operation: "read"}
	content, err := h.fileService.DownloadFile(request)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found", "detail": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", content)
}

// DownloadFile serves a file for download
// VULNERABLE: Path Traversal
func (h *FileHandler) DownloadFile(c *gin.Context) {
	filename := c.Query("name")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name parameter required"})
		return
	}

	// VULNERABLE: Path traversal via deep call graph
	request := &entity.FileRequest{Filename: filename, BaseDir: "./static", Operation: "download"}
	content, err := h.fileService.DownloadFile(request)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found", "detail": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filename)))
	c.Data(http.StatusOK, "application/octet-stream", content)
}

// UploadFile handles file uploads
// VULNERABLE: No file type validation, uses original filename
func (h *FileHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File required"})
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// VULNERABLE: Uses original filename without sanitization via deep call graph
	destPath, err := h.fileService.UploadFile(header.Filename, content, h.uploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded",
		"path":     destPath,
		"filename": header.Filename,
		"url":      "/uploads/" + header.Filename,
	})
}

// FetchURL fetches content from a URL
// VULNERABLE: SSRF via ExternalRequestService deep call graph
func (h *FileHandler) FetchURL(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL required"})
		return
	}

	// VULNERABLE: SSRF via deep call graph
	content, err := h.externalRequestService.FetchRemoteFile(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fetch failed", "detail": err.Error()})
		return
	}

	c.String(http.StatusOK, content)
}

// FetchURLPost fetches content from a URL (POST variant)
// VULNERABLE: SSRF
func (h *FileHandler) FetchURLPost(c *gin.Context) {
	url := c.PostForm("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL required"})
		return
	}

	// VULNERABLE: SSRF via deep call graph
	content, err := h.externalRequestService.FetchRemoteFile(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fetch failed", "detail": err.Error()})
		return
	}

	c.String(http.StatusOK, content)
}

// ImportImage imports an image from a URL
// VULNERABLE: SSRF + arbitrary file write
func (h *FileHandler) ImportImage(c *gin.Context) {
	url := c.Query("url")
	filename := c.Query("filename")

	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL required"})
		return
	}

	if filename == "" {
		filename = "imported_image"
	}

	// VULNERABLE: SSRF via deep call graph
	content, err := h.externalRequestService.FetchRemoteFile(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Import failed", "detail": err.Error()})
		return
	}

	_, err = h.fileService.UploadFile(filename, []byte(content), h.filesDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Save failed", "detail": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Image imported",
		"filename": filename,
	})
}

// IncludeRemote includes remote content
// VULNERABLE: RFI via ExternalRequestService deep call graph
func (h *FileHandler) IncludeRemote(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL required"})
		return
	}

	// VULNERABLE: RFI via deep call graph
	content, err := h.externalRequestService.FetchRemoteFile(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Include failed", "detail": err.Error()})
		return
	}

	// VULNERABLE: Content displayed without sanitization (XSS)
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, fmt.Sprintf("<html><body><h1>Included Content</h1><div>%s</div></body></html>", content))
}

// IncludeRemotePost includes remote content (POST variant)
// VULNERABLE: RFI
func (h *FileHandler) IncludeRemotePost(c *gin.Context) {
	url := c.PostForm("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL required"})
		return
	}

	content, err := h.externalRequestService.FetchRemoteFile(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Include failed"})
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, fmt.Sprintf("<html><body><h1>Included Content</h1><div>%s</div></body></html>", content))
}

// ExportData exports database table data
// VULNERABLE: Command Injection via SystemCommandService deep call graph
func (h *FileHandler) ExportData(c *gin.Context) {
	table := c.Query("table")
	filename := c.Query("filename")

	if table == "" || filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Table and filename required"})
		return
	}

	// VULNERABLE: Command injection via deep call graph
	target := fmt.Sprintf("sqlite3 /app/data/app.db 'SELECT * FROM %s' > /app/exports/%s.csv", table, filename)
	request := &entity.CommandRequest{Target: target, CommandType: "diagnostic"}
	output, err := h.commandService.ExecuteCommand(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "output": output})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Export completed",
		"filename": filename + ".csv",
	})
}

// FilesPage serves the files page
func (h *FileHandler) FilesPage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Files - Vulnerable App</title></head>
<body>
	<h1>File Operations (Path Traversal Demo)</h1>

	<h2>Read File</h2>
	<form method="GET" action="/api/files">
		<label>Filename: <input type="text" name="filename" placeholder="example.txt"></label>
		<button type="submit">Read</button>
	</form>

	<h2>Download File</h2>
	<form method="GET" action="/api/download">
		<label>Filename: <input type="text" name="name" placeholder="example.txt"></label>
		<button type="submit">Download</button>
	</form>

	<h2>Read Image</h2>
	<form method="GET" action="/api/images">
		<label>File: <input type="text" name="file" placeholder="logo.png"></label>
		<button type="submit">View</button>
	</form>

	<p>Hints:</p>
	<ul>
		<li>Try: ../../../etc/passwd</li>
		<li>Try: ....//....//....//etc/passwd</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// UploadPage serves the upload page
func (h *FileHandler) UploadPage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Upload - Vulnerable App</title></head>
<body>
	<h1>File Upload (Insecure Upload Demo)</h1>
	<form method="POST" action="/api/upload" enctype="multipart/form-data">
		<label>Select File: <input type="file" name="file"></label><br><br>
		<button type="submit">Upload</button>
	</form>

	<p>Hints:</p>
	<ul>
		<li>No file type validation - try uploading .php, .jsp, .exe</li>
		<li>Original filename preserved - try path injection</li>
		<li>No file size limit</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// FetchPage serves the URL fetch page
func (h *FileHandler) FetchPage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Fetch URL - Vulnerable App</title></head>
<body>
	<h1>URL Fetch (SSRF Demo)</h1>

	<h2>Fetch Content</h2>
	<form method="GET" action="/api/fetch">
		<label>URL: <input type="text" name="url" size="50" placeholder="http://example.com"></label>
		<button type="submit">Fetch (GET)</button>
	</form>

	<form method="POST" action="/api/fetch">
		<label>URL: <input type="text" name="url" size="50" placeholder="http://example.com"></label>
		<button type="submit">Fetch (POST)</button>
	</form>

	<h2>Import Image</h2>
	<form method="GET" action="/api/import-image">
		<label>URL: <input type="text" name="url" size="50" placeholder="http://example.com/image.jpg"></label>
		<label>Filename: <input type="text" name="filename" placeholder="imported"></label>
		<button type="submit">Import</button>
	</form>

	<p>Hints:</p>
	<ul>
		<li>Try: http://localhost:8080/api/debug</li>
		<li>Try: http://169.254.169.254/latest/meta-data/</li>
		<li>Try: file:///etc/passwd</li>
	</ul>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// IncludePage serves the include page
func (h *FileHandler) IncludePage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head><title>Include - Vulnerable App</title></head>
<body>
	<h1>Remote Include (RFI Demo)</h1>

	<form method="GET" action="/api/include">
		<label>URL: <input type="text" name="url" size="50" placeholder="http://example.com/content.html"></label>
		<button type="submit">Include (GET)</button>
	</form>

	<form method="POST" action="/api/include">
		<label>URL: <input type="text" name="url" size="50" placeholder="http://example.com/content.html"></label>
		<button type="submit">Include (POST)</button>
	</form>

	<p>Hint: Include remote HTML/JS content that will be rendered in the page</p>
</body>
</html>
`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
