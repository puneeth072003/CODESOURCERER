# CodeSourcerer Deployment Guide

> **Complete deployment documentation for CodeSourcerer microservices platform**
> Includes Docker Compose, Kubernetes, Helm, Sealed Secrets, and FluxCD GitOps configurations

## Table of Contents

- [Overview](#overview)
- [Directory Structure](#-directory-structure)
- [Quick Start](#-quick-start)
- [Deployment Methods](#-deployment-methods)
  - [Docker Compose](#option-1-docker-compose-development)
  - [Kubernetes with Helm](#option-2-helm-chart-production)
  - [GitOps with FluxCD](#option-3-fluxcd-gitops-production)
- [Secrets Management](#-secrets-management)
  - [Sealed Secrets Setup](#sealed-secrets-setup)
  - [Creating Sealed Secrets](#creating-sealed-secrets)
  - [Using Sealed Secrets](#using-sealed-secrets)
- [Configuration](#-configuration)
- [Security](#-security)
- [Testing & Monitoring](#-testing--monitoring)
- [Troubleshooting](#-troubleshooting)
- [Cleanup](#-cleanup)

---

## Overview

This deployment supports three primary methods:

1. **Docker Compose** - Local development and testing
2. **Helm** - Production Kubernetes deployments
3. **FluxCD** - GitOps-based continuous deployment

All methods support **Sealed Secrets** for secure secret management in Git repositories.

## 📁 Directory Structure

```
deployment/
├── README.md                        # This comprehensive guide
├── docker-compose/                  # Docker Compose deployment
│   └── docker-compose.example.yml   # Docker Compose example configuration
├── kubernetes/                      # Raw Kubernetes manifests
│   ├── deployment.yaml              # Kubernetes deployment manifests
│   ├── secrets.example.yaml         # Secrets template (safe to commit)
│   ├── secrets.yaml                 # Actual secrets (gitignored)
│   └── sealed-secrets.yaml          # Encrypted secrets (safe to commit)
├── helm/                            # Helm chart for Kubernetes
│   └── codesourcerer/               # Helm chart directory
│       ├── Chart.yaml               # Chart metadata
│       ├── values.yaml              # Default configuration values
│       ├── values-secrets.example.yaml  # Secrets values template
│       └── templates/               # Kubernetes templates
│           ├── secrets.yaml         # Conditional secrets template
│           ├── *-deployment.yaml    # Service deployments
│           └── *-service.yaml       # Service definitions
├── fluxcd/                          # FluxCD GitOps configuration
│   └── (to be created)              # Base and overlay configurations
└── scripts/                         # Deployment scripts (gitignored)
    ├── build-images.ps1             # Build Docker images
    ├── deploy-helm.ps1              # Advanced Helm deployment
    ├── deploy-helm-simple.ps1       # Simple Helm deployment
    ├── cleanup-all.ps1              # Complete cleanup script
    └── cleanup-quick.ps1            # Quick cleanup script
```

## 🚀 Quick Start

### Prerequisites

| Tool | Purpose | Installation |
|------|---------|--------------|
| **Docker & Docker Compose** | Container runtime | [Install Docker](https://docs.docker.com/get-docker/) |
| **Kubernetes cluster** | Container orchestration | k3d, minikube, or cloud provider |
| **Helm 3.x** | Kubernetes package manager | [Install Helm](https://helm.sh/docs/intro/install/) |
| **kubectl** | Kubernetes CLI | [Install kubectl](https://kubernetes.io/docs/tasks/tools/) |
| **kubeseal** (optional) | Sealed Secrets CLI | See [Sealed Secrets Setup](#sealed-secrets-setup) |
| **flux** (optional) | FluxCD CLI | [Install Flux](https://fluxcd.io/docs/installation/) |
| **PowerShell** | Script execution | Pre-installed on Windows |

### 30-Second Local Test

```powershell
# 1. Build images
./deployment/scripts/build-images.ps1

# 2. Start local cluster
k3d cluster create codesourcerer-test --api-port 0.0.0.0:6550

# 3. Fix kubeconfig (Windows Docker Desktop issue)
kubectl config set-cluster k3d-codesourcerer-test --server=https://127.0.0.1:6550

# 4. Deploy with Helm
./deployment/scripts/deploy-helm-simple.ps1 install -CreateNamespace -Namespace codesourcerer

# 5. Test
kubectl get pods -n codesourcerer
```

---

## 🚢 Deployment Methods

### Option 1: Docker Compose (Development)

**Best for:** Local development, quick testing, debugging

#### Steps

1. **Build images:**
   ```powershell
   ./deployment/scripts/build-images.ps1
   ```

2. **Create secrets file:**
   ```powershell
   # Copy example and fill in your values
   cp deployment/docker-compose/docker-compose.example.yml deployment/docker-compose/docker-compose.yml
   # Edit docker-compose.yml with your actual secrets
   ```

3. **Deploy services:**
   ```powershell
   cd deployment/docker-compose
   docker-compose up -d
   ```

4. **Test the deployment:**
   ```powershell
   curl http://localhost:3001/testsend
   ```

5. **View logs:**
   ```powershell
   docker-compose logs -f [service-name]
   ```

6. **Stop services:**
   ```powershell
   docker-compose down
   ```

### Option 2: Helm Chart (Production)

**Best for:** Production Kubernetes deployments, staging environments

#### Steps

1. **Build and push images (if using custom registry):**
   ```powershell
   ./deployment/scripts/build-images.ps1 -Registry your-registry.com -Push
   ```

2. **Create secrets (choose one method):**

   **Method A: Using Sealed Secrets (Recommended)**
   ```powershell
   # See "Sealed Secrets Setup" section below
   ```

   **Method B: Using values file**
   ```powershell
   # Copy example and fill in your values
   cp deployment/helm/codesourcerer/values-secrets.example.yaml values-secrets.yaml
   # Edit values-secrets.yaml with your actual secrets (DO NOT COMMIT)
   ```

3. **Install with Helm:**
   ```powershell
   # With sealed secrets
   ./deployment/scripts/deploy-helm.ps1 install -CreateNamespace -Namespace codesourcerer

   # With values file
   ./deployment/scripts/deploy-helm.ps1 install -CreateNamespace -Namespace codesourcerer -SecretsFile values-secrets.yaml
   ```

4. **Check status:**
   ```powershell
   ./deployment/scripts/deploy-helm.ps1 status -Namespace codesourcerer
   kubectl get pods -n codesourcerer
   ```

5. **Upgrade deployment:**
   ```powershell
   ./deployment/scripts/deploy-helm.ps1 upgrade -Namespace codesourcerer
   ```

### Option 3: FluxCD GitOps (Production)

**Best for:** Production environments, automated deployments, GitOps workflows

#### Prerequisites

- GitHub repository with this codebase
- Kubernetes cluster with FluxCD installed
- Sealed Secrets controller installed
- GitHub Personal Access Token

#### Steps

1. **Install Flux CLI:**
   ```bash
   # Linux/macOS
   curl -s https://fluxcd.io/install.sh | sudo bash

   # Windows (using Chocolatey)
   choco install flux

   # Verify
   flux --version
   ```

2. **Bootstrap FluxCD:**
   ```bash
   # Export GitHub token
   export GITHUB_TOKEN=<your-github-token>

   # Bootstrap FluxCD to your cluster
   flux bootstrap github \
     --owner=<github-username> \
     --repository=<repository-name> \
     --branch=main \
     --path=./deployment/fluxcd/overlays/production \
     --personal
   ```

3. **Install Sealed Secrets Controller:**
   ```bash
   kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.27.1/controller.yaml
   ```

4. **Create and commit sealed secrets:**
   ```bash
   # Create sealed secret
   kubeseal -f deployment/kubernetes/secrets.yaml \
     -w deployment/kubernetes/sealed-secrets.yaml \
     --namespace codesourcerer

   # Commit to Git
   git add deployment/kubernetes/sealed-secrets.yaml
   git commit -m "Add sealed secrets"
   git push
   ```

5. **Verify deployment:**
   ```bash
   # Check Flux status
   flux get all

   # Check application pods
   kubectl get pods -n codesourcerer
   ```

---

## 🔐 Secrets Management

### Overview

CodeSourcerer uses **Sealed Secrets** for production-grade secret management that is GitOps-ready. Sealed Secrets allows you to encrypt your Kubernetes secrets so they can be safely stored in Git repositories.

**Key Benefits:**
- ✅ Safe to commit encrypted secrets to Git
- ✅ GitOps-ready for FluxCD integration
- ✅ Automatic decryption in cluster
- ✅ No manual secret management needed
- ✅ Asymmetric encryption (only cluster can decrypt)

### Sealed Secrets Setup

#### 1. Install kubeseal CLI

**Windows (PowerShell):**
```powershell
$version = "0.27.1"
$url = "https://github.com/bitnami-labs/sealed-secrets/releases/download/v$version/kubeseal-$version-windows-amd64.tar.gz"
$outFile = "$env:TEMP\kubeseal.tar.gz"
Invoke-WebRequest -Uri $url -OutFile $outFile
tar -xzf $outFile -C $env:TEMP
Move-Item -Path "$env:TEMP\kubeseal.exe" -Destination "$env:USERPROFILE\kubeseal.exe" -Force

# Add to PATH (add this to your PowerShell profile for persistence)
$env:PATH += ";$env:USERPROFILE"

# Verify installation
kubeseal --version
```

**Linux/macOS:**
```bash
wget https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.27.1/kubeseal-0.27.1-linux-amd64.tar.gz
tar -xzf kubeseal-0.27.1-linux-amd64.tar.gz
sudo install -m 755 kubeseal /usr/local/bin/kubeseal

# Verify installation
kubeseal --version
```

#### 2. Install Sealed Secrets Controller

```bash
# Install the controller in your cluster
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.27.1/controller.yaml

# Wait for controller to be ready
kubectl wait --for=condition=ready pod -l name=sealed-secrets-controller -n kube-system --timeout=120s

# Verify installation
kubectl get pods -n kube-system | grep sealed-secrets
```

Expected output:
```
sealed-secrets-controller-xxxxx   1/1     Running   0          1m
```

### Creating Sealed Secrets

#### Method 1: From existing secrets.yaml

```bash
# Navigate to the kubernetes directory
cd deployment/kubernetes

# Create a sealed secret from your secrets.yaml
# IMPORTANT: Make sure your cluster is running first!
kubeseal -f secrets.yaml -w sealed-secrets.yaml

# The sealed-secrets.yaml can now be safely committed to Git
git add sealed-secrets.yaml
git commit -m "Add sealed secrets"
```

#### Method 2: From scratch

```bash
# Create a regular secret (don't commit this!)
kubectl create secret generic app-secrets \
  --from-literal=GEMINI_API_KEY=your-api-key \
  --from-literal=PAT_TOKEN=your-pat-token \
  --from-literal=APP_ID=your-app-id \
  --from-literal=INSTALLATION_ID=your-installation-id \
  --from-literal=BOT_EMAIL=your-bot-email \
  --from-literal=PRIVATE_KEY_PATH=your-private-key-path \
  --dry-run=client -o yaml > temp-secret.yaml

# Seal it
kubeseal -f temp-secret.yaml -w sealed-secrets.yaml

# Delete the temp file
rm temp-secret.yaml
```

#### Method 3: Using PowerShell

```powershell
# Set PATH to include kubeseal
$env:PATH += ";$env:USERPROFILE"

# Navigate to kubernetes directory
cd deployment\kubernetes

# Create sealed secret
kubeseal -f secrets.yaml -w sealed-secrets.yaml

# Verify the sealed secret was created
cat sealed-secrets.yaml
```

### Using Sealed Secrets

#### Step 1: Apply Sealed Secrets to Cluster

```bash
# Apply the sealed secret
kubectl apply -f deployment/kubernetes/sealed-secrets.yaml

# The controller will automatically decrypt it and create the actual secret
# Wait a few seconds, then verify the secret was created
kubectl get secrets app-secrets

# Check the secret details (values will be base64 encoded)
kubectl describe secret app-secrets
```

#### Step 2: Deploy with Helm

**Option A: Using External Secrets (Recommended)**

```powershell
# Make sure values.yaml has:
# secrets:
#   useExternalSecret: true
#   externalSecretName: "app-secrets"

# Deploy
./deployment/scripts/deploy-helm.ps1 install \
  -CreateNamespace \
  -Namespace codesourcerer
```

**Option B: Using Secrets Values File**

```powershell
# Create a custom values file (don't commit this!)
# values-local.yaml
@"
secrets:
  useExternalSecret: false
genai:
  env:
    GEMINI_API_KEY: "your-actual-api-key"
github:
  env:
    PAT_TOKEN: "your-pat-token"
    APP_ID: "your-app-id"
    INSTALLATION_ID: "your-installation-id"
    BOT_EMAIL: "your-bot-email"
    PRIVATE_KEY_PATH: "your-private-key-path"
"@ | Out-File values-local.yaml

# Deploy
./deployment/scripts/deploy-helm.ps1 install \
  -CreateNamespace \
  -Namespace codesourcerer \
  -SecretsFile values-local.yaml
```

### Secret Fields Reference

The following secrets are required:

| Secret Key | Description | Example |
|------------|-------------|---------|
| `GEMINI_API_KEY` | Google Gemini API key for AI functionality | `AIzaSy...` |
| `PAT_TOKEN` | GitHub Personal Access Token | `ghp_...` |
| `APP_ID` | GitHub App ID | `1028002` |
| `INSTALLATION_ID` | GitHub App Installation ID | `56113322` |
| `BOT_EMAIL` | Bot email address | `bot@example.com` |
| `PRIVATE_KEY_PATH` | Path to GitHub App private key | `keys/app.pem` |

### What's Safe to Commit?

✅ **Safe to commit:**
- `deployment/kubernetes/secrets.example.yaml` - Template with placeholders
- `deployment/kubernetes/sealed-secrets.yaml` - Encrypted secrets
- `deployment/helm/codesourcerer/values-secrets.example.yaml` - Template
- All documentation files

❌ **Never commit:**
- `deployment/kubernetes/secrets.yaml` - Actual secrets (gitignored)
- `deployment/helm/codesourcerer/values-secrets.yaml` - Actual values (gitignored)
- Any file with actual API keys, tokens, or passwords

---

## 🔧 Configuration

### Environment Variables

The following environment variables need to be configured:

#### Required for Gen-AI Service:
- `GEMINI_API_KEY`: Google Gemini API key for AI functionality

#### Required for GitHub Service:
- `PAT_TOKEN`: GitHub Personal Access Token
- `APP_ID`: GitHub App ID
- `INSTALLATION_ID`: GitHub App Installation ID
- `BOT_EMAIL`: Bot email address
- `PRIVATE_KEY_PATH`: Path to GitHub App private key

### Docker Compose Configuration

Edit `docker-compose/docker-compose.yml` to customize:
- Port mappings
- Environment variables
- Resource limits
- Volume mounts

### Helm Configuration

Create a custom values file or modify `helm/codesourcerer/values.yaml`:

```yaml
# custom-values.yaml
genai:
  enabled: true
  replicas: 1
  image:
    repository: codesourcerer/gen-ai
    tag: latest
  env:
    GEMINI_API_KEY: ""  # Set via sealed secrets or values file
  resources:
    requests:
      memory: "256Mi"
      cpu: "100m"
    limits:
      memory: "512Mi"
      cpu: "500m"

github:
  enabled: true
  replicas: 1
  image:
    repository: codesourcerer/github
    tag: latest
  env:
    PAT_TOKEN: ""  # Set via sealed secrets or values file
    APP_ID: ""
    INSTALLATION_ID: ""
    BOT_EMAIL: ""
    PRIVATE_KEY_PATH: ""
  service:
    type: LoadBalancer  # or NodePort for local clusters
    port: 3000

database:
  enabled: true
  replicas: 1
  image:
    repository: codesourcerer/database
    tag: latest

redis:
  enabled: true
  architecture: standalone
  auth:
    enabled: false

ingress:
  enabled: true
  className: "nginx"
  hosts:
    - host: codesourcerer.yourdomain.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: codesourcerer-tls
      hosts:
        - codesourcerer.yourdomain.com
```

Then deploy with:
```powershell
./deployment/scripts/deploy-helm-simple.ps1 install -ValuesFile custom-values.yaml
```

---

## 🧹 Cleanup

### Quick Cleanup (Docker Compose only)

```powershell
# Using cleanup script
./deployment/scripts/cleanup-quick.ps1

# Manual cleanup
cd deployment/docker-compose
docker-compose down -v
```

### Complete Cleanup (All resources)

```powershell
# Clean everything with confirmation
./deployment/scripts/cleanup-all.ps1 -All

# Clean everything without confirmation
./deployment/scripts/cleanup-all.ps1 -All -Force

# Clean specific resources
./deployment/scripts/cleanup-all.ps1 -Docker -Images
./deployment/scripts/cleanup-all.ps1 -Kubernetes
```

### Manual Cleanup Commands

#### Docker Compose
```bash
# Stop and remove containers, networks, volumes
docker-compose -f deployment/docker-compose/docker-compose.yml down -v

# Remove images
docker rmi codesourcerer/database:latest
docker rmi codesourcerer/gen-ai:latest
docker rmi codesourcerer/github:latest

# Remove all CodeSourcerer images
docker images | grep codesourcerer | awk '{print $3}' | xargs docker rmi -f
```

#### Kubernetes/Helm
```bash
# Uninstall Helm release
helm uninstall codesourcerer -n codesourcerer

# Delete namespace (this deletes all resources in the namespace)
kubectl delete namespace codesourcerer

# Delete sealed secrets controller (if needed)
kubectl delete -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.27.1/controller.yaml
```

#### k3d Cluster
```bash
# Delete k3d cluster
k3d cluster delete codesourcerer-test

# List all k3d clusters
k3d cluster list

# Delete all k3d clusters
k3d cluster delete --all
```

#### FluxCD
```bash
# Uninstall Flux (keeps CRDs)
flux uninstall

# Uninstall Flux and remove CRDs
flux uninstall --crds

# Remove Flux from Git repository
rm -rf deployment/fluxcd/flux-system
git add deployment/fluxcd/flux-system
git commit -m "Remove Flux system"
git push
```

### Complete System Reset

To completely reset everything:

```bash
# 1. Stop all deployments
helm uninstall codesourcerer -n codesourcerer 2>/dev/null || true
docker-compose -f deployment/docker-compose/docker-compose.yml down -v 2>/dev/null || true

# 2. Delete Kubernetes resources
kubectl delete namespace codesourcerer 2>/dev/null || true

# 3. Delete k3d cluster
k3d cluster delete codesourcerer-test 2>/dev/null || true

# 4. Remove Docker images
docker images | grep codesourcerer | awk '{print $3}' | xargs docker rmi -f 2>/dev/null || true

# 5. Clean Docker system
docker system prune -af --volumes

# 6. Verify cleanup
docker ps -a
kubectl get all -A
k3d cluster list
```

---

## 🏗️ Architecture

### System Components

The CodeSourcerer platform consists of four main components:

| Component | Technology | Port | Protocol | Purpose |
|-----------|-----------|------|----------|---------|
| **Redis** | Redis 7 Alpine | 6379 | TCP | Data persistence and caching |
| **Database Service** | Go + gRPC | 8080/8082 | gRPC | Data operations and storage interface |
| **Gen-AI Service** | Go + gRPC | 8081/8083 | gRPC | AI-powered test generation using Gemini |
| **GitHub Service** | Go + Gin | 3000/3001 | HTTP/gRPC | GitHub webhook handler and API |

### Service Communication Flow

```
┌─────────────────┐
│  GitHub Webhook │
└────────┬────────┘
         │ HTTP
         ▼
┌─────────────────────┐
│  GitHub Service     │
│  (Port 3000)        │
└──────┬──────┬───────┘
       │      │
       │ gRPC │ gRPC
       │      │
       ▼      ▼
┌──────────┐ ┌──────────────┐
│ Database │ │  Gen-AI      │
│ Service  │ │  Service     │
│ (8080)   │ │  (8081)      │
└────┬─────┘ └──────┬───────┘
     │              │
     │ TCP          │ HTTPS
     ▼              ▼
┌─────────┐   ┌──────────────┐
│  Redis  │   │ Google Gemini│
│  (6379) │   │     API      │
└─────────┘   └──────────────┘
```

### Communication Patterns

- **GitHub Service → Database Service**: gRPC calls for data persistence
- **GitHub Service → Gen-AI Service**: gRPC calls for test generation
- **Database Service → Redis**: Direct TCP connection for data storage
- **Gen-AI Service → Google Gemini API**: HTTPS REST API calls

### Deployment Topologies

#### Development (Docker Compose)
- All services on single host
- Direct container-to-container communication
- Port mapping to localhost
- Shared Docker network

#### Production (Kubernetes)
- Services distributed across nodes
- Kubernetes service discovery
- Internal ClusterIP services
- External LoadBalancer/Ingress for GitHub Service
- Horizontal pod autoscaling (optional)
- Multi-zone deployment (optional)

---

## 🔄 FluxCD GitOps Integration

### Overview

FluxCD enables GitOps-based continuous deployment where your Git repository is the single source of truth for your infrastructure and applications.

### Prerequisites

- Kubernetes cluster (k3d, minikube, or production)
- GitHub repository with this codebase
- GitHub Personal Access Token with repo permissions
- Flux CLI installed
- Sealed Secrets controller installed

### FluxCD Setup

#### 1. Install Flux CLI

```bash
# Linux/macOS
curl -s https://fluxcd.io/install.sh | sudo bash

# Windows (using Chocolatey)
choco install flux

# Verify
flux --version
```

#### 2. Check Prerequisites

```bash
# Check if cluster is ready for Flux
flux check --pre
```

#### 3. Bootstrap FluxCD

```bash
# Export GitHub token
export GITHUB_TOKEN=<your-github-token>

# Bootstrap FluxCD to your cluster
flux bootstrap github \
  --owner=<github-username> \
  --repository=<repository-name> \
  --branch=main \
  --path=./deployment/fluxcd/overlays/production \
  --personal
```

This will:
- Install Flux components in your cluster
- Create a deploy key in your GitHub repository
- Create Flux configuration files
- Set up automatic synchronization

#### 4. Verify Installation

```bash
# Check Flux components
flux check

# Check all Flux resources
flux get all

# Watch reconciliation
flux logs --all-namespaces --follow
```

### FluxCD Directory Structure

```
deployment/fluxcd/
├── base/                              # Base configurations
│   ├── namespace.yaml                 # Namespace definition
│   ├── helmrepository.yaml            # Helm repository source
│   └── kustomization.yaml             # Base kustomization
└── overlays/                          # Environment-specific overlays
    ├── development/
    │   ├── sealed-secrets.yaml        # Dev sealed secrets
    │   ├── helmrelease.yaml           # Dev Helm release config
    │   └── kustomization.yaml         # Dev kustomization
    ├── staging/
    │   ├── sealed-secrets.yaml        # Staging sealed secrets
    │   ├── helmrelease.yaml           # Staging Helm release config
    │   └── kustomization.yaml         # Staging kustomization
    └── production/
        ├── sealed-secrets.yaml        # Production sealed secrets
        ├── helmrelease.yaml           # Production Helm release config
        └── kustomization.yaml         # Production kustomization
```

### Creating FluxCD Configurations

#### Example HelmRelease

```yaml
# deployment/fluxcd/overlays/production/helmrelease.yaml
apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: codesourcerer
  namespace: codesourcerer
spec:
  interval: 5m
  chart:
    spec:
      chart: ./deployment/helm/codesourcerer
      sourceRef:
        kind: GitRepository
        name: codesourcerer
        namespace: flux-system
      interval: 1m
  values:
    secrets:
      useExternalSecret: true
      externalSecretName: "app-secrets"
    genai:
      enabled: true
      replicas: 2
    github:
      enabled: true
      replicas: 2
    database:
      enabled: true
      replicas: 1
    redis:
      enabled: true
```

### GitOps Workflow

1. **Make changes** to your deployment configuration in Git
2. **Commit and push** to your repository
3. **Flux automatically detects** changes (default: every 1 minute)
4. **Flux reconciles** the cluster state to match Git
5. **Monitor** the deployment with `flux get all`

### Useful Flux Commands

```bash
# Trigger immediate reconciliation
flux reconcile source git codesourcerer
flux reconcile helmrelease codesourcerer -n codesourcerer

# Suspend/resume reconciliation
flux suspend helmrelease codesourcerer -n codesourcerer
flux resume helmrelease codesourcerer -n codesourcerer

# Check status
flux get sources git
flux get helmreleases -A
flux get kustomizations -A

# View logs
flux logs --level=info
flux logs --kind=HelmRelease --name=codesourcerer -n codesourcerer

# Export current configuration
flux export source git codesourcerer
flux export helmrelease codesourcerer -n codesourcerer
```

---

## 🧪 Testing & Monitoring

### Health Checks

| Service | Health Check | Command |
|---------|--------------|---------|
| **GitHub Service** | HTTP endpoint | `curl http://localhost:3001/health` |
| **Database Service** | TCP connection | `nc -zv localhost 8082` |
| **Gen-AI Service** | TCP connection | `nc -zv localhost 8083` |
| **Redis** | Redis ping | `redis-cli ping` |

### Test Endpoints

```bash
# Test Gen-AI Integration
curl http://localhost:3001/testsend

# Test GitHub Integration
curl http://localhost:3001/testfinalizer

# Test with port-forward (Kubernetes)
kubectl port-forward svc/codesourcerer-github 3001:3000 -n codesourcerer
curl http://localhost:3001/testsend
```

### Monitoring

#### View Logs

```powershell
# Docker Compose
docker-compose -f deployment/docker-compose/docker-compose.yml logs -f [service-name]

# Kubernetes - All pods
kubectl logs -f -l app.kubernetes.io/name=codesourcerer -n codesourcerer

# Kubernetes - Specific service
kubectl logs -f deployment/codesourcerer-github -n codesourcerer
kubectl logs -f deployment/codesourcerer-database -n codesourcerer
kubectl logs -f deployment/codesourcerer-genai -n codesourcerer

# Kubernetes - Follow logs from all containers
kubectl logs -f -l app.kubernetes.io/instance=codesourcerer -n codesourcerer --all-containers
```

#### Resource Usage

```bash
# Check pod resource usage
kubectl top pods -n codesourcerer

# Check node resource usage
kubectl top nodes

# Describe pod for detailed info
kubectl describe pod <pod-name> -n codesourcerer
```

#### Service Status

```bash
# Check all pods
kubectl get pods -n codesourcerer

# Check services
kubectl get svc -n codesourcerer

# Check endpoints
kubectl get endpoints -n codesourcerer

# Check Helm release status
helm status codesourcerer -n codesourcerer
```

---

## 🔒 Security

### Security Best Practices Implemented

#### Container Security
- ✅ **Non-root user execution** - All containers run as non-root user (UID 1001)
- ✅ **Minimal Alpine Linux base images** - Reduced attack surface
- ✅ **No privileged containers** - Security contexts enforced
- ✅ **Resource limits enforced** - Prevents resource exhaustion attacks
- ✅ **Read-only root filesystem** (where applicable)
- ✅ **Security scanning** - Scan images before deployment

#### Secrets Management
- ✅ **Sealed Secrets** for GitOps-ready secret management (recommended for production)
  - Asymmetric encryption - only cluster can decrypt
  - Encrypted secrets safe to commit to Git
  - Automatic decryption by cluster controller
  - Namespace-scoped or cluster-wide encryption
- ✅ **Kubernetes Secrets** for sensitive data
- ✅ **Environment variable injection** - No hardcoded credentials
- ✅ **No hardcoded credentials** in images or code
- ✅ **Example files provided** (`.example.yaml`) for documentation
- ✅ **Comprehensive .gitignore** - Prevents accidental secret commits

#### Network Security
- ✅ **Internal service communication only** - Services not exposed externally by default
- ✅ **Ingress/LoadBalancer** for controlled external access
- ✅ **TLS/SSL support** for ingress
- ✅ **Network policies** (can be added for additional isolation)
- ✅ **Service mesh ready** (Istio/Linkerd compatible)

#### Access Control
- ✅ **RBAC enabled** - Role-based access control
- ✅ **Service accounts** - Dedicated service accounts per component
- ✅ **Pod security policies** - Enforced security standards
- ✅ **Namespace isolation** - Logical separation of environments

### Security Audit Checklist

Before deploying to production, verify:

- [ ] All secrets are managed via Sealed Secrets
- [ ] No hardcoded credentials in code or configuration
- [ ] All images scanned for vulnerabilities
- [ ] Resource limits set for all containers
- [ ] Network policies configured
- [ ] RBAC roles properly configured
- [ ] TLS enabled for all external endpoints
- [ ] Logging and monitoring configured
- [ ] Backup and disaster recovery plan in place

### Files Safe to Commit

✅ **Always safe:**
- `deployment/kubernetes/secrets.example.yaml`
- `deployment/kubernetes/sealed-secrets.yaml`
- `deployment/helm/codesourcerer/values.yaml`
- `deployment/helm/codesourcerer/values-secrets.example.yaml`
- All `*.example.yml` files
- All documentation files

❌ **Never commit:**
- `deployment/kubernetes/secrets.yaml`
- `deployment/helm/codesourcerer/values-secrets.yaml`
- `deployment/docker-compose/docker-compose.yml`
- Any file with actual API keys, tokens, or passwords
- Private keys (`.pem`, `.key` files)
- Any file in `deployment/scripts/` (gitignored)

---

## 🚨 Troubleshooting

### Common Issues

#### 1. Services Can't Communicate

**Symptoms:**
- Services timeout when calling each other
- gRPC connection errors
- DNS resolution failures

**Solutions:**
```bash
# Check service discovery (DNS resolution)
kubectl run -it --rm debug --image=busybox --restart=Never -- nslookup codesourcerer-database.codesourcerer.svc.cluster.local

# Verify port configurations
kubectl get svc -n codesourcerer

# Check network policies
kubectl get networkpolicies -n codesourcerer

# Test connectivity between pods
kubectl exec -it <github-pod> -n codesourcerer -- nc -zv codesourcerer-database 8080
```

#### 2. Images Not Found

**Symptoms:**
- `ImagePullBackOff` or `ErrImagePull` errors
- Pods stuck in `Pending` state

**Solutions:**
```bash
# Build images first
./deployment/scripts/build-images.ps1

# For k3d, import images to cluster
k3d image import codesourcerer/database:latest codesourcerer/genai:latest codesourcerer/github:latest -c codesourcerer-test

# Check image tags in values.yaml
cat deployment/helm/codesourcerer/values.yaml | grep -A 3 "image:"

# Verify registry access (if using private registry)
kubectl create secret docker-registry regcred \
  --docker-server=<your-registry> \
  --docker-username=<username> \
  --docker-password=<password> \
  -n codesourcerer
```

#### 3. Sealed Secrets Not Working

**Symptoms:**
- Secrets not created after applying SealedSecret
- Controller errors in logs

**Solutions:**
```bash
# Check sealed-secrets controller is running
kubectl get pods -n kube-system | grep sealed-secrets

# Check controller logs
kubectl logs -n kube-system -l name=sealed-secrets-controller

# Verify SealedSecret was created
kubectl get sealedsecrets -n codesourcerer

# Check if secret was unsealed
kubectl get secrets -n codesourcerer

# Re-create sealed secret with correct namespace
kubeseal -f deployment/kubernetes/secrets.yaml \
  -w deployment/kubernetes/sealed-secrets.yaml \
  --namespace codesourcerer
```

#### 4. API Key Errors

**Symptoms:**
- "Invalid API key" errors
- Authentication failures
- 401/403 HTTP errors

**Solutions:**
```bash
# Check secret exists
kubectl get secret app-secrets -n codesourcerer

# Verify secret contents (base64 encoded)
kubectl get secret app-secrets -n codesourcerer -o yaml

# Decode and verify a specific key
kubectl get secret app-secrets -n codesourcerer -o jsonpath='{.data.GEMINI_API_KEY}' | base64 -d

# Re-create secret if needed
kubectl delete secret app-secrets -n codesourcerer
kubectl apply -f deployment/kubernetes/sealed-secrets.yaml
```

#### 5. Resource Constraints

**Symptoms:**
- Pods stuck in `Pending` state
- OOMKilled errors
- CPU throttling

**Solutions:**
```bash
# Check node resources
kubectl top nodes

# Check pod resource usage
kubectl top pods -n codesourcerer

# Describe pod to see resource issues
kubectl describe pod <pod-name> -n codesourcerer

# Adjust resource requests/limits in values.yaml
# Then upgrade the release
helm upgrade codesourcerer deployment/helm/codesourcerer -n codesourcerer

# Scale down if needed
kubectl scale deployment codesourcerer-genai --replicas=1 -n codesourcerer
```

#### 6. k3d Cluster Connection Issues (Windows)

**Symptoms:**
- `kubectl` commands timeout
- "connection refused" errors
- Can't connect to `host.docker.internal`

**Solutions:**
```bash
# Delete and recreate cluster with specific API port
k3d cluster delete codesourcerer-test
k3d cluster create codesourcerer-test --api-port 0.0.0.0:6550

# Update kubeconfig to use localhost instead of host.docker.internal
kubectl config set-cluster k3d-codesourcerer-test --server=https://127.0.0.1:6550

# Verify connection
kubectl get nodes
```

#### 7. FluxCD Not Reconciling

**Symptoms:**
- Changes in Git not reflected in cluster
- Flux stuck on old version

**Solutions:**
```bash
# Check Flux status
flux check

# Check all Flux resources
flux get all

# Trigger manual reconciliation
flux reconcile source git codesourcerer
flux reconcile helmrelease codesourcerer -n codesourcerer

# Check Flux logs
flux logs --all-namespaces --level=error

# Suspend and resume to force refresh
flux suspend helmrelease codesourcerer -n codesourcerer
flux resume helmrelease codesourcerer -n codesourcerer
```

### Debug Commands

```bash
# Check pod status
kubectl get pods -n codesourcerer -o wide

# Describe problematic pod
kubectl describe pod <pod-name> -n codesourcerer

# Check pod events
kubectl get events -n codesourcerer --sort-by='.lastTimestamp'

# Check service endpoints
kubectl get endpoints -n codesourcerer

# Port forward for testing
kubectl port-forward svc/codesourcerer-github 3001:3000 -n codesourcerer

# Execute commands in pod
kubectl exec -it <pod-name> -n codesourcerer -- /bin/sh

# Check logs from all containers in a pod
kubectl logs <pod-name> -n codesourcerer --all-containers

# Check previous container logs (if crashed)
kubectl logs <pod-name> -n codesourcerer --previous
```

### Getting Help

If you're still experiencing issues:

1. **Check logs** - Most issues are visible in pod logs
2. **Describe resources** - Use `kubectl describe` for detailed information
3. **Check events** - Kubernetes events often show what went wrong
4. **Verify secrets** - Ensure all required secrets are present and correct
5. **Resource limits** - Check if pods have enough resources
6. **Network connectivity** - Test service-to-service communication

---

## 📚 Additional Resources

### Official Documentation

- [Docker Compose Documentation](https://docs.docker.com/compose/) - Container orchestration for development
- [Kubernetes Documentation](https://kubernetes.io/docs/) - Container orchestration platform
- [Helm Documentation](https://helm.sh/docs/) - Kubernetes package manager
- [k3d Documentation](https://k3d.io/) - Lightweight Kubernetes in Docker
- [FluxCD Documentation](https://fluxcd.io/docs/) - GitOps toolkit for Kubernetes
- [Sealed Secrets Documentation](https://github.com/bitnami-labs/sealed-secrets) - Encrypted Kubernetes secrets

### Related Tools

- [kubectl](https://kubernetes.io/docs/reference/kubectl/) - Kubernetes command-line tool
- [kubeseal](https://github.com/bitnami-labs/sealed-secrets#installation) - Sealed Secrets CLI
- [flux](https://fluxcd.io/docs/cmd/) - FluxCD CLI
- [minikube](https://minikube.sigs.k8s.io/) - Local Kubernetes alternative
- [kind](https://kind.sigs.k8s.io/) - Kubernetes in Docker alternative

### Learning Resources

- [Kubernetes Basics](https://kubernetes.io/docs/tutorials/kubernetes-basics/)
- [Helm Charts Guide](https://helm.sh/docs/chart_template_guide/)
- [GitOps Principles](https://www.gitops.tech/)
- [Sealed Secrets Tutorial](https://github.com/bitnami-labs/sealed-secrets#usage)

---

## 🤝 Contributing

When adding new services or modifying deployment:

### Checklist

- [ ] Update Docker Compose configuration (`deployment/docker-compose/docker-compose.example.yml`)
- [ ] Update Helm chart templates (`deployment/helm/codesourcerer/templates/`)
- [ ] Update Helm values (`deployment/helm/codesourcerer/values.yaml`)
- [ ] Update deployment scripts if needed
- [ ] Add any new secrets to `secrets.example.yaml`
- [ ] Test Docker Compose deployment
- [ ] Test Helm deployment
- [ ] Test with Sealed Secrets
- [ ] Update this README with new configuration options
- [ ] Update FluxCD configurations if applicable

### Testing Your Changes

```bash
# 1. Test Docker Compose
cd deployment/docker-compose
docker-compose up -d
# Verify services are running
docker-compose ps
docker-compose down

# 2. Test Helm deployment
k3d cluster create test-cluster --api-port 0.0.0.0:6550
kubectl config set-cluster k3d-test-cluster --server=https://127.0.0.1:6550
helm install codesourcerer deployment/helm/codesourcerer -n test --create-namespace
kubectl get pods -n test
helm uninstall codesourcerer -n test
k3d cluster delete test-cluster

# 3. Test with Sealed Secrets
# (Follow Sealed Secrets setup section)
```

---

## 📋 Summary

### Quick Reference

| Task | Command |
|------|---------|
| **Build images** | `./deployment/scripts/build-images.ps1` |
| **Deploy with Docker Compose** | `cd deployment/docker-compose && docker-compose up -d` |
| **Deploy with Helm** | `./deployment/scripts/deploy-helm.ps1 install -CreateNamespace -Namespace codesourcerer` |
| **Create sealed secret** | `kubeseal -f deployment/kubernetes/secrets.yaml -w deployment/kubernetes/sealed-secrets.yaml` |
| **Check pod status** | `kubectl get pods -n codesourcerer` |
| **View logs** | `kubectl logs -f deployment/codesourcerer-github -n codesourcerer` |
| **Port forward** | `kubectl port-forward svc/codesourcerer-github 3001:3000 -n codesourcerer` |
| **Cleanup** | `./deployment/scripts/cleanup-all.ps1 -All` |

### Deployment Decision Matrix

| Scenario | Recommended Method | Why |
|----------|-------------------|-----|
| **Local development** | Docker Compose | Fast iteration, easy debugging |
| **Testing Kubernetes configs** | k3d + Helm | Lightweight, fast setup |
| **Staging environment** | Kubernetes + Helm | Production-like environment |
| **Production** | Kubernetes + Helm + FluxCD | GitOps, automated deployments |
| **Multi-environment** | FluxCD with overlays | Consistent deployments across environments |

### Security Checklist

Before deploying to production:

- [ ] All secrets managed via Sealed Secrets
- [ ] No hardcoded credentials in code or configuration
- [ ] All images scanned for vulnerabilities
- [ ] Resource limits set for all containers
- [ ] Network policies configured
- [ ] RBAC roles properly configured
- [ ] TLS enabled for all external endpoints
- [ ] Logging and monitoring configured
- [ ] Backup and disaster recovery plan in place
- [ ] Security audit completed

---

## 📞 Support

For issues, questions, or contributions:

1. **Check this documentation** - Most common scenarios are covered
2. **Review troubleshooting section** - Common issues and solutions
3. **Check logs** - Most issues are visible in application logs
4. **Open an issue** - For bugs or feature requests

---

**Last Updated:** 2025-11-10
**Version:** 1.0.0
**Maintained by:** CodeSourcerer Team
