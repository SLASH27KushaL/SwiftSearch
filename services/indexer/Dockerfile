# Stage 1: Build the Go Binary
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git tzdata
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o indexer ./cmd/main.go

# Stage 2: Create the Minimal Image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app

COPY --from=builder /app/indexer .
USER appuser

CMD ["./indexer"]