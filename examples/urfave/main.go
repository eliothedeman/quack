package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eliothedeman/quack"
	"github.com/urfave/cli/v3"
)

// Simple command using UrfaveCommand interface for full cli.Command access
type GreetCmd struct {
	Name   string `short:"n" help:"Name to greet" default:"World"`
	Formal bool   `short:"f" help:"Use formal greeting"`
}

func (g *GreetCmd) Run(ctx context.Context, cmd *cli.Command) error {
	greeting := "Hello"
	if g.Formal {
		greeting = "Good day"
	}

	// Can access urfave-specific command features
	fmt.Printf("%s, %s!\n", greeting, g.Name)
	fmt.Printf("Command name: %s\n", cmd.Name)
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
	cmd, err := quack.BindUrfave("example", root)
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
