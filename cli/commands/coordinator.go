package commands

import (
	"flag"
	"fmt"
	"net"

	pb "openswarm/api/gen"
	co "openswarm/internal/coordinator"

	"google.golang.org/grpc"
)

func RunCoordinator(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing coordinator subcommand\n\n%s", coordinatorUsage())
	}

	switch args[0] {
	case "start":
		return runCoordinatorStart(args[1:])
	default:
		return fmt.Errorf("unknown coordinator subcommand %q\n\n%s", args[0], coordinatorUsage())
	}
}

func runCoordinatorStart(args []string) error {
	fs := flag.NewFlagSet("coordinator start", flag.ContinueOnError)
	addr := fs.String("addr", "localhost:8080", "Coordinator gRPC bind address")

	if err := fs.Parse(args); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", *addr, err)
	}

	server := grpc.NewServer()
	pb.RegisterCoordinatorServiceServer(server, co.NewCoordinatorGRPCServer())

	fmt.Printf("coordinator listening on %s\n", *addr)
	return server.Serve(lis)
}

func coordinatorUsage() string {
	return `Usage:
  openswarm-cli coordinator start [flags]

Flags:
  --addr string   Coordinator gRPC bind address (default "localhost:8080")`
}
