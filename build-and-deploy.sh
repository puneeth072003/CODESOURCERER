#!/bin/bash

# Build and Deploy Script for CodeSourcerer Services
set -e

echo "🚀 Building and deploying CodeSourcerer services to k3d..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if k3d is installed
if ! command -v k3d &> /dev/null; then
    print_error "k3d is not installed. Please install k3d first."
    exit 1
fi

# Check if kubectl is installed
if ! command -v kubectl &> /dev/null; then
    print_error "kubectl is not installed. Please install kubectl first."
    exit 1
fi

# Check if docker is running
if ! docker info &> /dev/null; then
    print_error "Docker is not running. Please start Docker first."
    exit 1
fi

# Create k3d cluster if it doesn't exist
CLUSTER_NAME="codesourcerer"
if ! k3d cluster list | grep -q "$CLUSTER_NAME"; then
    print_status "Creating k3d cluster: $CLUSTER_NAME"
    k3d cluster create $CLUSTER_NAME --port "3000:30000@loadbalancer" --port "8080:30080@loadbalancer" --port "8081:30081@loadbalancer"
    print_success "k3d cluster created successfully"
else
    print_status "k3d cluster '$CLUSTER_NAME' already exists"
    k3d cluster start $CLUSTER_NAME 2>/dev/null || true
fi

# Set kubectl context
kubectl config use-context k3d-$CLUSTER_NAME

# Build Docker images
print_status "Building Docker images..."

# Build database service
print_status "Building database service image..."
cd services/database
docker build -t codesourcerer/database:latest .
cd ../..
print_success "Database service image built"

# Build gen-ai service
print_status "Building gen-ai service image..."
cd services/gen-ai
docker build -t codesourcerer/gen-ai:latest .
cd ../..
print_success "Gen-AI service image built"

# Build github service
print_status "Building github service image..."
cd services/github
docker build -t codesourcerer/github:latest .
cd ../..
print_success "GitHub service image built"

# Import images to k3d cluster
print_status "Importing images to k3d cluster..."
k3d image import codesourcerer/database:latest -c $CLUSTER_NAME
k3d image import codesourcerer/gen-ai:latest -c $CLUSTER_NAME
k3d image import codesourcerer/github:latest -c $CLUSTER_NAME
print_success "Images imported to k3d cluster"

# Apply Kubernetes manifests
print_status "Applying Kubernetes manifests..."
kubectl apply -f deployment.yaml
print_success "Kubernetes manifests applied"

# Wait for deployments to be ready
print_status "Waiting for deployments to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/redis
kubectl wait --for=condition=available --timeout=300s deployment/database-service
kubectl wait --for=condition=available --timeout=300s deployment/genai-service
kubectl wait --for=condition=available --timeout=300s deployment/github-service

print_success "All deployments are ready!"

# Show status
print_status "Deployment status:"
kubectl get pods
kubectl get services

print_success "🎉 Deployment completed successfully!"
print_status "You can access the GitHub service at: http://localhost:3000"
print_status "To check logs, use: kubectl logs -f deployment/github-service"
print_status "To delete the cluster, use: k3d cluster delete $CLUSTER_NAME"
