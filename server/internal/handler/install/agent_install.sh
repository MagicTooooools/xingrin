#!/usr/bin/env bash
set -euo pipefail

TOKEN="{{.Token}}"
SERVER_URL="{{.ServerURL}}"
AGENT_IMAGE="{{.AgentImage}}"
AGENT_VERSION="{{.AgentVersion}}"
DEFAULT_WORKER_TOKEN="{{.WorkerToken}}"
WORKER_IMAGE="yyhuni/lunafox-worker"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "$1 is required" >&2
    exit 1
  fi
}

require_cmd curl

if ! command -v docker >/dev/null 2>&1; then
  echo "Docker is required. Install it first: https://docs.docker.com/engine/install/" >&2
  exit 1
fi

# Detect docker command (with or without sudo)
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

if [ -z "${WORKER_TOKEN:-}" ]; then
  WORKER_TOKEN="$DEFAULT_WORKER_TOKEN"
fi

if [ -z "${WORKER_TOKEN:-}" ]; then
  echo "WORKER_TOKEN is required (export WORKER_TOKEN=...)" >&2
  exit 1
fi

HOSTNAME="${LUNAFOX_HOSTNAME:-$(hostname)}"
DATA_DIR="${LUNAFOX_DATA_DIR:-/opt/lunafox}"

echo "Installing LunaFox Agent $AGENT_VERSION..."
echo "Registering agent..."
RESPONSE=$(curl -fsSL --connect-timeout 10 --max-time 30 \
  -X POST "$SERVER_URL/api/agents/register" \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\",\"hostname\":\"$HOSTNAME\",\"version\":\"$AGENT_VERSION\"}" 2>&1) || {
  echo "Registration failed: $RESPONSE" >&2
  exit 1
}
API_KEY="${RESPONSE//$'\r'/}"
API_KEY="${API_KEY//$'\n'/}"

if [ -z "$API_KEY" ]; then
  echo "Failed to obtain API key" >&2
  exit 1
fi

sudo mkdir -p "$DATA_DIR"

echo "Pulling agent image..."
$DOCKER_CMD pull "$AGENT_IMAGE:$AGENT_VERSION"

echo "Pulling worker image..."
$DOCKER_CMD pull "$WORKER_IMAGE:$AGENT_VERSION"

$DOCKER_CMD rm -f lunafox-agent >/dev/null 2>&1 || true
$DOCKER_CMD run -d --restart unless-stopped --name lunafox-agent \
  --hostname "$HOSTNAME" \
  -e SERVER_URL="$SERVER_URL" \
  -e API_KEY="$API_KEY" \
  -e WORKER_TOKEN="$WORKER_TOKEN" \
  -e AGENT_VERSION="$AGENT_VERSION" \
  -e MAX_TASKS="${MAX_TASKS:-5}" \
  -e CPU_THRESHOLD="${CPU_THRESHOLD:-85}" \
  -e MEM_THRESHOLD="${MEM_THRESHOLD:-85}" \
  -e DISK_THRESHOLD="${DISK_THRESHOLD:-90}" \
  -e LUNAFOX_HOSTNAME="$HOSTNAME" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$DATA_DIR:/opt/lunafox" \
  "$AGENT_IMAGE:$AGENT_VERSION" >/dev/null

echo "Agent installed and running (container: lunafox-agent)"
