"use client"

import React, { useState, useCallback } from "react"
import { useRouter } from "next/navigation"
import { useTranslations } from "next-intl"
import { createAllTargetsColumns, AllTargetsTranslations } from "@/components/target/all-targets-columns"
import { TargetsDataTable } from "@/components/target/targets-data-table"
import { AddTargetDialog } from "@/components/target/add-target-dialog"
import { InitiateScanDialog } from "@/components/scan/initiate-scan-dialog"
import { CreateScheduledScanDialog } from "@/components/scan/scheduled/create-scheduled-scan-dialog"
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
import { formatDate } from "@/lib/utils"
import { LoadingSpinner } from "@/components/loading-spinner"
import { DataTableSkeleton } from "@/components/ui/data-table-skeleton"
import { useTargets, useDeleteTarget, useBatchDeleteTargets } from "@/hooks/use-targets"
import type { Target } from "@/types/target.types"
import type { Organization } from "@/types/organization.types"

/**
 * All targets detail view component
 * Displays a list of all targets in the system, supports search, pagination, delete operations
 */
interface AllTargetsDetailViewProps {
  className?: string
  tableClassName?: string
  hideToolbar?: boolean
  hidePagination?: boolean
}

export function AllTargetsDetailView({
  className,
  tableClassName,
  hideToolbar,
  hidePagination,
}: AllTargetsDetailViewProps) {
  const router = useRouter()
  const tColumns = useTranslations("columns")
  const tTooltips = useTranslations("tooltips")
  const tCommon = useTranslations("common")
  const tConfirm = useTranslations("common.confirm")
  const tTarget = useTranslations("target")
  
  // Build translation object
  const translations: AllTargetsTranslations = {
    columns: {
      target: tColumns("target.target"),
      organization: tColumns("organization.organization"),
      addedOn: tColumns("target.addedOn"),
      lastScanned: tColumns("target.lastScanned"),
    },
    actions: {
      scheduleScan: tTooltips("scheduleScan"),
      delete: tCommon("actions.delete"),
      selectAll: tCommon("actions.selectAll"),
      selectRow: tCommon("actions.selectRow"),
    },
    tooltips: {
      targetDetails: tTooltips("targetDetails"),
      targetSummary: tTooltips("targetSummary"),
      initiateScan: tTooltips("initiateScan"),
      clickToCopy: tTooltips("clickToCopy"),
      copied: tTooltips("copied"),
    },
  }
  const [pagination, setPagination] = useState({ pageIndex: 0, pageSize: 10 })
  const [searchQuery, setSearchQuery] = useState("")
  const [typeFilter, setTypeFilter] = useState("")
  const [selectedTargets, setSelectedTargets] = useState<Target[]>([])
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [targetToDelete, setTargetToDelete] = useState<Target | null>(null)
  const [bulkDeleteDialogOpen, setBulkDeleteDialogOpen] = useState(false)
  const [shouldPrefetchOrgs, setShouldPrefetchOrgs] = useState(false)
  const [initiateScanDialogOpen, setInitiateScanDialogOpen] = useState(false)
  const [scheduleScanDialogOpen, setScheduleScanDialogOpen] = useState(false)
  const [targetToScan, setTargetToScan] = useState<Target | null>(null)
  const [targetToSchedule, setTargetToSchedule] = useState<Target | null>(null)

  // Handle pagination state change
  const handlePaginationChange = React.useCallback((newPagination: { pageIndex: number, pageSize: number }) => {
    setPagination(newPagination)
  }, [])

  const [isSearching, setIsSearching] = React.useState(false)

  const handleSearchChange = (value: string) => {
    setIsSearching(true)
    setSearchQuery(value)
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
  }

  const handleTypeFilterChange = (value: string) => {
    setTypeFilter(value)
    setPagination((prev) => ({ ...prev, pageIndex: 0 }))
  }

  // Use API hooks
  const { data, isLoading, isFetching, error } = useTargets(pagination.pageIndex + 1, pagination.pageSize, typeFilter || undefined, searchQuery || undefined)
  const deleteTargetMutation = useDeleteTarget()
  const batchDeleteMutation = useBatchDeleteTargets()

  const targets = data?.results || []
  const totalCount = data?.total || 0

  React.useEffect(() => {
    if (!isFetching && isSearching) {
      setIsSearching(false)
    }
  }, [isFetching, isSearching])

  // Handle add target
  const handleAddTarget = useCallback(() => {
    setIsAddDialogOpen(true)
  }, [])

  // Handle delete single target
  const handleDeleteTarget = useCallback((target: Target) => {
    setTargetToDelete(target)
    setDeleteDialogOpen(true)
  }, [])

  // Confirm delete target
  const confirmDelete = async () => {
    if (!targetToDelete) return

    try {
      await deleteTargetMutation.mutateAsync({ id: targetToDelete.id, name: targetToDelete.name })
      setDeleteDialogOpen(false)
      setTargetToDelete(null)
    } catch (error) {
      // Error already handled in hook
      console.error('Delete failed:', error)
    }
  }

  // Handle batch delete
  const handleBatchDelete = useCallback(() => {
    if (selectedTargets.length === 0) return
    setBulkDeleteDialogOpen(true)
  }, [selectedTargets])

  // Confirm batch delete
  const confirmBulkDelete = async () => {
    if (selectedTargets.length === 0) return

    try {
      await batchDeleteMutation.mutateAsync({
        ids: selectedTargets.map((t) => t.id),
      })
      setBulkDeleteDialogOpen(false)
      setSelectedTargets([])
    } catch (error) {
      // Error already handled in hook
      console.error('Batch delete failed:', error)
    }
  }

  // Handle initiate scan
  const handleInitiateScan = useCallback((target: Target) => {
    setTargetToScan(target)
    setInitiateScanDialogOpen(true)
  }, [])

  // Handle scheduled scan
  const handleScheduleScan = useCallback((target: Target) => {
    setTargetToSchedule(target)
    setScheduleScanDialogOpen(true)
  }, [])

  // Create table columns
  const columns = createAllTargetsColumns({
    formatDate,
    navigate: (path: string) => router.push(path),
    handleDelete: handleDeleteTarget,
    handleInitiateScan,
    handleScheduleScan,
    t: translations,
  })

  // Loading
  if (isLoading) {
    return (
      <DataTableSkeleton
        toolbarButtonCount={2}
        rows={6}
        columns={5}
      />
    )
  }

  // Error handling
  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <p className="text-destructive mb-2">{tCommon("status.error")}</p>
          <p className="text-sm text-muted-foreground">{error.message}</p>
        </div>
      </div>
    )
  }

  return (
    <>
      <TargetsDataTable
        data={targets}
        columns={columns}
        onAddNew={handleAddTarget}
        onAddHover={() => setShouldPrefetchOrgs(true)}
        onBulkDelete={handleBatchDelete}
        onSelectionChange={setSelectedTargets}
        searchPlaceholder={tTarget("name")}
        searchValue={searchQuery}
        onSearch={handleSearchChange}
        isSearching={isSearching}
        addButtonText={tTarget("addTarget")}
        // 分页相关属性
        pagination={pagination}
        onPaginationChange={handlePaginationChange}
        totalCount={totalCount}
        manualPagination={true}
        // 类型筛选
        typeFilter={typeFilter}
        onTypeFilterChange={handleTypeFilterChange}
        // Styling
        className={className}
        tableClassName={tableClassName}
        hideToolbar={hideToolbar}
        hidePagination={hidePagination}
      />

      {/* Add target dialog */}
      <AddTargetDialog
        onAdd={() => {
          setIsAddDialogOpen(false)
        }}
        open={isAddDialogOpen}
        onOpenChange={setIsAddDialogOpen}
        prefetchEnabled={shouldPrefetchOrgs}
      />

      {/* Delete confirmation dialog */}
      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{tConfirm("deleteTargetTitle")}</AlertDialogTitle>
            <AlertDialogDescription>
              {tConfirm("deleteTargetMessage", { name: targetToDelete?.name ?? "" })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{tCommon("actions.cancel")}</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              disabled={deleteTargetMutation.isPending}
            >
              {deleteTargetMutation.isPending ? (
                <>
                  <LoadingSpinner/>
                  {tConfirm("deleting")}
                </>
              ) : (
                tConfirm("confirmDelete")
              )}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Initiate scan dialog */}
      <InitiateScanDialog
        organization={
          targetToScan?.organizations && targetToScan.organizations.length > 0
            ? {
                id: targetToScan.organizations[0].id,
                name: targetToScan.organizations[0].name,
                targetCount: 1, // Current target
              } as Organization
            : null
        }
        targetId={targetToScan?.id}
        targetName={targetToScan?.name}
        open={initiateScanDialogOpen}
        onOpenChange={setInitiateScanDialogOpen}
        onSuccess={() => {
          setTargetToScan(null)
        }}
      />

      {/* Scheduled scan dialog */}
      <CreateScheduledScanDialog
        open={scheduleScanDialogOpen}
        onOpenChange={setScheduleScanDialogOpen}
        presetTargetId={targetToSchedule?.id}
        presetTargetName={targetToSchedule?.name}
        onSuccess={() => {
          setTargetToSchedule(null)
        }}
      />

      {/* Batch delete confirmation dialog */}
      <AlertDialog open={bulkDeleteDialogOpen} onOpenChange={setBulkDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{tConfirm("bulkDeleteTargetTitle")}</AlertDialogTitle>
            <AlertDialogDescription>
              {tConfirm("bulkDeleteTargetMessage", { count: selectedTargets.length })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          {/* Target list container - fixed max height with scroll support */}
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
              disabled={batchDeleteMutation.isPending}
            >
              {batchDeleteMutation.isPending ? (
                <>
                  <LoadingSpinner/>
                  {tConfirm("deleting")}
                </>
              ) : (
                tConfirm("deleteTargetCount", { count: selectedTargets.length })
              )}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  )
}
