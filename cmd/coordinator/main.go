package main

import (
	"fmt"
	"log"
	"net"
	pb "openswarm/api/gen"
	co "openswarm/internal/coordinator"
	se "openswarm/internal/session"
	wk "openswarm/internal/worker/agent"

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
	pb.RegisterWorkerServiceServer(
		grpc_server,
		wk.NewGRPCServer(),
	)
	pb.RegisterSessionServiceServer(
		grpc_server,
		se.NewGRPCServer(),
	)

	grpc_server.Serve(lis)

}
