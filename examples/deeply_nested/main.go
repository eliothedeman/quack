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
	X, Y, Z string
}

func (c) Run([]string) {

}

func (c) Help() string {
	return "the nested c command"
}

func main() {
	quack.Run("nested", quack.WithGroup(new(a)))
}
