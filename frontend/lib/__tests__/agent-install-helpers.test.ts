import { describe, expect, it } from "vitest"
import {
  buildAgentSnapshot,
  buildInstallCommand,
  detectNewAgent,
  isLocalOrigin,
  normalizeOrigin,
} from "@/lib/agent-install-helpers"
import type { Agent } from "@/types/agent.types"

const makeAgent = (overrides: Partial<Agent> = {}): Agent => ({
  id: 1,
  name: "agent-1",
  status: "offline",
  maxTasks: 1,
  cpuThreshold: 80,
  memThreshold: 80,
  diskThreshold: 80,
  health: { state: "ok" },
  createdAt: "2026-01-01T00:00:00Z",
  ...overrides,
})

describe("agent install helpers", () => {
  it("normalizes origins and detects local addresses", () => {
    expect(normalizeOrigin("https://example.com/")).toBe("https://example.com")
    expect(isLocalOrigin("http://localhost:3000")).toBe(true)
    expect(isLocalOrigin("https://example.com")).toBe(false)
    expect(isLocalOrigin("not-a-url")).toBe(false)
  })

  it("builds install commands consistently", () => {
    expect(buildInstallCommand("", "https://example.com")).toBe("")
    expect(buildInstallCommand("token", "")).toBe("")

    const command = buildInstallCommand("token", "https://example.com/")
    expect(command).toContain("/api/agents/install-script?token=token")
    expect(command).toContain('LUNAFOX_AGENT_REGISTER_URL="https://example.com/"')
  })

  it("detects new agents compared to baseline", () => {
    const baseline = buildAgentSnapshot([makeAgent()])
    const startAt = new Date("2026-01-01T00:00:00Z").getTime()

    const newAgent = makeAgent({
      id: 2,
      status: "online",
      createdAt: "2026-01-02T00:00:00Z",
      connectedAt: "2026-01-02T00:00:00Z",
      lastHeartbeat: "2026-01-02T00:00:00Z",
    })
    expect(detectNewAgent([newAgent], baseline, startAt)).toBe(true)

    const oldAgent = makeAgent({
      id: 3,
      status: "offline",
      createdAt: "2025-12-01T00:00:00Z",
    })
    expect(detectNewAgent([oldAgent], baseline, startAt)).toBe(false)
  })
})
