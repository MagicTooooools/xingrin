#!/usr/bin/env bash

lf_system_compose_plugin_path() {
  local paths=()
  local user_home=""
  local path_item=""
  local joined=""

  if [ -n "${SUDO_USER:-}" ]; then
    user_home="$(eval echo "~$SUDO_USER")"
  fi

  if [ -d "$HOME/.docker/cli-plugins" ]; then
    paths+=("$HOME/.docker/cli-plugins")
  fi
  if [ -n "$user_home" ] && [ -d "$user_home/.docker/cli-plugins" ]; then
    paths+=("$user_home/.docker/cli-plugins")
  fi

  for path_item in /usr/local/lib/docker/cli-plugins /usr/libexec/docker/cli-plugins /usr/lib/docker/cli-plugins; do
    if [ -d "$path_item" ]; then
      paths+=("$path_item")
    fi
  done

  for path_item in "${paths[@]}"; do
    joined="${joined:+$joined:}$path_item"
  done
  printf '%s' "$joined"
}

lf_system_generate_secret() {
  "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" run --rm alpine/openssl rand -hex 32
}

lf_system_require_secret_tool() {
  if command -v docker >/dev/null 2>&1; then
    return 0
  fi

  error "未检测到 docker，无法生成密钥"
  return 1
}

lf_system_ensure_docker() {
  if ! command -v docker >/dev/null 2>&1; then
    error "未检测到 docker 命令，请先安装 Docker。"
    return 1
  fi

  DOCKER_BIN="$(command -v docker)"
  if "$DOCKER_BIN" info >/dev/null 2>&1; then
    DOCKER_PREFIX=()
    return 0
  fi

  if command -v sudo >/dev/null 2>&1; then
    if sudo "$DOCKER_BIN" info >/dev/null 2>&1; then
      DOCKER_PREFIX=(sudo)
      return 0
    fi
    error "当前环境无法使用 sudo 访问 Docker（可能被禁用或需要权限）。"
  else
    error "未检测到 sudo，无法提升权限访问 Docker。"
  fi

  error "Docker 守护进程未运行或无权限访问。请确认 Docker 已启动且当前用户有权限访问 Docker socket。"
  return 1
}

lf_system_detect_compose() {
  local plugin_path=""
  local plugin_item=""

  if [ -z "$DOCKER_BIN" ]; then
    DOCKER_BIN="$(command -v docker || true)"
  fi
  if [ -z "$DOCKER_BIN" ]; then
    error "未检测到 docker 可执行文件"
    return 1
  fi

  plugin_path="$(lf_system_compose_plugin_path)"
  if [ -n "$plugin_path" ]; then
    COMPOSE_ENV=(env DOCKER_CLI_PLUGIN_PATH="$plugin_path")
  else
    COMPOSE_ENV=()
  fi

  if command -v docker-compose >/dev/null 2>&1; then
    COMPOSE_CMD=("${DOCKER_PREFIX[@]}" "$(command -v docker-compose)")
    return 0
  fi

  for plugin_item in ${plugin_path//:/ }; do
    if [ -x "$plugin_item/docker-compose" ]; then
      COMPOSE_CMD=("${DOCKER_PREFIX[@]}" "${COMPOSE_ENV[@]}" "$DOCKER_BIN" compose)
      return 0
    fi
  done

  error "未检测到 docker compose，请先安装。"
  return 1
}

lf_system_check_environment() {
  local mode="$1"
  local root_dir="$2"
  local compose_file="$3"

  info "系统环境校验..."

  if ! lf_system_ensure_docker; then
    return 1
  fi
  if ! lf_system_detect_compose; then
    return 1
  fi

  info "使用 docker 前缀: ${DOCKER_PREFIX[*]:-无}"
  info "使用 compose 命令: ${COMPOSE_CMD[*]}"

  if ! lf_system_require_secret_tool; then
    return 1
  fi

  if [ "$mode" != "dev" ] && [ ! -f "$root_dir/VERSION" ]; then
    error "未找到 VERSION 文件"
    return 1
  fi
  if [ ! -f "$compose_file" ]; then
    error "未找到 compose 文件: $compose_file"
    return 1
  fi

  success "环境校验通过"
}

lf_system_init_data_dir() {
  local data_dir="$1"

  if [ ! -d "$data_dir" ]; then
    info "创建数据目录: $data_dir"
    if ! mkdir -p "$data_dir"/{results,logs,wordlists,workspace}; then
      error "无法创建 $data_dir，请使用 sudo 运行此脚本。"
      return 1
    fi
  fi

  chmod -R 777 "$data_dir" || true
}

lf_system_generate_ssl_cert() {
  local docker_dir="$1"
  local ssl_dir="$docker_dir/nginx/ssl"
  local fullchain="$ssl_dir/fullchain.pem"
  local privkey="$ssl_dir/privkey.pem"

  if [ -f "$fullchain" ] && [ -f "$privkey" ]; then
    info "检测到已有 HTTPS 证书，跳过生成。"
    return 0
  fi

  info "生成自签 HTTPS 证书（localhost）..."
  mkdir -p "$ssl_dir"

  if ! "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" run --rm -v "$ssl_dir:/ssl" alpine/openssl \
    req -x509 -nodes -newkey rsa:2048 -days 365 \
    -keyout /ssl/privkey.pem \
    -out /ssl/fullchain.pem \
    -subj "/C=CN/ST=NA/L=NA/O=LunaFox/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"; then
    error "证书生成失败，请检查 Docker 与 openssl 镜像是否可用"
    return 1
  fi

  if [ ! -f "$fullchain" ] || [ ! -f "$privkey" ]; then
    error "证书生成失败，请手动放置证书到 $ssl_dir"
    return 1
  fi
}

lf_system_write_env_file() {
  local env_file="$1"
  local image_tag="$2"
  local go111module_value="$3"
  local goproxy_value="$4"
  local jwt_secret=""
  local worker_token=""

  info "生成 JWT 密钥..."
  if ! jwt_secret="$(lf_system_generate_secret)"; then
    error "JWT 密钥生成失败"
    return 1
  fi

  info "生成 Worker 令牌..."
  if ! worker_token="$(lf_system_generate_secret)"; then
    error "Worker 令牌生成失败"
    return 1
  fi

  cat > "$env_file" <<EOF_ENV
IMAGE_TAG=$image_tag
JWT_SECRET=$jwt_secret
WORKER_TOKEN=$worker_token
DB_HOST=postgres
DB_PASSWORD=postgres
REDIS_HOST=redis
DB_USER=postgres
DB_NAME=lunafox
DB_PORT=5432
REDIS_PORT=6379
GO111MODULE=$go111module_value
GOPROXY=$goproxy_value
EOF_ENV
}

lf_system_ensure_project_structure() {
  local docker_dir="$1"
  local compose_file="$2"

  if [ ! -d "$docker_dir" ] || [ ! -f "$compose_file" ]; then
    error "未找到 docker 目录或 compose 文件，请确认项目结构。"
    return 1
  fi
}

lf_system_resolve_version_tag() {
  local mode="$1"
  local root_dir="$2"
  local version_file="$root_dir/VERSION"
  local version_tag=""

  if [ "$mode" = "dev" ]; then
    printf '%s' "dev"
    return 0
  fi

  if [ ! -f "$version_file" ]; then
    error "未找到 VERSION 文件"
    return 1
  fi

  version_tag="$(tr -d '[:space:]' < "$version_file")"
  if [ -z "$version_tag" ]; then
    error "VERSION 文件为空"
    return 1
  fi

  printf '%s' "$version_tag"
}
