package quack

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/spf13/pflag"
)

func fmtUsage(w io.Writer, fs *pflag.FlagSet) {
	tw := tabwriter.NewWriter(w, 2, 2, 1, ' ', 0)
	fs.SortFlags = true
	var flags []string
	var options []string
	fs.VisitAll(func(f *pflag.Flag) {
		var line strings.Builder

		if f.Shorthand != "" {
			fmt.Fprintf(&line, "\t-%s,\t--%s", f.Shorthand, f.Name)
		} else {
			fmt.Fprintf(&line, "\t\t--%s", f.Name)
		}
		line.WriteByte('\t')

		valType := f.Value.Type()

		typeName, usage := pflag.UnquoteUsage(f)
		if typeName != "" {
			line.WriteString(typeName)
		}

		line.WriteByte('\t')
		if f.DefValue != "" {
			switch valType {
			case "string":
				fmt.Fprintf(&line, "(default='%s')", f.DefValue)
			case "bool":
				fmt.Fprintf(&line, "(default=%s)", f.DefValue)

			default:
				fmt.Fprintf(&line, "(default=%s)", f.DefValue)

			}

		}
		line.WriteByte('\t')
		line.WriteString(usage)

		switch valType {
		case "bool":
			flags = append(flags, line.String())
		default:
			options = append(options, line.String())
		}

	})

	for i, f := range flags {
		if i == 0 {
			fmt.Fprintln(tw, "Flags:\t\t\t\t\t")
		}
		fmt.Fprintln(tw, f)
	}
	for i, o := range options {
		if i == 0 {
			fmt.Fprintln(tw, "Options:\t\t\t\t\t")
		}
		fmt.Fprintln(tw, o)
	}

	tw.Flush()
}
