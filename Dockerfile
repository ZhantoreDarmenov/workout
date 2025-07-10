# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o workout ./cmd

# Runtime stage
FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /app/workout ./workout
COPY config ./config
EXPOSE 4001
ENV CONFIG_PATH=/app/config/config.yaml
CMD ["./workout"]