/**
 * 指纹管理 React Query Hooks
 */

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { FingerprintService } from "@/services/fingerprint.service"
import type { EholeFingerprint, GobyFingerprint, WappalyzerFingerprint, FingersFingerprint, FingerPrintHubFingerprint, ARLFingerprint } from "@/types/fingerprint.types"

type FingerprintListParams = { page?: number; pageSize?: number; filter?: string }

// Query Keys
export const fingerprintKeys = {
  all: ["fingerprints"] as const,
  stats: () => [...fingerprintKeys.all, "stats"] as const,
  ehole: {
    all: () => [...fingerprintKeys.all, "ehole"] as const,
    list: (params: FingerprintListParams) => [...fingerprintKeys.ehole.all(), "list", params] as const,
    detail: (id: number) => [...fingerprintKeys.ehole.all(), "detail", id] as const,
  },
  goby: {
    all: () => [...fingerprintKeys.all, "goby"] as const,
    list: (params: FingerprintListParams) => [...fingerprintKeys.goby.all(), "list", params] as const,
    detail: (id: number) => [...fingerprintKeys.goby.all(), "detail", id] as const,
  },
  wappalyzer: {
    all: () => [...fingerprintKeys.all, "wappalyzer"] as const,
    list: (params: FingerprintListParams) => [...fingerprintKeys.wappalyzer.all(), "list", params] as const,
    detail: (id: number) => [...fingerprintKeys.wappalyzer.all(), "detail", id] as const,
  },
  fingers: {
    all: () => [...fingerprintKeys.all, "fingers"] as const,
    list: (params: FingerprintListParams) => [...fingerprintKeys.fingers.all(), "list", params] as const,
    detail: (id: number) => [...fingerprintKeys.fingers.all(), "detail", id] as const,
  },
  fingerprinthub: {
    all: () => [...fingerprintKeys.all, "fingerprinthub"] as const,
    list: (params: FingerprintListParams) => [...fingerprintKeys.fingerprinthub.all(), "list", params] as const,
    detail: (id: number) => [...fingerprintKeys.fingerprinthub.all(), "detail", id] as const,
  },
  arl: {
    all: () => [...fingerprintKeys.all, "arl"] as const,
    list: (params: FingerprintListParams) => [...fingerprintKeys.arl.all(), "list", params] as const,
    detail: (id: number) => [...fingerprintKeys.arl.all(), "detail", id] as const,
  },
}

// ==================== EHole Hooks ====================

/**
 * 获取 EHole 指纹列表
 */
export function useEholeFingerprints(
  params: { page?: number; pageSize?: number; filter?: string } = {},
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: fingerprintKeys.ehole.list(params),
    queryFn: () => FingerprintService.getEholeFingerprints(params),
    ...options,
  })
}

/**
 * 获取 EHole 指纹详情
 */
export function useEholeFingerprint(id: number, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: fingerprintKeys.ehole.detail(id),
    queryFn: () => FingerprintService.getEholeFingerprint(id),
    enabled: id > 0 && options?.enabled !== false,
  })
}

/**
 * 创建 EHole 指纹
 */
export function useCreateEholeFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: Omit<EholeFingerprint, 'id' | 'createdAt'>) => 
      FingerprintService.createEholeFingerprint(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.ehole.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 更新 EHole 指纹
 */
export function useUpdateEholeFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<EholeFingerprint> }) =>
      FingerprintService.updateEholeFingerprint(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.ehole.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.ehole.detail(id) })
    },
  })
}

/**
 * 删除 EHole 指纹
 */
export function useDeleteEholeFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => FingerprintService.deleteEholeFingerprint(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.ehole.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 批量创建 EHole 指纹
 */
export function useBatchCreateEholeFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (fingerprints: Omit<EholeFingerprint, 'id' | 'createdAt'>[]) =>
      FingerprintService.batchCreateEholeFingerprints(fingerprints),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.ehole.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 文件导入 EHole 指纹
 */
export function useImportEholeFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (file: File) => FingerprintService.importEholeFingerprints(file),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.ehole.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 批量删除 EHole 指纹
 */
export function useBulkDeleteEholeFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (ids: number[]) => FingerprintService.bulkDeleteEholeFingerprints(ids),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.ehole.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 删除所有 EHole 指纹
 */
export function useDeleteAllEholeFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => FingerprintService.deleteAllEholeFingerprints(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.ehole.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

// ==================== Goby Hooks ====================

/**
 * 获取 Goby 指纹列表
 */
export function useGobyFingerprints(
  params: { page?: number; pageSize?: number; filter?: string } = {},
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: fingerprintKeys.goby.list(params),
    queryFn: () => FingerprintService.getGobyFingerprints(params),
    ...options,
  })
}

/**
 * 获取 Goby 指纹详情
 */
export function useGobyFingerprint(id: number, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: fingerprintKeys.goby.detail(id),
    queryFn: () => FingerprintService.getGobyFingerprint(id),
    enabled: id > 0 && options?.enabled !== false,
  })
}

/**
 * 创建 Goby 指纹
 */
export function useCreateGobyFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: Omit<GobyFingerprint, 'id' | 'createdAt'>) => 
      FingerprintService.createGobyFingerprint(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.goby.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 更新 Goby 指纹
 */
export function useUpdateGobyFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<GobyFingerprint> }) =>
      FingerprintService.updateGobyFingerprint(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.goby.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.goby.detail(id) })
    },
  })
}

/**
 * 删除 Goby 指纹
 */
export function useDeleteGobyFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => FingerprintService.deleteGobyFingerprint(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.goby.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 文件导入 Goby 指纹
 */
export function useImportGobyFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (file: File) => FingerprintService.importGobyFingerprints(file),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.goby.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 批量删除 Goby 指纹
 */
export function useBulkDeleteGobyFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (ids: number[]) => FingerprintService.bulkDeleteGobyFingerprints(ids),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.goby.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 删除所有 Goby 指纹
 */
export function useDeleteAllGobyFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => FingerprintService.deleteAllGobyFingerprints(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.goby.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

// ==================== Wappalyzer Hooks ====================

/**
 * 获取 Wappalyzer 指纹列表
 */
export function useWappalyzerFingerprints(
  params: { page?: number; pageSize?: number; filter?: string } = {},
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: fingerprintKeys.wappalyzer.list(params),
    queryFn: () => FingerprintService.getWappalyzerFingerprints(params),
    ...options,
  })
}

/**
 * 获取 Wappalyzer 指纹详情
 */
export function useWappalyzerFingerprint(id: number, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: fingerprintKeys.wappalyzer.detail(id),
    queryFn: () => FingerprintService.getWappalyzerFingerprint(id),
    enabled: id > 0 && options?.enabled !== false,
  })
}

/**
 * 创建 Wappalyzer 指纹
 */
export function useCreateWappalyzerFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: Omit<WappalyzerFingerprint, 'id' | 'createdAt'>) => 
      FingerprintService.createWappalyzerFingerprint(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.wappalyzer.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 更新 Wappalyzer 指纹
 */
export function useUpdateWappalyzerFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<WappalyzerFingerprint> }) =>
      FingerprintService.updateWappalyzerFingerprint(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.wappalyzer.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.wappalyzer.detail(id) })
    },
  })
}

/**
 * 删除 Wappalyzer 指纹
 */
export function useDeleteWappalyzerFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => FingerprintService.deleteWappalyzerFingerprint(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.wappalyzer.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 文件导入 Wappalyzer 指纹
 */
export function useImportWappalyzerFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (file: File) => FingerprintService.importWappalyzerFingerprints(file),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.wappalyzer.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 批量删除 Wappalyzer 指纹
 */
export function useBulkDeleteWappalyzerFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (ids: number[]) => FingerprintService.bulkDeleteWappalyzerFingerprints(ids),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.wappalyzer.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 删除所有 Wappalyzer 指纹
 */
export function useDeleteAllWappalyzerFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => FingerprintService.deleteAllWappalyzerFingerprints(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.wappalyzer.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

// ==================== Fingers Hooks ====================

/**
 * 获取 Fingers 指纹列表
 */
export function useFingersFingerprints(
  params: { page?: number; pageSize?: number; filter?: string } = {},
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: fingerprintKeys.fingers.list(params),
    queryFn: () => FingerprintService.getFingersFingerprints(params),
    ...options,
  })
}

/**
 * 获取 Fingers 指纹详情
 */
export function useFingersFingerprint(id: number, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: fingerprintKeys.fingers.detail(id),
    queryFn: () => FingerprintService.getFingersFingerprint(id),
    enabled: id > 0 && options?.enabled !== false,
  })
}

/**
 * 创建 Fingers 指纹
 */
export function useCreateFingersFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: Omit<FingersFingerprint, 'id' | 'createdAt'>) => 
      FingerprintService.createFingersFingerprint(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingers.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 更新 Fingers 指纹
 */
export function useUpdateFingersFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<FingersFingerprint> }) =>
      FingerprintService.updateFingersFingerprint(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingers.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingers.detail(id) })
    },
  })
}

/**
 * 删除 Fingers 指纹
 */
export function useDeleteFingersFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => FingerprintService.deleteFingersFingerprint(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingers.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 文件导入 Fingers 指纹
 */
export function useImportFingersFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (file: File) => FingerprintService.importFingersFingerprints(file),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingers.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 批量删除 Fingers 指纹
 */
export function useBulkDeleteFingersFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (ids: number[]) => FingerprintService.bulkDeleteFingersFingerprints(ids),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingers.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 删除所有 Fingers 指纹
 */
export function useDeleteAllFingersFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => FingerprintService.deleteAllFingersFingerprints(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingers.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

// ==================== FingerPrintHub Hooks ====================

/**
 * 获取 FingerPrintHub 指纹列表
 */
export function useFingerPrintHubFingerprints(
  params: { page?: number; pageSize?: number; filter?: string } = {},
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: fingerprintKeys.fingerprinthub.list(params),
    queryFn: () => FingerprintService.getFingerPrintHubFingerprints(params),
    ...options,
  })
}

/**
 * 获取 FingerPrintHub 指纹详情
 */
export function useFingerPrintHubFingerprint(id: number, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: fingerprintKeys.fingerprinthub.detail(id),
    queryFn: () => FingerprintService.getFingerPrintHubFingerprint(id),
    enabled: id > 0 && options?.enabled !== false,
  })
}

/**
 * 创建 FingerPrintHub 指纹
 */
export function useCreateFingerPrintHubFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: Omit<FingerPrintHubFingerprint, 'id' | 'createdAt'>) => 
      FingerprintService.createFingerPrintHubFingerprint(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingerprinthub.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 更新 FingerPrintHub 指纹
 */
export function useUpdateFingerPrintHubFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<FingerPrintHubFingerprint> }) =>
      FingerprintService.updateFingerPrintHubFingerprint(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingerprinthub.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingerprinthub.detail(id) })
    },
  })
}

/**
 * 删除 FingerPrintHub 指纹
 */
export function useDeleteFingerPrintHubFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => FingerprintService.deleteFingerPrintHubFingerprint(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingerprinthub.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 文件导入 FingerPrintHub 指纹
 */
export function useImportFingerPrintHubFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (file: File) => FingerprintService.importFingerPrintHubFingerprints(file),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingerprinthub.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 批量删除 FingerPrintHub 指纹
 */
export function useBulkDeleteFingerPrintHubFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (ids: number[]) => FingerprintService.bulkDeleteFingerPrintHubFingerprints(ids),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingerprinthub.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 删除所有 FingerPrintHub 指纹
 */
export function useDeleteAllFingerPrintHubFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => FingerprintService.deleteAllFingerPrintHubFingerprints(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.fingerprinthub.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

// ==================== ARL Hooks ====================

/**
 * 获取 ARL 指纹列表
 */
export function useARLFingerprints(
  params: { page?: number; pageSize?: number; filter?: string } = {},
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: fingerprintKeys.arl.list(params),
    queryFn: () => FingerprintService.getARLFingerprints(params),
    ...options,
  })
}

/**
 * 获取 ARL 指纹详情
 */
export function useARLFingerprint(id: number, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: fingerprintKeys.arl.detail(id),
    queryFn: () => FingerprintService.getARLFingerprint(id),
    enabled: id > 0 && options?.enabled !== false,
  })
}

/**
 * 创建 ARL 指纹
 */
export function useCreateARLFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (data: Omit<ARLFingerprint, 'id' | 'createdAt'>) => 
      FingerprintService.createARLFingerprint(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.arl.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 更新 ARL 指纹
 */
export function useUpdateARLFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: Partial<ARLFingerprint> }) =>
      FingerprintService.updateARLFingerprint(id, data),
    onSuccess: (_, { id }) => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.arl.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.arl.detail(id) })
    },
  })
}

/**
 * 删除 ARL 指纹
 */
export function useDeleteARLFingerprint() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: number) => FingerprintService.deleteARLFingerprint(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.arl.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 文件导入 ARL 指纹
 */
export function useImportARLFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (file: File) => FingerprintService.importARLFingerprints(file),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.arl.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 批量删除 ARL 指纹
 */
export function useBulkDeleteARLFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (ids: number[]) => FingerprintService.bulkDeleteARLFingerprints(ids),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.arl.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

/**
 * 删除所有 ARL 指纹
 */
export function useDeleteAllARLFingerprints() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: () => FingerprintService.deleteAllARLFingerprints(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.arl.all() })
      queryClient.invalidateQueries({ queryKey: fingerprintKeys.stats() })
    },
  })
}

// ==================== 统计 Hooks ====================

/**
 * 获取指纹库统计
 */
export function useFingerprintStats() {
  return useQuery({
    queryKey: fingerprintKeys.stats(),
    queryFn: () => FingerprintService.getStats(),
  })
}
