# Build stage
FROM golang:1.24.2-alpine AS builder
LABEL authors="obashion"

WORKDIR /build

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app ./app

WORKDIR /build/app
RUN go build -o /build/SyncBeat ./main.go


# Final stage
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /build/SyncBeat /app/SyncBeat
ENV GIN_MODE=release

# 実行
ENTRYPOINT ["/app/SyncBeat"]

EXPOSE 8000