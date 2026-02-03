"use client"

import { useMemo } from "react"
import { useTranslations } from "next-intl"
import {
  IconDotsVertical,
  IconSettings,
  IconTrash,
  IconAlertTriangle,
  IconActivity,
  IconMapPin,
  IconClock,
} from "@tabler/icons-react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {
  Card,
  CardContent,
  CardHeader,
} from "@/components/ui/card"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Separator } from "@/components/ui/separator"
import { Status, StatusIndicator } from "@/components/ui/shadcn-io/status"
import { useFormatDate, useFormatNumber } from "@/lib/i18n-format"
import { cn } from "@/lib/utils"
import type { Agent } from "@/types/agent.types"

function getHealthStyle(state: string) {
  const normalized = state.toLowerCase()
  if (normalized === "ok") {
    return "bg-[var(--success)]/10 text-[var(--success)] border-[var(--success)]/20"
  }
  if (normalized === "warning" || normalized === "warn") {
    return "bg-[var(--warning)]/10 text-[var(--warning)] border-[var(--warning)]/20"
  }
  if (normalized === "error" || normalized === "critical") {
    return "bg-[var(--error)]/10 text-[var(--error)] border-[var(--error)]/20"
  }
  return "bg-muted text-muted-foreground border-border"
}

function getStatusVariant(status: string) {
  if (status === "online") return "online"
  if (status === "offline") return "offline"
  return "maintenance"
}

function formatUptime(seconds?: number | null) {
  if (seconds === null || seconds === undefined) return "-"
  const total = Math.max(0, Math.floor(seconds))
  const minutes = Math.floor(total / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  if (days > 0) return `${days}d ${hours % 24}h`
  if (hours > 0) return `${hours}h ${minutes % 60}m`
  return `${minutes}m`
}

interface MetricProgressProps {
  label: string
  value: number
  threshold?: number
}

function MetricProgress({ label, value, threshold }: MetricProgressProps) {
  const percentage = Math.min(100, Math.max(0, value))

  const status = useMemo(() => {
    if (!threshold) return "normal"
    if (percentage >= threshold) return "critical"
    if (percentage >= threshold * 0.8) return "warning"
    return "normal"
  }, [percentage, threshold])

  const progressColor = useMemo(() => {
    if (status === "critical") return "bg-[var(--error)]"
    if (status === "warning") return "bg-[var(--warning)]"
    return "bg-[var(--success)]"
  }, [status])

  return (
    <div className="space-y-1.5">
      <div className="flex items-center justify-between text-xs">
        <span className="text-muted-foreground">{label}</span>
        <div className="flex items-center gap-1">
          <span className={cn(
            "font-medium tabular-nums",
            status === "critical" && "text-[var(--error)]",
            status === "warning" && "text-[var(--warning)]"
          )}>
            {percentage.toFixed(0)}%
          </span>
          {threshold && (
            <span className="text-muted-foreground text-[10px]">/ {threshold}%</span>
          )}
          {status !== "normal" && (
            <IconAlertTriangle className={cn(
              "h-3 w-3",
              status === "critical" ? "text-[var(--error)]" : "text-[var(--warning)]"
            )} />
          )}
        </div>
      </div>
      <div className="relative h-2 w-full bg-muted rounded-full overflow-hidden">
        <div
          className={cn("h-full transition-all duration-300", progressColor)}
          style={{ width: `${percentage}%` }}
        />
      </div>
    </div>
  )
}

interface AgentCardCompactProps {
  agent: Agent
  onConfig: (agent: Agent) => void
  onDelete: (agent: Agent) => void
}

export function AgentCardCompact({
  agent,
  onConfig,
  onDelete,
}: AgentCardCompactProps) {
  const t = useTranslations("settings.workers")
  const { formatDateTime } = useFormatDate()
  const formatNumber = useFormatNumber()


  const healthState = (agent.health?.state || "unknown").toLowerCase()
  const healthLabel = useMemo(() => {
    if (healthState === "ok") return t("health.ok")
    if (healthState === "warning" || healthState === "warn") return t("health.warning")
    if (healthState === "error" || healthState === "critical") return t("health.error")
    return t("health.unknown")
  }, [healthState, t])

  const heartbeat = agent.heartbeat

  // 检查是否有指标超过阈值
  const hasWarnings = useMemo(() => {
    if (!heartbeat) return false
    return (
      heartbeat.cpu >= agent.cpuThreshold ||
      heartbeat.mem >= agent.memThreshold ||
      heartbeat.disk >= agent.diskThreshold
    )
  }, [heartbeat, agent])

  // 计算最后心跳时间差（秒）
  const lastHeartbeatSeconds = useMemo(() => {
    if (!agent.lastHeartbeat) return null
    const now = Date.now()
    const lastHeartbeat = new Date(agent.lastHeartbeat).getTime()
    return Math.floor((now - lastHeartbeat) / 1000)
  }, [agent.lastHeartbeat])

  // 判断心跳是否过期（超过30秒）
  const isHeartbeatStale = lastHeartbeatSeconds !== null && lastHeartbeatSeconds > 30

  return (
    <Card className="transition-all duration-200 hover:shadow-md gap-0">
      <CardHeader className="pb-2">
        <div className="flex items-start justify-between gap-2 w-full min-w-0">
          <div className="flex items-center gap-2 min-w-0 flex-1">
            <Status status={getStatusVariant(agent.status)}>
              <StatusIndicator className={agent.status === "online" ? "animate-pulse" : ""} />
            </Status>
            <span className="font-medium truncate">{agent.name}</span>
          </div>
          <div className="flex items-center gap-2 shrink-0">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="icon" className="h-7 w-7">
                  <IconDotsVertical className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuLabel>{t("actions.title")}</DropdownMenuLabel>
                <DropdownMenuItem onClick={() => onConfig(agent)}>
                  <IconSettings className="h-4 w-4" />
                  {t("actions.config")}
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem variant="destructive" onClick={() => onDelete(agent)}>
                  <IconTrash className="h-4 w-4" />
                  {t("actions.delete")}
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-2.5 pt-0">
        {/* 基本信息区 */}
        <div className="space-y-2 text-xs">
          <div className="flex items-center gap-2 text-muted-foreground min-w-0">
            <IconMapPin className="h-3.5 w-3.5 shrink-0" />
            <span className="truncate">{agent.hostname || t("unknownHost")}</span>
          </div>
          <div className="flex items-center gap-2 text-muted-foreground min-w-0">
            <IconMapPin className="h-3.5 w-3.5 shrink-0 opacity-60" />
            <span className="truncate">{agent.ipAddress || t("unknownIp")}</span>
          </div>
          <div className="flex items-center gap-2">
            <IconClock className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
            <span className="text-muted-foreground">{t("metrics.lastHeartbeat")}:</span>
            <span className={cn("font-medium", isHeartbeatStale && "text-[var(--error)]")}>
              {formatDateTime(agent.lastHeartbeat)}
            </span>
            {agent.status === "online" && !isHeartbeatStale && (
              <IconActivity className="h-3.5 w-3.5 text-[var(--success)] animate-pulse" />
            )}
          </div>
          <div className="flex items-center gap-2 flex-wrap min-w-0">
            {agent.version && (
              <Badge
                variant="secondary"
                className="text-[10px] max-w-[12rem] truncate"
                title={agent.version}
              >
                {agent.version}
              </Badge>
            )}
            <Badge variant="outline" className={cn("text-[10px]", getHealthStyle(healthState))}>
              {healthLabel}
            </Badge>
            {hasWarnings && (
              <Badge variant="outline" className="bg-[var(--warning)]/10 text-[var(--warning)] border-[var(--warning)]/20 text-[10px]">
                <IconAlertTriangle className="h-3 w-3 mr-1" />
                {t("card.warning")}
              </Badge>
            )}
          </div>
        </div>

        <Separator />

        {/* 系统资源区 */}
        {heartbeat ? (
          <div className="space-y-3">
            <div className="space-y-2.5">
              <MetricProgress
                label={t("metrics.cpu")}
                value={heartbeat.cpu}
                threshold={agent.cpuThreshold}
              />
              <MetricProgress
                label={t("metrics.mem")}
                value={heartbeat.mem}
                threshold={agent.memThreshold}
              />
              <MetricProgress
                label={t("metrics.disk")}
                value={heartbeat.disk}
                threshold={agent.diskThreshold}
              />
            </div>
          </div>
        ) : (
          <div className="text-xs text-muted-foreground text-center py-4 border border-dashed rounded bg-muted/20">
            {t("card.waitingForHeartbeat")}
          </div>
        )}

        {heartbeat && (
          <>
            <Separator />

            {/* 任务和运行时间区 */}
            <div className="grid grid-cols-2 gap-3 text-xs">
              <div>
                <div className="text-muted-foreground mb-1">{t("metrics.tasks")}</div>
                <div className="font-medium text-base">
                  {formatNumber.formatInteger(heartbeat.tasks)}
                  <span className="text-sm text-muted-foreground">/{agent.maxTasks}</span>
                </div>
              </div>
              <div>
                <div className="text-muted-foreground mb-1">{t("list.uptime")}</div>
                <div className="font-medium text-base">{formatUptime(heartbeat.uptime)}</div>
              </div>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  )
}
