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
# LunaFox start script
#   ./start.sh        # default production mode
#   ./start.sh --dev  # development mode
# ==============================================================================

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DOCKER_DIR="$ROOT_DIR/docker"
ENV_FILE="$DOCKER_DIR/.env"
COMPOSE_DEV="$DOCKER_DIR/docker-compose.dev.yml"
COMPOSE_PROD="$DOCKER_DIR/docker-compose.yml"

MODE="prod"
VERSION_TAG="unknown"
COMPOSE_FILE=""
COMPOSE_CMD=()
COMPOSE_ENV=()
DOCKER_PREFIX=()
DOCKER_BIN=""
ENV_FILE_ARGS=()
PROFILE_ARGS=()

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
DIM='\033[2m'
RESET='\033[0m'

# Gradient colors
GRADIENT_1='\033[38;5;39m'
GRADIENT_2='\033[38;5;45m'
GRADIENT_3='\033[38;5;51m'
GRADIENT_4='\033[38;5;87m'
GRADIENT_5='\033[38;5;123m'

usage() {
  cat <<'EOF'
用法:
  ./start.sh        # 默认生产模式
  ./start.sh --dev  # 开发模式
EOF
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
  echo -e "  ${DIM}动作:${RESET} ${BLUE}启动服务${RESET}"
  echo ""
}

error() {
  echo -e "${RED}[ERROR]${RESET} $*" >&2
}

info() {
  echo -e "${CYAN}[INFO]${RESET} $*"
}

success() {
  echo -e "${GREEN}[OK]${RESET} $*"
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
  if [ ! -f "$ENV_FILE" ]; then
    error "未找到 docker/.env，请先运行 ./install.sh。"
    exit 1
  fi
  if ! grep -q "^IMAGE_TAG=" "$ENV_FILE"; then
    error "docker/.env 缺少 IMAGE_TAG"
    exit 1
  fi
  if [ "$MODE" = "prod" ]; then
    if ! grep -q "^JWT_SECRET=" "$ENV_FILE"; then
      error "docker/.env 缺少 JWT_SECRET"
      exit 1
    fi
    if ! grep -q "^WORKER_TOKEN=" "$ENV_FILE"; then
      error "docker/.env 缺少 WORKER_TOKEN"
      exit 1
    fi
    if ! grep -q "^DB_HOST=" "$ENV_FILE"; then
      error "docker/.env 缺少 DB_HOST"
      exit 1
    fi
    if ! grep -q "^DB_PASSWORD=" "$ENV_FILE"; then
      error "docker/.env 缺少 DB_PASSWORD"
      exit 1
    fi
    if ! grep -q "^REDIS_HOST=" "$ENV_FILE"; then
      error "docker/.env 缺少 REDIS_HOST"
      exit 1
    fi
  fi
  if [ "$MODE" != "dev" ] && [ ! -f "$ROOT_DIR/VERSION" ]; then
    error "未找到 VERSION 文件"
    exit 1
  fi
  if [ ! -f "$COMPOSE_FILE" ]; then
    error "未找到 compose 文件: $COMPOSE_FILE"
    exit 1
  fi
  info "使用 docker 前缀: ${DOCKER_PREFIX[*]:-无}"
  info "使用 compose 命令: ${COMPOSE_CMD[*]}"
  success "环境校验通过"
}

parse_args() {
  for arg in "$@"; do
    case "$arg" in
      --dev) MODE="dev" ;;
      -h|--help) usage; exit 0 ;;
      *) ;;
    esac
  done
}

set_compose_file() {
  if [ "$MODE" = "prod" ]; then
    COMPOSE_FILE="$COMPOSE_PROD"
  else
    COMPOSE_FILE="$COMPOSE_DEV"
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
    VERSION_TAG="unknown"
  fi
}

confirm_action() {
  echo -ne "${BOLD}${CYAN}[?] 确认开始启动？(y/N) ${RESET}"
  read -r confirm
  if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    error "已取消启动"
    exit 1
  fi
}

set_profile_args() {
  PROFILE_ARGS=()
  if [ "$MODE" = "prod" ]; then
    PROFILE_ARGS=(--profile local-db)
  fi
}

start_local_agent() {
  local stopped_agents=()
  while IFS= read -r container_name; do
    [ -n "$container_name" ] || continue
    stopped_agents+=("$container_name")
  done < <("${DOCKER_PREFIX[@]}" "$DOCKER_BIN" ps -a --filter status=exited --format "{{.Names}}" | grep -E '^lunafox-agent($|-)' || true)
  if [ "${#stopped_agents[@]}" -eq 0 ]; then
    return
  fi

  info "检测到已停止的本地 Agent 容器，正在启动: ${stopped_agents[*]}"
  "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" start "${stopped_agents[@]}" >/dev/null
  success "本地 Agent 已启动"
}

start_services() {
  ENV_FILE_ARGS=()
  if [ -f "$ENV_FILE" ]; then
    ENV_FILE_ARGS=(--env-file "$ENV_FILE")
  fi
  set_profile_args
  "${COMPOSE_CMD[@]}" "${ENV_FILE_ARGS[@]}" -f "$COMPOSE_FILE" "${PROFILE_ARGS[@]}" up -d
  start_local_agent
  success "服务已启动"
}

prewarm_frontend() {
  if [ "$MODE" != "dev" ]; then
    return
  fi
  if ! command -v curl >/dev/null 2>&1; then
    info "未检测到 curl，跳过前端预热"
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
      info "预热未完成: ${path}（可稍后访问触发编译）"
    fi
  done
}

main() {
  parse_args "$@"
  set_compose_file
  set_version_tag

  banner
  confirm_action

  if [ "$MODE" = "dev" ]; then
    info "模式：开发"
  else
    info "模式：生产"
  fi

  check_system
  start_services
  prewarm_frontend
}

main "$@"
