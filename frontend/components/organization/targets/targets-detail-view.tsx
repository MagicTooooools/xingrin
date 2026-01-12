"use client"

import React, { useState, useMemo } from "react"
import { AlertTriangle } from "lucide-react"
import { useRouter } from "next/navigation"
import { useTranslations, useLocale } from "next-intl"
import { toast } from "sonner"
import { TargetsDataTable } from "./targets-data-table"
import { createTargetColumns } from "./targets-columns"
import { AddTargetDialog } from "./add-target-dialog"
import { LoadingSpinner } from "@/components/loading-spinner"
import { DataTableSkeleton } from "@/components/ui/data-table-skeleton"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog"
import { useOrganization, useUnlinkTargetsFromOrganization } from "@/hooks/use-organizations"
import { useTargets } from "@/hooks/use-targets"
import { getDateLocale } from "@/lib/date-utils"
import type { Target } from "@/types/target.types"

/**
 * 组织目标详情视图组件（使用 React Query）
 * 用于显示和管理组织下的目标列表
 * 支持通过组织ID获取数据
 */
export function OrganizationTargetsDetailView({
  organizationId
}: {
  organizationId: string
}) {
  const [selectedTargets, setSelectedTargets] = useState<Target[]>([])
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [targetToDelete, setTargetToDelete] = useState<Target | null>(null)
  const [bulkDeleteDialogOpen, setBulkDeleteDialogOpen] = useState(false)

  // 国际化
  const tColumns = useTranslations("columns")
  const tCommon = useTranslations("common")
  const tTooltips = useTranslations("tooltips")
  const tTarget = useTranslations("target")
  const tConfirm = useTranslations("common.confirm")
  const tOrg = useTranslations("organization")
  const locale = useLocale()

  // 构建翻译对象
  const translations = useMemo(() => ({
    columns: {
      targetName: tColumns("target.target"),
      type: tColumns("common.type"),
    },
    actions: {
      selectAll: tCommon("actions.selectAll"),
      selectRow: tCommon("actions.selectRow"),
    },
    tooltips: {
      viewDetails: tTooltips("viewDetails"),
      unlinkTarget: tTooltips("unlinkTarget"),
      clickToCopy: tTooltips("clickToCopy"),
      copied: tTooltips("copied"),
    },
    types: {
      domain: tTarget("types.domain"),
      ip: tTarget("types.ip"),
      cidr: tTarget("types.cidr"),
    },
  }), [tColumns, tCommon, tTooltips, tTarget])

  // 使用解除关联 mutation
  const unlinkTargets = useUnlinkTargetsFromOrganization()

  // 分页状态
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  })

  // 使用 React Query 获取组织基本信息
  const {
    data: organization,
    isLoading: isLoadingOrg,
    error: orgError,
  } = useOrganization(parseInt(organizationId))

  // 使用 React Query 获取目标列表（过滤组织）
  const {
    data: targetsData,
    isLoading: isLoadingTargets,
    error: targetsError,
    refetch
  } = useTargets({
    page: pagination.pageIndex + 1,
    pageSize: pagination.pageSize,
    organizationId: parseInt(organizationId),
  })

  const isLoading = isLoadingOrg || isLoadingTargets
  const error = orgError || targetsError

  // 辅助函数 - 格式化日期
  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleString(getDateLocale(locale), {
      year: "numeric",
      month: "numeric",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
      hour12: false,
    })
  }

  // 导航函数（使用 Next.js 客户端路由）
  const router = useRouter()
  const navigate = (path: string) => {
    router.push(path)
  }

  // 处理解除关联目标
  const handleDeleteTarget = (target: Target) => {
    setTargetToDelete(target)
    setDeleteDialogOpen(true)
  }

  // 确认解除关联目标
  const confirmDelete = async () => {
    if (!targetToDelete) return

    setDeleteDialogOpen(false)
    const targetId = targetToDelete.id
    setTargetToDelete(null)

    // 调用解除关联 API
    unlinkTargets.mutate({
      organizationId: parseInt(organizationId),
      targetIds: [targetId]
    })
  }

  // 处理批量解除关联
  const handleBulkDelete = () => {
    if (selectedTargets.length === 0) {
      return
    }
    setBulkDeleteDialogOpen(true)
  }

  // 确认批量解除关联
  const confirmBulkDelete = async () => {
    if (selectedTargets.length === 0) return

    const targetIds = selectedTargets.map(target => target.id)

    setBulkDeleteDialogOpen(false)
    setSelectedTargets([])

    // 调用批量解除关联 API
    unlinkTargets.mutate({
      organizationId: parseInt(organizationId),
      targetIds
    })
  }

  // 处理添加目标
  const handleAddTarget = () => {
    setIsAddDialogOpen(true)
  }

  // 处理添加成功
  const handleAddSuccess = () => {
    setIsAddDialogOpen(false)
    // 刷新目标列表
    refetch()
  }

  // 处理分页变化
  const handlePaginationChange = (newPagination: { pageIndex: number; pageSize: number }) => {
    setPagination(newPagination)
    // 清空选中状态
    setSelectedTargets([])
  }

  // 创建列定义
  const targetColumns = useMemo(
    () =>
      createTargetColumns({
        formatDate,
        navigate,
        handleDelete: handleDeleteTarget,
        t: translations,
      }),
    [formatDate, navigate, handleDeleteTarget, translations]
  )

  // 错误状态
  if (error) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <div className="rounded-full bg-destructive/10 p-3 mb-4">
          <AlertTriangle className="h-10 w-10 text-destructive" />
        </div>
        <h3 className="text-lg font-semibold mb-2">{tCommon("status.error")}</h3>
        <p className="text-muted-foreground text-center mb-4">
          {error.message}
        </p>
        <button
          onClick={() => refetch()}
          className="px-4 py-2 bg-primary text-primary-foreground rounded-md hover:bg-primary/90"
        >
          {tCommon("actions.retry")}
        </button>
      </div>
    )
  }

  // 加载状态
  if (isLoading) {
    return (
      <DataTableSkeleton
        toolbarButtonCount={3}
        rows={6}
        columns={4}
      />
    )
  }

  if (!organization) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <p className="text-muted-foreground">{tOrg("notFound")}</p>
      </div>
    )
  }

  return (
    <>
      <TargetsDataTable
        data={targetsData?.targets || []}
        columns={targetColumns}
        onAddNew={handleAddTarget}
        onBulkDelete={handleBulkDelete}
        onSelectionChange={setSelectedTargets}
        searchPlaceholder={tColumns("target.target")}
        addButtonText={tTarget("addTarget")}
        pagination={pagination}
        setPagination={setPagination}
        paginationInfo={targetsData ? {
          total: targetsData.total,
          page: targetsData.page,
          pageSize: targetsData.pageSize,
          totalPages: targetsData.totalPages,
        } : undefined}
        onPaginationChange={handlePaginationChange}
      />

      {/* 添加目标对话框 */}
      <AddTargetDialog
        organizationId={parseInt(organizationId)}
        organizationName={organization.name}
        onAdd={handleAddSuccess}
        open={isAddDialogOpen}
        onOpenChange={setIsAddDialogOpen}
      />

      {/* 解除关联确认对话框 */}
      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{tConfirm("unlinkTitle")}</AlertDialogTitle>
            <AlertDialogDescription>
              {tConfirm("unlinkTargetMessage", { name: targetToDelete?.name ?? "" })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{tCommon("actions.cancel")}</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {tConfirm("confirmUnlink")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* 批量解除关联确认对话框 */}
      <AlertDialog open={bulkDeleteDialogOpen} onOpenChange={setBulkDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{tConfirm("bulkUnlinkTitle")}</AlertDialogTitle>
            <AlertDialogDescription>
              {tConfirm("bulkUnlinkTargetMessage", { count: selectedTargets.length })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          {/* 目标列表容器 - 固定最大高度并支持滚动 */}
          <div className="mt-2 p-2 bg-muted rounded-md max-h-96 overflow-y-auto">
            <ul className="text-sm space-y-1">
              {selectedTargets.map((target) => (
                <li key={target.id} className="flex items-center">
                  <span className="font-medium">{target.name}</span>
                  {target.description && (
                    <span className="text-muted-foreground ml-2">- {target.description}</span>
                  )}
                </li>
              ))}
            </ul>
          </div>
          <AlertDialogFooter>
            <AlertDialogCancel>{tCommon("actions.cancel")}</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmBulkDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {tConfirm("confirmUnlinkCount", { count: selectedTargets.length })}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}

export { OrganizationTargetsDetailView as TargetsDetailView }
