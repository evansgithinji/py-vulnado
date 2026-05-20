package entity

// CommandRequest represents a system command request DTO
type CommandRequest struct {
	Target      string
	CommandType string
	Args        map[string]string
}

// CommandPolicy defines the policy for command execution
type CommandPolicy struct {
	Name            string
	AllowedCommands []string
	Timeout         int
	AllowShell      bool // VULNERABLE: allows shell execution
}
