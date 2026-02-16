package commands

import (
	"fmt"
	"log"
	"net"
	pb "openswarm/api/gen"
	"openswarm/internal/coordinator"

	"google.golang.org/grpc"
)

type CoordinatorStartCmd struct {
	port *string
}

func runCoordinatorStartCmd(cmd CoordinatorStartCmd) {
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

	server.Serve(lis)
}

func Usage() string {
	return `
		usage: openswarm-cli coordinator start --addr [address]
	`
}
