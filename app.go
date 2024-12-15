package quack

import (
	"errors"
	"fmt"
	"reflect"
)

/*
	type node struct {
		name     string
		fn       any
		children []*node
	}

	type AppBuilder struct {
		node
	}

type Option func(*node)

	func App(name string) *AppBuilder {
		return &AppBuilder{node{name: name}}
	}

	type CmdBuilder struct {
		node
	}

	func Cmd[T any](name string, fn func(T)) Option {
		return func(n *node) {
			n.children = append(n.children, &node{
				name: name,
				fn:   fn,
			})
		}
	}
*/

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
	longNames  []string
	shortNames []rune
	present    bool
	v          T
}

func (b *baseDef[T]) long(l string) {
	b.longNames = append(b.longNames, l)
}

func (b *baseDef[T]) short(l rune) {
	b.shortNames = append(b.shortNames, l)
}

type Opt[T V] struct {
	baseDef[T]
}

func (Opt[T]) opt() {}

type opt interface {
	opt()
}

type Arg[T V] struct {
	baseDef[T]
}

func (a Arg[T]) arg() {}

type arg interface {
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
	fn       any
	input    any
	args     []arg
	opts     []opt
	flags    []flag
	reqArgs  []arg
	reqOpts  []opt
	reqFlags []flag
}

func extract[T any](fn func(T)) (*funcDesc, error) {
	t := new(T)
	v := reflect.ValueOf(t).Elem()
	tp := v.Type()
	if tp.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}
	out := &funcDesc{fn: fn, input: t}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		switch field := f.Addr().Interface().(type) {
		case required:
			switch field := field.required().(type) {
			case opt:
				out.reqOpts = append(out.reqOpts, field)
			case arg:
				out.reqArgs = append(out.reqArgs, field)
			case flag:
				out.reqFlags = append(out.reqFlags, field)
			}
		case opt:
			out.opts = append(out.opts, field)
		case arg:
			out.args = append(out.args, field)
		case flag:
			out.flags = append(out.flags, field)
		}
	}

	return out, nil
}
