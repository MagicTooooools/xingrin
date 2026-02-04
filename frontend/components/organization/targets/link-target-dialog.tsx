"use client"

import React, { useState, useRef, useMemo } from "react"
import { Plus, Target, Building2 } from "@/components/icons"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useTranslations } from "next-intl"

// 导入 UI 组件
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { LoadingSpinner } from "@/components/loading-spinner"
import { TargetValidator } from "@/lib/target-validator"
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form"

// 导入 React Query Hooks
import { useBatchCreateTargets } from "@/hooks/use-targets"

// 导入类型定义
import type { BatchCreateResponse } from "@/types/api-response.types"

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

  // 表单验证 Schema - 使用翻译
  const formSchema = useMemo(() => z.object({
    targets: z.string()
      .min(1, { message: t("validation.required") })
      .refine(
        (val) => {
          const lines = val.split('\n').map(l => l.trim()).filter(l => l.length > 0)
          return lines.length > 0
        },
        { message: t("validation.required") }
      ),
  }), [t])

  // 对话框开关状态 - 支持外部控制
  const [internalOpen, setInternalOpen] = useState(false)
  const open = externalOpen !== undefined ? externalOpen : internalOpen
  const setOpen = externalOnOpenChange || setInternalOpen
  
  // 行号列和输入框的 ref（用于同步滚动）
  const lineNumbersRef = useRef<HTMLDivElement | null>(null)
  const textareaRef = useRef<HTMLTextAreaElement | null>(null)
  
  // 使用 React Query 的批量创建目标 mutation
  const batchCreateTargets = useBatchCreateTargets()
  
  type FormValues = z.infer<typeof formSchema>
  
  // 初始化表单
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      targets: "",
    },
  })
  
  // 监听表单值变化
  const targetsText = form.watch("targets")

  // 实时验证目标
  const targetValidation = useMemo(() => {
    const lines = targetsText
      .split("\n")
      .map((s) => s.trim())
      .filter((s) => s.length > 0)

    if (lines.length === 0) {
      return {
        count: 0,
        invalid: []
      }
    }

    const results = TargetValidator.validateTargetBatch(lines)
    const invalid = results
      .filter((r) => !r.isValid)
      .map((r) => ({ index: r.index, originalTarget: r.originalTarget, error: r.error || t("validation.invalidFormat"), type: r.type }))
    
    return {
      count: lines.length,
      invalid
    }
  }, [targetsText, t])


  // 处理表单提交
  const onSubmit = (values: FormValues) => {
    // 检查是否有无效目标
    if (targetValidation.invalid.length > 0) {
      return
    }

    // 解析目标列表（每行一个目标）
    const targetList = values.targets
      .split("\n")
      .map(line => line.trim())
      .filter(line => line.length > 0)
      .map(name => ({
        name,
      }))

    if (targetList.length === 0) {
      return
    }

    // 使用 React Query mutation
    batchCreateTargets.mutate(
      {
        targets: targetList,
        organizationId: organizationId,
      },
      {
        onSuccess: (batchCreateResult) => {
          // 重置表单
          form.reset()
          
          // 关闭对话框
          setOpen(false)
          
          // 调用外部回调（如果提供）
          if (onAdd) {
            // 将批量创建结果适配为通用的 BatchCreateResponse 结构
            const adaptedResult: BatchCreateResponse = {
              message: batchCreateResult.message,
              requestedCount:
                batchCreateResult.createdCount +
                batchCreateResult.reusedCount +
                batchCreateResult.failedCount,
              createdCount: batchCreateResult.createdCount,
              existedCount: batchCreateResult.reusedCount,
              skippedCount: 0,
              skippedDomains: batchCreateResult.failedTargets.map((item) => ({
                name: item.name,
                reason: item.reason,
              })),
            }

            onAdd(adaptedResult)
          }
        }
      }
    )
  }

  // 处理对话框关闭
  const handleOpenChange = (newOpen: boolean) => {
    if (!batchCreateTargets.isPending) {
      setOpen(newOpen)
      if (!newOpen) {
        // 关闭时重置表单
        form.reset()
      }
    }
  }

  // 表单验证
  const isFormValid = form.formState.isValid && targetValidation.invalid.length === 0
  
  // 同步输入框和行号列的滚动
  const handleTextareaScroll = (e: React.UIEvent<HTMLTextAreaElement>) => {
    if (lineNumbersRef.current) {
      lineNumbersRef.current.scrollTop = e.currentTarget.scrollTop
    }
  }

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
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <Target />
            <span>{t("title")}</span>
          </DialogTitle>
          <DialogDescription>
            {t("description", { name: organizationName })}
          </DialogDescription>
        </DialogHeader>
        
        {/* 表单 */}
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="grid gap-4 py-4">
              {/* 目标输入框 - 支持多行，带行号 */}
              <FormField
                control={form.control}
                name="targets"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      {t("targetLabel")} <span className="text-destructive">*</span>
                    </FormLabel>
                    <FormControl>
                      <div className="relative border rounded-md overflow-hidden bg-background">
                        <div className="flex h-[324px]">
                          {/* 行号列 - 固定显示15行 */}
                          <div className="flex-shrink-0 w-12 bg-muted/30 border-r select-none overflow-hidden">
                            <div 
                              ref={lineNumbersRef}
                              className="py-3 px-2 text-right font-mono text-xs text-muted-foreground leading-[1.4] h-full overflow-y-auto scrollbar-hide"
                            >
                              {Array.from({ length: Math.max(field.value.split('\n').length, 15) }, (_, i) => (
                                <div key={i + 1} className="h-[20px]">
                                  {i + 1}
                                </div>
                              ))}
                            </div>
                          </div>
                          {/* 输入框 - 固定高度显示15行 */}
                          <Textarea
                            {...field}
                            ref={(e) => {
                              field.ref(e)
                              textareaRef.current = e
                            }}
                            onScroll={handleTextareaScroll}
                            placeholder={t("placeholder")}
                            disabled={batchCreateTargets.isPending}
                            className="font-mono h-full overflow-y-auto resize-none border-0 focus-visible:ring-0 focus-visible:ring-offset-0 leading-[1.4] text-sm py-3"
                            style={{ lineHeight: '20px' }}
                          />
                        </div>
                      </div>
                    </FormControl>
                    <FormDescription>
                      {t("targetCount", { count: targetValidation.count })}
                      {targetValidation.invalid.length > 0 && (
                        <span className="text-destructive ml-2">
                          | {t("invalidCount", { count: targetValidation.invalid.length })}
                        </span>
                      )}
                    </FormDescription>
                    {targetValidation.invalid.length > 0 && (
                      <div className="text-xs text-destructive">
                        {t("invalidExample", { line: targetValidation.invalid[0].index + 1, target: targetValidation.invalid[0].originalTarget, error: targetValidation.invalid[0].error })}
                      </div>
                    )}
                    <FormMessage />
                  </FormItem>
                )}
              />

            {/* 所属组织（只读显示） */}
            <div className="grid gap-2">
              <Label className="flex items-center space-x-2">
                <Building2 />
                <span>{t("organizationLabel")}</span>
              </Label>
              <div className="flex items-center gap-2 px-3 py-2 border rounded-md bg-muted/50">
                <Building2 className="h-4 w-4 text-muted-foreground" />
                <span className="font-medium">{organizationName}</span>
              </div>
            </div>
            </div>
          
          {/* 对话框底部按钮 */}
          <DialogFooter>
            <Button 
              type="button" 
              variant="outline" 
              onClick={() => handleOpenChange(false)}
              disabled={batchCreateTargets.isPending}
            >
              {tCommon("actions.cancel")}
            </Button>
            <Button 
              type="submit" 
              disabled={batchCreateTargets.isPending || !isFormValid}
            >
              {batchCreateTargets.isPending ? (
                <>
                  <LoadingSpinner/>
                  {t("creating")}
                </>
              ) : (
                <>
                  <Plus />
                  {t("createTarget")}
                </>
              )}
            </Button>
          </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
