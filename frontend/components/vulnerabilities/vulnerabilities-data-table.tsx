"use client"

import * as React from "react"
import type { ColumnDef } from "@tanstack/react-table"
import { useTranslations } from "next-intl"
import { CheckCircle, Circle } from "lucide-react"
import { UnifiedDataTable } from "@/components/ui/data-table"
import { PREDEFINED_FIELDS, type FilterField } from "@/components/common/smart-filter-input"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import type { Vulnerability } from "@/types/vulnerability.types"
import type { PaginationInfo } from "@/types/common.types"
import type { DownloadOption } from "@/types/data-table.types"

// Review filter type
export type ReviewFilter = "all" | "pending" | "reviewed"

// Vulnerability page filter fields
const VULNERABILITY_FILTER_FIELDS: FilterField[] = [
  { key: "type", label: "Type", description: "Vulnerability type" },
  PREDEFINED_FIELDS.severity,
  { key: "source", label: "Source", description: "Scanner source" },
  PREDEFINED_FIELDS.url,
]

// Vulnerability page examples
const VULNERABILITY_FILTER_EXAMPLES = [
  'type="xss" || type="sqli"',
  'severity="critical" || severity="high"',
  'source="nuclei" && severity="high"',
  'type="xss" && url="/api/*"',
]

interface VulnerabilitiesDataTableProps {
  data: Vulnerability[]
  columns: ColumnDef<Vulnerability>[]
  filterValue?: string
  onFilterChange?: (value: string) => void
  pagination?: { pageIndex: number; pageSize: number }
  setPagination?: React.Dispatch<React.SetStateAction<{ pageIndex: number; pageSize: number }>>
  paginationInfo?: PaginationInfo
  onPaginationChange?: (pagination: { pageIndex: number; pageSize: number }) => void
  onBulkDelete?: () => void
  onSelectionChange?: (selectedRows: Vulnerability[]) => void
  onDownloadAll?: () => void
  onDownloadSelected?: () => void
  hideToolbar?: boolean
  // Review status props
  reviewFilter?: ReviewFilter
  onReviewFilterChange?: (filter: ReviewFilter) => void
  pendingCount?: number
  selectedRows?: Vulnerability[]
  onBulkMarkAsReviewed?: () => void
  onBulkMarkAsPending?: () => void
}

export function VulnerabilitiesDataTable({
  data = [],
  columns,
  filterValue,
  onFilterChange,
  pagination,
  setPagination,
  paginationInfo,
  onPaginationChange,
  onBulkDelete,
  onSelectionChange,
  onDownloadAll,
  onDownloadSelected,
  hideToolbar = false,
  reviewFilter = "all",
  onReviewFilterChange,
  pendingCount = 0,
  selectedRows = [],
  onBulkMarkAsReviewed,
  onBulkMarkAsPending,
}: VulnerabilitiesDataTableProps) {
  const t = useTranslations("common.status")
  const tDownload = useTranslations("common.download")
  const tActions = useTranslations("common.actions")
  const tVuln = useTranslations("vulnerabilities")
  
  // Handle smart filter search
  const handleFilterSearch = (rawQuery: string) => {
    onFilterChange?.(rawQuery)
  }

  // Download options
  const downloadOptions: DownloadOption[] = []
  if (onDownloadAll) {
    downloadOptions.push({
      key: "all",
      label: tDownload("all"),
      onClick: onDownloadAll,
    })
  }
  if (onDownloadSelected) {
    downloadOptions.push({
      key: "selected",
      label: tDownload("selected"),
      onClick: onDownloadSelected,
      disabled: (count) => count === 0,
    })
  }

  // Custom toolbar content for review filter tabs and bulk actions
  const reviewToolbarContent = (
    <div className="flex items-center gap-4">
      {/* Review filter tabs */}
      {onReviewFilterChange && (
        <Tabs value={reviewFilter} onValueChange={(v) => onReviewFilterChange(v as ReviewFilter)}>
          <TabsList className="h-8">
            <TabsTrigger value="all" className="text-xs px-3 h-7">
              {tVuln("reviewStatus.all")}
            </TabsTrigger>
            <TabsTrigger value="pending" className="text-xs px-3 h-7">
              {tVuln("reviewStatus.pending")}
              {pendingCount > 0 && (
                <Badge variant="secondary" className="ml-1.5 h-5 px-1.5 text-xs">
                  {pendingCount}
                </Badge>
              )}
            </TabsTrigger>
            <TabsTrigger value="reviewed" className="text-xs px-3 h-7">
              {tVuln("reviewStatus.reviewed")}
            </TabsTrigger>
          </TabsList>
        </Tabs>
      )}

      {/* Bulk review actions */}
      {selectedRows.length > 0 && (
        <div className="flex items-center gap-2">
          {onBulkMarkAsReviewed && (
            <Button
              variant="outline"
              size="sm"
              className="h-8"
              onClick={onBulkMarkAsReviewed}
            >
              <CheckCircle className="h-4 w-4 mr-1" />
              {tVuln("markAsReviewed")}
            </Button>
          )}
          {onBulkMarkAsPending && (
            <Button
              variant="outline"
              size="sm"
              className="h-8"
              onClick={onBulkMarkAsPending}
            >
              <Circle className="h-4 w-4 mr-1" />
              {tVuln("markAsPending")}
            </Button>
          )}
        </div>
      )}
    </div>
  )

  return (
    <UnifiedDataTable
      data={data}
      columns={columns}
      getRowId={(row) => String(row.id)}
      // Pagination
      pagination={pagination}
      setPagination={setPagination}
      paginationInfo={paginationInfo}
      onPaginationChange={onPaginationChange}
      // Smart filter
      searchMode="smart"
      searchValue={filterValue}
      onSearch={handleFilterSearch}
      filterFields={VULNERABILITY_FILTER_FIELDS}
      filterExamples={VULNERABILITY_FILTER_EXAMPLES}
      // Selection
      onSelectionChange={onSelectionChange}
      // Bulk operations
      onBulkDelete={onBulkDelete}
      bulkDeleteLabel={tActions("delete")}
      showAddButton={false}
      // Download
      downloadOptions={downloadOptions.length > 0 ? downloadOptions : undefined}
      // Toolbar
      hideToolbar={hideToolbar}
      toolbarLeft={reviewToolbarContent}
      // Empty state
      emptyMessage={t("noData")}
    />
  )
}
