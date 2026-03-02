#!/usr/bin/env bash
set -euo pipefail

# Bump version across all npm packages.
# Usage: ./bump-version.sh <new-version>

if [ $# -ne 1 ]; then
  echo "Usage: $0 <new-version>"
  echo "Example: $0 0.2.0"
  exit 1
fi

NEW_VERSION="$1"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

PACKAGES=(
  "cli"
  "cli-darwin-arm64"
  "cli-darwin-x64"
  "cli-linux-arm64"
  "cli-linux-x64"
)

for pkg in "${PACKAGES[@]}"; do
  file="${SCRIPT_DIR}/${pkg}/package.json"
  node -e "
    const fs = require('fs');
    const p = JSON.parse(fs.readFileSync('${file}', 'utf8'));
    p.version = '${NEW_VERSION}';
    if (p.optionalDependencies) {
      for (const k of Object.keys(p.optionalDependencies)) {
        p.optionalDependencies[k] = '${NEW_VERSION}';
      }
    }
    fs.writeFileSync('${file}', JSON.stringify(p, null, 2) + '\n');
  "
  echo "  ${pkg}/package.json → ${NEW_VERSION}"
done

echo ""
echo "All packages bumped to v${NEW_VERSION}"
