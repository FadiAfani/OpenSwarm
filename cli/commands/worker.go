package commands

import (
	"context"
	"flag"
	"fmt"
	"time"

	pb "openswarm/api/gen"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunWorker(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing worker subcommand\n\n%s", workerUsage())
	}

	switch args[0] {
	case "register":
		return runWorkerRegister(args[1:])
	default:
		return fmt.Errorf("unknown worker subcommand %q\n\n%s", args[0], workerUsage())
	}
}

func runWorkerRegister(args []string) error {
	fs := flag.NewFlagSet("worker register", flag.ContinueOnError)

	coordinatorAddr := fs.String("coordinator-addr", "localhost:8080", "Coordinator gRPC address")
	workerID := fs.String("worker-id", "", "Worker UUID (optional, auto-detected if omitted)")
	gpuModel := fs.String("gpu-model", "", "GPU model (optional, auto-detected if omitted)")
	gpuModelAlias := fs.String("gpu", "", "GPU model (alias for --gpu-model)")
	vram := fs.Int("vram", 0, "VRAM in GB (optional, auto-detected if omitted)")
	timeout := fs.Duration("timeout", 5*time.Second, "Request timeout")

	if err := fs.Parse(args); err != nil {
		return err
	}

	hostMeta := detectHostMetadata()

	resolvedWorkerID := pickString(*workerID, hostMeta.WorkerID)
	resolvedGPUModel := pickFirstNonEmpty(*gpuModel, *gpuModelAlias, hostMeta.GPUModel, "unspecified")
	resolvedVRAM := pickPositiveOrDefault(*vram, int(hostMeta.VRAMGB), 1)

	if _, err := uuid.Parse(resolvedWorkerID); err != nil {
		return fmt.Errorf("resolved worker_id must be a valid UUID: %w", err)
	}
	if *coordinatorAddr == "" {
		return fmt.Errorf("--coordinator-addr is required")
	}

	conn, err := grpc.NewClient(*coordinatorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to coordinator: %w", err)
	}
	defer conn.Close()

	client := pb.NewWorkerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	resp, err := client.RegisterWorker(ctx, &pb.RegisterWorkerRequest{
		WorkerId:      resolvedWorkerID,
		GpuModel:      resolvedGPUModel,
		Vram:          int32(resolvedVRAM),
		CoordinatorId: *coordinatorAddr,
	})
	if err != nil {
		return fmt.Errorf("register worker RPC failed: %w", err)
	}

	fmt.Printf(
		"worker registered: accepted=%t worker_id=%s gpu_model=%q vram_gb=%d coordinator_id=%q\n",
		resp.GetAccepted(),
		resolvedWorkerID,
		resolvedGPUModel,
		resolvedVRAM,
		*coordinatorAddr,
	)
	return nil
}

func workerUsage() string {
	return `Usage:
  openswarm-cli worker register [flags]

Flags:
  --coordinator-addr string Coordinator gRPC address (default "localhost:8080")
  --worker-id string        Worker UUID (optional, defaults to host-derived ID)
  --gpu string              GPU model (alias for --gpu-model)
  --gpu-model string        GPU model (optional, auto-detected, fallback "unspecified")
  --vram int                VRAM in GB (optional, auto-detected, fallback 1)
  --timeout duration        Request timeout (default 5s)`
}

func pickFirstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func pickString(primary string, fallback string) string {
	if primary != "" {
		return primary
	}
	return fallback
}

func pickInt(primary int, fallback int) int {
	if primary > 0 {
		return primary
	}
	return fallback
}

func pickPositiveOrDefault(primary int, fallback int, defaultValue int) int {
	if primary > 0 {
		return primary
	}
	if fallback > 0 {
		return fallback
	}
	return defaultValue
}
