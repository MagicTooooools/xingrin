"use client"

import { useQuery, useMutation, useQueryClient, keepPreviousData } from "@tanstack/react-query"
import { useTranslations } from "next-intl"
import { toast } from "sonner"

import { VulnerabilityService } from "@/services/vulnerability.service"
import type {
  Vulnerability,
  VulnerabilitySeverity,
  GetVulnerabilitiesParams,
} from "@/types/vulnerability.types"
import type { PaginationInfo } from "@/types/common.types"

export const vulnerabilityKeys = {
  all: ["vulnerabilities"] as const,
  list: (params: GetVulnerabilitiesParams, filter?: string) =>
    [...vulnerabilityKeys.all, "list", params, filter] as const,
  byScan: (scanId: number, params: GetVulnerabilitiesParams, filter?: string) =>
    [...vulnerabilityKeys.all, "scan", scanId, params, filter] as const,
  byTarget: (targetId: number, params: GetVulnerabilitiesParams, filter?: string) =>
    [...vulnerabilityKeys.all, "target", targetId, params, filter] as const,
  stats: () => [...vulnerabilityKeys.all, "stats"] as const,
  statsByTarget: (targetId: number) => [...vulnerabilityKeys.all, "stats", "target", targetId] as const,
}

/** 获取所有漏洞 */
export function useAllVulnerabilities(
  params?: GetVulnerabilitiesParams,
  options?: { enabled?: boolean },
  filter?: string,
) {
  const defaultParams: GetVulnerabilitiesParams = {
    page: 1,
    pageSize: 10,
    ...params,
  }

  return useQuery({
    queryKey: vulnerabilityKeys.list(defaultParams, filter),
    queryFn: () => VulnerabilityService.getAllVulnerabilities(defaultParams, filter),
    enabled: options?.enabled ?? true,
    select: (response: any) => {
      const items = (response?.results ?? []) as any[]

      const vulnerabilities: Vulnerability[] = items.map((item) => {
        let severity = (item.severity || "info") as
          | VulnerabilitySeverity
          | "unknown"
        if (severity === "unknown") {
          severity = "info"
        }

        let cvssScore: number | undefined
        if (typeof item.cvssScore === "number") {
          cvssScore = item.cvssScore
        } else if (item.cvssScore != null) {
          const num = Number(item.cvssScore)
          cvssScore = Number.isNaN(num) ? undefined : num
        }

        const createdAt: string = item.createdAt

        return {
          id: item.id,
          vulnType: item.vulnType || "unknown",
          url: item.url || "",
          description: item.description || "",
          severity: severity as VulnerabilitySeverity,
          source: item.source || "scan",
          cvssScore,
          rawOutput: item.rawOutput || {},
          isReviewed: item.isReviewed ?? false,
          reviewedAt: item.reviewedAt ?? null,
          createdAt,
        }
      })

      const pagination: PaginationInfo = {
        total: response?.total ?? 0,
        page: response?.page ?? defaultParams.page ?? 1,
        pageSize:
          response?.pageSize ??
          response?.page_size ??
          defaultParams.pageSize ??
          10,
        totalPages:
          response?.totalPages ??
          response?.total_pages ??
          0,
      }

      return { vulnerabilities, pagination }
    },
    placeholderData: keepPreviousData,
  })
}

export function useScanVulnerabilities(
  scanId: number,
  params?: GetVulnerabilitiesParams,
  options?: { enabled?: boolean },
  filter?: string,
) {
  const defaultParams: GetVulnerabilitiesParams = {
    page: 1,
    pageSize: 10,
    ...params,
  }

  return useQuery({
    queryKey: vulnerabilityKeys.byScan(scanId, defaultParams, filter),
    queryFn: () =>
      VulnerabilityService.getVulnerabilitiesByScanId(scanId, defaultParams, filter),
    enabled: options?.enabled !== undefined ? options.enabled : !!scanId,
    select: (response: any) => {
      const items = (response?.results ?? []) as any[]

      const vulnerabilities: Vulnerability[] = items.map((item) => {
        let severity = (item.severity || "info") as
          | VulnerabilitySeverity
          | "unknown"
        if (severity === "unknown") {
          severity = "info"
        }

        let cvssScore: number | undefined
        if (typeof item.cvssScore === "number") {
          cvssScore = item.cvssScore
        } else if (item.cvssScore != null) {
          const num = Number(item.cvssScore)
          cvssScore = Number.isNaN(num) ? undefined : num
        }

        const createdAt: string = item.createdAt

        return {
          id: item.id,
          vulnType: item.vulnType || "unknown",
          url: item.url || "",
          description: item.description || "",
          severity: severity as VulnerabilitySeverity,
          source: item.source || "scan",
          cvssScore,
          rawOutput: item.rawOutput || {},
          isReviewed: item.isReviewed ?? false,
          reviewedAt: item.reviewedAt ?? null,
          createdAt,
        }
      })

      const pagination: PaginationInfo = {
        total: response?.total ?? 0,
        page: response?.page ?? defaultParams.page ?? 1,
        pageSize:
          response?.pageSize ??
          response?.page_size ??
          defaultParams.pageSize ??
          10,
        totalPages:
          response?.totalPages ??
          response?.total_pages ??
          0,
      }

      return { vulnerabilities, pagination }
    },
    placeholderData: keepPreviousData,
  })
}

export function useTargetVulnerabilities(
  targetId: number,
  params?: GetVulnerabilitiesParams,
  options?: { enabled?: boolean },
  filter?: string,
) {
  const defaultParams: GetVulnerabilitiesParams = {
    page: 1,
    pageSize: 10,
    ...params,
  }

  return useQuery({
    queryKey: vulnerabilityKeys.byTarget(targetId, defaultParams, filter),
    queryFn: () =>
      VulnerabilityService.getVulnerabilitiesByTargetId(targetId, defaultParams, filter),
    enabled: options?.enabled !== undefined ? options.enabled : !!targetId,
    select: (response: any) => {
      const items = (response?.results ?? []) as any[]

      const vulnerabilities: Vulnerability[] = items.map((item) => {
        let severity = (item.severity || "info") as
          | VulnerabilitySeverity
          | "unknown"
        if (severity === "unknown") {
          severity = "info"
        }

        let cvssScore: number | undefined
        if (typeof item.cvssScore === "number") {
          cvssScore = item.cvssScore
        } else if (item.cvssScore != null) {
          const num = Number(item.cvssScore)
          cvssScore = Number.isNaN(num) ? undefined : num
        }

        const createdAt: string = item.createdAt

        return {
          id: item.id,
          vulnType: item.vulnType || "unknown",
          url: item.url || "",
          description: item.description || "",
          severity: severity as VulnerabilitySeverity,
          source: item.source || "scan",
          target: item.target ?? targetId,
          cvssScore,
          rawOutput: item.rawOutput || {},
          isReviewed: item.isReviewed ?? false,
          reviewedAt: item.reviewedAt ?? null,
          createdAt,
        }
      })

      const pagination: PaginationInfo = {
        total: response?.total ?? 0,
        page: response?.page ?? defaultParams.page ?? 1,
        pageSize:
          response?.pageSize ??
          response?.page_size ??
          defaultParams.pageSize ??
          10,
        totalPages:
          response?.totalPages ??
          response?.total_pages ??
          0,
      }

      return { vulnerabilities, pagination }
    },
    placeholderData: keepPreviousData,
  })
}

/** Mark a single vulnerability as reviewed */
export function useMarkAsReviewed() {
  const queryClient = useQueryClient()
  const t = useTranslations("vulnerabilities")

  return useMutation({
    mutationFn: (id: number) => VulnerabilityService.markAsReviewed(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: vulnerabilityKeys.all })
      toast.success(t("reviewSuccess"))
    },
    onError: () => {
      toast.error(t("reviewError"))
    },
  })
}

/** Mark a single vulnerability as pending (unreview) */
export function useMarkAsUnreviewed() {
  const queryClient = useQueryClient()
  const t = useTranslations("vulnerabilities")

  return useMutation({
    mutationFn: (id: number) => VulnerabilityService.markAsUnreviewed(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: vulnerabilityKeys.all })
      toast.success(t("unreviewSuccess"))
    },
    onError: () => {
      toast.error(t("unreviewError"))
    },
  })
}

/** Bulk mark vulnerabilities as reviewed */
export function useBulkMarkAsReviewed() {
  const queryClient = useQueryClient()
  const t = useTranslations("vulnerabilities")

  return useMutation({
    mutationFn: (ids: number[]) => VulnerabilityService.bulkMarkAsReviewed(ids),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: vulnerabilityKeys.all })
      toast.success(t("bulkReviewSuccess", { count: data.updatedCount }))
    },
    onError: () => {
      toast.error(t("bulkReviewError"))
    },
  })
}

/** Bulk mark vulnerabilities as pending (unreview) */
export function useBulkMarkAsUnreviewed() {
  const queryClient = useQueryClient()
  const t = useTranslations("vulnerabilities")

  return useMutation({
    mutationFn: (ids: number[]) => VulnerabilityService.bulkMarkAsUnreviewed(ids),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: vulnerabilityKeys.all })
      toast.success(t("bulkUnreviewSuccess", { count: data.updatedCount }))
    },
    onError: () => {
      toast.error(t("bulkUnreviewError"))
    },
  })
}

/** Get global vulnerability stats */
export function useVulnerabilityStats(options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: vulnerabilityKeys.stats(),
    queryFn: () => VulnerabilityService.getStats(),
    enabled: options?.enabled ?? true,
  })
}

/** Get vulnerability stats by target ID */
export function useTargetVulnerabilityStats(targetId: number, options?: { enabled?: boolean }) {
  return useQuery({
    queryKey: vulnerabilityKeys.statsByTarget(targetId),
    queryFn: () => VulnerabilityService.getStatsByTargetId(targetId),
    enabled: options?.enabled !== undefined ? options.enabled : !!targetId,
  })
}
