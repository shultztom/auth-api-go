#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

# Optional first arg overrides the tag (e.g. ./helmDeploy.sh 1.0.0 to redeploy an older version)
VERSION=${1:-$(tr -d '[:space:]' < VERSION)}

if [ -z "$VERSION" ]; then
  echo "No version given and VERSION file is empty" >&2
  exit 1
fi

gcloud auth print-access-token | helm registry login -u oauth2accesstoken \
--password-stdin https://us-central1-docker.pkg.dev

helm upgrade --install auth-api-go oci://us-central1-docker.pkg.dev/shultzlab/shultzlab-helm/http-deploy --version 0.2.2 -f .helm/deploy/values.yaml --set image.tag="$VERSION"
