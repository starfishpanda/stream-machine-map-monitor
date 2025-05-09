FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY proto/ ./proto/

COPY ./ws-proxy/ ./ws-proxy/

WORKDIR /app/ws-proxy
RUN go build -x -o ws-proxy main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ws-proxy/ws-proxy .

EXPOSE 3001

CMD ["./ws-proxy"]