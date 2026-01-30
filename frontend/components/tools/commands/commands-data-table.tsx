"use client"

import * as React from "react"
import type { ColumnDef } from "@tanstack/react-table"
import { useTranslations } from "next-intl"
import { Input } from "@/components/ui/input"
import { UnifiedDataTable } from "@/components/ui/data-table"

interface CommandsDataTableProps<TData extends { id: number; displayName?: string }> {
  columns: ColumnDef<TData, unknown>[]
  data: TData[]
  onBulkDelete?: (selectedIds: number[]) => void
  onAdd?: () => void
}

export function CommandsDataTable<TData extends { id: number; displayName?: string }>({
  columns,
  data,
  onBulkDelete,
  onAdd,
}: CommandsDataTableProps<TData>) {
  const t = useTranslations("tools.commands")
  const tCommon = useTranslations("common")
  
  // Local search state
  const [searchValue, setSearchValue] = React.useState("")
  const [selectedRows, setSelectedRows] = React.useState<TData[]>([])

  // Filter data (local filtering)
  const filteredData = React.useMemo(() => {
    if (!searchValue) return data
    return data.filter((item) => {
      const displayName = item.displayName || ""
      return displayName.toLowerCase().includes(searchValue.toLowerCase())
    })
  }, [data, searchValue])

  // Handle bulk delete
  const handleBulkDelete = () => {
    if (onBulkDelete && selectedRows.length > 0) {
      const selectedIds = selectedRows.map((row) => row.id)
      onBulkDelete(selectedIds)
    }
  }

  return (
    <UnifiedDataTable
      data={filteredData}
      columns={columns}
      getRowId={(row) => String(row.id)}
      onSelectionChange={setSelectedRows}
      onBulkDelete={onBulkDelete ? handleBulkDelete : undefined}
      onAddNew={onAdd}
      addButtonLabel={tCommon("actions.add")}
      bulkDeleteLabel={tCommon("actions.delete")}
      emptyMessage={tCommon("status.noData")}
      toolbarLeft={
        <Input
          placeholder={t("searchPlaceholder")}
          value={searchValue}
          onChange={(e) => setSearchValue(e.target.value)}
          className="max-w-sm h-8"
        />
      }
    />
  )
}
