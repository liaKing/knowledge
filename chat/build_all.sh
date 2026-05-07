#!/usr/bin/env bash
set -euo pipefail

APP_NAME="litellm_chat"
ENTRY_FILE="litellm_chat.go"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="${SCRIPT_DIR}/dist"
ENTRY_PATH="${SCRIPT_DIR}/${ENTRY_FILE}"

# Override unavailable local proxy defaults (e.g. goproxy.cn).
: "${GOPROXY:=https://proxy.golang.org,direct}"
export GOPROXY

if [[ ! -f "${ENTRY_PATH}" ]]; then
  echo "Error: ${ENTRY_PATH} not found."
  exit 1
fi

mkdir -p "${DIST_DIR}"

build_target() {
  local goos="$1"
  local goarch="$2"
  local ext=""
  local output_name

  if [[ "${goos}" == "windows" ]]; then
    ext=".exe"
  fi

  output_name="${APP_NAME}-${goos}-${goarch}${ext}"
  echo "Building ${output_name} ..."

  # First attempt: use configured GOPROXY.
  if ! (
    cd "${SCRIPT_DIR}"
    CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" \
      go build -ldflags="-s -w" -o "${DIST_DIR}/${output_name}" "${ENTRY_FILE}"
  ); then
    echo "Primary build failed, retry with GOPROXY=direct ..."
    # Second attempt: bypass proxy.
    if ! (
      cd "${SCRIPT_DIR}"
      CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" GOPROXY=direct \
        go build -ldflags="-s -w" -o "${DIST_DIR}/${output_name}" "${ENTRY_FILE}"
    ); then
      echo "Direct mode failed, retry with GOSUMDB=off ..."
      # Last resort for restricted networks (less secure, but practical).
      (
        cd "${SCRIPT_DIR}"
        CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" GOPROXY=direct GOSUMDB=off \
          go build -ldflags="-s -w" -o "${DIST_DIR}/${output_name}" "${ENTRY_FILE}"
      )
    fi
  fi
}

build_target linux amd64
build_target linux arm64
build_target darwin amd64
build_target darwin arm64
build_target windows amd64
build_target windows arm64

echo
echo "Build completed. Artifacts are in ${DIST_DIR}:"
ls -lh "${DIST_DIR}"
