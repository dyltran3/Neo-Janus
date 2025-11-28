.PHONY: help build run test clean docker-build docker-up docker-down lint fmt

# Variables
GO_PROJECT=3_janus_core
GO_BIN=$(GO_PROJECT)/bin/server.exe
GO_CMD=go

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

# --- Go Backend ---

build: ## Build Go backend
	@echo "Building Neo-Janus Core..."
	@cd $(GO_PROJECT) && $(GO_CMD) mod tidy && $(GO_CMD) build -o bin/server.exe ./cmd/server/
	@echo "✅ Build complete: $(GO_BIN)"

run: build ## Build and run the server
	@echo "Starting Neo-Janus Core server..."
	@./$(GO_BIN)

test: ## Run all tests
	@echo "Running tests..."
	@cd $(GO_PROJECT) && $(GO_CMD) test -v ./internal/...
	@echo "✅ Tests passed!"

test-cover: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@cd $(GO_PROJECT) && $(GO_CMD) test -v -coverprofile=coverage.out ./internal/...
	@cd $(GO_PROJECT) && $(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: $(GO_PROJECT)/coverage.html"

bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@cd $(GO_PROJECT) && $(GO_CMD) test -bench=. -benchmem ./internal/...

fmt: ## Format Go code
	@echo "Formatting code..."
	@cd $(GO_PROJECT) && $(GO_CMD) fmt ./...
	@echo "✅ Code formatted"

lint: ## Run linter
	@echo "Running linter..."
	@cd $(GO_PROJECT) && golangci-lint run ./... || $(GO_CMD) vet ./...

# --- Python Red Agent ---

attack: ## Run Red Agent attack campaign
	@echo "Launching Red Team attack campaign..."
	@cd 2_red_agent && python auto_attack.py 10

attack-intense: ## Run intensive attack campaign
	@echo "Launching intensive attack campaign..."
	@cd 2_red_agent && python auto_attack.py 100

# --- Docker ---

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t neo-janus:latest .
	@echo "✅ Docker image built"

docker-up: ## Start services with docker-compose
	@echo "Starting services..."
	@docker-compose up -d
	@echo "✅ Services started"
	@echo "API: http://localhost:8080"

docker-down: ## Stop services
	@echo "Stopping services..."
	@docker-compose down
	@echo "✅ Services stopped"

docker-logs: ## Show Docker logs
	@docker-compose logs -f janus-core

# --- Utilities ---

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(GO_PROJECT)/bin/
	@rm -rf $(GO_PROJECT)/coverage.out $(GO_PROJECT)/coverage.html
	@rm -rf ./data/logs/*.log ./data/vaccine/*.json
	@echo "✅ Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@cd $(GO_PROJECT) && $(GO_CMD) mod download
	@echo "✅ Dependencies downloaded"

dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	@cd $(GO_PROJECT) && $(GO_CMD) mod download
	@mkdir -p data/logs data/vaccine
	@echo "✅ Development environment ready"

serve: run ## Alias for 'run'

# --- Full workflow ---

all: clean build test ## Full build and test cycle

start: docker-up ## Alias for 'docker-up'

stop: docker-down ## Alias for 'docker-down'

full-test: test attack ## Run all tests including attack simulation

# --- Info ---

info: ## Show project info
	@echo "╔══════════════════════════════════════╗"
	@echo "║      Neo-Janus Project Info          ║"
	@echo "╚══════════════════════════════════════╝"
	@echo "Language:  Go + Python"
	@echo "Backend:   $(GO_PROJECT)"
	@echo "Binary:    $(GO_BIN)"
	@echo "API:       http://localhost:8080"
	@echo "Health:    http://localhost:8080/health"
	@echo "Analyze:   POST http://localhost:8080/api/analyze"
	@echo ""
	@echo "Commands:"
	@echo "  make build     - Build backend"
	@echo "  make run       - Run server"
	@echo "  make test      - Run tests"
	@echo "  make attack    - Run attack campaign"
	@echo "  make docker-up - Start with Docker"
	@echo ""

version: ## Show version info
	@echo "Go version: $$(go version)"
	@python --version 2>&1 || echo "Python not found"
	@docker --version 2>&1 || echo "Docker not found"
