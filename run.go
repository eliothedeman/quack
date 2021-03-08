package quack

import (
	"fmt"
	"os"
)

func handleRunError(err error) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}

type options struct {
	Unit
	colorEnabled bool
	args         []string
}

// Option allows the user to customize the running of commands and groups
type Option func(*options)

// WithArgs will override the args used by quack instead of os.Args
func WithArgs(args []string) Option {
	return func(o *options) {
		o.args = args
	}
}

// WithCommand sets the command to run. If called again, or after WithGroup, the last call will be respected.
func WithCommand(c Command) Option {
	return func(o *options) {
		o.Unit = c
	}
}

// WithGroup sets the group to run. If called again, or after WithGroup, the last call will be respected.
func WithGroup(g Group) Option {
	return func(o *options) {
		o.Unit = g
	}
}

// Run with the given options
func Run(name string, opts ...Option) {
	config := options{
		args: os.Args[1:],
	}
	for _, o := range opts {
		o(&config)
	}

	handleRunError(config.run(name))
}

// run a command or find a subcommand
func (o *options) run(name string) error {
	u := o.Unit
	raw := o.args
	err := validateUnit(u)
	if err != nil {
		return err
	}

	switch u := u.(type) {
	case Command:

		if d, ok := u.(Defaulter); ok {
			d.Default()
		}
		fs := getFlags(name, u)

		if hasHelpArg(raw, fs.ShorthandLookup("h") == nil) {
			printHelp(name, u)
			return nil
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

		if hasHelpArg(raw, len(raw) == 1) || len(raw) == 0 {
			printHelp(name, u)
			return nil
		}
		subs := u.SubCommands()
		s, ok := subs[raw[0]]
		if !ok {
			return fmt.Errorf("unable to find subcommand %s\n%s", raw[0], fmtHelp(name, u))
		}
		child := options{
			Unit:         s,
			args:         raw[1:],
			colorEnabled: o.colorEnabled,
		}
		return child.run(raw[0])
	default:
		return fmt.Errorf("unknown type %T is not Command or Group", u)
	}
	return nil
}
