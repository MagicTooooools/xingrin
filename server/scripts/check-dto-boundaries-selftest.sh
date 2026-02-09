#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHECK_SCRIPT="$ROOT_DIR/scripts/check-dto-boundaries.sh"
TMP_MODULE="$ROOT_DIR/internal/modules/.dto_guard_selftest_tmp"
SNAPSHOT_TMP="$ROOT_DIR/internal/modules/snapshot/dto/_selftest_violation.go"

cleanup() {
  rm -rf "$TMP_MODULE"
  rm -f "$SNAPSHOT_TMP"
}
trap cleanup EXIT

if [[ ! -x "$CHECK_SCRIPT" ]]; then
  echo "❌ 未找到可执行守卫脚本: $CHECK_SCRIPT"
  exit 1
fi

echo "[1/2] 校验 dto 模型文件违规能被拦截..."
mkdir -p "$TMP_MODULE/demo/dto"
cat > "$TMP_MODULE/demo/dto/target_models.go" <<'CASE1'
package dto

import shared "github.com/yyhuni/lunafox/server/internal/dto"

type IllegalAlias = shared.PaginationQuery
CASE1

if bash "$CHECK_SCRIPT" >/tmp/dto-guard-selftest-case1.log 2>&1; then
  echo "❌ 自测失败：dto 模型文件违规未被拦截"
  cat /tmp/dto-guard-selftest-case1.log
  exit 1
fi

echo "[2/2] 校验 snapshot/dto 跨模块依赖违规能被拦截..."
rm -rf "$TMP_MODULE"
cat > "$SNAPSHOT_TMP" <<'CASE2'
package dto

import _ "github.com/yyhuni/lunafox/server/internal/modules/asset/dto"
CASE2

if bash "$CHECK_SCRIPT" >/tmp/dto-guard-selftest-case2.log 2>&1; then
  echo "❌ 自测失败：snapshot/dto 违规依赖未被拦截"
  cat /tmp/dto-guard-selftest-case2.log
  exit 1
fi

echo "✅ DTO 守卫自测通过（两类违规均可被正确拦截）"
