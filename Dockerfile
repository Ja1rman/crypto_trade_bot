# syntax=docker/dockerfile:1

FROM golang:1.24.2-bullseye AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o crypto_trading main.go

FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y libc6 ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/crypto_trading .

EXPOSE 2112

ENTRYPOINT ["/app/crypto_trading"]
