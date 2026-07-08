package usecase

import (
	"fmt"
	"os"
	"os/exec"
)

// XmlDocumentProcessor processes XML documents using xmllint
type XmlDocumentProcessor struct {
	tempDir string
}

// NewXmlDocumentProcessor creates a new XmlDocumentProcessor
func NewXmlDocumentProcessor(tempDir string) *XmlDocumentProcessor {
	os.MkdirAll(tempDir, 0755)
	return &XmlDocumentProcessor{tempDir: tempDir}
}

// Parse parses XML content using xmllint with the given arguments
// VULNERABLE: Parses XML with entity-resolving parser (XXE)
func (p *XmlDocumentProcessor) Parse(parserArgs []string, xmlContent string) (string, error) {
	tmpFile, err := os.CreateTemp(p.tempDir, "xml-*.xml")
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

	args := append(parserArgs, tmpFile.Name())
	cmd := exec.Command("xmllint", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), fmt.Errorf("XML parsing error: %v", err)
	}
	return string(output), nil
}

// Validate validates XML content using xmllint
// VULNERABLE: Validates XML with entity-resolving parser (XXE)
func (p *XmlDocumentProcessor) Validate(parserArgs []string, xmlContent string) (string, error) {
	tmpFile, err := os.CreateTemp(p.tempDir, "xml-*.xml")
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

	args := append(parserArgs, "--noout", tmpFile.Name())
	cmd := exec.Command("xmllint", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return string(output), err
	}
	return "XML is valid", nil
}
