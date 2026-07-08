package usecase

import "os/exec"

// ShellExecutor executes shell commands
type ShellExecutor struct{}

// NewShellExecutor creates a new ShellExecutor
func NewShellExecutor() *ShellExecutor {
	return &ShellExecutor{}
}

// Execute executes a shell command and returns the output
// VULNERABLE: Executes shell command with sh -c (shell=true equivalent)
func (e *ShellExecutor) Execute(command string) (string, error) {
	output, err := exec.Command("sh", "-c", command).CombinedOutput()
	return string(output), err
}
