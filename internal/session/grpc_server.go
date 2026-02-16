package session

import (
	"context"
	"sync"
	"time"

	api "openswarm/api/gen"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type pipelineRecord struct {
	Workers []uuid.UUID
}

type record struct {
	ID         uuid.UUID
	ConsumerID string
	ModelID    string
	Pipeline   pipelineRecord
	CreatedAt  time.Time
}

type GRPCServer struct {
	api.UnimplementedSessionServiceServer
	mu       sync.RWMutex
	sessions map[string]*record
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{
		mu:       sync.RWMutex{},
		sessions: make(map[string]*record),
	}
}

func (s *GRPCServer) CreateSession(ctx context.Context, req *api.CreateSessionRequest) (*api.CreateSessionResponse, error) {
	_ = ctx
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	if req.ConsumerId == "" {
		return nil, status.Error(codes.InvalidArgument, "consumer_id is required")
	}
	if req.ModelId == "" {
		return nil, status.Error(codes.InvalidArgument, "model_id is required")
	}

	sessionUUID := uuid.New()
	sessionID := sessionUUID.String()

	var workers []uuid.UUID
	if req.Pipeline != nil {
		workers = make([]uuid.UUID, 0, len(req.Pipeline.WorkerIds))
		for _, workerID := range req.Pipeline.WorkerIds {
			parsedID, err := uuid.Parse(workerID)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "pipeline.worker_ids contains invalid UUID")
			}
			workers = append(workers, parsedID)
		}
	}

	created := record{
		ID:         sessionUUID,
		ConsumerID: req.ConsumerId,
		ModelID:    req.ModelId,
		Pipeline: pipelineRecord{
			Workers: workers,
		},
		CreatedAt: time.Now().UTC(),
	}

	s.mu.Lock()
	s.sessions[sessionID] = &created
	s.mu.Unlock()

	return &api.CreateSessionResponse{SessionId: sessionID}, nil
}

func (s *GRPCServer) DestroySession(ctx context.Context, req *api.DestroySessionRequest) (*api.DestroySessionResponse, error) {
	_ = ctx
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	s.mu.Lock()
	_, existed := s.sessions[req.SessionId]
	if existed {
		delete(s.sessions, req.SessionId)
	}
	s.mu.Unlock()

	return &api.DestroySessionResponse{Destroyed: existed}, nil
}

func (s *GRPCServer) GetSession(ctx context.Context, req *api.GetSessionRequest) (*api.GetSessionResponse, error) {
	_ = ctx
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "missing request")
	}
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	s.mu.RLock()
	session, ok := s.sessions[req.SessionId]
	s.mu.RUnlock()
	if !ok {
		return &api.GetSessionResponse{Found: false}, nil
	}

	workerIDs := make([]string, 0, len(session.Pipeline.Workers))
	for _, workerID := range session.Pipeline.Workers {
		workerIDs = append(workerIDs, workerID.String())
	}

	return &api.GetSessionResponse{
		Found:      true,
		SessionId:  session.ID.String(),
		ConsumerId: session.ConsumerID,
		ModelId:    session.ModelID,
		Pipeline: &api.Pipeline{
			WorkerIds: workerIDs,
		},
	}, nil
}
