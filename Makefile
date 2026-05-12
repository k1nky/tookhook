.PHONY: build clean test lint run docker-build docker-up docker-down

# Binary
BINARY_NAME=tookhook
MAIN_PATH=./cmd/tookhook

# Build flags
VERSION?=dev
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-w -s -X main.buildVersion=$(VERSION) -X main.buildDate=$(BUILD_DATE) -X main.buildCommit=$(GIT_COMMIT)"

# Docker
DOCKER_IMAGE=tookhook2
DOCKER_TAG?=latest

# Build the binary
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

# View test coverage
coverage: test
	go tool cover -html=coverage.out

# Run linter
lint:
	golangci-lint run ./...

# Run the application locally
run: build
	./$(BINARY_NAME) serve --config config/config.example.yaml

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Docker build
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Start with docker-compose
docker-up:
	docker-compose up -d

# Stop docker-compose
docker-down:
	docker-compose down

# View logs
docker-logs:
	docker-compose logs -f

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Check for issues
check: fmt test lint

# Development setup
dev: deps build