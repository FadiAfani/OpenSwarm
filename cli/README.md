## CLI

Current commands:

- `openswarm-cli coordinator start`
- `openswarm-cli worker register`

Start coordinator:

```bash
go run ./cmd/cli coordinator start --addr localhost:8080
```

Example:

```bash
go run ./cmd/cli worker register \
  --coordinator-addr localhost:8080
```

The CLI auto-detects worker metadata from host (worker ID, GPU model, VRAM).
`--coordinator-addr` is used for dialing and sent as coordinator ID in the RPC.
If GPU detection is unavailable, it falls back to `gpu_model="unspecified"` and `vram=1`.
