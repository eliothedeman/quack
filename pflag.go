package quack

import (
	"fmt"
	"strconv"
)

/*
func toFlag(field baseInner, fs *pflag.FlagSet) {
	switch v := field.ptr().(type) {
	case *int:
		fs.IntVarP(v, field.long(), string(field.short()), 0, "")
	default:
		panic(v)
	}
}

func nodeToFlagSet(n *node) *pflag.FlagSet {
	fs := pflag.NewFlagSet(n.name, pflag.ContinueOnError)
	for _, o := range n.desc.opts {
		toFlag(o, fs)
	}
	return fs
}
*/

type argument interface {
	parseFromArgs(args []string) error
}

type Value interface {
	Parse(string) error
	Format() string
}

type integerValue[T int64] struct {
	val T
}

func (i *integerValue[T]) Parse(s string) error {
	x, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	i.val = int64(x)

	return nil
}

type inner[T V] struct {
	longName      string
	shortName     string
	usage         string
	defaultString string
	present       bool
	val           T
}

type Opt[T V] struct {
	inner[T]
}

type (
	Int    int64
	Float  float64
	String string
)

func (i *Int) Parse(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	*i = Int(v)
	return nil
}

func (i *Int) String() string {
	return fmt.Sprintf("%d", i)
}

/*

func (f *Float) Parse(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*f = Float(v)
	return nil
}

func (s *String) Parse(x string) error {
	*s = String(x)
	return nil
}
*/
