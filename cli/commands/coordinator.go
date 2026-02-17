package commands

import (
	"fmt"
	"log"
	"net"
	pb "openswarm/api/gen"
	"openswarm/internal/coordinator"
	"openswarm/internal/worker/agent"
	"strings"

	"google.golang.org/grpc"
)

type CoordinatorStartCmd struct {
	port *string
}

func ParseCoordinatorStartCmd(in string) (CoordinatorStartCmd, error) {
	cmd := CoordinatorStartCmd{}
	tokens := strings.Fields(in)
	if len(tokens) < 3 {
		return cmd, fmt.Errorf("command is too short")
	}

	if tokens[2] != "start" {
		return cmd, fmt.Errorf("expected 'start' but got %q", tokens[2])
	}

	for i := 3; i < len(tokens); i++ {
		switch tokens[i] {
		case "--port":
			if i+1 >= len(tokens) {
				return cmd, fmt.Errorf("missing value for port")
			}
			port := tokens[i+1]
			cmd.port = &port
			i++
		default:
			return cmd, fmt.Errorf("unknown option %q", tokens[i])
		}
	}

	return cmd, nil
}

func RunCoordinatorStartCmd(cmd CoordinatorStartCmd) {
	port := "8081"

	if cmd.port != nil {
		port = *cmd.port
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterCoordinatorServiceServer(server, coordinator.NewCoordinatorGRPCServer())
	pb.RegisterWorkerServiceServer(server, agent.NewGRPCServer())

	server.Serve(lis)
}

func Usage() string {
	return `
		usage: openswarm-cli coordinator start --addr [address]
	`
}
