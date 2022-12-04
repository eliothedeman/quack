package quack

import (
	"errors"
	"fmt"
)

var (
	// ErrWrongType is returned when a any does not have the the correct underlying type.
	ErrWrongType = errors.New("wrong type")
)

// validate that a unit is either a command or group. If Group this will be done recursivly.
func validateany(u any) error {
	switch u := u.(type) {
	case Command:
	case SimpleCommand:
	case Group:
		for _, v := range u.SubCommands() {
			err := validateany(v)
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

// SimpleCommand is a command that doesn't care about the raw args
type SimpleCommand interface {
	Run()
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
