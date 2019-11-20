package quack

import "testing"

func TestBasicCommand(t *testing.T) {
	called := false
	Run("test", WithCommand(CommandFunc(func([]string) {
		called = true
	})))
	if !called {
		t.Fail()
	}
}

func TestCommandFunc(t *testing.T) {
	called := false
	Run("test", WithGroup(CommandMap{
		"my_name": CommandFunc(func([]string) {
			called = true
		}),
	}),
		WithArgs([]string{"my_name"}),
	)

	if !called {
		t.Fail()
	}
}
