package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"goapp/internal/domain/entity"
	"goapp/internal/domain/service"
	"goapp/internal/domain/usecase"
)

// XmlHandler handles XML-related HTTP requests
type XmlHandler struct {
	xmlProcessingService *service.XmlProcessingService
	xmlUseCase           *usecase.XmlUseCase // kept for additional endpoints
}

// NewXmlHandler creates a new XmlHandler
func NewXmlHandler(xmlProcessingService *service.XmlProcessingService, xmlUseCase *usecase.XmlUseCase) *XmlHandler {
	return &XmlHandler{
		xmlProcessingService: xmlProcessingService,
		xmlUseCase:           xmlUseCase,
	}
}

// RegisterRoutes registers XML routes
func (h *XmlHandler) RegisterRoutes(r *gin.Engine) {
	// XML parsing endpoints
	r.POST("/api/xml/parse", h.ParseXML)
	r.GET("/api/xml/parse", h.ParseXMLGet)
	r.POST("/api/xml/format", h.FormatXML)
	r.POST("/api/xml/validate", h.ValidateXML)

	// XPath query endpoint
	r.POST("/api/xpath", h.ExecuteXPath)
	r.GET("/api/xpath", h.ExecuteXPathGet)

	// XSLT transformation
	r.POST("/api/xml/transform", h.TransformXML)

	// Business logic endpoints
	r.POST("/api/products/import", h.ImportProductFeed)
	r.POST("/api/config/import", h.ImportUserConfig)
	r.POST("/api/reports/generate", h.GenerateReport)

	// File and URL based parsing
	r.GET("/api/xml/file", h.ParseFromFile)
	r.GET("/api/xml/url", h.ParseFromURL)

	// HTML pages for testing
	r.GET("/xml", h.XmlPage)
	r.GET("/xml/xpath", h.XPathPage)
	r.GET("/xml/products", h.ProductImportPage)
	r.GET("/xml/config", h.ConfigImportPage)
}

// ParseXML parses XML content (POST)
// VULNERABLE: XXE via XmlProcessingService deep call graph
func (h *XmlHandler) ParseXML(c *gin.Context) {
	xmlContent, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	if len(xmlContent) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "XML content required"})
		return
	}

	// VULNERABLE: XXE via deep call graph
	request := &entity.XmlProcessingRequest{
		XmlContent: string(xmlContent),
		Operation:  "parse",
		ConfigName: "default",
	}
	result, err := h.xmlProcessingService.ProcessXml(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "XML parsing failed",
			"detail": err.Error(),
			"output": result,
		})
		return
	}

	c.String(http.StatusOK, result)
}

// ParseXMLGet parses XML from query parameter (GET)
// VULNERABLE: XXE
func (h *XmlHandler) ParseXMLGet(c *gin.Context) {
	xmlContent := c.Query("xml")
	if xmlContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "xml parameter required"})
		return
	}

	// VULNERABLE: XXE via deep call graph
	request := &entity.XmlProcessingRequest{
		XmlContent: xmlContent,
		Operation:  "parse",
		ConfigName: "default",
	}
	result, err := h.xmlProcessingService.ProcessXml(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "XML parsing failed",
			"detail": err.Error(),
			"output": result,
		})
		return
	}

	c.String(http.StatusOK, result)
}

// FormatXML formats XML with custom options
// VULNERABLE: XXE + Command Injection combo
func (h *XmlHandler) FormatXML(c *gin.Context) {
	xmlContent, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	format := c.Query("format")

	// VULNERABLE: XXE + Command Injection (falls back to xmlUseCase for format param)
	result, err := h.xmlUseCase.ParseXMLWithFormat(string(xmlContent), format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "XML formatting failed",
			"detail": err.Error(),
			"output": result,
		})
		return
	}

	c.String(http.StatusOK, result)
}

// ValidateXML validates XML syntax
// VULNERABLE: XXE via deep call graph
func (h *XmlHandler) ValidateXML(c *gin.Context) {
	xmlContent, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// VULNERABLE: XXE via deep call graph
	request := &entity.XmlProcessingRequest{
		XmlContent: string(xmlContent),
		Operation:  "validate",
		ConfigName: "validate",
	}
	result, err := h.xmlProcessingService.ValidateXml(request)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"valid":   false,
			"message": result,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"message": result,
	})
}

// ExecuteXPath executes XPath query (POST)
// VULNERABLE: XXE + XPath Injection
func (h *XmlHandler) ExecuteXPath(c *gin.Context) {
	xmlContent, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	xpath := c.Query("xpath")
	if xpath == "" {
		xpath = c.PostForm("xpath")
	}

	if xpath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "xpath parameter required"})
		return
	}

	// VULNERABLE: XXE + XPath injection (uses xmlUseCase for xpath)
	result, err := h.xmlUseCase.ExecuteXPath(string(xmlContent), xpath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "XPath query failed",
			"detail": err.Error(),
			"output": result,
		})
		return
	}

	c.String(http.StatusOK, result)
}

// ExecuteXPathGet executes XPath query (GET)
// VULNERABLE: XXE + XPath Injection
func (h *XmlHandler) ExecuteXPathGet(c *gin.Context) {
	xmlContent := c.Query("xml")
	xpath := c.Query("xpath")

	if xmlContent == "" || xpath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "xml and xpath parameters required"})
		return
	}

	// VULNERABLE: XXE + XPath injection
	result, err := h.xmlUseCase.ExecuteXPath(xmlContent, xpath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "XPath query failed",
			"detail": err.Error(),
			"output": result,
		})
		return
	}

	c.String(http.StatusOK, result)
}

// TransformXML transforms XML using XSLT
// VULNERABLE: XXE + XSLT injection
func (h *XmlHandler) TransformXML(c *gin.Context) {
	xmlContent := c.PostForm("xml")
	xsltContent := c.PostForm("xslt")

	if xmlContent == "" || xsltContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "xml and xslt parameters required"})
		return
	}

	result, err := h.xmlUseCase.TransformXML(xmlContent, xsltContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "XSLT transformation failed",
			"detail": err.Error(),
			"output": result,
		})
		return
	}

	c.String(http.StatusOK, result)
}

// ImportProductFeed imports product feed from XML
// VULNERABLE: XXE via deep call graph
func (h *XmlHandler) ImportProductFeed(c *gin.Context) {
	feedXML, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// VULNERABLE: XXE via deep call graph
	request := &entity.XmlProcessingRequest{
		XmlContent: string(feedXML),
		Operation:  "parse",
		ConfigName: "default",
	}
	result, err := h.xmlProcessingService.ProcessXml(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Product feed import failed",
			"detail":  err.Error(),
			"partial": result,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product feed imported successfully",
		"data":    result,
	})
}

// ImportUserConfig imports user configuration from XML
// VULNERABLE: XXE via deep call graph
func (h *XmlHandler) ImportUserConfig(c *gin.Context) {
	configXML, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// VULNERABLE: XXE via deep call graph
	request := &entity.XmlProcessingRequest{
		XmlContent: string(configXML),
		Operation:  "parse",
		ConfigName: "default",
	}
	result, err := h.xmlProcessingService.ProcessXml(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Config import failed",
			"detail":  err.Error(),
			"partial": result,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User configuration imported successfully",
		"config":  result,
	})
}

// GenerateReport generates a report from XML data
// VULNERABLE: XXE
func (h *XmlHandler) GenerateReport(c *gin.Context) {
	reportXML, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	outputFormat := c.Query("format")
	if outputFormat == "" {
		outputFormat = "standard"
	}

	result, err := h.xmlUseCase.GenerateReport(string(reportXML), outputFormat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Report generation failed",
			"detail": err.Error(),
			"output": result,
		})
		return
	}

	c.String(http.StatusOK, result)
}

// ParseFromFile parses XML from a file
// VULNERABLE: XXE + Path Traversal
func (h *XmlHandler) ParseFromFile(c *gin.Context) {
	filePath := c.Query("file")
	if filePath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file parameter required"})
		return
	}

	result, err := h.xmlUseCase.ParseFromFile(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to parse file",
			"detail": err.Error(),
		})
		return
	}

	c.String(http.StatusOK, result)
}

// ParseFromURL parses XML from a URL
// VULNERABLE: XXE + SSRF
func (h *XmlHandler) ParseFromURL(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter required"})
		return
	}

	result, err := h.xmlUseCase.ParseFromURL(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Failed to parse URL",
			"detail": err.Error(),
		})
		return
	}

	c.String(http.StatusOK, result)
}

// XmlPage serves the XML parsing test page (unchanged)
func (h *XmlHandler) XmlPage(c *gin.Context) {
	html := `<!DOCTYPE html><html><head><title>XML Parser - Vulnerable App</title></head><body><h1>XML Parser (XXE Vulnerability Demo)</h1><form method="POST" action="/api/xml/parse"><textarea name="xml" rows="10" cols="60"><?xml version="1.0"?><product><name>Sample</name><price>99.99</price></product></textarea><br><button type="submit">Parse XML</button></form><p><a href="/xml/xpath">XPath</a> | <a href="/xml/products">Products</a> | <a href="/xml/config">Config</a></p></body></html>`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// XPathPage serves the XPath query test page (unchanged)
func (h *XmlHandler) XPathPage(c *gin.Context) {
	html := `<!DOCTYPE html><html><head><title>XPath Query - Vulnerable App</title></head><body><h1>XPath Query (XXE + XPath Injection Demo)</h1><form method="POST" action="/api/xpath"><textarea name="xml" rows="10" cols="60"><?xml version="1.0"?><users><user id="1"><username>admin</username><password>admin123</password></user></users></textarea><br><input type="text" name="xpath" value="//user[@id='1']/username/text()" size="60"><br><button type="submit">Execute XPath</button></form></body></html>`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// ProductImportPage serves the product import page (unchanged)
func (h *XmlHandler) ProductImportPage(c *gin.Context) {
	html := `<!DOCTYPE html><html><head><title>Product Feed Import</title></head><body><h1>Product Feed Import (XXE in Business Context)</h1><form method="POST" action="/api/products/import"><textarea name="xml" rows="10" cols="60"><?xml version="1.0"?><products><product id="1"><name>Laptop</name><price>999.99</price></product></products></textarea><br><button type="submit">Import</button></form></body></html>`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// ConfigImportPage serves the config import page (unchanged)
func (h *XmlHandler) ConfigImportPage(c *gin.Context) {
	html := `<!DOCTYPE html><html><head><title>Config Import</title></head><body><h1>User Configuration Import (XXE in Business Context)</h1><form method="POST" action="/api/config/import"><textarea name="xml" rows="10" cols="60"><?xml version="1.0"?><config><settings><theme>dark</theme></settings></config></textarea><br><button type="submit">Import</button></form></body></html>`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
