#!/bin/bash

set -e

echo "🧹 Cleaning up School Schedule from Kubernetes..."

# Delete in reverse order
kubectl delete -f ingress/network-policies.yaml --ignore-not-found=true
kubectl delete -f ingress/ingress.yaml --ignore-not-found=true
kubectl delete -f services/rbac.yaml --ignore-not-found=true
kubectl delete -f services/iam.yaml --ignore-not-found=true
kubectl delete -f monitoring/otel-collector.yaml --ignore-not-found=true
kubectl delete -f infrastructure/kafka.yaml --ignore-not-found=true
kubectl delete -f infrastructure/redis.yaml --ignore-not-found=true
kubectl delete -f infrastructure/postgres.yaml --ignore-not-found=true
kubectl delete -f infrastructure/secrets.yaml --ignore-not-found=true
kubectl delete -f infrastructure/configmaps.yaml --ignore-not-found=true
kubectl delete -f namespaces/namespaces.yaml --ignore-not-found=true

echo "✅ Cleanup completed!"
