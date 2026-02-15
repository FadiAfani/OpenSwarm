package coordinator

import (
	"context"
	"sync"
	"time"

	api "openswarm/api/gen"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Store struct {
	mu      sync.RWMutex
	tasks   map[string]*Task
	workers map[string]*Worker
}

type WorkerStatus struct {
}

type Worker struct {
	ID       uuid.UUID
	GpuModel string
	Vram     uint
	Status   WorkerStatus
}

func NewStore() *Store {
	return &Store{
		tasks:   make(map[string]*Task),
		workers: make(map[string]*Worker),
	}
}

type CoordinatorGRPCServer struct {
	api.UnimplementedCoordinatorServiceServer
	mu      sync.RWMutex
	tasks   map[string][]*Task
	workers map[string]*Worker
	sp      SamplingPlan
}

func NewCoordinatorGRPCServer() *CoordinatorGRPCServer {
	return &CoordinatorGRPCServer{
		mu:      sync.RWMutex{},
		tasks:   make(map[string][]*Task),
		workers: make(map[string]*Worker),
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

	return &api.RelayResponse{}, nil
}

func (s *CoordinatorGRPCServer) GetCoordinatorHealth(ctx context.Context, req *api.GetCoordinatorHealthRequest) (*api.CoordinatorHealthResponse, error) {
	_ = ctx
	_ = req
	return &api.CoordinatorHealthResponse{Status: true}, nil
}

func (s *CoordinatorGRPCServer) RegisterWorker(ctx context.Context, req *api.RegisterWorkerRequest) (*api.RegisterWorkerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	if err := ValidateRegisterWorkerRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.mu.Lock()
	worker := Worker{
		ID:       uuid.MustParse(req.WorkerId),
		GpuModel: req.GpuModel,
		Vram:     uint(req.Vram),
		Status:   WorkerStatus{},
	}
	s.workers[req.WorkerId] = &worker
	s.mu.Unlock()
	return &api.RegisterWorkerResponse{Accepted: true}, nil
}

func (s *CoordinatorGRPCServer) UnregisterWorker(ctx context.Context, req *api.UnregisterWorkerRequest) (*api.UnregisterWorkerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	if err := ValidateUnregisterWorkerRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.mu.Lock()
	delete(s.workers, req.WorkerId)
	s.mu.Unlock()
	return &api.UnregisterWorkerResponse{}, nil
}

func (s *CoordinatorGRPCServer) ListWorkers(ctx context.Context, req *api.ListWorkersRequest) (*api.ListWorkersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	if err := ValidateListWorkersRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.mu.RLock()
	workers := make([]*api.Worker, 0, len(s.workers))
	for _, worker := range s.workers {
		workers = append(workers, &api.Worker{
			Id:       worker.ID.String(),
			GpuModel: worker.GpuModel,
			Vram:     int32(worker.Vram),
		})
	}
	s.mu.RUnlock()
	return &api.ListWorkersResponse{Workers: workers}, nil
}
