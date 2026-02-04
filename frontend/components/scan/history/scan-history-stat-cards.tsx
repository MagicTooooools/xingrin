"use client"

import { useTranslations } from "next-intl"
import { Card, CardAction, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Skeleton } from "@/components/ui/skeleton"
import { 
  IconRadar, 
  IconPlayerPlay, 
  IconBug,
  IconStack2
} from "@/components/icons"
import { useScanStatistics } from "@/hooks/use-scans"

function StatCard({
  title,
  value,
  icon,
  loading,
  footer,
  badgeText,
}: {
  title: string
  value: string | number
  icon: React.ReactNode
  loading?: boolean
  footer: string
  badgeText: string
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
            {value}
          </CardTitle>
        )}
        <CardAction>
          <Badge variant="outline">{badgeText}</Badge>
        </CardAction>
      </CardHeader>
      <CardFooter className="flex-col items-start gap-1.5 text-sm">
        <div className="text-muted-foreground">{footer}</div>
      </CardFooter>
    </Card>
  )
}

export function ScanHistoryStatCards() {
  const { data, isLoading } = useScanStatistics()
  const t = useTranslations("scan.history.stats")

  return (
    <div className="*:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card grid grid-cols-1 gap-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:shadow-xs @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
      <StatCard
        title={t("totalScans")}
        value={data?.total ?? 0}
        icon={<IconRadar className="size-4" />}
        loading={isLoading}
        footer={t("allScanTasks")}
        badgeText={t("all")}
      />
      <StatCard
        title={t("running")}
        value={data?.running ?? 0}
        icon={<IconPlayerPlay className="size-4" />}
        loading={isLoading}
        footer={t("runningScans")}
        badgeText={t("all")}
      />
      <StatCard
        title={t("vulnsFound")}
        value={data?.totalVulns ?? 0}
        icon={<IconBug className="size-4" />}
        loading={isLoading}
        footer={t("completedScansFound")}
        badgeText={t("all")}
      />
      <StatCard
        title={t("assetsFound")}
        value={data?.totalAssets ?? 0}
        icon={<IconStack2 className="size-4" />}
        loading={isLoading}
        footer={t("assetTypes")}
        badgeText={t("all")}
      />
    </div>
  )
}
