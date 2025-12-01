package quack

import (
	"context"
	"fmt"
	"reflect"

	"github.com/urfave/cli/v3"
)

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
		cmd.Commands = append(cmd.Commands, s.toUrfaveCommand())
	}

	// Set action
	if c.run != nil {
		// Check if target implements UrfaveCommand with v3 signature
		type urfaveCommandV3 interface {
			Run(ctx context.Context, cmd *cli.Command) error
		}

		if urfaveCmd, ok := c.target.(urfaveCommandV3); ok {
			cmd.Action = func(ctx context.Context, cliCmd *cli.Command) error {
				// Parse flags from command into struct fields
				if err := c.parseUrfaveFlags(cliCmd); err != nil {
					return err
				}
				// Parse positional arguments
				args := cliCmd.Args().Slice()
				if err := c.parsePositionalArgs(args); err != nil {
					return err
				}
				// Validate options if command doesn't implement Validator
				if err := c.validateOptions(); err != nil {
					return err
				}
				// Call the UrfaveCommand's Run method directly
				return urfaveCmd.Run(ctx, cliCmd)
			}
		} else {
			// Use the standard run function for other command types
			originalRun := c.run
			cmd.Action = func(ctx context.Context, cliCmd *cli.Command) error {
				// Parse flags from command into struct fields
				if err := c.parseUrfaveFlags(cliCmd); err != nil {
					return err
				}
				// Parse positional arguments
				args := cliCmd.Args().Slice()
				if err := c.parsePositionalArgs(args); err != nil {
					return err
				}
				// Validate options if command doesn't implement Validator
				if err := c.validateOptions(); err != nil {
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
			return &cli.UintSliceFlag{
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

// parseUrfaveFlags reads flag values from the cli.Command and assigns them to the struct fields
func (c *node) parseUrfaveFlags(cmd *cli.Command) error {
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
				values := cmd.StringSlice(name)
				if values != nil {
					v.Set(reflect.ValueOf(values))
				}
			case reflect.Int:
				values := cmd.IntSlice(name)
				if values != nil {
					v.Set(reflect.ValueOf(values))
				}
			case reflect.Int64:
				values := cmd.Int64Slice(name)
				if values != nil {
					v.Set(reflect.ValueOf(values))
				}
			case reflect.Uint:
				values := cmd.UintSlice(name)
				if values != nil {
					v.Set(reflect.ValueOf(values))
				}
			case reflect.Uint64:
				values := cmd.Uint64Slice(name)
				if values != nil {
					v.Set(reflect.ValueOf(values))
				}
			case reflect.Float64:
				values := cmd.Float64Slice(name)
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
			v.SetBool(cmd.Bool(name))
		case reflect.Int:
			v.SetInt(int64(cmd.Int(name)))
		case reflect.Int64:
			v.SetInt(cmd.Int64(name))
		case reflect.Uint:
			v.SetUint(uint64(cmd.Uint(name)))
		case reflect.Uint64:
			v.SetUint(cmd.Uint64(name))
		case reflect.Float64:
			v.SetFloat(cmd.Float64(name))
		case reflect.String:
			v.SetString(cmd.String(name))
		default:
			return fmt.Errorf("unsupported type for flag %s: %v", name, v.Kind())
		}
	}
	return nil
}

// BindUrfave binds a structure to a *cli.Command (and sub-commands)
func BindUrfave(name string, root any) (*cli.Command, error) {
	rn := new(node)
	err := rn.fromStruct(name, root)
	if err != nil {
		return nil, err
	}
	return rn.toUrfaveCommand(), nil
}

// MustBindUrfave will panic if BindUrfave returns an error
func MustBindUrfave(name string, root any) *cli.Command {
	cmd, err := BindUrfave(name, root)
	if err != nil {
		panic(err)
	}
	return cmd
}
