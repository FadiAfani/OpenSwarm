package commands

import (
	"fmt"
	"log"
	"net"
	pb "openswarm/api/gen"
	"openswarm/internal/worker/agent"

	"google.golang.org/grpc"
)

type WorkerStartCmd struct {
	port *string
}

func RunWorkerStartCmd(cmd WorkerStartCmd) {
	port := "8082"
	if cmd.port != nil {
		port = *cmd.port
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterWorkerServiceServer(server, agent.NewGRPCServer())

	server.Serve(lis)

}
