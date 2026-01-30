"use client"

import React, { useState } from "react"
import { toast } from "sonner"
import { useTranslations } from "next-intl"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import {
  Dropzone,
  DropzoneContent,
  DropzoneEmptyState,
} from "@/components/ui/dropzone"
import {
  useImportEholeFingerprints,
  useImportGobyFingerprints,
  useImportWappalyzerFingerprints,
  useImportFingersFingerprints,
  useImportFingerPrintHubFingerprints,
  useImportARLFingerprints,
} from "@/hooks/use-fingerprints"

type FingerprintType = "ehole" | "goby" | "wappalyzer" | "fingers" | "fingerprinthub" | "arl"

interface ImportFingerprintDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  onSuccess?: () => void
  fingerprintType?: FingerprintType
  acceptedFileTypes?: string
}

export function ImportFingerprintDialog({
  open,
  onOpenChange,
  onSuccess,
  fingerprintType = "ehole",
  acceptedFileTypes,
}: ImportFingerprintDialogProps) {
  const [files, setFiles] = useState<File[]>([])
  const t = useTranslations("tools.fingerprints")
  const tCommon = useTranslations("common.actions")
  const tToast = useTranslations("toast")
  
  const eholeImportMutation = useImportEholeFingerprints()
  const gobyImportMutation = useImportGobyFingerprints()
  const wappalyzerImportMutation = useImportWappalyzerFingerprints()
  const fingersImportMutation = useImportFingersFingerprints()
  const fingerprinthubImportMutation = useImportFingerPrintHubFingerprints()
  const arlImportMutation = useImportARLFingerprints()

  const getErrorMessage = (error: unknown): string =>
    error instanceof Error ? error.message : ""

  type FingerprintConfig = {
    title: string
    description: string
    formatHint: string
    validate: (json: unknown) => { valid: boolean; error?: string }
  }

  // Fingerprint type configuration
  const FINGERPRINT_CONFIG: Record<FingerprintType, FingerprintConfig> = {
    ehole: {
      title: t("import.eholeTitle"),
      description: t("import.eholeDesc"),
      formatHint: t.raw("import.eholeFormatHint") as string,
      validate: (json) => {
        if (!json || typeof json !== "object") {
          return { valid: false, error: t("import.eholeInvalidFields") }
        }
        const obj = json as Record<string, unknown>
        if (!obj.fingerprint) {
          return { valid: false, error: t("import.eholeInvalidMissing") }
        }
        if (!Array.isArray(obj.fingerprint)) {
          return { valid: false, error: t("import.eholeInvalidArray") }
        }
        if (obj.fingerprint.length === 0) {
          return { valid: false, error: t("import.emptyData") }
        }
        const first = obj.fingerprint[0] as Record<string, unknown>
        if (!first.cms || !first.keyword) {
          return { valid: false, error: t("import.eholeInvalidFields") }
        }
        return { valid: true }
      },
    },
    goby: {
      title: t("import.gobyTitle"),
      description: t("import.gobyDesc"),
      formatHint: t.raw("import.gobyFormatHint") as string,
      validate: (json) => {
        // Support both array and object formats
        if (Array.isArray(json)) {
          if (json.length === 0) {
            return { valid: false, error: t("import.emptyData") }
          }
          const first = json[0] as Record<string, unknown>
          if (!first.product || !first.rule) {
            return { valid: false, error: t("import.gobyInvalidFields") }
          }
        } else if (typeof json === "object" && json !== null) {
          if (Object.keys(json).length === 0) {
            return { valid: false, error: t("import.emptyData") }
          }
        } else {
          return { valid: false, error: t("import.gobyInvalidFormat") }
        }
        return { valid: true }
      },
    },
    wappalyzer: {
      title: t("import.wappalyzerTitle"),
      description: t("import.wappalyzerDesc"),
      formatHint: t.raw("import.wappalyzerFormatHint") as string,
      validate: (json) => {
        // Support array format
        if (Array.isArray(json)) {
          if (json.length === 0) {
            return { valid: false, error: t("import.emptyData") }
          }
          return { valid: true }
        }
        // Support object format (apps or technologies)
        if (!json || typeof json !== "object") {
          return { valid: false, error: t("import.wappalyzerInvalidFormat") }
        }
        const obj = json as Record<string, unknown>
        const apps = obj.apps || obj.technologies
        if (apps) {
          if (typeof apps !== "object" || Array.isArray(apps)) {
            return { valid: false, error: t("import.wappalyzerInvalidApps") }
          }
          if (Object.keys(apps).length === 0) {
            return { valid: false, error: t("import.emptyData") }
          }
          return { valid: true }
        }
        // Direct object format
        if (typeof json === "object" && json !== null) {
          if (Object.keys(json).length === 0) {
            return { valid: false, error: t("import.emptyData") }
          }
          return { valid: true }
        }
        return { valid: false, error: t("import.wappalyzerInvalidFormat") }
      },
    },
    fingers: {
      title: t("import.fingersTitle"),
      description: t("import.fingersDesc"),
      formatHint: t.raw("import.fingersFormatHint") as string,
      validate: (json) => {
        if (!Array.isArray(json)) {
          return { valid: false, error: t("import.fingersInvalidArray") }
        }
        if (json.length === 0) {
          return { valid: false, error: t("import.emptyData") }
        }
        const first = json[0] as Record<string, unknown>
        if (!first.name || !first.rule) {
          return { valid: false, error: t("import.fingersInvalidFields") }
        }
        return { valid: true }
      },
    },
    fingerprinthub: {
      title: t("import.fingerprinthubTitle"),
      description: t("import.fingerprinthubDesc"),
      formatHint: t.raw("import.fingerprinthubFormatHint") as string,
      validate: (json) => {
        if (!Array.isArray(json)) {
          return { valid: false, error: t("import.fingerprinthubInvalidArray") }
        }
        if (json.length === 0) {
          return { valid: false, error: t("import.emptyData") }
        }
        const first = json[0] as Record<string, unknown>
        if (!first.id || !first.info) {
          return { valid: false, error: t("import.fingerprinthubInvalidFields") }
        }
        return { valid: true }
      },
    },
    arl: {
      title: t("import.arlTitle"),
      description: t("import.arlDesc"),
      formatHint: t.raw("import.arlFormatHint") as string,
      validate: (json) => {
        // ARL supports both YAML and JSON, validation is done on backend
        if (!Array.isArray(json)) {
          return { valid: false, error: t("import.arlInvalidArray") }
        }
        if (json.length === 0) {
          return { valid: false, error: t("import.emptyData") }
        }
        const first = json[0] as Record<string, unknown>
        if (!first.name || !first.rule) {
          return { valid: false, error: t("import.arlInvalidFields") }
        }
        return { valid: true }
      },
    },
  }

  const config = FINGERPRINT_CONFIG[fingerprintType]
  
  const importMutation = {
    ehole: eholeImportMutation,
    goby: gobyImportMutation,
    wappalyzer: wappalyzerImportMutation,
    fingers: fingersImportMutation,
    fingerprinthub: fingerprinthubImportMutation,
    arl: arlImportMutation,
  }[fingerprintType]

  const parseAcceptedFileTypes = (value?: string): string[] => {
    if (!value) return []
    return value
      .split(",")
      .map((item) => item.trim())
      .filter(Boolean)
  }

  const buildAcceptConfig = (extensions: string[]): Record<string, string[]> => {
    const accept: Record<string, string[]> = {}
    const pushExt = (key: string, ext: string) => {
      if (!accept[key]) accept[key] = []
      accept[key].push(ext)
    }

    extensions.forEach((rawExt) => {
      const normalized = rawExt.startsWith(".") ? rawExt : `.${rawExt}`
      if (normalized === ".json") {
        pushExt("application/json", normalized)
        return
      }
      if (normalized === ".yaml" || normalized === ".yml") {
        pushExt("application/x-yaml", normalized)
        pushExt("text/yaml", normalized)
        return
      }
      pushExt("application/octet-stream", normalized)
    })

    return accept
  }

  // Determine accepted file types based on fingerprint type
  const getAcceptConfig = (): Record<string, string[]> => {
    const customExtensions = parseAcceptedFileTypes(acceptedFileTypes)
    if (customExtensions.length > 0) {
      return buildAcceptConfig(customExtensions)
    }
    if (fingerprintType === "arl") {
      return { 
        "application/json": [".json"],
        "application/x-yaml": [".yaml", ".yml"],
        "text/yaml": [".yaml", ".yml"],
      }
    }
    return { "application/json": [".json"] }
  }

  const handleDrop = (acceptedFiles: File[]) => {
    setFiles(acceptedFiles)
  }

  const handleImport = async () => {
    if (files.length === 0) {
      toast.error(tToast("selectFileFirst"))
      return
    }

    const file = files[0]
    const isYamlFile = file.name.endsWith('.yaml') || file.name.endsWith('.yml')

    // Skip frontend validation for YAML files (ARL), let backend handle it
    if (!isYamlFile) {
      // Frontend basic validation for JSON files
      try {
        const text = await file.text()
        let json: unknown

        // Try standard JSON first
        try {
          json = JSON.parse(text)
        } catch {
          // If standard JSON fails, try JSONL format (for goby)
          if (fingerprintType === "goby") {
            const lines = text.trim().split('\n').filter(line => line.trim())
            if (lines.length === 0) {
              toast.error(t("import.emptyData"))
              return
            }
            // Parse each line as JSON
            json = lines.map((line, index) => {
              try {
                return JSON.parse(line)
              } catch {
                throw new Error(`Line ${index + 1}: Invalid JSON`)
              }
            })
          } else {
            throw new Error("Invalid JSON")
          }
        }

        const validation = config.validate(json)
        if (!validation.valid) {
          toast.error(validation.error)
          return
        }
      } catch (e) {
        toast.error(getErrorMessage(e) || tToast("invalidJsonFile"))
        return
      }
    }

    // Validation passed, submit to backend
    try {
      const result = await importMutation.mutateAsync(file)
      toast.success(t("import.importSuccessDetail", { created: result.created, failed: result.failed }))
      setFiles([])
      onOpenChange(false)
      onSuccess?.()
    } catch (error) {
      toast.error(getErrorMessage(error) || tToast("importFailed"))
    }
  }

  const handleClose = (open: boolean) => {
    if (!open) {
      setFiles([])
    }
    onOpenChange(open)
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{config.title}</DialogTitle>
          <DialogDescription>
            {config.description}
          </DialogDescription>
        </DialogHeader>

        <div className="py-4">
          <Dropzone
            src={files}
            onDrop={handleDrop}
            accept={getAcceptConfig()}
            maxFiles={1}
            maxSize={50 * 1024 * 1024}  // 50MB
            onError={(error) => toast.error(error.message)}
          >
            <DropzoneEmptyState />
            <DropzoneContent />
          </Dropzone>

          <p className="text-xs text-muted-foreground mt-3">
            {t("import.supportedFormat")}{" "}
            <code className="bg-muted px-1 rounded">
              {config.formatHint}
            </code>
          </p>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => handleClose(false)}>
            {tCommon("cancel")}
          </Button>
          <Button
            onClick={handleImport}
            disabled={files.length === 0 || importMutation.isPending}
          >
            {importMutation.isPending ? t("import.importing") : tCommon("import")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
