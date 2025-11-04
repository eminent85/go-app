# Go App Helm Chart

This Helm chart deploys the Go application to Kubernetes with comprehensive configuration options.

## Features

- **Horizontal Pod Autoscaling (HPA)** - Automatically scale based on CPU and memory utilization
- **Health Checks** - Built-in liveness, readiness, and startup probes
- **Flexible Ingress** - Support for standard Kubernetes Ingress or Istio
- **Security** - Pod security contexts, network policies, RBAC
- **High Availability** - Pod anti-affinity rules and pod disruption budgets
- **Monitoring** - ServiceMonitor for Prometheus Operator
- **Istio Support** - VirtualService, DestinationRule, Gateway, and PeerAuthentication resources

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- (Optional) Istio 1.10+ if using Istio features
- (Optional) Prometheus Operator if using ServiceMonitor

## Installing the Chart

### Basic Installation

```bash
helm install go-app ./helm/go-app
```

### With Custom Values

```bash
helm install go-app ./helm/go-app -f myvalues.yaml
```

### Install with Ingress

```bash
helm install go-app ./helm/go-app \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=myapp.example.com \
  --set ingress.hosts[0].paths[0].path=/ \
  --set ingress.hosts[0].paths[0].pathType=Prefix
```

### Install with Istio

```bash
helm install go-app ./helm/go-app \
  --set istio.enabled=true \
  --set istio.virtualService.enabled=true \
  --set istio.virtualService.hosts[0]=myapp.example.com
```

## Configuration

The following table lists the configurable parameters and their default values.

### General Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `2` |
| `image.repository` | Image repository | `your-registry/go-app` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `image.tag` | Image tag (defaults to chart appVersion) | `""` |
| `nameOverride` | Override chart name | `""` |
| `fullnameOverride` | Override full name | `""` |

### Service Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `80` |
| `service.targetPort` | Container port | `8080` |
| `service.annotations` | Service annotations | `{}` |

### Autoscaling Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `autoscaling.enabled` | Enable HPA | `true` |
| `autoscaling.minReplicas` | Minimum replicas | `2` |
| `autoscaling.maxReplicas` | Maximum replicas | `10` |
| `autoscaling.targetCPUUtilizationPercentage` | Target CPU utilization | `75` |
| `autoscaling.targetMemoryUtilizationPercentage` | Target memory utilization | `80` |

### Ingress Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ingress.enabled` | Enable ingress | `false` |
| `ingress.className` | Ingress class name | `""` |
| `ingress.annotations` | Ingress annotations | `{}` |
| `ingress.hosts` | Ingress hosts configuration | See values.yaml |

### Istio Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `istio.enabled` | Enable Istio integration | `false` |
| `istio.virtualService.enabled` | Create VirtualService | `false` |
| `istio.destinationRule.enabled` | Create DestinationRule | `false` |
| `istio.peerAuthentication.enabled` | Create PeerAuthentication | `false` |
| `istio.gateway.enabled` | Create Gateway | `false` |

### Resource Limits

| Parameter | Description | Default |
|-----------|-------------|---------|
| `resources.limits.cpu` | CPU limit | `500m` |
| `resources.limits.memory` | Memory limit | `512Mi` |
| `resources.requests.cpu` | CPU request | `250m` |
| `resources.requests.memory` | Memory request | `256Mi` |

### Monitoring

| Parameter | Description | Default |
|-----------|-------------|---------|
| `serviceMonitor.enabled` | Enable Prometheus ServiceMonitor | `false` |
| `serviceMonitor.interval` | Scrape interval | `30s` |
| `serviceMonitor.path` | Metrics path | `/metrics` |

## Examples

### Example 1: Basic Deployment

```yaml
replicaCount: 3
image:
  repository: myregistry/go-app
  tag: "1.0.0"
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 200m
    memory: 256Mi
```

### Example 2: With Standard Ingress (NGINX)

```yaml
ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: myapp.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: myapp-tls
      hosts:
        - myapp.example.com
```

### Example 3: With Istio Service Mesh

```yaml
istio:
  enabled: true
  virtualService:
    enabled: true
    gateways:
      - istio-system/main-gateway
    hosts:
      - myapp.example.com
    http:
      - match:
          - uri:
              prefix: /
        route:
          - destination:
              host: go-app
              port:
                number: 80
        timeout: 30s
        retries:
          attempts: 3
          perTryTimeout: 10s
  destinationRule:
    enabled: true
    trafficPolicy:
      connectionPool:
        tcp:
          maxConnections: 100
        http:
          http1MaxPendingRequests: 50
          maxRequestsPerConnection: 2
      loadBalancer:
        simple: LEAST_CONN
  peerAuthentication:
    enabled: true
    mtls:
      mode: STRICT
```

### Example 4: Production Configuration with HA

```yaml
replicaCount: 3
autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 20
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 75

resources:
  requests:
    cpu: 500m
    memory: 512Mi
  limits:
    cpu: 1000m
    memory: 1Gi

podDisruptionBudget:
  enabled: true
  minAvailable: 2

affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
            - key: app.kubernetes.io/name
              operator: In
              values:
                - go-app
        topologyKey: kubernetes.io/hostname

topologySpreadConstraints:
  - maxSkew: 1
    topologyKey: topology.kubernetes.io/zone
    whenUnsatisfiable: DoNotSchedule
    labelSelector:
      matchLabels:
        app.kubernetes.io/name: go-app

serviceMonitor:
  enabled: true
  interval: 15s

networkPolicy:
  enabled: true
```

## Upgrading

```bash
helm upgrade go-app ./helm/go-app -f myvalues.yaml
```

## Uninstalling

```bash
helm uninstall go-app
```

## Health Checks

The application exposes the following health endpoints:

- `/health` - General health check
- `/health/ready` - Readiness probe endpoint
- `/health/live` - Liveness probe endpoint
- `/metrics` - Prometheus metrics endpoint

## Troubleshooting

### Check pod status

```bash
kubectl get pods -l app.kubernetes.io/name=go-app
```

### View logs

```bash
kubectl logs -l app.kubernetes.io/name=go-app -f
```

### Check HPA status

```bash
kubectl get hpa
```

### Verify Istio configuration

```bash
kubectl get virtualservices
kubectl get destinationrules
kubectl get peerauthentications
```

## Contributing

Please ensure all Helm chart changes are tested before submitting.

## License

See the main project LICENSE file.
