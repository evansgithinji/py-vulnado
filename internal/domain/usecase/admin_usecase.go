package usecase

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// AdminUseCase handles admin operations
type AdminUseCase struct{}

// NewAdminUseCase creates a new AdminUseCase
func NewAdminUseCase() *AdminUseCase {
	return &AdminUseCase{}
}

// GetDebugInfo returns debug information
// VULNERABLE: Information disclosure - exposes environment variables
func (uc *AdminUseCase) GetDebugInfo() string {
	var result strings.Builder

	result.WriteString("=== Debug Information ===\n")
	result.WriteString(fmt.Sprintf("Database: %s\n", "/app/data/app.db"))
	result.WriteString(fmt.Sprintf("Working Dir: %s\n", os.Getenv("PWD")))
	result.WriteString("\n=== Environment Variables ===\n")

	// VULNERABLE: Exposes all environment variables
	for _, env := range os.Environ() {
		result.WriteString(fmt.Sprintf("%s\n", env))
	}

	return result.String()
}

// GenerateInvoicePDF generates an invoice PDF
// VULNERABLE: Command injection via orderID and format
func (uc *AdminUseCase) GenerateInvoicePDF(orderID, format string) (string, error) {
	// VULNERABLE: Command injection
	cmd := fmt.Sprintf("echo 'Invoice for Order #%s' | wkhtmltopdf - invoice_%s.%s",
		orderID, orderID, format)
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("generation failed: %s - %v", string(output), err)
	}
	return string(output), nil
}

// ExecuteConfig executes a configuration command
// VULNERABLE: Command injection via config parameter
func (uc *AdminUseCase) ExecuteConfig(config string) (string, error) {
	// VULNERABLE: Insecure deserialization/command execution
	if strings.Contains(config, "cmd:") {
		cmd := strings.TrimPrefix(config, "cmd:")
		output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
		if err != nil {
			return "", err
		}
		return string(output), nil
	}
	return "Config applied", nil
}

// BackupDatabase creates a database backup
// VULNERABLE: Command injection via backup name
func (uc *AdminUseCase) BackupDatabase(backupName string) (string, error) {
	// VULNERABLE: Command injection
	cmd := fmt.Sprintf("cp /app/data/app.db /app/backups/%s.db", backupName)
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("backup failed: %s - %v", string(output), err)
	}
	return "Backup created: " + backupName + ".db", nil
}

// RestoreDatabase restores a database backup
// VULNERABLE: Command injection via backup name
func (uc *AdminUseCase) RestoreDatabase(backupName string) (string, error) {
	// VULNERABLE: Command injection
	cmd := fmt.Sprintf("cp /app/backups/%s.db /app/data/app.db", backupName)
	output, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("restore failed: %s - %v", string(output), err)
	}
	return "Database restored from: " + backupName, nil
}
