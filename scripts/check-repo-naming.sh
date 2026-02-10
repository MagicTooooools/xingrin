#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

bash "$ROOT_DIR/server/scripts/check-naming-conventions.sh"
node "$ROOT_DIR/scripts/check-frontend-naming.mjs"

echo "✅ 全仓命名规范检查通过"
