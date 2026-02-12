"use client"

import React, { useMemo } from "react"
import dynamic from "next/dynamic"
import { useRouter } from "next/navigation"
import { useTranslations, useLocale } from "next-intl"
import { createScanHistoryColumns } from "./scan-history-columns"
import { getDateLocale } from "@/lib/date-utils"
import type { ScanRecord, ScanStatus } from "@/types/scan.types"
import type { ColumnDef } from "@tanstack/react-table"
import { DataTableSkeleton } from "@/components/ui/data-table-skeleton"
import { useScans } from "@/hooks/use-scans"
import { useSearchState } from "@/hooks/_shared/use-search-state"
import { buildPaginationInfo, normalizePagination } from "@/hooks/_shared/pagination"
import { ScanHistoryDialogs } from "@/components/scan/history/scan-history-list-dialogs"
import { useScanHistoryActions } from "@/components/scan/history/scan-history-list-state"

const ScanHistoryDataTable = dynamic(
  () => import("./scan-history-data-table").then((mod) => mod.ScanHistoryDataTable),
  {
    ssr: false,
    loading: () => <DataTableSkeleton rows={6} columns={6} withPadding />,
  }
)

const ScanProgressDialog = dynamic(
  () => import("@/components/scan/scan-progress-dialog").then((mod) => mod.ScanProgressDialog),
  { ssr: false }
)

/**
 * Scan history list component
 * Used to display and manage scan history records
 */
interface ScanHistoryListProps {
  hideToolbar?: boolean
  targetId?: number  // Filter by target ID
  pageSize?: number  // Custom page size
  hideTargetColumn?: boolean  // Hide target column (useful when showing scans for a specific target)
  pageSizeOptions?: number[]  // Custom page size options
  hidePagination?: boolean  // Hide pagination completely
}

export function ScanHistoryList({ hideToolbar = false, targetId, pageSize: customPageSize, hideTargetColumn = false, pageSizeOptions, hidePagination = false }: ScanHistoryListProps) {
  // Internationalization
  const tColumns = useTranslations("columns")
  const tCommon = useTranslations("common")
  const tTooltips = useTranslations("tooltips")
  const tScan = useTranslations("scan")
  const tToast = useTranslations("toast")
  const tConfirm = useTranslations("common.confirm")
  const locale = useLocale()

  // Build translation object
  const translations = useMemo(() => ({
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
      stop: tCommon("actions.stop"),
      stopScanPending: tScan("stopScanPending"),
      delete: tCommon("actions.delete"),
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
  
  // Pagination state
  const [pagination, setPagination] = React.useState({
    pageIndex: 0,
    pageSize: customPageSize || 10,
  })

  // Search state
  const [searchQuery, setSearchQuery] = React.useState("")
  
  // Status filter state
  const [statusFilter, setStatusFilter] = React.useState<ScanStatus | "all">("all")

  const handleStatusFilterChange = (status: ScanStatus | "all") => {
    setStatusFilter(status)
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
  }
  
  // Get scan list data
  const { data, isLoading, isFetching, error, refetch } = useScans({
    page: pagination.pageIndex + 1, // API page numbers start from 1
    pageSize: pagination.pageSize,
    search: searchQuery || undefined,
    target: targetId,
    status: statusFilter === "all" ? undefined : statusFilter,
  })
  const { isSearching, handleSearchChange } = useSearchState({
    isFetching,
    setSearchValue: setSearchQuery,
    onResetPage: () => setPagination((prev) => ({ ...prev, pageIndex: 0 })),
  })
  const paginationInfo = data
    ? buildPaginationInfo({
      ...normalizePagination(data, pagination.pageIndex + 1, pagination.pageSize),
      minTotalPages: 1,
    })
    : undefined
  
  // Scan list data
  const scans = data?.results || []
  
  const {
    selectedScans,
    setSelectedScans,
    deleteDialogOpen,
    setDeleteDialogOpen,
    scanToDelete,
    bulkDeleteDialogOpen,
    setBulkDeleteDialogOpen,
    stopDialogOpen,
    setStopDialogOpen,
    scanToStop,
    progressDialogOpen,
    setProgressDialogOpen,
    progressData,
    handleDeleteScan,
    confirmDelete,
    handleBulkDelete,
    confirmBulkDelete,
    handleStopScan,
    confirmStop,
    handleViewProgress,
  } = useScanHistoryActions({ tToast })

  // Helper function - format date
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

  // Navigation function
  const router = useRouter()
  const navigate = React.useCallback((path: string) => {
    router.push(path)
  }, [router])

  // Handle pagination change
  const handlePaginationChange = (newPagination: { pageIndex: number; pageSize: number }) => {
    setPagination(newPagination)
  }

  // Create column definitions
  const scanColumns = useMemo(
    () =>
      createScanHistoryColumns({
        formatDate,
        navigate,
        handleDelete: handleDeleteScan,
        handleStop: handleStopScan,
        handleViewProgress,
        statusClickable: false,
        t: translations,
        hideTargetColumn,
      }),
    [formatDate, navigate, translations, hideTargetColumn]
  )

  // Error handling
  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <p className="text-destructive mb-4">{tScan("history.loadFailed")}</p>
        <button
          onClick={() => {
            void refetch()
          }}
          className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
        >
          {tScan("history.retry")}
        </button>
      </div>
    )
  }

  // Loading state
  if (isLoading) {
    return (
      <DataTableSkeleton
        toolbarButtonCount={2}
        rows={6}
        columns={6}
        withPadding={false}
      />
    )
  }

  return (
    <>
      <ScanHistoryDataTable
        data={scans}
        columns={scanColumns as ColumnDef<ScanRecord>[]}
        onBulkDelete={hideToolbar ? undefined : handleBulkDelete}
        onSelectionChange={setSelectedScans}
        searchPlaceholder={tScan("history.searchPlaceholder")}
        searchValue={searchQuery}
        onSearch={handleSearchChange}
        isSearching={isSearching}
        pagination={pagination}
        setPagination={setPagination}
        paginationInfo={paginationInfo}
        onPaginationChange={handlePaginationChange}
        hideToolbar={hideToolbar}
        pageSizeOptions={pageSizeOptions}
        hidePagination={hidePagination}
        statusFilter={statusFilter}
        onStatusFilterChange={handleStatusFilterChange}
      />

      <ScanHistoryDialogs
        tConfirm={tConfirm}
        tCommon={tCommon}
        deleteDialogOpen={deleteDialogOpen}
        setDeleteDialogOpen={setDeleteDialogOpen}
        scanToDelete={scanToDelete}
        onConfirmDelete={confirmDelete}
        bulkDeleteDialogOpen={bulkDeleteDialogOpen}
        setBulkDeleteDialogOpen={setBulkDeleteDialogOpen}
        selectedScans={selectedScans}
        onConfirmBulkDelete={confirmBulkDelete}
        stopDialogOpen={stopDialogOpen}
        setStopDialogOpen={setStopDialogOpen}
        scanToStop={scanToStop}
        onConfirmStop={confirmStop}
      />

      {/* Scan progress dialog */}
      {progressDialogOpen ? (
        <ScanProgressDialog
          open={progressDialogOpen}
          onOpenChange={setProgressDialogOpen}
          data={progressData}
        />
      ) : null}
    </>
  )
}
