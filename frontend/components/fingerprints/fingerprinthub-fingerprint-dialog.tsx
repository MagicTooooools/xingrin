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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  useCreateFingerPrintHubFingerprint,
  useUpdateFingerPrintHubFingerprint,
} from "@/hooks/use-fingerprints"
import type { FingerPrintHubFingerprint, FingerPrintHubHttpMatcher } from "@/types/fingerprint.types"
import { useTranslations } from "next-intl"

interface FingerPrintHubFingerprintDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  fingerprint?: FingerPrintHubFingerprint | null
  onSuccess?: () => void
}

interface FormData {
  fpId: string
  name: string
  author: string
  tags: string
  severity: string
  metadata: string
  http: string
  sourceFile: string
}

const SEVERITY_OPTIONS = [
  { value: "info", label: "Info" },
  { value: "low", label: "Low" },
  { value: "medium", label: "Medium" },
  { value: "high", label: "High" },
  { value: "critical", label: "Critical" },
]

export function FingerPrintHubFingerprintDialog({
  open,
  onOpenChange,
  fingerprint,
  onSuccess,
}: FingerPrintHubFingerprintDialogProps) {
  const isEdit = !!fingerprint
  const t = useTranslations("tools.fingerprints")
  const tCommon = useTranslations("common.actions")

  const createMutation = useCreateFingerPrintHubFingerprint()
  const updateMutation = useUpdateFingerPrintHubFingerprint()

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    defaultValues: {
      fpId: "",
      name: "",
      author: "",
      tags: "",
      severity: "info",
      metadata: "{}",
      http: "[]",
      sourceFile: "",
    },
  })

  useEffect(() => {
    if (fingerprint) {
      reset({
        fpId: fingerprint.fpId,
        name: fingerprint.name,
        author: fingerprint.author || "",
        tags: fingerprint.tags || "",
        severity: fingerprint.severity || "info",
        metadata: JSON.stringify(fingerprint.metadata || {}, null, 2),
        http: JSON.stringify(fingerprint.http || [], null, 2),
        sourceFile: fingerprint.sourceFile || "",
      })
    } else {
      reset({
        fpId: "",
        name: "",
        author: "",
        tags: "",
        severity: "info",
        metadata: "{}",
        http: "[]",
        sourceFile: "",
      })
    }
  }, [fingerprint, reset])

  const getErrorMessage = (error: unknown): string =>
    error instanceof Error ? error.message : ""

  const onSubmit = async (data: FormData) => {
    // Parse metadata JSON
    let metadataObj: Record<string, unknown>
    try {
      metadataObj = JSON.parse(data.metadata)
      if (typeof metadataObj !== "object" || Array.isArray(metadataObj)) {
        toast.error(t("form.metadataObjectRequired"))
        return
      }
    } catch {
      toast.error(t("form.metadataJsonInvalid"))
      return
    }

    // Parse http JSON
    let httpArray: FingerPrintHubHttpMatcher[]
    try {
      httpArray = JSON.parse(data.http)
      if (!Array.isArray(httpArray)) {
        toast.error(t("form.httpArrayRequired"))
        return
      }
    } catch {
      toast.error(t("form.httpJsonInvalid"))
      return
    }

    const payload = {
      fpId: data.fpId.trim(),
      name: data.name.trim(),
      author: data.author.trim(),
      tags: data.tags.trim(),
      severity: data.severity,
      metadata: metadataObj,
      http: httpArray,
      sourceFile: data.sourceFile.trim(),
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

  const severity = watch("severity")

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[700px] max-h-[85vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{isEdit ? t("fingerprinthub.editTitle") : t("fingerprinthub.addTitle")}</DialogTitle>
          <DialogDescription>
            {isEdit ? t("fingerprinthub.editDesc") : t("fingerprinthub.addDesc")}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {/* FP ID & Name */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="fpId">{t("form.fpId")} *</Label>
              <Input
                id="fpId"
                placeholder={t("form.fpIdPlaceholder")}
                {...register("fpId", { required: t("form.fpIdRequired") })}
              />
              {errors.fpId && (
                <p className="text-sm text-destructive">{errors.fpId.message}</p>
              )}
            </div>

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
          </div>

          {/* Author & Severity */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="author">{t("form.author")}</Label>
              <Input
                id="author"
                placeholder={t("form.authorPlaceholder")}
                {...register("author")}
              />
            </div>

            <div className="space-y-2">
              <Label>{t("form.severity")}</Label>
              <Select value={severity} onValueChange={(v) => setValue("severity", v)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {SEVERITY_OPTIONS.map((opt) => (
                    <SelectItem key={opt.value} value={opt.value}>
                      {opt.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Tags */}
          <div className="space-y-2">
            <Label htmlFor="tags">{t("form.tags")}</Label>
            <Input
              id="tags"
              placeholder={t("form.tagsPlaceholder")}
              {...register("tags")}
            />
          </div>

          {/* HTTP Matchers (JSON) */}
          <div className="space-y-2">
            <Label htmlFor="http">{t("form.http")} *</Label>
            <Textarea
              id="http"
              placeholder={t("form.httpPlaceholder")}
              className="font-mono text-sm min-h-[150px]"
              {...register("http", { required: t("form.httpRequired") })}
            />
            {errors.http && (
              <p className="text-sm text-destructive">{errors.http.message}</p>
            )}
          </div>

          {/* Metadata (JSON) */}
          <div className="space-y-2">
            <Label htmlFor="metadata">{t("form.metadata")}</Label>
            <Textarea
              id="metadata"
              placeholder={t("form.metadataPlaceholder")}
              className="font-mono text-sm min-h-[80px]"
              {...register("metadata")}
            />
          </div>

          {/* Source File */}
          <div className="space-y-2">
            <Label htmlFor="sourceFile">{t("form.sourceFile")}</Label>
            <Input
              id="sourceFile"
              placeholder={t("form.sourceFilePlaceholder")}
              {...register("sourceFile")}
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
