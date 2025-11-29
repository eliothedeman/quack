package quack

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

var ()

const (
	helpTag       = "help"
	defaultTag    = "default"
	shortTag      = "short"
	longTag       = "long"
	ignoreTag     = "ignore"
	positionalTag = "positional"
	repeatedTag   = "repeated"
)

type option struct {
	Name       string
	Target     reflect.Value
	Help       string
	Default    string
	Short      string
	Long       string
	Ignore     bool
	Positional bool
	Repeated   bool
}

func (o *option) fmtBuffer(w io.Writer) {
	fmt.Fprintf(
		w,
		" (%s target:%+v help:%s default:%s short:%s long:%s ignore:%t positional:%t repeated:%t)",
		o.Name,
		o.Target.Type(),
		o.Help,
		o.Default,
		o.Short,
		o.Long,
		o.Ignore,
		o.Positional,
		o.Repeated,
	)
}

func optionFromField(field reflect.StructField) option {
	opt := option{
		Name: fieldNameToArg(field.Name),
	}
	tags := field.Tag
	opt.Help = tags.Get(helpTag)
	opt.Short = tags.Get(shortTag)
	opt.Long = tags.Get(longTag)
	opt.Default = tags.Get(defaultTag)
	_, opt.Ignore = tags.Lookup(ignoreTag)
	_, opt.Positional = tags.Lookup(positionalTag)
	_, opt.Repeated = tags.Lookup(repeatedTag)

	return opt
}

// node is a tree structure of commands that binds a structure to an abstract tree for cli applications.
type node struct {
	name              string
	long              string
	short             string
	run               func(*cobra.Command, []string)
	options           []option
	positionalOptions []option
	subcommands       []*node
}

func (c *node) toCobra() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.name,
		Long:  c.long,
		Short: c.short,
	}
	flags := cmd.Flags()
	for _, o := range c.options {
		o.setFlag(flags)
	}

	for _, s := range c.subcommands {
		cmd.AddCommand(s.toCobra())
	}

	// Wrap the run function to handle positional arguments
	if c.run != nil {
		originalRun := c.run
		cmd.Run = func(cobraCmd *cobra.Command, args []string) {
			// Parse positional arguments
			if err := c.parsePositionalArgs(args); err != nil {
				cobraCmd.PrintErr(err)
				return
			}
			originalRun(cobraCmd, args)
		}
	}

	return cmd
}

// parsePositionalArgs parses positional arguments and assigns them to the appropriate fields
func (c *node) parsePositionalArgs(args []string) error {
	argIndex := 0
	for _, opt := range c.positionalOptions {
		if argIndex >= len(args) {
			// Not enough arguments provided
			if opt.Default == "" {
				return fmt.Errorf("missing required positional argument: %s", opt.Name)
			}
			// Use default value
			if err := opt.parseValue(opt.Default); err != nil {
				return fmt.Errorf("failed to parse default value for %s: %w", opt.Name, err)
			}
			continue
		}

		if opt.Repeated {
			// Consume all remaining arguments
			for argIndex < len(args) {
				if err := opt.appendValue(args[argIndex]); err != nil {
					return fmt.Errorf("failed to parse positional argument %s: %w", opt.Name, err)
				}
				argIndex++
			}
		} else {
			// Single positional argument
			if err := opt.parseValue(args[argIndex]); err != nil {
				return fmt.Errorf("failed to parse positional argument %s: %w", opt.Name, err)
			}
			argIndex++
		}
	}
	return nil
}

func (c *node) fmtBuffer(depth int, w io.Writer) {
	padding := strings.Repeat(" ", depth)
	fmt.Fprintf(w, "\n%s(%s", padding, c.name)
	for _, o := range c.options {
		fmt.Fprintf(w, "\n%s", padding)
		o.fmtBuffer(w)
	}

	for _, s := range c.subcommands {
		s.fmtBuffer(depth+1, w)
	}
	fmt.Fprintf(w, "\n%s)", strings.Repeat(" ", depth))
}

type mapWrapper struct {
	m Map
}

func (m *mapWrapper) SubCommands() Map {
	return m.m
}

func (m *mapWrapper) Run(cmd *cobra.Command) {
	cmd.Help()
}

func (c *node) fromStruct(name string, target any) error {
	if m, ok := target.(Map); ok {
		target = &mapWrapper{m}
	}
	v := reflect.Indirect(reflect.ValueOf(target))
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return fmt.Errorf(
			"%w only structures can be commands got %T",
			ErrInvalidType,
			target,
		)
	}
	c.name = name

	if helper, ok := target.(Helper); ok {
		c.long = helper.Help()
	}
	c.short = c.long
	if sh, ok := target.(ShortHelper); ok {
		c.short = sh.ShortHelp()
	}

	switch target := target.(type) {
	case Command:
		c.run = func(c *cobra.Command, s []string) {
			target.Run(s)
		}
	case SimpleCommand:
		c.run = func(*cobra.Command, []string) {
			target.Run()
		}
	case CobraCommand:
		c.run = target.Run
	case Group:
		c.run = func(c *cobra.Command, s []string) {
			c.Help()
		}
		for name, s := range target.SubCommands() {
			cn := new(node)
			if err := cn.fromStruct(name, s); err != nil {
				return err
			}
			c.subcommands = append(c.subcommands, cn)
		}
	default:
		return fmt.Errorf("%w. must impliment quack.(Command|SimpleCommand|Group|SubCommander)", ErrNotACommand)
	}

	v.FieldByNameFunc(func(name string) bool {
		f := v.FieldByName(name)
		sf, ok := t.FieldByName(name)
		if !ok {
			panic("wtf")
		}
		// embedded structs will be iterated over later
		if sf.Anonymous || !sf.IsExported() {
			return false
		}
		opt := optionFromField(sf)
		opt.Target = f

		// Separate positional and regular options
		if opt.Positional {
			c.positionalOptions = append(c.positionalOptions, opt)
		} else {
			c.options = append(c.options, opt)
		}

		return false
	})

	return nil
}

// BindCobra a structure to a *cobra.Command (and sub-commands)
func BindCobra(name string, root any) (*cobra.Command, error) {
	rn := new(node)
	err := rn.fromStruct(name, root)
	if err != nil {
		return nil, err
	}
	return rn.toCobra(), err
}

// MustBindCobra will panic if BindCobra returns an error
func MustBindCobra(name string, root any) *cobra.Command {
	cmd, err := BindCobra(name, root)
	if err != nil {
		panic(err)
	}
	return cmd
}
