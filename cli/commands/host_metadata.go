package commands

import (
	"encoding/json"
	"math"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type hostMetadata struct {
	WorkerID string
	GPUModel string
	VRAMGB   int32
}

func detectHostMetadata() hostMetadata {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		hostname = "localhost"
	}

	workerID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte(hostname+"-"+runtime.GOOS+"-"+runtime.GOARCH)).String()
	gpuModel, vramGB := detectGPUInfo()
	return hostMetadata{
		WorkerID: workerID,
		GPUModel: gpuModel,
		VRAMGB:   vramGB,
	}
}

func detectGPUInfo() (string, int32) {
	if runtime.GOOS == "darwin" {
		if model, vram := detectGPUInfoDarwin(); model != "" && vram > 0 {
			return model, vram
		}
	}

	if model, vram := detectGPUInfoNvidiaSMI(); model != "" && vram > 0 {
		return model, vram
	}

	return "", 0
}

func detectGPUInfoNvidiaSMI() (string, int32) {
	out, err := exec.Command("nvidia-smi", "--query-gpu=name,memory.total", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return "", 0
	}

	line := strings.TrimSpace(strings.Split(string(out), "\n")[0])
	if line == "" {
		return "", 0
	}
	parts := strings.Split(line, ",")
	if len(parts) < 2 {
		return "", 0
	}

	model := strings.TrimSpace(parts[0])
	mbText := strings.TrimSpace(parts[1])
	mb, err := strconv.ParseFloat(mbText, 64)
	if err != nil {
		return "", 0
	}

	gb := int32(math.Ceil(mb / 1024.0))
	if gb < 1 {
		gb = 1
	}
	return model, gb
}

func detectGPUInfoDarwin() (string, int32) {
	out, err := exec.Command("system_profiler", "SPDisplaysDataType", "-json").Output()
	if err != nil {
		return "", 0
	}

	var payload map[string][]map[string]any
	if err := json.Unmarshal(out, &payload); err != nil {
		return "", 0
	}

	entries := payload["SPDisplaysDataType"]
	for _, entry := range entries {
		model := firstString(entry, "sppci_model", "_name", "spdisplays_vendor")
		vramText := firstString(entry, "spdisplays_vram", "spdisplays_vram_shared", "spdisplays_device-memory")
		vram := parseVRAMToGB(vramText)

		if model == "" {
			continue
		}

		if vram > 0 {
			return model, vram
		}
	}

	return "", 0
}

func firstString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		v, ok := m[k]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if ok {
			s = strings.TrimSpace(s)
			if s != "" {
				return s
			}
		}
	}
	return ""
}

func parseVRAMToGB(s string) int32 {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0
	}

	re := regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)\s*(TB|GB|MB)`) //nolint:gocritic
	matches := re.FindStringSubmatch(s)
	if len(matches) == 3 {
		value, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0
		}
		switch matches[2] {
		case "TB":
			return int32(math.Ceil(value * 1024.0))
		case "GB":
			return int32(math.Ceil(value))
		case "MB":
			return int32(math.Ceil(value / 1024.0))
		}
	}

	if plain, err := strconv.ParseFloat(s, 64); err == nil {
		return int32(math.Ceil(plain / 1024.0))
	}

	return 0
}
