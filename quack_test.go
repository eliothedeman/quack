package quack

import (
	"errors"
	"os"
	"testing"
)

type TRoot struct {
	t TSub `ignore:"true"`
}

type TSub struct {
	Name    string
	Called  bool
	private int
}

func (t *TSub) Run([]string) {
	t.Called = true
}

func (t *TRoot) SubCommands() Map {
	return Map{
		"hello": &t.t,
	}
}

func TestRunGroup(t *testing.T) {
	r := &TRoot{}
	Run("test", WithGroup(r), WithArgs(
		[]string{"hello", "--name", "wilson"},
	))

	if r.t.Called != true {
		t.Fatal()
	}
}

func TestRunGroupWithOSArgs(t *testing.T) {
	r := &TRoot{}

	os.Args = []string{"./mycmd", "hello", "--name", "not-wison"}
	Run("test", WithGroup(r), WithArgs(
		[]string{"hello", "--name", "wilson"},
	))

	if r.t.Name == "not-wilson" {
		t.Errorf("r.t.Name == not-wilson")
	}

	if r.t.Called != true {
		t.Fatal()
	}
}

func TestRunCommand(t *testing.T) {
	r := &TSub{}
	os.Args = []string{"./bin", "hello", "--name", "wilson"}
	Run("test", WithCommand(r))

	if r.Called != true {
		t.Fatal()
	}
}

func Test_hasHelpArgNoOverride(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{
			"help no flags",
			[]string{
				"help",
			},
			false,
		},
		{
			"single help",
			[]string{
				"--help",
			},
			true,
		},
		{
			"subcommand help",
			[]string{
				"sub", "--help",
			},
			false,
		},
		{
			"multple flags",
			[]string{
				"--sub", "--help",
			},
			true,
		},
		{
			"help first",
			[]string{
				"--help", "next",
			},
			true,
		},
		{
			"short first",
			[]string{
				"-h", "next",
			},
			false,
		},
		{
			"no args",
			[]string{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasHelpArg(tt.args, false); got != tt.want {
				t.Errorf("hasHelpArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasHelpArgWithOverride(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{
			"help no flags",
			[]string{
				"help",
			},
			false,
		},
		{
			"single help",
			[]string{
				"--help",
			},
			true,
		},
		{
			"subcommand help",
			[]string{
				"sub", "--help",
			},
			false,
		},
		{
			"multple flags",
			[]string{
				"--sub", "--help",
			},
			true,
		},
		{
			"help first",
			[]string{
				"--help", "next",
			},
			true,
		},
		{
			"short first",
			[]string{
				"-h", "next",
			},
			true,
		},
		{
			"no args",
			[]string{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasHelpArg(tt.args, true); got != tt.want {
				t.Errorf("hasHelpArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

type valid struct {
	Val string
}

func (v *valid) Validate() error {
	if v.Val == "invalid" {
		return errors.New("invalid")
	}
	return nil
}

func (v valid) Run([]string) {

}
func TestValidator(t *testing.T) {
	v := new(valid)
	err := run("invalid", v, []string{"--val", "invalid"})
	if err == nil {
		t.Error("Invalid command should return validation error")
		return
	}

	err = run("valid", v, []string{"--val", "yup"})
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
}

func TestValidateUnit(t *testing.T) {
	err := run("bad", Map{
		"bad_cmd": 10,
	}, []string{})
	if !errors.Is(err, ErrWrongType) {
		t.Errorf("[%s] is wrong type of error", err)
	}
}
