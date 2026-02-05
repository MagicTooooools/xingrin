"use client"

import * as React from "react"
import type { ColumnDef } from "@tanstack/react-table"
import { useTranslations } from "next-intl"
import { IconSearch, IconLoader2 } from "@/components/icons"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { UnifiedDataTable } from "@/components/ui/data-table/unified-data-table"
import type { ScheduledScan } from "@/types/scheduled-scan.types"
import type { PaginationInfo } from "@/types/common.types"

interface ScheduledScanDataTableProps {
  data: ScheduledScan[]
  columns: ColumnDef<ScheduledScan>[]
  onAddNew?: () => void
  searchPlaceholder?: string
  searchValue?: string
  onSearch?: (value: string) => void
  isSearching?: boolean
  addButtonText?: string
  // Server-side pagination related
  page?: number
  pageSize?: number
  total?: number
  totalPages?: number
  onPageChange?: (page: number) => void
  onPageSizeChange?: (pageSize: number) => void
}

/**
 * Scheduled scan data table component
 * Uses UnifiedDataTable unified component
 */
export function ScheduledScanDataTable({
  data = [],
  columns,
  onAddNew,
  searchPlaceholder,
  searchValue,
  onSearch,
  isSearching = false,
  addButtonText,
  page = 1,
  pageSize = 10,
  total = 0,
  totalPages = 1,
  onPageChange,
  onPageSizeChange,
}: ScheduledScanDataTableProps) {
  const t = useTranslations("common.status")
  const tScan = useTranslations("scan.scheduled")
  
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

  // Convert to pagination format required by UnifiedDataTable
  const pagination = { pageIndex: page - 1, pageSize }
  const paginationInfo: PaginationInfo = {
    total,
    totalPages,
    page,
    pageSize,
  }

  const handlePaginationChange = (newPagination: { pageIndex: number; pageSize: number }) => {
    if (newPagination.pageSize !== pageSize && onPageSizeChange) {
      onPageSizeChange(newPagination.pageSize)
    }
    if (newPagination.pageIndex !== page - 1 && onPageChange) {
      onPageChange(newPagination.pageIndex + 1)
    }
  }

  return (
    <UnifiedDataTable
      data={data}
      columns={columns}
      getRowId={(row) => String(row.id)}
      // Pagination
      pagination={pagination}
      paginationInfo={paginationInfo}
      onPaginationChange={handlePaginationChange}
      // Selection
      enableRowSelection={false}
      // Bulk operations
      showBulkDelete={false}
      onAddNew={onAddNew}
      addButtonLabel={addButtonText || tScan("createTitle")}
      // Empty state
      emptyMessage={t("noData")}
      // Custom search box
      toolbarLeft={
        <div className="flex items-center space-x-2">
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
        </div>
      }
    />
  )
}
