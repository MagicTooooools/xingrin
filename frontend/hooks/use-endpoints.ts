"use client"

import { useQuery, keepPreviousData } from "@tanstack/react-query"
import { useResourceMutation } from '@/hooks/_shared/create-resource-mutation'
import { createResourceKeys } from "@/hooks/_shared/query-keys"
import { normalizePagination } from "@/hooks/_shared/pagination"
import {
  getAssetDeletedCount,
  resolveAssetBulkCreateToast,
} from "@/hooks/_shared/asset-mutation-helpers"
import { EndpointService } from "@/services/endpoint.service"
import type { 
  Endpoint, 
  CreateEndpointRequest,
  GetEndpointsRequest,
  GetEndpointsResponse,
  BatchDeleteEndpointsRequest
} from "@/types/endpoint.types"

type EndpointPageResponse = {
  results?: Endpoint[]
  total?: number
  page?: number
  pageSize?: number
  totalPages?: number
  page_size?: number
  total_pages?: number
}

// Query Keys
const endpointKeyBase = createResourceKeys("endpoints", {
  list: (params: GetEndpointsRequest) => params,
  detail: (id: number) => id,
})

export const endpointKeys = {
  ...endpointKeyBase,
  byTarget: (targetId: number, params: GetEndpointsRequest) =>
    [...endpointKeyBase.all, 'target', targetId, params] as const,
  bySubdomain: (subdomainId: number, params: GetEndpointsRequest) =>
    [...endpointKeyBase.all, 'subdomain', subdomainId, params] as const,
  byScan: (scanId: number, params: GetEndpointsRequest) =>
    [...endpointKeyBase.all, 'scan', scanId, params] as const,
}

function endpointAllInvalidates() {
  return [{ queryKey: endpointKeys.all }]
}

// 获取单个 Endpoint 详情
export function useEndpoint(id: number) {
  return useQuery({
    queryKey: endpointKeys.detail(id),
    queryFn: () => EndpointService.getEndpointById(id),
    select: (response) => {
      // RESTful 标准：直接返回数据
      return response as Endpoint
    },
    enabled: !!id,
  })
}

// 获取 Endpoint 列表
export function useEndpoints(params?: GetEndpointsRequest) {
  const defaultParams: GetEndpointsRequest = {
    page: 1,
    pageSize: 10,
    ...params
  }
  
  return useQuery({
    queryKey: endpointKeys.list(defaultParams),
    queryFn: () => EndpointService.getEndpoints(defaultParams),
    select: (response) => {
      // RESTful 标准：直接返回数据
      return response as GetEndpointsResponse
    },
  })
}

// 根据目标ID获取 Endpoint 列表（使用专用路由）
export function useEndpointsByTarget(targetId: number, params?: Omit<GetEndpointsRequest, 'targetId'>, filter?: string) {
  const defaultParams: GetEndpointsRequest = {
    page: 1,
    pageSize: 10,
    ...params
  }
  
  return useQuery({
    queryKey: [...endpointKeys.byTarget(targetId, defaultParams), filter],
    queryFn: () => EndpointService.getEndpointsByTargetId(targetId, defaultParams, filter),
    select: (response) => {
      // RESTful 标准：直接返回数据
      return response as GetEndpointsResponse
    },
    enabled: !!targetId,
    placeholderData: keepPreviousData,
  })
}

// 根据扫描ID获取 Endpoint 列表（历史快照）
export function useScanEndpoints(scanId: number, params?: Omit<GetEndpointsRequest, 'targetId'>, options?: { enabled?: boolean }, filter?: string) {
  const defaultParams: GetEndpointsRequest = {
    page: 1,
    pageSize: 10,
    ...params,
  }

  return useQuery({
    queryKey: [...endpointKeys.byScan(scanId, defaultParams), filter],
    queryFn: () => EndpointService.getEndpointsByScanId(scanId, defaultParams, filter),
    enabled: options?.enabled !== undefined ? options.enabled : !!scanId,
    select: (response: EndpointPageResponse) => {
      // 后端使用通用分页格式：results/total/page/pageSize/totalPages
      return {
        endpoints: response.results || [],
        pagination: normalizePagination(response, defaultParams.page, defaultParams.pageSize),
      }
    },
    placeholderData: keepPreviousData,
  })
}

// 创建 Endpoint（完全自动化）
export function useCreateEndpoint() {
  return useResourceMutation({
    mutationFn: (data: {
      endpoints: Array<CreateEndpointRequest>
    }) => EndpointService.createEndpoints(data),
    loadingToast: {
      key: 'common.status.creating',
      params: {},
      id: 'create-endpoint',
    },
    invalidate: endpointAllInvalidates(),
    onSuccess: ({ data, toast }) => {
      const { createdCount, existedCount } = data

      if (existedCount > 0) {
        toast.warning('toast.asset.endpoint.create.partialSuccess', {
          success: createdCount,
          skipped: existedCount
        })
      } else {
        toast.success('toast.asset.endpoint.create.success', { count: createdCount })
      }
    },
    errorFallbackKey: 'toast.asset.endpoint.create.error',
  })
}

// 删除单个 Endpoint
export function useDeleteEndpoint() {
  return useResourceMutation({
    mutationFn: (id: number) => EndpointService.deleteEndpoint(id),
    loadingToast: {
      key: 'common.status.deleting',
      params: {},
      id: (id) => `delete-endpoint-${id}`,
    },
    invalidate: endpointAllInvalidates(),
    onSuccess: ({ toast }) => {
      toast.success('toast.asset.endpoint.delete.success')
    },
    errorFallbackKey: 'toast.asset.endpoint.delete.error',
  })
}

// 批量删除 Endpoint
export function useBatchDeleteEndpoints() {
  return useResourceMutation({
    mutationFn: (data: BatchDeleteEndpointsRequest) => EndpointService.batchDeleteEndpoints(data),
    loadingToast: {
      key: 'common.status.batchDeleting',
      params: {},
      id: 'batch-delete-endpoints',
    },
    invalidate: endpointAllInvalidates(),
    onSuccess: ({ data, toast }) => {
      toast.success('toast.asset.endpoint.delete.bulkSuccess', {
        count: getAssetDeletedCount(data),
      })
    },
    errorFallbackKey: 'toast.asset.endpoint.delete.error',
  })
}

// 批量创建端点（绑定到目标）
export function useBulkCreateEndpoints() {
  return useResourceMutation({
    mutationFn: (data: { targetId: number; urls: string[] }) =>
      EndpointService.bulkCreateEndpoints(data.targetId, data.urls),
    loadingToast: {
      key: 'common.status.batchCreating',
      params: {},
      id: 'bulk-create-endpoints',
    },
    invalidate: [
      ({ variables }) => ({
        queryKey: endpointKeys.byTarget(variables.targetId, {}),
        exact: false,
        refetchType: 'active',
      }),
      {
        queryKey: endpointKeys.all,
        exact: false,
        refetchType: 'active',
      },
    ],
    onSuccess: ({ data, toast }) => {
      const toastPayload = resolveAssetBulkCreateToast(data.createdCount, {
        success: 'toast.asset.endpoint.create.success',
        partial: 'toast.asset.endpoint.create.partialSuccess',
      })
      if (toastPayload.variant === 'success') {
        toast.success(toastPayload.key, toastPayload.params)
      } else {
        toast.warning(toastPayload.key, toastPayload.params)
      }
    },
    errorFallbackKey: 'toast.asset.endpoint.create.error',
  })
}
