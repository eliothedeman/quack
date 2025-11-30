package main

import (
	"fmt"
	"log"
	"os"

	"github.com/eliothedeman/quack"
	"github.com/urfave/cli/v2"
)

// Simple command using UrfaveCommand interface for full cli.Context access
type GreetCmd struct {
	Name   string `short:"n" help:"Name to greet" default:"World"`
	Formal bool   `short:"f" help:"Use formal greeting"`
}

func (g *GreetCmd) Run(ctx *cli.Context) error {
	greeting := "Hello"
	if g.Formal {
		greeting = "Good day"
	}

	// Can access urfave-specific context features
	fmt.Printf("%s, %s!\n", greeting, g.Name)
	fmt.Printf("App name from context: %s\n", ctx.App.Name)
	return nil
}

// Standard command using Command interface
type CountCmd struct {
	Files []string `arg:"1" help:"Files to count"`
}

func (c *CountCmd) Run(args []string) {
	fmt.Printf("Counting %d files:\n", len(c.Files))
	for i, file := range c.Files {
		fmt.Printf("  %d. %s\n", i+1, file)
	}
}

func (c *CountCmd) Help() string {
	return "Count files"
}

// Root command with subcommands
type Root struct{}

func (r *Root) SubCommands() quack.Map {
	return quack.Map{
		"greet": &GreetCmd{},
		"count": &CountCmd{},
	}
}

func main() {
	root := &Root{}
	app, err := quack.BindUrfave("example", root)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
