"use client"

import { useMemo, useRef } from "react"
import { AnsiLogViewer } from "@/components/settings/system-logs/ansi-log-viewer"
import type { ScanLog } from "@/services/scan.service"

interface ScanLogListProps {
  logs: ScanLog[]
  loading?: boolean
}

/**
 * 格式化时间为 HH:mm:ss
 */
function formatTime(isoString: string): string {
  try {
    const date = new Date(isoString)
    const h = String(date.getHours()).padStart(2, '0')
    const m = String(date.getMinutes()).padStart(2, '0')
    const s = String(date.getSeconds()).padStart(2, '0')
    return `${h}:${m}:${s}`
  } catch {
    return isoString
  }
}

/**
 * 扫描日志列表组件
 * 复用 AnsiLogViewer 组件
 */
export function ScanLogList({ logs, loading }: ScanLogListProps) {
  // 稳定的 content 引用，只有内容真正变化时才更新
  const contentRef = useRef('')
  const lastLogCountRef = useRef(0)
  const lastLogIdRef = useRef<number | null>(null)
  const firstLogIdRef = useRef<number | null>(null)
  
  // 将日志转换为纯文本格式
  const content = useMemo(() => {
    if (logs.length === 0) {
      contentRef.current = ''
      lastLogCountRef.current = 0
      lastLogIdRef.current = null
      firstLogIdRef.current = null
      return ''
    }
    
    // 检查是否真正需要更新
    const lastLog = logs[logs.length - 1]
    const firstLog = logs[0]

    const shouldRebuild =
      lastLogIdRef.current === null ||
      logs.length < lastLogCountRef.current ||
      (firstLogIdRef.current !== null && firstLog?.id !== firstLogIdRef.current)

    if (!shouldRebuild) {
      if (
        logs.length === lastLogCountRef.current &&
        lastLog?.id === lastLogIdRef.current
      ) {
        // 日志没有变化，返回缓存的 content
        return contentRef.current
      }

      const lastIndex = logs.findIndex((log) => log.id === lastLogIdRef.current)
      if (lastIndex !== -1) {
        const newLogs = logs.slice(lastIndex + 1)
        if (newLogs.length > 0) {
          const appended = newLogs.map(log => {
            const time = formatTime(log.createdAt)
            const levelTag = log.level.toUpperCase()
            return `[${time}] [${levelTag}] ${log.content}`
          }).join('\n')
          contentRef.current = contentRef.current
            ? `${contentRef.current}\n${appended}`
            : appended
        }
        lastLogCountRef.current = logs.length
        lastLogIdRef.current = lastLog?.id ?? null
        firstLogIdRef.current = firstLog?.id ?? null
        return contentRef.current
      }
    }

    // 需要完整重建（初次加载、截断、乱序或无法增量更新）
    const newContent = logs.map(log => {
      const time = formatTime(log.createdAt)
      const levelTag = log.level.toUpperCase()
      return `[${time}] [${levelTag}] ${log.content}`
    }).join('\n')

    contentRef.current = newContent
    lastLogCountRef.current = logs.length
    lastLogIdRef.current = lastLog?.id ?? null
    firstLogIdRef.current = firstLog?.id ?? null
    return newContent
  }, [logs])
  
  if (loading && logs.length === 0) {
    return (
      <div className="h-full flex items-center justify-center bg-[#1e1e1e] text-[#808080]">
        加载中...
      </div>
    )
  }
  
  if (logs.length === 0) {
    return (
      <div className="h-full flex items-center justify-center bg-[#1e1e1e] text-[#808080]">
        暂无日志
      </div>
    )
  }
  
  return (
    <div className="h-full">
      <AnsiLogViewer content={content} />
    </div>
  )
}
