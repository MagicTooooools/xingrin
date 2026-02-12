"use client"

import React, { useCallback, useMemo, useState } from "react"
import { AlertTriangle } from "@/components/icons"
import { useTranslations, useLocale } from "next-intl"
import { DirectoriesDataTable } from "./directories-data-table"
import { createDirectoryColumns } from "./directories-columns"
import { DataTableSkeleton } from "@/components/ui/data-table-skeleton"
import { Button } from "@/components/ui/button"
import { useTargetDirectories, useScanDirectories } from "@/hooks/use-directories"
import { useTarget } from "@/hooks/use-targets"
import { useSearchState } from "@/hooks/_shared/use-search-state"
import { buildPaginationInfo, normalizePagination } from "@/hooks/_shared/pagination"
import { DirectoryService } from "@/services/directory.service"
import { BulkAddUrlsDialog } from "@/components/common/bulk-add-urls-dialog"
import { ConfirmDialog } from "@/components/ui/confirm-dialog"
import { getDateLocale } from "@/lib/date-utils"
import { escapeCSV, formatDateForCSV } from "@/lib/csv-utils"
import { downloadBlob } from "@/lib/download-utils"
import type { TargetType } from "@/lib/url-validator"
import type { Directory } from "@/types/directory.types"
import { toast } from "sonner"

export function DirectoriesView({
  targetId,
  scanId,
}: {
  targetId?: number
  scanId?: number
}) {
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  })
  const [selectedDirectories, setSelectedDirectories] = useState<Directory[]>([])
  const [bulkAddDialogOpen, setBulkAddDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [isDeleting, setIsDeleting] = useState(false)

  const [filterQuery, setFilterQuery] = useState("")

  // Internationalization
  const tColumns = useTranslations("columns")
  const tCommon = useTranslations("common")
  const tToast = useTranslations("toast")
  const tStatus = useTranslations("common.status")
  const locale = useLocale()

  // Build translation object
  const translations = useMemo(() => ({
    columns: {
      url: tColumns("common.url"),
      status: tColumns("common.status"),
      length: tColumns("directory.length"),
      words: tColumns("directory.words"),
      lines: tColumns("directory.lines"),
      contentType: tColumns("endpoint.contentType"),
      duration: tColumns("directory.duration"),
      createdAt: tColumns("common.createdAt"),
    },
    actions: {
      selectAll: tCommon("actions.selectAll"),
      selectRow: tCommon("actions.selectRow"),
    },
  }), [tColumns, tCommon])

  // Get target info (for URL matching validation)
  const { data: target } = useTarget(targetId || 0, { enabled: !!targetId })

  const targetQuery = useTargetDirectories(
    targetId || 0,
    {
      page: pagination.pageIndex + 1,
      pageSize: pagination.pageSize,
      filter: filterQuery || undefined,
    },
    { enabled: !!targetId }
  )

  const scanQuery = useScanDirectories(
    scanId || 0,
    {
      page: pagination.pageIndex + 1,
      pageSize: pagination.pageSize,
      filter: filterQuery || undefined,
    },
    { enabled: !!scanId }
  )

  const activeQuery = targetId ? targetQuery : scanQuery
  const { data, isLoading, isFetching, error, refetch } = activeQuery
  const { isSearching, handleSearchChange: handleFilterChange } = useSearchState({
    isFetching,
    setSearchValue: setFilterQuery,
    onResetPage: () => setPagination((prev) => ({ ...prev, pageIndex: 0 })),
  })

  const formatDate = useCallback((dateString: string) => {
    return new Date(dateString).toLocaleString(getDateLocale(locale), {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    })
  }, [locale])

  const columns = useMemo(
    () =>
      createDirectoryColumns({
        formatDate,
        t: translations,
      }),
    [formatDate, translations]
  )

  const directories: Directory[] = useMemo(() => {
    if (!data?.results) return []
    return data.results
  }, [data])

  const paginationInfo = data
    ? buildPaginationInfo({
      ...normalizePagination(data, pagination.pageIndex + 1, pagination.pageSize),
      minTotalPages: 1,
    })
    : undefined

  const handleSelectionChange = useCallback((selectedRows: Directory[]) => {
    setSelectedDirectories(selectedRows)
  }, [])

  // Convert nanoseconds to milliseconds
  const formatDurationNsToMs = (durationNs: number | null | undefined): string => {
    if (durationNs === null || durationNs === undefined) return ''
    return String(Math.floor(durationNs / 1_000_000))
  }

  // Generate CSV content
  const generateCSV = (items: Directory[]): string => {
    const BOM = '\ufeff'
    const headers = [
      'url', 'status', 'content_length', 'words',
      'lines', 'content_type', 'duration', 'created_at'
    ]
    
    const rows = items.map(item => [
      escapeCSV(item.url),
      escapeCSV(item.status),
      escapeCSV(item.contentLength),
      escapeCSV(item.words),
      escapeCSV(item.lines),
      escapeCSV(item.contentType),
      escapeCSV(formatDurationNsToMs(item.duration)),
      escapeCSV(formatDateForCSV(item.createdAt))
    ].join(','))
    
    return BOM + [headers.join(','), ...rows].join('\n')
  }

  // Handle download all directories
  const handleDownloadAll = async () => {
    try {
      let blob: Blob | null = null

      if (scanId) {
        const data = await DirectoryService.exportDirectoriesByScanId(scanId)
        blob = data
      } else if (targetId) {
        const data = await DirectoryService.exportDirectoriesByTargetId(targetId)
        blob = data
      } else {
        if (!directories || directories.length === 0) {
          return
        }
        const csvContent = generateCSV(directories)
        blob = new Blob([csvContent], { type: "text/csv;charset=utf-8" })
      }

      if (!blob) return

      const prefix = scanId ? `scan-${scanId}` : targetId ? `target-${targetId}` : "directories"
      downloadBlob(blob, `${prefix}-directories-${Date.now()}.csv`)
    } catch {
      toast.error(tToast("downloadFailed"))
    }
  }

  // Handle download selected directories
  const handleDownloadSelected = () => {
    if (selectedDirectories.length === 0) {
      return
    }
    const csvContent = generateCSV(selectedDirectories)
    const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8" })
    const prefix = scanId ? `scan-${scanId}` : targetId ? `target-${targetId}` : "directories"
    downloadBlob(blob, `${prefix}-directories-selected-${Date.now()}.csv`)
  }

  // Handle bulk delete
  const handleBulkDelete = async () => {
    if (selectedDirectories.length === 0) return
    
    setIsDeleting(true)
    try {
      const ids = selectedDirectories.map(d => d.id)
      const result = await DirectoryService.bulkDelete(ids)
      toast.success(tToast("deleteSuccess", { count: result.deletedCount }))
      setSelectedDirectories([])
      setDeleteDialogOpen(false)
      refetch()
    } catch {
      toast.error(tToast("deleteFailed"))
    } finally {
      setIsDeleting(false)
    }
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="rounded-full bg-destructive/10 p-3 mb-4">
          <AlertTriangle className="h-10 w-10 text-destructive" />
        </div>
        <h3 className="text-lg font-semibold mb-2">{tStatus("error")}</h3>
        <p className="text-muted-foreground text-center mb-4">
          {tStatus("error")}
        </p>
        <Button onClick={() => refetch()}>{tCommon("actions.retry")}</Button>
      </div>
    )
  }

  if (isLoading && !data) {
    return (
      <DataTableSkeleton
        toolbarButtonCount={2}
        rows={6}
        columns={5}
      />
    )
  }

  return (
    <>
      <DirectoriesDataTable
        data={directories}
        columns={columns}
        filterValue={filterQuery}
        onFilterChange={handleFilterChange}
        isSearching={isSearching}
        pagination={pagination}
        setPagination={setPagination}
        paginationInfo={paginationInfo}
        onPaginationChange={setPagination}
        onSelectionChange={handleSelectionChange}
        onDownloadAll={handleDownloadAll}
        onDownloadSelected={handleDownloadSelected}
        onBulkDelete={targetId ? () => setDeleteDialogOpen(true) : undefined}
        onBulkAdd={targetId ? () => setBulkAddDialogOpen(true) : undefined}
      />

      {/* Bulk add dialog */
      /* ... */
      }
      {targetId && (
        <BulkAddUrlsDialog
          targetId={targetId}
          assetType="directory"
          targetName={target?.name}
          targetType={target?.type as TargetType}
          open={bulkAddDialogOpen}
          onOpenChange={setBulkAddDialogOpen}
          onSuccess={() => refetch()}
        />
      )}

      {/* Delete confirmation dialog */}
      <ConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        title={tCommon("actions.confirmDelete")}
        description={tCommon("actions.deleteConfirmMessage", { count: selectedDirectories.length })}
        onConfirm={handleBulkDelete}
        loading={isDeleting}
        variant="destructive"
      />
    </>
  )
}
