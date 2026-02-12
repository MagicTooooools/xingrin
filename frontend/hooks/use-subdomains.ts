"use client"

import { useQuery, keepPreviousData } from "@tanstack/react-query"
import { useResourceMutation } from '@/hooks/_shared/create-resource-mutation'
import { createResourceKeys } from "@/hooks/_shared/query-keys"
import { normalizePagination } from "@/hooks/_shared/pagination"
import {
  getSubdomainBatchDeleteCount,
  getSubdomainBatchDeleteFromOrgCount,
  resolveSubdomainCreateToast,
} from "@/hooks/_shared/subdomain-mutation-helpers"
import { SubdomainService } from "@/services/subdomain.service"
import { OrganizationService } from "@/services/organization.service"
import type { GetAllSubdomainsParams } from "@/types/subdomain.types"
import type { PaginationParams } from "@/types/common.types"

// Query Keys
export const subdomainKeys = createResourceKeys("subdomains", {
  list: (params: PaginationParams & { organizationId?: string }) => params,
  detail: (id: number) => id,
})

function subdomainCascadeInvalidates() {
  return [
    { queryKey: subdomainKeys.all },
    { queryKey: ['targets'] },
    { queryKey: ['scans'] },
    { queryKey: ['organizations'] },
  ]
}

// 获取单个子域名详情
export function useSubdomain(id: number) {
  return useQuery({
    queryKey: subdomainKeys.detail(id),
    queryFn: () => SubdomainService.getSubdomainById(id),
    enabled: !!id,
  })
}

// 获取组织的子域名列表
export function useOrganizationSubdomains(
  organizationId: number,
  params?: { page?: number; pageSize?: number },
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: ['organizations', 'detail', organizationId, 'subdomains', {
      page: params?.page,
      pageSize: params?.pageSize,
    }],
    queryFn: () => SubdomainService.getSubdomainsByOrgId(organizationId, {
      page: params?.page || 1,
      pageSize: params?.pageSize || 10,
    }),
    enabled: options?.enabled !== undefined ? options.enabled : true,
    select: (response) => ({
      domains: response.domains || [],
      pagination: normalizePagination(response, params?.page ?? 1, params?.pageSize ?? 10),
    }),
  })
}

// 创建子域名（绑定到资产）
export function useCreateSubdomain() {
  return useResourceMutation({
    mutationFn: (data: { domains: Array<{ name: string }>; assetId: number }) =>
      SubdomainService.createSubdomains(data),
    loadingToast: {
      key: 'common.status.creating',
      params: {},
      id: 'create-subdomain',
    },
    invalidate: [
      { queryKey: subdomainKeys.all },
      { queryKey: ['assets'] },
    ],
    onSuccess: ({ data, toast }) => {
      const toastPayload = resolveSubdomainCreateToast(data)
      if (toastPayload.variant === "warning") {
        toast.warning(toastPayload.key, toastPayload.params)
      } else {
        toast.success(toastPayload.key, toastPayload.params)
      }
    },
    errorFallbackKey: 'toast.asset.subdomain.create.error',
  })
}

// 从组织中移除子域名
export function useDeleteSubdomainFromOrganization() {
  return useResourceMutation({
    mutationFn: (data: { organizationId: number; targetId: number }) =>
      OrganizationService.unlinkTargetsFromOrganization({
        organizationId: data.organizationId,
        targetIds: [data.targetId],
      }),
    loadingToast: {
      key: 'common.status.removing',
      params: {},
      id: ({ organizationId, targetId }) => `delete-${organizationId}-${targetId}`,
    },
    invalidate: [
      { queryKey: subdomainKeys.all },
      { queryKey: ['organizations'] },
    ],
    onSuccess: ({ toast }) => {
      toast.success('toast.asset.subdomain.delete.success')
    },
    errorFallbackKey: 'toast.asset.subdomain.delete.error',
  })
}

// 批量从组织中移除子域名
export function useBatchDeleteSubdomainsFromOrganization() {
  return useResourceMutation({
    mutationFn: (data: { organizationId: number; domainIds: number[] }) => 
      SubdomainService.batchDeleteSubdomainsFromOrganization(data),
    loadingToast: {
      key: 'common.status.batchRemoving',
      params: {},
      id: ({ organizationId }) => `batch-delete-${organizationId}`,
    },
    invalidate: [
      { queryKey: subdomainKeys.all },
      { queryKey: ['organizations'] },
    ],
    onSuccess: ({ data, toast }) => {
      const successCount = getSubdomainBatchDeleteFromOrgCount(data)
      toast.success('toast.asset.subdomain.delete.bulkSuccess', { count: successCount })
    },
    errorFallbackKey: 'toast.asset.subdomain.delete.error',
  })
}

// 删除单个子域名（使用单独的 DELETE API）
export function useDeleteSubdomain() {
  return useResourceMutation({
    mutationFn: (id: number) => SubdomainService.deleteSubdomain(id),
    loadingToast: {
      key: 'common.status.deleting',
      params: {},
      id: (id) => `delete-subdomain-${id}`,
    },
    invalidate: subdomainCascadeInvalidates(),
    onSuccess: ({ toast }) => {
      toast.success('toast.asset.subdomain.delete.success')
    },
    errorFallbackKey: 'toast.asset.subdomain.delete.error',
  })
}

// 批量删除子域名（使用统一的批量删除接口）
export function useBatchDeleteSubdomains() {
  return useResourceMutation({
    mutationFn: (ids: number[]) => SubdomainService.batchDeleteSubdomains(ids),
    loadingToast: {
      key: 'common.status.batchDeleting',
      params: {},
      id: 'batch-delete-subdomains',
    },
    invalidate: subdomainCascadeInvalidates(),
    onSuccess: ({ data, toast }) => {
      toast.success('toast.asset.subdomain.delete.bulkSuccess', {
        count: getSubdomainBatchDeleteCount(data),
      })
    },
    errorFallbackKey: 'toast.asset.subdomain.delete.error',
  })
}

// 更新子域名
export function useUpdateSubdomain() {
  return useResourceMutation({
    mutationFn: ({ id, data }: { id: number; data: { name?: string; description?: string } }) =>
      SubdomainService.updateSubdomain({ id, ...data }),
    loadingToast: {
      key: 'common.status.updating',
      params: {},
      id: ({ id }) => `update-subdomain-${id}`,
    },
    invalidate: [
      { queryKey: subdomainKeys.all },
      { queryKey: ['organizations'] },
    ],
    onSuccess: ({ toast }) => {
      toast.success('common.status.updateSuccess')
    },
    errorFallbackKey: 'common.status.updateFailed',
  })
}

// 获取所有子域名列表
export function useAllSubdomains(
  params: GetAllSubdomainsParams = {},
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: ['subdomains', 'all', { page: params.page, pageSize: params.pageSize }],
    queryFn: () => SubdomainService.getAllSubdomains(params),
    select: (response) => ({
      domains: response.domains || [],
      pagination: normalizePagination(response, params.page ?? 1, params.pageSize ?? 10),
    }),
    enabled: options?.enabled !== undefined ? options.enabled : true,
  })
}

// 获取目标的子域名列表
export function useTargetSubdomains(
  targetId: number,
  params?: { page?: number; pageSize?: number; filter?: string },
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: ['targets', targetId, 'subdomains', { page: params?.page, pageSize: params?.pageSize, filter: params?.filter }],
    queryFn: () => SubdomainService.getSubdomainsByTargetId(targetId, params),
    enabled: options?.enabled !== undefined ? options.enabled : !!targetId,
    placeholderData: keepPreviousData,
  })
}

// 获取扫描的子域名列表
export function useScanSubdomains(
  scanId: number,
  params?: { page?: number; pageSize?: number; filter?: string },
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: ['scans', scanId, 'subdomains', { page: params?.page, pageSize: params?.pageSize, filter: params?.filter }],
    queryFn: () => SubdomainService.getSubdomainsByScanId(scanId, params),
    enabled: options?.enabled !== undefined ? options.enabled : !!scanId,
    placeholderData: keepPreviousData,
  })
}

// 批量创建子域名（绑定到目标）
export function useBulkCreateSubdomains() {
  return useResourceMutation({
    mutationFn: (data: { targetId: number; subdomains: string[] }) =>
      SubdomainService.bulkCreateSubdomains(data.targetId, data.subdomains),
    loadingToast: {
      key: 'common.status.batchCreating',
      params: {},
      id: 'bulk-create-subdomains',
    },
    invalidate: [
      ({ variables }) => ({
        queryKey: ['targets', variables.targetId, 'subdomains'],
        exact: false,
        refetchType: 'active',
      }),
      {
        queryKey: subdomainKeys.all,
        exact: false,
        refetchType: 'active',
      },
    ],
    onSuccess: ({ data, toast }) => {
      const { createdCount, skippedCount = 0, invalidCount = 0, mismatchedCount = 0 } = data
      const totalSkipped = skippedCount + invalidCount + mismatchedCount

      if (totalSkipped > 0) {
        toast.warning('toast.asset.subdomain.create.partialSuccess', {
          success: createdCount,
          skipped: totalSkipped
        })
      } else if (createdCount > 0) {
        toast.success('toast.asset.subdomain.create.success', { count: createdCount })
      } else {
        toast.warning('toast.asset.subdomain.create.partialSuccess', {
          success: 0,
          skipped: 0
        })
      }
    },
    errorFallbackKey: 'toast.asset.subdomain.create.error',
  })
}
