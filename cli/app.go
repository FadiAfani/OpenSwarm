package cli

import (
	"fmt"

	"openswarm/cli/commands"
)

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing command\n\n%s", Usage())
	}

	switch args[0] {
	case "worker":
		return commands.RunWorker(args[1:])
	case "coordinator":
		return commands.RunCoordinator(args[1:])
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], Usage())
	}
}

func Usage() string {
	return `Usage:
  openswarm-cli <command> [subcommand] [flags]

Commands:
  coordinator start Start coordinator server
  worker register   Register a worker with coordinator`
}
