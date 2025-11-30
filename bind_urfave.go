package quack

import (
	"fmt"
	"reflect"

	"github.com/urfave/cli/v2"
)

// toUrfaveApp converts a node to a *cli.App
func (c *node) toUrfaveApp() *cli.App {
	app := &cli.App{
		Name:  c.name,
		Usage: c.short,
	}
	if c.long != "" {
		app.Description = c.long
	}

	// Set flags
	app.Flags = c.toUrfaveFlags()

	// Set commands
	for _, s := range c.subcommands {
		app.Commands = append(app.Commands, s.toUrfaveCommand())
	}

	// Set action
	if c.run != nil {
		// Check if target implements UrfaveCommand
		// We need to check if it implements the interface with *cli.Context parameter
		type urfaveCommandWithContext interface {
			Run(ctx *cli.Context) error
		}

		if urfaveCmd, ok := c.target.(urfaveCommandWithContext); ok {
			app.Action = func(ctx *cli.Context) error {
				// Parse flags from context into struct fields
				if err := c.parseUrfaveFlags(ctx); err != nil {
					return err
				}
				// Parse positional arguments from ctx.Args()
				args := ctx.Args().Slice()
				if err := c.parsePositionalArgs(args); err != nil {
					return err
				}
				// Call the UrfaveCommand's Run method directly
				return urfaveCmd.Run(ctx)
			}
		} else {
			// Use the standard run function for other command types
			originalRun := c.run
			app.Action = func(ctx *cli.Context) error {
				// Parse flags from context into struct fields
				if err := c.parseUrfaveFlags(ctx); err != nil {
					return err
				}
				// Parse positional arguments from ctx.Args()
				args := ctx.Args().Slice()
				if err := c.parsePositionalArgs(args); err != nil {
					return err
				}
				// Call the original run function with nil cobra command since we're in urfave context
				originalRun(nil, args)
				return nil
			}
		}
	}

	return app
}

// toUrfaveCommand converts a node to a *cli.Command
func (c *node) toUrfaveCommand() *cli.Command {
	cmd := &cli.Command{
		Name:  c.name,
		Usage: c.short,
	}
	if c.long != "" {
		cmd.Description = c.long
	}

	// Set flags
	cmd.Flags = c.toUrfaveFlags()

	// Set subcommands
	for _, s := range c.subcommands {
		cmd.Subcommands = append(cmd.Subcommands, s.toUrfaveCommand())
	}

	// Set action
	if c.run != nil {
		// Check if target implements UrfaveCommand
		// We need to check if it implements the interface with *cli.Context parameter
		type urfaveCommandWithContext interface {
			Run(ctx *cli.Context) error
		}

		if urfaveCmd, ok := c.target.(urfaveCommandWithContext); ok {
			cmd.Action = func(ctx *cli.Context) error {
				// Parse flags from context into struct fields
				if err := c.parseUrfaveFlags(ctx); err != nil {
					return err
				}
				// Parse positional arguments from ctx.Args()
				args := ctx.Args().Slice()
				if err := c.parsePositionalArgs(args); err != nil {
					return err
				}
				// Call the UrfaveCommand's Run method directly
				return urfaveCmd.Run(ctx)
			}
		} else {
			// Use the standard run function for other command types
			originalRun := c.run
			cmd.Action = func(ctx *cli.Context) error {
				// Parse flags from context into struct fields
				if err := c.parseUrfaveFlags(ctx); err != nil {
					return err
				}
				// Parse positional arguments from ctx.Args()
				args := ctx.Args().Slice()
				if err := c.parsePositionalArgs(args); err != nil {
					return err
				}
				// Call the original run function with nil cobra command since we're in urfave context
				originalRun(nil, args)
				return nil
			}
		}
	}

	return cmd
}

// toUrfaveFlags converts the node's options to urfave/cli flags
func (c *node) toUrfaveFlags() []cli.Flag {
	var flags []cli.Flag
	for _, o := range c.options {
		if flag := o.toUrfaveFlag(); flag != nil {
			flags = append(flags, flag)
		}
	}
	return flags
}

// toUrfaveFlag converts an option to a urfave/cli flag
func (o *option) toUrfaveFlag() cli.Flag {
	if o.Ignore {
		return nil
	}

	name := o.Name
	usage := o.Help

	// Build aliases (short flags)
	var aliases []string
	if o.Short != "" {
		aliases = append(aliases, o.Short)
	}

	v := o.Target

	// Handle slice types (automatically repeated)
	if v.Kind() == reflect.Slice {
		elemType := v.Type().Elem()
		switch elemType.Kind() {
		case reflect.String:
			return &cli.StringSliceFlag{
				Name:    name,
				Aliases: aliases,
				Usage:   usage,
			}
		case reflect.Int:
			return &cli.IntSliceFlag{
				Name:    name,
				Aliases: aliases,
				Usage:   usage,
			}
		case reflect.Int64:
			return &cli.Int64SliceFlag{
				Name:    name,
				Aliases: aliases,
				Usage:   usage,
			}
		case reflect.Uint:
			return &cli.Uint64SliceFlag{
				Name:    name,
				Aliases: aliases,
				Usage:   usage,
			}
		case reflect.Uint64:
			return &cli.Uint64SliceFlag{
				Name:    name,
				Aliases: aliases,
				Usage:   usage,
			}
		case reflect.Float64:
			return &cli.Float64SliceFlag{
				Name:    name,
				Aliases: aliases,
				Usage:   usage,
			}
		default:
			panic(fmt.Sprintf("Unable to handle slice type for urfave flag: %v", elemType.Kind()))
		}
	}

	// Handle non-slice types
	switch v.Kind() {
	case reflect.Bool:
		boolVal := o.Default == "true"
		return &cli.BoolFlag{
			Name:    name,
			Aliases: aliases,
			Usage:   usage,
			Value:   boolVal,
		}
	case reflect.Int:
		intVal := 0
		if o.Default != "" {
			fmt.Sscanf(o.Default, "%d", &intVal)
		}
		return &cli.IntFlag{
			Name:    name,
			Aliases: aliases,
			Usage:   usage,
			Value:   intVal,
		}
	case reflect.Int64:
		var intVal int64 = 0
		if o.Default != "" {
			fmt.Sscanf(o.Default, "%d", &intVal)
		}
		return &cli.Int64Flag{
			Name:    name,
			Aliases: aliases,
			Usage:   usage,
			Value:   intVal,
		}
	case reflect.Uint:
		var uintVal uint = 0
		if o.Default != "" {
			fmt.Sscanf(o.Default, "%d", &uintVal)
		}
		return &cli.UintFlag{
			Name:    name,
			Aliases: aliases,
			Usage:   usage,
			Value:   uintVal,
		}
	case reflect.Uint64:
		var uintVal uint64 = 0
		if o.Default != "" {
			fmt.Sscanf(o.Default, "%d", &uintVal)
		}
		return &cli.Uint64Flag{
			Name:    name,
			Aliases: aliases,
			Usage:   usage,
			Value:   uintVal,
		}
	case reflect.Float64:
		var floatVal float64 = 0
		if o.Default != "" {
			fmt.Sscanf(o.Default, "%f", &floatVal)
		}
		return &cli.Float64Flag{
			Name:    name,
			Aliases: aliases,
			Usage:   usage,
			Value:   floatVal,
		}
	case reflect.String:
		return &cli.StringFlag{
			Name:    name,
			Aliases: aliases,
			Usage:   usage,
			Value:   o.Default,
		}
	default:
		panic(fmt.Sprintf("Unable to handle type for urfave flag: %v", v.Kind()))
	}
}

// parseUrfaveFlags reads flag values from the cli.Context and assigns them to the struct fields
func (c *node) parseUrfaveFlags(ctx *cli.Context) error {
	for _, opt := range c.options {
		if opt.Ignore {
			continue
		}

		v := opt.Target
		name := opt.Name

		// Handle slice types
		if v.Kind() == reflect.Slice {
			elemType := v.Type().Elem()
			switch elemType.Kind() {
			case reflect.String:
				values := ctx.StringSlice(name)
				if values != nil {
					v.Set(reflect.ValueOf(values))
				}
			case reflect.Int:
				values := ctx.IntSlice(name)
				if values != nil {
					v.Set(reflect.ValueOf(values))
				}
			case reflect.Int64:
				values := ctx.Int64Slice(name)
				if values != nil {
					v.Set(reflect.ValueOf(values))
				}
			case reflect.Uint, reflect.Uint64:
				values := ctx.Uint64Slice(name)
				if values != nil {
					// Convert []uint64 to []uint if needed
					if elemType.Kind() == reflect.Uint {
						uintSlice := make([]uint, len(values))
						for i, val := range values {
							uintSlice[i] = uint(val)
						}
						v.Set(reflect.ValueOf(uintSlice))
					} else {
						v.Set(reflect.ValueOf(values))
					}
				}
			case reflect.Float64:
				values := ctx.Float64Slice(name)
				if values != nil {
					v.Set(reflect.ValueOf(values))
				}
			default:
				return fmt.Errorf("unsupported slice type for flag %s: %v", name, elemType.Kind())
			}
			continue
		}

		// Handle non-slice types
		switch v.Kind() {
		case reflect.Bool:
			v.SetBool(ctx.Bool(name))
		case reflect.Int:
			v.SetInt(int64(ctx.Int(name)))
		case reflect.Int64:
			v.SetInt(ctx.Int64(name))
		case reflect.Uint:
			v.SetUint(uint64(ctx.Uint(name)))
		case reflect.Uint64:
			v.SetUint(ctx.Uint64(name))
		case reflect.Float64:
			v.SetFloat(ctx.Float64(name))
		case reflect.String:
			v.SetString(ctx.String(name))
		default:
			return fmt.Errorf("unsupported type for flag %s: %v", name, v.Kind())
		}
	}
	return nil
}

// BindUrfave binds a structure to a *cli.App (and sub-commands)
func BindUrfave(name string, root any) (*cli.App, error) {
	rn := new(node)
	err := rn.fromStruct(name, root)
	if err != nil {
		return nil, err
	}
	return rn.toUrfaveApp(), nil
}

// MustBindUrfave will panic if BindUrfave returns an error
func MustBindUrfave(name string, root any) *cli.App {
	app, err := BindUrfave(name, root)
	if err != nil {
		panic(err)
	}
	return app
}
