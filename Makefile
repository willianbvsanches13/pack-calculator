.PHONY: build run test clean docker-build docker-run dev frontend backend

BINARY_NAME=pack-calculator
DOCKER_IMAGE_FRONTEND=pack-calculator-frontend
DOCKER_IMAGE_BACKEND=pack-calculator-backend
PORT=80

# Build both frontend and backend
build: frontend backend

# Build only backend
backend:
	go build -o $(BINARY_NAME) ./cmd/server

# Build only frontend
frontend:
	cd web && npm ci && npm run build

# Run backend (assumes frontend is built)
run: backend
	./$(BINARY_NAME)

# Development mode - run backend and frontend separately
dev:
	@echo "Run in separate terminals:"
	@echo "  Terminal 1: make dev-backend"
	@echo "  Terminal 2: make dev-frontend"

dev-backend:
	go run ./cmd/server

dev-frontend:
	cd web && npm run dev

# Tests
test:
	go test -v ./...

test-coverage:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

bench:
	go test -bench=. -benchmem ./internal/calculator/

# Lint
lint:
	go vet ./...
	cd web && npm run lint

# Clean
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	rm -rf web/dist
	rm -rf web/node_modules

# Docker - build individual images
docker-build-frontend:
	docker build -t $(DOCKER_IMAGE_FRONTEND) -f Dockerfile.frontend .

docker-build-backend:
	docker build -t $(DOCKER_IMAGE_BACKEND) -f Dockerfile.backend .

docker-build: docker-build-frontend docker-build-backend

# Docker Compose
docker-up:
	docker-compose up --build

docker-up-detached:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f
