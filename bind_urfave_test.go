package quack

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

type simpleUrfaveCmd struct {
	NoHelp string
	AnInt  int `help:"I am an int"`
}

func (s *simpleUrfaveCmd) Run([]string) {
}

func (s *simpleUrfaveCmd) Help() string {
	return "longer help message"
}

type positionalUrfaveCmd struct {
	Source string `arg:"1"`
	Target string `arg:"2"`
}

func (p *positionalUrfaveCmd) Run([]string) {
}

type repeatedFlagUrfaveCmd struct {
	Files []string
}

func (r *repeatedFlagUrfaveCmd) Run([]string) {
}

type repeatedPositionalUrfaveCmd struct {
	Files []string `arg:"1"`
}

func (r *repeatedPositionalUrfaveCmd) Run([]string) {
}

type urfaveContextCmd struct {
	Name string
}

func (u *urfaveContextCmd) Run(ctx context.Context, cmd *cli.Command) error {
	// Can access urfave-specific command here
	return nil
}

func TestBindUrfave(t *testing.T) {
	simple := new(simpleUrfaveCmd)
	tests := []struct {
		name string
		in   any
		err  error
	}{
		{
			"bad_type",
			0,
			ErrInvalidType,
		},
		{
			"not a command",
			struct{}{},
			ErrNotACommand,
		},
		{
			"simple",
			simple,
			nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, err := BindUrfave(test.name, test.in)
			if test.err == nil {
				assert.Nil(t, err)
				assert.NotNil(t, app)
				assert.Equal(t, test.name, app.Name)
				return
			}
			assert.ErrorIs(t, err, test.err)
			assert.Nil(t, app)
		})
	}
}

func TestUrfavePositionalArgs(t *testing.T) {
	t.Run("basic_positional", func(t *testing.T) {
		cmd := new(positionalUrfaveCmd)
		app, err := BindUrfave("copy", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, app)

		// Simulate running the command with positional args
		err = app.Run(context.Background(), []string{"copy", "file1.txt", "file2.txt"})
		assert.Nil(t, err)
		assert.Equal(t, "file1.txt", cmd.Source)
		assert.Equal(t, "file2.txt", cmd.Target)
	})

	t.Run("repeated_positional", func(t *testing.T) {
		cmd := new(repeatedPositionalUrfaveCmd)
		app, err := BindUrfave("list", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, app)

		// Simulate running the command with multiple positional args
		err = app.Run(context.Background(), []string{"list", "file1.txt", "file2.txt", "file3.txt"})
		assert.Nil(t, err)
		assert.Equal(t, []string{"file1.txt", "file2.txt", "file3.txt"}, cmd.Files)
	})
}

func TestUrfaveRepeatedFlags(t *testing.T) {
	t.Run("repeated_flag", func(t *testing.T) {
		cmd := new(repeatedFlagUrfaveCmd)
		app, err := BindUrfave("process", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, app)

		// Simulate running the command with repeated flags
		err = app.Run(context.Background(), []string{"process", "--files", "file1.txt", "--files", "file2.txt"})
		assert.Nil(t, err)
		assert.Equal(t, []string{"file1.txt", "file2.txt"}, cmd.Files)
	})
}

func TestUrfaveCommand(t *testing.T) {
	t.Run("urfave_context_command", func(t *testing.T) {
		cmd := new(urfaveContextCmd)
		app, err := BindUrfave("test", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, app)

		// Simulate running the command
		err = app.Run(context.Background(), []string{"test", "--name", "testname"})
		assert.Nil(t, err)
		assert.Equal(t, "testname", cmd.Name)
	})
}

type cmdWithDefaults struct {
	Port int    `default:"8080"`
	Host string `default:"localhost"`
}

func (c *cmdWithDefaults) Run([]string) {
}

func TestUrfaveWithDefaults(t *testing.T) {
	t.Run("default_values", func(t *testing.T) {
		cmd := &cmdWithDefaults{}
		app, err := BindUrfave("server", cmd)
		assert.Nil(t, err)
		assert.NotNil(t, app)
	})
}

type subCmd struct {
	Value int
}

func (s *subCmd) Run([]string) {
}

type rootCmd struct {
}

func (r *rootCmd) SubCommands() Map {
	return Map{
		"sub": &subCmd{},
	}
}

func TestUrfaveSubcommands(t *testing.T) {
	t.Run("subcommands", func(t *testing.T) {
		root := &rootCmd{}
		app, err := BindUrfave("root", root)
		assert.Nil(t, err)
		assert.NotNil(t, app)
		assert.Len(t, app.Commands, 1)
		assert.Equal(t, "sub", app.Commands[0].Name)
	})
}

