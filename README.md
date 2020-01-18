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

_main.go_

```go
var root = quack.Map{
	"build":	builtin.Build,
	"test": 	builtin.Test,
}

func main() {
	quack.Run("go", quack.WithGroup(root))
}
```

```
go run main.go -h
Usage: go <cmd> [args]
	build 	"build go source files"
	test 	"test go source files"
```
