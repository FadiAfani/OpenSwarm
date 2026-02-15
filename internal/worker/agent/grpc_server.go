package agent

import (
	"context"
	api "openswarm/api/gen"
	"sync"

	"google.golang.org/grpc"
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
	if err := ValidateConnectToCoordinatorRequest(req); err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(req.Addr)
	if err != nil {
		return nil, err
	}

	s.co = api.NewCoordinatorServiceClient(conn)

	wreq := &api.RegisterWorkerRequest{}

	res, err := s.co.RegisterWorker(ctx, wreq)

	if err != nil {
		return nil, err
	}

	if !res.Accepted {
		return nil, status.Error(codes.Unknown, "registration request was denied")
	}

	return &api.ConnectToCoordinatorResponse{Accepted: true}, nil
}
