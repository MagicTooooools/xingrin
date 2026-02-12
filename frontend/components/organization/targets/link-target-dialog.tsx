"use client"

import React from "react"
import { Plus } from "@/components/icons"
import { useTranslations } from "next-intl"

import { Button } from "@/components/ui/button"
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog"
import { Form } from "@/components/ui/form"

import type { BatchCreateResponse } from "@/types/api-response.types"
import { useLinkTargetDialogState } from "@/components/organization/targets/link-target-dialog-state"
import {
  LinkTargetDialogFooter,
  LinkTargetDialogHeader,
  LinkTargetInputSection,
  LinkTargetOrganizationSection,
} from "@/components/organization/targets/link-target-dialog-sections"

// 组件属性类型定义
interface LinkTargetDialogProps {
  organizationId: number                                     // 组织ID（固定，不可修改）
  organizationName: string                                   // 组织名称
  onAdd?: (result: BatchCreateResponse) => void              // 添加成功回调，返回批量创建的统计信息
  open?: boolean                                             // 外部控制对话框开关状态
  onOpenChange?: (open: boolean) => void                     // 外部控制对话框开关回调
}

/**
 * 关联目标对话框组件（使用 React Query）
 * 
 * 功能特性：
 * 1. 批量输入目标并关联到组织
 * 2. 自动创建不存在的目标
 * 3. 自动管理提交状态
 * 4. 自动错误处理和成功提示
 * 5. 固定组织ID，不可修改
 */
export function LinkTargetDialog({ 
  organizationId,
  organizationName,
  onAdd,
  open: externalOpen, 
  onOpenChange: externalOnOpenChange,
}: LinkTargetDialogProps) {
  const t = useTranslations("organization.linkTarget")
  const tCommon = useTranslations("common")
  const tTarget = useTranslations("target")
  const {
    form,
    open,
    handleOpenChange,
    lineNumbersRef,
    textareaRef,
    targetValidation,
    isFormValid,
    handleTextareaScroll,
    batchCreateTargets,
    onSubmit,
  } = useLinkTargetDialogState({
    organizationId,
    onAdd,
    open: externalOpen,
    onOpenChange: externalOnOpenChange,
    t,
  })

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      {/* 触发按钮 - 仅在非外部控制时显示 */}
      {externalOpen === undefined && (
        <DialogTrigger asChild>
          <Button size="sm" variant="secondary">
            <Plus />
            {tTarget("addTarget")}
          </Button>
        </DialogTrigger>
      )}
      
      {/* 对话框内容 */}
      <DialogContent className="sm:max-w-[650px] max-h-[90vh] overflow-y-auto">
        <LinkTargetDialogHeader organizationName={organizationName} t={t} />
        
        {/* 表单 */}
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="grid gap-4 py-4">
              <LinkTargetInputSection
                t={t}
                formControl={form.control}
                name="targets"
                lineNumbersRef={lineNumbersRef}
                textareaRef={textareaRef}
                onScroll={handleTextareaScroll}
                isPending={batchCreateTargets.isPending}
                targetValidation={targetValidation}
              />

              <LinkTargetOrganizationSection organizationName={organizationName} t={t} />
            </div>
          
          {/* 对话框底部按钮 */}
          <LinkTargetDialogFooter
            tCommon={tCommon}
            t={t}
            onCancel={() => handleOpenChange(false)}
            isPending={batchCreateTargets.isPending}
            isFormValid={isFormValid}
          />
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
