#!/bin/bash

set -e

echo "🚀 Deploying School Schedule to Kubernetes..."

# Create namespaces first
echo "📦 Creating namespaces..."
kubectl apply -f namespaces/namespaces.yaml

# Deploy infrastructure
echo "🏗️ Deploying infrastructure..."
kubectl apply -f infrastructure/configmaps.yaml
kubectl apply -f infrastructure/secrets.yaml
kubectl apply -f infrastructure/postgres.yaml
kubectl apply -f infrastructure/redis.yaml
kubectl apply -f infrastructure/kafka.yaml

# Wait for infrastructure to be ready
echo "⏳ Waiting for infrastructure to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres-iam -n school-schedule-infrastructure --timeout=300s
kubectl wait --for=condition=ready pod -l app=postgres-rbac -n school-schedule-infrastructure --timeout=300s
kubectl wait --for=condition=ready pod -l app=redis-cluster -n school-schedule-infrastructure --timeout=300s
kubectl wait --for=condition=ready pod -l app=kafka -n school-schedule-infrastructure --timeout=300s

# Deploy monitoring
echo "📊 Deploying monitoring..."
kubectl apply -f monitoring/otel-collector.yaml

# Deploy services
echo "🔧 Deploying services..."
kubectl apply -f services/iam.yaml
kubectl apply -f services/rbac.yaml

# Deploy ingress
echo "🌐 Deploying ingress..."
kubectl apply -f ingress/ingress.yaml
kubectl apply -f ingress/network-policies.yaml

echo "✅ Deployment completed!"
echo ""
echo "📋 Useful commands:"
echo "  kubectl get pods -A"
echo "  kubectl logs -f deployment/iam-service -n school-schedule-services"
echo "  kubectl logs -f deployment/rbac-service -n school-schedule-services"
echo "  kubectl port-forward service/iam-service 8080:8080 -n school-schedule-services"
echo "  kubectl port-forward service/rbac-service 8080:8080 -n school-schedule-services"
