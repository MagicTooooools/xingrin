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
  pagination?: { pageIndex: number; pageSize: number }
  setPagination?: React.Dispatch<React.SetStateAction<{ pageIndex: number; pageSize: number }>>
  paginationInfo?: PaginationInfo
  onPaginationChange?: (pagination: { pageIndex: number; pageSize: number }) => void
  typeFilter?: string
  onTypeFilterChange?: (value: string) => void
}

/**
 * 目标数据表格组件 (organization 版本)
 * 使用 UnifiedDataTable 统一组件
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
  setPagination: setExternalPagination,
  paginationInfo,
  onPaginationChange,
  typeFilter,
  onTypeFilterChange,
}: TargetsDataTableProps) {
  const t = useTranslations("common.status")
  const tTarget = useTranslations("target")
  const tTooltips = useTranslations("tooltips")
  const tCommon = useTranslations("common")
  
  const {
    value: localSearchValue,
    setValue: setLocalSearchValue,
    submit: handleSearchSubmit,
  } = useSimpleSearchState({ searchValue, onSearch })

  return (
    <UnifiedDataTable
      data={data}
      columns={columns}
      getRowId={(row) => String(row.id)}
      state={{
        pagination: externalPagination,
        setPagination: setExternalPagination,
        paginationInfo,
        onPaginationChange,
        onSelectionChange,
      }}
      actions={{
        showBulkDelete: !!onBulkDelete,
        onBulkDelete,
        bulkDeleteLabel: tTooltips("unlinkTarget"),
        showAddButton: !!onAddNew,
        onAddNew,
        onAddHover,
        addButtonLabel: addButtonText || tTarget("addTarget"),
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
                  <SelectValue placeholder={tCommon("actions.filter")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{tCommon("actions.all")}</SelectItem>
                  <SelectItem value="domain">{tTarget("types.domain")}</SelectItem>
                  <SelectItem value="ip">{tTarget("types.ip")}</SelectItem>
                  <SelectItem value="cidr">{tTarget("types.cidr")}</SelectItem>
                </SelectContent>
              </Select>
            ) : null}
          />
        ),
      }}
    />
  )
}
