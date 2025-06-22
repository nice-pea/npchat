.PHONY: test vet lint check run build

# Get number of CPU cores minus 1 for parallel execution
CORES := $(shell echo $$(( $$(nproc) - 1 )))

test:
	echo "Running go test..."
	go test ./...

vet:
	echo "Running go vet..."
	go vet ./...

lint:
	echo "Running golangci-lint with $(CORES) workers..."
# The -j parameter for golangci-lint will use all available CPU cores minus one (to avoid overloading your system)
	golangci-lint run -v -j $(CORES)

# Combined target to run both vet and lint
check: vet lint

run:
	echo "Running npc main package..."
	go run ./cmd/npchat/main.go

build:
	echo "Running npchat main package..."
	CGO_ENABLED=0 GOOS=linux go build -o bin/npchat cmd/npchat/main.go