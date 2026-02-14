package coordinator

import (
	"time"
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
