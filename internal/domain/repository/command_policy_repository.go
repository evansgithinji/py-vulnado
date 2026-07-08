package repository

import "goapp/internal/domain/entity"

// CommandPolicyRepository stores command execution policies
type CommandPolicyRepository struct {
	policies map[string]*entity.CommandPolicy
}

// NewCommandPolicyRepository creates a new CommandPolicyRepository
func NewCommandPolicyRepository() *CommandPolicyRepository {
	return &CommandPolicyRepository{
		policies: map[string]*entity.CommandPolicy{
			"ping": {
				Name:            "ping",
				AllowedCommands: []string{"ping"},
				Timeout:         10,
				AllowShell:      true, // VULNERABLE
			},
			"backup": {
				Name:            "backup",
				AllowedCommands: []string{"tar"},
				Timeout:         60,
				AllowShell:      true, // VULNERABLE
			},
			"diagnostic": {
				Name:            "diagnostic",
				AllowedCommands: []string{"echo"},
				Timeout:         30,
				AllowShell:      true, // VULNERABLE
			},
			"port_check": {
				Name:            "port_check",
				AllowedCommands: []string{"nc"},
				Timeout:         10,
				AllowShell:      true, // VULNERABLE
			},
		},
	}
}

// GetPolicy retrieves a command policy by name
func (r *CommandPolicyRepository) GetPolicy(name string) *entity.CommandPolicy {
	if policy, ok := r.policies[name]; ok {
		return policy
	}
	return r.policies["ping"]
}
