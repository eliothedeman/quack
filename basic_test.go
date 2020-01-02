package quack

import "testing"

import . "github.com/stretchr/testify/assert"

func TestBasicCommand(t *testing.T) {
	called := false
	Run("test", WithCommand(Func(func([]string) {
		called = true
	})))
	if !called {
		t.Fail()
	}
}

func TestCommandFunc(t *testing.T) {
	called := false
	Run("test", WithGroup(Map{
		"my_name": Func(func([]string) {
			called = true
		}),
	}),
		WithArgs([]string{"my_name"}),
	)

	if !called {
		t.Fail()
	}
}

func TestGroupWithHelp(t *testing.T) {
	called := false
	Run("test", WithGroup(Map{
		"my_name": Func(func([]string) {
			called = true
		}),
	}.WithHelp("my test command")),
		WithArgs([]string{"my_name"}),
	)

	if !called {
		t.Fail()
	}
}

type fmtHelper struct {
	Map
	help string
}

func (f fmtHelper) Help() string {
	return f.help
}
func TestGroupWithHelpSameAsStruct(t *testing.T) {
	m := Map{
		"my_name": Func(func([]string) {
		}),
	}

	helpStr := "my help string"
	a := fmtHelp("tt", m.WithHelp(helpStr))
	b := fmtHelp("tt", fmtHelper{Map: m, help: helpStr})
	Equal(t, a, b)
	Contains(t, a, helpStr)
}
