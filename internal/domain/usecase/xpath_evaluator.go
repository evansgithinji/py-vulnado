package usecase

import (
	"fmt"
	"os"
	"os/exec"

	"goapp/internal/domain/entity"
)

// XPathExpressionBuilder constructs XPath query strings.
type XPathExpressionBuilder struct{}

// NewXPathExpressionBuilder creates a new XPathExpressionBuilder.
func NewXPathExpressionBuilder() *XPathExpressionBuilder {
	return &XPathExpressionBuilder{}
}

// BuildAuthQuery builds an XPath authentication query.
// VULNERABLE: XPath Injection (CWE-643) - string concatenation in XPath expression
func (b *XPathExpressionBuilder) BuildAuthQuery(username, password string) string {
	// VULNERABLE: XPath Injection - string concatenation
	return "/users/user[username='" + username + "' and password='" + password + "']/name"
}

// BuildRawQuery returns the query string as-is for direct evaluation.
// VULNERABLE: XPath Injection (CWE-643) - user-controlled XPath expression
func (b *XPathExpressionBuilder) BuildRawQuery(query string) string {
	// VULNERABLE: Direct user input as XPath expression
	return query
}

// XPathEvaluator executes XPath expressions against XML documents.
type XPathEvaluator struct{}

// NewXPathEvaluator creates a new XPathEvaluator.
func NewXPathEvaluator() *XPathEvaluator {
	return &XPathEvaluator{}
}

// Evaluate writes the XML document to a temp file and runs xmllint --xpath.
func (e *XPathEvaluator) Evaluate(doc *entity.XmlDocument, xpathExpr string) (string, error) {
	tmpFile, err := os.CreateTemp(doc.TempDir, "xpath-*.xml")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(doc.Content)
	if err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write XML: %v", err)
	}
	tmpFile.Close()

	// VULNERABLE: XPath expression passed directly to xmllint
	cmd := exec.Command("xmllint", "--xpath", xpathExpr, tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("XPath query failed: %s", string(output))
	}

	return string(output), nil
}
