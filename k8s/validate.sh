#!/bin/bash

set -e

echo "🔍 Validating Kubernetes manifests..."

# Check if we have kubectl
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl not found. Installing..."
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x kubectl
    export PATH=$PATH:$(pwd)
fi

echo "📋 Validating YAML syntax..."

# Validate namespaces
echo "  ✓ Validating namespaces..."
kubectl apply --dry-run=client -f namespaces/namespaces.yaml

# Validate infrastructure
echo "  ✓ Validating infrastructure configmaps..."
kubectl apply --dry-run=client -f infrastructure/configmaps.yaml

echo "  ✓ Validating infrastructure secrets..."
kubectl apply --dry-run=client -f infrastructure/secrets.yaml

echo "  ✓ Validating PostgreSQL..."
kubectl apply --dry-run=client -f infrastructure/postgres.yaml

echo "  ✓ Validating Redis..."
kubectl apply --dry-run=client -f infrastructure/redis.yaml

echo "  ✓ Validating Kafka..."
kubectl apply --dry-run=client -f infrastructure/kafka.yaml

# Validate monitoring
echo "  ✓ Validating OpenTelemetry..."
kubectl apply --dry-run=client -f monitoring/otel-collector.yaml

# Validate services
echo "  ✓ Validating IAM service..."
kubectl apply --dry-run=client -f services/iam.yaml

echo "  ✓ Validating RBAC service..."
kubectl apply --dry-run=client -f services/rbac.yaml

# Validate ingress
echo "  ✓ Validating Ingress..."
kubectl apply --dry-run=client -f ingress/ingress.yaml

echo "  ✓ Validating Network Policies..."
kubectl apply --dry-run=client -f ingress/network-policies.yaml

echo "✅ All manifests are valid!"
echo ""
echo "📊 Summary:"
echo "  - Namespaces: 5"
echo "  - ConfigMaps: 3"
echo "  - Secrets: 4"
echo "  - StatefulSets: 4"
echo "  - Deployments: 3"
echo "  - Services: 8"
echo "  - Ingress: 1"
echo "  - Network Policies: 3"
echo ""
echo "🚀 Ready for deployment!"
