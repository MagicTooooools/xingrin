#!/usr/bin/env bash
if [ -z "${BASH_VERSION:-}" ]; then
  SCRIPT_PATH="$0"
  case "$SCRIPT_PATH" in
    /*|*/*) ;;
    *) SCRIPT_PATH="./$SCRIPT_PATH" ;;
  esac
  exec /usr/bin/env bash "$SCRIPT_PATH" "$@"
fi
set -e

# ==============================================================================
# LunaFox 安装脚本
#   ./install.sh        # 默认生产模式
#   ./install.sh --dev  # 开发模式
# ==============================================================================

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DOCKER_DIR="$ROOT_DIR/docker"
ENV_FILE="$DOCKER_DIR/.env"
COMPOSE_DEV="$DOCKER_DIR/docker-compose.dev.yml"
COMPOSE_PROD="$DOCKER_DIR/docker-compose.yml"
DATA_DIR="/opt/lunafox"

MODE="prod"
VERSION_TAG="unknown"
COMPOSE_FILE=""
COMPOSE_CMD=()
COMPOSE_ENV=()
DOCKER_PREFIX=()
DOCKER_BIN=""
TOTAL_STEPS=8
GO111MODULE_VALUE="on"
GOPROXY_VALUE="https://proxy.golang.org,direct"
ALLOW_PULL="false"

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
DIM='\033[2m'
RESET='\033[0m'

# 渐变色定义
GRADIENT_1='\033[38;5;39m'
GRADIENT_2='\033[38;5;45m'
GRADIENT_3='\033[38;5;51m'
GRADIENT_4='\033[38;5;87m'
GRADIENT_5='\033[38;5;123m'

# 特殊效果
UNDERLINE='\033[4m'

usage() {
  cat <<'EOF'
用法:
  ./install.sh        # 默认生产模式
  ./install.sh --dev  # 开发模式
  ./install.sh --goproxy                  # 使用 https://goproxy.cn,direct
  ./install.sh --dev --allow-pull         # 开发模式允许从仓库拉镜像
EOF
}

info() {
  echo -e "${CYAN}ℹ${RESET}  $*"
}

warn() {
  echo -e "${YELLOW}⚠${RESET}  $*" >&2
}

error() {
  echo -e "${RED}✗${RESET}  $*" >&2
}

success() {
  echo -e "${GREEN}✓${RESET}  $*"
}

step_header() {
  local current=$1
  local total=$2
  local title=$3
  echo ""
  echo -e "${BOLD}${CYAN}[$current/$total]${RESET} ${BOLD}${title}${RESET}"
}

success_animation() {
  local message=$1
  echo ""
  echo -e "${GREEN}${BOLD}✓ ${message}${RESET}"
}

banner() {
  clear
  echo ""
  echo -e "  ${GRADIENT_1}██${GRADIENT_2}╗     ${GRADIENT_2}██${GRADIENT_2}╗   ${GRADIENT_2}██${GRADIENT_2}╗${GRADIENT_3}███${GRADIENT_3}╗   ${GRADIENT_3}██${GRADIENT_3}╗${GRADIENT_3}█████${GRADIENT_4}╗ ${GRADIENT_4}███████${GRADIENT_5}╗ ${GRADIENT_5}██████${GRADIENT_5}╗ ${GRADIENT_5}██${GRADIENT_5}╗  ${GRADIENT_5}██${GRADIENT_5}╗${RESET}"
  echo -e "  ${GRADIENT_1}██${GRADIENT_2}║     ${GRADIENT_2}██${GRADIENT_2}║   ${GRADIENT_2}██${GRADIENT_2}║${GRADIENT_3}████${GRADIENT_3}╗  ${GRADIENT_3}██${GRADIENT_3}║${GRADIENT_3}██${GRADIENT_4}╔══${GRADIENT_4}██${GRADIENT_4}╗${GRADIENT_4}██${GRADIENT_4}╔════╝${GRADIENT_4}██${GRADIENT_5}╔═══${GRADIENT_5}██${GRADIENT_5}╗${GRADIENT_5}╚${GRADIENT_5}██${GRADIENT_5}╗${GRADIENT_5}██${GRADIENT_5}╔╝${RESET}"
  echo -e "  ${GRADIENT_2}██${GRADIENT_2}║     ${GRADIENT_2}██${GRADIENT_2}║   ${GRADIENT_2}██${GRADIENT_2}║${GRADIENT_3}██${GRADIENT_3}╔${GRADIENT_3}██${GRADIENT_3}╗ ${GRADIENT_3}██${GRADIENT_3}║${GRADIENT_4}███████${GRADIENT_4}║${GRADIENT_4}█████${GRADIENT_5}╗  ${GRADIENT_5}██${GRADIENT_5}║   ${GRADIENT_5}██${GRADIENT_5}║ ${GRADIENT_5}╚███${GRADIENT_5}╔╝${RESET}"
  echo -e "  ${GRADIENT_2}██${GRADIENT_2}║     ${GRADIENT_2}██${GRADIENT_2}║   ${GRADIENT_2}██${GRADIENT_2}║${GRADIENT_3}██${GRADIENT_3}║${GRADIENT_3}╚${GRADIENT_3}██${GRADIENT_3}╗${GRADIENT_3}██${GRADIENT_3}║${GRADIENT_4}██${GRADIENT_4}╔══${GRADIENT_4}██${GRADIENT_4}║${GRADIENT_5}██${GRADIENT_5}╔══╝  ${GRADIENT_5}██${GRADIENT_5}║   ${GRADIENT_5}██${GRADIENT_5}║ ${GRADIENT_5}██${GRADIENT_5}╔${GRADIENT_5}██${GRADIENT_5}╗${RESET}"
  echo -e "  ${GRADIENT_3}███████${GRADIENT_3}╗${GRADIENT_3}╚${GRADIENT_3}██████${GRADIENT_4}╔╝${GRADIENT_4}██${GRADIENT_4}║ ${GRADIENT_4}╚████${GRADIENT_4}║${GRADIENT_5}██${GRADIENT_5}║  ${GRADIENT_5}██${GRADIENT_5}║${GRADIENT_5}██${GRADIENT_5}║     ${GRADIENT_5}╚${GRADIENT_5}██████${GRADIENT_5}╔╝${GRADIENT_5}██${GRADIENT_5}╔╝ ${GRADIENT_5}██${GRADIENT_5}╗${RESET}"
  echo -e "  ${GRADIENT_3}╚══════╝${GRADIENT_4} ╚═════╝ ${GRADIENT_4}╚═╝  ╚═══╝${GRADIENT_5}╚═╝  ╚═╝${GRADIENT_5}╚═╝      ╚═════╝ ╚═╝  ╚═╝${RESET}"
  echo ""
  echo -e "  ${BOLD}${CYAN}🦊 LunaFox${RESET} ${DIM}·${RESET} ${BOLD}开源安全扫描平台${RESET}"
  echo -e "  ${DIM}版本:${RESET} ${YELLOW}${VERSION_TAG}${RESET}  ${DIM}模式:${RESET} ${MAGENTA}${MODE}${RESET}"
  echo ""
}

extract_json_field() {
  local field="$1"
  grep -o "\"${field}\":\"[^\"]*\"" | head -n1 | cut -d: -f2 | tr -d '"'
}

compose_plugin_path() {
  local paths=()
  local user_home=""
  if [ -n "${SUDO_USER:-}" ]; then
    user_home="$(eval echo "~$SUDO_USER")"
  fi
  if [ -d "$HOME/.docker/cli-plugins" ]; then
    paths+=("$HOME/.docker/cli-plugins")
  fi
  if [ -n "$user_home" ] && [ -d "$user_home/.docker/cli-plugins" ]; then
    paths+=("$user_home/.docker/cli-plugins")
  fi
  for p in /usr/local/lib/docker/cli-plugins /usr/libexec/docker/cli-plugins /usr/lib/docker/cli-plugins; do
    if [ -d "$p" ]; then
      paths+=("$p")
    fi
  done
  local joined=""
  for p in "${paths[@]}"; do
    joined="${joined:+$joined:}$p"
  done
  echo "$joined"
}

generate_secret() {
  "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" run --rm alpine/openssl rand -hex 32
}

require_secret_tool() {
  if command -v docker; then
    return 0
  fi
  error "未检测到 docker，无法生成密钥"
  exit 1
}

ensure_docker() {
  if ! command -v docker; then
    error "未检测到 docker 命令，请先安装 Docker。"
    exit 1
  fi
  DOCKER_BIN="$(command -v docker)"
  if "$DOCKER_BIN" info; then
    DOCKER_PREFIX=()
    return 0
  fi
  if command -v sudo; then
    if sudo "$DOCKER_BIN" info; then
      DOCKER_PREFIX=(sudo)
      return 0
    fi
    error "当前环境无法使用 sudo 访问 Docker（可能被禁用或需要权限）。"
  else
    error "未检测到 sudo，无法提升权限访问 Docker。"
  fi
  error "Docker 守护进程未运行或无权限访问。请确认 Docker 已启动且当前用户有权限访问 Docker socket。"
  exit 1
}

detect_compose() {
  if [ -z "$DOCKER_BIN" ]; then
    DOCKER_BIN="$(command -v docker || true)"
  fi
  if [ -z "$DOCKER_BIN" ]; then
    error "未检测到 docker 可执行文件"
    exit 1
  fi
  local plugin_path
  plugin_path="$(compose_plugin_path)"
  if [ -n "$plugin_path" ]; then
    COMPOSE_ENV=(env DOCKER_CLI_PLUGIN_PATH="$plugin_path")
  else
    COMPOSE_ENV=()
  fi
  if command -v docker-compose; then
    COMPOSE_CMD=("${DOCKER_PREFIX[@]}" "$(command -v docker-compose)")
    return 0
  fi
  for p in ${plugin_path//:/ }; do
    if [ -x "$p/docker-compose" ]; then
      COMPOSE_CMD=("${DOCKER_PREFIX[@]}" "${COMPOSE_ENV[@]}" "$DOCKER_BIN" compose)
      return 0
    fi
  done
  error "未检测到 docker compose，请先安装。"
  exit 1
}

check_system() {
  info "系统环境校验..."
  ensure_docker
  detect_compose
  info "使用 docker 前缀: ${DOCKER_PREFIX[*]:-无}"
  info "使用 compose 命令: ${COMPOSE_CMD[*]}"
  require_secret_tool
  if [ "$MODE" != "dev" ] && [ ! -f "$ROOT_DIR/VERSION" ]; then
    error "未找到 VERSION 文件"
    exit 1
  fi
  if [ ! -f "$COMPOSE_FILE" ]; then
    error "未找到 compose 文件: $COMPOSE_FILE"
    exit 1
  fi
  success "环境校验通过"
}

init_data_dir() {
  if [ ! -d "$DATA_DIR" ]; then
    info "创建数据目录: $DATA_DIR"
    if ! mkdir -p "$DATA_DIR"/{results,logs,wordlists,workspace}; then
      error "无法创建 $DATA_DIR，请使用 sudo 运行此脚本。"
      exit 1
    fi
  fi
  chmod -R 777 "$DATA_DIR" || true
}

generate_ssl_cert() {
  local ssl_dir="$DOCKER_DIR/nginx/ssl"
  local fullchain="$ssl_dir/fullchain.pem"
  local privkey="$ssl_dir/privkey.pem"

  if [ -f "$fullchain" ] && [ -f "$privkey" ]; then
    info "检测到已有 HTTPS 证书，跳过生成。"
    return 0
  fi

  info "生成自签 HTTPS 证书（localhost）..."
  mkdir -p "$ssl_dir"

  "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" run --rm -v "$ssl_dir:/ssl" alpine/openssl \
    req -x509 -nodes -newkey rsa:2048 -days 365 \
    -keyout /ssl/privkey.pem \
    -out /ssl/fullchain.pem \
    -subj "/C=CN/ST=NA/L=NA/O=LunaFox/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

  if [ ! -f "$fullchain" ] || [ ! -f "$privkey" ]; then
    error "证书生成失败，请手动放置证书到 $ssl_dir"
    exit 1
  fi
}

write_env_file() {
  info "生成 JWT 密钥..."
  JWT_SECRET="$(generate_secret)"
  info "生成 Worker 令牌..."
  WORKER_TOKEN="$(generate_secret)"
  IMAGE_TAG="$VERSION_TAG"

  cat > "$ENV_FILE" <<EOF
IMAGE_TAG=$IMAGE_TAG
JWT_SECRET=$JWT_SECRET
WORKER_TOKEN=$WORKER_TOKEN
DB_HOST=postgres
DB_PASSWORD=postgres
REDIS_HOST=redis
DB_USER=postgres
DB_NAME=lunafox
DB_PORT=5432
REDIS_PORT=6379
GO111MODULE=$GO111MODULE_VALUE
GOPROXY=$GOPROXY_VALUE
EOF
}

build_dev_images() {
  if [ "$MODE" != "dev" ]; then
    info "生产模式跳过本地构建"
    return
  fi

  export DOCKER_BUILDKIT=1
  export COMPOSE_DOCKER_CLI_BUILD=1
  export BUILDKIT_PROGRESS=plain
  info "已启用 BuildKit 加速构建"

  CACHE_BASE="${HOME}/.cache/lunafox-buildx"
  AGENT_CACHE="${CACHE_BASE}/agent"
  WORKER_CACHE="${CACHE_BASE}/worker"
  mkdir -p "$AGENT_CACHE" "$WORKER_CACHE"
  info "Buildx 缓存目录: $CACHE_BASE"

  info "检测 buildx driver..."
  BUILDX_INSPECT_FILE="$(mktemp)"
  "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" buildx inspect | tee "$BUILDX_INSPECT_FILE"
  BUILDX_DRIVER="$(awk -F': *' '/^Driver:/ {print $2; exit}' "$BUILDX_INSPECT_FILE")"
  rm -f "$BUILDX_INSPECT_FILE"

  if [ -z "$BUILDX_DRIVER" ]; then
    warn "未检测到 buildx driver，默认禁用本地缓存导出"
    ENABLE_CACHE_EXPORT="false"
  elif [ "$BUILDX_DRIVER" = "docker" ]; then
    warn "检测到 buildx driver=docker，不支持 cache export，已禁用本地缓存导出"
    ENABLE_CACHE_EXPORT="false"
  else
    info "buildx driver: $BUILDX_DRIVER（启用本地缓存导出）"
    ENABLE_CACHE_EXPORT="true"
  fi

  info "并行构建 agent/worker 镜像 (buildx bake)..."
  BAKE_FILE="$(mktemp)"
  {
    cat <<EOF
group "default" {
  targets = ["agent", "worker"]
}

target "agent" {
  context = "$ROOT_DIR"
  dockerfile = "$ROOT_DIR/agent/Dockerfile"
  tags = ["yyhuni/lunafox-agent:dev"]
  build-args = {
    BUILDKIT_INLINE_CACHE = "1"
    GO111MODULE = "$GO111MODULE_VALUE"
    GOPROXY = "$GOPROXY_VALUE"
  }
EOF
    if [ "$ENABLE_CACHE_EXPORT" = "true" ]; then
      cat <<EOF
  cache-from = ["type=local,src=$AGENT_CACHE"]
  cache-to = ["type=local,dest=$AGENT_CACHE,mode=max"]
EOF
    fi
    cat <<EOF
}

target "worker" {
  context = "$ROOT_DIR/worker"
  dockerfile = "$ROOT_DIR/worker/Dockerfile"
  tags = ["yyhuni/lunafox-worker:dev"]
  build-args = {
    BUILDKIT_INLINE_CACHE = "1"
    GO111MODULE = "$GO111MODULE_VALUE"
    GOPROXY = "$GOPROXY_VALUE"
  }
EOF
    if [ "$ENABLE_CACHE_EXPORT" = "true" ]; then
      cat <<EOF
  cache-from = ["type=local,src=$WORKER_CACHE"]
  cache-to = ["type=local,dest=$WORKER_CACHE,mode=max"]
EOF
    fi
    cat <<EOF
}
EOF
  } > "$BAKE_FILE"

  "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" buildx bake -f "$BAKE_FILE" --progress=plain --load
  rm -f "$BAKE_FILE"
  success "Agent/Worker 镜像已构建"
}

start_services() {
  ENV_FILE_ARGS=()
  if [ -f "$ENV_FILE" ]; then
    ENV_FILE_ARGS=(--env-file "$ENV_FILE")
  fi
  PROFILE_ARG=""
  if [ "$MODE" = "prod" ]; then
    PROFILE_ARG="--profile local-db"
  fi

  info "正在启动服务（自动重建所有镜像）..."
  echo ""
  if ! "${COMPOSE_CMD[@]}" "${ENV_FILE_ARGS[@]}" -f "$COMPOSE_FILE" ${PROFILE_ARG:+$PROFILE_ARG} up -d --build --force-recreate; then
    error "服务启动失败"
    exit 1
  fi
  success "Docker 服务启动成功"
}

wait_for_health() {
  step_header 7 $TOTAL_STEPS "等待服务就绪"
  HEALTH_URL="https://localhost:8083/health"
  FRONTEND_URL="https://localhost:8083/"
  MAX_ATTEMPTS=120
  INTERVAL=2
  BACKEND_TIMEOUT=2
  FRONTEND_TIMEOUT=6

  if ! command -v curl >/dev/null 2>&1; then
    warn "未检测到 curl，跳过健康检查"
    return
  fi

  local spin='⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
  local i=0
  local backend_ready=0
  local frontend_ready=0
  local backend_status="..."
  local frontend_status="..."

  for attempt in $(seq 1 $MAX_ATTEMPTS); do
    if [ "$backend_ready" -eq 0 ]; then
      if curl -ksf --max-time "$BACKEND_TIMEOUT" "$HEALTH_URL" >/dev/null 2>&1; then
        backend_ready=1
        backend_status="OK"
      fi
    fi

    if [ "$frontend_ready" -eq 0 ]; then
      if curl -ksf --max-time "$FRONTEND_TIMEOUT" "$FRONTEND_URL" >/dev/null 2>&1; then
        frontend_ready=1
        frontend_status="OK"
      fi
    fi

    if [ "$backend_ready" -eq 1 ] && [ "$frontend_ready" -eq 1 ]; then
      printf "\r%50s\r"
      success "服务已就绪"
      return
    fi
    i=$(( (i+1) % 10 ))
    printf "\r${CYAN}${spin:$i:1}${RESET}  等待服务启动... (%d/%d) [后端:%s 前端:%s]" "$attempt" "$MAX_ATTEMPTS" "$backend_status" "$frontend_status"
    sleep $INTERVAL
  done

  echo ""
  error "服务启动超时，请检查日志: ${COMPOSE_CMD[*]} -f $COMPOSE_FILE logs"
  exit 1
}

prewarm_frontend() {
  if [ "$MODE" != "dev" ]; then
    return
  fi
  if ! command -v curl >/dev/null 2>&1; then
    warn "未检测到 curl，跳过前端预热"
    return
  fi

  local base_url="https://localhost:8083"
  local paths=("/zh/login" "/zh/dashboard/")
  local max_attempts=30
  local interval=2
  local timeout=6

  info "预热前端页面（开发模式）..."
  for path in "${paths[@]}"; do
    local url="${base_url}${path}"
    local warmed=0
    for attempt in $(seq 1 "$max_attempts"); do
      local status
      status=$(curl -ksS -o /dev/null -w "%{http_code}" --max-time "$timeout" "$url" || echo "000")
      if [ "$status" -ge 200 ] && [ "$status" -lt 500 ]; then
        success "预热完成: ${path}"
        warmed=1
        break
      fi
      sleep "$interval"
    done
    if [ "$warmed" -eq 0 ]; then
      warn "预热未完成: ${path}（可稍后访问触发编译）"
    fi
  done
}

print_summary() {
  echo ""
  success_animation "安装完成!"

  echo ""
  echo -e "${BOLD}${CYAN}🌐 访问地址:${RESET} ${UNDERLINE}https://localhost:8083/${RESET}"
  echo -e "${BOLD}${MAGENTA}🧰 默认账号:${RESET} ${BOLD}admin${RESET}"
  echo -e "${BOLD}${MAGENTA}🔐 默认密码:${RESET} ${BOLD}admin${RESET}"
  echo ""
  echo -e "${BOLD}${YELLOW}🐳 镜像说明:${RESET}"
  echo -e "   ${DIM}服务镜像:${RESET} ${GREEN}yyhuni/lunafox-server:${IMAGE_TAG}${RESET}"
  echo -e "   ${DIM}前端镜像:${RESET} ${GREEN}yyhuni/lunafox-frontend:${IMAGE_TAG}${RESET}"
  echo -e "   ${DIM}网关镜像:${RESET} ${GREEN}yyhuni/lunafox-nginx:${IMAGE_TAG}${RESET}"
  echo -e "   ${DIM}Agent 镜像:${RESET} ${GREEN}yyhuni/lunafox-agent:${IMAGE_TAG}${RESET}"
  echo -e "   ${DIM}Worker 镜像:${RESET} ${GREEN}yyhuni/lunafox-worker:${IMAGE_TAG}${RESET}"
  echo ""
  echo -e "${BOLD}${YELLOW}📋 常用命令:${RESET}"
  echo -e "   ${DIM}停止服务:${RESET} ./stop.sh"
  echo -e "   ${DIM}重启服务:${RESET} ./restart.sh"
  echo -e "   ${DIM}启动服务:${RESET} ./start.sh"
  echo ""
  echo -e "${BOLD}${YELLOW}📦 日志/镜像:${RESET}"
  echo -e "   ${DIM}查看日志(全部):${RESET} ${COMPOSE_CMD[*]} -f ${COMPOSE_FILE} logs -f"
  echo -e "   ${DIM}查看日志(后端):${RESET} ${COMPOSE_CMD[*]} -f ${COMPOSE_FILE} logs -f server"
  echo -e "   ${DIM}查看日志(前端):${RESET} ${COMPOSE_CMD[*]} -f ${COMPOSE_FILE} logs -f frontend"
  echo -e "   ${DIM}查看日志(网关):${RESET} ${COMPOSE_CMD[*]} -f ${COMPOSE_FILE} logs -f nginx"
  echo -e "   ${DIM}查看日志(数据库):${RESET} ${COMPOSE_CMD[*]} -f ${COMPOSE_FILE} logs -f postgres"
  echo -e "   ${DIM}查看日志(Redis):${RESET} ${COMPOSE_CMD[*]} -f ${COMPOSE_FILE} logs -f redis"
  echo -e "   ${DIM}查看日志(Agent):${RESET} docker logs -f lunafox-agent"
  echo -e "   ${DIM}查看镜像:${RESET} docker images | grep lunafox"
  echo ""
}

install_agent_worker() {
  step_header 8 $TOTAL_STEPS "安装本地 Agent/Worker"
  SERVER_URL="${INSTALL_SERVER_URL:-https://localhost:8083}"
  AGENT_SERVER_URL="${LUNAFOX_AGENT_SERVER_URL:-http://server:8080}"
  REGISTER_URL_ENV=("LUNAFOX_REGISTER_URL=${LUNAFOX_REGISTER_URL:-$SERVER_URL}")
  AGENT_URL_ENV=("LUNAFOX_AGENT_SERVER_URL=$AGENT_SERVER_URL")
  NETWORK_NAME="${LUNAFOX_NETWORK:-lunafox_network}"
  NETWORK_ENV=("LUNAFOX_NETWORK=$NETWORK_NAME")
  THRESHOLD_ENV=(
    "MAX_TASKS=${MAX_TASKS:-10}"
    "CPU_THRESHOLD=${CPU_THRESHOLD:-80}"
    "MEM_THRESHOLD=${MEM_THRESHOLD:-80}"
    "DISK_THRESHOLD=${DISK_THRESHOLD:-85}"
  )
  ADMIN_USER="admin"
  ADMIN_PASS="admin"

  if ! command -v curl >/dev/null 2>&1; then
    error "未检测到 curl，无法自动注册本地 Agent"
    echo -e "${DIM}请先安装 curl 后重试。${RESET}"
    exit 1
  fi

  info "检查服务是否可用..."
  HEALTH_URL="$SERVER_URL/health"
  local health_ready=0
  for attempt in $(seq 1 20); do
    if curl -ksf --max-time 3 "$HEALTH_URL" >/dev/null 2>&1; then
      health_ready=1
      break
    fi
    sleep 2
  done
  if [ "$health_ready" -ne 1 ]; then
    error "服务未就绪，无法继续注册 Agent"
    echo -e "${DIM}请求地址: $HEALTH_URL${RESET}"
    echo -e "${DIM}请检查服务日志: ${COMPOSE_CMD[*]} -f $COMPOSE_FILE logs${RESET}"
    exit 1
  fi

  if [ -n "$NETWORK_NAME" ] && [ "$NETWORK_NAME" != "off" ] && [ "$NETWORK_NAME" != "none" ]; then
    if ! "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" network inspect "$NETWORK_NAME" >/dev/null 2>&1; then
      info "创建 Docker 网络: $NETWORK_NAME"
      if ! "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" network create "$NETWORK_NAME" >/dev/null 2>&1; then
        error "无法创建 Docker 网络: $NETWORK_NAME"
        exit 1
      fi
    fi
  else
    if echo "$AGENT_SERVER_URL" | grep -q "server:8080"; then
      warn "当前网络为 off，但 Agent 连接地址使用了 server:8080，可能无法解析"
    fi
  fi

  info "获取访问令牌..."
  LOGIN_BODY="$(mktemp)"
  LOGIN_STATUS=$(curl -ksS -o "$LOGIN_BODY" -w "%{http_code}" -X POST "$SERVER_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$ADMIN_USER\",\"password\":\"$ADMIN_PASS\"}" || echo "000")
  echo "$(cat "$LOGIN_BODY")"
  ACCESS_TOKEN=$(extract_json_field accessToken < "$LOGIN_BODY")
  if [ "$LOGIN_STATUS" != "200" ] || [ -z "$ACCESS_TOKEN" ]; then
    error "无法获取 accessToken，默认账号可能已修改"
    echo -e "${DIM}请求地址: $SERVER_URL/api/auth/login${RESET}"
    LOGIN_MSG=$(extract_json_field message < "$LOGIN_BODY")
    if [ -n "$LOGIN_MSG" ]; then
      echo -e "${DIM}错误信息: $LOGIN_MSG${RESET}"
    fi
    echo -e "${DIM}请在前端「设置 → Workers」里手动生成安装命令，或确认 admin/admin 可用。${RESET}"
    rm -f "$LOGIN_BODY"
    exit 1
  fi
  rm -f "$LOGIN_BODY"

  echo ""
  info "创建 Agent 注册令牌..."
  TOKEN_BODY="$(mktemp)"
  TOKEN_STATUS=$(curl -ksS -o "$TOKEN_BODY" -w "%{http_code}" -X POST "$SERVER_URL/api/registration-tokens" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -H "Content-Type: application/json" || echo "000")
  echo "$(cat "$TOKEN_BODY")"
  REG_TOKEN=$(extract_json_field token < "$TOKEN_BODY")
  if { [ "$TOKEN_STATUS" != "200" ] && [ "$TOKEN_STATUS" != "201" ]; } || [ -z "$REG_TOKEN" ]; then
    error "无法获取注册令牌"
    echo -e "${DIM}请求地址: $SERVER_URL/api/registration-tokens${RESET}"
    TOKEN_MSG=$(extract_json_field message < "$TOKEN_BODY")
    if [ -n "$TOKEN_MSG" ]; then
      echo -e "${DIM}错误信息: $TOKEN_MSG${RESET}"
    fi
    echo -e "${DIM}请确认服务已就绪，或在前端手动生成安装命令。${RESET}"
    rm -f "$TOKEN_BODY"
    exit 1
  fi
  rm -f "$TOKEN_BODY"

  WORKER_TOKEN="$(grep -E "^WORKER_TOKEN=" "$ENV_FILE" | cut -d'=' -f2 | tr -d '\"' || true)"
  if [ -z "$WORKER_TOKEN" ]; then
    error "docker/.env 缺少 WORKER_TOKEN"
    exit 1
  fi
  export WORKER_TOKEN

  echo ""
  info "下载并执行 Agent 安装脚本..."
  PULL_ENV=()
  if [ "$MODE" = "dev" ] && [ "$ALLOW_PULL" != "true" ]; then
    PULL_ENV=("AGENT_SKIP_PULL=1")
  fi
  if [ ${#PULL_ENV[@]} -gt 0 ]; then
    if ! (set -o pipefail; curl -kfsSL "$SERVER_URL/api/agents/install.sh?token=$REG_TOKEN" | env "${REGISTER_URL_ENV[@]}" "${AGENT_URL_ENV[@]}" "${NETWORK_ENV[@]}" "${THRESHOLD_ENV[@]}" "${PULL_ENV[@]}" bash); then
      error "Agent 安装脚本执行失败"
      echo -e "${DIM}请求地址: $SERVER_URL/api/agents/install.sh?token=$REG_TOKEN${RESET}"
      echo -e "${DIM}注册地址: ${REGISTER_URL_ENV[*]}${RESET}"
      echo -e "${DIM}Agent 连接地址: ${AGENT_URL_ENV[*]}${RESET}"
      echo -e "${DIM}网络配置: ${NETWORK_ENV[*]}${RESET}"
      echo -e "${DIM}请检查服务端是否可达、Docker 是否可用，以及网络配置是否正确。${RESET}"
      exit 1
    fi
  else
    if ! (set -o pipefail; curl -kfsSL "$SERVER_URL/api/agents/install.sh?token=$REG_TOKEN" | env "${REGISTER_URL_ENV[@]}" "${AGENT_URL_ENV[@]}" "${NETWORK_ENV[@]}" "${THRESHOLD_ENV[@]}" bash); then
      error "Agent 安装脚本执行失败"
      echo -e "${DIM}请求地址: $SERVER_URL/api/agents/install.sh?token=$REG_TOKEN${RESET}"
      echo -e "${DIM}注册地址: ${REGISTER_URL_ENV[*]}${RESET}"
      echo -e "${DIM}Agent 连接地址: ${AGENT_URL_ENV[*]}${RESET}"
      echo -e "${DIM}网络配置: ${NETWORK_ENV[*]}${RESET}"
      echo -e "${DIM}请检查服务端是否可达、Docker 是否可用，以及网络配置是否正确。${RESET}"
      exit 1
    fi
  fi
  success "Agent/Worker 安装完成"
}

parse_args() {
  for arg in "$@"; do
    case "$arg" in
      --dev) MODE="dev" ;;
      --goproxy) GOPROXY_VALUE="https://goproxy.cn,direct" ;;
      --allow-pull) ALLOW_PULL="true" ;;
      -h|--help) usage; exit 0 ;;
      *) ;;
    esac
  done
}

set_compose_file() {
  if [ "$MODE" = "dev" ]; then
    COMPOSE_FILE="$COMPOSE_DEV"
  else
    COMPOSE_FILE="$COMPOSE_PROD"
  fi
}

ensure_project_structure() {
  if [ ! -d "$DOCKER_DIR" ] || [ ! -f "$COMPOSE_FILE" ]; then
    error "未找到 docker 目录或 compose 文件，请确认项目结构。"
    exit 1
  fi
}

set_version_tag() {
  if [ "$MODE" = "dev" ]; then
    VERSION_TAG="dev"
    return
  fi
  if [ -f "$ROOT_DIR/VERSION" ]; then
    VERSION_TAG="$(tr -d '[:space:]' < "$ROOT_DIR/VERSION")"
  else
    error "未找到 VERSION 文件"
    exit 1
  fi
}

main() {
  parse_args "$@"
  set_compose_file
  ensure_project_structure
  set_version_tag

  banner
  echo -e "${BOLD}开始安装 LunaFox 安全扫描平台${RESET}"
  echo -e "${DIM}将创建数据目录、生成证书并启动服务${RESET}"
  echo ""

  step_header 1 $TOTAL_STEPS "系统环境校验"
  check_system

  step_header 2 $TOTAL_STEPS "初始化数据目录"
  init_data_dir
  success "数据目录已就绪"

  step_header 3 $TOTAL_STEPS "生成 HTTPS 证书"
  generate_ssl_cert
  success "证书已就绪"

  step_header 4 $TOTAL_STEPS "生成环境配置"
  write_env_file
  success "配置文件已生成"

  step_header 5 $TOTAL_STEPS "准备 Agent/Worker 镜像"
  build_dev_images

  step_header 6 $TOTAL_STEPS "启动 Docker 服务"
  start_services

  wait_for_health
  prewarm_frontend
  install_agent_worker
  print_summary
}

main "$@"
