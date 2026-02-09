/**
 * 扫描日志轮询 Hook
 * 
 * 功能：
 * - 初始加载获取全部日志
 * - 增量轮询获取新日志（3s 间隔）
 * - 扫描结束后停止轮询
 */

import { useState, useEffect, useCallback, useRef } from 'react'
import { getScanLogs, type ScanLog } from '@/services/scan.service'

interface UseScanLogsOptions {
  scanId: number
  enabled?: boolean
  pollingInterval?: number  // 默认 3000ms
  maxLogs?: number  // 默认 5000，<=0 表示不限制
}

interface UseScanLogsReturn {
  logs: ScanLog[]
  loading: boolean
  refetch: () => void
}

export function useScanLogs({
  scanId,
  enabled = true,
  pollingInterval = 3000,
  maxLogs = 5000,
}: UseScanLogsOptions): UseScanLogsReturn {
  const [logs, setLogs] = useState<ScanLog[]>([])
  const [loading, setLoading] = useState(false)
  const lastLogIDRef = useRef<number | null>(null)
  const isMounted = useRef(true)
  
  const clampLogs = useCallback((items: ScanLog[]) => {
    if (!maxLogs || maxLogs <= 0) return items
    return items.length > maxLogs ? items.slice(-maxLogs) : items
  }, [maxLogs])

  const fetchLogs = useCallback(async (incremental = false) => {
    if (!enabled || !isMounted.current) return
    
    setLoading(true)
    try {
      const params: { limit: number; afterId?: number } = { limit: 200 }
      if (incremental && lastLogIDRef.current !== null) {
        params.afterId = lastLogIDRef.current
      }
      
      const response = await getScanLogs(scanId, params)
      const newLogs = response.results
      if (newLogs.length > 0) {
        lastLogIDRef.current = newLogs[newLogs.length - 1].id
      }
      
      if (!isMounted.current) return
      
      if (newLogs.length > 0) {
        if (incremental) {
          // 按 ID 去重，防止 React Strict Mode 或竞态条件导致的重复
          setLogs(prev => {
            const existingIds = new Set(prev.map(l => l.id))
            const uniqueNewLogs = newLogs.filter(l => !existingIds.has(l.id))
            if (uniqueNewLogs.length === 0) return prev
            return clampLogs([...prev, ...uniqueNewLogs])
          })
        } else {
          setLogs(clampLogs(newLogs))
        }
      }
    } catch (error) {
      void error
    } finally {
      if (isMounted.current) {
        setLoading(false)
      }
    }
  }, [scanId, enabled, clampLogs])
  
  // 初始加载
  useEffect(() => {
    isMounted.current = true
    if (enabled) {
      // 重置状态
      setLogs([])
      lastLogIDRef.current = null
      fetchLogs(false)
    }
    return () => {
      isMounted.current = false
    }
  }, [scanId, enabled, fetchLogs])
  
  // 轮询
  useEffect(() => {
    if (!enabled) return
    // pollingInterval <= 0 表示禁用轮询（避免 setInterval(0) 导致高频请求/卡顿）
    if (!pollingInterval || pollingInterval <= 0) return

    const interval = setInterval(() => {
      fetchLogs(true) // 增量查询
    }, pollingInterval)

    return () => clearInterval(interval)
  }, [enabled, pollingInterval, fetchLogs])
  
  const refetch = useCallback(() => {
    setLogs([])
    lastLogIDRef.current = null
    fetchLogs(false)
  }, [fetchLogs])

  // 当 maxLogs 变化时，主动裁剪缓存，避免长时间运行的内存占用增长
  useEffect(() => {
    if (!maxLogs || maxLogs <= 0) return
    setLogs(prev => (prev.length > maxLogs ? prev.slice(-maxLogs) : prev))
  }, [maxLogs])
  
  return { logs, loading, refetch }
}
