package coordinator

import (
	"errors"

	api "openswarm/api/gen"

	"github.com/google/uuid"
)

func ValidateInitializeTaskRequest(req *api.InitTaskRequest) error {
	if req.SessionId == "" {
		return errors.New("session_id is required")
	}
	if req.Prompt == "" {
		return errors.New("prompt is required")
	}
	return nil
}

func ValidateRelayRequest(req *api.RelayRequest) error {
	if req == nil {
		return errors.New("missing request")
	}
	_, err := uuid.Parse(req.WorkerId)
	if err != nil {
		return errors.New("worker_id is not a valid UUID")
	}
	return nil
}
