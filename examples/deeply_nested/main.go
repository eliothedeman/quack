package main

import "github.com/eliothedeman/quack"

type a struct {
}

func (a) SubCommands() quack.Map {
	return quack.Map{
		"b": new(b),
	}
}

type b struct {
}

func (b) SubCommands() quack.Map {
	return quack.Map{
		"c": new(c),
	}
}

type c struct {
	XX string `default:"YYY" short:"x"`
	Y  int    `help:"this is a help message"`
	Z  bool   `default:"true"`
}

func (c) Run([]string) {

}

func (c) Help() string {
	return "the nested c command"
}

func main() {
	cmd, _ := quack.Bind("nested", new(a))
	cmd.Execute()
}
