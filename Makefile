.PHONY: build clean test install lint help

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
GOFLAGS := -v

build:
	@echo "🔨 Building CTX $(VERSION)..."
	@go build $(GOFLAGS) $(LDFLAGS) -o bin/ctx ./cmd/ctx

build-all:
	@echo "🔨 Building for all platforms..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/ctx-linux-amd64 ./cmd/ctx
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/ctx-linux-arm64 ./cmd/ctx
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/ctx-darwin-amd64 ./cmd/ctx
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/ctx-darwin-arm64 ./cmd/ctx
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/ctx-windows-amd64.exe ./cmd/ctx

install:
	@echo "📦 Installing CTX..."
	@go install $(LDFLAGS) ./cmd/ctx

test:
	@echo "🧪 Running tests..."
	@go test -v -cover ./...

test-coverage:
	@echo "📊 Test coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

lint:
	@echo "🔍 Linting..."
	@go vet ./...
	@test -z "$$(gofmt -l .)" || (echo "Files need formatting:"; gofmt -l .; exit 1)

clean:
	@echo "🧹 Cleaning..."
	@rm -rf bin/ coverage.out
	@go clean

bench:
	@echo "⏱️  Running benchmarks..."
	@go test -bench=. -benchmem ./...

deps:
	@echo "📚 Updating dependencies..."
	@go mod tidy
	@go mod verify

release: test build-all
	@echo "🚀 Creating release $(VERSION)..."
	@mkdir -p release
	@cp bin/* release/
	@cd release && sha256sum * > checksums.txt

help:
	@echo "CTX Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build          Build for current platform"
	@echo "  make build-all      Build for all platforms"
	@echo "  make install        Install to GOPATH/bin"
	@echo "  make test           Run tests"
	@echo "  make test-coverage  Run tests with coverage"
	@echo "  make lint           Run linters"
	@echo "  make clean          Clean build files"
	@echo "  make bench          Run benchmarks"
	@echo "  make deps           Update dependencies"
	@echo "  make release        Create release"
	@echo "  make help           Show this help"