FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git tzdata
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o search-api ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app

COPY --from=builder /app/search-api .
USER appuser

EXPOSE 8080
CMD ["./search-api"]