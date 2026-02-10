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

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
for lib in "$SCRIPT_DIR/lib"/*.sh; do
  if [ -f "$lib" ]; then
    # shellcheck source=/dev/null
    . "$lib"
  fi
done

# ==============================================================================
# LunaFox install script
#   ./install.sh        # default production mode
#   ./install.sh --dev  # development mode
# ==============================================================================

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
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

# Special effects
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



check_system() {
  if ! lf_system_check_environment "$MODE" "$ROOT_DIR" "$COMPOSE_FILE"; then
    exit 1
  fi
}

init_data_dir() {
  if ! lf_system_init_data_dir "$DATA_DIR"; then
    exit 1
  fi
}

generate_ssl_cert() {
  if ! lf_system_generate_ssl_cert "$DOCKER_DIR"; then
    exit 1
  fi
}

write_env_file() {
  IMAGE_TAG="$VERSION_TAG"
  if ! lf_system_write_env_file "$ENV_FILE" "$IMAGE_TAG" "$GO111MODULE_VALUE" "$GOPROXY_VALUE"; then
    exit 1
  fi
}

build_dev_images() {
  if ! lf_compose_build_dev_images "$MODE" "$ROOT_DIR" "$GO111MODULE_VALUE" "$GOPROXY_VALUE"; then
    exit 1
  fi
}

start_services() {
  if ! lf_compose_start_services "$MODE" "$ENV_FILE" "$COMPOSE_FILE"; then
    exit 1
  fi
}

wait_for_health() {
  step_header 7 $TOTAL_STEPS "等待服务就绪"
  if ! lf_health_wait_for_services "https://localhost:8083/health" "https://localhost:8083/" 120 2 2 6 "$COMPOSE_FILE"; then
    exit 1
  fi
}

prewarm_frontend() {
  if ! lf_health_prewarm_frontend "$MODE" "https://localhost:8083" 30 2 6; then
    exit 1
  fi
}

print_summary() {
  lf_ui_print_install_summary "$IMAGE_TAG" "$COMPOSE_FILE"
}

install_agent_worker() {
  step_header 8 $TOTAL_STEPS "安装本地 Agent/Worker"

  local server_url="${INSTALL_SERVER_URL:-https://localhost:8083}"
  local agent_server_url=""
  local agent_register_url=""
  local network_name=""

  agent_server_url="$(lf_agent_server_url_value "http://server:8080")"
  agent_register_url="$(lf_agent_register_url_value "$server_url")"
  network_name="$(lf_agent_network_name_value "lunafox_network")"

  if ! lf_agent_install_local_worker "$MODE" "$ALLOW_PULL" "$server_url" "$agent_server_url" "$agent_register_url" "$network_name" "$ENV_FILE" "admin" "admin"; then
    exit 1
  fi
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
  if ! lf_system_ensure_project_structure "$DOCKER_DIR" "$COMPOSE_FILE"; then
    exit 1
  fi
}

set_version_tag() {
  if ! VERSION_TAG="$(lf_system_resolve_version_tag "$MODE" "$ROOT_DIR")"; then
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
