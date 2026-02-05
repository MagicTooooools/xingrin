"use client"

import React from "react"
import Link from "next/link"
import dynamic from "next/dynamic"
import { useTranslations, useLocale } from "next-intl"
import {
  Globe,
  Network,
  Server,
  Link2,
  FolderOpen,
  Camera,
  AlertTriangle,
  Clock,
  Calendar,
  ChevronRight,
  CheckCircle2,
  XCircle,
  Loader2,
  Cpu,
  HardDrive,
} from "@/components/icons"
import {
  IconCircleCheck,
  IconCircleX,
  IconClock,
} from "@/components/icons"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { Badge } from "@/components/ui/badge"
import { Switch } from "@/components/ui/switch"
import { Label } from "@/components/ui/label"
import { Separator } from "@/components/ui/separator"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useScan } from "@/hooks/use-scans"
import { useScanLogs } from "@/hooks/use-scan-logs"
import { ScanLogList } from "@/components/scan/scan-log-list"
import { cn } from "@/lib/utils"
import type { ScanRecord, ScanStatus, StageProgressItem, StageStatus } from "@/types/scan.types"

// Dynamic import for YamlEditor (only loaded when config tab is active)
const YamlEditor = dynamic(() => import('@/components/ui/yaml-editor').then(m => ({ default: m.YamlEditor })), {
  loading: () => <div className="flex items-center justify-center h-full text-muted-foreground text-sm">加载编辑器中...</div>,
  ssr: false
})

interface ScanOverviewProps {
  scanId: number
}

/**
 * Scan overview component
 * Displays statistics cards for the scan results
 */
// Pulsing dot animation
function PulsingDot({ className }: { className?: string }) {
  return (
    <span className={cn("relative flex h-3 w-3", className)}>
      <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-current opacity-75" />
      <span className="relative inline-flex h-3 w-3 rounded-full bg-current" />
    </span>
  )
}

// Stage status icon
function StageStatusIcon({ status }: { status: StageStatus }) {
  switch (status) {
    case "completed":
      return <IconCircleCheck className="h-5 w-5 text-[var(--success)]" />
    case "running":
      return <PulsingDot className="text-[var(--warning)]" />
    case "failed":
      return <IconCircleX className="h-5 w-5 text-[var(--error)]" />
    case "cancelled":
      return <IconCircleX className="h-5 w-5 text-muted-foreground" />
    default:
      return <IconClock className="h-5 w-5 text-muted-foreground" />
  }
}

// Format duration (seconds -> readable string)
function formatStageDuration(seconds?: number): string | undefined {
  if (seconds === undefined || seconds === null) return undefined
  if (seconds < 1) return "<1s"
  if (seconds < 60) return `${Math.round(seconds)}s`
  const minutes = Math.floor(seconds / 60)
  const secs = Math.round(seconds % 60)
  return secs > 0 ? `${minutes}m ${secs}s` : `${minutes}m`
}

// Status priority for sorting (lower = higher priority)
const STAGE_STATUS_PRIORITY: Record<StageStatus, number> = {
  running: 0,
  pending: 1,
  completed: 2,
  failed: 3,
  cancelled: 4,
}

// Status style configuration (consistent with scan-history-columns)
const SCAN_STATUS_STYLES: Record<string, string> = {
  running: "bg-[var(--warning)]/10 text-[var(--warning)] border-[var(--warning)]/20",
  cancelled: "bg-muted/10 text-muted-foreground border-muted/20",
  completed: "bg-[var(--success)]/10 text-[var(--success)] border-[var(--success)]/20",
  failed: "bg-[var(--error)]/10 text-[var(--error)] border-[var(--error)]/20",
  pending: "bg-[var(--info)]/10 text-[var(--info)] border-[var(--info)]/20",
}

/**
 * Format date helper function
 */
function formatDate(dateString: string | undefined, locale: string): string {
  if (!dateString) return "-"
  const localeStr = locale === 'zh' ? 'zh-CN' : 'en-US'
  return new Date(dateString).toLocaleString(localeStr, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  })
}

/**
 * Calculate duration between two dates
 */
function formatDuration(startedAt: string | undefined, completedAt: string | undefined): string {
  if (!startedAt) return "-"
  const start = new Date(startedAt)
  const end = completedAt ? new Date(completedAt) : new Date()
  const diffMs = end.getTime() - start.getTime()
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMins / 60)
  const remainingMins = diffMins % 60

  if (diffHours > 0) {
    return `${diffHours}h ${remainingMins}m`
  }
  return `${diffMins}m`
}

/**
 * Get status icon configuration
 */
function getStatusIcon(status: string) {
  switch (status) {
    case "completed":
      return { icon: CheckCircle2, animate: false }
    case "running":
      return { icon: Loader2, animate: true }
    case "failed":
      return { icon: XCircle, animate: false }
    case "cancelled":
      return { icon: XCircle, animate: false }
    case "pending":
      return { icon: Loader2, animate: true }
    default:
      return { icon: Clock, animate: false }
  }
}

export function ScanOverview({ scanId }: ScanOverviewProps) {
  const t = useTranslations("scan.history.overview")
  const tStatus = useTranslations("scan.history.status")
  const tProgress = useTranslations("scan.progress")
  const locale = useLocale()

  const { data: scan, isLoading, error } = useScan(scanId)

  type LegacyVulnerabilitySummary = {
    total?: number
    critical?: number
    high?: number
    medium?: number
    low?: number
  }

  type LegacyScanSummary = {
    subdomains?: number
    websites?: number
    endpoints?: number
    ips?: number
    directories?: number
    screenshots?: number
    vulnerabilities?: LegacyVulnerabilitySummary
  }

  type ScanRecordWithLegacy = ScanRecord & {
    summary?: LegacyScanSummary
    startedAt?: string
    completedAt?: string
  }

  const scanWithLegacy = scan as ScanRecordWithLegacy | undefined

  // Memoize isRunning to avoid unnecessary recalculations
  const isRunning = React.useMemo(
    () => scan?.status === 'running' || scan?.status === 'pending',
    [scan?.status]
  )

  // Auto-refresh state (default: on when running)
  const [autoRefresh, setAutoRefresh] = React.useState(true)

  // Tab state for logs/config
  const [activeTab, setActiveTab] = React.useState<'logs' | 'config'>('logs')

  // Logs hook
  const { logs, loading: logsLoading } = useScanLogs({
    scanId,
    enabled: !!scan && activeTab === 'logs',
    pollingInterval: isRunning && autoRefresh && activeTab === 'logs' ? 3000 : 0,
  })

  // Memoize derived values to avoid unnecessary recalculations
  const summary = React.useMemo(() => {
    const stats = scan?.cachedStats
    const legacy = scanWithLegacy?.summary
    return {
      subdomains: stats?.subdomainsCount ?? legacy?.subdomains ?? 0,
      websites: stats?.websitesCount ?? legacy?.websites ?? 0,
      endpoints: stats?.endpointsCount ?? legacy?.endpoints ?? 0,
      ips: stats?.ipsCount ?? legacy?.ips ?? 0,
      directories: stats?.directoriesCount ?? legacy?.directories ?? 0,
      screenshots: stats?.screenshotsCount ?? legacy?.screenshots ?? 0,
    }
  }, [scan, scanWithLegacy])

  const vulnSummary = React.useMemo(() => {
    const stats = scan?.cachedStats
    const legacy = scanWithLegacy?.summary
    return (
      legacy?.vulnerabilities || {
        total: stats?.vulnsTotal ?? 0,
        critical: stats?.vulnsCritical ?? 0,
        high: stats?.vulnsHigh ?? 0,
        medium: stats?.vulnsMedium ?? 0,
        low: stats?.vulnsLow ?? 0,
      }
    )
  }, [scan, scanWithLegacy])

  const status = (scan?.status ?? "pending") as ScanStatus
  const statusIconConfig = React.useMemo(() => getStatusIcon(status), [status])
  const StatusIcon = statusIconConfig.icon
  const statusStyle = SCAN_STATUS_STYLES[status] || "bg-muted text-muted-foreground"
  const startedAt = React.useMemo(
    () => scanWithLegacy?.startedAt || scan?.createdAt,
    [scan, scanWithLegacy]
  )
  const completedAt = React.useMemo(() => scanWithLegacy?.completedAt, [scanWithLegacy])

  const assetCards = React.useMemo(
    () => [
      {
        title: t("cards.websites"),
        value: summary.websites || 0,
        icon: Globe,
        code: "DAT-WEB",
        href: `/scan/history/${scanId}/websites/`,
      },
      {
        title: t("cards.subdomains"),
        value: summary.subdomains || 0,
        icon: Network,
        code: "DAT-SUB",
        href: `/scan/history/${scanId}/subdomain/`,
      },
      {
        title: t("cards.ips"),
        value: summary.ips || 0,
        icon: Server,
        code: "DAT-IP",
        href: `/scan/history/${scanId}/ip-addresses/`,
      },
      {
        title: t("cards.urls"),
        value: summary.endpoints || 0,
        icon: Link2,
        code: "DAT-URL",
        href: `/scan/history/${scanId}/endpoints/`,
      },
      {
        title: t("cards.directories"),
        value: summary.directories || 0,
        icon: FolderOpen,
        code: "DAT-DIR",
        href: `/scan/history/${scanId}/directories/`,
      },
      {
        title: t("cards.screenshots"),
        value: summary.screenshots || 0,
        icon: Camera,
        code: "DAT-SCR",
        href: `/scan/history/${scanId}/screenshots/`,
      },
    ],
    [summary, scanId, t]
  )

  if (isLoading) {
    return (
      <div className="space-y-6">
        {/* Stats cards skeleton */}
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {[...Array(6)].map((_, i) => (
            <Card key={i}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="h-4 w-4" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-8 w-16" />
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    )
  }

  if (error || !scan) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <AlertTriangle className="h-10 w-10 text-destructive mb-4" />
        <p className="text-muted-foreground">{t("loadError")}</p>
      </div>
    )
  }

  return (
    <div className="flex flex-col gap-6 flex-1 min-h-0">
      {/* Scan info + Status */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-6 text-sm text-muted-foreground">
          {/* Started at */}
          <div className="flex items-center gap-1.5">
            <Calendar className="h-4 w-4" />
            <span>{t("startedAt")}: {formatDate(startedAt, locale)}</span>
          </div>
          {/* Duration */}
          <div className="flex items-center gap-1.5">
            <Clock className="h-4 w-4" />
            <span>{t("duration")}: {formatDuration(startedAt, completedAt)}</span>
          </div>
          {/* Engine */}
          {scan.engineNames && scan.engineNames.length > 0 && (
            <div className="flex items-center gap-1.5">
              <Cpu className="h-4 w-4" />
              <span>{scan.engineNames.join(", ")}</span>
            </div>
          )}
          {/* Worker */}
          {scan.workerName && (
            <div className="flex items-center gap-1.5">
              <HardDrive className="h-4 w-4" />
              <span>{scan.workerName}</span>
            </div>
          )}
        </div>
        {/* Status badge */}
        <Badge variant="outline" className={statusStyle}>
          <StatusIcon className={`h-3.5 w-3.5 mr-1.5 ${statusIconConfig.animate ? 'animate-spin' : ''}`} />
          {tStatus(scan.status)}
        </Badge>
      </div>

      {/* Asset statistics cards */}
      <div>
        <h3 className="text-lg font-semibold mb-4">{t("assetsTitle")}</h3>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-6">
          {assetCards.map((card) => (
            <Link key={card.title} href={card.href} className="block">
              <div
                className="group relative p-4 hover:bg-accent/5 transition-all duration-300 cursor-pointer"
                style={{ background: "var(--card)" }}
              >
                <div className="absolute inset-0 border border-border/40 group-hover:border-primary/30 transition-colors" />
                <div className="absolute top-0 right-0 h-2 w-2 border-r border-t border-primary/50" />
                <div className="absolute bottom-0 left-0 h-2 w-2 border-l border-b border-primary/50" />

                <div className="relative z-10">
                  <div className="flex justify-between items-start mb-2">
                    <div className="text-[10px] font-mono text-muted-foreground bg-muted px-1.5 py-0.5 rounded-sm">
                      {card.code}
                    </div>
                    <card.icon className="h-4 w-4 text-muted-foreground/70 group-hover:text-primary transition-colors" />
                  </div>

                  <div className="text-3xl font-light tracking-tight text-foreground group-hover:translate-x-1 transition-transform duration-300">
                    {card.value.toLocaleString()}
                  </div>

                  <div className="mt-2 flex items-center gap-2">
                    <div className="h-px flex-1 bg-border border-t border-dashed border-muted-foreground/20" />
                    <span className="text-[11px] text-foreground/85 font-mono uppercase tracking-wider">
                      {card.title}
                    </span>
                  </div>
                </div>
              </div>
            </Link>
          ))}
        </div>
      </div>

      {/* Stage Progress + Logs - Left-Right Split Layout */}
      <div className="grid gap-4 md:grid-cols-[280px_1fr] flex-1 min-h-0">
        {/* Left Column: Stage Progress + Vulnerability Stats */}
        <div className="flex flex-col gap-4 min-h-0">
          {/* Stage Progress */}
          <Card className="flex-1 min-h-0">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
              <CardTitle className="text-sm font-medium">{t("stagesTitle")}</CardTitle>
              {scan.stageProgress && (
                <span className="text-xs text-muted-foreground">
                  {Object.values(scan.stageProgress).filter((p) => p.status === "completed").length}/
                  {Object.keys(scan.stageProgress).length} {t("stagesCompleted")}
                </span>
              )}
            </CardHeader>
            <CardContent className="pt-0 flex flex-col flex-1 min-h-0">
              {scan.stageProgress && Object.keys(scan.stageProgress).length > 0 ? (
                <div className="space-y-1 flex-1 min-h-0 overflow-y-auto pr-1">
                  {(Object.entries(scan.stageProgress) as Array<[string, StageProgressItem]>)
                    .toSorted(([, a], [, b]) => {
                      const priorityA = STAGE_STATUS_PRIORITY[a.status as StageStatus] ?? 99
                      const priorityB = STAGE_STATUS_PRIORITY[b.status as StageStatus] ?? 99
                      if (priorityA !== priorityB) {
                        return priorityA - priorityB
                      }
                      return (a.order ?? 0) - (b.order ?? 0)
                    })
                    .map(([stageName, stageProgress]) => {
                      const isRunning = stageProgress.status === "running"
                      return (
                        <div
                          key={stageName}
                          className={cn(
                            "flex items-center justify-between py-2 px-2 rounded-md transition-colors text-sm",
                            isRunning && "bg-[var(--warning)]/10 border border-[var(--warning)]/30",
                            stageProgress.status === "completed" && "text-muted-foreground",
                            stageProgress.status === "failed" && "bg-[var(--error)]/10 text-[var(--error)]",
                            stageProgress.status === "cancelled" && "text-muted-foreground",
                          )}
                        >
                          <div className="flex items-center gap-2 min-w-0">
                            <StageStatusIcon status={stageProgress.status} />
                            <span className={cn("truncate", isRunning && "font-medium text-foreground")}>
                              {tProgress(`stages.${stageName}`)}
                            </span>
                          </div>
                          <span className="text-xs text-muted-foreground font-mono shrink-0 ml-2">
                            {stageProgress.status === "completed" && stageProgress.duration
                              ? formatStageDuration(stageProgress.duration)
                              : stageProgress.status === "running"
                                ? tProgress("stage_running")
                                : stageProgress.status === "pending"
                                  ? "--"
                                  : ""}
                          </span>
                        </div>
                      )
                    })}
                </div>
              ) : (
                <div className="text-sm text-muted-foreground text-center py-4">
                  {t("noStages")}
                </div>
              )}
            </CardContent>
          </Card>

          {/* Vulnerability Stats - Compact */}
          <Link href={`/scan/history/${scanId}/vulnerabilities/`} className="block">
            <Card className="hover:border-primary/50 transition-colors cursor-pointer">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">{t("vulnerabilitiesTitle")}</CardTitle>
                <ChevronRight className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent className="pt-0">
                <div className="flex items-center gap-3 flex-wrap">
                  <div className="flex items-center gap-1.5">
                    <div className="w-2.5 h-2.5 rounded-full bg-[var(--error)]" />
                    <span className="text-sm font-medium">{vulnSummary.critical}</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <div className="w-2.5 h-2.5 rounded-full bg-[var(--error)]/70" />
                    <span className="text-sm font-medium">{vulnSummary.high}</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <div className="w-2.5 h-2.5 rounded-full bg-[var(--warning)]" />
                    <span className="text-sm font-medium">{vulnSummary.medium}</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <div className="w-2.5 h-2.5 rounded-full bg-[var(--info)]" />
                    <span className="text-sm font-medium">{vulnSummary.low}</span>
                  </div>
                  <span className="text-xs text-muted-foreground ml-auto">
                    {t("totalVulns", { count: vulnSummary.total ?? 0 })}
                  </span>
                </div>
              </CardContent>
            </Card>
          </Link>
        </div>

        {/* Right Column: Logs / Config */}
        <div className="flex flex-col min-h-0 rounded-lg overflow-hidden border">
          {/* Tab Header */}
          <div className="flex items-center justify-between px-3 py-2 bg-muted/30 border-b shrink-0">
            <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as 'logs' | 'config')}>
            <TabsList variant="minimal-tab">
              <TabsTrigger variant="minimal-tab" value="logs">
                {t("logsTitle")}
              </TabsTrigger>
              <TabsTrigger variant="minimal-tab" value="config">
                {t("configTitle")}
              </TabsTrigger>
            </TabsList>
            </Tabs>
            {/* Auto-refresh toggle (only for logs tab when running) */}
            {activeTab === 'logs' && isRunning && (
              <div className="flex items-center gap-2">
                <Switch
                  id="log-auto-refresh"
                  checked={autoRefresh}
                  onCheckedChange={setAutoRefresh}
                  className="scale-75"
                />
                <Label htmlFor="log-auto-refresh" className="text-xs cursor-pointer">
                  {t("autoRefresh")}
                </Label>
              </div>
            )}
          </div>
          
          {/* Tab Content */}
          <div className="flex-1 min-h-0">
            {activeTab === 'logs' ? (
              <ScanLogList logs={logs} loading={logsLoading} />
            ) : (
              <div className="h-full">
                {scan.yamlConfiguration ? (
                  <YamlEditor
                    value={scan.yamlConfiguration}
                    onChange={() => {}}
                    disabled={true}
                    className="h-full"
                  />
                ) : (
                  <div className="flex items-center justify-center h-full text-muted-foreground text-sm">
                    {t("noConfig")}
                  </div>
                )}
              </div>
            )}
          </div>
          
          {/* Bottom status bar (only for logs tab) */}
          {activeTab === 'logs' && (
            <div className="flex items-center px-4 py-2 bg-muted/50 border-t text-xs text-muted-foreground shrink-0">
              <span>{logs.length} 条记录</span>
              {isRunning && autoRefresh && (
                <>
                  <Separator orientation="vertical" className="h-3 mx-3" />
                  <span className="flex items-center gap-1.5">
                    <span className="size-1.5 rounded-full bg-[var(--success)] animate-pulse" />
                    每 3 秒刷新
                  </span>
                </>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
