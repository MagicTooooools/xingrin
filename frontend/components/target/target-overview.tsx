"use client"

import React, { useState } from "react"
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
  ShieldAlert,
  AlertTriangle,
  Clock,
  Calendar,
  ChevronRight,
  CheckCircle2,
  PauseCircle,
  Loader2,
} from "@/components/icons"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { Button } from "@/components/ui/button"
import { useTarget } from "@/hooks/use-targets"
import { useScheduledScans } from "@/hooks/use-scheduled-scans"
import { ScanHistoryList } from "@/components/scan/history/scan-history-list"
import type { TargetDetail } from "@/types/target.types"

// Dynamic import for InitiateScanDialog (only loaded when dialog is opened)
const InitiateScanDialog = dynamic(() => import('@/components/scan/initiate-scan-dialog').then(m => ({ default: m.InitiateScanDialog })), {
  ssr: false
})

interface TargetOverviewProps {
  targetId: number
}

/**
 * Format date helper function
 */
function formatDate(dateString: string | undefined, locale: string): string {
  if (!dateString) return "-"
  return new Date(dateString).toLocaleString(locale === 'zh' ? 'zh-CN' : 'en-US', {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  })
}

/**
 * Format short date for scheduled scans
 */
function formatShortDate(
  dateString: string | undefined,
  locale: string,
  todayText: string,
  tomorrowText: string
): string {
  if (!dateString) return "-"
  const date = new Date(dateString)
  const now = new Date()
  const tomorrow = new Date(now)
  tomorrow.setDate(tomorrow.getDate() + 1)

  const localeStr = locale === 'zh' ? 'zh-CN' : 'en-US'

  // Check if it's today
  if (date.toDateString() === now.toDateString()) {
    return todayText + " " + date.toLocaleTimeString(localeStr, {
      hour: "2-digit",
      minute: "2-digit",
    })
  }
  // Check if it's tomorrow
  if (date.toDateString() === tomorrow.toDateString()) {
    return tomorrowText + " " + date.toLocaleTimeString(localeStr, {
      hour: "2-digit",
      minute: "2-digit",
    })
  }
  // Otherwise show date
  return date.toLocaleString(localeStr, {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  })
}

/**
 * Target overview component
 * Displays statistics cards for the target
 */
export function TargetOverview({ targetId }: TargetOverviewProps) {
  const t = useTranslations("pages.targetDetail.overview")
  const tInitiate = useTranslations("scan.initiate")
  const locale = useLocale()

  const [scanDialogOpen, setScanDialogOpen] = useState(false)

  const { data: target, isLoading, error } = useTarget(targetId)
  const { data: scheduledScansData, isLoading: isLoadingScans } = useScheduledScans({
    targetId,
    pageSize: 5
  })

  const targetSummary = (target as TargetDetail | undefined)?.summary

  // Memoize summary and vulnerability data to avoid unnecessary recalculations
  const summary = React.useMemo(() => ({
    subdomains: targetSummary?.subdomains ?? 0,
    websites: targetSummary?.websites ?? 0,
    endpoints: targetSummary?.endpoints ?? 0,
    ips: targetSummary?.ips ?? 0,
    directories: targetSummary?.directories ?? 0,
    screenshots: targetSummary?.screenshots ?? 0,
  }), [targetSummary])

  const vulnSummary = React.useMemo(
    () => targetSummary?.vulnerabilities || { total: 0, critical: 0, high: 0, medium: 0, low: 0 },
    [targetSummary]
  )

  // Memoize asset cards array to avoid recreation on every render
  const assetCards = React.useMemo(
    () => [
      {
        title: t("cards.websites"),
        value: summary.websites || 0,
        icon: Globe,
        code: "DAT-WEB",
        href: `/target/${targetId}/websites/`,
      },
      {
        title: t("cards.subdomains"),
        value: summary.subdomains || 0,
        icon: Network,
        code: "DAT-SUB",
        href: `/target/${targetId}/subdomain/`,
      },
      {
        title: t("cards.ips"),
        value: summary.ips || 0,
        icon: Server,
        code: "DAT-IP",
        href: `/target/${targetId}/ip-addresses/`,
      },
      {
        title: t("cards.urls"),
        value: summary.endpoints || 0,
        icon: Link2,
        code: "DAT-URL",
        href: `/target/${targetId}/endpoints/`,
      },
      {
        title: t("cards.directories"),
        value: summary.directories || 0,
        icon: FolderOpen,
        code: "DAT-DIR",
        href: `/target/${targetId}/directories/`,
      },
      {
        title: t("cards.screenshots"),
        value: summary.screenshots || 0,
        icon: Camera,
        code: "DAT-SCR",
        href: `/target/${targetId}/screenshots/`,
      },
    ],
    [summary, targetId, t]
  )

  // Memoize derived values to avoid unnecessary recalculations
  const scheduledScans = React.useMemo(() => scheduledScansData?.results || [], [scheduledScansData?.results])
  const totalScheduledScans = React.useMemo(() => scheduledScansData?.total || 0, [scheduledScansData?.total])
  const enabledScans = React.useMemo(() => scheduledScans.filter(s => s.isEnabled), [scheduledScans])

  // Get next execution time from enabled scans (memoized)
  const nextExecution = React.useMemo(() => {
    const enabledWithNextRun = enabledScans.filter(s => s.nextRunTime)
    if (enabledWithNextRun.length === 0) return null

    const sorted = enabledWithNextRun.toSorted((a, b) =>
      new Date(a.nextRunTime!).getTime() - new Date(b.nextRunTime!).getTime()
    )
    return sorted[0]
  }, [enabledScans])

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

  if (error || !target) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <AlertTriangle className="h-10 w-10 text-destructive mb-4" />
        <p className="text-muted-foreground">{t("loadError")}</p>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Target info + Initiate Scan button */}
      <div className="flex flex-wrap items-end justify-between gap-3 -mt-4">
        <div className="flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-muted-foreground">
          <div className="flex items-center gap-1.5">
            <Calendar className="h-4 w-4" />
            <span>{t("createdAt")}: {formatDate(target.createdAt, locale)}</span>
          </div>
          <div className="flex items-center gap-1.5">
            <Clock className="h-4 w-4" />
            <span>{t("lastScanned")}: {formatDate(target.lastScannedAt, locale)}</span>
          </div>
        </div>
        <Button 
          onClick={() => setScanDialogOpen(true)}
          className="
            relative overflow-hidden group 
            bg-background border border-destructive/50 text-destructive
            hover:border-destructive
            min-w-[140px]
            shrink-0 self-end
          "
        >
          <span className="absolute inset-y-0 left-0 w-[2px] bg-destructive transition-all duration-300 group-hover:w-full opacity-10" />
          <span className="absolute bottom-0 right-0 h-[2px] w-0 bg-destructive transition-all duration-300 group-hover:w-full" />
          <span className="relative z-10 flex items-center gap-2">
            <Loader2 className="h-4 w-4 animate-spin hidden group-hover:block" />
            <span className="group-hover:hidden">{t("initiateScan")}</span>
            <span className="hidden group-hover:block">{tInitiate("initiating")}</span>
          </span>
        </Button>
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

      {/* Scheduled Scans + Vulnerability Statistics (Two columns) */}
      <div className="grid gap-4 md:grid-cols-2">
        {/* Scheduled Scans Card */}
        <Card className="flex flex-col">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <div className="flex items-center gap-2">
              <Clock className="h-4 w-4 text-muted-foreground" />
              <CardTitle className="text-sm font-medium">{t("scheduledScans.title")}</CardTitle>
            </div>
            <Link href={`/target/${targetId}/settings/`}>
              <Button variant="ghost" size="sm" className="h-7 text-xs">
                {t("scheduledScans.manage")}
                <ChevronRight className="h-3 w-3 ml-1" />
              </Button>
            </Link>
          </CardHeader>
          <CardContent className="flex-1 flex flex-col">
            {isLoadingScans ? (
              <div className="space-y-2">
                <Skeleton className="h-4 w-32" />
                <Skeleton className="h-4 w-48" />
              </div>
            ) : totalScheduledScans === 0 ? (
              <div className="flex-1 flex flex-col items-center justify-center">
                <Clock className="h-8 w-8 text-muted-foreground/50 mb-2" />
                <p className="text-sm text-muted-foreground">{t("scheduledScans.empty")}</p>
                <Link href={`/target/${targetId}/settings/`}>
                  <Button variant="link" size="sm" className="mt-1">
                    {t("scheduledScans.createFirst")}
                  </Button>
                </Link>
              </div>
            ) : (
              <div className="space-y-3">
                {/* Stats row */}
                <div className="flex items-center gap-4 text-sm">
                  <div>
                    <span className="text-muted-foreground">{t("scheduledScans.configured")}: </span>
                    <span className="font-medium">{totalScheduledScans}</span>
                  </div>
                  <div>
                    <span className="text-muted-foreground">{t("scheduledScans.enabled")}: </span>
                    <span className="font-medium text-green-600">{enabledScans.length}</span>
                  </div>
                </div>
                
                {/* Next execution */}
                {nextExecution && (
                  <div className="text-sm">
                    <span className="text-muted-foreground">{t("scheduledScans.nextRun")}: </span>
                    <span className="font-medium">
                      {formatShortDate(
                        nextExecution.nextRunTime,
                        locale,
                        t("scheduledScans.today"),
                        t("scheduledScans.tomorrow")
                      )}
                    </span>
                  </div>
                )}

                {/* Task list - max 2 items */}
                <div className="space-y-2 pt-2 border-t">
                  {scheduledScans.slice(0, 2).map((scan) => (
                    <div key={scan.id} className="flex items-center gap-2 text-sm">
                      {scan.isEnabled ? (
                        <CheckCircle2 className="h-3.5 w-3.5 text-green-500 shrink-0" />
                      ) : (
                        <PauseCircle className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
                      )}
                      <span className={`truncate ${!scan.isEnabled ? 'text-muted-foreground' : ''}`}>
                        {scan.name}
                      </span>
                    </div>
                  ))}
                  {totalScheduledScans > 2 && (
                    <p className="text-xs text-muted-foreground">
                      {t("scheduledScans.more", { count: totalScheduledScans - 2 })}
                    </p>
                  )}
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Vulnerability Statistics Card */}
        <Link href={`/target/${targetId}/vulnerabilities/`} className="block">
          <Card className="h-full hover:border-primary/50 transition-colors cursor-pointer flex flex-col">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
              <div className="flex items-center gap-2">
              <ShieldAlert className="h-4 w-4 text-muted-foreground" />
                <CardTitle className="text-sm font-medium">{t("vulnerabilitiesTitle")}</CardTitle>
              </div>
              <Button variant="ghost" size="sm" className="h-7 text-xs">
                {t("viewAll")}
                <ChevronRight className="h-3 w-3 ml-1" />
              </Button>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Total count */}
              <div className="flex items-baseline gap-2">
                <span className="text-3xl font-bold">{vulnSummary.total}</span>
                <span className="text-sm text-muted-foreground">{t("cards.vulnerabilities")}</span>
              </div>

              {/* Severity breakdown */}
              <div className="grid grid-cols-2 gap-3">
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 rounded-full bg-red-500" />
                  <span className="text-sm text-muted-foreground">{t("severity.critical")}</span>
                  <span className="text-sm font-medium ml-auto">{vulnSummary.critical}</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 rounded-full bg-orange-500" />
                  <span className="text-sm text-muted-foreground">{t("severity.high")}</span>
                  <span className="text-sm font-medium ml-auto">{vulnSummary.high}</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 rounded-full bg-yellow-500" />
                  <span className="text-sm text-muted-foreground">{t("severity.medium")}</span>
                  <span className="text-sm font-medium ml-auto">{vulnSummary.medium}</span>
                </div>
                <div className="flex items-center gap-2">
                  <div className="w-3 h-3 rounded-full bg-blue-500" />
                  <span className="text-sm text-muted-foreground">{t("severity.low")}</span>
                  <span className="text-sm font-medium ml-auto">{vulnSummary.low}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </Link>
      </div>

      {/* Scan history */}
      <div>
        <h3 className="text-lg font-semibold mb-4">{t("scanHistoryTitle")}</h3>
        <ScanHistoryList targetId={targetId} hideToolbar pageSize={5} hideTargetColumn pageSizeOptions={[5, 10, 20, 50, 100]} />
      </div>

      {/* Initiate Scan Dialog */}
      <InitiateScanDialog
        open={scanDialogOpen}
        onOpenChange={setScanDialogOpen}
        targetId={targetId}
        targetName={target.name}
      />
    </div>
  )
}
