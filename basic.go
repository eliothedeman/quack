package quack

// Func turns a function into a command.
type Func func([]string)

func (f Func) Run(args []string) {
	f(args)
}

type helpWrapper string

func (h helpWrapper) Help() string {
	return string(h)
}

type funcWithHelp struct {
	Func
	helpWrapper
}

// WithHelp creates a wrapper that has a help function for the function
func (f Func) WithHelp(help string) Command {
	return funcWithHelp{
		Func:        f,
		helpWrapper: helpWrapper(help),
	}
}

// Map is a wrapper around a group of commands. No Need to define a struct
type Map map[string]any

// WithHelp will return the map with a "Helper" interface attached
func (m Map) WithHelp(help string) Group {
	return &mapWithHelp{
		Map:         m,
		helpWrapper: helpWrapper(help),
	}
}

// SubCommands returns the commands the basic group wraps around.
func (m Map) SubCommands() Map {
	return m
}

type mapWithHelp struct {
	Map
	helpWrapper
}
