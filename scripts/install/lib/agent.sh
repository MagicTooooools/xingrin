#!/usr/bin/env bash

LF_AGENT_LAST_ENDPOINT=""
LF_AGENT_LAST_STATUS=""
LF_AGENT_LAST_MESSAGE=""
LF_AGENT_LAST_STAGE=""

lf_agent_health_endpoint() {
  local server_url="$1"
  printf '%s' "$server_url/health"
}

lf_agent_login_endpoint() {
  local server_url="$1"
  printf '%s' "$server_url/api/auth/login"
}

lf_agent_registration_token_endpoint() {
  local server_url="$1"
  printf '%s' "$server_url/api/agents/registration-tokens"
}

lf_agent_install_endpoint() {
  local server_url="$1"
  local token="$2"
  printf '%s' "$server_url/api/agents/install-script?token=$token"
}

lf_agent_extract_json_field() {
  local field="$1"
  grep -o "\"${field}\":\"[^\"]*\"" | head -n1 | cut -d: -f2 | tr -d '"'
}

lf_agent_wait_for_health() {
  local health_url="$1"
  local max_attempts="${2:-20}"
  local interval_seconds="${3:-2}"
  local curl_timeout="${4:-3}"
  local attempt=""

  LF_AGENT_LAST_ENDPOINT="$health_url"
  LF_AGENT_LAST_STATUS=""
  LF_AGENT_LAST_MESSAGE=""
  LF_AGENT_LAST_STAGE="health"

  for attempt in $(seq 1 "$max_attempts"); do
    if curl -ksf --max-time "$curl_timeout" "$health_url" >/dev/null 2>&1; then
      LF_AGENT_LAST_STATUS="ready"
      return 0
    fi
    sleep "$interval_seconds"
  done

  LF_AGENT_LAST_STATUS="timeout"
  LF_AGENT_LAST_MESSAGE="服务未就绪"
  return 1
}

lf_agent_login_access_token() {
  local server_url="$1"
  local username="$2"
  local password="$3"
  local response_file=""
  local access_token=""

  LF_AGENT_LAST_ENDPOINT="$(lf_agent_login_endpoint "$server_url")"
  LF_AGENT_LAST_STATUS=""
  LF_AGENT_LAST_MESSAGE=""
  LF_AGENT_LAST_STAGE="login"

  response_file="$(mktemp)"
  LF_AGENT_LAST_STATUS=$(curl -ksS -o "$response_file" -w "%{http_code}" -X POST "$LF_AGENT_LAST_ENDPOINT" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$username\",\"password\":\"$password\"}" || echo "000")

  cat "$response_file" >&2
  LF_AGENT_LAST_MESSAGE="$(lf_agent_extract_json_field message < "$response_file")"
  access_token="$(lf_agent_extract_json_field accessToken < "$response_file")"
  rm -f "$response_file"

  if [ "$LF_AGENT_LAST_STATUS" = "200" ] && [ -n "$access_token" ]; then
    printf '%s' "$access_token"
    return 0
  fi
  return 1
}

lf_agent_create_registration_token() {
  local server_url="$1"
  local access_token="$2"
  local response_file=""
  local registration_token=""

  LF_AGENT_LAST_ENDPOINT="$(lf_agent_registration_token_endpoint "$server_url")"
  LF_AGENT_LAST_STATUS=""
  LF_AGENT_LAST_MESSAGE=""
  LF_AGENT_LAST_STAGE="registration-token"

  response_file="$(mktemp)"
  LF_AGENT_LAST_STATUS=$(curl -ksS -o "$response_file" -w "%{http_code}" -X POST "$LF_AGENT_LAST_ENDPOINT" \
    -H "Authorization: Bearer $access_token" \
    -H "Content-Type: application/json" || echo "000")

  cat "$response_file" >&2
  LF_AGENT_LAST_MESSAGE="$(lf_agent_extract_json_field message < "$response_file")"
  registration_token="$(lf_agent_extract_json_field token < "$response_file")"
  rm -f "$response_file"

  if { [ "$LF_AGENT_LAST_STATUS" = "200" ] || [ "$LF_AGENT_LAST_STATUS" = "201" ]; } && [ -n "$registration_token" ]; then
    printf '%s' "$registration_token"
    return 0
  fi
  return 1
}

lf_agent_issue_registration_token() {
  local server_url="$1"
  local username="$2"
  local password="$3"
  local access_token=""
  local registration_token=""

  LF_AGENT_LAST_STAGE="login"
  if ! access_token="$(lf_agent_login_access_token "$server_url" "$username" "$password")"; then
    return 1
  fi

  LF_AGENT_LAST_STAGE="registration-token"
  if ! registration_token="$(lf_agent_create_registration_token "$server_url" "$access_token")"; then
    return 1
  fi

  LF_AGENT_LAST_STAGE="registration-token-issued"
  printf '%s' "$registration_token"
}

lf_agent_print_login_failure_hint() {
  echo -e "${DIM}请求地址: ${LF_AGENT_LAST_ENDPOINT}${RESET}"
  if [ -n "${LF_AGENT_LAST_MESSAGE}" ]; then
    echo -e "${DIM}错误信息: ${LF_AGENT_LAST_MESSAGE}${RESET}"
  fi
  echo -e "${DIM}请在前端「设置 → Workers」里手动生成安装命令，或确认 admin/admin 可用。${RESET}"
}

lf_agent_print_registration_token_failure_hint() {
  echo -e "${DIM}请求地址: ${LF_AGENT_LAST_ENDPOINT}${RESET}"
  if [ -n "${LF_AGENT_LAST_MESSAGE}" ]; then
    echo -e "${DIM}错误信息: ${LF_AGENT_LAST_MESSAGE}${RESET}"
  fi
  echo -e "${DIM}请确认服务已就绪，或在前端手动生成安装命令。${RESET}"
}

lf_agent_print_issue_registration_token_failure_hint() {
  case "${LF_AGENT_LAST_STAGE}" in
    login)
      lf_agent_print_login_failure_hint
      ;;
    registration-token)
      lf_agent_print_registration_token_failure_hint
      ;;
    *)
      echo -e "${DIM}请求地址: ${LF_AGENT_LAST_ENDPOINT}${RESET}"
      if [ -n "${LF_AGENT_LAST_MESSAGE}" ]; then
        echo -e "${DIM}错误信息: ${LF_AGENT_LAST_MESSAGE}${RESET}"
      fi
      ;;
  esac
}

lf_agent_print_install_failure_hint() {
  local install_script_url="$1"
  local register_url="$2"
  local agent_server_url="$3"
  local network_name="$4"

  echo -e "${DIM}请求地址: $install_script_url${RESET}"
  echo -e "${DIM}注册地址: LUNAFOX_AGENT_REGISTER_URL=$register_url${RESET}"
  echo -e "${DIM}Agent 连接地址: LUNAFOX_AGENT_SERVER_URL=$agent_server_url${RESET}"
  echo -e "${DIM}网络配置: LUNAFOX_AGENT_DOCKER_NETWORK=$network_name${RESET}"
  echo -e "${DIM}请检查服务端是否可达、Docker 是否可用，以及网络配置是否正确。${RESET}"
}

lf_agent_install_local_worker() {
  local mode="$1"
  local allow_pull="$2"
  local server_url="$3"
  local agent_server_url="$4"
  local agent_register_url="$5"
  local network_name="$6"
  local env_file="$7"
  local admin_user="${8:-admin}"
  local admin_pass="${9:-admin}"
  local health_url=""
  local reg_token=""
  local install_script_url=""
  local worker_token=""
  local -a install_env=()

  if ! command -v curl >/dev/null 2>&1; then
    error "未检测到 curl，无法自动注册本地 Agent"
    echo -e "${DIM}请先安装 curl 后重试。${RESET}"
    return 1
  fi

  info "检查服务是否可用..."
  health_url="$(lf_agent_health_endpoint "$server_url")"
  if ! lf_agent_wait_for_health "$health_url" 20 2 3; then
    error "服务未就绪，无法继续注册 Agent"
    echo -e "${DIM}请求地址: ${LF_AGENT_LAST_ENDPOINT}${RESET}"
    echo -e "${DIM}请检查服务日志: ${COMPOSE_CMD[*]} -f $COMPOSE_FILE logs${RESET}"
    return 1
  fi

  if lf_docker_network_enabled "$network_name"; then
    if ! lf_docker_ensure_network "$network_name"; then
      error "${LF_DOCKER_LAST_ERROR}"
      return 1
    fi
  else
    if echo "$agent_server_url" | grep -q "server:8080"; then
      warn "当前网络为 off，但 Agent 连接地址使用了 server:8080，可能无法解析"
    fi
  fi

  info "获取 Agent 注册令牌..."
  if ! reg_token="$(lf_agent_issue_registration_token "$server_url" "$admin_user" "$admin_pass")"; then
    if [ "${LF_AGENT_LAST_STAGE}" = "login" ]; then
      error "无法获取 accessToken，默认账号可能已修改"
    else
      error "无法获取注册令牌"
    fi
    lf_agent_print_issue_registration_token_failure_hint
    return 1
  fi

  install_script_url="$(lf_agent_install_endpoint "$server_url" "$reg_token")"
  if ! worker_token="$(lf_worker_token_from_env_file "$env_file")"; then
    error "${LF_ENV_LAST_ERROR}"
    return 1
  fi
  export WORKER_TOKEN="$worker_token"

  echo ""
  info "下载并执行 Agent 安装脚本..."
  if ! lf_agent_collect_install_env_array install_env "$mode" "$allow_pull" "$server_url" "$agent_server_url" "$network_name"; then
    error "${LF_ENV_LAST_ERROR}"
    return 1
  fi

  if ! lf_agent_execute_install_script "$install_script_url" "${install_env[@]}"; then
    error "Agent 安装脚本执行失败"
    lf_agent_print_install_failure_hint "$install_script_url" "$agent_register_url" "$agent_server_url" "$network_name"
    return 1
  fi

  success "Agent/Worker 安装完成"
}

lf_agent_execute_install_script() {
  local install_script_url="$1"
  shift
  (
    set -o pipefail
    curl -kfsSL "$install_script_url" | env "$@" bash
  )
}
