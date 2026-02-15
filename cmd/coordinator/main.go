package main

import (
	"fmt"
	"log"
	"net"
	pb "openswarm/api/gen"
	co "openswarm/internal/coordinator"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:8080"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpc_server := grpc.NewServer()
	pb.RegisterCoordinatorServiceServer(
		grpc_server,
		co.NewCoordinatorGRPCServer(),
	)

	grpc_server.Serve(lis)

}
