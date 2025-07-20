# Agent Guidelines

## Build/Test Commands
- `go build -v .` - Build the package
- `go test ./...` - Run all tests
- `go test -run TestName ./...` - Run specific test
- `go test -v ./...` - Verbose test output

## Code Style
- **Imports**: Standard library first, third-party second, local last
- **Formatting**: Use `gofmt` (handled automatically by most editors)
- **Naming**: CamelCase for exported, camelCase for unexported
- **Error handling**: Return errors explicitly, use `fmt.Errorf` with context
- **Types**: Interface-based design (Command, SimpleCommand, CobraCommand, Group)
- **Tags**: Use struct tags: `help:"description"`, `default:"value"`, `short:"x"`
- **Interfaces**: Implement appropriate interfaces (Command, SimpleCommand, Group, Validator, Defaulter, Parser, Helper)
- **Structure**: Commands are structs with exported fields for CLI flags

## Rules
- Put your planning steps in plan.md These files are ignored by the git repo so will only
  be available for one coding session. Use this file as your temporary memory store.
