# Neo-Janus: Quick Commands Reference

## 🚀 Build & Run

### Local Development
```bash
# Build backend
cd 3_janus_core
go mod tidy
go build -o bin/server.exe ./cmd/server/

# Run server
./bin/server.exe

# Or with go run (auto-rebuild)
go run ./cmd/server/main.go
```

### With Make
```bash
# Build
make build

# Run
make run

# All-in-one
make all
```

### Docker
```bash
# Build image
make docker-build

# Start services
make docker-up

# Stop services
make docker-down

# View logs
make docker-logs
```

---

## 🧪 Testing

### Run All Tests
```bash
cd 3_janus_core

# Run all tests
go test -v ./internal/...

# Run specific package tests
go test -v ./internal/api
go test -v ./internal/vaccine
go test -v ./internal/logger

# Run with coverage
go test -v -coverprofile=coverage.out ./internal/...
go tool cover -html=coverage.out
```

### With Make
```bash
make test        # Run tests
make test-cover  # With coverage report
make bench       # Benchmarks
```

---

## 🔴 Red Team Attack

### Launch Attack Campaign
```bash
cd 2_red_agent

# Basic attack (10 payloads)
python auto_attack.py

# Custom intensity
python auto_attack.py 50
python auto_attack.py 100

# With make
make attack
make attack-intense
```

---

## 📊 Monitor & Debug

### Health Check
```bash
curl http://localhost:8080/health
```

### API Endpoint
```bash
# Analyze request
curl -X POST http://localhost:8080/api/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "input": "malicious test input",
    "source": "USER"
  }'

# Response example
{
  "status": "BLOCKED",
  "risk_score": 0.95,
  "message": "🛡️ NEO-JANUS: Content blocked."
}
```

### View Logs
```bash
# Real-time logs
tail -f data/logs/core.log

# With Docker
docker-compose logs -f janus-core

# With make
make docker-logs
```

### Check Vaccine Data
```bash
# List generated patches
ls -la data/vaccine/

# View patch content
cat data/vaccine/vaccine_*.json | jq .
```

---

## 🛠️ Development Workflow

### Code Formatting
```bash
cd 3_janus_core

# Format Go code
go fmt ./...

# With make
make fmt
```

### Linting
```bash
cd 3_janus_core

# Using go vet
go vet ./...

# Or with golangci-lint (if installed)
golangci-lint run ./...

# With make
make lint
```

### Clean Build
```bash
# Remove artifacts
make clean

# Full clean + rebuild
make clean && make build
```

---

## 📝 File Management

### Project Structure
```
Neo-Janus/
├── 1_blue_sentinel/         # Defense AI model
├── 2_red_agent/             # Attack simulator
│   └── attack_lib/
│       ├── fuzzer.py        # Payload generator
│       └── payloads.txt     # Base attack patterns
├── 3_janus_core/            # Backend (Go)
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── api/             # HTTP handlers
│   │   ├── vaccine/         # Vaccine manager
│   │   └── logger/          # Logging system
│   └── go.mod
├── 4_frontend/              # Dashboard (Streamlit)
├── data/
│   ├── logs/                # Application logs
│   └── vaccine/             # Vaccine patches
├── config.yaml              # Configuration
├── Makefile                 # Build automation
├── Dockerfile               # Container image
└── docker-compose.yml       # Multi-container setup
```

### Create Directories
```bash
mkdir -p data/logs data/vaccine
```

---

## 🔧 Configuration

### Edit config.yaml
```yaml
server:
  port: 8080

vaccine:
  enabled: true
  trigger_count: 5

logging:
  level: "info"
  output_dir: "./data/logs/"
```

### Environment Variables
```bash
# (Currently using config.yaml, future support planned)
export LOG_LEVEL=debug
export SERVER_PORT=8000
```

---

## 🐛 Troubleshooting

### Port Already In Use
```bash
# Check what's using port 8080
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Kill process (graceful via SIGTERM)
kill -TERM <PID>

# Or use different port
sed -i 's/port: 8080/port: 9090/' config.yaml
```

### Build Failed
```bash
# Clean and rebuild
cd 3_janus_core
go clean
go mod tidy
go build -o bin/server.exe ./cmd/server/
```

### Tests Failing
```bash
# Run with verbose output
go test -v -run TestName ./internal/package

# Check for race conditions
go test -race ./internal/...

# View coverage
go tool cover -html=coverage.out
```

### Docker Issues
```bash
# Rebuild image
docker-compose down
docker system prune -a
make docker-build
make docker-up
```

---

## 📊 Performance Tips

### Optimize Build
```bash
# Release build (smaller binary)
cd 3_janus_core
CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/server.exe ./cmd/server/
```

### Profile Performance
```bash
# CPU profile
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./internal/...
go tool pprof cpu.prof

# Memory profile
go tool pprof mem.prof
```

---

## 📚 Documentation Files

- **README.md** - Project overview
- **BEST_PRACTICES.md** - Code patterns & roadmap
- **OPTIMIZATION_SUMMARY.md** - Recent improvements
- **Makefile** - Build commands
- **config.yaml** - Application configuration

---

## 🔗 Useful Links

- Go Official: https://golang.org
- Python Best Practices: https://pep8.org
- REST API: https://restfulapi.net
- Docker: https://docker.com
- Security: https://owasp.org

