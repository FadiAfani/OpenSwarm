package agent

import (
	"context"
	"sync"

	api "openswarm/api/gen"

	"openswarm/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Status struct{}

type Record struct {
	ID       models.UUID
	GpuModel string
	Vram     uint
	Status   Status
}

type GRPCServer struct {
	api.UnimplementedWorkerServiceServer
	mu      sync.RWMutex
	workers map[models.UUID]*Record
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{
		mu:      sync.RWMutex{},
		workers: make(map[models.UUID]*Record),
	}
}

func (s *GRPCServer) RegisterWorker(ctx context.Context, req *api.RegisterWorkerRequest) (*api.RegisterWorkerResponse, error) {
	_ = ctx
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	if req.WorkerId == "" {
		return nil, status.Error(codes.InvalidArgument, "worker_id is required")
	}
	// gpu_model and vram are optional; use defaults when empty/zero
	gpuModel := req.GpuModel
	if gpuModel == "" {
		gpuModel = "unspecified"
	}

	id := models.NewUUID()

	s.mu.Lock()
	worker := Record{
		ID:       id,
		GpuModel: gpuModel,
		Vram:     uint(req.Vram),
		Status:   Status{},
	}
	s.workers[id] = &worker
	s.mu.Unlock()

	return &api.RegisterWorkerResponse{Accepted: true}, nil
}

func (s *GRPCServer) UnregisterWorker(ctx context.Context, req *api.UnregisterWorkerRequest) (*api.UnregisterWorkerResponse, error) {
	_ = ctx
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	if err := ValidateUnregisterWorkerRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := models.Parse(req.WorkerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "worker_id is not a valid UUID")
	}

	s.mu.Lock()
	delete(s.workers, id)
	s.mu.Unlock()

	return &api.UnregisterWorkerResponse{}, nil
}

func (s *GRPCServer) GetWorker(ctx context.Context, req *api.GetWorkerRequest) (*api.GetWorkerResponse, error) {
	_ = ctx
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	if err := ValidateGetWorkerRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id := models.UUID(req.WorkerId)
	s.mu.RLock()
	worker, ok := s.workers[id]
	s.mu.RUnlock()
	if !ok {
		return &api.GetWorkerResponse{Found: false}, nil
	}

	return &api.GetWorkerResponse{
		Found:    true,
		WorkerId: id.String(),
		GpuModel: worker.GpuModel,
		Vram:     int32(worker.Vram),
	}, nil
}

func (s *GRPCServer) ListWorkers(ctx context.Context, req *api.ListWorkersRequest) (*api.ListWorkersResponse, error) {
	_ = ctx
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

func (s *GRPCServer) Heartbeat(ctx context.Context, req *api.HeartbeatRequest) (*api.HeartbeatResponse, error) {
	_ = ctx
	if req == nil || req.WorkerId == "" {
		return nil, status.Error(codes.InvalidArgument, "worker_id is required")
	}
	return &api.HeartbeatResponse{Accepted: true}, nil
}

func (s *GRPCServer) ConnectToCoordinator(ctx context.Context, req *api.ConnectToCoordinatorRequest) (*api.ConnectToCoordinatorResponse, error) {
	_ = ctx
	if req == nil || req.Addr == "" {
		return nil, status.Error(codes.InvalidArgument, "addr is required")
	}
	return &api.ConnectToCoordinatorResponse{Accepted: true}, nil
}
