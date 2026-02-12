import { useQuery, keepPreviousData } from "@tanstack/react-query"
import {
  getScans,
  getScan,
  getScanStatistics,
} from "@/services/scan.service"
import type { GetScansParams } from "@/types/scan.types"
import { scanKeys } from "./keys"

export function useScans(params: GetScansParams = { page: 1, pageSize: 10 }) {
  return useQuery({
    queryKey: scanKeys.list(params),
    queryFn: () => getScans(params),
    placeholderData: keepPreviousData,
  })
}

export function useRunningScans(page = 1, pageSize = 10) {
  return useScans({ page, pageSize, status: "running" })
}

/**
 * 获取目标的扫描历史
 */
export function useTargetScans(targetId: number, pageSize = 5) {
  return useQuery({
    queryKey: scanKeys.target(targetId, pageSize),
    queryFn: () => getScans({ target: targetId, pageSize }),
    enabled: !!targetId,
  })
}

export function useScan(id: number) {
  return useQuery({
    queryKey: scanKeys.detail(id),
    queryFn: () => getScan(id),
    enabled: !!id,
  })
}

/**
 * 获取扫描统计数据
 */
export function useScanStatistics() {
  return useQuery({
    queryKey: scanKeys.statistics(),
    queryFn: getScanStatistics,
  })
}
