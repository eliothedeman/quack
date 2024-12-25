package quack

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
)

type node struct {
	name     string
	desc     *funcDesc
	children []*node
}

func isOptionOrFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

// find the terminal command. The return value is a node if found, and the arguments that
// haven't been consumed yet.
func (n *node) find(args []string) (*node, []string) {
	if len(args) > 0 {
		if isOptionOrFlag(args[0]) {
			return nil, nil
		}
		if args[0] == n.name {
			return n, args[1:]
		}
	}
	return nil, nil
}

type Application struct {
	name string
	err  error
	root node
}

func (a *Application) Run() error {
	return a.RunArgs(os.Args)
}

func (a *Application) RunArgs(string []string) error {
	return nil
}

func App(name string, commands ...Fn) *Application {
	b := &Application{name: name}
	for _, c := range commands {
		err := c(&b.root)
		if err != nil {
			b.err = err
			break
		}
	}
	return b
}

type Fn func(*node) error

func Cmd[T any](name string, fn func(args T), subcommands ...Fn) Fn {
	return func(parent *node) error {
		def, err := extract(name, fn)
		if err != nil {
			return err
		}
		this := &node{
			name:     name,
			desc:     def,
			children: nil,
		}
		parent.children = append(parent.children, this)
		for _, cmd := range subcommands {
			err = cmd(this)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

var (
	ErrInvalidFunction = errors.New("invalid function")
	ErrVariadic        = fmt.Errorf("variadic function not supported %w", ErrInvalidFunction)
	ErrTooManyParams   = fmt.Errorf("only a single parameter is supported %w", ErrInvalidFunction)
	ErrNotStruct       = fmt.Errorf("commands must take a single structure parameter")
)

type params interface {
	parse(string) error
}

type option interface {
	long(string)
	short(rune)
	fieldName(string)
}

type baseDef[T V] struct {
	longName      string
	shortName     string
	help          string
	defaultString string
	present       bool
	val           T
}

func (b *baseDef[T]) parseTags(t *reflect.StructField) {
	b.longName = t.Tag.Get("long")
	if b.longName == "" {
		b.longName = t.Name
	}
	b.shortName = t.Tag.Get("short")
	b.help = t.Tag.Get("help")
	b.defaultString = t.Tag.Get("default")
}

func (b *baseDef[T]) Val() T {
	return b.val
}

func (b *baseDef[T]) ptr() any {
	return &b.val
}

func (b *baseDef[T]) long() string {
	return b.longName
}

func (b *baseDef[T]) short() string {
	return b.shortName
}

var _ baseInner = new(baseDef[int])

type baseInner interface {
	ptr() any
	long() string
	short() string
	parseTags(t *reflect.StructField)
}

type Opt[T V] struct {
	baseDef[T]
}

func (Opt[T]) opt() {}

type opt interface {
	baseInner
	opt()
}

type Arg[T V] struct {
	baseDef[T]
}

func (a Arg[T]) arg() {}

type arg interface {
	baseInner
	arg()
}
type Flag struct {
	baseDef[bool]
}

func (Flag) flag() {}

type flag interface {
	flag()
}

type Req[T any] struct {
	v T
}

func (r Req[T]) required() any {
	return r.v
}

type required interface {
	required() any
}

type funcDesc struct {
	fn       reflect.Value
	fnArgPtr reflect.Value
	flagSet  *pflag.FlagSet
	args     []arg
	opts     []opt
	flags    []flag
	reqArgs  []arg
	reqOpts  []opt
	reqFlags []flag
}

func (f *funcDesc) call(args []string) {
	f.fn.Call([]reflect.Value{f.fnArgPtr})
}

func extract[T any](name string, fn func(T)) (*funcDesc, error) {
	t := new(T)
	v := reflect.ValueOf(t)
	tp := v.Elem().Type()
	if tp.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	out := &funcDesc{fn: reflect.ValueOf(fn), fnArgPtr: v, flagSet: pflag.NewFlagSet(name, pflag.ContinueOnError)}
	for i := 0; i < v.Elem().NumField(); i++ {
		f := v.Elem().Field(i)
		ft := v.Elem().Type().Field(i)

		if base, ok := f.Addr().Interface().(baseInner); ok {
			base.parseTags(&ft)
		} else {
			continue
		}

		switch field := f.Addr().Interface().(type) {
		case opt:
			out.opts = append(out.opts, field)
		case arg:
			out.args = append(out.args, field)
		case flag:
			out.flags = append(out.flags, field)
		default:
			log.Panicf("Bad type %T", field)
		}
	}

	return out, nil
}
