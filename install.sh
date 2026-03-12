#!/bin/sh
# Install lookit - downloads the latest release binary for your platform.
#
# Usage:
#   curl -sSL https://raw.githubusercontent.com/Benjamin-Connelly/lookit/master/install.sh | sh
#   curl -sSL https://raw.githubusercontent.com/Benjamin-Connelly/lookit/master/install.sh | sh -s -- --dir /usr/local/bin
#   curl -sSL https://raw.githubusercontent.com/Benjamin-Connelly/lookit/master/install.sh | sh -s -- --version v0.2.0

set -e

REPO="Benjamin-Connelly/lookit"
INSTALL_DIR="${HOME}/.local/bin"
VERSION=""

usage() {
    echo "Usage: install.sh [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --dir DIR        Install directory (default: ~/.local/bin)"
    echo "  --version VER    Install specific version (default: latest)"
    echo "  --help           Show this help"
    exit 0
}

while [ $# -gt 0 ]; do
    case "$1" in
        --dir)     INSTALL_DIR="$2"; shift 2 ;;
        --version) VERSION="$2"; shift 2 ;;
        --help)    usage ;;
        *)         echo "Unknown option: $1"; usage ;;
    esac
done

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    linux)  OS="linux" ;;
    darwin) OS="darwin" ;;
    *)      echo "Error: unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *)             echo "Error: unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest version if not specified
if [ -z "$VERSION" ]; then
    VERSION=$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep '"tag_name"' | head -1 | sed 's/.*"tag_name": *"\([^"]*\)".*/\1/')
    if [ -z "$VERSION" ]; then
        echo "Error: could not determine latest version"
        exit 1
    fi
fi

BINARY="lookit-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY}"
CHECKSUM_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"

echo "Installing lookit ${VERSION} (${OS}/${ARCH})..."

# Create install directory
mkdir -p "$INSTALL_DIR"

# Download binary
TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

echo "Downloading ${URL}..."
if ! curl -fsSL -o "${TMPDIR}/${BINARY}" "$URL"; then
    echo "Error: download failed. Check that ${VERSION} exists and has a ${OS}/${ARCH} binary."
    echo "Releases: https://github.com/${REPO}/releases"
    exit 1
fi

# Verify checksum
echo "Verifying checksum..."
if curl -fsSL -o "${TMPDIR}/checksums.txt" "$CHECKSUM_URL" 2>/dev/null; then
    expected=$(grep "${BINARY}$" "${TMPDIR}/checksums.txt" | awk '{print $1}')
    if [ -n "$expected" ]; then
        if command -v sha256sum >/dev/null 2>&1; then
            actual=$(sha256sum "${TMPDIR}/${BINARY}" | awk '{print $1}')
        elif command -v shasum >/dev/null 2>&1; then
            actual=$(shasum -a 256 "${TMPDIR}/${BINARY}" | awk '{print $1}')
        else
            echo "Warning: no sha256 tool found, skipping verification"
            actual="$expected"
        fi
        if [ "$expected" != "$actual" ]; then
            echo "Error: checksum mismatch"
            echo "  expected: $expected"
            echo "  got:      $actual"
            exit 1
        fi
        echo "Checksum verified."
    else
        echo "Warning: binary not found in checksums.txt, skipping verification"
    fi
else
    echo "Warning: could not download checksums, skipping verification"
fi

# Install
mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/lookit"
chmod +x "${INSTALL_DIR}/lookit"

echo ""
echo "Installed lookit ${VERSION} to ${INSTALL_DIR}/lookit"

# Check if install dir is in PATH
case ":${PATH}:" in
    *":${INSTALL_DIR}:"*) ;;
    *)
        echo ""
        echo "Add ${INSTALL_DIR} to your PATH:"
        echo "  echo 'export PATH=\"${INSTALL_DIR}:\$PATH\"' >> ~/.bashrc"
        echo ""
        ;;
esac

echo "Run 'lookit --help' to get started."
