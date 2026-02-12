"use client"

import * as React from "react"
import { Filter } from "@/components/icons"
import { useTranslations } from "next-intl"
import { UnifiedDataTable } from "@/components/ui/data-table/unified-data-table"
import { SimpleSearchToolbar } from "@/components/ui/data-table/simple-search-toolbar"
import { useSimpleSearchState } from "@/components/ui/data-table/use-simple-search"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import type { ColumnDef } from "@tanstack/react-table"
import type { Target } from "@/types/target.types"
import type { PaginationInfo } from "@/types/common.types"
import { buildPaginationInfo } from "@/hooks/_shared/pagination"

interface TargetsDataTableProps {
  data: Target[]
  columns: ColumnDef<Target>[]
  onAddNew?: () => void
  onAddHover?: () => void
  onBulkDelete?: () => void
  onSelectionChange?: (selectedRows: Target[]) => void
  searchPlaceholder?: string
  searchValue?: string
  onSearch?: (value: string) => void
  isSearching?: boolean
  addButtonText?: string
  // Pagination related props
  pagination?: { pageIndex: number, pageSize: number }
  onPaginationChange?: (pagination: { pageIndex: number, pageSize: number }) => void
  totalCount?: number
  manualPagination?: boolean
  // Type filter
  typeFilter?: string
  onTypeFilterChange?: (value: string) => void
  // Styling
  className?: string
  tableClassName?: string
  hideToolbar?: boolean
  hidePagination?: boolean
}

/**
 * Targets data table component (target version)
 * Uses UnifiedDataTable unified component
 */
export function TargetsDataTable({
  data = [],
  columns,
  onAddNew,
  onAddHover,
  onBulkDelete,
  onSelectionChange,
  searchPlaceholder,
  searchValue,
  onSearch,
  isSearching = false,
  addButtonText,
  pagination: externalPagination,
  onPaginationChange,
  totalCount,
  manualPagination = false,
  typeFilter,
  onTypeFilterChange,
  className,
  tableClassName,
  hideToolbar = false,
  hidePagination = false,
}: TargetsDataTableProps) {
  const t = useTranslations("common.status")
  const tActions = useTranslations("common.actions")
  const tTarget = useTranslations("target")
  
  // Internal pagination state
  const [internalPagination, setInternalPagination] = React.useState<{ pageIndex: number, pageSize: number }>({
    pageIndex: 0,
    pageSize: 10,
  })

  const {
    value: localSearchValue,
    setValue: setLocalSearchValue,
    submit: handleSearchSubmit,
  } = useSimpleSearchState({ searchValue, onSearch })

  const pagination = externalPagination || internalPagination

  // Handle pagination state change
  const handlePaginationChange = (newPagination: { pageIndex: number, pageSize: number }) => {
    if (onPaginationChange) {
      onPaginationChange(newPagination)
    } else {
      setInternalPagination(newPagination)
    }
  }

  // Build paginationInfo
  const paginationInfo: PaginationInfo | undefined =
    manualPagination && totalCount !== undefined
      ? buildPaginationInfo({
        total: totalCount ?? 0,
        page: pagination.pageIndex + 1,
        pageSize: pagination.pageSize,
        minTotalPages: 1,
      })
      : undefined

  return (
    <UnifiedDataTable
      data={data}
      columns={columns}
      getRowId={(row) => String(row.id)}
      state={{
        pagination,
        setPagination: onPaginationChange ? undefined : setInternalPagination,
        paginationInfo,
        onPaginationChange: handlePaginationChange,
        onSelectionChange,
      }}
      actions={{
        onBulkDelete,
        bulkDeleteLabel: tActions("delete"),
        onAddNew,
        onAddHover,
        addButtonLabel: addButtonText || tTarget("addTarget"),
        showAddButton: !!onAddNew,
      }}
      ui={{
        emptyMessage: t("noData"),
        toolbarLeft: (
          <SimpleSearchToolbar
            value={localSearchValue}
            onChange={setLocalSearchValue}
            onSubmit={handleSearchSubmit}
            loading={isSearching}
            placeholder={searchPlaceholder || tTarget("title")}
            after={onTypeFilterChange ? (
              <Select value={typeFilter || "all"} onValueChange={(value) => onTypeFilterChange(value === "all" ? "" : value)}>
                <SelectTrigger size="sm" className="w-auto">
                  <Filter className="h-4 w-4" />
                  <SelectValue placeholder={tActions("filter")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{tActions("all")}</SelectItem>
                  <SelectItem value="domain">{tTarget("types.domain")}</SelectItem>
                  <SelectItem value="ip">{tTarget("types.ip")}</SelectItem>
                  <SelectItem value="cidr">{tTarget("types.cidr")}</SelectItem>
                </SelectContent>
              </Select>
            ) : null}
          />
        ),
        className,
        tableClassName,
        hideToolbar,
        hidePagination,
      }}
    />
  )
}
