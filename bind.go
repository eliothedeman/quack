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
	helpTag    = "help"
	defaultTag = "default"
	shortTag   = "short"
	longTag    = "long"
	ignoreTag  = "ignore"
)

type option struct {
	Name    string
	Target  reflect.Value
	Help    string
	Default string
	Short   string
	Long    string
	Ignore  bool
}

func (o *option) fmtBuffer(w io.Writer) {
	fmt.Fprintf(
		w,
		" (%s target:%+v help:%s default:%s short:%s long:%s ignore:%t)",
		o.Name,
		o.Target.Type(),
		o.Help,
		o.Default,
		o.Short,
		o.Long,
		o.Ignore,
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

	return opt
}

// node is a tree structure of commands that binds a structure to an abstract tree for cli applications.
type node struct {
	name        string
	long        string
	short       string
	run         func(*cobra.Command, []string)
	options     []option
	subcommands []*node
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
	cmd.Run = c.run
	return cmd
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
		c.options = append(c.options, opt)

		return false
	})

	return nil
}

// Bind a structure to a *cobra.Command (and sub-commands)
func Bind(name string, root any) (*cobra.Command, error) {
	rn := new(node)
	err := rn.fromStruct(name, root)
	if err != nil {
		return nil, err
	}
	return rn.toCobra(), err
}

// MustBind will panic if Bind returns an error
func MustBind(name string, root any) *cobra.Command {
	cmd, err := Bind(name, root)
	if err != nil {
		panic(err)
	}
	return cmd
}
