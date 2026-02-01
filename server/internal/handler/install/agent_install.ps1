$ErrorActionPreference = "Stop"
Set-StrictMode -Version Latest

$Token = "{{.Token}}"
$ServerUrl = "{{.ServerURL}}"
$RegisterUrl = if ($env:LUNAFOX_REGISTER_URL) { $env:LUNAFOX_REGISTER_URL } else { "" }
$AgentServerUrl = if ($env:LUNAFOX_AGENT_SERVER_URL) { $env:LUNAFOX_AGENT_SERVER_URL } else { $RegisterUrl }
$NetworkName = if ($env:LUNAFOX_NETWORK) { $env:LUNAFOX_NETWORK } else { "off" }
$AgentImage = "{{.AgentImage}}"
$AgentVersion = "{{.AgentVersion}}"
$DefaultWorkerToken = "{{.WorkerToken}}"
$WorkerImage = "yyhuni/lunafox-worker"

if (-not $RegisterUrl) {
  Write-Error "LUNAFOX_REGISTER_URL is required (no defaults)."
  exit 1
}

Write-Output "Configuration:"
Write-Output "Register URL: $RegisterUrl"
Write-Output "Agent server URL: $AgentServerUrl"
Write-Output "Network: $NetworkName"

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

if (-not $env:AGENT_HOSTNAME) {
  $env:AGENT_HOSTNAME = $env:COMPUTERNAME
}

$body = @{ token = $Token; hostname = $env:AGENT_HOSTNAME; version = $AgentVersion } | ConvertTo-Json -Compress
$invokeParams = @{ Method = "Post"; Uri = "$RegisterUrl/api/agents/register"; ContentType = "application/json"; Body = $body; TimeoutSec = 30 }
try {
  [System.Net.ServicePointManager]::ServerCertificateValidationCallback = { $true }
} catch { }
if ($PSVersionTable.PSVersion.Major -ge 7) {
  $invokeParams.SkipCertificateCheck = $true
}
try {
  $apiKey = Invoke-RestMethod @invokeParams
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

$dataDir = if ($env:AGENT_DATA_DIR) { $env:AGENT_DATA_DIR } else { "C:\lunafox" }
if (-not (Test-Path $dataDir)) {
  New-Item -ItemType Directory -Path $dataDir | Out-Null
}

Write-Output "Installing LunaFox Agent $AgentVersion..."
Write-Output "Pulling agent image..."
docker pull "$AgentImage`:$AgentVersion"

Write-Output "Pulling worker image..."
docker pull "$WorkerImage`:$AgentVersion"

docker rm -f lunafox-agent 2>$null | Out-Null

$networkArgs = @()
if ($NetworkName -and $NetworkName -notin @("off", "none")) {
  docker network inspect $NetworkName 2>$null | Out-Null
  if ($LASTEXITCODE -eq 0) {
    $networkArgs = @("--network", $NetworkName)
  } else {
    Write-Warning "Docker network '$NetworkName' not found, using default bridge."
  }
}

$dockerArgs = @(
  "run", "-d", "--restart", "unless-stopped", "--name", "lunafox-agent"
)
if ($networkArgs.Count -gt 0) { $dockerArgs += $networkArgs }
$dockerArgs += @(
  "--hostname", $env:AGENT_HOSTNAME,
  "-e", "SERVER_URL=$AgentServerUrl",
  "-e", "API_KEY=$apiKey",
  "-e", "WORKER_TOKEN=$env:WORKER_TOKEN",
  "-e", "AGENT_VERSION=$AgentVersion",
  "-e", "MAX_TASKS=$maxTasks",
  "-e", "CPU_THRESHOLD=$cpuThreshold",
  "-e", "MEM_THRESHOLD=$memThreshold",
  "-e", "DISK_THRESHOLD=$diskThreshold",
  "-e", "AGENT_HOSTNAME=$env:AGENT_HOSTNAME",
  "-v", "//var/run/docker.sock:/var/run/docker.sock",
  "-v", "$dataDir`:/opt/lunafox",
  "$AgentImage`:$AgentVersion"
)

docker @dockerArgs | Out-Null

Write-Output "Agent installed and running (container: lunafox-agent)"
