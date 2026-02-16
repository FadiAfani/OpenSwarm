package session

import (
	"sync"
	"time"

	api "openswarm/api/gen"

	"github.com/google/uuid"
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
