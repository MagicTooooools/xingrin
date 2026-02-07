#!/usr/bin/env bash

LF_ENV_LAST_ERROR=""

lf_env_bool() {
  case "${1:-}" in
    1|true|TRUE|yes|YES|on|ON) return 0 ;;
    *) return 1 ;;
  esac
}

lf_agent_server_url_value() {
  local default_url="${1:-http://server:8080}"
  printf '%s' "${LUNAFOX_AGENT_SERVER_URL:-$default_url}"
}

lf_agent_register_url_value() {
  local server_url="$1"
  printf '%s' "${LUNAFOX_AGENT_REGISTER_URL:-$server_url}"
}

lf_agent_network_name_value() {
  local default_network="${1:-lunafox_network}"
  lf_trim "${LUNAFOX_AGENT_DOCKER_NETWORK:-$default_network}"
}

lf_agent_max_tasks_value() {
  printf '%s' "${LUNAFOX_AGENT_MAX_TASKS:-10}"
}

lf_agent_cpu_threshold_value() {
  printf '%s' "${LUNAFOX_AGENT_CPU_THRESHOLD:-80}"
}

lf_agent_mem_threshold_value() {
  printf '%s' "${LUNAFOX_AGENT_MEM_THRESHOLD:-80}"
}

lf_agent_disk_threshold_value() {
  printf '%s' "${LUNAFOX_AGENT_DISK_THRESHOLD:-85}"
}

lf_agent_should_skip_pull() {
  local mode="$1"
  local allow_pull="$2"
  [ "$mode" = "dev" ] && ! lf_env_bool "$allow_pull"
}

lf_agent_install_env_lines() {
  local mode="$1"
  local allow_pull="$2"
  local server_url="$3"
  local agent_server_url="$4"
  local network_name="$5"

  printf '%s\n' "LUNAFOX_AGENT_REGISTER_URL=$(lf_agent_register_url_value "$server_url")"
  printf '%s\n' "LUNAFOX_AGENT_SERVER_URL=$agent_server_url"
  printf '%s\n' "LUNAFOX_AGENT_DOCKER_NETWORK=$network_name"
  printf '%s\n' "LUNAFOX_AGENT_USE_LOCAL_LIMITS=1"
  printf '%s\n' "LUNAFOX_AGENT_MAX_TASKS=$(lf_agent_max_tasks_value)"
  printf '%s\n' "LUNAFOX_AGENT_CPU_THRESHOLD=$(lf_agent_cpu_threshold_value)"
  printf '%s\n' "LUNAFOX_AGENT_MEM_THRESHOLD=$(lf_agent_mem_threshold_value)"
  printf '%s\n' "LUNAFOX_AGENT_DISK_THRESHOLD=$(lf_agent_disk_threshold_value)"

  if lf_agent_should_skip_pull "$mode" "$allow_pull"; then
    printf '%s\n' "LUNAFOX_AGENT_SKIP_PULL=1"
  fi
}

lf_agent_collect_install_env_array() {
  local output_array_name="$1"
  local mode="$2"
  local allow_pull="$3"
  local server_url="$4"
  local agent_server_url="$5"
  local network_name="$6"
  local env_line=""
  local env_escaped=""

  LF_ENV_LAST_ERROR=""

  if [ -z "$output_array_name" ]; then
    LF_ENV_LAST_ERROR="未提供安装环境变量数组名称"
    return 1
  fi

  if [ "${BASH_VERSINFO[0]:-0}" -ge 4 ]; then
    local -n output_array_ref="$output_array_name"
    output_array_ref=()
    while IFS= read -r env_line; do
      if [ -n "$env_line" ]; then
        output_array_ref+=("$env_line")
      fi
    done < <(lf_agent_install_env_lines "$mode" "$allow_pull" "$server_url" "$agent_server_url" "$network_name")
    return 0
  fi

  eval "$output_array_name=()"
  while IFS= read -r env_line; do
    if [ -n "$env_line" ]; then
      env_escaped="$env_line"
      env_escaped="${env_escaped//\\/\\\\}"
      env_escaped="${env_escaped//\"/\\\"}"
      env_escaped="${env_escaped//\$/\\$}"
      eval "$output_array_name+=(\"$env_escaped\")"
    fi
  done < <(lf_agent_install_env_lines "$mode" "$allow_pull" "$server_url" "$agent_server_url" "$network_name")
}

lf_worker_token_from_env_file() {
  local env_file="$1"
  local worker_token=""

  LF_ENV_LAST_ERROR=""

  if [ ! -f "$env_file" ]; then
    LF_ENV_LAST_ERROR="未找到环境文件: $env_file"
    return 1
  fi

  worker_token="$(grep -E '^WORKER_TOKEN=' "$env_file" | head -n1 | cut -d'=' -f2- | tr -d '"' | tr -d "'" || true)"
  worker_token="$(lf_trim "$worker_token")"

  if [ -z "$worker_token" ]; then
    LF_ENV_LAST_ERROR="docker/.env 缺少 WORKER_TOKEN"
    return 1
  fi

  printf '%s' "$worker_token"
}
