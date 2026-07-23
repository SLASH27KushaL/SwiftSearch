# ==========================================
# Stage 1: Build Stage
# ==========================================
FROM golang:alpine AS builder

# Install SSL certificates and build tools
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory inside container
WORKDIR /app

# Copy go module files and download dependencies first for fast layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files from services/spider
COPY . .

# Build a fully static Go binary (CGO_ENABLED=0 disables C libraries)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o spider \
    ./cmd/main.go

# ==========================================
# Stage 2: Minimal Runtime Stage (~25MB Image)
# ==========================================
FROM alpine:3.19

# CA certificates are required to make outbound HTTPS requests during crawling
RUN apk --no-cache add ca-certificates tzdata

# Security best practice: Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy executable from build stage
COPY --from=builder /app/spider .

# Set ownership
RUN chown appuser:appgroup /app/spider

# Run as non-root user
USER appuser

# Environment defaults (can be overridden by docker-compose)
ENV WORKER_COUNT=10 \
    MAX_DEPTH=3 \
    CRAWL_DELAY_MS=500 \
    REDIS_ADDR=redis:6379 \
    MONGO_URI=mongodb://mongo:27017

ENTRYPOINT ["/app/spider"]