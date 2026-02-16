package agent

import (
	"errors"

	api "openswarm/api/gen"
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
	_ = req
	return nil
}
