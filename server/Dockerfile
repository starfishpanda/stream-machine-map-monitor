FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN go build -o machine-map-server server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/stream-machine-map-monitor-server .

EXPOSE 50051

CMD ["./stream-machine-map-monitor-server"]