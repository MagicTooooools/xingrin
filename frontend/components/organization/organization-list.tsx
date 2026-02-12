"use client"

import {
  OrganizationListDialogs,
  OrganizationListErrorState,
  OrganizationListSkeleton,
  OrganizationListTable,
} from "./organization-list-sections"
import { useOrganizationListState } from "./organization-list-state"

/**
 * 组织列表组件（使用 React Query）
 * 
 * 功能特性：
 * 1. 统一的 Loading 状态管理
 * 2. 自动缓存和重新验证
 * 3. 乐观更新
 * 4. 自动错误处理
 * 5. 更好的用户体验
 */
export function OrganizationList() {
  const state = useOrganizationListState()

  if (state.error) {
    return (
      <OrganizationListErrorState
        error={state.error}
        onRetry={state.refetch}
        tCommon={state.tCommon}
      />
    )
  }

  if (state.isLoading || !state.data) {
    return <OrganizationListSkeleton />
  }

  return (
    <div className="space-y-4">
      <OrganizationListTable state={state} />
      <OrganizationListDialogs state={state} />
    </div>
  )
}
