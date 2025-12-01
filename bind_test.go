package quack

import (
	"fmt"
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
	Source string `arg:"1"`
	Target string `arg:"2"`
}

func (p *positionalCmd) Run(*cobra.Command, []string) {
}

type repeatedFlagCmd struct {
	Files []string
}

func (r *repeatedFlagCmd) Run(*cobra.Command, []string) {
}

type repeatedPositionalCmd struct {
	Files []string `arg:"1"`
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
	Target string `arg:"2"`
	Source string `arg:"1"`
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

		// Verify that arg positions are respected, not field order
		cobraCmd.SetArgs([]string{"source.txt", "target.txt"})
		err = cobraCmd.Execute()
		assert.Nil(t, err)
		assert.Equal(t, "source.txt", cmd.Source) // arg 1
		assert.Equal(t, "target.txt", cmd.Target) // arg 2
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

// Test option validation

// validatingOption is a custom type that implements Validator
type validatingOption string

func (v validatingOption) Validate() error {
	if v == "invalid" {
		return fmt.Errorf("option value is invalid")
	}
	return nil
}

// cmdWithValidatingOption has an option that implements Validator but the command doesn't
type cmdWithValidatingOption struct {
	Option validatingOption
}

func (c *cmdWithValidatingOption) Run(*cobra.Command, []string) {
}

// cmdWithValidatingOptionAndCommandValidator has both option and command validation
type cmdWithValidatingOptionAndCommandValidator struct {
	Option validatingOption
}

func (c *cmdWithValidatingOptionAndCommandValidator) Run(*cobra.Command, []string) {
}

func (c *cmdWithValidatingOptionAndCommandValidator) Validate() error {
	// Command level validation - options should not be validated individually
	return fmt.Errorf("command validation error")
}

func TestOptionValidation(t *testing.T) {
	t.Run("valid_option", func(t *testing.T) {
		cmd := &cmdWithValidatingOption{Option: "valid"}
		cobraCmd, err := BindCobra("test", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, cobraCmd)

		cobraCmd.SetArgs([]string{"--option", "valid"})
		err = cobraCmd.Execute()
		assert.Nil(t, err)
	})

	t.Run("invalid_option", func(t *testing.T) {
		cmd := &cmdWithValidatingOption{Option: "valid"}
		cobraCmd, err := BindCobra("test", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, cobraCmd)

		cobraCmd.SetArgs([]string{"--option", "invalid"})
		err = cobraCmd.Execute()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "validation failed for option")
	})

	t.Run("command_implements_validator_skips_option_validation", func(t *testing.T) {
		cmd := &cmdWithValidatingOptionAndCommandValidator{Option: "invalid"}
		cobraCmd, err := BindCobra("test", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, cobraCmd)

		// Even though option is "invalid", it should not be validated
		// because the command itself implements Validator
		cobraCmd.SetArgs([]string{"--option", "invalid"})
		err = cobraCmd.Execute()
		// Should succeed without option validation error
		assert.Nil(t, err)
	})
}
