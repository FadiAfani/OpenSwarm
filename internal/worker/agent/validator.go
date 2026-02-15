package agent

import (
	api "openswarm/api/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateConnectToCoordinatorRequest(req *api.ConnectToCoordinatorRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "missing request")
	}

	if req.Addr == "" {
		return status.Error(codes.InvalidArgument, "invalid coordinator address")
	}

	return nil

}
