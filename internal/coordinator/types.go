package coordinator

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	SessionID string
	Prompt    string
	CreatedAt time.Time
}

type SamplingStrategy int32

const (
	SamplingStrategyUnspecified SamplingStrategy = iota
	SamplingStrategyGreedy
	SamplingStrategyTopP
	SamplingStrategyTopK
	SamplingStrategyTopPTopK
)

type SamplingPlan struct {
	Strategy          SamplingStrategy
	Temperature       float32
	TopP              float32
	TopK              int32
	RepetitionPenalty float32
	MinNewTokens      int32
	MaxNewTokens      int32
	Seed              int32
	Deterministic     bool
}

type Pipeline struct {
	Workers []uuid.UUID
}

type Session struct {
	ID         uuid.UUID
	ConsumerID string
	ModelID    string
	Pipeline   Pipeline
	CreatedAt  time.Time
}
