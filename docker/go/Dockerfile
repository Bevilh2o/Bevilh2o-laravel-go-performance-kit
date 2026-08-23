# Stage 1: Build binary
FROM golang:alpine AS builder

WORKDIR /app

# Download dependencies
COPY go/go.mod go/go.sum ./
RUN go mod download

# Copy source code
COPY go/ ./

# Compile statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/ingestion ./cmd/ingestion

# Stage 2: Minimal production runtime image
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /bin/ingestion /app/ingestion

EXPOSE 8080

ENTRYPOINT ["/app/ingestion"]