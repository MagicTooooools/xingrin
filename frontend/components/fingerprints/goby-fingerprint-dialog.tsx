"use client"

import React, { useEffect } from "react"
import { useForm, useFieldArray } from "react-hook-form"
import { toast } from "sonner"
import { IconPlus, IconTrash } from "@tabler/icons-react"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Checkbox } from "@/components/ui/checkbox"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  useCreateGobyFingerprint,
  useUpdateGobyFingerprint,
} from "@/hooks/use-fingerprints"
import type { GobyFingerprint, GobyRule } from "@/types/fingerprint.types"
import { useTranslations } from "next-intl"

interface GobyFingerprintDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  fingerprint?: GobyFingerprint | null
  onSuccess?: () => void
}

interface FormData {
  name: string
  logic: string
  rule: GobyRule[]
}

const LABEL_OPTIONS = [
  { value: "title", label: "title" },
  { value: "header", label: "header" },
  { value: "body", label: "body" },
  { value: "server", label: "server" },
  { value: "banner", label: "banner" },
  { value: "port", label: "port" },
  { value: "protocol", label: "protocol" },
  { value: "cert", label: "cert" },
]

export function GobyFingerprintDialog({
  open,
  onOpenChange,
  fingerprint,
  onSuccess,
}: GobyFingerprintDialogProps) {
  const isEdit = !!fingerprint
  const t = useTranslations("tools.fingerprints")
  const tCommon = useTranslations("common.actions")
  const tColumns = useTranslations("columns.fingerprint")

  const createMutation = useCreateGobyFingerprint()
  const updateMutation = useUpdateGobyFingerprint()

  const {
    register,
    handleSubmit,
    reset,
    control,
    setValue,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    defaultValues: {
      name: "",
      logic: "a",
      rule: [{ label: "title", feature: "", is_equal: true }],
    },
  })

  const { fields, append, remove } = useFieldArray({
    control,
    name: "rule",
  })

  useEffect(() => {
    if (fingerprint) {
      reset({
        name: fingerprint.name,
        logic: fingerprint.logic,
        rule: fingerprint.rule.length > 0 
          ? fingerprint.rule 
          : [{ label: "title", feature: "", is_equal: true }],
      })
    } else {
      reset({
        name: "",
        logic: "a",
        rule: [{ label: "title", feature: "", is_equal: true }],
      })
    }
  }, [fingerprint, reset])

  const onSubmit = async (data: FormData) => {
    if (data.rule.length === 0) {
      toast.error(t("form.logicRequired"))
      return
    }

    const payload = {
      name: data.name.trim(),
      logic: data.logic.trim(),
      rule: data.rule,
    }

    try {
      if (isEdit && fingerprint) {
        await updateMutation.mutateAsync({ id: fingerprint.id, data: payload })
        toast.success(t("toast.updateSuccess"))
      } else {
        await createMutation.mutateAsync(payload)
        toast.success(t("toast.createSuccess"))
      }
      onOpenChange(false)
      onSuccess?.()
    } catch (error) {
      const message = error instanceof Error ? error.message : ""
      toast.error(message || (isEdit ? t("toast.updateFailed") : t("toast.createFailed")))
    }
  }

  const addRule = () => {
    append({ label: "title", feature: "", is_equal: true })
  }

  const watchedRules = watch("rule")

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{isEdit ? t("goby.editTitle") : t("goby.addTitle")}</DialogTitle>
          <DialogDescription>
            {isEdit ? t("goby.editDesc") : t("goby.addDesc")}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {/* Product name & Logic expression */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="name">{tColumns("name") || "Name"} *</Label>
              <Input
                id="name"
                placeholder={t("form.namePlaceholder")}
                {...register("name", { required: t("form.nameRequired") })}
              />
              {errors.name && (
                <p className="text-sm text-destructive">{errors.name.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="logic">{tColumns("logic")} *</Label>
              <Input
                id="logic"
                placeholder={t("form.logicPlaceholder")}
                {...register("logic", { required: t("form.logicRequired") })}
              />
              {errors.logic && (
                <p className="text-sm text-destructive">{errors.logic.message}</p>
              )}
            </div>
          </div>

          {/* Rule list */}
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <Label>{tColumns("rules")} *</Label>
              <Button type="button" variant="outline" size="sm" onClick={addRule}>
                <IconPlus className="h-4 w-4" />
                {tCommon("add")}
              </Button>
            </div>
            
            <div className="space-y-2">
              {fields.map((field, index) => (
                <div key={field.id} className="flex items-center gap-2 p-3 border rounded-md bg-muted/30">
                  <div className="w-24">
                    <Select 
                      value={watchedRules[index]?.label || "title"} 
                      onValueChange={(v) => setValue(`rule.${index}.label`, v)}
                    >
                      <SelectTrigger className="h-8">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {LABEL_OPTIONS.map((opt) => (
                          <SelectItem key={opt.value} value={opt.value}>
                            {opt.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="flex-1">
                    <Input
                      {...register(`rule.${index}.feature` as const, { required: true })}
                      placeholder={t("form.featurePlaceholder")}
                      className="h-8"
                    />
                  </div>
                  <div className="flex items-center gap-1">
                    <Checkbox
                      checked={watchedRules[index]?.is_equal ?? true}
                      onCheckedChange={(checked) => setValue(`rule.${index}.is_equal`, !!checked)}
                    />
                    <span className="text-xs text-muted-foreground">Match</span>
                  </div>
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => remove(index)}
                    disabled={fields.length <= 1}
                    className="h-8 w-8 p-0"
                  >
                    <IconTrash className="h-4 w-4 text-destructive" />
                  </Button>
                </div>
              ))}
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              {tCommon("cancel")}
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? "..." : isEdit ? tCommon("save") : tCommon("create")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
