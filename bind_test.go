package quack

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type simpleCmd struct {
	NoHelp string
	AnInt  int `help:"I am an int"`
}

func (s *simpleCmd) Run(*cobra.Command, []string) {

}

func (s *simpleCmd) Help() string {
	return "longer help message"
}

func sanatize(s string) string {
	remove := []string{
		"\n",
		"\t",
		" ",
	}
	for _, x := range remove {
		s = strings.ReplaceAll(s, x, "")
	}
	return s
}

func TestBind(t *testing.T) {
	simple := new(simpleCmd)
	tests := []struct {
		name  string
		in    any
		usage string
		err   error
	}{
		{
			"bad_type",
			0,
			"",
			ErrInvalidType,
		},
		{
			"not a command",
			struct{}{},
			"",
			ErrNotACommand,
		},
		{
			"simple",
			simple,
			`
Usage:
  simple [flags]

Flags:
 --an-int int       I am an int
 --no-help string
`,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cmd, err := Bind(test.name, test.in)
			if test.err == nil {
				assert.Nil(t, err)
				cmd.UsageString()
				assert.Equal(t,
					sanatize(test.usage), sanatize(cmd.UsageString()))
				return
			}
			assert.ErrorIs(t, err, test.err)
			assert.Nil(t, cmd)
		})

	}
}
