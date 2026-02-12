"use client"

import * as React from "react"
import { useTranslations } from "next-intl"
import { UnifiedDataTable } from "@/components/ui/data-table/unified-data-table"
import { SimpleSearchToolbar } from "@/components/ui/data-table/simple-search-toolbar"
import { useSimpleSearchState } from "@/components/ui/data-table/use-simple-search"
import type { OrganizationDataTableProps } from "@/types/organization.types"

export function OrganizationDataTable({
  data,
  columns,
  onAddNew,
  onBulkDelete,
  onSelectionChange,
  searchPlaceholder,
  searchValue,
  onSearch,
  isSearching,
  pagination: externalPagination,
  paginationInfo,
  onPaginationChange,
}: OrganizationDataTableProps) {
  const t = useTranslations("organization")
  const tActions = useTranslations("common.actions")
  const {
    value: localSearchValue,
    setValue: setLocalSearchValue,
    submit: handleSearchSubmit,
  } = useSimpleSearchState({ searchValue, onSearch })

  // 默认排序
  const defaultSorting = [{ id: "createdAt", desc: true }]

  return (
    <UnifiedDataTable
      data={data}
      columns={columns}
      getRowId={(row) => String(row.id)}
      state={{
        pagination: externalPagination,
        paginationInfo,
        onPaginationChange,
        onSelectionChange,
        defaultSorting,
      }}
      actions={{
        onBulkDelete,
        bulkDeleteLabel: tActions("delete"),
        onAddNew,
        addButtonLabel: t("addOrganization"),
      }}
      ui={{
        emptyMessage: t("noResults"),
        toolbarLeft: (
          <SimpleSearchToolbar
            value={localSearchValue}
            onChange={setLocalSearchValue}
            onSubmit={handleSearchSubmit}
            loading={isSearching}
            placeholder={searchPlaceholder ?? t("searchPlaceholder")}
          />
        ),
      }}
    />
  )
}
