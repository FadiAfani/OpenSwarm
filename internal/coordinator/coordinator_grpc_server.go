package coordinator

import (
	"context"
	"sync"
	"time"

	api "openswarm/api/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CoordinatorGRPCServer struct {
	api.UnimplementedCoordinatorServiceServer
	mu    sync.RWMutex
	tasks map[string][]*Task
	sp    SamplingPlan
}

func NewCoordinatorGRPCServer() *CoordinatorGRPCServer {
	return &CoordinatorGRPCServer{
		mu:    sync.RWMutex{},
		tasks: make(map[string][]*Task),
		sp: SamplingPlan{
			Strategy:          SamplingStrategyGreedy,
			Temperature:       0,
			TopP:              0,
			TopK:              0,
			RepetitionPenalty: 1.0,
			MinNewTokens:      0,
			MaxNewTokens:      256,
			Seed:              0,
			Deterministic:     true,
		},
	}
}

func (s *CoordinatorGRPCServer) InitializeTask(ctx context.Context, req *api.InitTaskRequest) (*api.InitTaskResponse, error) {
	_ = ctx
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}
	if req.Prompt == "" {
		return nil, status.Error(codes.InvalidArgument, "prompt is required")
	}

	s.mu.Lock()
	task := Task{
		SessionID: req.SessionId,
		Prompt:    req.Prompt,
		CreatedAt: time.Now().UTC(),
	}
	s.tasks[req.SessionId] = append(s.tasks[req.SessionId], &task)
	s.mu.Unlock()

	return &api.InitTaskResponse{}, nil
}

func (s *CoordinatorGRPCServer) RelayForward(ctx context.Context, req *api.RelayRequest) (*api.RelayResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	if err := ValidateRelayRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_ = ctx
	return &api.RelayResponse{}, nil
}

func (s *CoordinatorGRPCServer) GetCoordinatorHealth(ctx context.Context, req *api.GetCoordinatorHealthRequest) (*api.CoordinatorHealthResponse, error) {
	_ = ctx
	_ = req
	return &api.CoordinatorHealthResponse{Status: true}, nil
}
