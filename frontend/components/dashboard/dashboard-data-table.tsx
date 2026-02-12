"use client"

import * as React from "react"
import { IconBug, IconRadar } from "@/components/icons"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Skeleton } from "@/components/ui/skeleton"
import { UnifiedDataTable } from "@/components/ui/data-table/unified-data-table"
import { useAllVulnerabilities } from "@/hooks/use-vulnerabilities"
import { useScans } from "@/hooks/use-scans"
import { useResourceMutation } from "@/hooks/_shared/create-resource-mutation"
import { VulnerabilityDetailDialog } from "@/components/vulnerabilities/vulnerability-detail-dialog"
import { createVulnerabilityColumns } from "@/components/vulnerabilities/vulnerabilities-columns"
import { createScanHistoryColumns } from "@/components/scan/history/scan-history-columns"
import { ScanProgressDialog, buildScanProgressData, type ScanProgressData } from "@/components/scan/scan-progress-dialog"
import { getScan } from "@/services/scan.service"
import { useRouter } from "next/navigation"
import { toast } from "sonner"
import { deleteScan, stopScan } from "@/services/scan.service"
import { useTranslations, useLocale } from "next-intl"
import { getDateLocale } from "@/lib/date-utils"
import { buildPaginationInfo, normalizePagination } from "@/hooks/_shared/pagination"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import type { Vulnerability } from "@/types/vulnerability.types"
import type { ScanRecord } from "@/types/scan.types"
import type { PaginationInfo } from "@/types/common.types"

export function DashboardDataTable() {
  const router = useRouter()
  const t = useTranslations()
  const locale = useLocale()
  const [activeTab, setActiveTab] = React.useState("scans")
  
  // 漏洞详情弹窗
  const [selectedVuln, setSelectedVuln] = React.useState<Vulnerability | null>(null)
  const [vulnDialogOpen, setVulnDialogOpen] = React.useState(false)
  
  // 扫描进度弹窗
  const [progressData, setProgressData] = React.useState<ScanProgressData | null>(null)
  const [progressDialogOpen, setProgressDialogOpen] = React.useState(false)
  
  // 删除确认弹窗
  const [deleteDialogOpen, setDeleteDialogOpen] = React.useState(false)
  const [scanToDelete, setScanToDelete] = React.useState<ScanRecord | null>(null)
  
  // 停止确认弹窗
  const [stopDialogOpen, setStopDialogOpen] = React.useState(false)
  const [scanToStop, setScanToStop] = React.useState<ScanRecord | null>(null)
  
  // 分页状态
  const [vulnPagination, setVulnPagination] = React.useState({ pageIndex: 0, pageSize: 10 })
  const [scanPagination, setScanPagination] = React.useState({ pageIndex: 0, pageSize: 10 })

  // 获取漏洞数据
  const vulnQuery = useAllVulnerabilities({
    page: vulnPagination.pageIndex + 1,
    pageSize: vulnPagination.pageSize,
  })
  
  // 获取扫描数据
  const scanQuery = useScans({
    page: scanPagination.pageIndex + 1,
    pageSize: scanPagination.pageSize,
  })

  // 删除扫描的 mutation
  const deleteMutation = useResourceMutation({
    mutationFn: deleteScan,
    invalidate: [{ queryKey: ['scans'] }],
    skipDefaultErrorHandler: true,
  })

  // 停止扫描的 mutation
  const stopMutation = useResourceMutation({
    mutationFn: stopScan,
    invalidate: [{ queryKey: ['scans'] }],
    skipDefaultErrorHandler: true,
  })

  const vulnerabilities = vulnQuery.data?.vulnerabilities ?? []
  const scans = scanQuery.data?.results ?? []

  // 格式化日期
  const formatDate = React.useCallback((dateString: string): string => {
    return new Date(dateString).toLocaleString(getDateLocale(locale), {
      year: "numeric",
      month: "numeric",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    })
  }, [locale])

  // 点击漏洞行
  const handleVulnRowClick = React.useCallback((vuln: Vulnerability) => {
    setSelectedVuln(vuln)
    setVulnDialogOpen(true)
  }, [])

  // 漏洞列定义
  const vulnColumns = React.useMemo(
    () => createVulnerabilityColumns({
      formatDate,
      handleViewDetail: handleVulnRowClick,
      t: {
        columns: {
          severity: t('columns.vulnerability.severity'),
          source: t('columns.vulnerability.source'),
          vulnType: t('columns.vulnerability.vulnType'),
          url: t('columns.common.url'),
          createdAt: t('columns.common.createdAt'),
        },
        actions: {
          details: t('common.actions.details'),
          selectAll: t('common.actions.selectAll'),
          selectRow: t('common.actions.selectRow'),
        },
        tooltips: {
          vulnDetails: t('tooltips.vulnDetails'),
          reviewed: t('tooltips.reviewed'),
          pending: t('tooltips.pending'),
        },
        severity: {
          critical: t('severity.critical'),
          high: t('severity.high'),
          medium: t('severity.medium'),
          low: t('severity.low'),
          info: t('severity.info'),
        },
      },
    }),
    [formatDate, handleVulnRowClick, t]
  )

  // 扫描进度查看
  const handleViewProgress = React.useCallback(async (scan: ScanRecord) => {
    try {
      const fullScan = await getScan(scan.id)
      const data = buildScanProgressData(fullScan)
      setProgressData(data)
      setProgressDialogOpen(true)
    } catch (error) {
      void error
    }
  }, [])

  // 处理删除扫描
  const handleDelete = React.useCallback((scan: ScanRecord) => {
    setScanToDelete(scan)
    setDeleteDialogOpen(true)
  }, [])

  // 确认删除
  const confirmDelete = async () => {
    if (!scanToDelete) return
    setDeleteDialogOpen(false)
    try {
      await deleteMutation.mutateAsync(scanToDelete.id)
      toast.success(t('common.status.success'))
    } catch {
      toast.error(t('common.status.error'))
    } finally {
      setScanToDelete(null)
    }
  }

  // 处理停止扫描
  const handleStop = React.useCallback((scan: ScanRecord) => {
    setScanToStop(scan)
    setStopDialogOpen(true)
  }, [])

  // 确认停止
  const confirmStop = async () => {
    if (!scanToStop) return
    setStopDialogOpen(false)
    try {
      await stopMutation.mutateAsync(scanToStop.id)
      toast.success(t('common.status.success'))
    } catch {
      toast.error(t('common.status.error'))
    } finally {
      setScanToStop(null)
    }
  }

  // 扫描列定义
  const scanColumns = React.useMemo(
    () => createScanHistoryColumns({
      formatDate,
      navigate: (path: string) => router.push(path),
      handleDelete,
      handleStop,
      handleViewProgress,
      t: {
        columns: {
          target: t('columns.scanHistory.target'),
          summary: t('columns.scanHistory.summary'),
          engineName: t('columns.scanHistory.engineName'),
          workerName: t('columns.scanHistory.workerName'),
          createdAt: t('columns.common.createdAt'),
          status: t('columns.common.status'),
          progress: t('columns.scanHistory.progress'),
        },
        actions: {
          snapshot: t('common.actions.snapshot'),
          stop: t('scan.stopScan'),
          stopScanPending: t('scan.stopScanPending'),
          delete: t('common.actions.delete'),
          selectAll: t('common.actions.selectAll'),
          selectRow: t('common.actions.selectRow'),
        },
        tooltips: {
          targetDetails: t('tooltips.targetDetails'),
          viewProgress: t('tooltips.viewProgress'),
        },
        status: {
          cancelled: t('common.status.cancelled'),
          completed: t('common.status.completed'),
          failed: t('common.status.failed'),
          pending: t('common.status.pending'),
          running: t('common.status.running'),
        },
        summary: {
          subdomains: t('columns.scanHistory.subdomains'),
          websites: t('columns.scanHistory.websites'),
          ipAddresses: t('columns.scanHistory.ipAddresses'),
          endpoints: t('columns.scanHistory.endpoints'),
          vulnerabilities: t('columns.scanHistory.vulnerabilities'),
        },
      },
    }),
    [formatDate, router, handleViewProgress, handleDelete, handleStop, t]
  )

  const vulnPaginationInfo: PaginationInfo = buildPaginationInfo({
    ...normalizePagination(
      vulnQuery.data?.pagination,
      vulnPagination.pageIndex + 1,
      vulnPagination.pageSize
    ),
    minTotalPages: 1,
  })

  const scanPaginationInfo: PaginationInfo = buildPaginationInfo({
    ...normalizePagination(
      scanQuery.data,
      scanPagination.pageIndex + 1,
      scanPagination.pageSize
    ),
    minTotalPages: 1,
  })

  const tabsToolbar = (
    <TabsList>
      <TabsTrigger value="scans" className="gap-1.5">
        <IconRadar className="h-4 w-4" />
        {t('navigation.scanHistory')}
      </TabsTrigger>
      <TabsTrigger value="vulnerabilities" className="gap-1.5">
        <IconBug className="h-4 w-4" />
        {t('navigation.vulnerabilities')}
      </TabsTrigger>
    </TabsList>
  )

  return (
    <>
      <VulnerabilityDetailDialog
        vulnerability={selectedVuln}
        open={vulnDialogOpen}
        onOpenChange={setVulnDialogOpen}
      />
      {progressData && (
        <ScanProgressDialog
          open={progressDialogOpen}
          onOpenChange={setProgressDialogOpen}
          data={progressData}
        />
      )}

      {/* 删除确认对话框 */}
      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t('common.confirm.deleteTitle')}</AlertDialogTitle>
            <AlertDialogDescription>
              {t('common.confirm.deleteMessage')}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t('common.actions.cancel')}</AlertDialogCancel>
            <AlertDialogAction 
              onClick={confirmDelete} 
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {t('common.actions.delete')}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* 停止扫描确认对话框 */}
      <AlertDialog open={stopDialogOpen} onOpenChange={setStopDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{t('common.confirm.title')}</AlertDialogTitle>
            <AlertDialogDescription>
              {t('common.confirm.deleteMessage')}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{t('common.actions.cancel')}</AlertDialogCancel>
            <AlertDialogAction 
              onClick={confirmStop} 
              className="bg-primary text-primary-foreground hover:bg-primary/90"
            >
              {t('scan.stopScan')}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
      
      <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
        {/* 漏洞表格 */}
        <TabsContent value="vulnerabilities" className="mt-0">
          {vulnQuery.isLoading ? (
            <div className="space-y-2">
              {[...Array(5)].map((_, i) => <Skeleton key={i} className="h-12 w-full" />)}
            </div>
          ) : (
            <UnifiedDataTable
              data={vulnerabilities}
              columns={vulnColumns}
              getRowId={(row) => String(row.id)}
              state={{
                pagination: vulnPagination,
                onPaginationChange: setVulnPagination,
                paginationInfo: vulnPaginationInfo,
              }}
              behavior={{
                enableRowSelection: false,
              }}
              actions={{
                showAddButton: false,
                showBulkDelete: false,
              }}
              ui={{
                emptyMessage: t('common.status.noData'),
                toolbarLeft: tabsToolbar,
              }}
            />
          )}
        </TabsContent>

        {/* 扫描历史表格 */}
        <TabsContent value="scans" className="mt-0">
          {scanQuery.isLoading ? (
            <div className="space-y-2">
              {[...Array(5)].map((_, i) => <Skeleton key={i} className="h-12 w-full" />)}
            </div>
          ) : (
            <UnifiedDataTable
              data={scans}
              columns={scanColumns}
              getRowId={(row) => String(row.id)}
              state={{
                pagination: scanPagination,
                onPaginationChange: setScanPagination,
                paginationInfo: scanPaginationInfo,
              }}
              behavior={{
                enableRowSelection: false,
                enableAutoColumnSizing: true,
              }}
              actions={{
                showAddButton: false,
                showBulkDelete: false,
              }}
              ui={{
                emptyMessage: t('common.status.noData'),
                toolbarLeft: tabsToolbar,
              }}
            />
          )}
        </TabsContent>
      </Tabs>
    </>
  )
}
