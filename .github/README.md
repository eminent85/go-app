# GitHub Actions CI/CD

This directory contains all GitHub Actions workflows, scripts, and configuration for the CI/CD pipeline.

## Directory Structure

```
.github/
├── workflows/
│   ├── ci.yml                    # Continuous Integration (test, lint, build)
│   ├── docker.yml                # Docker image build and push
│   ├── release.yml               # GoReleaser for application releases
│   ├── helm-publish.yml          # Helm chart packaging and publishing
│   └── production-deploy.yml     # Production deployment workflow
├── scripts/
│   └── bump-helm-version.sh      # Helm version bumping utility
├── helm-version.txt              # Current Helm chart version
├── WORKFLOWS.md                  # Detailed workflow documentation
├── QUICKSTART.md                 # Quick reference guide
└── README.md                     # This file
```

## Quick Links

- **[Quick Start Guide](./QUICKSTART.md)** - Common tasks and commands
- **[Detailed Workflows Documentation](./WORKFLOWS.md)** - Complete workflow reference
- **[Helm Chart Documentation](../helm/go-app/README.md)** - Helm chart configuration
- **[Helm Installation Guide](../helm/INSTALL.md)** - Helm deployment examples

## Workflows Overview

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| **CI** | Push/PR to main/develop | Run tests, lint, build, security scans |
| **Docker** | Push to main, version tags | Build and push Docker images |
| **Release** | Version tags (v*) | Create GitHub releases with binaries |
| **Helm Publish** | Helm changes, manual | Package and publish Helm charts |
| **Production Deploy** | Push to main, manual | Deploy to Kubernetes |

## Version Management

### Application Version
- **Location:** Git tags (e.g., `v1.2.3`)
- **Format:** Semantic versioning with `v` prefix
- **Managed by:** Git tags + GoReleaser

### Helm Chart Version
- **Location:** `.github/helm-version.txt`
- **Format:** Semantic versioning (e.g., `0.1.0`)
- **Managed by:** Automated workflow or manual script
- **Independent from application version**

## Getting Started

### 1. First Time Setup

```bash
# Ensure you have required secrets configured
# GitHub Settings → Secrets and variables → Actions

# Required:
# - KUBECONFIG (base64 encoded)

# Optional:
# - CODECOV_TOKEN
# - SLACK_WEBHOOK
```

### 2. Deploy to Production

```bash
# Automatic on merge to main
git push origin main

# Manual via GitHub UI
# Actions → Production Deployment → Run workflow
```

### 3. Publish Helm Chart

```bash
# Automatic on helm changes
git add helm/
git commit -m "feat: update helm configuration"
git push origin main

# Manual version bump
./.github/scripts/bump-helm-version.sh minor
git push origin main
```

### 4. Create Application Release

```bash
# Tag and push
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

## Workflow Execution Flow

### On Push to Main
```
1. CI Workflow
   ├─ Test (Go 1.21, 1.22)
   ├─ Lint (golangci-lint)
   ├─ Build
   ├─ Security Scan
   └─ Helm Lint

2. Docker Workflow
   └─ Build & Push Multi-Arch Image

3. Helm Publish (if helm/** changed)
   ├─ Auto-bump version
   ├─ Package chart
   ├─ Push to OCI registry
   └─ Create GitHub release

4. Production Deploy
   ├─ Build Docker
   ├─ Deploy to Kubernetes
   └─ Run smoke tests
```

### On Version Tag (v*)
```
1. Docker Workflow
   └─ Build & Push with version tags

2. Release Workflow
   └─ GoReleaser (binaries + GitHub release)
```

## Common Tasks

See [QUICKSTART.md](./QUICKSTART.md) for detailed commands.

**Deploy to production:**
```bash
git push origin main  # Automatic
```

**Bump Helm version:**
```bash
./.github/scripts/bump-helm-version.sh patch
```

**Rollback deployment:**
```bash
helm rollback go-app -n production
```

**View logs:**
```bash
kubectl logs -n production -l app.kubernetes.io/name=go-app -f
```

## Configuration

### Kubernetes Cluster Access

The production deployment workflow requires Kubernetes access. Configure based on your provider:

**Generic (using KUBECONFIG secret):**
```bash
cat ~/.kube/config | base64 -w 0
# Add as KUBECONFIG secret in GitHub
```

**AWS EKS:**
Add to workflow in `production-deploy.yml`:
```yaml
- name: Configure AWS credentials
  uses: aws-actions/configure-aws-credentials@v4
  with:
    aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
    aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
    aws-region: us-east-1

- name: Configure kubectl
  run: aws eks update-kubeconfig --name my-cluster --region us-east-1
```

**GCP GKE:**
```yaml
- name: Authenticate to Google Cloud
  uses: google-github-actions/auth@v2
  with:
    credentials_json: ${{ secrets.GCP_SA_KEY }}

- name: Configure kubectl
  run: gcloud container clusters get-credentials my-cluster --region us-central1
```

**Azure AKS:**
```yaml
- name: Azure Login
  uses: azure/login@v2
  with:
    creds: ${{ secrets.AZURE_CREDENTIALS }}

- name: Configure kubectl
  run: az aks get-credentials --resource-group my-rg --name my-cluster
```

## Monitoring

### Workflow Status
- Check: GitHub → Actions tab
- Status badges available in [QUICKSTART.md](./QUICKSTART.md)

### Deployment Status
```bash
# Check pods
kubectl get pods -n production

# Check HPA
kubectl get hpa -n production

# Check recent events
kubectl get events -n production --sort-by='.lastTimestamp'
```

## Troubleshooting

### Workflow Failures

1. **CI failures:** Check test output in Actions logs
2. **Docker build failures:** Check Dockerfile and dependencies
3. **Helm publish failures:** Run `helm lint ./helm/go-app` locally
4. **Deployment failures:** Check Kubernetes events and pod logs

### Common Issues

**"Failed to pull image"**
```bash
# Ensure image pull secret exists
kubectl get secret ghcr-secret -n production

# Recreate if needed (done automatically by workflow)
```

**"Deployment timeout"**
```bash
# Check pod status
kubectl describe pod <pod-name> -n production

# Check events
kubectl get events -n production
```

**"Helm chart not found"**
```bash
# Verify chart in registry
helm show chart oci://ghcr.io/<username>/charts/go-app --version <version>

# Pull manually to test
helm pull oci://ghcr.io/<username>/charts/go-app --version <version>
```

## Best Practices

1. **Test locally before pushing**
   - Run tests: `go test ./...`
   - Lint code: `golangci-lint run`
   - Lint Helm: `helm lint ./helm/go-app`

2. **Use semantic versioning**
   - App: Major.Minor.Patch (v1.2.3)
   - Helm: Major.Minor.Patch (0.1.0)

3. **Deploy to staging first**
   - Test critical changes in staging
   - Use manual workflow dispatch

4. **Monitor deployments**
   - Watch pod rollout status
   - Check application metrics
   - Review logs

5. **Use protected branches**
   - Require PR reviews for main
   - Require status checks to pass

## Security

### Secrets Management
- Never commit secrets to Git
- Use Kubernetes Secrets for sensitive data
- Rotate credentials regularly
- Use least-privilege access

### Image Scanning
- Trivy scans run automatically
- Review security findings in GitHub Security tab
- Address critical vulnerabilities promptly

### Network Security
- NetworkPolicies configured in Helm chart
- Istio mTLS available for service mesh
- Enable as needed for your environment

## Contributing

When adding new workflows:
1. Document in [WORKFLOWS.md](./WORKFLOWS.md)
2. Add common tasks to [QUICKSTART.md](./QUICKSTART.md)
3. Test thoroughly before merging
4. Use descriptive job and step names

## Support

- **Workflow issues:** Check [WORKFLOWS.md](./WORKFLOWS.md)
- **Helm issues:** Check [helm/go-app/README.md](../helm/go-app/README.md)
- **Deployment issues:** Check [helm/INSTALL.md](../helm/INSTALL.md)
- **General questions:** Open an issue in the repository

## License

See the main project LICENSE file.
