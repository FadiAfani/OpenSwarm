package agent

import (
	"context"
	api "openswarm/api/gen"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCWorkerServer struct {
	api.UnimplementedWorkerServiceServer

	mu sync.RWMutex
	co api.CoordinatorServiceClient
}

func NewGRPCWorkerServer() *GRPCWorkerServer {
	return &GRPCWorkerServer{
		mu: sync.RWMutex{},
	}
}

func (s *GRPCWorkerServer) ConnectToCoordinator(ctx context.Context, req *api.ConnectToCoordinatorRequest) (*api.ConnectToCoordinatorResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Unimplemented")
}
