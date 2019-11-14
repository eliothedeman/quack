package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/eliothedeman/quack"
)

var (
	_ quack.Group   = new(cmds)
	_ quack.Command = new(ls)
)

type cmds struct {
}

func (c cmds) Help() string {
	return "A collection of commands"
}

func (c cmds) SubCommands() map[string]quack.Unit {
	return map[string]quack.Unit{
		"ls": new(ls),
	}
}

type ls struct {
	Path  string `short:"p" help:"path to list" default:"."`
	Files bool   `short:"f" help:"have files or not" default:"true"`
}

func (ls) Help() string {
	return "list the dir"
}

func (l *ls) Default() {
	l.Path = "./"
}

func (l *ls) Validate() error {
	if l.Path == "wrong" {
		return fmt.Errorf("%s is not a valid path", l.Path)
	}
	return nil
}

func (l ls) Run([]string) {
	d, err := ioutil.ReadDir(l.Path)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range d {
		fmt.Println(d.Name())
	}
}

func main() {
	quack.Run(
		"example",
		quack.WithArgs(os.Args),
		quack.WithGroup(new(cmds)),
	)
}
