package main

import (
	"fmt"
	"os"
	"strings"

	"openswarm/cli/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage())
		os.Exit(1)
	}

	// Parsers expect "openswarm-cli <subcommand> <verb> ..." so token index 2 is "start" or "register"
	in := "openswarm-cli " + strings.Join(os.Args[1:], " ")

	switch os.Args[1] {
	case "coordinator":
		switch os.Args[2] {
		case "start":
			cmd, err := commands.ParseCoordinatorStartCmd(in)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			commands.RunCoordinatorStartCmd(cmd)
		default:
			fmt.Fprintf(os.Stderr, "unknown coordinator subcommand %q\n", os.Args[2])
			fmt.Fprintln(os.Stderr, usage())
			os.Exit(1)
		}
	case "worker":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "worker subcommand required: start | register")
			fmt.Fprintln(os.Stderr, usage())
			os.Exit(1)
		}
		switch os.Args[2] {
		case "start":
			cmd, err := commands.ParseWorkerStartCmd(in)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			commands.RunWorkerStartCmd(cmd)
		case "register":
			cmd, err := commands.ParseRegisterWorkerCmd(in)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			commands.RunRegisterWorkerCmd(cmd)

		case "list-workers":
			cmd, err := commands.ParseListWorkersCmd(in)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			commands.RunListWorkersCmd(cmd)
		default:
			fmt.Fprintf(os.Stderr, "unknown worker subcommand %q\n", os.Args[2])
			fmt.Fprintln(os.Stderr, usage())
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		fmt.Fprintln(os.Stderr, usage())
		os.Exit(1)
	}
}

func usage() string {
	return `usage: openswarm-cli <command> [options]

Commands:
  coordinator start   Start the coordinator server (e.g. --port 8081)
  worker start        Start a worker server (e.g. --port 8082)
  worker register     Register this worker with a coordinator (e.g. --coordinator-addr <addr>)
  worker list-workers Lists all registered workers with a coordinator (e.g. --coordinator-addr <addr>)
`
}
