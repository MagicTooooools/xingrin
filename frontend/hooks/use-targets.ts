/**
 * Targets Hooks - 目标管理相关 hooks
 */
import { useQuery, useMutation, useQueryClient, keepPreviousData } from '@tanstack/react-query'
import { useToastMessages } from '@/lib/toast-helpers'
import { getErrorCode } from '@/lib/response-parser'
import {
  getTargets,
  getTargetById,
  createTarget,
  updateTarget,
  deleteTarget,
  batchDeleteTargets,
  batchCreateTargets,
  getTargetOrganizations,
  linkTargetOrganizations,
  unlinkTargetOrganizations,
  getTargetEndpoints,
  getTargetBlacklist,
  updateTargetBlacklist,
} from '@/services/target.service'
import type {
  CreateTargetRequest,
  UpdateTargetRequest,
  BatchDeleteTargetsRequest,
  BatchCreateTargetsRequest,
} from '@/types/target.types'

/**
 * 获取所有目标列表
 * 支持两种调用方式：
 * 1. useTargets(page, pageSize, type, filter) - 直接传参数
 * 2. useTargets({ page, pageSize, organizationId, filter }, options) - 传对象
 */
export function useTargets(
  pageOrParams: number | { page?: number; pageSize?: number; organizationId?: number; filter?: string } = 1,
  pageSizeOrOptions: number | { enabled?: boolean } = 10,
  type?: string,
  filter?: string
) {
  // 处理参数：支持对象参数或独立参数
  let actualPage: number
  let actualPageSize: number
  let actualOrgId: number | undefined
  let actualFilter: string | undefined
  let actualType: string | undefined
  let enabled: boolean = true

  if (typeof pageOrParams === 'object') {
    // 对象参数方式
    actualPage = pageOrParams.page || 1
    actualPageSize = pageOrParams.pageSize || 10
    actualOrgId = pageOrParams.organizationId
    actualFilter = pageOrParams.filter
    actualType = undefined
    // 第二个参数是 options
    if (typeof pageSizeOrOptions === 'object') {
      enabled = pageSizeOrOptions.enabled !== false
    }
  } else {
    // 独立参数方式
    actualPage = pageOrParams
    actualPageSize = typeof pageSizeOrOptions === 'number' ? pageSizeOrOptions : 10
    actualOrgId = undefined
    actualFilter = filter
    actualType = type
  }

  return useQuery({
    queryKey: ['targets', { page: actualPage, pageSize: actualPageSize, organizationId: actualOrgId, filter: actualFilter, type: actualType }],
    queryFn: () => getTargets(actualPage, actualPageSize, actualFilter, actualType),
    enabled,
    select: (response) => {
      // 如果指定了 organizationId，过滤结果
      if (actualOrgId) {
        const filteredResults = response.results.filter(target => 
          target.organizations?.some(org => org.id === actualOrgId)
        )
        return {
          ...response,
          results: filteredResults,
          total: filteredResults.length,
          // 为兼容性添加额外字段
          count: filteredResults.length,  // 兼容字段
          targets: filteredResults,
          page: actualPage,
          pageSize: actualPageSize,
          totalPages: Math.ceil(filteredResults.length / actualPageSize),
        }
      }
      
      // 否则直接返回原始响应，并添加兼容字段
      return {
        ...response,
        targets: response.results,
        // 后端返回 total，不是 count
        count: response.total,  // 兼容字段，使用 total 值
        // 保持原有字段
        total: response.total,
        page: response.page,
        pageSize: response.pageSize,
        totalPages: response.totalPages,
      }
    },
    placeholderData: keepPreviousData,
  })
}

/**
 * 获取单个目标详情
 */
export function useTarget(id: number, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: ['targets', id],
    queryFn: () => getTargetById(id),
    enabled: options?.enabled !== undefined ? options.enabled : !!id,
  })
}

/**
 * 创建目标
 */
export function useCreateTarget() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: (data: CreateTargetRequest) => createTarget(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['targets'] })
      toastMessages.success('toast.target.create.success')
    },
    onError: (error: any) => {
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.target.create.error')
    },
  })
}

/**
 * 更新目标
 */
export function useUpdateTarget() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: UpdateTargetRequest }) =>
      updateTarget(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['targets'] })
      queryClient.invalidateQueries({ queryKey: ['targets', variables.id] })
      toastMessages.success('toast.target.update.success')
    },
    onError: (error: any) => {
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.target.update.error')
    },
  })
}

/**
 * 删除目标（RESTful 204 No Content）
 */
export function useDeleteTarget() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: ({ id, name }: { id: number; name: string }) => deleteTarget(id),
    onMutate: ({ id }) => {
      toastMessages.loading('common.status.deleting', {}, `delete-target-${id}`)
    },
    onSuccess: (_response, { id, name }) => {
      toastMessages.dismiss(`delete-target-${id}`)
      toastMessages.success('toast.target.delete.success', { name })
      
      queryClient.invalidateQueries({ queryKey: ['targets'] })
      queryClient.invalidateQueries({ queryKey: ['organizations'] })
    },
    onError: (error: any, { id }) => {
      toastMessages.dismiss(`delete-target-${id}`)
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.target.delete.error')
    },
  })
}

/**
 * 批量删除目标
 */
export function useBatchDeleteTargets() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: (data: BatchDeleteTargetsRequest) => batchDeleteTargets(data),
    onSuccess: (response) => {
      queryClient.invalidateQueries({ queryKey: ['targets'] })
      toastMessages.success('toast.target.delete.bulkSuccess', { count: response.deletedCount })
    },
    onError: (error: any) => {
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.target.delete.error')
    },
  })
}

/**
 * 批量创建目标
 */
export function useBatchCreateTargets() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: (data: BatchCreateTargetsRequest) => batchCreateTargets(data),
    onSuccess: (response) => {
      queryClient.invalidateQueries({ queryKey: ['targets'] })
      queryClient.invalidateQueries({ queryKey: ['organizations'] })
      toastMessages.success('toast.target.create.bulkSuccess', { count: response.createdCount || 0 })
    },
    onError: (error: any) => {
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.target.create.error')
    },
  })
}

/**
 * 获取目标的组织列表
 */
export function useTargetOrganizations(targetId: number, page = 1, pageSize = 10) {
  return useQuery({
    queryKey: ['targets', targetId, 'organizations', page, pageSize],
    queryFn: () => getTargetOrganizations(targetId, page, pageSize),
    enabled: !!targetId,
  })
}

/**
 * 关联目标与组织
 */
export function useLinkTargetOrganizations() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: ({ targetId, organizationIds }: { targetId: number; organizationIds: number[] }) =>
      linkTargetOrganizations(targetId, organizationIds),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['targets', variables.targetId, 'organizations'] })
      queryClient.invalidateQueries({ queryKey: ['targets', variables.targetId] })
      toastMessages.success('toast.target.link.success')
    },
    onError: (error: any) => {
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.target.link.error')
    },
  })
}

/**
 * 取消目标与组织的关联
 */
export function useUnlinkTargetOrganizations() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: ({ targetId, organizationIds }: { targetId: number; organizationIds: number[] }) =>
      unlinkTargetOrganizations(targetId, organizationIds),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['targets', variables.targetId, 'organizations'] })
      queryClient.invalidateQueries({ queryKey: ['targets', variables.targetId] })
      toastMessages.success('toast.target.unlink.success')
    },
    onError: (error: any) => {
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.target.unlink.error')
    },
  })
}

/**
 * 获取目标的端点列表
 */
export function useTargetEndpoints(
  targetId: number,
  params?: {
    page?: number
    pageSize?: number
    filter?: string
  },
  options?: {
    enabled?: boolean
  }
) {
  return useQuery({
    queryKey: ['targets', 'detail', targetId, 'endpoints', {
      page: params?.page,
      pageSize: params?.pageSize,
      filter: params?.filter,
    }],
    queryFn: () => getTargetEndpoints(targetId, params?.page || 1, params?.pageSize || 10, params?.filter),
    enabled: options?.enabled !== undefined ? options.enabled : !!targetId,
    select: (response: any) => {
      // 后端使用通用分页格式：results/total/page/pageSize/totalPages
      return {
        endpoints: response.results || response.endpoints || [],
        pagination: {
          total: response.total || 0,
          page: response.page || 1,
          pageSize: response.pageSize || response.page_size || 10,
          totalPages: response.totalPages || response.total_pages || 0,
        }
      }
    },
  })
}

/**
 * 获取目标的黑名单规则
 */
export function useTargetBlacklist(targetId: number) {
  return useQuery({
    queryKey: ['targets', targetId, 'blacklist'],
    queryFn: () => getTargetBlacklist(targetId),
    enabled: !!targetId,
  })
}

/**
 * 更新目标的黑名单规则
 */
export function useUpdateTargetBlacklist() {
  const queryClient = useQueryClient()
  const toastMessages = useToastMessages()

  return useMutation({
    mutationFn: ({ targetId, patterns }: { targetId: number; patterns: string[] }) =>
      updateTargetBlacklist(targetId, patterns),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['targets', variables.targetId, 'blacklist'] })
      toastMessages.success('toast.blacklist.save.success')
    },
    onError: (error: any) => {
      toastMessages.errorFromCode(getErrorCode(error?.response?.data), 'toast.blacklist.save.error')
    },
  })
}

