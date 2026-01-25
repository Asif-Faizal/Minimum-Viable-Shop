FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY account ./account

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/account-server ./account/cmd/account

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates curl

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/account-server .

# Create non-root user
RUN addgroup -g 1000 app && adduser -D -u 1000 -G app app
USER app

EXPOSE 50051

HEALTHCHECK --interval=15s --timeout=5s --retries=3 --start-period=10s \
    CMD curl -f http://localhost:50051/health || exit 1

CMD ["./account-server"]

