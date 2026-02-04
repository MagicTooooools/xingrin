"use client"

import React, { useState, useRef, useMemo } from "react"
import { Plus, Building2, Target } from "@/components/icons"
import { useForm } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import * as z from "zod"
import { useTranslations } from "next-intl"

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
import { Input } from "@/components/ui/input"
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

import { useCreateOrganization } from "@/hooks/use-organizations"
import { batchCreateTargets } from "@/services/target.service"

import type { Organization } from "@/types/organization.types"

interface AddOrganizationDialogProps {
  onAdd?: (organization: Organization) => void
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

export function AddOrganizationDialog({ 
  onAdd, 
  open: externalOpen, 
  onOpenChange: externalOnOpenChange 
}: AddOrganizationDialogProps) {
  const t = useTranslations("organization.dialog")
  const tValidation = useTranslations("organization.validation")

  const formSchema = z.object({
    name: z.string()
      .min(2, { message: tValidation("nameMin", { min: 2 }) })
      .max(50, { message: tValidation("nameMax", { max: 50 }) }),
    description: z.string().max(200, { message: tValidation("descMax", { max: 200 }) }).optional(),
    targets: z.string().optional(),
  })

  type FormValues = z.infer<typeof formSchema>

  const [internalOpen, setInternalOpen] = useState(false)
  const open = externalOpen !== undefined ? externalOpen : internalOpen
  const setOpen = externalOnOpenChange || setInternalOpen

  const lineNumbersRef = useRef<HTMLDivElement | null>(null)
  const textareaRef = useRef<HTMLTextAreaElement | null>(null)

  const createOrganization = useCreateOrganization()
  const [isCreatingTargets, setIsCreatingTargets] = useState(false)
  
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      description: "",
      targets: "",
    },
  })

  const targetsText = form.watch("targets") || ""

  const targetValidation = useMemo(() => {
    const lines = targetsText
      .split("\n")
      .map((s) => s.trim())
      .filter((s) => s.length > 0)

    if (lines.length === 0) {
      return { count: 0, invalid: [] }
    }

    const results = TargetValidator.validateTargetBatch(lines)
    const invalid = results
      .filter((r) => !r.isValid)
      .map((r) => ({ index: r.index, originalTarget: r.originalTarget, error: r.error || tValidation("targetInvalid"), type: r.type }))
    
    return { count: lines.length, invalid }
  }, [targetsText, tValidation])

  const handleTextareaScroll = (e: React.UIEvent<HTMLTextAreaElement>) => {
    if (lineNumbersRef.current) {
      lineNumbersRef.current.scrollTop = e.currentTarget.scrollTop
    }
  }

  const onSubmit = async (values: FormValues) => {
    if (targetValidation.invalid.length > 0) return

    createOrganization.mutate(
      {
        name: values.name.trim(),
        description: values.description?.trim() || "",
      },
      {
        onSuccess: async (newOrganization) => {
          if (values.targets && values.targets.trim()) {
            const targetList = values.targets
              .split("\n")
              .map(line => line.trim())
              .filter(line => line.length > 0)
              .map(name => ({ name }))

            if (targetList.length > 0) {
              setIsCreatingTargets(true)
              try {
                // Call service directly to avoid double toast
                await batchCreateTargets({ targets: targetList, organizationId: newOrganization.id })
              } finally {
                setIsCreatingTargets(false)
              }
            }
          }
          form.reset()
          setOpen(false)
          if (onAdd) onAdd(newOrganization)
        }
      }
    )
  }

  const handleOpenChange = (newOpen: boolean) => {
    if (!createOrganization.isPending && !isCreatingTargets) {
      setOpen(newOpen)
      if (!newOpen) form.reset()
    }
  }

  const isFormValid = form.formState.isValid && targetValidation.invalid.length === 0
  const isSubmitting = createOrganization.isPending || isCreatingTargets

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      {externalOpen === undefined && (
        <DialogTrigger asChild>
          <Button size="sm">
            <Plus />
            {t("addButton")}
          </Button>
        </DialogTrigger>
      )}
      
      <DialogContent className="sm:max-w-[650px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <Building2 />
            <span>{t("addTitle")}</span>
          </DialogTitle>
          <DialogDescription>{t("addDesc")}</DialogDescription>
        </DialogHeader>
        
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)}>
            <div className="grid gap-4 py-4">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      {t("orgName")} <span className="text-destructive">*</span>
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("orgNamePlaceholder")}
                        disabled={isSubmitting}
                        maxLength={50}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      {t("characters", { count: field.value.length, max: 50 })}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              
              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t("orgDesc")}</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder={t("orgDescPlaceholder")}
                        disabled={isSubmitting}
                        rows={3}
                        maxLength={200}
                        {...field}
                      />
                    </FormControl>
                    <FormDescription>
                      {t("characters", { count: (field.value || "").length, max: 200 })}
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="targets"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="flex items-center space-x-2">
                      <Target className="h-4 w-4" />
                      <span>{t("addTargets")}</span>
                    </FormLabel>
                    <FormControl>
                      <div className="relative border rounded-md overflow-hidden bg-background">
                        <div className="flex h-[324px]">
                          <div className="flex-shrink-0 w-12 bg-muted/30 border-r select-none overflow-hidden">
                            <div 
                              ref={lineNumbersRef}
                              className="py-3 px-2 text-right font-mono text-xs text-muted-foreground leading-[1.4] h-full overflow-y-auto scrollbar-hide"
                            >
                              {Array.from({ length: Math.max(field.value?.split('\n').length || 1, 15) }, (_, i) => (
                                <div key={i + 1} className="h-[20px]">{i + 1}</div>
                              ))}
                            </div>
                          </div>
                          <Textarea
                            {...field}
                            ref={(e) => { field.ref(e); textareaRef.current = e }}
                            onScroll={handleTextareaScroll}
                            placeholder={t("targetsPlaceholder")}
                            disabled={isSubmitting}
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
                        {t("invalidExample", { 
                          line: targetValidation.invalid[0].index + 1, 
                          target: targetValidation.invalid[0].originalTarget, 
                          error: targetValidation.invalid[0].error 
                        })}
                      </div>
                    )}
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            
            <DialogFooter>
              <Button type="button" variant="outline" onClick={() => handleOpenChange(false)} disabled={isSubmitting}>
                {t("cancel")}
              </Button>
              <Button type="submit" disabled={isSubmitting || !isFormValid}>
                {isSubmitting ? (
                  <>
                    <LoadingSpinner/>
                    {createOrganization.isPending ? t("creating") : t("creatingTargets")}
                  </>
                ) : (
                  <>
                    <Plus />
                    {t("create")}
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
