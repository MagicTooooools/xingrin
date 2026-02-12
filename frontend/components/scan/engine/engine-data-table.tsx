"use client"

import * as React from "react"
import type { ColumnDef } from "@tanstack/react-table"
import { useTranslations } from "next-intl"
import { UnifiedDataTable } from "@/components/ui/data-table/unified-data-table"
import { SimpleSearchToolbar } from "@/components/ui/data-table/simple-search-toolbar"
import type { ScanEngine } from "@/types/engine.types"

// Component props type definitions
interface EngineDataTableProps {
  data: ScanEngine[]
  columns: ColumnDef<ScanEngine>[]
  onAddNew?: () => void
  searchPlaceholder?: string
  searchColumn?: string
  addButtonText?: string
}

/**
 * Scan engine data table component
 * Uses UnifiedDataTable unified component
 */
export function EngineDataTable({
  data = [],
  columns,
  onAddNew,
  searchPlaceholder,
  addButtonText,
}: EngineDataTableProps) {
  const t = useTranslations("common.status")
  const tEngine = useTranslations("scan.engine")
  
  // Local search state
  const [searchValue, setSearchValue] = React.useState("")

  // Filter data (local filtering)
  const filteredData = React.useMemo(() => {
    if (!searchValue) return data
    return data.filter((item) => {
      const name = item.name || ""
      return name.toLowerCase().includes(searchValue.toLowerCase())
    })
  }, [data, searchValue])

  return (
    <UnifiedDataTable
      data={filteredData}
      columns={columns}
      getRowId={(row) => String(row.id)}
      behavior={{ enableRowSelection: false }}
      actions={{
        onAddNew,
        addButtonLabel: addButtonText || tEngine("createEngine"),
        showBulkDelete: false,
      }}
      ui={{
        emptyMessage: t("noData"),
        toolbarLeft: (
          <SimpleSearchToolbar
            value={searchValue}
            onChange={setSearchValue}
            placeholder={searchPlaceholder || tEngine("searchPlaceholder")}
            showButton={false}
          />
        ),
      }}
    />
  )
}
