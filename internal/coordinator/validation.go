package coordinator

import (
	"errors"

	api "openswarm/api/gen"
)

func validateRegisterWorkerRequest(req *api.RegisterWorkerRequest) error {
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

func validateUnregisterWorkerRequest(req *api.UnregisterWorkerRequest) error {
	if req.WorkerId == "" {
		return errors.New("worker_id is required")
	}
	return nil
}

func validateGetWorkerRequest(req *api.GetWorkerRequest) error {
	if req.WorkerId == "" {
		return errors.New("worker_id is required")
	}
	return nil
}

func validateListWorkersRequest(req *api.ListWorkersRequest) error {
	return nil
}

func validateInitializeTaskRequest(req *api.InitTaskRequest) error {
	if req.SessionId == "" {
		return errors.New("session_id is required")
	}
	if req.Prompt == "" {
		return errors.New("prompt is required")
	}
	return nil
}
