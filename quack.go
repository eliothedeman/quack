package quack

import (
	"errors"
	"fmt"
)

var (
	// ErrWrongType is returned when a Unit does not have the the correct underlying type.
	ErrWrongType = errors.New("wrong type")
)

// Unit is a placeholder for commands and groups.
type Unit interface{}

// validate that a unit is either a command or group. If Group this will be done recursivly.
func validateUnit(u Unit) error {
	switch u := u.(type) {
	case Command:
		break
	case Group:
		for _, v := range u.SubCommands() {
			err := validateUnit(v)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("%T is neither Command or Group: %w", u, ErrWrongType)
	}
	return nil
}

// Command is a runnable command that doesn't have sub commands
type Command interface {
	Run([]string)
}

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
