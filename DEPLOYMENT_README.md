# CodeSourcerer Kubernetes Deployment

This directory contains Docker images and Kubernetes deployment configurations for the CodeSourcerer services.

## Services

1. **Database Service** (Port 8080) - gRPC server with Redis backend
2. **Gen-AI Service** (Port 8081) - gRPC server with Gemini API integration
3. **GitHub Service** (Port 3000) - HTTP REST API with Gin framework

## Prerequisites

Before deploying, ensure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/)
- [k3d](https://k3d.io/v5.4.6/#installation)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)

## Quick Start

### Option 1: Automated Deployment (Recommended)

#### Linux/macOS:
```bash
chmod +x build-and-deploy.sh
./build-and-deploy.sh
```

#### Windows (PowerShell):
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
.\build-and-deploy.ps1
```

### Option 2: Manual Deployment

1. **Create k3d cluster:**
```bash
k3d cluster create codesourcerer --port "3000:30000@loadbalancer" --port "8080:30080@loadbalancer" --port "8081:30081@loadbalancer"
```

2. **Build Docker images:**
```bash
# Database service
cd services/database
docker build -t codesourcerer/database:latest .
cd ../..

# Gen-AI service
cd services/gen-ai
docker build -t codesourcerer/gen-ai:latest .
cd ../..

# GitHub service
cd services/github
docker build -t codesourcerer/github:latest .
cd ../..
```

3. **Import images to k3d:**
```bash
k3d image import codesourcerer/database:latest -c codesourcerer
k3d image import codesourcerer/gen-ai:latest -c codesourcerer
k3d image import codesourcerer/github:latest -c codesourcerer
```

4. **Deploy to Kubernetes:**
```bash
kubectl apply -f deployment.yaml
```

5. **Wait for deployments:**
```bash
kubectl wait --for=condition=available --timeout=300s deployment/redis
kubectl wait --for=condition=available --timeout=300s deployment/database-service
kubectl wait --for=condition=available --timeout=300s deployment/genai-service
kubectl wait --for=condition=available --timeout=300s deployment/github-service
```

## Configuration

### Environment Variables

The deployment uses ConfigMaps and Secrets for configuration:

#### ConfigMap (app-config):
- `DATABASE_PORT`: 8080
- `GENAI_PORT`: 8081
- `GITHUB_PORT`: 3000
- `DATABASE_URL`: redis://redis-service:6379

#### Secrets (app-secrets):
Before deploying, update the base64 encoded values in `deployment.yaml`:

```bash
# Encode your secrets
echo -n "your-actual-gemini-api-key" | base64
echo -n "your-actual-pat-token" | base64
echo -n "your-actual-app-id" | base64
echo -n "your-actual-installation-id" | base64
echo -n "your-actual-bot-email" | base64
```

Replace the placeholder values in the Secret section of `deployment.yaml`.

## Accessing Services

- **GitHub Service**: http://localhost:3000
  - Webhook endpoint: http://localhost:3000/webhook
  - Test endpoints: 
    - http://localhost:3000/testsend
    - http://localhost:3000/testfinalizer

- **Database Service**: Internal gRPC service (port 8080)
- **Gen-AI Service**: Internal gRPC service (port 8081)
- **Redis**: Internal service (port 6379)

## Monitoring and Troubleshooting

### Check deployment status:
```bash
kubectl get pods
kubectl get services
kubectl get deployments
```

### View logs:
```bash
# GitHub service logs
kubectl logs -f deployment/github-service

# Database service logs
kubectl logs -f deployment/database-service

# Gen-AI service logs
kubectl logs -f deployment/genai-service

# Redis logs
kubectl logs -f deployment/redis
```

### Debug pod issues:
```bash
# Describe a pod
kubectl describe pod <pod-name>

# Get events
kubectl get events --sort-by=.metadata.creationTimestamp
```

## Cleanup

To remove the deployment:
```bash
kubectl delete -f deployment.yaml
k3d cluster delete codesourcerer
```

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   GitHub App    │    │   Gen-AI API    │    │   Database      │
│   (Port 3000)   │◄──►│   (Port 8081)   │◄──►│   (Port 8080)   │
│                 │    │                 │    │                 │
│ • Webhook       │    │ • gRPC Server   │    │ • gRPC Server   │
│ • REST API      │    │ • Gemini API    │    │ • Redis Client  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
                                               ┌─────────────────┐
                                               │     Redis       │
                                               │   (Port 6379)   │
                                               │                 │
                                               │ • Key-Value     │
                                               │ • Caching       │
                                               └─────────────────┘
```

## Notes

- All services run as non-root users for security
- Health checks are configured for all services
- Resource limits are set to prevent resource exhaustion
- Services communicate internally via Kubernetes DNS
- Only the GitHub service is exposed externally via LoadBalancer
