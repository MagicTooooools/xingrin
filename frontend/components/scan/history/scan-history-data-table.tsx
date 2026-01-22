"use client"

import * as React from "react"
import type { ColumnDef } from "@tanstack/react-table"
import { useTranslations } from "next-intl"
import { IconSearch, IconLoader2 } from "@tabler/icons-react"
import { Filter } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { UnifiedDataTable } from "@/components/ui/data-table"
import type { ScanRecord, ScanStatus } from "@/types/scan.types"
import type { PaginationInfo } from "@/types/common.types"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

interface ScanHistoryDataTableProps {
  data: ScanRecord[]
  columns: ColumnDef<ScanRecord>[]
  onAddNew?: () => void
  onBulkDelete?: () => void
  onSelectionChange?: (selectedRows: ScanRecord[]) => void
  searchPlaceholder?: string
  searchValue?: string
  onSearch?: (value: string) => void
  isSearching?: boolean
  addButtonText?: string
  pagination?: { pageIndex: number; pageSize: number }
  setPagination?: React.Dispatch<React.SetStateAction<{ pageIndex: number; pageSize: number }>>
  paginationInfo?: PaginationInfo
  onPaginationChange?: (pagination: { pageIndex: number; pageSize: number }) => void
  hideToolbar?: boolean
  hidePagination?: boolean
  pageSizeOptions?: number[]
  statusFilter?: ScanStatus | "all"
  onStatusFilterChange?: (status: ScanStatus | "all") => void
}

/**
 * Scan history data table component
 * Uses UnifiedDataTable unified component
 */
export function ScanHistoryDataTable({
  data = [],
  columns,
  onAddNew,
  onBulkDelete,
  onSelectionChange,
  searchPlaceholder,
  searchValue,
  onSearch,
  isSearching = false,
  addButtonText,
  pagination: externalPagination,
  setPagination: setExternalPagination,
  paginationInfo,
  onPaginationChange,
  hideToolbar = false,
  hidePagination = false,
  pageSizeOptions,
  statusFilter = "all",
  onStatusFilterChange,
}: ScanHistoryDataTableProps) {
  const t = useTranslations("common.status")
  const tScan = useTranslations("scan.history")
  const tActions = useTranslations("common.actions")
  
  // Search local state
  const [localSearchValue, setLocalSearchValue] = React.useState(searchValue || "")

  React.useEffect(() => {
    setLocalSearchValue(searchValue || "")
  }, [searchValue])

  const handleSearchSubmit = () => {
    if (onSearch) {
      onSearch(localSearchValue)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      handleSearchSubmit()
    }
  }

  // Status options
  const statusOptions: { value: ScanStatus | "all"; label: string }[] = [
    { value: "all", label: tScan("allStatus") },
    { value: "running", label: t("running") },
    { value: "completed", label: t("completed") },
    { value: "failed", label: t("failed") },
    { value: "pending", label: t("pending") },
    { value: "cancelled", label: t("cancelled") },
  ]

  return (
    <UnifiedDataTable
      data={data}
      columns={columns}
      getRowId={(row) => String(row.id)}
      // Pagination
      pagination={externalPagination}
      setPagination={setExternalPagination}
      paginationInfo={paginationInfo}
      onPaginationChange={onPaginationChange}
      hidePagination={hidePagination}
      pageSizeOptions={pageSizeOptions}
      // Selection
      onSelectionChange={onSelectionChange}
      // Bulk operations
      onBulkDelete={onBulkDelete}
      bulkDeleteLabel={tActions("delete")}
      onAddNew={onAddNew}
      addButtonLabel={addButtonText || tScan("title")}
      // Toolbar
      hideToolbar={hideToolbar}
      // Empty state
      emptyMessage={t("noData")}
      // Auto column sizing
      enableAutoColumnSizing
      // Custom search box and status filter
      toolbarLeft={
        <div className="flex items-center gap-2">
          <Input
            placeholder={searchPlaceholder || tScan("searchPlaceholder")}
            value={localSearchValue}
            onChange={(e) => setLocalSearchValue(e.target.value)}
            onKeyDown={handleKeyDown}
            className="h-8 max-w-sm"
          />
          <Button variant="outline" size="sm" onClick={handleSearchSubmit} disabled={isSearching}>
            {isSearching ? (
              <IconLoader2 className="h-4 w-4 animate-spin" />
            ) : (
              <IconSearch className="h-4 w-4" />
            )}
          </Button>
          {onStatusFilterChange && (
            <Select
              value={statusFilter}
              onValueChange={(value) => onStatusFilterChange(value as ScanStatus | "all")}
            >
              <SelectTrigger size="sm" className="w-auto">
                <Filter className="h-4 w-4" />
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {statusOptions.map((option) => (
                  <SelectItem key={option.value} value={option.value}>
                    {option.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
        </div>
      }
    />
  )
}
