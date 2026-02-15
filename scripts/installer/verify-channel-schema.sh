#!/usr/bin/env bash
set -euo pipefail

EXPECTED_SCHEMA_VERSION="${EXPECTED_SCHEMA_VERSION:-1}"

error() {
  echo "✗ $*" >&2
}

read_env_value() {
  local file="$1"
  local key="$2"
  awk -F= -v key="$key" '$1 == key { print substr($0, index($0, "=") + 1); exit }' "$file" | tr -d '\r'
}

require_key() {
  local file="$1"
  local key="$2"
  local value
  value="$(read_env_value "$file" "$key")"
  if [ -z "$value" ]; then
    error "[$file] missing required key: $key"
    exit 1
  fi
}

validate_file() {
  local file="$1"
  if [ ! -f "$file" ]; then
    error "file not found: $file"
    exit 1
  fi

  require_key "$file" "SCHEMA_VERSION"
  local schema
  schema="$(read_env_value "$file" "SCHEMA_VERSION")"
  if [ "$schema" != "$EXPECTED_SCHEMA_VERSION" ]; then
    error "[$file] incompatible SCHEMA_VERSION: $schema (expected: $EXPECTED_SCHEMA_VERSION)"
    exit 1
  fi

  require_key "$file" "VERSION"
  require_key "$file" "GITHUB_BASE_URL"
  require_key "$file" "GITEE_BASE_URL"

  local platforms=(
    "LINUX_AMD64"
    "LINUX_ARM64"
    "DARWIN_AMD64"
    "DARWIN_ARM64"
  )
  local platform=""
  for platform in "${platforms[@]}"; do
    require_key "$file" "${platform}_ASSET"
    require_key "$file" "${platform}_SHA256"
  done

  require_key "$file" "IMAGE_REGISTRY"
  require_key "$file" "IMAGE_NAMESPACE"
  require_key "$file" "AGENT_IMAGE"
  require_key "$file" "WORKER_IMAGE"
}

if [ "$#" -lt 1 ]; then
  error "usage: $0 <channel-env> [channel-env...]"
  exit 1
fi

for file in "$@"; do
  validate_file "$file"
done

echo "✓ channel schema validated (${EXPECTED_SCHEMA_VERSION})"
