#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MODULE_ROOT="$ROOT_DIR/internal/modules"

if [[ ! -d "$MODULE_ROOT" ]]; then
  echo "ℹ️ 未找到模块目录，跳过 layer 依赖检查: $MODULE_ROOT"
  exit 0
fi

VIOLATIONS=""
append_violation() {
  local title="$1"
  local body="$2"
  if [[ -n "$VIOLATIONS" ]]; then
    VIOLATIONS+=$'\n\n'
  fi
  VIOLATIONS+="$title"
  VIOLATIONS+=$'\n'
  VIOLATIONS+="$body"
}

collect_files() {
  local pattern="$1"
  find "$MODULE_ROOT" -type f -path "$pattern" ! -name '*_test.go' | sort
}

# 1) handler 层禁止依赖 legacy/service、legacy model 与 persistence 实现
HANDLER_FILES=()
while IFS= read -r file; do HANDLER_FILES+=("$file"); done < <(collect_files '*/handler/*.go')
while IFS= read -r file; do HANDLER_FILES+=("$file"); done < <(collect_files '*/handler/*/*.go')
if [[ ${#HANDLER_FILES[@]} -gt 0 ]]; then
  output="$(rg -n --no-heading \
    -e 'internal/modules/.+/service' \
    -e 'internal/modules/.+/model' \
    -e 'internal/modules/.+/repository/persistence' \
    "${HANDLER_FILES[@]}" || true)"
  if [[ -n "$output" ]]; then
    append_violation "禁止 handler 层直接依赖 service/legacy model/repository-persistence" "$output"
  fi
fi

# 2) application（核心）禁止回退到 web/中间件层
APP_CORE_FILES=()
while IFS= read -r file; do APP_CORE_FILES+=("$file"); done < <(
  find "$MODULE_ROOT" -type f -path '*/application/*.go' ! -name '*_test.go' ! -name 'facade_*.go' | sort
)
if [[ ${#APP_CORE_FILES[@]} -gt 0 ]]; then
  output="$(rg -n --no-heading -e 'internal/modules/.+/(handler|router)' -e 'internal/middleware/' "${APP_CORE_FILES[@]}" || true)"
  if [[ -n "$output" ]]; then
    append_violation "禁止 application 核心层依赖 handler/router/middleware" "$output"
  fi
fi

# 2.1) 已 strict 收口模块：application 禁止依赖 repository 实现层
STRICT_MODULES=(agent security catalog identity asset scan snapshot)
for module in "${STRICT_MODULES[@]}"; do
  APP_FILES=()
  while IFS= read -r file; do APP_FILES+=("$file"); done < <(
    find "$MODULE_ROOT/$module/application" -type f -name '*.go' ! -name '*_test.go' | sort 2>/dev/null || true
  )
  if [[ ${#APP_FILES[@]} -eq 0 ]]; then
    continue
  fi
  output="$(rg -n --no-heading -e "internal/modules/$module/repository" "${APP_FILES[@]}" || true)"
  if [[ -n "$output" ]]; then
    append_violation "禁止 $module application 依赖 repository（含 persistence）" "$output"
  fi
done

# 3) repository 层禁止依赖 handler/router
REPO_FILES=()
while IFS= read -r file; do REPO_FILES+=("$file"); done < <(collect_files '*/repository/*.go')
if [[ ${#REPO_FILES[@]} -gt 0 ]]; then
  output="$(rg -n --no-heading -e 'internal/modules/.+/(handler|router)' "${REPO_FILES[@]}" || true)"
  if [[ -n "$output" ]]; then
    append_violation "禁止 repository 层依赖 handler/router" "$output"
  fi
fi

if [[ -n "$VIOLATIONS" ]]; then
  echo "❌ layer 依赖检查失败"
  echo "$VIOLATIONS"
  exit 1
fi

echo "✅ layer 依赖检查通过"
