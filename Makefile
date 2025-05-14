test:
	go test ./...

lint:
	go vet ./... && golangci-lint run -v -j $(( $(nproc) - 1))

run:
	go run ./cmd/npc/main.go