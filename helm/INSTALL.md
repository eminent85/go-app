# Helm Chart Installation Guide

This guide provides step-by-step instructions for deploying the Go application using Helm.

## Quick Start

### 1. Basic Installation

```bash
# Install with default values
helm install go-app ./helm/go-app

# Install in a specific namespace
helm install go-app ./helm/go-app --namespace my-app --create-namespace
```

### 2. Verify Installation

```bash
# Check deployment status
kubectl get pods -l app.kubernetes.io/name=go-app

# Check all resources
kubectl get all -l app.kubernetes.io/name=go-app

# View release information
helm status go-app
```

## Deployment Scenarios

### Scenario 1: Standard Kubernetes with NGINX Ingress

```bash
helm install go-app ./helm/go-app \
  --set image.repository=myregistry/go-app \
  --set image.tag=1.0.0 \
  --set ingress.enabled=true \
  --set ingress.className=nginx \
  --set ingress.hosts[0].host=myapp.example.com \
  --set ingress.hosts[0].paths[0].path=/ \
  --set ingress.hosts[0].paths[0].pathType=Prefix
```

### Scenario 2: With Istio Service Mesh

```bash
helm install go-app ./helm/go-app \
  --values ./helm/go-app/values-istio.yaml \
  --set image.repository=myregistry/go-app \
  --set image.tag=1.0.0 \
  --set istio.virtualService.hosts[0]=myapp.example.com
```

### Scenario 3: Production Deployment

```bash
helm install go-app ./helm/go-app \
  --values ./helm/go-app/values-production.yaml \
  --set image.repository=myregistry/go-app \
  --set image.tag=1.0.0 \
  --namespace production \
  --create-namespace
```

### Scenario 4: Development Environment

```bash
helm install go-app ./helm/go-app \
  --set replicaCount=1 \
  --set autoscaling.enabled=false \
  --set resources.requests.cpu=100m \
  --set resources.requests.memory=128Mi \
  --set resources.limits.cpu=200m \
  --set resources.limits.memory=256Mi \
  --set env[0].name=ENVIRONMENT \
  --set env[0].value=development \
  --namespace dev \
  --create-namespace
```

## Configuration Examples

### Using ConfigMap for Configuration

Create a custom values file `my-values.yaml`:

```yaml
configMap:
  enabled: true
  data:
    APP_CONFIG: |
      {
        "feature_flags": {
          "new_feature": true
        }
      }

envFrom:
  - configMapRef:
      name: go-app
```

Install with:

```bash
helm install go-app ./helm/go-app -f my-values.yaml
```

### Using Secrets

Create secrets separately for security:

```bash
# Create a secret
kubectl create secret generic go-app-secrets \
  --from-literal=API_KEY=your-secret-key \
  --from-literal=DATABASE_PASSWORD=your-db-password

# Install referencing the secret
helm install go-app ./helm/go-app \
  --set envFrom[0].secretRef.name=go-app-secrets
```

### Enable Prometheus Monitoring

```bash
helm install go-app ./helm/go-app \
  --set serviceMonitor.enabled=true \
  --set serviceMonitor.labels.release=prometheus
```

### Configure Network Policies

```bash
helm install go-app ./helm/go-app \
  --set networkPolicy.enabled=true
```

## Istio-Specific Deployments

### Basic Istio Setup

```bash
# Ensure Istio is installed
kubectl get namespace istio-system

# Install with Istio enabled
helm install go-app ./helm/go-app \
  --set istio.enabled=true \
  --set istio.virtualService.enabled=true \
  --set istio.virtualService.gateways[0]=istio-system/main-gateway \
  --set istio.virtualService.hosts[0]=myapp.example.com
```

### Istio with Circuit Breaking

```bash
helm install go-app ./helm/go-app \
  --set istio.enabled=true \
  --set istio.virtualService.enabled=true \
  --set istio.destinationRule.enabled=true \
  --set istio.destinationRule.trafficPolicy.connectionPool.tcp.maxConnections=100 \
  --set istio.destinationRule.trafficPolicy.outlierDetection.consecutiveErrors=5
```

### Istio with mTLS

```bash
helm install go-app ./helm/go-app \
  --set istio.enabled=true \
  --set istio.peerAuthentication.enabled=true \
  --set istio.peerAuthentication.mtls.mode=STRICT
```

### Complete Istio Setup with Gateway

```bash
helm install go-app ./helm/go-app \
  --set istio.enabled=true \
  --set istio.gateway.enabled=true \
  --set istio.virtualService.enabled=true \
  --set istio.destinationRule.enabled=true \
  --set istio.peerAuthentication.enabled=true
```

## Upgrading

### Upgrade with New Image Version

```bash
helm upgrade go-app ./helm/go-app \
  --set image.tag=1.1.0 \
  --reuse-values
```

### Upgrade with New Values File

```bash
helm upgrade go-app ./helm/go-app \
  -f new-values.yaml
```

### Upgrade and Wait for Readiness

```bash
helm upgrade go-app ./helm/go-app \
  --set image.tag=1.1.0 \
  --wait \
  --timeout 5m
```

## Rolling Back

```bash
# View release history
helm history go-app

# Rollback to previous version
helm rollback go-app

# Rollback to specific revision
helm rollback go-app 2
```

## Testing

### Lint the Chart

```bash
helm lint ./helm/go-app
```

### Dry Run Installation

```bash
helm install go-app ./helm/go-app --dry-run --debug
```

### Template Rendering

```bash
# Render all templates
helm template go-app ./helm/go-app

# Render with values
helm template go-app ./helm/go-app -f my-values.yaml

# Show only specific template
helm template go-app ./helm/go-app --show-only templates/deployment.yaml
```

## Troubleshooting

### Check Generated Resources

```bash
# View all resources created by Helm
helm get manifest go-app

# Check specific resource
kubectl describe deployment go-app
kubectl describe hpa go-app
```

### View Logs

```bash
# View all pod logs
kubectl logs -l app.kubernetes.io/name=go-app --all-containers=true

# Follow logs
kubectl logs -l app.kubernetes.io/name=go-app -f
```

### Debug Failed Installation

```bash
# Get detailed information
helm status go-app

# Check events
kubectl get events --sort-by='.lastTimestamp'

# Describe pods
kubectl describe pods -l app.kubernetes.io/name=go-app
```

### Test Connectivity

```bash
# Port forward to test locally
kubectl port-forward svc/go-app 8080:80

# Test health endpoints
curl http://localhost:8080/health
curl http://localhost:8080/health/ready
curl http://localhost:8080/health/live
```

## Uninstalling

```bash
# Uninstall release
helm uninstall go-app

# Uninstall from specific namespace
helm uninstall go-app --namespace production
```

## Advanced Configuration

### Custom Resource Limits per Environment

Create environment-specific value files:

**values-dev.yaml:**
```yaml
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 200m
    memory: 256Mi
autoscaling:
  enabled: false
replicaCount: 1
```

**values-prod.yaml:**
```yaml
resources:
  requests:
    cpu: 500m
    memory: 512Mi
  limits:
    cpu: 1000m
    memory: 1Gi
autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 20
```

### Using with CI/CD

Example GitHub Actions workflow snippet:

```yaml
- name: Deploy to Kubernetes
  run: |
    helm upgrade --install go-app ./helm/go-app \
      --set image.tag=${{ github.sha }} \
      --set image.repository=${{ env.REGISTRY }}/go-app \
      --namespace production \
      --create-namespace \
      --wait \
      --timeout 5m
```

### Multi-Cluster Deployment

```bash
# Deploy to cluster 1
kubectl config use-context cluster-1
helm install go-app ./helm/go-app -f values-cluster1.yaml

# Deploy to cluster 2
kubectl config use-context cluster-2
helm install go-app ./helm/go-app -f values-cluster2.yaml
```

## Best Practices

1. **Always use version tags** instead of `latest` for production
2. **Use separate values files** for different environments
3. **Enable resource limits** to prevent resource exhaustion
4. **Enable HPA** for production workloads
5. **Use PodDisruptionBudgets** to maintain availability during updates
6. **Enable monitoring** with ServiceMonitor or annotations
7. **Use Network Policies** to restrict traffic
8. **Test in staging** before deploying to production
9. **Use `--dry-run`** to preview changes before applying
10. **Keep sensitive data in Secrets**, not in values files

## Support

For issues or questions:
- Check the [README.md](./go-app/README.md) for detailed configuration options
- Review the [values.yaml](./go-app/values.yaml) for all available settings
- Check application logs: `kubectl logs -l app.kubernetes.io/name=go-app`
