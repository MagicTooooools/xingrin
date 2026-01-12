import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query'
import { useToastMessages } from '@/lib/toast-helpers'
import { getErrorCode } from '@/lib/response-parser'
import { WebsiteService } from '@/services/website.service'
import { api } from '@/lib/api-client'
import type { WebSite, WebSiteListResponse } from '@/types/website.types'

// API 服务函数
const websiteService = {
  // 获取目标的网站列表
  getTargetWebSites: async (
    targetId: number,
    params: { page: number; pageSize: number; filter?: string }
  ): Promise<WebSiteListResponse> => {
    const response = await api.get<WebSiteListResponse>(
      `/targets/${targetId}/websites/`,
      { params }
    )
    return response.data
  },

  // 获取扫描的网站列表
  getScanWebSites: async (
    scanId: number,
    params: { page: number; pageSize: number; filter?: string }
  ): Promise<WebSiteListResponse> => {
    const response = await api.get<WebSiteListResponse>(
      `/scans/${scanId}/websites/`,
      { params }
    )
    return response.data
  },

  // 批量删除网站（支持单个或多个）
  bulkDeleteWebSites: async (ids: number[]): Promise<{
    message: string
    deletedCount: number
    requestedIds: number[]
    cascadeDeleted: Record<string, number>
  }> => {
    const response = await api.post('/websites/bulk-delete/', { ids })
    return response.data
  },

  // 删除单个网站（使用单独的 DELETE API）
  deleteWebSite: async (websiteId: number): Promise<{
    message: string
    websiteId: number
    websiteUrl: string
    deletedCount: number
    deletedWebSites: string[]
    detail: {
      phase1: string
      phase2: string
    }
  }> => {
    const response = await api.delete(`/websites/${websiteId}/`)
    return response.data
  },
}

// 获取目标的网站列表
export function useTargetWebSites(
  targetId: number,
  params: { page: number; pageSize: number; filter?: string },
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: ['target-websites', targetId, params],
    queryFn: () => websiteService.getTargetWebSites(targetId, params),
    enabled: options?.enabled ?? true,
    placeholderData: keepPreviousData,
  })
}

// 获取扫描的网站列表
export function useScanWebSites(
  scanId: number,
  params: { page: number; pageSize: number; filter?: string },
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: ['scan-websites', scanId, params],
    queryFn: () => websiteService.getScanWebSites(scanId, params),
    enabled: options?.enabled ?? true,
    placeholderData: keepPreviousData,
  })
}

// 删除单个网站（使用单独的 DELETE API）
export function useDeleteWebSite() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: websiteService.deleteWebSite,
    onMutate: (id) => {
      toastMessages.loading('common.status.deleting', {}, `delete-website-${id}`)
    },
    onSuccess: (response, id) => {
      toastMessages.dismiss(`delete-website-${id}`)
      toastMessages.success('toast.asset.website.delete.success')
      
      queryClient.invalidateQueries({ queryKey: ['target-websites'] })
      queryClient.invalidateQueries({ queryKey: ['scan-websites'] })
      queryClient.invalidateQueries({ queryKey: ['targets'] })
      queryClient.invalidateQueries({ queryKey: ['scans'] })
    },
    onError: (error: any, id) => {
      toastMessages.dismiss(`delete-website-${id}`)
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.asset.website.delete.error')
    },
  })
}

// 批量删除网站（使用统一的批量删除接口）
export function useBulkDeleteWebSites() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: websiteService.bulkDeleteWebSites,
    onMutate: () => {
      toastMessages.loading('common.status.batchDeleting', {}, 'bulk-delete-websites')
    },
    onSuccess: (response) => {
      toastMessages.dismiss('bulk-delete-websites')
      toastMessages.success('toast.asset.website.delete.bulkSuccess', { count: response.deletedCount })
      
      queryClient.invalidateQueries({ queryKey: ['target-websites'] })
      queryClient.invalidateQueries({ queryKey: ['scan-websites'] })
      queryClient.invalidateQueries({ queryKey: ['targets'] })
      queryClient.invalidateQueries({ queryKey: ['scans'] })
    },
    onError: (error: any) => {
      toastMessages.dismiss('bulk-delete-websites')
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.asset.website.delete.error')
    },
  })
}


// 批量创建网站（绑定到目标）
export function useBulkCreateWebsites() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: (data: { targetId: number; urls: string[] }) =>
      WebsiteService.bulkCreateWebsites(data.targetId, data.urls),
    onMutate: async () => {
      toastMessages.loading('common.status.batchCreating', {}, 'bulk-create-websites')
    },
    onSuccess: (response, { targetId }) => {
      toastMessages.dismiss('bulk-create-websites')
      const { createdCount } = response
      
      if (createdCount > 0) {
        toastMessages.success('toast.asset.website.create.success', { count: createdCount })
      } else {
        toastMessages.warning('toast.asset.website.create.partialSuccess', { success: 0, skipped: 0 })
      }
      
      queryClient.invalidateQueries({
        queryKey: ['target-websites', targetId],
        exact: false,
        refetchType: 'active',
      })
      queryClient.invalidateQueries({
        queryKey: ['scan-websites'],
        exact: false,
        refetchType: 'active',
      })
      queryClient.invalidateQueries({
        queryKey: ['targets', targetId],
        refetchType: 'active',
      })
    },
    onError: (error: any) => {
      toastMessages.dismiss('bulk-create-websites')
      console.error('Failed to bulk create websites:', error)
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.asset.website.create.error')
    },
  })
}
