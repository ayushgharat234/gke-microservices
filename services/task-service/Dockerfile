# Build Stage
FROM golang:1.23.3-alpine AS builder

WORKDIR /app
COPY . .  
RUN go mod tidy && go build -o task-service ./cmd

# Run Stage 
FROM alpine:latest 
RUN apk --no-cache add ca-certificates

WORKDIR /app 
COPY --from=builder /app/task-service .

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
USER appuser

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./task-service"]