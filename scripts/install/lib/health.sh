#!/usr/bin/env bash

lf_health_wait_for_services() {
  local health_url="$1"
  local frontend_url="$2"
  local max_attempts="${3:-120}"
  local interval_seconds="${4:-2}"
  local backend_timeout="${5:-2}"
  local frontend_timeout="${6:-6}"
  local compose_file="$7"
  local spin='⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
  local i=0
  local backend_ready=0
  local frontend_ready=0
  local backend_status="..."
  local frontend_status="..."
  local attempt=""

  if ! command -v curl >/dev/null 2>&1; then
    warn "未检测到 curl，跳过健康检查"
    return 0
  fi

  for attempt in $(seq 1 "$max_attempts"); do
    if [ "$backend_ready" -eq 0 ]; then
      if curl -ksf --max-time "$backend_timeout" "$health_url" >/dev/null 2>&1; then
        backend_ready=1
        backend_status="OK"
      fi
    fi

    if [ "$frontend_ready" -eq 0 ]; then
      if curl -ksf --max-time "$frontend_timeout" "$frontend_url" >/dev/null 2>&1; then
        frontend_ready=1
        frontend_status="OK"
      fi
    fi

    if [ "$backend_ready" -eq 1 ] && [ "$frontend_ready" -eq 1 ]; then
      printf "\r%50s\r"
      success "服务已就绪"
      return 0
    fi

    i=$(( (i + 1) % 10 ))
    printf "\r${CYAN}${spin:$i:1}${RESET}  等待服务启动... (%d/%d) [后端:%s 前端:%s]" "$attempt" "$max_attempts" "$backend_status" "$frontend_status"
    sleep "$interval_seconds"
  done

  echo ""
  error "服务启动超时，请检查日志: ${COMPOSE_CMD[*]} -f $compose_file logs"
  return 1
}

lf_health_prewarm_frontend() {
  local mode="$1"
  local base_url="${2:-https://localhost:8083}"
  local max_attempts="${3:-30}"
  local interval_seconds="${4:-2}"
  local timeout_seconds="${5:-6}"
  local -a paths=("/zh/login" "/zh/dashboard/")
  local path=""
  local url=""
  local warmed="0"
  local attempt=""
  local status=""

  if [ "$mode" != "dev" ]; then
    return 0
  fi

  if ! command -v curl >/dev/null 2>&1; then
    warn "未检测到 curl，跳过前端预热"
    return 0
  fi

  info "预热前端页面（开发模式）..."
  for path in "${paths[@]}"; do
    url="${base_url}${path}"
    warmed="0"

    for attempt in $(seq 1 "$max_attempts"); do
      status="$(curl -ksS -o /dev/null -w "%{http_code}" --max-time "$timeout_seconds" "$url" || echo "000")"
      if [ "$status" -ge 200 ] && [ "$status" -lt 500 ]; then
        success "预热完成: ${path}"
        warmed="1"
        break
      fi
      sleep "$interval_seconds"
    done

    if [ "$warmed" -eq 0 ]; then
      warn "预热未完成: ${path}（可稍后访问触发编译）"
    fi
  done
}
