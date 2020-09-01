# quack

![GoDoc](https://godoc.org/github.com/eliothedeman/quack?status.svg)
![Test](https://github.com/eliothedeman/quack/workflows/Test/badge.svg)

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

### A simple set of sub commands

_examples/deeply_nested/main.go_

```go
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
	quack.Run("nested", quack.WithGroup(new(a)))
}

```

```
go run examples/deeply_nested/main.go b -h
Usage:    b <cmd> [args]
      c the nested c command


go run examples/deeply_nested/main.go b c -h
Usage:    c [args]
        the nested c command
Flags:                                   
             --z         (default=true)  
Options:                                 
         -x, --xx string (default='YYY') 
             --y  int    (default=0)     this is a help message
```
