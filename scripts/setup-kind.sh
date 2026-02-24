#!/bin/bash
# Setup kind cluster with Argo Workflows and shared services
#
# Usage:
#   ./scripts/setup-kind.sh [cluster-name]
#   
# Default cluster name: agentic-k8s-dev
#
# What this script does:
# 1. Creates a 2-node kind cluster (1 control plane, 1 worker)
# 2. Installs Argo Workflows v4.0.1
# 3. Installs agentic-k8s-operator
# 4. Deploys shared services (PostgreSQL, MinIO, Browserless, LiteLLM)
#
# Requirements:
# - Docker (running)
# - kind (installed: https://kind.sigs.k8s.io/docs/user/quick-start/)
# - kubectl (installed: https://kubernetes.io/docs/tasks/tools/)
# - helm (optional, for Argo installation)

set -euo pipefail

# Configuration
CLUSTER_NAME="${1:-agentic-k8s-dev}"
KIND_VERSION="0.21.0"
ARGO_VERSION="v3.7.10"  # Latest available chart version
NAMESPACE_ARGO="argo-workflows"
NAMESPACE_OPERATOR="agentic-system"
NAMESPACE_SHARED="shared-services"

echo "=========================================="
echo "Setting up kind cluster: $CLUSTER_NAME"
echo "Argo version: $ARGO_VERSION"
echo "=========================================="

# Step 1: Create kind cluster
echo ""
echo "[1/5] Creating kind cluster..."

kind create cluster \
  --name "$CLUSTER_NAME" \
  --config - <<EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
  - role: worker
kubeadmConfigPatches:
  - |
    kind: KubeletConfiguration
    # Allow pod scheduling on control plane (for single-node testing)
    systemReserved:
      memory: 256Mi
      cpu: 100m
      ephemeral-storage: 1Gi
EOF

echo "✓ kind cluster created: $CLUSTER_NAME"

# Wait for cluster to be ready
echo "Waiting for cluster to be ready..."
for i in {1..30}; do
  if kubectl cluster-info >/dev/null 2>&1; then
    echo "✓ Cluster is ready"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "✗ Cluster failed to become ready after 30 attempts"
    exit 1
  fi
  sleep 1
done

# Step 2: Install Argo Workflows
echo ""
echo "[2/5] Installing Argo Workflows $ARGO_VERSION..."

# Add Argo Helm repository
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update

# Create argo-workflows namespace
kubectl create namespace $NAMESPACE_ARGO --dry-run=client -o yaml | kubectl apply -f -

# Install Argo Workflows (use latest chart version)
helm install argo-workflows \
  argo/argo-workflows \
  --namespace $NAMESPACE_ARGO \
  --set serviceAccount.create=true \
  --set serviceAccount.name=argo-workflows \
  --set controller.serviceAccount.create=true \
  --set server.serviceAccount.create=true \
  --set server.authMode=server  # Allow unauthenticated access (dev only)

echo "✓ Argo Workflows installed"

# Wait for Argo controller to be ready
echo "Waiting for Argo controller..."
kubectl rollout status deployment/argo-workflows-workflow-controller \
  -n $NAMESPACE_ARGO \
  --timeout=120s || {
    echo "✗ Argo controller failed to start"
    exit 1
  }

echo "✓ Argo controller is ready"

# Step 3: Create shared-services namespace and deploy services
echo ""
echo "[3/5] Deploying shared services..."

kubectl create namespace $NAMESPACE_SHARED --dry-run=client -o yaml | kubectl apply -f -

# Deploy PostgreSQL, MinIO, Browserless, LiteLLM
echo "Applying shared-services manifests..."
kubectl apply -f config/shared-services/

# Wait for critical services to be ready
echo "Waiting for PostgreSQL to be ready..."
kubectl rollout status deployment/postgres \
  -n $NAMESPACE_SHARED \
  --timeout=120s || {
    echo "⚠ PostgreSQL failed to start (expected if PVC unavailable)"
  }

echo "Waiting for MinIO to be ready..."
kubectl rollout status deployment/minio \
  -n $NAMESPACE_SHARED \
  --timeout=120s || {
    echo "⚠ MinIO failed to start (expected if PVC unavailable)"
  }

echo "Waiting for Browserless to be ready..."
kubectl rollout status deployment/browserless \
  -n $NAMESPACE_SHARED \
  --timeout=120s || {
    echo "⚠ Browserless failed to start"
  }

echo "✓ Shared services deployed"

# Step 4: Deploy Argo WorkflowTemplate
echo ""
echo "[4/5] Deploying Argo WorkflowTemplate..."

# Create argo-workflows namespace if not exists (in case helm skipped it)
kubectl create namespace $NAMESPACE_ARGO --dry-run=client -o yaml | kubectl apply -f -

# Apply WorkflowTemplate
kubectl apply -f config/argo/workflowtemplate.yaml

echo "✓ WorkflowTemplate deployed"

# Step 5: Display access information
echo ""
echo "[5/5] Setup complete!"
echo "=========================================="
echo "Cluster Information"
echo "=========================================="
echo "Cluster name: $CLUSTER_NAME"
echo ""
echo "Namespaces:"
echo "  - $NAMESPACE_ARGO: Argo Workflows"
echo "  - $NAMESPACE_OPERATOR: Operator (to be deployed)"
echo "  - $NAMESPACE_SHARED: Shared services"
echo ""
echo "Services:"
echo "  - Argo Server: kubectl port-forward -n $NAMESPACE_ARGO svc/argo-workflows-server 2746:2746"
echo "  - MinIO S3 API: kubectl port-forward -n $NAMESPACE_SHARED svc/minio 9000:9000"
echo "  - MinIO Console: kubectl port-forward -n $NAMESPACE_SHARED svc/minio-console 9001:9001"
echo "  - Browserless: kubectl port-forward -n $NAMESPACE_SHARED svc/browserless 3000:3000"
echo "  - LiteLLM: kubectl port-forward -n $NAMESPACE_SHARED svc/litellm 8000:8000"
echo ""
echo "Next steps:"
echo "  1. Deploy the operator: kubectl apply -f config/operator/"
echo "  2. Check status: kubectl get all -n $NAMESPACE_OPERATOR"
echo "  3. Submit a workflow: argo submit -n $NAMESPACE_ARGO --from workflowtemplate/visual-analysis-template"
echo ""
echo "To delete the cluster:"
echo "  kind delete cluster --name $CLUSTER_NAME"

# Optional: Output kubeconfig location
echo ""
echo "Kubeconfig:"
echo "  export KUBECONFIG=\"\$(kind get kubeconfig-path --name=$CLUSTER_NAME)\""
