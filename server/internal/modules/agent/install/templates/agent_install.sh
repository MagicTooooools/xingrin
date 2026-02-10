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

TOKEN="{{.Token}}"
SERVER_URL="{{.ServerURL}}"
REGISTER_URL="${LUNAFOX_AGENT_REGISTER_URL:-}"
AGENT_SERVER_URL="${LUNAFOX_AGENT_SERVER_URL:-$REGISTER_URL}"
NETWORK_NAME="${LUNAFOX_AGENT_DOCKER_NETWORK:-off}"
AGENT_IMAGE="{{.AgentImage}}"
AGENT_VERSION="{{.AgentVersion}}"
DEFAULT_WORKER_TOKEN="{{.WorkerToken}}"
WORKER_IMAGE="yyhuni/lunafox-worker"
SKIP_PULL="${LUNAFOX_AGENT_SKIP_PULL:-}"
LOCAL_AGENT_CONFIG="${LUNAFOX_AGENT_USE_LOCAL_LIMITS:-}"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "$1 is required" >&2
    exit 1
  fi
}

require_cmd curl

if [ -z "$REGISTER_URL" ]; then
  echo "LUNAFOX_AGENT_REGISTER_URL is required (no defaults)." >&2
  exit 1
fi

echo "Configuration:"
echo "Register URL: $REGISTER_URL"
echo "Agent server URL: $AGENT_SERVER_URL"
echo "Network: $NETWORK_NAME"

curl_opts=("-fsSL" "--connect-timeout" "10" "--max-time" "30" "-k")

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker is required. Install it first: https://docs.docker.com/engine/install/" >&2
  exit 1
fi

DOCKER_CMD="docker"
if ! docker info >/dev/null 2>&1; then
  if sudo docker info >/dev/null 2>&1; then
    DOCKER_CMD="sudo docker"
  else
    echo "Docker daemon is not running or is not accessible." >&2
    echo "Please start Docker and ensure the current user can access the Docker daemon." >&2
    exit 1
  fi
fi

NETWORK_ARGS=()
if [ -n "$NETWORK_NAME" ] && [ "$NETWORK_NAME" != "off" ] && [ "$NETWORK_NAME" != "none" ]; then
  if $DOCKER_CMD network inspect "$NETWORK_NAME" >/dev/null 2>&1; then
    NETWORK_ARGS=(--network "$NETWORK_NAME")
  else
    echo "Docker network '$NETWORK_NAME' not found, using default bridge." >&2
  fi
fi

image_exists() {
  $DOCKER_CMD image inspect "$1" >/dev/null 2>&1
}

if [ -z "${WORKER_TOKEN:-}" ]; then
  WORKER_TOKEN="$DEFAULT_WORKER_TOKEN"
fi

if [ -z "${WORKER_TOKEN:-}" ]; then
  echo "WORKER_TOKEN is required (export WORKER_TOKEN=...)" >&2
  exit 1
fi

HOSTNAME="${AGENT_HOSTNAME:-$(hostname)}"
DATA_DIR="${AGENT_DATA_DIR:-/opt/lunafox}"

echo "Installing LunaFox Agent $AGENT_VERSION..."
echo "Registering agent..."
MAX_TASKS_VALUE=""
CPU_THRESHOLD_VALUE=""
MEM_THRESHOLD_VALUE=""
DISK_THRESHOLD_VALUE=""

if [ "$LOCAL_AGENT_CONFIG" = "1" ] || [ "$LOCAL_AGENT_CONFIG" = "true" ]; then
  MAX_TASKS_VALUE="${LUNAFOX_AGENT_MAX_TASKS:-10}"
  CPU_THRESHOLD_VALUE="${LUNAFOX_AGENT_CPU_THRESHOLD:-80}"
  MEM_THRESHOLD_VALUE="${LUNAFOX_AGENT_MEM_THRESHOLD:-80}"
  DISK_THRESHOLD_VALUE="${LUNAFOX_AGENT_DISK_THRESHOLD:-85}"
fi

REGISTER_PAYLOAD=$(printf '{"token":"%s","hostname":"%s","version":"%s"' "$TOKEN" "$HOSTNAME" "$AGENT_VERSION")
if [ -n "$MAX_TASKS_VALUE" ]; then
  REGISTER_PAYLOAD=$(printf '%s,"maxTasks":%s,"cpuThreshold":%s,"memThreshold":%s,"diskThreshold":%s' \
    "$REGISTER_PAYLOAD" "$MAX_TASKS_VALUE" "$CPU_THRESHOLD_VALUE" "$MEM_THRESHOLD_VALUE" "$DISK_THRESHOLD_VALUE")
fi
REGISTER_PAYLOAD=$(printf '%s}' "$REGISTER_PAYLOAD")

RESPONSE=$(curl "${curl_opts[@]}" \
  -X POST "$REGISTER_URL/api/agent/register" \
  -H "Content-Type: application/json" \
  -d "$REGISTER_PAYLOAD" 2>&1) || {
  echo "Registration failed: $RESPONSE" >&2
  exit 1
}
API_KEY="$(printf '%s' "$RESPONSE" | sed -n 's/.*"apiKey"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' | head -n1)"

if [ -z "$API_KEY" ]; then
  API_KEY="${RESPONSE//$'\r'/}"
  API_KEY="${API_KEY//$'\n'/}"
fi

if [ -z "$API_KEY" ]; then
  echo "Failed to obtain API key" >&2
  exit 1
fi

sudo mkdir -p "$DATA_DIR"

echo "Pulling agent image..."
if [ "$SKIP_PULL" = "1" ] || [ "$SKIP_PULL" = "true" ]; then
  if image_exists "$AGENT_IMAGE:$AGENT_VERSION"; then
    echo "Using local agent image: $AGENT_IMAGE:$AGENT_VERSION"
  else
    echo "Local agent image not found and LUNAFOX_AGENT_SKIP_PULL is set." >&2
    exit 1
  fi
else
  if image_exists "$AGENT_IMAGE:$AGENT_VERSION"; then
    echo "Local agent image exists, skip pull: $AGENT_IMAGE:$AGENT_VERSION"
  else
    $DOCKER_CMD pull "$AGENT_IMAGE:$AGENT_VERSION"
  fi
fi

echo "Pulling worker image..."
if [ "$SKIP_PULL" = "1" ] || [ "$SKIP_PULL" = "true" ]; then
  if image_exists "$WORKER_IMAGE:$AGENT_VERSION"; then
    echo "Using local worker image: $WORKER_IMAGE:$AGENT_VERSION"
  else
    echo "Local worker image not found and LUNAFOX_AGENT_SKIP_PULL is set." >&2
    exit 1
  fi
else
  if image_exists "$WORKER_IMAGE:$AGENT_VERSION"; then
    echo "Local worker image exists, skip pull: $WORKER_IMAGE:$AGENT_VERSION"
  else
    $DOCKER_CMD pull "$WORKER_IMAGE:$AGENT_VERSION"
  fi
fi

$DOCKER_CMD rm -f lunafox-agent >/dev/null 2>&1 || true
$DOCKER_CMD run -d --restart unless-stopped --name lunafox-agent \
  "${NETWORK_ARGS[@]}" \
  --hostname "$HOSTNAME" \
  -e SERVER_URL="$AGENT_SERVER_URL" \
  -e API_KEY="$API_KEY" \
  -e WORKER_TOKEN="$WORKER_TOKEN" \
  -e AGENT_VERSION="$AGENT_VERSION" \
  -e AGENT_HOSTNAME="$HOSTNAME" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$DATA_DIR:/opt/lunafox" \
  "$AGENT_IMAGE:$AGENT_VERSION" >/dev/null

echo "Agent installed and running (container: lunafox-agent)"
