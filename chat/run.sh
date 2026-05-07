#!/usr/bin/env bash
set -euo pipefail

APP_NAME="litellm_chat"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="${SCRIPT_DIR}/dist"

os="$(uname -s)"
arch="$(uname -m)"

case "${os}" in
  Darwin) goos="darwin" ;;
  Linux) goos="linux" ;;
  *)
    echo "Unsupported OS: ${os}"
    echo "Please run manually with a binary from ${DIST_DIR}/"
    exit 1
    ;;
esac

case "${arch}" in
  x86_64|amd64) goarch="amd64" ;;
  arm64|aarch64) goarch="arm64" ;;
  *)
    echo "Unsupported architecture: ${arch}"
    echo "Please run manually with a matching binary from ${DIST_DIR}/"
    exit 1
    ;;
esac

bin_path="${DIST_DIR}/${APP_NAME}-${goos}-${goarch}"

if [[ ! -x "${bin_path}" ]]; then
  echo "Binary not found or not executable: ${bin_path}"
  echo "Run ${SCRIPT_DIR}/build_all.sh first."
  exit 1
fi

echo "Running ${bin_path} ..."
exec "${bin_path}"
