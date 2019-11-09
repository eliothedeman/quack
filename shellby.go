package quack

import (
	"fmt"
	"os"
)

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

func handleRunError(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}

// RunCommand runs a simple command with no subcommands.
func RunCommand(name string, c Command) {
	args := os.Args
	if len(args) > 0 {
		args = args[1:]
	}
	handleRunError(run(name, c, args))
}

// RunGroup runs a group of commands
func RunGroup(name string, g Group) {
	args := os.Args
	if len(args) > 0 {
		args = args[1:]
	}
	handleRunError(run(name, g, args))
}

// run a command or find a subcommand
func run(name string, u Unit, raw []string) error {

	switch u := u.(type) {
	case Command:
		if d, ok := u.(Defaulter); ok {
			d.Default()
		}
		fs := getFlags(name, u)

		if hasHelpArg(raw, fs.ShorthandLookup("h") == nil) {
			return helpError(name, u)
		}

		err := fs.Parse(raw)

		if err != nil {
			return fmt.Errorf("parsing error: %w", err)
		}

		if v, ok := u.(Validator); ok {
			err := v.Validate()
			if err != nil {
				return fmt.Errorf("validation error: %w", err)
			}
		}
		u.Run(raw)
	case Group:

		if hasHelpArg(raw, false) || len(raw) == 0 {
			return helpError(name, u)
		}
		subs := u.SubCommands()
		s, ok := subs[raw[0]]
		if !ok {
			return fmt.Errorf("unable to find subcommand %s\n%s", raw[0], fmtHelp(name, u))
		}
		return run(raw[0], s, raw[1:])
	default:
		return fmt.Errorf("unknown type %T is not Command or Group", u)
	}
	return nil
}

func keys(m map[string]Unit) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
