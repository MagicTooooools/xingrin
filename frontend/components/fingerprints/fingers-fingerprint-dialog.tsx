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
import { Checkbox } from "@/components/ui/checkbox"
import { Textarea } from "@/components/ui/textarea"
import {
  useCreateFingersFingerprint,
  useUpdateFingersFingerprint,
} from "@/hooks/use-fingerprints"
import type { FingersFingerprint, FingersRule } from "@/types/fingerprint.types"
import { useTranslations } from "next-intl"

interface FingersFingerprintDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  fingerprint?: FingersFingerprint | null
  onSuccess?: () => void
}

interface FormData {
  name: string
  link: string
  rule: string
  tag: string
  focus: boolean
  defaultPort: string
}

export function FingersFingerprintDialog({
  open,
  onOpenChange,
  fingerprint,
  onSuccess,
}: FingersFingerprintDialogProps) {
  const isEdit = !!fingerprint
  const t = useTranslations("tools.fingerprints")
  const tCommon = useTranslations("common.actions")

  const createMutation = useCreateFingersFingerprint()
  const updateMutation = useUpdateFingersFingerprint()

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    defaultValues: {
      name: "",
      link: "",
      rule: "[]",
      tag: "",
      focus: false,
      defaultPort: "",
    },
  })

  useEffect(() => {
    if (fingerprint) {
      reset({
        name: fingerprint.name,
        link: fingerprint.link || "",
        rule: JSON.stringify(fingerprint.rule || [], null, 2),
        tag: (fingerprint.tag || []).join(", "),
        focus: fingerprint.focus || false,
        defaultPort: (fingerprint.defaultPort || []).join(", "),
      })
    } else {
      reset({
        name: "",
        link: "",
        rule: "[]",
        tag: "",
        focus: false,
        defaultPort: "",
      })
    }
  }, [fingerprint, reset])

  const getErrorMessage = (error: unknown): string =>
    error instanceof Error ? error.message : ""

  const onSubmit = async (data: FormData) => {
    // Parse rule JSON
    let ruleArray: FingersRule[]
    try {
      ruleArray = JSON.parse(data.rule)
      if (!Array.isArray(ruleArray)) {
        toast.error(t("form.ruleArrayRequired"))
        return
      }
    } catch {
      toast.error(t("form.ruleJsonInvalid"))
      return
    }

    // Parse tags
    const tagArray = data.tag
      .split(",")
      .map((t) => t.trim())
      .filter((t) => t.length > 0)

    // Parse ports
    const portArray = data.defaultPort
      .split(",")
      .map((p) => parseInt(p.trim(), 10))
      .filter((p) => !isNaN(p))

    const payload = {
      name: data.name.trim(),
      link: data.link.trim(),
      rule: ruleArray,
      tag: tagArray,
      focus: data.focus,
      defaultPort: portArray,
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

  const focus = watch("focus")

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{isEdit ? t("fingers.editTitle") : t("fingers.addTitle")}</DialogTitle>
          <DialogDescription>
            {isEdit ? t("fingers.editDesc") : t("fingers.addDesc")}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {/* Name */}
          <div className="space-y-2">
            <Label htmlFor="name">{t("form.name")} *</Label>
            <Input
              id="name"
              placeholder={t("form.namePlaceholder")}
              {...register("name", { required: t("form.nameRequired") })}
            />
            {errors.name && (
              <p className="text-sm text-destructive">{errors.name.message}</p>
            )}
          </div>

          {/* Link */}
          <div className="space-y-2">
            <Label htmlFor="link">{t("form.link")}</Label>
            <Input
              id="link"
              placeholder={t("form.linkPlaceholder")}
              {...register("link")}
            />
          </div>

          {/* Rule (JSON) */}
          <div className="space-y-2">
            <Label htmlFor="rule">{t("form.rule")} *</Label>
            <Textarea
              id="rule"
              placeholder={t("form.rulePlaceholder")}
              className="font-mono text-sm min-h-[120px]"
              {...register("rule", { required: t("form.ruleRequired") })}
            />
            {errors.rule && (
              <p className="text-sm text-destructive">{errors.rule.message}</p>
            )}
          </div>

          {/* Tags & Focus */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="tag">{t("form.tag")}</Label>
              <Input
                id="tag"
                placeholder={t("form.tagPlaceholder")}
                {...register("tag")}
              />
            </div>

            <div className="space-y-2">
              <Label>{t("form.mark")}</Label>
              <div className="flex items-center space-x-2 h-9">
                <Checkbox
                  id="focus"
                  checked={focus}
                  onCheckedChange={(checked) => setValue("focus", !!checked)}
                />
                <Label htmlFor="focus" className="cursor-pointer font-normal">
                  {t("form.focusLabel")}
                </Label>
              </div>
            </div>
          </div>

          {/* Default Ports */}
          <div className="space-y-2">
            <Label htmlFor="defaultPort">{t("form.defaultPort")}</Label>
            <Input
              id="defaultPort"
              placeholder={t("form.defaultPortPlaceholder")}
              {...register("defaultPort")}
            />
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
