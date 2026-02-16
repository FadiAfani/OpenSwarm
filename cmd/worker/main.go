package worker

import (
	"log"
	"net"
	pb "openswarm/api/gen"
	worker "openswarm/internal/worker/agent"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpc_server := grpc.NewServer()

	pb.RegisterWorkerServiceServer(
		grpc_server,
		worker.NewGRPCServer(),
	)

	grpc_server.Serve(lis)

}
