import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query'
import { useToastMessages } from '@/lib/toast-helpers'
import { getErrorCode } from '@/lib/response-parser'
import { DirectoryService } from '@/services/directory.service'
import { api } from '@/lib/api-client'
import type { Directory, DirectoryListResponse } from '@/types/directory.types'

// API 服务函数
const directoryService = {
  // 获取目标的目录列表
  getTargetDirectories: async (
    targetId: number,
    params: { page: number; pageSize: number; filter?: string }
  ): Promise<DirectoryListResponse> => {
    const response = await api.get<DirectoryListResponse>(
      `/targets/${targetId}/directories/`,
      { params }
    )
    return response.data
  },

  // 获取扫描的目录列表
  getScanDirectories: async (
    scanId: number,
    params: { page: number; pageSize: number; filter?: string }
  ): Promise<DirectoryListResponse> => {
    const response = await api.get<DirectoryListResponse>(
      `/scans/${scanId}/directories/`,
      { params }
    )
    return response.data
  },

  // 批量删除目录（支持单个或多个）
  bulkDeleteDirectories: async (ids: number[]): Promise<{
    message: string
    deletedCount: number
    requestedIds: number[]
    cascadeDeleted: Record<string, number>
  }> => {
    const response = await api.post<{
      message: string
      deletedCount: number
      requestedIds: number[]
      cascadeDeleted: Record<string, number>
    }>('/directories/bulk-delete/', { ids })
    return response.data
  },

  // 删除单个目录（使用单独的 DELETE API）
  deleteDirectory: async (directoryId: number): Promise<{
    message: string
    directoryId: number
    directoryUrl: string
    deletedCount: number
    deletedDirectories: string[]
    detail: {
      phase1: string
      phase2: string
    }
  }> => {
    const response = await api.delete<{
      message: string
      directoryId: number
      directoryUrl: string
      deletedCount: number
      deletedDirectories: string[]
      detail: {
        phase1: string
        phase2: string
      }
    }>(`/directories/${directoryId}/`)
    return response.data
  },
}

// 获取目标的目录列表
export function useTargetDirectories(
  targetId: number,
  params: { page: number; pageSize: number; filter?: string },
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: ['target-directories', targetId, params],
    queryFn: () => directoryService.getTargetDirectories(targetId, params),
    enabled: options?.enabled ?? true,
    placeholderData: keepPreviousData,
  })
}

// 获取扫描的目录列表
export function useScanDirectories(
  scanId: number,
  params: { page: number; pageSize: number; filter?: string },
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: ['scan-directories', scanId, params],
    queryFn: () => directoryService.getScanDirectories(scanId, params),
    enabled: options?.enabled ?? true,
    placeholderData: keepPreviousData,
  })
}

// 删除单个目录（使用单独的 DELETE API）
export function useDeleteDirectory() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: directoryService.deleteDirectory,
    onMutate: (id) => {
      toastMessages.loading('common.status.deleting', {}, `delete-directory-${id}`)
    },
    onSuccess: (response, id) => {
      toastMessages.dismiss(`delete-directory-${id}`)
      toastMessages.success('toast.asset.directory.delete.success')
      
      queryClient.invalidateQueries({ queryKey: ['target-directories'] })
      queryClient.invalidateQueries({ queryKey: ['scan-directories'] })
      queryClient.invalidateQueries({ queryKey: ['targets'] })
      queryClient.invalidateQueries({ queryKey: ['scans'] })
    },
    onError: (error: any, id) => {
      toastMessages.dismiss(`delete-directory-${id}`)
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.asset.directory.delete.error')
    },
  })
}

// 批量删除目录（使用统一的批量删除接口）
export function useBulkDeleteDirectories() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: directoryService.bulkDeleteDirectories,
    onMutate: () => {
      toastMessages.loading('common.status.batchDeleting', {}, 'bulk-delete-directories')
    },
    onSuccess: (response) => {
      toastMessages.dismiss('bulk-delete-directories')
      toastMessages.success('toast.asset.directory.delete.bulkSuccess', { count: response.deletedCount })
      
      queryClient.invalidateQueries({ queryKey: ['target-directories'] })
      queryClient.invalidateQueries({ queryKey: ['scan-directories'] })
      queryClient.invalidateQueries({ queryKey: ['targets'] })
      queryClient.invalidateQueries({ queryKey: ['scans'] })
    },
    onError: (error: any) => {
      toastMessages.dismiss('bulk-delete-directories')
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.asset.directory.delete.error')
    },
  })
}


// 批量创建目录（绑定到目标）
export function useBulkCreateDirectories() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: (data: { targetId: number; urls: string[] }) =>
      DirectoryService.bulkCreateDirectories(data.targetId, data.urls),
    onMutate: async () => {
      toastMessages.loading('common.status.batchCreating', {}, 'bulk-create-directories')
    },
    onSuccess: (response, { targetId }) => {
      toastMessages.dismiss('bulk-create-directories')
      const { createdCount } = response
      
      if (createdCount > 0) {
        toastMessages.success('toast.asset.directory.create.success', { count: createdCount })
      } else {
        toastMessages.warning('toast.asset.directory.create.partialSuccess', { success: 0, skipped: 0 })
      }
      
      queryClient.invalidateQueries({
        queryKey: ['target-directories', targetId],
        exact: false,
        refetchType: 'active',
      })
      queryClient.invalidateQueries({
        queryKey: ['scan-directories'],
        exact: false,
        refetchType: 'active',
      })
    },
    onError: (error: any) => {
      toastMessages.dismiss('bulk-create-directories')
      console.error('Failed to bulk create directories:', error)
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.asset.directory.create.error')
    },
  })
}
