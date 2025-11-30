package quack

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var ()

const (
	helpTag     = "help"
	defaultTag  = "default"
	shortTag    = "short"
	longTag     = "long"
	ignoreTag   = "ignore"
	argTag      = "arg"
	repeatedTag = "repeated"
)

type option struct {
	Name     string
	Target   reflect.Value
	Help     string
	Default  string
	Short    string
	Long     string
	Ignore   bool
	Arg      int  // 0 means not a positional arg, >0 means positional argument at that index
	Repeated bool
}

func (o *option) fmtBuffer(w io.Writer) {
	fmt.Fprintf(
		w,
		" (%s target:%+v help:%s default:%s short:%s long:%s ignore:%t arg:%d repeated:%t)",
		o.Name,
		o.Target.Type(),
		o.Help,
		o.Default,
		o.Short,
		o.Long,
		o.Ignore,
		o.Arg,
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
	_, opt.Repeated = tags.Lookup(repeatedTag)

	// Parse arg tag
	if argStr := tags.Get(argTag); argStr != "" {
		arg, err := strconv.Atoi(argStr)
		if err != nil || arg < 1 {
			panic(fmt.Sprintf("invalid arg value '%s' for field %s: must be a positive integer", argStr, field.Name))
		}
		opt.Arg = arg
	}

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
	target            any // Store the original target for framework-specific handling
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

		// Check if this is a slice type (repeated argument)
		if opt.Target.Kind() == reflect.Slice {
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
	c.target = target // Store the target for framework-specific handling

	if helper, ok := target.(Helper); ok {
		c.long = helper.Help()
	}
	c.short = c.long
	if sh, ok := target.(ShortHelper); ok {
		c.short = sh.ShortHelp()
	}

	// Check if target has a Run method that might be UrfaveCommand
	// We check this using reflection to avoid import dependencies
	hasUrfaveRun := false
	targetValue := reflect.ValueOf(target)
	runMethod := targetValue.MethodByName("Run")
	if runMethod.IsValid() {
		methodType := runMethod.Type()
		// Check if it matches the UrfaveCommand signature: Run(ctx *Something) error
		// NumIn() == 1 means one parameter (the receiver is not counted for method values)
		// NumOut() == 1 means one return value
		if methodType.NumIn() == 1 && methodType.NumOut() == 1 {
			// Check if the return type is error
			errorType := reflect.TypeOf((*error)(nil)).Elem()
			if methodType.Out(0).Implements(errorType) {
				hasUrfaveRun = true
			}
		}
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
	case UrfaveCommand:
		// UrfaveCommand is handled differently in urfave binding
		// We set a placeholder run function here
		c.run = func(*cobra.Command, []string) {
			// This will be overridden in toUrfaveApp/toUrfaveCommand
		}
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
		// Check if it implements UrfaveCommand pattern via reflection
		if hasUrfaveRun {
			c.run = func(*cobra.Command, []string) {
				// This will be overridden in toUrfaveApp/toUrfaveCommand
			}
		} else {
			return fmt.Errorf("%w. must impliment quack.(Command|SimpleCommand|Group|SubCommander)", ErrNotACommand)
		}
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

		// Separate positional args and named options
		if opt.Arg > 0 {
			c.positionalOptions = append(c.positionalOptions, opt)
		} else {
			c.options = append(c.options, opt)
		}

		return false
	})

	// Sort positional args by their arg position
	sort.Slice(c.positionalOptions, func(i, j int) bool {
		return c.positionalOptions[i].Arg < c.positionalOptions[j].Arg
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
