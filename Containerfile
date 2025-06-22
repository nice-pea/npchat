# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o npchat cmd/npchat/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/npchat ./
CMD ["./npchat"]
ENTRYPOINT ["./npchat"]