package quack

import (
	"github.com/spf13/pflag"
)

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
