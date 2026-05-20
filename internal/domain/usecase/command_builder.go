package usecase

import (
	"fmt"

	"goapp/internal/domain/entity"
)

// CommandBuilder builds shell commands from policy and target
type CommandBuilder struct{}

// NewCommandBuilder creates a new CommandBuilder
func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{}
}

// BuildCommand builds a shell command based on the policy and target
// VULNERABLE: String concatenation for shell command building
func (b *CommandBuilder) BuildCommand(target string, policy *entity.CommandPolicy) string {
	switch policy.Name {
	case "ping":
		return fmt.Sprintf("ping -c 1 %s", target)
	case "backup":
		return fmt.Sprintf("tar czf %s", target)
	case "diagnostic":
		return fmt.Sprintf("echo 'Running diagnostic: %s'", target)
	case "port_check":
		return fmt.Sprintf("nc -zv -w 3 %s 2>&1 || true", target)
	default:
		return fmt.Sprintf("echo %s", target)
	}
}

// BuildBackupCommand builds a backup shell command
// VULNERABLE: String concatenation
func (b *CommandBuilder) BuildBackupCommand(name, backupsDir string) string {
	return fmt.Sprintf("tar czf %s/%s.tar.gz -C /app/data .", backupsDir, name)
}

// BuildFileInfoCommand builds a file info shell command
// VULNERABLE: String concatenation
func (b *CommandBuilder) BuildFileInfoCommand(filesDir, filename string) string {
	return fmt.Sprintf("file %s/%s", filesDir, filename)
}
