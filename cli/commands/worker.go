package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	pb "openswarm/api/gen"
	"openswarm/internal/models"
	"openswarm/internal/worker/agent"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WorkerStartCmd struct {
	port *string
}

type RegisterWorkerCmd struct {
	coAddr   string
	gpuModel *string
	vram     *int32
}

type ListWorkersCmd struct {
	coAddr string
}

func ParseWorkerStartCmd(in string) (WorkerStartCmd, error) {
	cmd := WorkerStartCmd{}
	tokens := strings.Fields(in)
	if len(tokens) < 3 {
		return cmd, fmt.Errorf("command is too short")
	}

	if tokens[2] != "start" {
		return cmd, fmt.Errorf("expected 'start' but got %q", tokens[2])
	}

	for i := 3; i < len(tokens); i++ {
		switch tokens[i] {
		case "--port":
			if i+1 >= len(tokens) {
				return cmd, fmt.Errorf("missing value for --port")
			}
			port := tokens[i+1]
			cmd.port = &port
			i++
		default:
			return cmd, fmt.Errorf("unknown option %q", tokens[i])
		}
	}

	return cmd, nil
}

func ParseRegisterWorkerCmd(in string) (RegisterWorkerCmd, error) {
	cmd := RegisterWorkerCmd{}
	tokens := strings.Fields(in)
	if len(tokens) < 3 {
		return cmd, errors.New("command is too short")
	}

	if tokens[2] != "register" {
		return cmd, fmt.Errorf("expected 'register' but got %q", tokens[2])
	}

	for i := 3; i < len(tokens); i++ {
		switch tokens[i] {
		case "--coordinator-addr":
			if len(tokens) <= i+1 {
				return cmd, fmt.Errorf("missing value for --coordinator-addr")
			}
			cmd.coAddr = tokens[i+1]
			i++
		case "--gpu-model":
			if len(tokens) <= i+1 {
				return cmd, fmt.Errorf("missing value for --gpu-model")
			}
			v := tokens[i+1]
			cmd.gpuModel = &v
			i++
		case "--vram":
			if len(tokens) <= i+1 {
				return cmd, fmt.Errorf("missing value for --vram")
			}
			var n int32
			if _, err := fmt.Sscanf(tokens[i+1], "%d", &n); err != nil {
				return cmd, fmt.Errorf("invalid --vram: %w", err)
			}
			cmd.vram = &n
			i++
		default:
			return cmd, fmt.Errorf("unknown option %q", tokens[i])
		}
	}

	return cmd, nil
}

func ParseListWorkersCmd(in string) (ListWorkersCmd, error) {
	cmd := ListWorkersCmd{}
	tokens := strings.Fields(in)
	if len(tokens) < 3 {
		return cmd, errors.New("command is too short")
	}

	for i := 3; i < len(tokens); i++ {
		switch tokens[i] {
		case "--coordinator-addr":
			if len(tokens) <= i+1 {
				return cmd, fmt.Errorf("missing value for --coordinator-addr")
			}
			cmd.coAddr = tokens[i+1]
			i++
		default:
			return cmd, fmt.Errorf("unknown option %q", tokens[i])
		}
	}

	return cmd, nil
}

func RunWorkerStartCmd(cmd WorkerStartCmd) {
	port := "8082"
	if cmd.port != nil {
		port = *cmd.port
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterWorkerServiceServer(server, agent.NewGRPCServer())

	server.Serve(lis)

}

func RunRegisterWorkerCmd(cmd RegisterWorkerCmd) {
	conn, err := grpc.NewClient(cmd.coAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to coordinator: %v", err)
		defer conn.Close()
	}

	client := pb.NewWorkerServiceClient(conn)
	req := &pb.RegisterWorkerRequest{
		WorkerId: string(models.NewUUID()),
		GpuModel: "unspecified",
		Vram:     0,
	}
	if cmd.gpuModel != nil {
		req.GpuModel = *cmd.gpuModel
	}
	if cmd.vram != nil {
		req.Vram = *cmd.vram
	}
	resp, err := client.RegisterWorker(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to register worker: %v", err)
	}
	if !resp.Accepted {
		log.Fatalf("Coordinator rejected worker registration")
	}

	fmt.Println("Worker registered successfully")
}

func RunListWorkersCmd(cmd ListWorkersCmd) {
	conn, err := grpc.NewClient(cmd.coAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to coordinator: %v", err)
		defer conn.Close()
	}

	client := pb.NewWorkerServiceClient(conn)
	req := &pb.ListWorkersRequest{}

	resp, err := client.ListWorkers(context.Background(), req)

	if err != nil {
		log.Fatalf("Failed to list workers: %v", err)
	}

	workers := resp.Workers
	fmt.Printf("Workers (%d)\n", len(workers))
	if len(workers) == 0 {
		fmt.Println("No workers are currently registered.")
		return
	}

	for i, w := range workers {
		layerRange := "unassigned"
		if w.Layers != nil {
			layerRange = fmt.Sprintf("%d-%d", w.Layers.First, w.Layers.Last)
		}

		fmt.Printf(
			"\n[%d] %s\n  Status: %s\n  GPU:    %s\n  VRAM:   %d\n  Layers: %s\n",
			i+1,
			w.Id,
			map[bool]string{true: "online", false: "offline"}[w.OnlineStatus],
			w.GpuModel,
			w.Vram,
			layerRange,
		)
	}
}
