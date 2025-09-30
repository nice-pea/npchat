# Этап сборки исполняемого файла внутри контейнера
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION="dev"
ARG COMMIT="unknown"

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -ldflags "-X main.version=${VERSION} -X main.buildDate=$(date +"%Y-%m-%dT%H:%M:%SZ") -X main.commit=${COMMIT}" -o npchat cmd/npchat/main.go

# Финальный образ
FROM alpine:3.22.1
WORKDIR /app
COPY --from=builder /app/npchat ./
# Для возможности предевать аргументы не указывая исполняемый файл
ENTRYPOINT ["./npchat"]
# Запуск приложения при старте контейнера
CMD ["./npchat"]
