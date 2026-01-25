FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY order ./order
COPY account ./account
COPY catalog ./catalog

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/order-server ./order/cmd/order

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/order-server .

# Create non-root user
RUN addgroup -g 1000 app && adduser -D -u 1000 -G app app
USER app

EXPOSE 50053

CMD ["./order-server"]
