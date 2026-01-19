"use client"

import { memo } from "react"
import { useAssetStatistics } from "@/hooks/use-dashboard"
import { Card, CardAction, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { IconTarget, IconStack2, IconBug, IconPlayerPlay, IconTrendingUp, IconTrendingDown } from "@tabler/icons-react"
import { useTranslations } from "next-intl"
import { useLocale } from "next-intl"

const TrendBadge = memo(function TrendBadge({ change }: { change: number }) {
  if (change === 0) return null
  
  const isPositive = change > 0
  return (
    <Badge 
      variant="outline" 
      className={isPositive 
        ? "text-[#238636] dark:text-[#3fb950] border-[#238636]/20 bg-[#238636]/10" 
        : "text-[#da3633] dark:text-[#f85149] border-[#da3633]/20 bg-[#da3633]/10"
      }
    >
      {isPositive ? <IconTrendingUp className="size-3 mr-1" /> : <IconTrendingDown className="size-3 mr-1" />}
      {isPositive ? '+' : ''}{change}
    </Badge>
  )
})

const StatCard = memo(function StatCard({
  title,
  value,
  change,
  icon,
  footer,
  loading,
}: {
  title: string
  value: string | number
  change?: number
  icon: React.ReactNode
  footer: string
  loading?: boolean
}) {
  return (
    <Card className="@container/card">
      <CardHeader>
        <CardDescription className="flex items-center gap-2">
          {icon}
          {title}
        </CardDescription>
        {loading ? (
          <Skeleton className="h-8 w-24" />
        ) : (
          <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
            {typeof value === 'number' ? value.toLocaleString() : value}
          </CardTitle>
        )}
        {!loading && change !== undefined && (
          <CardAction>
            <TrendBadge change={change} />
          </CardAction>
        )}
      </CardHeader>
      <CardFooter className="flex-col items-start gap-1.5 text-sm">
        <div className="text-muted-foreground">{footer}</div>
      </CardFooter>
    </Card>
  )
})

function formatUpdateTime(dateStr: string | null, locale: string, noDataText: string) {
  if (!dateStr) return noDataText
  const date = new Date(dateStr)
  return date.toLocaleString(locale === 'zh' ? 'zh-CN' : 'en-US', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function DashboardStatCards() {
  const { data, isLoading } = useAssetStatistics()
  const t = useTranslations("dashboard.statCards")
  const locale = useLocale()

  return (
    <div className="flex flex-col gap-2 px-4 lg:px-6">
      <div className="grid grid-cols-1 gap-4 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
        <StatCard
          title={t("assetsFound")}
          value={data?.totalAssets ?? 0}
          change={data?.changeAssets}
          icon={<IconStack2 className="size-4" />}
          loading={isLoading}
          footer={t("assetsFooter")}
        />
        <StatCard
          title={t("vulnsFound")}
          value={data?.totalVulns ?? 0}
          change={data?.changeVulns}
          icon={<IconBug className="size-4" />}
          loading={isLoading}
          footer={t("vulnsFooter")}
        />
        <StatCard
          title={t("monitoredTargets")}
          value={data?.totalTargets ?? 0}
          change={data?.changeTargets}
          icon={<IconTarget className="size-4" />}
          loading={isLoading}
          footer={t("targetsFooter")}
        />
        <StatCard
          title={t("runningScans")}
          value={data?.runningScans ?? 0}
          icon={<IconPlayerPlay className="size-4" />}
          loading={isLoading}
          footer={t("scansFooter")}
        />
      </div>
      <div className="flex items-center gap-3 mt-1 -mb-2 text-xs text-muted-foreground">
        <div className="flex-1 border-t" />
        <span>{t("updatedAt", { time: formatUpdateTime(data?.updatedAt ?? null, locale, t("noData")) })}</span>
      </div>
    </div>
  )
}
