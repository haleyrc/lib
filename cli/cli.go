// Package cli provides a minimal framework for building command-line
// applications with named subcommands.
//
// The CLI type is generic over an App parameter, which allows passing
// application-level dependencies (e.g. database connections, configuration)
// through to command handlers without global state.
//
// CLIs can be nested by using [CLI.Command] to convert a child CLI into a
// [Command] that can be registered on a parent.
package cli

import (
	"context"
	"fmt"
	"log"
	"sort"
)

// Exit codes returned by [CLI.Run].
const (
	ExitCodeSuccess = 0
	ExitCodeError   = 1
	ExitCodeUsage   = 2
)

// Command is a named subcommand with a description and a run function. The
// description is displayed in help output. The run function receives the
// context, the application value, and any remaining arguments after the command
// name.
type Command[App any] struct {
	Description string
	Run         func(ctx context.Context, app App, args []string) error
}

// CLI is a command-line application with a set of named subcommands. The App
// type parameter represents application-level dependencies shared across
// commands.
type CLI[App any] struct {
	Name     string
	Commands map[string]Command[App]
}

// Run parses the first element of args as a command name and dispatches to the
// matching [Command]. It returns an exit code suitable for use with
// [os.Exit]. If no command is given or the command is "help", it prints usage
// information and returns [ExitCodeUsage].
func (cli *CLI[App]) Run(ctx context.Context, app App, args []string) int {
	command, rest := getCommand(args)

	switch command {
	case "", "help":
		cli.printHelp()
		return ExitCodeUsage
	}

	cmd, ok := cli.Commands[command]
	if !ok {
		log.Printf("Unknown command: %s\n", command)
		cli.printHelp()
		return ExitCodeError
	}

	if err := cmd.Run(ctx, app, rest); err != nil {
		log.Println("ERROR:", err)
		return ExitCodeError
	}

	return ExitCodeSuccess
}

// Command wraps this CLI as a [Command] with the given description, allowing
// it to be registered as a subcommand of a parent [CLI]. This enables nested
// command hierarchies.
func (cli *CLI[App]) Command(desc string) Command[App] {
	return Command[App]{
		Description: desc,
		Run: func(ctx context.Context, app App, args []string) error {
			if code := cli.Run(ctx, app, args); code != ExitCodeSuccess {
				return fmt.Errorf("%s failed with exit code %d", cli.Name, code)
			}
			return nil
		},
	}
}

func (cli *CLI[App]) printHelp() {
	fmt.Printf("Usage: %s <command> [arguments]\n\n", cli.Name)
	fmt.Println("Commands:")

	names := make([]string, 0, len(cli.Commands))
	for name := range cli.Commands {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Printf("  %-12s %s\n", name, cli.Commands[name].Description)
	}
}

func getCommand(args []string) (string, []string) {
	if len(args) == 0 {
		return "", nil
	}
	if len(args) == 1 {
		return args[0], nil
	}
	return args[0], args[1:]
}
