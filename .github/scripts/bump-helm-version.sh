#!/bin/bash
set -e

# Helm chart version bumping script
# Usage: ./bump-helm-version.sh [major|minor|patch]

BUMP_TYPE="${1:-patch}"
VERSION_FILE=".github/helm-version.txt"
CHART_FILE="helm/go-app/Chart.yaml"

if [ ! -f "$VERSION_FILE" ]; then
  echo "Error: Version file $VERSION_FILE not found"
  exit 1
fi

if [ ! -f "$CHART_FILE" ]; then
  echo "Error: Chart file $CHART_FILE not found"
  exit 1
fi

# Read current version
CURRENT_VERSION=$(cat "$VERSION_FILE")
echo "Current Helm chart version: $CURRENT_VERSION"

# Parse version components
IFS='.' read -r major minor patch <<< "$CURRENT_VERSION"

# Bump version based on type
case $BUMP_TYPE in
  major)
    major=$((major + 1))
    minor=0
    patch=0
    ;;
  minor)
    minor=$((minor + 1))
    patch=0
    ;;
  patch)
    patch=$((patch + 1))
    ;;
  *)
    echo "Error: Invalid bump type '$BUMP_TYPE'. Use major, minor, or patch."
    exit 1
    ;;
esac

NEW_VERSION="$major.$minor.$patch"
echo "New Helm chart version: $NEW_VERSION"

# Update version file
echo "$NEW_VERSION" > "$VERSION_FILE"
echo "Updated $VERSION_FILE"

# Update Chart.yaml
if [[ "$OSTYPE" == "darwin"* ]]; then
  # macOS
  sed -i '' "s/^version:.*/version: $NEW_VERSION/" "$CHART_FILE"
else
  # Linux
  sed -i "s/^version:.*/version: $NEW_VERSION/" "$CHART_FILE"
fi
echo "Updated $CHART_FILE"

echo "Helm chart version bumped to $NEW_VERSION"
