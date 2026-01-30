$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$Token = "{{.Token}}"
$ServerUrl = "{{.ServerURL}}"
$AgentImage = "{{.AgentImage}}"
$AgentVersion = "{{.AgentVersion}}"
$DefaultWorkerToken = "{{.WorkerToken}}"
$WorkerImage = "yyhuni/lunafox-worker"

if (-not $env:WORKER_TOKEN) {
  $env:WORKER_TOKEN = $DefaultWorkerToken
}

if (-not $env:WORKER_TOKEN) {
  Write-Error "WORKER_TOKEN is required (set WORKER_TOKEN=...)"
  exit 1
}

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
  Write-Error "Docker is required. Install it first: https://docs.docker.com/desktop/install/windows-install/"
  exit 1
}

docker info 2>$null | Out-Null
if ($LASTEXITCODE -ne 0) {
  Write-Error "Docker daemon is not running or is not accessible. Please start Docker Desktop (or the Docker service)."
  exit 1
}

if (-not $env:LUNAFOX_HOSTNAME) {
  $env:LUNAFOX_HOSTNAME = $env:COMPUTERNAME
}

$body = @{ token = $Token; hostname = $env:LUNAFOX_HOSTNAME; version = $AgentVersion } | ConvertTo-Json -Compress
try {
  $apiKey = Invoke-RestMethod -Method Post -Uri "$ServerUrl/api/agents/register" -ContentType "application/json" -Body $body -TimeoutSec 30
  $apiKey = $apiKey.Trim()
} catch {
  Write-Error "Registration failed: $_"
  exit 1
}

if (-not $apiKey) {
  Write-Error "Failed to register agent"
  exit 1
}

$maxTasks = if ($env:MAX_TASKS) { $env:MAX_TASKS } else { "5" }
$cpuThreshold = if ($env:CPU_THRESHOLD) { $env:CPU_THRESHOLD } else { "85" }
$memThreshold = if ($env:MEM_THRESHOLD) { $env:MEM_THRESHOLD } else { "85" }
$diskThreshold = if ($env:DISK_THRESHOLD) { $env:DISK_THRESHOLD } else { "90" }

$dataDir = if ($env:LUNAFOX_DATA_DIR) { $env:LUNAFOX_DATA_DIR } else { "C:\lunafox" }
if (-not (Test-Path $dataDir)) {
  New-Item -ItemType Directory -Path $dataDir | Out-Null
}

Write-Output "Installing LunaFox Agent $AgentVersion..."
Write-Output "Pulling agent image..."
docker pull "$AgentImage`:$AgentVersion"

Write-Output "Pulling worker image..."
docker pull "$WorkerImage`:$AgentVersion"

docker rm -f lunafox-agent 2>$null | Out-Null
docker run -d --restart unless-stopped --name lunafox-agent `
  --hostname $env:LUNAFOX_HOSTNAME `
  -e SERVER_URL="$ServerUrl" `
  -e API_KEY="$apiKey" `
  -e WORKER_TOKEN="$env:WORKER_TOKEN" `
  -e AGENT_VERSION="$AgentVersion" `
  -e MAX_TASKS="$maxTasks" `
  -e CPU_THRESHOLD="$cpuThreshold" `
  -e MEM_THRESHOLD="$memThreshold" `
  -e DISK_THRESHOLD="$diskThreshold" `
  -e LUNAFOX_HOSTNAME="$env:LUNAFOX_HOSTNAME" `
  -v //var/run/docker.sock:/var/run/docker.sock `
  -v "$dataDir`:/opt/lunafox" `
  "$AgentImage`:$AgentVersion" | Out-Null

Write-Output "Agent installed and running (container: lunafox-agent)"
