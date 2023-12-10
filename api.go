package quack

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	// ErrInvalidType will be returned when a type that doesn't implement the correct interfaces is passed to a bidning function
	ErrInvalidType = errors.New("invlaid type")
	// ErrNotACommand will be returned when a binding target doesn't implement one of the command interfaces.
	ErrNotACommand = errors.New("not a command")
)

// Command is a runnable command that doesn't have sub commands
type Command interface {
	Run([]string)
}

// SimpleCommand is a command that doesn't care about the raw args
type SimpleCommand interface {
	Run()
}

// CobraCommand is the a comand that implements the cobra.Command.Run interface.
// This is useful when you need lower level access to things like global options or the raw cli args.
type CobraCommand interface {
	Run(cmd *cobra.Command, args []string)
}

// Map to cut down on repetitive use of map[string]any
type Map = map[string]any

// Group is a set of subcommands or sub groups.
type Group interface {
	SubCommands() Map
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

// ShortHelper implements a less verbose help.
type ShortHelper interface {
	ShortHelp() string
}
