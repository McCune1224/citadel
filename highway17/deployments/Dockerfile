# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache git make npm

# Copy module files
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate templ components
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate

# Build CSS with npm
RUN npm install && \
    npm run css

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dashboard ./cmd/dashboard

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary and assets from builder
COPY --from=builder /build/dashboard /app/
COPY --from=builder /build/migrations /app/migrations/
COPY --from=builder /build/web/static /app/web/static/

# Create non-root user
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app && \
    chown -R app:app /app

USER app

EXPOSE 8080

ENV APP_PORT=8080 \
    APP_ENV=production \
    LOG_LEVEL=info

CMD ["./dashboard"]
