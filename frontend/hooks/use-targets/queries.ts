import { useQuery, keepPreviousData, type UseQueryResult } from "@tanstack/react-query"
import {
  resolveTargetsQueryInput,
  selectTargetsResponse,
  type TargetSelectResponse,
  type UseTargetsOptions,
  type UseTargetsParams,
} from "@/hooks/_shared/targets-helpers"
import { normalizePagination } from "@/hooks/_shared/pagination"
import {
  getTargets,
  getTargetById,
  getTargetOrganizations,
  getTargetEndpoints,
  getTargetBlacklist,
} from "@/services/target.service"
import type { TargetsResponse } from "@/types/target.types"
import type { Endpoint } from "@/types/endpoint.types"
import { targetKeys } from "./keys"

type EndpointListResponse = {
  results?: Endpoint[]
  endpoints?: Endpoint[]
  total?: number
  page?: number
  pageSize?: number
  totalPages?: number
  page_size?: number
  total_pages?: number
}

export type UseTargetsResult = UseQueryResult<TargetSelectResponse, Error>

/**
 * 获取所有目标列表
 * 支持两种调用方式：
 * 1. useTargets(page, pageSize, type, filter) - 直接传参数
 * 2. useTargets({ page, pageSize, organizationId, filter }, options) - 传对象
 */
export function useTargets(
  params: UseTargetsParams,
  options?: UseTargetsOptions
): UseTargetsResult

export function useTargets(
  page?: number,
  pageSize?: number,
  type?: string,
  filter?: string
): UseTargetsResult

export function useTargets(
  pageOrParams: number | UseTargetsParams = 1,
  pageSizeOrOptions: number | UseTargetsOptions = 10,
  type?: string,
  filter?: string
): UseTargetsResult {
  const resolved = resolveTargetsQueryInput(
    pageOrParams,
    pageSizeOrOptions,
    type,
    filter
  )

  return useQuery<TargetsResponse, Error, TargetSelectResponse>({
    queryKey: targetKeys.list({
      page: resolved.page,
      pageSize: resolved.pageSize,
      organizationId: resolved.organizationId,
      filter: resolved.filter,
      type: resolved.type,
    }),
    queryFn: () => getTargets(resolved.page, resolved.pageSize, resolved.filter, resolved.type),
    enabled: resolved.enabled,
    select: (response) => selectTargetsResponse(response, {
      page: resolved.page,
      pageSize: resolved.pageSize,
      organizationId: resolved.organizationId,
    }),
    placeholderData: keepPreviousData,
  })
}

/**
 * 获取单个目标详情
 */
export function useTarget(id: number, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: targetKeys.detail(id),
    queryFn: () => getTargetById(id),
    enabled: options?.enabled !== undefined ? options.enabled : !!id,
  })
}

/**
 * 获取目标的组织列表
 */
export function useTargetOrganizations(targetId: number, page = 1, pageSize = 10) {
  return useQuery({
    queryKey: targetKeys.organizations(targetId, page, pageSize),
    queryFn: () => getTargetOrganizations(targetId, page, pageSize),
    enabled: !!targetId,
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
)
{
  return useQuery({
    queryKey: targetKeys.endpoints(targetId, {
      page: params?.page,
      pageSize: params?.pageSize,
      filter: params?.filter,
    }),
    queryFn: () => getTargetEndpoints(targetId, params?.page || 1, params?.pageSize || 10, params?.filter),
    enabled: options?.enabled !== undefined ? options.enabled : !!targetId,
    select: (response: EndpointListResponse) => {
      return {
        endpoints: response.results || response.endpoints || [],
        pagination: normalizePagination(response, params?.page ?? 1, params?.pageSize ?? 10),
      }
    },
  })
}

/**
 * 获取目标的黑名单规则
 */
export function useTargetBlacklist(targetId: number) {
  return useQuery({
    queryKey: targetKeys.blacklist(targetId),
    queryFn: () => getTargetBlacklist(targetId),
    enabled: !!targetId,
  })
}
