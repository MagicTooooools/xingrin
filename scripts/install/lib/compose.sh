#!/usr/bin/env bash

lf_compose_build_dev_images() {
  local mode="$1"
  local root_dir="$2"
  local go111module_value="$3"
  local goproxy_value="$4"
  local cache_base=""
  local agent_cache=""
  local worker_cache=""
  local buildx_inspect_file=""
  local buildx_driver=""
  local enable_cache_export="false"
  local bake_file=""

  if [ "$mode" != "dev" ]; then
    info "生产模式跳过本地构建"
    return 0
  fi

  export DOCKER_BUILDKIT=1
  export COMPOSE_DOCKER_CLI_BUILD=1
  export BUILDKIT_PROGRESS=plain
  info "已启用 BuildKit 加速构建"

  cache_base="${HOME}/.cache/lunafox-buildx"
  agent_cache="${cache_base}/agent"
  worker_cache="${cache_base}/worker"
  mkdir -p "$agent_cache" "$worker_cache"
  info "Buildx 缓存目录: $cache_base"

  info "检测 buildx driver..."
  buildx_inspect_file="$(mktemp)"
  "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" buildx inspect | tee "$buildx_inspect_file"
  buildx_driver="$(awk -F': *' '/^Driver:/ {print $2; exit}' "$buildx_inspect_file")"
  rm -f "$buildx_inspect_file"

  if [ -z "$buildx_driver" ]; then
    warn "未检测到 buildx driver，默认禁用本地缓存导出"
    enable_cache_export="false"
  elif [ "$buildx_driver" = "docker" ]; then
    warn "检测到 buildx driver=docker，不支持 cache export，已禁用本地缓存导出"
    enable_cache_export="false"
  else
    info "buildx driver: $buildx_driver（启用本地缓存导出）"
    enable_cache_export="true"
  fi

  info "并行构建 agent/worker 镜像 (buildx bake)..."
  bake_file="$(mktemp)"
  {
    cat <<EOF_BAKE
group "default" {
  targets = ["agent", "worker"]
}

target "agent" {
  context = "$root_dir"
  dockerfile = "$root_dir/agent/Dockerfile"
  tags = ["yyhuni/lunafox-agent:dev"]
  build-args = {
    BUILDKIT_INLINE_CACHE = "1"
    GO111MODULE = "$go111module_value"
    GOPROXY = "$goproxy_value"
  }
EOF_BAKE
    if [ "$enable_cache_export" = "true" ]; then
      cat <<EOF_BAKE
  cache-from = ["type=local,src=$agent_cache"]
  cache-to = ["type=local,dest=$agent_cache,mode=max"]
EOF_BAKE
    fi
    cat <<EOF_BAKE
}

target "worker" {
  context = "$root_dir/worker"
  dockerfile = "$root_dir/worker/Dockerfile"
  tags = ["yyhuni/lunafox-worker:dev"]
  build-args = {
    BUILDKIT_INLINE_CACHE = "1"
    GO111MODULE = "$go111module_value"
    GOPROXY = "$goproxy_value"
  }
EOF_BAKE
    if [ "$enable_cache_export" = "true" ]; then
      cat <<EOF_BAKE
  cache-from = ["type=local,src=$worker_cache"]
  cache-to = ["type=local,dest=$worker_cache,mode=max"]
EOF_BAKE
    fi
    cat <<'EOF_BAKE'
}
EOF_BAKE
  } > "$bake_file"

  if ! "${DOCKER_PREFIX[@]}" "$DOCKER_BIN" buildx bake -f "$bake_file" --progress=plain --load; then
    rm -f "$bake_file"
    error "Agent/Worker 镜像构建失败"
    return 1
  fi

  rm -f "$bake_file"
  success "Agent/Worker 镜像已构建"
}

lf_compose_start_services() {
  local mode="$1"
  local env_file="$2"
  local compose_file="$3"
  local profile_arg=""
  local -a env_file_args=()

  if [ -f "$env_file" ]; then
    env_file_args=(--env-file "$env_file")
  fi

  if [ "$mode" = "prod" ]; then
    profile_arg="--profile local-db"
  fi

  info "正在启动服务（自动重建所有镜像）..."
  echo ""
  if ! "${COMPOSE_CMD[@]}" "${env_file_args[@]}" -f "$compose_file" ${profile_arg:+$profile_arg} up -d --build --force-recreate; then
    error "服务启动失败"
    return 1
  fi

  success "Docker 服务启动成功"
}
