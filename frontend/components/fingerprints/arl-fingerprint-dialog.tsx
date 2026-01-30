"use client"

import React, { useEffect } from "react"
import { useForm } from "react-hook-form"
import { toast } from "sonner"
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
import { Textarea } from "@/components/ui/textarea"
import {
  useCreateARLFingerprint,
  useUpdateARLFingerprint,
} from "@/hooks/use-fingerprints"
import type { ARLFingerprint } from "@/types/fingerprint.types"
import { useTranslations } from "next-intl"

interface ARLFingerprintDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  fingerprint?: ARLFingerprint | null
  onSuccess?: () => void
}

interface FormData {
  name: string
  rule: string
}

export function ARLFingerprintDialog({
  open,
  onOpenChange,
  fingerprint,
  onSuccess,
}: ARLFingerprintDialogProps) {
  const isEdit = !!fingerprint
  const t = useTranslations("tools.fingerprints")
  const tCommon = useTranslations("common.actions")

  const createMutation = useCreateARLFingerprint()
  const updateMutation = useUpdateARLFingerprint()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    defaultValues: {
      name: "",
      rule: "",
    },
  })

  useEffect(() => {
    if (fingerprint) {
      reset({
        name: fingerprint.name,
        rule: fingerprint.rule,
      })
    } else {
      reset({
        name: "",
        rule: "",
      })
    }
  }, [fingerprint, reset])

  const getErrorMessage = (error: unknown): string =>
    error instanceof Error ? error.message : ""

  const onSubmit = async (data: FormData) => {
    const payload = {
      name: data.name.trim(),
      rule: data.rule.trim(),
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
      toast.error(getErrorMessage(error) || (isEdit ? t("toast.updateFailed") : t("toast.createFailed")))
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{isEdit ? t("arl.editTitle") : t("arl.addTitle")}</DialogTitle>
          <DialogDescription>
            {isEdit ? t("arl.editDesc") : t("arl.addDesc")}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {/* Name */}
          <div className="space-y-2">
            <Label htmlFor="name">{t("form.name")} *</Label>
            <Input
              id="name"
              placeholder={t("form.arlNamePlaceholder")}
              {...register("name", { required: t("form.nameRequired") })}
            />
            {errors.name && (
              <p className="text-sm text-destructive">{errors.name.message}</p>
            )}
          </div>

          {/* Rule */}
          <div className="space-y-2">
            <Label htmlFor="rule">{t("form.arlRule")} *</Label>
            <Textarea
              id="rule"
              placeholder={t("form.arlRulePlaceholder")}
              className="font-mono text-sm min-h-[120px]"
              {...register("rule", { required: t("form.arlRuleRequired") })}
            />
            <p className="text-xs text-muted-foreground">
              {t("form.arlRuleHint")}
            </p>
            {errors.rule && (
              <p className="text-sm text-destructive">{errors.rule.message}</p>
            )}
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
