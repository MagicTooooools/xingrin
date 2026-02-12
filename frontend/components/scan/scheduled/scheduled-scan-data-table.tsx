"use client"

import * as React from "react"
import type { ColumnDef } from "@tanstack/react-table"
import { useTranslations } from "next-intl"
import { UnifiedDataTable } from "@/components/ui/data-table/unified-data-table"
import { SimpleSearchToolbar } from "@/components/ui/data-table/simple-search-toolbar"
import { useSimpleSearchState } from "@/components/ui/data-table/use-simple-search"
import type { ScheduledScan } from "@/types/scheduled-scan.types"
import type { PaginationInfo } from "@/types/common.types"
import { buildPaginationInfo } from "@/hooks/_shared/pagination"

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
  
  const {
    value: localSearchValue,
    setValue: setLocalSearchValue,
    submit: handleSearchSubmit,
  } = useSimpleSearchState({ searchValue, onSearch })

  // Convert to pagination format required by UnifiedDataTable
  const pagination = { pageIndex: page - 1, pageSize }
  const paginationInfo: PaginationInfo = buildPaginationInfo({
    total,
    page,
    pageSize,
    totalPages,
    minTotalPages: 1,
  })

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
      state={{
        pagination,
        paginationInfo,
        onPaginationChange: handlePaginationChange,
      }}
      behavior={{ enableRowSelection: false }}
      actions={{
        showBulkDelete: false,
        onAddNew,
        addButtonLabel: addButtonText || tScan("createTitle"),
      }}
      ui={{
        emptyMessage: t("noData"),
        toolbarLeft: (
          <SimpleSearchToolbar
            value={localSearchValue}
            onChange={setLocalSearchValue}
            onSubmit={handleSearchSubmit}
            loading={isSearching}
            placeholder={searchPlaceholder || tScan("searchPlaceholder")}
          />
        ),
      }}
    />
  )
}
