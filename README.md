# herdctl

A lightweight CLI for managing multi-service local dev environments via a single config file.

---

## Installation

```bash
go install github.com/yourusername/herdctl@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/herdctl/releases).

---

## Usage

Define your services in a `herd.yaml` file at the root of your project:

```yaml
services:
  api:
    cmd: "go run ./cmd/api"
    port: 8080
  worker:
    cmd: "go run ./cmd/worker"
  postgres:
    image: "postgres:15"
    port: 5432
```

Then manage your environment with a single command:

```bash
# Start all services
herdctl up

# Start specific services
herdctl up api postgres

# Stop all running services
herdctl down

# View aggregated logs
herdctl logs

# Follow logs for a specific service
herdctl logs api --follow

# Check service status
herdctl status
```

---

## Commands

| Command | Description |
|---|---|
| `up [services]` | Start all or specified services |
| `down` | Stop all running services |
| `logs [service]` | Tail logs for all or one service |
| `status` | Show current state of all services |
| `restart [service]` | Restart all or a specific service |

### Flags

| Flag | Command | Description |
|---|---|---|
| `--follow`, `-f` | `logs` | Stream logs continuously |
| `--timeout`, `-t` | `down`, `restart` | Seconds to wait before force-killing (default: 10) |

---

## License

[MIT](LICENSE)
