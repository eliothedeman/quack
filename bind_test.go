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

type positionalCmd struct {
	Source string `position:"1"`
	Target string `position:"2"`
}

func (p *positionalCmd) Run(*cobra.Command, []string) {
}

type repeatedFlagCmd struct {
	Files []string `repeated:""`
}

func (r *repeatedFlagCmd) Run(*cobra.Command, []string) {
}

type repeatedPositionalCmd struct {
	Files []string `position:"1" repeated:""`
}

func (r *repeatedPositionalCmd) Run(*cobra.Command, []string) {
}

func TestBindCobra(t *testing.T) {
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
			cmd, err := BindCobra(test.name, test.in)
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

type outOfOrderPositionalCmd struct {
	Target string `position:"2"`
	Source string `position:"1"`
}

func (o *outOfOrderPositionalCmd) Run(*cobra.Command, []string) {
}

func TestPositionalArgs(t *testing.T) {
	t.Run("basic_positional", func(t *testing.T) {
		cmd := new(positionalCmd)
		cobraCmd, err := BindCobra("copy", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, cobraCmd)

		// Simulate running the command with positional args
		cobraCmd.SetArgs([]string{"file1.txt", "file2.txt"})
		err = cobraCmd.Execute()
		assert.Nil(t, err)
		assert.Equal(t, "file1.txt", cmd.Source)
		assert.Equal(t, "file2.txt", cmd.Target)
	})

	t.Run("out_of_order_positional", func(t *testing.T) {
		cmd := new(outOfOrderPositionalCmd)
		cobraCmd, err := BindCobra("copy", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, cobraCmd)

		// Verify that positions are respected, not field order
		cobraCmd.SetArgs([]string{"source.txt", "target.txt"})
		err = cobraCmd.Execute()
		assert.Nil(t, err)
		assert.Equal(t, "source.txt", cmd.Source) // position 1
		assert.Equal(t, "target.txt", cmd.Target) // position 2
	})

	t.Run("repeated_positional", func(t *testing.T) {
		cmd := new(repeatedPositionalCmd)
		cobraCmd, err := BindCobra("list", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, cobraCmd)

		// Simulate running the command with multiple positional args
		cobraCmd.SetArgs([]string{"file1.txt", "file2.txt", "file3.txt"})
		err = cobraCmd.Execute()
		assert.Nil(t, err)
		assert.Equal(t, []string{"file1.txt", "file2.txt", "file3.txt"}, cmd.Files)
	})
}

func TestRepeatedFlags(t *testing.T) {
	t.Run("repeated_flag", func(t *testing.T) {
		cmd := new(repeatedFlagCmd)
		cobraCmd, err := BindCobra("process", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, cobraCmd)

		// Simulate running the command with repeated flags
		cobraCmd.SetArgs([]string{"--files", "file1.txt", "--files", "file2.txt"})
		err = cobraCmd.Execute()
		assert.Nil(t, err)
		assert.Equal(t, []string{"file1.txt", "file2.txt"}, cmd.Files)
	})
}
