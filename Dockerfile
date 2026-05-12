# Build stage
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /tookhook ./cmd/tookhook

# Final stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /tookhook /app/tookhook

# Copy default config
COPY config/config.example.yaml /app/config/config.yaml

# Create non-root user
RUN adduser -D -u 1000 tookhook
USER tookhook

EXPOSE 8080

ENTRYPOINT ["/app/tookhook"]
CMD ["serve", "--config", "/app/config/config.yaml"]