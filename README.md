# GKE Microservices Architecture

A microservices architecture designed for Google Kubernetes Engine (GKE) with four services working together to process tasks asynchronously.

## Architecture Overview

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │───▶│ API Gateway │───▶│Task Service │
└─────────────┘    └─────────────┘    └─────────────┘
                                              │
                                              ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│Notification │◀───│Worker Service│◀───│   Redis    │
│  Service    │    └─────────────┘    └─────────────┘
└─────────────┘
```

## Services

### 1. API Gateway (Port: 9090)
- Entry point for external requests
- Routes `/create-task` requests to Task Service
- Health endpoint: `/health`

### 2. Task Service (Port: 8080)
- Creates tasks with unique UUIDs
- Stores tasks in Redis queue
- Health endpoint: `/health`
- Readiness endpoint: `/readiness`

### 3. Worker Service
- Background task processor
- Polls Redis queue for tasks
- Processes tasks and sends notifications
- No HTTP endpoints (background service)

### 4. Notification Service (Port: 8083)
- Receives and logs task notifications
- Health endpoint: `/health`
- Notification endpoint: `/notify`

## Docker Setup

### Prerequisites
- Docker
- Docker Compose
- curl (for testing)

### Local Testing with Docker

1. **Build and start all services:**
   ```bash
   docker-compose up --build
   ```

2. **Run the test script:**
   ```bash
   chmod +x test-services.sh
   ./test-services.sh
   ```

3. **Manual testing:**
   ```bash
   # Test health endpoints
   curl http://localhost:9090/health
   curl http://localhost:8080/health
   curl http://localhost:8083/health

   # Create a task
   curl -X POST http://localhost:9090/create-task \
     -H "Content-Type: application/json" \
     -d '{"title": "My Task", "status": "pending"}'
   ```

4. **Check logs:**
   ```bash
   # All services
   docker-compose logs -f

   # Specific service
   docker-compose logs -f api-gateway
   docker-compose logs -f task-service
   docker-compose logs -f worker-service
   docker-compose logs -f notification-service
   ```

### Individual Service Testing

#### API Gateway
```bash
# Build
docker build -t api-gateway ./services/api-gateway

# Run
docker run -p 9090:9090 -e TASK_SERVICE_URL=http://host.docker.internal:8080 api-gateway
```

#### Task Service
```bash
# Build
docker build -t task-service ./services/task-service

# Run (requires Redis)
docker run -p 8080:8080 -e REDIS_ADDR=host.docker.internal:6379 task-service
```

#### Worker Service
```bash
# Build
docker build -t worker-service ./services/worker-service

# Run (requires Redis)
docker run -e REDIS_ADDR=host.docker.internal:6379 worker-service
```

#### Notification Service
```bash
# Build
docker build -t notification-service ./services/notification-service

# Run
docker run -p 8083:8083 notification-service
```

## Docker Features

### ✅ Multi-stage Builds
All services use multi-stage builds to reduce image size:
- Build stage: Compiles Go code
- Run stage: Minimal Alpine image with only the binary

### ✅ Health Checks
All HTTP services include health check endpoints:
- API Gateway: `/health`
- Task Service: `/health`
- Notification Service: `/health`

### ✅ Security
- Non-root user execution
- Minimal Alpine base images
- CA certificates for HTTPS support

### ✅ Local Testing
- Complete docker-compose setup
- Automated test script
- Service dependency management
- Network isolation

## Environment Variables

| Service | Variable | Default | Description |
|---------|----------|---------|-------------|
| API Gateway | `TASK_SERVICE_URL` | `http://localhost:8080` | Task service endpoint |
| Task Service | `REDIS_ADDR` | `localhost:6379` | Redis server address |
| Task Service | `REDIS_PASS` | `` | Redis password |
| Worker Service | `REDIS_ADDR` | `localhost:6379` | Redis server address |
| Worker Service | `REDIS_PASS` | `` | Redis password |

## Troubleshooting

### Common Issues

1. **Services not starting:**
   ```bash
   docker-compose logs [service-name]
   ```

2. **Redis connection issues:**
   - Ensure Redis is running: `docker-compose up redis`
   - Check Redis logs: `docker-compose logs redis`

3. **Health check failures:**
   - Wait for services to fully start
   - Check if ports are available: `netstat -tulpn | grep :8080`

4. **Task processing not working:**
   - Check worker service logs: `docker-compose logs worker-service`
   - Verify Redis queue: `docker exec -it [redis-container] redis-cli LLEN task_queue`

### Performance Testing

```bash
# Load test with multiple tasks
for i in {1..10}; do
  curl -X POST http://localhost:9090/create-task \
    -H "Content-Type: application/json" \
    -d "{\"title\": \"Task $i\", \"status\": \"pending\"}"
done
```

## Next Steps for GKE Deployment

1. Create Kubernetes manifests
2. Set up Redis cluster
3. Configure service mesh (Istio)
4. Add monitoring and logging
5. Set up CI/CD pipeline 