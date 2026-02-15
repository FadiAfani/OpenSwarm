package coordinator

import (
	"errors"

	api "openswarm/api/gen"

	"github.com/google/uuid"
)

func ValidateRegisterWorkerRequest(req *api.RegisterWorkerRequest) error {
	if req.WorkerId == "" {
		return errors.New("worker_id is required")
	}
	if req.GpuModel == "" {
		return errors.New("gpu_model is required")
	}
	if req.Vram == 0 {
		return errors.New("vram is required")
	}
	if req.CoordinatorId == "" {
		return errors.New("coordinator_id is required")
	}
	return nil
}

func ValidateUnregisterWorkerRequest(req *api.UnregisterWorkerRequest) error {
	if req.WorkerId == "" {
		return errors.New("worker_id is required")
	}
	return nil
}

func ValidateGetWorkerRequest(req *api.GetWorkerRequest) error {
	if req.WorkerId == "" {
		return errors.New("worker_id is required")
	}
	return nil
}

func ValidateListWorkersRequest(req *api.ListWorkersRequest) error {
	return nil
}

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
