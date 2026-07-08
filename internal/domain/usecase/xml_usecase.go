package usecase

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// XmlUseCase handles XML operations
type XmlUseCase struct {
	tempDir string
}

// NewXmlUseCase creates a new XmlUseCase
func NewXmlUseCase(tempDir string) *XmlUseCase {
	// Ensure temp directory exists
	os.MkdirAll(tempDir, 0755)
	return &XmlUseCase{
		tempDir: tempDir,
	}
}

// ParseXML parses XML content and resolves external entities
// VULNERABLE: XXE - Uses xmllint with --noent which processes external entities
func (uc *XmlUseCase) ParseXML(xmlContent string) (string, error) {
	// VULNERABLE: Write XML to temp file
	tmpFile, err := os.CreateTemp(uc.tempDir, "xml-*.xml")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(xmlContent)
	if err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write XML: %v", err)
	}
	tmpFile.Close()

	// VULNERABLE: xmllint with --noent flag processes external entities
	// This allows XXE attacks including:
	// - File disclosure via file:// protocol
	// - SSRF via http:// protocol
	// - Denial of Service via billion laughs attack
	cmd := exec.Command("xmllint", "--noent", tmpFile.Name())
	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), fmt.Errorf("XML parsing error: %v", err)
	}

	return string(output), nil
}

// ParseXMLWithFormat formats XML output
// VULNERABLE: XXE + Command Injection combo
func (uc *XmlUseCase) ParseXMLWithFormat(xmlContent, format string) (string, error) {
	// Write XML to temp file
	tmpFile, err := os.CreateTemp(uc.tempDir, "xml-*.xml")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(xmlContent)
	if err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write XML: %v", err)
	}
	tmpFile.Close()

	// VULNERABLE: Command injection via format parameter
	// VULNERABLE: XXE via xmllint --noent
	cmdStr := fmt.Sprintf("xmllint --noent --format %s", tmpFile.Name())
	if format != "" {
		cmdStr = fmt.Sprintf("xmllint --noent %s %s", format, tmpFile.Name())
	}

	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

// ExecuteXPath executes XPath query on XML
// VULNERABLE: XXE + XPath Injection
func (uc *XmlUseCase) ExecuteXPath(xmlContent, xpath string) (string, error) {
	// Write XML to temp file
	tmpFile, err := os.CreateTemp(uc.tempDir, "xml-*.xml")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(xmlContent)
	if err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write XML: %v", err)
	}
	tmpFile.Close()

	// VULNERABLE: XXE via xmllint
	// VULNERABLE: XPath injection - no sanitization of xpath parameter
	// VULNERABLE: Command injection via unsanitized xpath
	cmd := exec.Command("xmllint", "--noent", "--xpath", xpath, tmpFile.Name())
	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), fmt.Errorf("XPath query error: %v", err)
	}

	return string(output), nil
}

// ParseProductFeed parses product feed XML
// VULNERABLE: XXE in business context (product feed processing)
func (uc *XmlUseCase) ParseProductFeed(feedXML string) (map[string]interface{}, error) {
	// VULNERABLE: Parse XML with external entities
	parsedXML, err := uc.ParseXML(feedXML)
	if err != nil {
		return nil, err
	}

	// Simple extraction (in real scenario, this would parse the XML properly)
	result := make(map[string]interface{})
	result["raw_xml"] = parsedXML
	result["status"] = "processed"

	// Extract some basic info if possible
	if strings.Contains(parsedXML, "<product>") {
		result["type"] = "product_feed"
	}
	if strings.Contains(parsedXML, "<name>") {
		result["has_products"] = true
	}

	return result, nil
}

// ParseUserConfig parses user configuration XML
// VULNERABLE: XXE in business context (user config processing)
func (uc *XmlUseCase) ParseUserConfig(configXML string) (map[string]string, error) {
	// VULNERABLE: Parse XML with external entities
	parsedXML, err := uc.ParseXML(configXML)
	if err != nil {
		return nil, err
	}

	config := make(map[string]string)
	config["raw"] = parsedXML
	config["status"] = "loaded"

	return config, nil
}

// ValidateXML validates XML syntax
// VULNERABLE: XXE via xmllint validation
func (uc *XmlUseCase) ValidateXML(xmlContent string) (bool, string, error) {
	// Write XML to temp file
	tmpFile, err := os.CreateTemp(uc.tempDir, "xml-*.xml")
	if err != nil {
		return false, "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(xmlContent)
	if err != nil {
		tmpFile.Close()
		return false, "", fmt.Errorf("failed to write XML: %v", err)
	}
	tmpFile.Close()

	// VULNERABLE: xmllint processes external entities during validation
	cmd := exec.Command("xmllint", "--noent", "--noout", tmpFile.Name())
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false, string(output), err
	}

	return true, "XML is valid", nil
}

// TransformXML transforms XML using XSLT
// VULNERABLE: XXE + XSLT injection via xsltproc
func (uc *XmlUseCase) TransformXML(xmlContent, xsltContent string) (string, error) {
	// Write XML to temp file
	xmlFile, err := os.CreateTemp(uc.tempDir, "xml-*.xml")
	if err != nil {
		return "", fmt.Errorf("failed to create XML temp file: %v", err)
	}
	defer os.Remove(xmlFile.Name())

	_, err = xmlFile.WriteString(xmlContent)
	if err != nil {
		xmlFile.Close()
		return "", fmt.Errorf("failed to write XML: %v", err)
	}
	xmlFile.Close()

	// Write XSLT to temp file
	xsltFile, err := os.CreateTemp(uc.tempDir, "xslt-*.xsl")
	if err != nil {
		return "", fmt.Errorf("failed to create XSLT temp file: %v", err)
	}
	defer os.Remove(xsltFile.Name())

	_, err = xsltFile.WriteString(xsltContent)
	if err != nil {
		xsltFile.Close()
		return "", fmt.Errorf("failed to write XSLT: %v", err)
	}
	xsltFile.Close()

	// VULNERABLE: xsltproc processes external entities in both XML and XSLT
	cmd := exec.Command("xsltproc", xsltFile.Name(), xmlFile.Name())
	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), fmt.Errorf("XSLT transformation error: %v", err)
	}

	return string(output), nil
}

// ParseFromFile parses XML from a file path
// VULNERABLE: XXE + Path Traversal combo
func (uc *XmlUseCase) ParseFromFile(filePath string) (string, error) {
	// VULNERABLE: No path sanitization - allows reading arbitrary files
	xmlContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// VULNERABLE: Parse with XXE
	return uc.ParseXML(string(xmlContent))
}

// ParseFromURL parses XML from a URL
// VULNERABLE: XXE + SSRF combo
func (uc *XmlUseCase) ParseFromURL(url string) (string, error) {
	// VULNERABLE: SSRF - fetches from arbitrary URL
	// VULNERABLE: Command injection via URL parameter
	cmdStr := fmt.Sprintf("curl -s '%s'", url)
	cmd := exec.Command("sh", "-c", cmdStr)
	xmlContent, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %v", err)
	}

	// VULNERABLE: Parse with XXE
	return uc.ParseXML(string(xmlContent))
}

// GenerateReport generates a report from XML data
// VULNERABLE: XXE in report generation
func (uc *XmlUseCase) GenerateReport(reportXML, outputFormat string) (string, error) {
	// Write XML to temp file
	tmpFile, err := os.CreateTemp(uc.tempDir, "report-*.xml")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(reportXML)
	if err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write XML: %v", err)
	}
	tmpFile.Close()

	// VULNERABLE: XXE + Command Injection
	cmdStr := fmt.Sprintf("xmllint --noent --format %s | head -n 100", tmpFile.Name())
	if outputFormat == "compact" {
		cmdStr = fmt.Sprintf("xmllint --noent --noblanks %s", tmpFile.Name())
	}

	cmd := exec.Command("sh", "-c", cmdStr)
	output, err := cmd.CombinedOutput()

	return string(output), err
}

// CleanTempFiles removes old temporary files
func (uc *XmlUseCase) CleanTempFiles() error {
	files, err := filepath.Glob(filepath.Join(uc.tempDir, "xml-*"))
	if err != nil {
		return err
	}

	for _, f := range files {
		os.Remove(f)
	}

	files, err = filepath.Glob(filepath.Join(uc.tempDir, "xslt-*"))
	if err != nil {
		return err
	}

	for _, f := range files {
		os.Remove(f)
	}

	return nil
}
