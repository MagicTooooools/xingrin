import type { Agent } from "@/types/agent.types"

export type AgentSnapshot = {
  status: string
  connectedAt?: string | null
  lastHeartbeat?: string | null
  heartbeatUpdatedAt?: string | null
  createdAt?: string | null
}

export const normalizeOrigin = (value: string): string => value.replace(/\/+$/, "")

export const isLocalOrigin = (origin: string): boolean => {
  try {
    const url = new URL(origin)
    return url.hostname === "localhost" || url.hostname === "127.0.0.1" || url.hostname === "::1"
  } catch {
    return false
  }
}

export const buildInstallCommand = (token: string, registerUrl: string): string => {
  if (!token) return ""
  const trimmedUrl = registerUrl.trim()
  if (!trimmedUrl) return ""

  const scriptBaseUrl = normalizeOrigin(trimmedUrl)
  return `curl -kfsSL "${scriptBaseUrl}/api/agents/install-script?token=${token}" | LUNAFOX_AGENT_REGISTER_URL="${trimmedUrl}" bash`
}

const toTimestamp = (value?: string | null) => {
  if (!value) return 0
  const timestamp = new Date(value).getTime()
  return Number.isNaN(timestamp) ? 0 : timestamp
}

const isAfterStart = (value: string | null | undefined, startAt: number) =>
  toTimestamp(value) > startAt

export const buildAgentSnapshot = (agents: Agent[]) => {
  const snapshot = new Map<number, AgentSnapshot>()
  agents.forEach((agent) => {
    snapshot.set(agent.id, {
      status: agent.status,
      connectedAt: agent.connectedAt ?? null,
      lastHeartbeat: agent.lastHeartbeat ?? null,
      heartbeatUpdatedAt: agent.heartbeat?.updatedAt ?? null,
      createdAt: agent.createdAt ?? null,
    })
  })
  return snapshot
}

export const detectNewAgent = (
  agents: Agent[],
  baseline: Map<number, AgentSnapshot>,
  startAt: number
) => {
  for (const agent of agents) {
    if (baseline.has(agent.id)) continue

    const recent =
      isAfterStart(agent.createdAt, startAt) ||
      isAfterStart(agent.connectedAt, startAt) ||
      isAfterStart(agent.lastHeartbeat, startAt) ||
      isAfterStart(agent.heartbeat?.updatedAt, startAt)
    if (!recent) continue

    if (
      agent.status !== "online" &&
      !isAfterStart(agent.lastHeartbeat, startAt) &&
      !isAfterStart(agent.heartbeat?.updatedAt, startAt)
    ) {
      continue
    }

    return true
  }
  return false
}
