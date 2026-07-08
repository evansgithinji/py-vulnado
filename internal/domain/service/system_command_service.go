package service

import (
	"fmt"

	"goapp/internal/domain/entity"
	"goapp/internal/domain/repository"
	"goapp/internal/domain/usecase"
)

// SystemCommandService orchestrates system command execution with deep call graph
type SystemCommandService struct {
	policyRepo     *repository.CommandPolicyRepository
	commandBuilder *usecase.CommandBuilder
	shellExecutor  *usecase.ShellExecutor
}

// NewSystemCommandService creates a new SystemCommandService
func NewSystemCommandService(
	policyRepo *repository.CommandPolicyRepository,
	commandBuilder *usecase.CommandBuilder,
	shellExecutor *usecase.ShellExecutor,
) *SystemCommandService {
	return &SystemCommandService{
		policyRepo:     policyRepo,
		commandBuilder: commandBuilder,
		shellExecutor:  shellExecutor,
	}
}

// ExecuteCommand executes a system command through the deep call graph
// Deep call graph for Command Injection (depth 5):
// Handler -> SystemCommandService.ExecuteCommand()
//   -> CommandPolicyRepository.GetPolicy() -> CommandPolicy
//   -> CommandBuilder.BuildCommand() -> command string (VULNERABLE)
//   -> ShellExecutor.Execute() -> output (VULNERABLE)
func (s *SystemCommandService) ExecuteCommand(request *entity.CommandRequest) (string, error) {
	policy := s.policyRepo.GetPolicy(request.CommandType)
	command := s.commandBuilder.BuildCommand(request.Target, policy)
	return s.shellExecutor.Execute(command)
}

// ExecuteBackup executes a backup command through the deep call graph
func (s *SystemCommandService) ExecuteBackup(name, backupsDir string) (string, error) {
	_ = s.policyRepo.GetPolicy("backup")
	command := s.commandBuilder.BuildBackupCommand(name, backupsDir)
	_, err := s.shellExecutor.Execute(command)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Backup '%s' created successfully", name), nil
}

// ExecuteFileInfo executes a file info command through the deep call graph
func (s *SystemCommandService) ExecuteFileInfo(filesDir, filename string) (string, error) {
	_ = s.policyRepo.GetPolicy("diagnostic")
	command := s.commandBuilder.BuildFileInfoCommand(filesDir, filename)
	return s.shellExecutor.Execute(command)
}
