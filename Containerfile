# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION="dev"

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -ldflags "-X main.version=${VERSION} -X main.buildDate=$(date +"%Y-%m-%dT%H:%M:%SZ")" -o npchat cmd/npchat/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/npchat ./
CMD ["./npchat"]
ENTRYPOINT ["./npchat"]