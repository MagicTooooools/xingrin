"use client"

import { ScanLogListContent, ScanLogListEmptyState, ScanLogListLoadingState } from "@/components/scan/scan-log-list-sections"
import { useScanLogListState } from "@/components/scan/scan-log-list-state"

import type { ScanLog } from "@/services/scan.service"

interface ScanLogListProps {
  logs: ScanLog[]
  loading?: boolean
}

/**
 * 扫描日志列表组件
 * 复用 AnsiLogViewer 组件
 */
export function ScanLogList({ logs, loading }: ScanLogListProps) {
  const state = useScanLogListState({ logs })

  if (loading && logs.length === 0) {
    return <ScanLogListLoadingState />
  }

  if (logs.length === 0) {
    return <ScanLogListEmptyState />
  }

  return <ScanLogListContent content={state.content} />
}
