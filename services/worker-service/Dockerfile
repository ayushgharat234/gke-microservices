# Build Stage
FROM golang:1.23.3-alpine AS builder 
WORKDIR /app 
COPY . .
RUN go mod tidy && go build -o work-service ./cmd 

# Run Stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /app/work-service .

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup
USER appuser

# Note: Worker service doesn't expose HTTP port, so no health check needed
# It's a background service that processes Redis queue

CMD ["./work-service"] 