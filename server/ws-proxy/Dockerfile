FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN go build -o ws-proxy server/ws-proxy/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ws-proxy .

EXPOSE 3001

CMD ["./ws-proxy"]