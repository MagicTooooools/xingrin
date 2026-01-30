"use client"

import React, { useEffect } from "react"
import { useForm } from "react-hook-form"
import { useTranslations } from "next-intl"
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
  useCreateWappalyzerFingerprint,
  useUpdateWappalyzerFingerprint,
} from "@/hooks/use-fingerprints"
import type { WappalyzerFingerprint } from "@/types/fingerprint.types"

interface WappalyzerFingerprintDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  fingerprint?: WappalyzerFingerprint | null
  onSuccess?: () => void
}

interface FormData {
  name: string
  cats: string
  description: string
  website: string
  cpe: string
  cookies: string
  headers: string
  scriptSrc: string
  js: string
  meta: string
  html: string
  implies: string
}

export function WappalyzerFingerprintDialog({
  open,
  onOpenChange,
  fingerprint,
  onSuccess,
}: WappalyzerFingerprintDialogProps) {
  const t = useTranslations("tools.fingerprints")
  const tCommon = useTranslations("common.actions")
  const isEdit = !!fingerprint

  const createMutation = useCreateWappalyzerFingerprint()
  const updateMutation = useUpdateWappalyzerFingerprint()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormData>({
    defaultValues: {
      name: "",
      cats: "",
      description: "",
      website: "",
      cpe: "",
      cookies: "{}",
      headers: "{}",
      scriptSrc: "",
      js: "",
      meta: "{}",
      html: "",
      implies: "",
    },
  })

  useEffect(() => {
    if (fingerprint) {
      reset({
        name: fingerprint.name,
        cats: fingerprint.cats?.join(", ") || "",
        description: fingerprint.description || "",
        website: fingerprint.website || "",
        cpe: fingerprint.cpe || "",
        cookies: JSON.stringify(fingerprint.cookies || {}, null, 2),
        headers: JSON.stringify(fingerprint.headers || {}, null, 2),
        scriptSrc: fingerprint.scriptSrc?.join(", ") || "",
        js: fingerprint.js?.join(", ") || "",
        meta: JSON.stringify(fingerprint.meta || {}, null, 2),
        html: fingerprint.html?.join(", ") || "",
        implies: fingerprint.implies?.join(", ") || "",
      })
    } else {
      reset({
        name: "",
        cats: "",
        description: "",
        website: "",
        cpe: "",
        cookies: "{}",
        headers: "{}",
        scriptSrc: "",
        js: "",
        meta: "{}",
        html: "",
        implies: "",
      })
    }
  }, [fingerprint, reset])

  const parseArray = (str: string): string[] => {
    return str.split(",").map(s => s.trim()).filter(s => s.length > 0)
  }

  const parseNumberArray = (str: string): number[] => {
    return str.split(",").map(s => parseInt(s.trim(), 10)).filter(n => !isNaN(n))
  }

  const parseJson = (str: string): Record<string, unknown> => {
    try {
      return JSON.parse(str)
    } catch {
      return {}
    }
  }

  const parseStringRecord = (str: string): Record<string, string> => {
    const raw = parseJson(str)
    return Object.fromEntries(
      Object.entries(raw).filter(([, value]) => typeof value === "string")
    ) as Record<string, string>
  }

  const parseStringArrayRecord = (str: string): Record<string, string[]> => {
    const raw = parseJson(str)
    const output: Record<string, string[]> = {}
    for (const [key, value] of Object.entries(raw)) {
      if (Array.isArray(value)) {
        const items = value.filter((item): item is string => typeof item === "string")
        if (items.length > 0) {
          output[key] = items
        }
      }
    }
    return output
  }

  const getErrorMessage = (error: unknown): string =>
    error instanceof Error ? error.message : ""

  const onSubmit = async (data: FormData) => {
    const payload = {
      name: data.name.trim(),
      cats: parseNumberArray(data.cats),
      description: data.description.trim(),
      website: data.website.trim(),
      cpe: data.cpe.trim(),
      cookies: parseStringRecord(data.cookies),
      headers: parseStringRecord(data.headers),
      scriptSrc: parseArray(data.scriptSrc),
      js: parseArray(data.js),
      meta: parseStringArrayRecord(data.meta),
      html: parseArray(data.html),
      implies: parseArray(data.implies),
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
          <DialogTitle>{isEdit ? t("wappalyzer.editTitle") : t("wappalyzer.addTitle")}</DialogTitle>
          <DialogDescription>
            {isEdit ? t("wappalyzer.editDesc") : t("wappalyzer.addDesc")}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {/* Basic information */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="name">{t("form.appNamePlaceholder").split("：")[0]} *</Label>
              <Input
                id="name"
                placeholder={t("form.appNamePlaceholder")}
                {...register("name", { required: t("form.appNameRequired") })}
              />
              {errors.name && (
                <p className="text-sm text-destructive">{errors.name.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="cats">{t("category")}</Label>
              <Input
                id="cats"
                placeholder={t("form.catsPlaceholder")}
                {...register("cats")}
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="website">{tCommon("website")}</Label>
              <Input
                id="website"
                placeholder="https://example.com"
                {...register("website")}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="cpe">CPE</Label>
              <Input
                id="cpe"
                placeholder="cpe:/a:vendor:product"
                {...register("cpe")}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">{tCommon("description")}</Label>
            <Textarea
              id="description"
              placeholder={t("form.descPlaceholder")}
              rows={2}
              {...register("description")}
            />
          </div>

          {/* Detection rules */}
          <div className="space-y-1">
            <Label className="text-sm font-medium">{t("form.detectionRules")}</Label>
            <p className="text-xs text-muted-foreground">{t("form.detectionRulesHint")}</p>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="cookies">{t("form.cookies")}</Label>
              <Textarea
                id="cookies"
                placeholder='{"name": "pattern"}'
                rows={2}
                className="font-mono text-xs"
                {...register("cookies")}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="headers">{t("form.headers")}</Label>
              <Textarea
                id="headers"
                placeholder='{"X-Powered-By": "pattern"}'
                rows={2}
                className="font-mono text-xs"
                {...register("headers")}
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="scriptSrc">{t("form.scriptUrl")}</Label>
              <Input
                id="scriptSrc"
                placeholder="pattern1, pattern2"
                className="font-mono text-xs"
                {...register("scriptSrc")}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="js">{t("form.jsVariables")}</Label>
              <Input
                id="js"
                placeholder="window.var1, window.var2"
                className="font-mono text-xs"
                {...register("js")}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="meta">{t("form.metaTags")} (JSON)</Label>
            <Textarea
              id="meta"
              placeholder='{"generator": ["pattern"]}'
              rows={2}
              className="font-mono text-xs"
              {...register("meta")}
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="html">{t("form.htmlContent")}</Label>
              <Input
                id="html"
                placeholder="pattern1, pattern2"
                className="font-mono text-xs"
                {...register("html")}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="implies">{t("form.implies")}</Label>
              <Input
                id="implies"
                placeholder="PHP, MySQL"
                {...register("implies")}
              />
            </div>
          </div>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              {tCommon("cancel")}
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? tCommon("saving") : isEdit ? tCommon("update") : tCommon("create")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
