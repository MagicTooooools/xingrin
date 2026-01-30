"use client"

import * as React from "react"
import { useRouter } from "next/navigation"
import { useTranslations, useLocale } from "next-intl"
import { ScanHistoryDataTable } from "@/components/scan/history/scan-history-data-table"
import { createScanHistoryColumns } from "@/components/scan/history/scan-history-columns"
import { useScans } from "@/hooks/use-scans"
import { DataTableSkeleton } from "@/components/ui/data-table-skeleton"
import { getDateLocale } from "@/lib/date-utils"
import type { ScanRecord } from "@/types/scan.types"
import type { ColumnDef } from "@tanstack/react-table"

export function DashboardScanHistory() {
  const [pagination, setPagination] = React.useState({ pageIndex: 0, pageSize: 5 })
  const router = useRouter()
  const locale = useLocale()

  // 国际化
  const tColumns = useTranslations("columns")
  const tCommon = useTranslations("common")
  const tTooltips = useTranslations("tooltips")
  const tScan = useTranslations("scan")

  // 构建翻译对象
  const translations = React.useMemo(() => ({
    columns: {
      target: tColumns("scanHistory.target"),
      summary: tColumns("scanHistory.summary"),
      engineName: tColumns("scanHistory.engineName"),
      workerName: tColumns("scanHistory.workerName"),
      createdAt: tColumns("common.createdAt"),
      status: tColumns("common.status"),
      progress: tColumns("scanHistory.progress"),
    },
    actions: {
      snapshot: tCommon("actions.snapshot"),
      stopScan: tScan("stopScan"),
      delete: tCommon("actions.delete"),
      openMenu: tCommon("actions.openMenu"),
      selectAll: tCommon("actions.selectAll"),
      selectRow: tCommon("actions.selectRow"),
    },
    tooltips: {
      targetDetails: tTooltips("targetDetails"),
      viewProgress: tTooltips("viewProgress"),
    },
    status: {
      cancelled: tCommon("status.cancelled"),
      completed: tCommon("status.completed"),
      failed: tCommon("status.failed"),
      pending: tCommon("status.pending"),
      running: tCommon("status.running"),
    },
    summary: {
      subdomains: tColumns("scanHistory.subdomains"),
      websites: tColumns("scanHistory.websites"),
      ipAddresses: tColumns("scanHistory.ipAddresses"),
      endpoints: tColumns("scanHistory.endpoints"),
      vulnerabilities: tColumns("scanHistory.vulnerabilities"),
    },
  }), [tColumns, tCommon, tTooltips, tScan])

  const { data, isLoading } = useScans({
    page: pagination.pageIndex + 1,
    pageSize: pagination.pageSize,
    status: 'running',
  })

  const formatDate = React.useCallback((dateString: string) => new Date(dateString).toLocaleString(getDateLocale(locale), { hour12: false }), [locale])
  const navigate = React.useCallback((path: string) => router.push(path), [router])
  const handleDelete = React.useCallback(() => {}, [])
  const handleStop = React.useCallback(() => {
    // 仪表盘列表暂时不提供停止逻辑，实现时可在此调用对应的停止扫描接口
  }, [])

  const columns = React.useMemo(
    () => createScanHistoryColumns({ formatDate, navigate, handleDelete, handleStop, t: translations }) as ColumnDef<ScanRecord>[],
    [formatDate, navigate, handleDelete, handleStop, translations]
  )

  if (isLoading && !data) {
    return (
      <DataTableSkeleton
        withPadding={false}
        toolbarButtonCount={2}
        rows={4}
        columns={3}
      />
    )
  }

  const paginationInfo = data
    ? { total: data.total, page: data.page, pageSize: data.pageSize, totalPages: data.totalPages }
    : undefined

  return (
    <ScanHistoryDataTable
      data={data?.results ?? []}
      columns={columns}
      hideToolbar
      hidePagination
      pagination={pagination}
      setPagination={setPagination}
      paginationInfo={paginationInfo}
      onPaginationChange={setPagination}
    />
  )
}
