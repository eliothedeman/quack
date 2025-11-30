# quack

![GoDoc](https://godoc.org/github.com/eliothedeman/quack?status.svg)
![Test](https://github.com/eliothedeman/quack/actions/workflows/test.yml/badge.svg)

An interface driven cli lib for go.

## Overview

Have you ever said

> Why can't writing CLIs be as easy as defining a struct and writing a function?

Then quack is for you.

## Features

- 🏗️ **Struct-based CLI** - Define commands using simple Go structs
- 🎯 **Positional arguments** - Use `arg:"1"`, `arg:"2"` tags for positional args
- 🔁 **Repeated arguments** - Slices are automatically treated as variadic
- 🏷️ **Named options** - Support for short (`-f`) and long (`--file`) flags
- 🪆 **Nested sub-commands** - Easy command hierarchies
- 🔌 **Multiple frameworks** - Works with [Cobra](https://github.com/spf13/cobra) and [urfave/cli](https://github.com/urfave/cli) v2

## Framework Support

### Integration with [spf13/cobra](https://github.com/spf13/cobra)
The `BindCobra` API creates a `cobra.Command` from the given structure. This allows for easy integration
with existing CLIs that use this framework.

```go
cmd, err := quack.BindCobra("myapp", new(MyCommand))
```

### Integration with [urfave/cli](https://github.com/urfave/cli) v2
The `BindUrfave` API creates a `cli.App` from the given structure, providing seamless integration
with urfave/cli v2.

```go
app, err := quack.BindUrfave("myapp", new(MyCommand))
```

For urfave-specific features like accessing `*cli.Context`, your command can implement the `UrfaveCommand` interface:

```go
type MyCmd struct {
    Name string
}

func (m *MyCmd) Run(ctx *cli.Context) error {
    // Access urfave-specific context features
    fmt.Println(ctx.App.Name)
    return nil
}
```

### Other framework support
Both Cobra and urfave/cli are now fully supported! Feel free to file an issue if you'd like support for additional frameworks.

## Examples

### A simple command

_main.go_

```go
type ToHex struct {
	Input int
}

func (t ToHex) Run() {
	fmt.Printf("%x", t.Input)
}

func main() {
	quack.MustBind("tohex", new(ToHex)).Execute()
}
```

Can now be run

```
go run main.go --input 12334
302e
```

### Positional arguments

Use the `arg` tag to specify positional arguments:

```go
type CopyCmd struct {
	Source string `arg:"1"`
	Target string `arg:"2"`
}

func (c *CopyCmd) Run([]string) {
	fmt.Printf("Copying %s to %s\n", c.Source, c.Target)
}

func main() {
	quack.MustBind("copy", new(CopyCmd)).Execute()
}
```

```bash
$ go run main.go copy source.txt target.txt
Copying source.txt to target.txt
```

### Repeated arguments (slices)

Slices are automatically treated as variadic arguments:

```go
type CompileCmd struct {
	Files []string `arg:"1"`  // Consumes all remaining args
}

func (c *CompileCmd) Run([]string) {
	fmt.Printf("Compiling: %v\n", c.Files)
}
```

```bash
$ go run main.go compile file1.go file2.go file3.go
Compiling: [file1.go file2.go file3.go]
```

### Mixed flags and positional args

```go
type BuildCmd struct {
	Verbose bool     `short:"v" help:"Enable verbose output"`
	Output  string   `short:"o" help:"Output file"`
	Files   []string `arg:"1" help:"Source files"`
}

func (b *BuildCmd) Run([]string) {
	if b.Verbose {
		fmt.Printf("Building %v -> %s\n", b.Files, b.Output)
	}
	// ... build logic
}
```

```bash
$ go run main.go build -v -o app.bin main.go utils.go
Building [main.go utils.go] -> app.bin
```

### Repeated flags

Slices work as flags too:

```go
type ServerCmd struct {
	Port    int      `short:"p" default:"8080"`
	Allowed []string `short:"a" help:"Allowed IPs"`
}
```

```bash
$ go run main.go server -a 192.168.1.1 -a 10.0.0.1
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
	quack.MustBind("nested", new(a))
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

## Available Struct Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `arg:"N"` | Positional argument at position N (1-indexed) | `arg:"1"` |
| `short:"x"` | Short flag name | `short:"v"` for `-v` |
| `long:"name"` | Long flag name (auto-generated from field name if not specified) | `long:"verbose"` |
| `default:"value"` | Default value | `default:"8080"` |
| `help:"text"` | Help text for the option | `help:"Port to listen on"` |
| `ignore:""` | Ignore this field | `ignore:""` |

**Note:** Slice types are automatically treated as repeated/variadic - no special tag needed!

