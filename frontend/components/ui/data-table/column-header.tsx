"use client"

import type { Column } from "@tanstack/react-table"
import { ChevronUp, ChevronDown, ChevronsUpDown } from "@/components/icons"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"

interface DataTableColumnHeaderProps<TData, TValue> {
  column: Column<TData, TValue>
  title: string
  className?: string
}

/**
 * Unified column header component
 * 
 * Features:
 * - Supports sort indicators
 * - Click to toggle sort direction
 * - Unified styling (both sortable and non-sortable headers have consistent font size)
 */
export function DataTableColumnHeader<TData, TValue>({
  column,
  title,
  className,
}: DataTableColumnHeaderProps<TData, TValue>) {
  // 不可排序的列：使用与可排序列相同的字体大小
  if (!column.getCanSort()) {
    return <div className={cn("text-left text-sm", className)}>{title}</div>
  }

  const sorted = column.getIsSorted()

  return (
    <Button
      variant="ghost"
      size="sm"
      className={cn("-ml-3 h-8 data-[state=open]:bg-accent", className)}
      onClick={() => column.toggleSorting(sorted === "asc")}
    >
      <span className="whitespace-nowrap">{title}</span>
      {sorted === "desc" ? (
        <ChevronDown className="ml-2 h-4 w-4 shrink-0" />
      ) : sorted === "asc" ? (
        <ChevronUp className="ml-2 h-4 w-4 shrink-0" />
      ) : (
        <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0" />
      )}
    </Button>
  )
}
