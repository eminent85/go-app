# CI/CD Quick Start Guide

Quick reference for common CI/CD tasks.

## Common Tasks

### 1. Deploy to Production

**Automatic (on merge to main):**
```bash
git checkout main
git pull
git merge feature-branch
git push origin main
# Automatically triggers full CI/CD pipeline
```

**Manual Deployment:**
```bash
# Go to: Actions → Production Deployment → Run workflow
# Select options:
#   - app_version: v1.2.3
#   - helm_version: 0.2.0
#   - namespace: production
#   - environment: production
```

---

### 2. Publish Helm Chart

**Automatic (on Helm changes):**
```bash
# Edit Helm charts
vim helm/go-app/values.yaml

git add helm/
git commit -m "feat: add new configuration option"
git push origin main
# Auto-bumps patch version (e.g., 0.1.0 → 0.1.1)
```

**Manual Version Bump:**
```bash
# Go to: Actions → Helm Chart Publish → Run workflow
# Select version bump: major, minor, or patch
```

**Local Version Bump:**
```bash
# Patch: 0.1.0 → 0.1.1
./.github/scripts/bump-helm-version.sh patch

# Minor: 0.1.0 → 0.2.0
./.github/scripts/bump-helm-version.sh minor

# Major: 0.1.0 → 1.0.0
./.github/scripts/bump-helm-version.sh major

git add .
git commit -m "chore: bump helm version"
git push origin main
```

---

### 3. Create Application Release

```bash
# Create and push a version tag
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3

# This triggers:
#   - Docker build with version tags
#   - GoReleaser for binaries
```

---

### 4. Deploy to Staging

```bash
# Go to: Actions → Production Deployment → Run workflow
# Set:
#   - namespace: staging
#   - environment: staging
```

---

### 5. Rollback Deployment

**Via Helm:**
```bash
# Check release history
helm history go-app -n production

# Rollback to previous release
helm rollback go-app -n production

# Rollback to specific revision
helm rollback go-app 5 -n production
```

**Via Kubernetes:**
```bash
# Rollback deployment
kubectl rollout undo deployment/go-app -n production

# Rollback to specific revision
kubectl rollout undo deployment/go-app --to-revision=3 -n production
```

**Via GitHub Actions:**
```bash
# Go to: Actions → Production Deployment → Run workflow
# Set previous versions:
#   - app_version: previous-tag
#   - helm_version: previous-version
```

---

### 6. Check Deployment Status

```bash
# Check pods
kubectl get pods -n production -l app.kubernetes.io/name=go-app

# Check deployment
kubectl get deployment go-app -n production

# Check HPA
kubectl get hpa go-app -n production

# Check logs
kubectl logs -n production -l app.kubernetes.io/name=go-app -f

# Check events
kubectl get events -n production --sort-by='.lastTimestamp'
```

---

### 7. Install Helm Chart Locally

```bash
# From OCI registry
helm install go-app oci://ghcr.io/<username>/charts/go-app --version 0.1.0

# From local files
helm install go-app ./helm/go-app

# With custom values
helm install go-app ./helm/go-app -f my-values.yaml

# Dry run (test without installing)
helm install go-app ./helm/go-app --dry-run --debug
```

---

### 8. Update Existing Deployment

```bash
# Upgrade with new values
helm upgrade go-app ./helm/go-app -f new-values.yaml -n production

# Upgrade with new image
helm upgrade go-app ./helm/go-app \
  --set image.tag=v1.2.3 \
  --reuse-values \
  -n production

# Upgrade from OCI registry
helm upgrade go-app oci://ghcr.io/<username>/charts/go-app \
  --version 0.2.0 \
  -n production
```

---

### 9. Test Helm Chart Locally

```bash
# Lint chart
helm lint ./helm/go-app

# Template rendering (default values)
helm template go-app ./helm/go-app

# Template with production values
helm template go-app ./helm/go-app -f ./helm/go-app/values-production.yaml

# Template with Istio values
helm template go-app ./helm/go-app -f ./helm/go-app/values-istio.yaml

# Show specific template
helm template go-app ./helm/go-app --show-only templates/deployment.yaml
```

---

### 10. View Helm Chart in Registry

```bash
# Show chart metadata
helm show chart oci://ghcr.io/<username>/charts/go-app --version 0.1.0

# Show chart values
helm show values oci://ghcr.io/<username>/charts/go-app --version 0.1.0

# Show all chart information
helm show all oci://ghcr.io/<username>/charts/go-app --version 0.1.0

# Pull chart locally
helm pull oci://ghcr.io/<username>/charts/go-app --version 0.1.0
```

---

## Pre-Deployment Checklist

Before deploying to production:

- [ ] Tests pass locally (`go test ./...`)
- [ ] Code linted (`golangci-lint run`)
- [ ] Helm chart linted (`helm lint ./helm/go-app`)
- [ ] Helm templates validated (`helm template go-app ./helm/go-app`)
- [ ] Changes tested in staging environment
- [ ] Database migrations tested (if applicable)
- [ ] Environment variables configured
- [ ] Secrets created in Kubernetes
- [ ] Resource limits appropriate for workload
- [ ] HPA thresholds reviewed
- [ ] Monitoring/alerts configured

---

## Troubleshooting Commands

```bash
# Check workflow runs
gh run list
gh run view <run-id>
gh run watch

# Check pod issues
kubectl describe pod <pod-name> -n production
kubectl logs <pod-name> -n production --previous

# Check service endpoints
kubectl get endpoints go-app -n production

# Port forward for local testing
kubectl port-forward -n production svc/go-app 8080:80

# Test health endpoints
curl http://localhost:8080/health
curl http://localhost:8080/health/ready
curl http://localhost:8080/health/live
curl http://localhost:8080/metrics

# Check HPA metrics
kubectl describe hpa go-app -n production
kubectl top pods -n production

# Check image pull issues
kubectl get events -n production | grep "Failed to pull"
```

---

## Required Secrets

Set these in GitHub repository settings (Settings → Secrets and variables → Actions):

1. **`KUBECONFIG`** (required for deployment)
   ```bash
   cat ~/.kube/config | base64 -w 0
   # Add output as GitHub secret
   ```

2. **`CODECOV_TOKEN`** (optional, for private repos)
   - Get from codecov.io

3. **Notification secrets** (optional)
   - `SLACK_WEBHOOK`
   - `DISCORD_WEBHOOK`

---

## Environment URLs

After deployment, access your application:

**Production:**
- URL: https://myapp.example.com (configure in values)
- Health: https://myapp.example.com/health
- Metrics: https://myapp.example.com/metrics

**Staging:**
- URL: https://staging.myapp.example.com
- Health: https://staging.myapp.example.com/health

---

## Workflow Status Badges

Add to your README.md:

```markdown
![CI](https://github.com/<username>/go-app/actions/workflows/ci.yml/badge.svg)
![Docker](https://github.com/<username>/go-app/actions/workflows/docker.yml/badge.svg)
![Helm](https://github.com/<username>/go-app/actions/workflows/helm-publish.yml/badge.svg)
![Deploy](https://github.com/<username>/go-app/actions/workflows/production-deploy.yml/badge.svg)
```

---

## Support

- Full documentation: [WORKFLOWS.md](./WORKFLOWS.md)
- Helm chart docs: [helm/go-app/README.md](../helm/go-app/README.md)
- Installation guide: [helm/INSTALL.md](../helm/INSTALL.md)
