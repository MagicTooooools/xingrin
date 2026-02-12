"use client"

import {
  TargetsDetailViewDialogs,
  TargetsDetailViewEmptyState,
  TargetsDetailViewErrorState,
  TargetsDetailViewLoadingState,
  TargetsDetailViewTable,
} from "./targets-detail-view-sections"
import { useTargetsDetailViewState } from "./targets-detail-view-state"

/**
 * 组织目标详情视图组件（使用 React Query）
 * 用于显示和管理组织下的目标列表
 * 支持通过组织ID获取数据
 */
export function OrganizationTargetsDetailView({
  organizationId,
}: {
  organizationId: string
}) {
  const state = useTargetsDetailViewState({ organizationId })

  if (state.error) {
    return (
      <TargetsDetailViewErrorState
        error={state.error}
        onRetry={state.refetch}
        tCommon={state.tCommon}
      />
    )
  }

  if (state.isLoading) {
    return <TargetsDetailViewLoadingState />
  }

  if (!state.organization) {
    return <TargetsDetailViewEmptyState tOrg={state.tOrg} />
  }

  return (
    <>
      <TargetsDetailViewTable state={state} />
      <TargetsDetailViewDialogs state={state} />
    </>
  )
}

export { OrganizationTargetsDetailView as TargetsDetailView }
