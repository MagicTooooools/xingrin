"use client"

import * as React from "react"
import type { ColumnDef } from "@tanstack/react-table"
import { useTranslations } from "next-intl"
import { ChevronDown, CheckCircle, Circle, X } from "lucide-react"
import { UnifiedDataTable } from "@/components/ui/data-table"
import { PREDEFINED_FIELDS, type FilterField } from "@/components/common/smart-filter-input"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuCheckboxItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import type { Vulnerability, VulnerabilitySeverity } from "@/types/vulnerability.types"
import type { PaginationInfo } from "@/types/common.types"
import type { DownloadOption } from "@/types/data-table.types"

// Review filter type
export type ReviewFilter = "all" | "pending" | "reviewed"

// Severity filter type
export type SeverityFilter = VulnerabilitySeverity | "all"

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
  // New: severity filter
  severityFilter?: SeverityFilter
  onSeverityFilterChange?: (filter: SeverityFilter) => void
  // New: source filter
  sourceFilter?: string
  onSourceFilterChange?: (source: string) => void
  availableSources?: string[]
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
  severityFilter = "all",
  onSeverityFilterChange,
  sourceFilter = "all",
  onSourceFilterChange,
  availableSources = [],
}: VulnerabilitiesDataTableProps) {
  const t = useTranslations("common.status")
  const tDownload = useTranslations("common.download")
  const tActions = useTranslations("common.actions")
  const tVuln = useTranslations("vulnerabilities")
  const tSeverity = useTranslations("severity")
  
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

  // Severity options for dropdown
  const severityOptions: SeverityFilter[] = ["all", "critical", "high", "medium", "low", "info"]

  // Right toolbar content - bulk actions, filters and review tabs
  const rightToolbarContent = (
    <>
      {/* Severity dropdown filter */}
      {onSeverityFilterChange && (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="sm" className="h-8">
              {severityFilter === "all" ? tVuln("severity") : tSeverity(severityFilter)}
              <ChevronDown className="ml-1 h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {severityOptions.map((sev) => (
              <DropdownMenuCheckboxItem
                key={sev}
                checked={severityFilter === sev}
                onCheckedChange={() => onSeverityFilterChange(sev)}
              >
                {sev === "all" ? tVuln("reviewStatus.all") : tSeverity(sev)}
              </DropdownMenuCheckboxItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      )}

      {/* Source dropdown filter */}
      {onSourceFilterChange && availableSources.length > 0 && (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="sm" className="h-8">
              {sourceFilter === "all" ? tVuln("source") : sourceFilter}
              <ChevronDown className="ml-1 h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuCheckboxItem
              checked={sourceFilter === "all"}
              onCheckedChange={() => onSourceFilterChange("all")}
            >
              {tVuln("reviewStatus.all")}
            </DropdownMenuCheckboxItem>
            {availableSources.map((src) => (
              <DropdownMenuCheckboxItem
                key={src}
                checked={sourceFilter === src}
                onCheckedChange={() => onSourceFilterChange(src)}
              >
                {src}
              </DropdownMenuCheckboxItem>
            ))}
          </DropdownMenuContent>
        </DropdownMenu>
      )}

      {/* Review filter tabs */}
      {onReviewFilterChange && (
        <Tabs value={reviewFilter} onValueChange={(v) => onReviewFilterChange(v as ReviewFilter)}>
          <TabsList>
            <TabsTrigger value="all">
              {tVuln("reviewStatus.all")}
            </TabsTrigger>
            <TabsTrigger value="pending">
              {tVuln("reviewStatus.pending")}
              {pendingCount > 0 && (
                <Badge variant="secondary" className="ml-1.5 h-5 min-w-5 rounded-full px-1.5 text-xs">
                  {pendingCount}
                </Badge>
              )}
            </TabsTrigger>
            <TabsTrigger value="reviewed">
              {tVuln("reviewStatus.reviewed")}
            </TabsTrigger>
          </TabsList>
        </Tabs>
      )}
    </>
  )

  // Floating action bar for bulk operations
  const floatingActionBar = selectedRows.length > 0 && (onBulkMarkAsReviewed || onBulkMarkAsPending) && (
    <div className="fixed bottom-6 left-[calc(50vw+var(--sidebar-width,14rem)/2)] -translate-x-1/2 z-50 animate-in slide-in-from-bottom-4 fade-in duration-200">
      <div className="flex items-center gap-3 bg-background border rounded-lg shadow-lg px-4 py-2.5">
        <span className="text-sm text-muted-foreground">
          {tVuln("selected", { count: selectedRows.length })}
        </span>
        <div className="h-4 w-px bg-border" />
        {onBulkMarkAsReviewed && (
          <Button
            variant="outline"
            size="sm"
            onClick={onBulkMarkAsReviewed}
            className="h-8"
          >
            <CheckCircle className="h-4 w-4 mr-1.5" />
            {tVuln("markAsReviewed")}
          </Button>
        )}
        {onBulkMarkAsPending && (
          <Button
            variant="outline"
            size="sm"
            onClick={onBulkMarkAsPending}
            className="h-8"
          >
            <Circle className="h-4 w-4 mr-1.5" />
            {tVuln("markAsPending")}
          </Button>
        )}
        {onSelectionChange && (
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onSelectionChange([])}
            className="h-8 w-8 text-muted-foreground hover:text-foreground"
          >
            <X className="h-4 w-4" />
          </Button>
        )}
      </div>
    </div>
  )

  return (
    <>
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
        toolbarRight={rightToolbarContent}
        // Empty state
        emptyMessage={t("noData")}
      />
      {floatingActionBar}
    </>
  )
}
