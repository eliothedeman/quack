package quack

// Func turns a function into a command.
type Func func([]string)

// Run the basic command wrapper
func (f Func) Run(args []string) {
	f(args)
}

// Map is a wrapper around a group of commands. No Need to define a struct
type Map map[string]Unit

// WithHelp will return the map with a "Helper" interface attached
func (m Map) WithHelp(help string) Group {
	return &mapWithHelp{
		Map:  m,
		help: help,
	}
}

// SubCommands returns the commands the basic group wraps around.
func (m Map) SubCommands() Map {
	return m
}

type mapWithHelp struct {
	Map
	help string
}

func (m *mapWithHelp) Help() string {
	return m.help
}
