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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  useCreateEholeFingerprint,
  useUpdateEholeFingerprint,
} from "@/hooks/use-fingerprints"
import type { EholeFingerprint } from "@/types/fingerprint.types"
import { useTranslations } from "next-intl"

interface EholeFingerprintDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  fingerprint?: EholeFingerprint | null
  onSuccess?: () => void
}

interface FormData {
  cms: string
  method: string
  location: string
  keyword: string
  type: string
  isImportant: boolean
}

const METHOD_OPTIONS = [
  { value: "keyword", label: "keyword" },
  { value: "faviconhash", label: "faviconhash" },
  { value: "icon_hash", label: "icon_hash" },
  { value: "header", label: "header" },
]

const LOCATION_OPTIONS = [
  { value: "body", label: "body" },
  { value: "header", label: "header" },
  { value: "title", label: "title" },
]

export function EholeFingerprintDialog({
  open,
  onOpenChange,
  fingerprint,
  onSuccess,
}: EholeFingerprintDialogProps) {
  const isEdit = !!fingerprint
  const t = useTranslations("tools.fingerprints")
  const tCommon = useTranslations("common.actions")
  const tColumns = useTranslations("columns.fingerprint")

  const createMutation = useCreateEholeFingerprint()
  const updateMutation = useUpdateEholeFingerprint()

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    defaultValues: {
      cms: "",
      method: "keyword",
      location: "body",
      keyword: "",
      type: "-",
      isImportant: false,
    },
  })

  useEffect(() => {
    if (fingerprint) {
      reset({
        cms: fingerprint.cms,
        method: fingerprint.method,
        location: fingerprint.location,
        keyword: fingerprint.keyword.join(", "),
        type: fingerprint.type || "-",
        isImportant: fingerprint.isImportant,
      })
    } else {
      reset({
        cms: "",
        method: "keyword",
        location: "body",
        keyword: "",
        type: "-",
        isImportant: false,
      })
    }
  }, [fingerprint, reset])

  const getErrorMessage = (error: unknown): string =>
    error instanceof Error ? error.message : ""

  const onSubmit = async (data: FormData) => {
    const keywordArray = data.keyword
      .split(",")
      .map((k) => k.trim())
      .filter((k) => k.length > 0)

    if (keywordArray.length === 0) {
      toast.error(t("form.keywordRequired"))
      return
    }

    const payload = {
      cms: data.cms.trim(),
      method: data.method,
      location: data.location,
      keyword: keywordArray,
      type: data.type || "-",
      isImportant: data.isImportant,
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

  const method = watch("method")
  const location = watch("location")
  const isImportant = watch("isImportant")

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{isEdit ? t("ehole.editTitle") : t("ehole.addTitle")}</DialogTitle>
          <DialogDescription>
            {isEdit ? t("ehole.editDesc") : t("ehole.addDesc")}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {/* CMS 名称 */}
          <div className="space-y-2">
            <Label htmlFor="cms">{tColumns("cms")} *</Label>
            <Input
              id="cms"
              placeholder={t("form.cmsPlaceholder")}
              {...register("cms", { required: t("form.cmsRequired") })}
            />
            {errors.cms && (
              <p className="text-sm text-destructive">{errors.cms.message}</p>
            )}
          </div>

          {/* 匹配方式 & 匹配位置 */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label>{tColumns("method")} *</Label>
              <Select value={method} onValueChange={(v) => setValue("method", v)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {METHOD_OPTIONS.map((opt) => (
                    <SelectItem key={opt.value} value={opt.value}>
                      {opt.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label>{t("form.location")} *</Label>
              <Select value={location} onValueChange={(v) => setValue("location", v)}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {LOCATION_OPTIONS.map((opt) => (
                    <SelectItem key={opt.value} value={opt.value}>
                      {opt.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* 关键词 */}
          <div className="space-y-2">
            <Label htmlFor="keyword">{tColumns("keyword")} *</Label>
            <Input
              id="keyword"
              placeholder={t("form.keywordPlaceholder")}
              {...register("keyword", { required: t("form.keywordRequired") })}
            />
            {errors.keyword && (
              <p className="text-sm text-destructive">{errors.keyword.message}</p>
            )}
          </div>

          {/* 类型 & 重点资产 */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="type">{tColumns("type") || "Type"}</Label>
              <Input
                id="type"
                placeholder={t("form.typePlaceholder")}
                {...register("type")}
              />
            </div>

            <div className="space-y-2">
              <Label>{t("form.mark")}</Label>
              <div className="flex items-center space-x-2 h-9">
                <Checkbox
                  id="isImportant"
                  checked={isImportant}
                  onCheckedChange={(checked) => setValue("isImportant", !!checked)}
                />
                <Label htmlFor="isImportant" className="cursor-pointer font-normal">
                  {tColumns("important")}
                </Label>
              </div>
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
