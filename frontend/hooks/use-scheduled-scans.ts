import { useQuery, keepPreviousData } from '@tanstack/react-query'
import {
  getScheduledScans,
  getScheduledScan,
  createScheduledScan,
  updateScheduledScan,
  deleteScheduledScan,
  toggleScheduledScan,
} from '@/services/scheduled-scan.service'
import { useResourceMutation } from '@/hooks/_shared/create-resource-mutation'
import { createResourceKeys } from "@/hooks/_shared/query-keys"
import { handleScheduledScanMutationSuccess } from '@/hooks/_shared/scheduled-scan-mutation-helpers'
import { getErrorCode, getErrorResponseData } from '@/lib/response-parser'
import type {
  CreateScheduledScanRequest,
  UpdateScheduledScanRequest,
  GetScheduledScansResponse,
  ScheduledScan,
} from '@/types/scheduled-scan.types'

// Query Keys
export const scheduledScanKeys = createResourceKeys("scheduled-scans", {
  list: (params: {
    page?: number
    pageSize?: number
    search?: string
    targetId?: number
    organizationId?: number
  }) => params,
  detail: (id: number) => id,
})

/**
 * 获取定时扫描列表
 */
export function useScheduledScans(params: {
  page?: number
  pageSize?: number
  search?: string
  targetId?: number
  organizationId?: number
} = { page: 1, pageSize: 10 }) {
  return useQuery({
    queryKey: scheduledScanKeys.list(params),
    queryFn: () => getScheduledScans(params),
    placeholderData: keepPreviousData,
  })
}

/**
 * 获取定时扫描详情
 */
export function useScheduledScan(id: number) {
  return useQuery({
    queryKey: scheduledScanKeys.detail(id),
    queryFn: () => getScheduledScan(id),
    enabled: !!id,
  })
}

/**
 * 创建定时扫描
 */
export function useCreateScheduledScan() {
  return useResourceMutation({
    mutationFn: (data: CreateScheduledScanRequest) => createScheduledScan(data),
    invalidate: [{ queryKey: scheduledScanKeys.all }],
    onSuccess: ({ data: response, toast }) => {
      handleScheduledScanMutationSuccess({
        response,
        onSuccess: () => {
          // 使用 i18n 消息显示成功提示
          toast.success('toast.scheduledScan.create.success')
        },
      })
    },
    errorFallbackKey: 'toast.scheduledScan.create.error',
  })
}

/**
 * 更新定时扫描
 */
export function useUpdateScheduledScan() {
  return useResourceMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateScheduledScanRequest }) =>
      updateScheduledScan(id, data),
    invalidate: [
      { queryKey: scheduledScanKeys.all },
      { queryKey: scheduledScanKeys.details() },
    ],
    onSuccess: ({ data: response, toast }) => {
      handleScheduledScanMutationSuccess({
        response,
        onSuccess: () => {
          // 使用 i18n 消息显示成功提示
          toast.success('toast.scheduledScan.update.success')
        },
      })
    },
    errorFallbackKey: 'toast.scheduledScan.update.error',
  })
}

/**
 * 删除定时扫描
 */
export function useDeleteScheduledScan() {
  return useResourceMutation({
    mutationFn: (id: number) => deleteScheduledScan(id),
    invalidate: [{ queryKey: scheduledScanKeys.all }],
    onSuccess: ({ data: response, toast }) => {
      handleScheduledScanMutationSuccess({
        response,
        onSuccess: () => {
          // 使用 i18n 消息显示成功提示
          toast.success('toast.scheduledScan.delete.success')
        },
      })
    },
    errorFallbackKey: 'toast.scheduledScan.delete.error',
  })
}

/**
 * 切换定时扫描启用状态
 * 使用乐观更新，避免重新获取数据导致列表重新排序
 */
export function useToggleScheduledScan() {
  return useResourceMutation({
    mutationFn: ({ id, isEnabled }: { id: number; isEnabled: boolean }) =>
      toggleScheduledScan(id, isEnabled),
    onMutate: async ({ id, isEnabled }, context) => {
      const { queryClient } = context
      // 取消正在进行的查询
      await queryClient.cancelQueries({ queryKey: scheduledScanKeys.all })

      // 获取当前缓存的所有 scheduled-scans 查询
      const previousQueries = queryClient.getQueriesData({ queryKey: scheduledScanKeys.all })

      // 乐观更新所有匹配的查询缓存
      queryClient.setQueriesData(
        { queryKey: scheduledScanKeys.all },
        (old: GetScheduledScansResponse | undefined) => {
          if (!old?.results) return old
          return {
            ...old,
            results: old.results.map((item: ScheduledScan) =>
              item.id === id ? { ...item, isEnabled } : item
            ),
          }
        }
      )

      // 返回上下文用于回滚
      return { previousQueries }
    },
    onSuccess: ({ data: response, variables: { isEnabled }, toast }) => {
      handleScheduledScanMutationSuccess({
        response,
        onSuccess: () => {
          // 使用 i18n 消息显示成功提示
          if (isEnabled) {
            toast.success('toast.scheduledScan.toggle.enabled')
          } else {
            toast.success('toast.scheduledScan.toggle.disabled')
          }
        },
      })
      // 不调用 invalidateQueries，保持当前排序
    },
    onError: ({ error, context, toast, queryClient }) => {
      // 回滚到之前的状态
      if (context?.previousQueries) {
        context.previousQueries.forEach(([queryKey, data]) => {
          queryClient.setQueryData(queryKey, data)
        })
      }
      const errorCode = getErrorCode(getErrorResponseData(error))
      if (errorCode) {
        toast.errorFromCode(errorCode)
      } else {
        toast.error('toast.scheduledScan.toggle.error')
      }
    },
  })
}
