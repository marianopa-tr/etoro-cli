#!/usr/bin/env bash
set -euo pipefail

# eToro CLI installer
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/marianopa-tr/etoro-cli/main/scripts/install.sh | bash
#
# Environment overrides:
#   ETORO_VERSION        Release tag (default: latest)
#   ETORO_INSTALL_DIR    Install directory (default: OS-dependent)
#   ETORO_INSTALL_URL    Direct URL to a .tar.gz asset (skips auto-detection)
#   ETORO_INSTALL_OS     Override OS detection  (darwin, linux)
#   ETORO_INSTALL_ARCH   Override arch detection (amd64, arm64)
#   ETORO_SKIP_CHECKSUM  Set to 1 to skip SHA256 verification

VERSION="${ETORO_VERSION:-latest}"
REPO="marianopa-tr/etoro-cli"
BASE_URL="https://github.com/${REPO}/releases"
CUSTOM_URL="${ETORO_INSTALL_URL:-}"
SKIP_CHECKSUM="${ETORO_SKIP_CHECKSUM:-0}"

# --- helpers ----------------------------------------------------------------

info()  { printf '  \033[1;34m>\033[0m %s\n' "$*"; }
ok()    { printf '  \033[1;32m✓\033[0m %s\n' "$*"; }
warn()  { printf '  \033[1;33m!\033[0m %s\n' "$*" >&2; }
fail()  { printf '  \033[1;31m✗\033[0m %s\n' "$*" >&2; exit 1; }

detect_os() {
  if [[ -n "${ETORO_INSTALL_OS:-}" ]]; then echo "$ETORO_INSTALL_OS"; return; fi
  case "$(uname -s)" in
    Darwin) echo "darwin" ;;
    Linux)  echo "linux"  ;;
    *)      fail "Unsupported OS: $(uname -s). Set ETORO_INSTALL_OS to override." ;;
  esac
}

detect_arch() {
  if [[ -n "${ETORO_INSTALL_ARCH:-}" ]]; then echo "$ETORO_INSTALL_ARCH"; return; fi
  case "$(uname -m)" in
    x86_64|amd64)   echo "amd64" ;;
    arm64|aarch64)  echo "arm64" ;;
    *)              fail "Unsupported arch: $(uname -m). Set ETORO_INSTALL_ARCH to override." ;;
  esac
}

default_install_dir() {
  if [[ -n "${ETORO_INSTALL_DIR:-}" ]]; then echo "$ETORO_INSTALL_DIR"; return; fi
  case "$(uname -s)" in
    Darwin) echo "/usr/local/bin" ;;
    *)      echo "$HOME/.local/bin" ;;
  esac
}

download() {
  local url="$1" out="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$out" 2>/dev/null && return 0
    return 1
  fi
  if command -v wget >/dev/null 2>&1; then
    wget -qO "$out" "$url" 2>/dev/null && return 0
    return 1
  fi
  fail "curl or wget is required."
}

sha256_verify() {
  local file="$1" expected="$2"
  local actual
  if command -v sha256sum >/dev/null 2>&1; then
    actual="$(sha256sum "$file" | awk '{print $1}')"
  elif command -v shasum >/dev/null 2>&1; then
    actual="$(shasum -a 256 "$file" | awk '{print $1}')"
  else
    warn "No sha256sum or shasum found; skipping checksum verification."
    return 0
  fi
  if [[ "$actual" != "$expected" ]]; then
    fail "Checksum mismatch!\n  Expected: ${expected}\n  Got:      ${actual}"
  fi
}

# --- main -------------------------------------------------------------------

os_name="$(detect_os)"
arch_name="$(detect_arch)"
install_dir="$(default_install_dir)"
asset="etoro_${os_name}_${arch_name}.tar.gz"

printf '\n  \033[1meToro CLI Installer\033[0m\n\n'
info "OS: ${os_name}  Arch: ${arch_name}  Version: ${VERSION}"

# Resolve download URL
if [[ -n "$CUSTOM_URL" ]]; then
  url="$CUSTOM_URL"
elif [[ "$VERSION" == "latest" ]]; then
  url="${BASE_URL}/latest/download/${asset}"
else
  url="${BASE_URL}/download/${VERSION}/${asset}"
fi

# Check for existing installation
if [[ -f "${install_dir}/etoro" ]]; then
  existing="$("${install_dir}/etoro" version 2>/dev/null || echo "unknown")"
  warn "Existing installation found: ${install_dir}/etoro (${existing})"
  info "It will be overwritten."
fi

# Download to temp dir
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

info "Downloading ${asset}..."
if ! download "$url" "$tmp_dir/$asset"; then
  echo ""
  fail "Download failed: ${url}
  Possible causes:
    - No release assets published yet (run: make release)
    - GitHub rate limiting (try again in a few minutes)
    - Invalid version tag: ${VERSION}
  You can also build from source:
    git clone https://github.com/${REPO}.git && cd etoro-cli && make install"
fi

# Checksum verification
if [[ "$SKIP_CHECKSUM" != "1" ]]; then
  checksum_found=0

  # Try GoReleaser-style checksums.txt first (one file, all hashes)
  if [[ "$VERSION" == "latest" ]]; then
    checksums_url="${BASE_URL}/latest/download/checksums.txt"
  else
    checksums_url="${BASE_URL}/download/${VERSION}/checksums.txt"
  fi
  if download "$checksums_url" "$tmp_dir/checksums.txt"; then
    expected="$(grep "${asset}" "$tmp_dir/checksums.txt" | awk '{print $1}')"
    if [[ -n "$expected" ]]; then
      info "Verifying SHA256 checksum..."
      sha256_verify "$tmp_dir/$asset" "$expected"
      ok "Checksum verified."
      checksum_found=1
    fi
  fi

  # Fallback: per-file .sha256 sidecar
  if [[ "$checksum_found" -eq 0 ]]; then
    sidecar_url="${url%.tar.gz}.sha256"
    if download "$sidecar_url" "$tmp_dir/checksum"; then
      expected="$(awk '{print $1}' "$tmp_dir/checksum")"
      info "Verifying SHA256 checksum..."
      sha256_verify "$tmp_dir/$asset" "$expected"
      ok "Checksum verified."
    else
      warn "No checksum file found; skipping verification."
    fi
  fi
fi

# Extract and install
info "Installing to ${install_dir}..."
tar -xzf "$tmp_dir/$asset" -C "$tmp_dir"
mkdir -p "$install_dir"
rm -f "$install_dir/etoro"
cp "$tmp_dir/etoro" "$install_dir/etoro"
chmod +x "$install_dir/etoro"

# Verify
if "${install_dir}/etoro" version >/dev/null 2>&1; then
  installed_version="$("${install_dir}/etoro" version 2>/dev/null)"
  ok "Installed: ${install_dir}/etoro (${installed_version})"
else
  ok "Installed: ${install_dir}/etoro"
fi

# PATH hint
case ":$PATH:" in
  *":${install_dir}:"*) ;;
  *)
    echo ""
    warn "${install_dir} is not in your PATH."
    info "Add it with:"
    info "  echo 'export PATH=\"${install_dir}:\$PATH\"' >> ~/.zshrc"
    ;;
esac

# Next steps
echo ""
info "Get started:"
info "  etoro setup        # configure API keys"
info "  etoro status       # check connectivity"
info "  etoro shell        # interactive REPL"
echo ""
