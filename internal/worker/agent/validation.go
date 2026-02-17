package agent

import (
	"errors"

	api "openswarm/api/gen"
)

func ValidateRegisterWorkerRequest(req *api.RegisterWorkerRequest) error {
	if req.WorkerId == "" {
		return errors.New("worker_id is required")
	}
	// gpu_model, vram, and coordinator_id are optional
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
	_ = req
	return nil
}
