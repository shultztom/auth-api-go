#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

IMAGE=us-central1-docker.pkg.dev/shultzlab/docker-gcp/auth-api-go
VERSION=$(tr -d '[:space:]' < VERSION)

if [ -z "$VERSION" ]; then
  echo "VERSION file is empty" >&2
  exit 1
fi

# Never overwrite a version that's already been pushed
if docker manifest inspect "$IMAGE:$VERSION" >/dev/null 2>&1; then
  echo "$IMAGE:$VERSION already exists in the registry. Bump VERSION first." >&2
  exit 1
fi

docker build -t "$IMAGE:$VERSION" -t "$IMAGE:latest" .
docker push "$IMAGE:$VERSION"
docker push "$IMAGE:latest"
