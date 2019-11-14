# quack

[![GoDoc](https://godoc.org/github.com/eliothedeman/quack?status.svg)](https://godoc.org/github.com/eliothedeman/quack)

An interface driven cli lib for go.

## Overview

Have you ever said

> Why can't writing CLIs be as easy as defining a struct and writing a function?

Then quack is for you.

### A simple command

_main.go_

```go
type ToHex struct {
	Input int
}

func (t ToHex) Run([]args) {
	fmt.Printf("%x", t.Input)
}

func main() {
	quack.Run("tohex", quack.WithCommand(&ToHex{}))
}
```

Can now be run

```
go run main.go --input 12334
302e
```

## Interfaces

Quack uses interfaces to let you specialize the behavior of your cli. This lets you use as much or as little of the library as you want.

```go
// Command is a runnable command that doesn't have sub commands
type Command interface {
    Run([]string)
}

// Group is a set of subcommands or sub groups.
type Group interface {
    SubCommands() map[string]Unit
}

// Unit is a placeholder for commands and groups.
type Unit interface{}

// Validator is a command or argument that wants to be validated.
type Validator interface {
	Validate() error
}

// Defaulter can set up the default arguments of a command
type Defaulter interface {
	Default()
}

// Parser is an argument that wants to parse itself.
type Parser interface {
	Parse(string) error
}

// Helper returns usage information for a command or group.
type Helper interface {
	Help() string
}
```

## Structs

Quack uses structs and struct tags for argument parsing. Before a command is run, it's arguments will be parsed and fill the struct the command is run on.

```go
type ls struct {
	Path  string `short:"p" help:"path to list" default:"."`
}
```

In this example, "Path" can be refered to as `-p` or `--path` and has a default value of "`.`"

## Example

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/eliothedeman/quack"
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
	if l.Path == "invalid" {
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
		quack.WithGroup(new(cmds))
	)
}
```
