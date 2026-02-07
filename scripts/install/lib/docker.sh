#!/usr/bin/env bash

LF_DOCKER_LAST_NETWORK=""
LF_DOCKER_LAST_ERROR=""

lf_docker_network_enabled() {
  local network_name="${1:-}"
  [ -n "$network_name" ] && [ "$network_name" != "off" ] && [ "$network_name" != "none" ]
}

lf_docker_ensure_network() {
  local network_name="${1:-}"

  LF_DOCKER_LAST_NETWORK="$network_name"
  LF_DOCKER_LAST_ERROR=""

  if ! lf_docker_network_enabled "$network_name"; then
    return 0
  fi

  if "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" network inspect "$network_name" >/dev/null 2>&1; then
    return 0
  fi

  if "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" network create "$network_name" >/dev/null 2>&1; then
    return 0
  fi

  LF_DOCKER_LAST_ERROR="无法创建 Docker 网络: $network_name"
  return 1
}
