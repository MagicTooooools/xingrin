"use client"

import React, { useState, useMemo } from "react"
import { useTranslations, useLocale } from "next-intl"
import { VulnerabilitiesDataTable, type ReviewFilter } from "./vulnerabilities-data-table"
import { createVulnerabilityColumns } from "./vulnerabilities-columns"
import { VulnerabilityDetailDialog } from "./vulnerability-detail-dialog"
import { DataTableSkeleton } from "@/components/ui/data-table-skeleton"
import { getDateLocale } from "@/lib/date-utils"
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
import {
  useScanVulnerabilities,
  useTargetVulnerabilities,
  useAllVulnerabilities,
  useMarkAsReviewed,
  useMarkAsUnreviewed,
  useBulkMarkAsReviewed,
  useBulkMarkAsUnreviewed,
  useVulnerabilityStats,
  useTargetVulnerabilityStats,
} from "@/hooks/use-vulnerabilities"

interface VulnerabilitiesDetailViewProps {
  /** Used in scan history page: view vulnerabilities by scan dimension */
  scanId?: number
  /** Used in target detail page: view vulnerabilities by target dimension */
  targetId?: number
  /** Hide toolbar (search, column controls, etc.) */
  hideToolbar?: boolean
}

export function VulnerabilitiesDetailView({
  scanId,
  targetId,
  hideToolbar = false,
}: VulnerabilitiesDetailViewProps) {
  const [selectedVulnerabilities, setSelectedVulnerabilities] = useState<Vulnerability[]>([])
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [vulnerabilityToDelete, setVulnerabilityToDelete] = useState<Vulnerability | null>(null)
  const [bulkDeleteDialogOpen, setBulkDeleteDialogOpen] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [detailDialogOpen, setDetailDialogOpen] = useState(false)
  const [selectedVulnerability, setSelectedVulnerability] = useState<Vulnerability | null>(null)
  const [reviewFilter, setReviewFilter] = useState<ReviewFilter>("all")

  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  })

  // Internationalization
  const tColumns = useTranslations("columns")
  const tCommon = useTranslations("common")
  const tTooltips = useTranslations("tooltips")
  const tSeverity = useTranslations("severity")
  const tConfirm = useTranslations("common.confirm")
  const locale = useLocale()

  // Build translation object
  const translations = useMemo(() => ({
    columns: {
      status: tColumns("common.status"),
      severity: tColumns("vulnerability.severity"),
      source: tColumns("vulnerability.source"),
      vulnType: tColumns("vulnerability.vulnType"),
      url: tColumns("common.url"),
      createdAt: tColumns("common.createdAt"),
    },
    actions: {
      details: tCommon("actions.details"),
      selectAll: tCommon("actions.selectAll"),
      selectRow: tCommon("actions.selectRow"),
    },
    tooltips: {
      vulnDetails: tTooltips("vulnDetails"),
      reviewed: tTooltips("reviewed"),
      pending: tTooltips("pending"),
    },
    severity: {
      critical: tSeverity("critical"),
      high: tSeverity("high"),
      medium: tSeverity("medium"),
      low: tSeverity("low"),
      info: tSeverity("info"),
    },
  }), [tColumns, tCommon, tTooltips, tSeverity])

  // Review mutations
  const markAsReviewed = useMarkAsReviewed()
  const markAsUnreviewed = useMarkAsUnreviewed()
  const bulkMarkAsReviewed = useBulkMarkAsReviewed()
  const bulkMarkAsUnreviewed = useBulkMarkAsUnreviewed()

  // Smart filter query
  const [filterQuery, setFilterQuery] = useState("")

  const handleFilterChange = (value: string) => {
    setFilterQuery(value)
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
  }

  const handleReviewFilterChange = (filter: ReviewFilter) => {
    setReviewFilter(filter)
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
  }

  // Convert review filter to API parameter
  const isReviewedParam = reviewFilter === "all" ? undefined : reviewFilter === "reviewed"

  const paginationParams = {
    page: pagination.pageIndex + 1,
    pageSize: pagination.pageSize,
    isReviewed: isReviewedParam,
  }

  // Load by scan dimension (scan history page)
  const scanQuery = useScanVulnerabilities(
    scanId ?? 0,
    paginationParams,
    { enabled: !!scanId },
    filterQuery || undefined,
  )

  // Load by target dimension (target detail page)
  const targetQuery = useTargetVulnerabilities(
    targetId ?? 0,
    paginationParams,
    { enabled: !!targetId && !scanId },
    filterQuery || undefined,
  )

  // Get all vulnerabilities (global vulnerabilities page)
  const allQuery = useAllVulnerabilities(
    paginationParams,
    { enabled: !scanId && !targetId },
    filterQuery || undefined,
  )

  // Select which query to use based on parameters
  const activeQuery = scanId ? scanQuery : targetId ? targetQuery : allQuery
  const isQueryLoading = activeQuery.isLoading

  const vulnerabilities = activeQuery.data?.vulnerabilities ?? []
  const paginationInfo = activeQuery.data?.pagination ?? {
    total: vulnerabilities.length,
    page: pagination.pageIndex + 1,
    pageSize: pagination.pageSize,
    totalPages: 1,
  }

  // Get pending count from stats API (only for global and target pages, not scan)
  const globalStatsQuery = useVulnerabilityStats({ enabled: !scanId && !targetId })
  const targetStatsQuery = useTargetVulnerabilityStats(targetId ?? 0, { enabled: !!targetId && !scanId })
  
  // Use stats API pendingCount, fallback to 0 for scan page (no review feature)
  const pendingCount = scanId 
    ? 0 
    : (targetId ? targetStatsQuery.data?.pendingCount : globalStatsQuery.data?.pendingCount) ?? 0


  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleString(getDateLocale(locale), {
      year: "numeric",
      month: "numeric",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    })
  }

  const navigate = (path: string) => {
    console.log("Navigate to:", path)
  }

  const handleViewDetail = (vulnerability: Vulnerability) => {
    setSelectedVulnerability(vulnerability)
    setDetailDialogOpen(true)
  }

  const handleDeleteVulnerability = (vulnerability: Vulnerability) => {
    setVulnerabilityToDelete(vulnerability)
    setDeleteDialogOpen(true)
  }

  const confirmDelete = async () => {
    if (!vulnerabilityToDelete) return

    setDeleteDialogOpen(false)
    setIsLoading(true)

    setTimeout(() => {
      console.log("Delete vulnerability:", vulnerabilityToDelete.id)
      setVulnerabilityToDelete(null)
      setIsLoading(false)
    }, 1000)
  }

  const handleBulkDelete = () => {
    if (selectedVulnerabilities.length === 0) {
      return
    }
    setBulkDeleteDialogOpen(true)
  }

  const confirmBulkDelete = async () => {
    if (selectedVulnerabilities.length === 0) return

    const deletedIds = selectedVulnerabilities.map(vulnerability => vulnerability.id)

    setBulkDeleteDialogOpen(false)
    setIsLoading(true)

    setTimeout(() => {
      console.log("Bulk delete vulnerabilities:", deletedIds)
      setSelectedVulnerabilities([])
      setIsLoading(false)
    }, 1000)
  }

  const handlePaginationChange = (newPagination: { pageIndex: number; pageSize: number }) => {
    setPagination(newPagination)
  }

  // Handle toggle review status for single vulnerability
  const handleToggleReview = (vulnerability: Vulnerability) => {
    if (vulnerability.isReviewed) {
      markAsUnreviewed.mutate(vulnerability.id)
    } else {
      markAsReviewed.mutate(vulnerability.id)
    }
  }

  // Handle bulk mark as reviewed
  const handleBulkMarkAsReviewed = () => {
    if (selectedVulnerabilities.length === 0) return
    const ids = selectedVulnerabilities.map(v => v.id)
    bulkMarkAsReviewed.mutate(ids, {
      onSuccess: () => {
        setSelectedVulnerabilities([])
      },
    })
  }

  // Handle bulk mark as pending
  const handleBulkMarkAsPending = () => {
    if (selectedVulnerabilities.length === 0) return
    const ids = selectedVulnerabilities.map(v => v.id)
    bulkMarkAsUnreviewed.mutate(ids, {
      onSuccess: () => {
        setSelectedVulnerabilities([])
      },
    })
  }

  // Handle download all vulnerabilities
  const handleDownloadAll = () => {
    // TODO: Implement download all vulnerabilities functionality
    console.log('Download all vulnerabilities')
  }

  // Handle download selected vulnerabilities
  const handleDownloadSelected = () => {
    // TODO: Implement download selected vulnerabilities functionality
    console.log('Download selected vulnerabilities:', selectedVulnerabilities)
    if (selectedVulnerabilities.length === 0) {
      return
    }
  }

  const vulnerabilityColumns = useMemo(
    () =>
      createVulnerabilityColumns({
        formatDate,
        handleViewDetail,
        onToggleReview: handleToggleReview,
        t: translations,
      }),
    [handleViewDetail, handleToggleReview, translations]
  )

  if ((isLoading || isQueryLoading) && !activeQuery.data) {
    return (
      <DataTableSkeleton
        toolbarButtonCount={2}
        rows={6}
        columns={6}
      />
    )
  }

  return (
    <>
      <VulnerabilityDetailDialog
        vulnerability={selectedVulnerability}
        open={detailDialogOpen}
        onOpenChange={setDetailDialogOpen}
      />

      <VulnerabilitiesDataTable
        data={vulnerabilities}
        columns={vulnerabilityColumns}
        filterValue={filterQuery}
        onFilterChange={handleFilterChange}
        pagination={pagination}
        setPagination={setPagination}
        paginationInfo={{
          total: paginationInfo.total,
          page: paginationInfo.page,
          pageSize: paginationInfo.pageSize,
          totalPages: paginationInfo.totalPages,
        }}
        onPaginationChange={handlePaginationChange}
        onSelectionChange={setSelectedVulnerabilities}
        hideToolbar={hideToolbar}
        reviewFilter={reviewFilter}
        onReviewFilterChange={handleReviewFilterChange}
        pendingCount={pendingCount}
        selectedRows={selectedVulnerabilities}
        onBulkMarkAsReviewed={handleBulkMarkAsReviewed}
        onBulkMarkAsPending={handleBulkMarkAsPending}
      />

      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{tConfirm("deleteTitle")}</AlertDialogTitle>
            <AlertDialogDescription>
              {tConfirm("deleteVulnMessage", { name: vulnerabilityToDelete?.vulnType ?? "" })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{tCommon("actions.cancel")}</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {tCommon("actions.delete")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <AlertDialog open={bulkDeleteDialogOpen} onOpenChange={setBulkDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{tConfirm("bulkDeleteTitle")}</AlertDialogTitle>
            <AlertDialogDescription>
              {tConfirm("bulkDeleteVulnMessage", { count: selectedVulnerabilities.length })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className="mt-2 p-2 bg-muted rounded-md max-h-96 overflow-y-auto">
            <ul className="text-sm space-y-1">
              {selectedVulnerabilities.map((vulnerability) => (
                <li key={vulnerability.id} className="flex items-center">
                  <span className="font-medium">{vulnerability.vulnType}</span>
                </li>
              ))}
            </ul>
          </div>
          <AlertDialogFooter>
            <AlertDialogCancel>{tCommon("actions.cancel")}</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmBulkDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {tConfirm("deleteVulnCount", { count: selectedVulnerabilities.length })}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
