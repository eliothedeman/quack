package quack

// Func turns a function into a command.
type Func func([]string)

// Run the basic command wrapper
func (f Func) Run(args []string) {
	f(args)
}

// Map is a wrapper around a group of commands. No Need to define a struct
type Map map[string]Unit

// SubCommands returns the commands the basic group wraps around.
func (m Map) SubCommands() Map {
	return m
}
