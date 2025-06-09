#!/bin/bash
# frontend/deployments/scripts/deploy.sh

#!/bin/bash
set -e

# Build and push Docker images
echo "Building and pushing Docker images..."
docker-compose -f deployments/docker/docker-compose.yml build
docker push autosysadmin/backend:latest
docker push autosysadmin/frontend:latest

# Apply Kubernetes configurations
echo "Deploying to Kubernetes..."
kubectl apply -f deployments/kubernetes/api-deployment.yaml
kubectl apply -f deployments/kubernetes/frontend-deployment.yaml
kubectl apply -f deployments/kubernetes/ingress.yaml

echo "Deployment complete!"