#!/usr/bin/env bash
set -euo pipefail

# Publish all npm packages for the current release.
# Prerequisites:
#   1. Run `make release` first to produce dist/ tarballs
#   2. Be logged in to npm: `npm login`
#   3. You must own the etoro-cli* package names on npm

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
VERSION=$(node -p "require('${SCRIPT_DIR}/cli/package.json').version")

echo "Publishing etoro-cli v${VERSION}"
echo ""

PLATFORMS=(
  "cli-darwin-arm64"
  "cli-darwin-x64"
  "cli-linux-arm64"
  "cli-linux-x64"
  "cli-windows-arm64"
  "cli-windows-x64"
)

for platform in "${PLATFORMS[@]}"; do
  dir="${SCRIPT_DIR}/${platform}"
  pkg_version=$(node -p "require('${dir}/package.json').version")
  if [ "$pkg_version" != "$VERSION" ]; then
    echo "ERROR: ${platform}/package.json version (${pkg_version}) != ${VERSION}"
    exit 1
  fi
done

# Extract fresh binaries from dist/ into each platform package
echo "Extracting binaries from dist/..."
tar -xzf "${SCRIPT_DIR}/../dist/etoro_darwin_arm64.tar.gz" -C "${SCRIPT_DIR}/cli-darwin-arm64/bin/"
tar -xzf "${SCRIPT_DIR}/../dist/etoro_darwin_amd64.tar.gz"  -C "${SCRIPT_DIR}/cli-darwin-x64/bin/"
tar -xzf "${SCRIPT_DIR}/../dist/etoro_linux_arm64.tar.gz"   -C "${SCRIPT_DIR}/cli-linux-arm64/bin/"
tar -xzf "${SCRIPT_DIR}/../dist/etoro_linux_amd64.tar.gz"   -C "${SCRIPT_DIR}/cli-linux-x64/bin/"
unzip -o "${SCRIPT_DIR}/../dist/etoro_windows_arm64.zip"    -d "${SCRIPT_DIR}/cli-windows-arm64/bin/" >/dev/null
unzip -o "${SCRIPT_DIR}/../dist/etoro_windows_amd64.zip"    -d "${SCRIPT_DIR}/cli-windows-x64/bin/"  >/dev/null
chmod +x "${SCRIPT_DIR}"/cli-darwin-*/bin/etoro "${SCRIPT_DIR}"/cli-linux-*/bin/etoro

# Verify binaries exist
for platform in "${PLATFORMS[@]}"; do
  if [[ "${platform}" == cli-windows-* ]]; then
    bin_name="etoro.exe"
  else
    bin_name="etoro"
  fi
  bin_path="${SCRIPT_DIR}/${platform}/bin/${bin_name}"
  test -f "${bin_path}" || {
    echo "ERROR: ${bin_path} missing"; exit 1;
  }
done

# Publish platform packages first
for platform in "${PLATFORMS[@]}"; do
  echo ""
  echo "Publishing etoro-cli-${platform#cli-}..."
  (cd "${SCRIPT_DIR}/${platform}" && npm publish --access public)
done

# Publish main package last (depends on platform packages)
echo ""
echo "Publishing etoro-cli..."
(cd "${SCRIPT_DIR}/cli" && npm publish --access public)

echo ""
echo "Done! All packages published for v${VERSION}"
