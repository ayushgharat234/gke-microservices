# GKE Microservices - Helm Charts

This directory contains Helm charts for deploying a microservices architecture on Google Kubernetes Engine (GKE). The system consists of four main services and a Redis database cluster that work together to process tasks asynchronously.

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │───▶│  Task Service   │───▶│     Redis       │
│     :9090       │    │     :8080       │    │  (Master/Slave) │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                        │                        │
         │                        │                        ▼
         ▼                        ▼              ┌─────────────────┐
┌─────────────────┐    ┌─────────────────┐      │ Worker Service  │
│   External      │    │  Notification   │◀─────│     (Queue      │
│   Client        │    │   Service       │      │   Processor)    │
│                 │    │     :8083       │      └─────────────────┘
└─────────────────┘    └─────────────────┘
```

## Services Description

### 1. API Gateway (`api-gateway`)
- **Port**: 9090 (exposed via LoadBalancer)
- **Purpose**: Entry point for external clients, routes requests to backend services
- **Key Features**:
  - Health check endpoint (`/health`)
  - Request forwarding to task-service (`/create-task`)
  - Load balancing across multiple replicas

### 2. Task Service (`task-service`)
- **Port**: 8080 (internal ClusterIP)
- **Purpose**: Manages task creation and queuing
- **Key Features**:
  - Creates tasks with unique IDs
  - Stores tasks in Redis queue for processing
  - Health and readiness probes
  - Connects to Redis master for persistence

### 3. Worker Service (`worker-service`)
- **Purpose**: Background task processor
- **Key Features**:
  - Polls Redis queue for pending tasks
  - Processes tasks asynchronously
  - Sends completion notifications
  - No exposed ports (internal service only)

### 4. Notification Service (`notification-service`)
- **Port**: 8083 (internal ClusterIP)
- **Purpose**: Handles task completion notifications
- **Key Features**:
  - Receives notifications from worker service
  - Logs task status updates
  - Can be extended for email/SMS notifications

## Data Flow

1. **Task Creation**:
   ```
   Client → API Gateway → Task Service → Redis Queue
   ```

2. **Task Processing**:
   ```
   Redis Queue → Worker Service → Notification Service
   ```

3. **Complete Flow**:
   ```
   POST /create-task → Task queued → Worker processes → Notification sent
   ```

## Networking Configuration

### Service Types and Ports

| Service | Type | Internal Port | External Access | Purpose |
|---------|------|---------------|-----------------|---------|
| api-gateway | LoadBalancer | 9090 | :9090 | External entry point |
| task-service | ClusterIP | 8080 | Internal only | Task management |
| notification-service | ClusterIP | 8083 | Internal only | Notifications |
| worker-service | None | N/A | Background only | Task processing |

### Redis Configuration

The system uses a Redis cluster with the following services:
- `redis-master`: Primary Redis instance for writes (port 6379)
- `redis-replicas`: Read replicas for scaling reads (port 6379)
- `redis-headless`: Headless service for StatefulSet discovery

**Important**: All services connect to `redis-master:6379` for consistency.

### Internal Service Communication

```yaml
# API Gateway → Task Service
taskServiceUrl: "http://task-service:8080"

# Task Service → Redis
redisAddr: "redis-master:6379"

# Worker Service → Redis
redisAddr: "redis-master:6379"

# Worker Service → Notification Service
notificationServiceUrl: "http://notification-service:8083"
```

## Deployment Instructions

### Prerequisites
- Kubernetes cluster (GKE, minikube, Docker Desktop, etc.)
- Helm 3.x installed
- Redis deployed (using Bitnami Redis chart or similar)

### Deploy Redis (if not already deployed)
```bash
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install redis bitnami/redis
```

### Deploy Services
```bash
# Deploy all services
helm install api-gateway ./charts/api-gateway
helm install task-service ./charts/task-service
helm install notification-service ./charts/notification-service
helm install worker-service ./charts/worker-service
```

### Verify Deployment
```bash
kubectl get pods
kubectl get services
```

## Accessing the Application

### Local Development (Port Forwarding)
```bash
# Forward API Gateway port
kubectl port-forward service/api-gateway 9090:9090

# Test the API
curl -X POST http://localhost:9090/create-task \
  -H "Content-Type: application/json" \
  -d '{"title": "My test task"}'
```

### Production (LoadBalancer)
```bash
# Get external IP
kubectl get service api-gateway

# Access via external IP
curl -X POST http://<EXTERNAL-IP>:9090/create-task \
  -H "Content-Type: application/json" \
  -d '{"title": "Production task"}'
```

## Configuration Parameters

### API Gateway (`api-gateway/values.yaml`)
```yaml
replicaCount: 2                    # Number of replicas
service:
  type: LoadBalancer               # Service type for external access
  port: 9090                       # External port
environment:
  taskServiceUrl: "http://task-service:8080"  # Backend service URL
```

### Task Service (`task-service/values.yaml`)
```yaml
replicaCount: 2                    # Number of replicas
service:
  type: ClusterIP                  # Internal service only
  port: 8080
environment:
  redisAddr: "redis-master:6379"   # Redis connection string
  redisPass: ""                    # Redis password (empty for no auth)
```

### Worker Service (`worker-service/values.yaml`)
```yaml
replicaCount: 2                    # Number of worker instances
environment:
  redisAddr: "redis-master:6379"   # Redis connection
  notificationServiceUrl: "http://notification-service:8083"
```

### Notification Service (`notification-service/values.yaml`)
```yaml
replicaCount: 2                    # Number of replicas
service:
  type: ClusterIP                  # Internal service only
  port: 8083
```

## Health Checks and Monitoring

### Liveness and Readiness Probes
- **Task Service**: 
  - Liveness: `/health` (checks basic service health)
  - Readiness: `/readiness` (checks Redis connectivity)
- **API Gateway**: 
  - Liveness: `/health`

### Monitoring Endpoints
```bash
# Check API Gateway health
curl http://localhost:9090/health

# Check individual service health (via port-forward)
kubectl port-forward service/task-service 8080:8080
curl http://localhost:8080/health
```

## Scaling

### Manual Scaling
```bash
# Scale API Gateway
kubectl scale deployment api-gateway --replicas=3

# Scale Task Service
kubectl scale deployment task-service --replicas=4

# Scale Worker Service for more processing power
kubectl scale deployment worker-service --replicas=5
```

### Auto-scaling (Optional)
Each chart includes HPA configuration that can be enabled:
```yaml
autoscaling:
  enabled: true                    # Enable HPA
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
```

## Troubleshooting

### Common Issues

1. **Pods in CrashLoopBackOff**
   - Check Redis connectivity: `kubectl logs <pod-name>`
   - Verify Redis service is running: `kubectl get services | grep redis`

2. **API Gateway Connection Refused**
   - Check if port-forwarding is active
   - Verify LoadBalancer external IP: `kubectl get service api-gateway`

3. **Tasks Not Processing**
   - Check worker service logs: `kubectl logs deployment/worker-service`
   - Verify Redis queue: Connect to Redis and check `LLEN task_queue`

### Debugging Commands
```bash
# Check pod status
kubectl get pods

# View logs
kubectl logs deployment/task-service
kubectl logs deployment/worker-service
kubectl logs deployment/notification-service

# Check service connectivity
kubectl exec -it <pod-name> -- nslookup redis-master
kubectl exec -it <pod-name> -- curl http://task-service:8080/health

# Monitor Redis
kubectl exec -it redis-master-0 -- redis-cli
> LLEN task_queue
> MONITOR
```

## Security Considerations

### Current Security Features
- **Non-root containers**: All services run as non-root user (UID 1000)
- **Resource limits**: CPU and memory limits prevent resource exhaustion
- **Network policies**: Services communicate only through defined ports

### Production Recommendations
1. Enable Redis authentication (`redisPass`)
2. Use TLS for inter-service communication
3. Implement network policies for pod-to-pod communication
4. Use secrets for sensitive configuration
5. Enable RBAC for service accounts

## Development and Customization

### Adding New Endpoints
1. Modify the API Gateway router to add new routes
2. Update service configurations in `values.yaml`
3. Redeploy using `helm upgrade`

### Extending Worker Processing
1. Modify `worker-service/internal/worker.go`
2. Add business logic for different task types
3. Update notification handling as needed

### Custom Notifications
1. Extend `notification-service/internal/handler.go`
2. Add integrations for email, SMS, webhooks
3. Configure external service credentials

This microservices architecture provides a scalable, resilient foundation for task processing with clear separation of concerns and proper Kubernetes-native deployment patterns.
