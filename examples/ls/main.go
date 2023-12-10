package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eliothedeman/quack"
)

var root = quack.Map{
	"ls": &ls{},
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
	d, err := os.ReadDir(l.Path)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range d {
		fmt.Println(d.Name())
	}
}

func main() {
	cmd, err := quack.Bind("ls", root)
	if err != nil {
		log.Fatal(err)
	}
	cmd.Execute()
}
