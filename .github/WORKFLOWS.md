# GitHub Actions Workflows Documentation

This document describes all GitHub Actions workflows configured for this project.

## Overview

The project uses multiple CI/CD workflows to automate testing, building, publishing, and deploying the Go application and its Helm charts.

## Workflows

### 1. CI Workflow (`ci.yml`)

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**Jobs:**
1. **Test** - Runs tests across Go versions 1.21 and 1.22
   - Downloads dependencies
   - Runs `go vet`
   - Runs tests with race detection and coverage
   - Uploads coverage to Codecov

2. **Lint** - Runs golangci-lint
   - Uses golangci-lint with latest version
   - 5-minute timeout

3. **Build** - Builds the application
   - Compiles the binary
   - Uploads artifact for 7 days

4. **Security** - Security scanning
   - Runs Gosec security scanner
   - Uploads results to GitHub Security

5. **Helm Lint** - Validates Helm charts
   - Lints the Helm chart
   - Validates templates with all value files

**Status:** Required for merging PRs

---

### 2. Docker Build and Push (`docker.yml`)

**Triggers:**
- Push to `main` branch
- Push of version tags (`v*`)
- Pull requests to `main` branch

**Features:**
- Multi-architecture builds (amd64, arm64)
- Pushes to GitHub Container Registry (ghcr.io)
- Automatic tagging:
  - Branch name (e.g., `main`)
  - Git SHA (e.g., `main-abc1234`)
  - Semantic versions (e.g., `v1.2.3`, `1.2`, `1`)
  - `latest` tag for main branch
- Trivy vulnerability scanning
- Build cache optimization

**Registry:** `ghcr.io/<username>/go-app`

---

### 3. Release (`release.yml`)

**Triggers:**
- Push of version tags (`v*`)

**Features:**
- Uses GoReleaser for creating releases
- Builds binaries for multiple platforms
- Creates GitHub releases with binaries attached

**Permissions:**
- `contents: write` - for creating releases
- `packages: write` - for publishing packages

---

### 4. Helm Chart Publish (`helm-publish.yml`)

**Triggers:**
- Push to `main` branch when Helm files change
- Manual workflow dispatch with version bump type selection

**Version Control:**
- Independent versioning from application
- Stored in `.github/helm-version.txt`
- Supports semantic versioning (major, minor, patch)

**Jobs:**
1. **Package and Publish Helm Chart**
   - Reads current Helm version
   - Calculates new version (auto-patch or manual selection)
   - Updates `Chart.yaml` and version file
   - Lints and packages the chart
   - Pushes to OCI registry at `ghcr.io/<username>/charts/go-app`
   - Creates GitHub release with tag `helm-v<version>`
   - Commits version bump back to repository

**Manual Usage:**
```bash
# Trigger manually via GitHub UI
Actions → Helm Chart Publish → Run workflow → Select version bump type
```

**Automatic Usage:**
```bash
# Automatically triggers on Helm chart changes
git add helm/
git commit -m "feat: update helm chart configuration"
git push origin main
```

**Chart Registry:** `oci://ghcr.io/<username>/charts/go-app`

---

### 5. Production Deployment (`production-deploy.yml`)

**Triggers:**
- Push to `main` branch (automatic deployment)
- Manual workflow dispatch with configuration options

**Manual Inputs:**
- `app_version` - Docker image tag (default: latest commit SHA)
- `helm_version` - Helm chart version (default: from `.github/helm-version.txt`)
- `namespace` - Kubernetes namespace (default: `production`)
- `environment` - Environment name (production or staging)

**Jobs:**

1. **Prepare** - Determines deployment configuration
   - Resolves app and Helm versions
   - Sets namespace and environment

2. **Build Docker** - Builds and pushes Docker image (only on push events)
   - Multi-architecture build
   - Pushes to ghcr.io
   - Runs Trivy security scan

3. **Deploy to Kubernetes** - Deploys using Helm
   - Configures kubectl access
   - Creates namespace if needed
   - Pulls Helm chart from OCI registry
   - Creates image pull secret
   - Deploys with production values:
     - 3 replicas minimum
     - HPA enabled (3-20 replicas)
     - Resource limits configured
     - PodDisruptionBudget enabled
     - ServiceMonitor enabled
   - Verifies deployment
   - Runs smoke tests on health endpoints

4. **Notify Deployment** - Sends deployment status notifications
   - Success/failure notifications
   - Can be extended with Slack, email, etc.

**Prerequisites:**
- Kubernetes cluster access configured
- `KUBECONFIG` secret set in GitHub
- Container registry credentials

**Example Manual Deployment:**
```bash
# Via GitHub UI:
Actions → Production Deployment → Run workflow
  app_version: v1.2.3
  helm_version: 0.2.0
  namespace: production
  environment: production
```

---

## Helm Version Management

### Version File
- **Location:** `.github/helm-version.txt`
- **Format:** Semantic version (e.g., `0.1.0`)
- **Purpose:** Track Helm chart version independently from application version

### Version Bumping

#### Automatic (via workflow)
Changes to Helm files trigger automatic patch version bump:
```bash
git add helm/
git commit -m "fix: update deployment resources"
git push origin main
# Automatically bumps from 0.1.0 → 0.1.1
```

#### Manual (via script)
```bash
# Patch bump (0.1.0 → 0.1.1)
./.github/scripts/bump-helm-version.sh patch

# Minor bump (0.1.0 → 0.2.0)
./.github/scripts/bump-helm-version.sh minor

# Major bump (0.1.0 → 1.0.0)
./.github/scripts/bump-helm-version.sh major
```

#### Manual (via GitHub Actions)
```bash
# Via GitHub UI:
Actions → Helm Chart Publish → Run workflow → Select version bump type
```

---

## Secrets Required

Configure these secrets in GitHub repository settings:

### Required Secrets

1. **`GITHUB_TOKEN`**
   - **Auto-provided by GitHub Actions**
   - Used for: Container registry, GitHub releases

2. **`KUBECONFIG`** (for production deployment)
   - **Base64-encoded kubeconfig file**
   - Used for: Kubernetes cluster access
   - Generate:
     ```bash
     cat ~/.kube/config | base64 -w 0
     ```

### Optional Secrets

3. **`CODECOV_TOKEN`** (optional)
   - For private repository coverage uploads
   - Get from: codecov.io

4. **Notification Secrets** (optional)
   - `SLACK_WEBHOOK` - For Slack notifications
   - `DISCORD_WEBHOOK` - For Discord notifications
   - Configure in `notify-deployment` job

---

## Environment Configuration

### GitHub Environments

Configure environments in repository settings:

1. **Production Environment**
   - Protection rules recommended
   - Required reviewers
   - Deployment branches: `main` only

2. **Staging Environment**
   - Less restrictive rules
   - For testing before production

---

## Complete Deployment Flow

### On Main Branch Merge

```mermaid
Push to main
  ├─→ CI Workflow (test, lint, build, security, helm-lint)
  ├─→ Docker Workflow (build & push image)
  ├─→ Helm Publish (if helm files changed)
  └─→ Production Deploy
       ├─→ Build Docker (multi-arch)
       ├─→ Deploy to Kubernetes
       └─→ Notify deployment status
```

### On Version Tag

```mermaid
Push tag v*
  ├─→ Docker Workflow (build & push with version tags)
  └─→ Release Workflow (GoReleaser)
```

---

## Usage Examples

### Deploy Specific Versions

```yaml
# Deploy specific app and helm versions to staging
Actions → Production Deployment → Run workflow
  app_version: v1.2.3
  helm_version: 0.2.0
  namespace: staging
  environment: staging
```

### Publish New Helm Chart Version

```yaml
# Bump helm chart minor version
Actions → Helm Chart Publish → Run workflow
  version_bump: minor
```

### Emergency Rollback

```bash
# Rollback to previous Helm release
kubectl rollout undo deployment/go-app -n production

# Or deploy specific version via workflow
Actions → Production Deployment → Run workflow
  app_version: previous-version-tag
  helm_version: previous-chart-version
```

---

## Troubleshooting

### Failed Deployment

1. Check workflow logs in Actions tab
2. Verify Kubernetes access:
   ```bash
   kubectl get nodes
   ```
3. Check pod status:
   ```bash
   kubectl get pods -n production
   kubectl logs -n production -l app.kubernetes.io/name=go-app
   ```

### Helm Chart Issues

1. Lint locally:
   ```bash
   helm lint ./helm/go-app
   ```
2. Test template rendering:
   ```bash
   helm template go-app ./helm/go-app --debug
   ```
3. Check chart in registry:
   ```bash
   helm show chart oci://ghcr.io/<username>/charts/go-app --version <version>
   ```

### Version Conflicts

If version file and Chart.yaml are out of sync:
```bash
# Manually sync versions
./.github/scripts/bump-helm-version.sh patch
git add .github/helm-version.txt helm/go-app/Chart.yaml
git commit -m "chore: sync helm versions"
git push origin main
```

---

## Best Practices

1. **Always test Helm changes locally** before pushing
   ```bash
   helm lint ./helm/go-app
   helm template go-app ./helm/go-app
   ```

2. **Use semantic versioning** for both app and Helm
   - App: `v1.2.3` (Git tags)
   - Helm: `0.1.0` (Independent versioning)

3. **Test in staging first** for critical changes
   ```bash
   # Deploy to staging
   helm upgrade --install go-app ./helm/go-app -n staging
   ```

4. **Monitor deployments** in production
   - Check pod status
   - Review metrics
   - Monitor logs

5. **Use manual approval** for production deployments
   - Configure in GitHub Environment protection rules

6. **Keep Helm chart backward compatible**
   - Don't break existing deployments
   - Use proper version bumps (major for breaking changes)

---

## Workflow Diagram

```
┌─────────────────────────────────────────────────────────┐
│                    Push to Main                         │
└─────────┬───────────────────────────────────────────────┘
          │
          ├─→ CI (Test, Lint, Build, Security, Helm Lint)
          │
          ├─→ Docker (Build & Push Multi-Arch Image)
          │
          ├─→ Helm Publish (if helm/** changed)
          │   ├─→ Auto-bump version (patch)
          │   ├─→ Package chart
          │   ├─→ Push to OCI registry
          │   └─→ Create GitHub release
          │
          └─→ Production Deploy
              ├─→ Prepare (determine versions)
              ├─→ Build Docker (if push event)
              ├─→ Deploy to K8s
              │   ├─→ Pull Helm chart
              │   ├─→ Apply configuration
              │   ├─→ Verify deployment
              │   └─→ Run smoke tests
              └─→ Notify deployment status
```

---

## Support

For issues or questions:
- Check workflow run logs in Actions tab
- Review this documentation
- Check Helm chart [README](../helm/go-app/README.md)
- Check [INSTALL guide](../helm/INSTALL.md)
