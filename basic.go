package quack

// CommandFunc turns a function into a command.
type CommandFunc func([]string)

// Run the basic command wrapper
func (c CommandFunc) Run(args []string) {
	c(args)
}

// CommandMap is a wrapper around a group of commands. No Need to define a struct
type CommandMap map[string]Unit

// SubCommands returns the commands the basic group wraps around.
func (c CommandMap) SubCommands() map[string]Unit {
	return map[string]Unit(c)
}
