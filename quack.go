package quack

// Unit is a placeholder for commands and groups.
type Unit interface{}

// Command is a runnable command that doesn't have sub commands
type Command interface {
	Run([]string)
}

// Group is a set of subcommands or sub groups.
type Group interface {
	SubCommands() map[string]Unit
}

// Validator is a command or argument that wants to be validated.
type Validator interface {
	Validate() error
}

// Defaulter can set up the default arguments of a command
type Defaulter interface {
	Default()
}

// Parser is an argument that wants to parse itself.
type Parser interface {
	Parse(string) error
}

// Helper returns usage information for a command or group.
type Helper interface {
	Help() string
}
